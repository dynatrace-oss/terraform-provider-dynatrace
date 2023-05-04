/**
* @license
* Copyright 2020 Dynatrace LLC
*
* Licensed under the Apache License, Version 2.0 (the "License");
* you may not use this file except in compliance with the License.
* You may obtain a copy of the License at
*
*     http://www.apache.org/licenses/LICENSE-2.0
*
* Unless required by applicable law or agreed to in writing, software
* distributed under the License is distributed on an "AS IS" BASIS,
* WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
* See the License for the specific language governing permissions and
* limitations under the License.
 */

package export

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"sort"
	"sync"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v2/entity"
	entitysettings "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v2/entity/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings/services/cache"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/shutdown"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/provider/version"
	"github.com/spf13/afero"
)

type Environment struct {
	mu           sync.Mutex
	OutputFolder string
	Credentials  *settings.Credentials
	Modules      map[ResourceType]*Module
	Flags        Flags
	ResArgs      map[string][]string
}

func (me *Environment) DataSource(id string) *DataSource {
	for _, module := range me.Modules {
		if dataSource, found := module.DataSources[id]; found {
			return dataSource
		}
	}
	service := cache.Read(entity.Service(me.Credentials))
	var entity entitysettings.Entity
	if err := service.Get(id, &entity); err == nil {
		return &DataSource{ID: *entity.EntityId, Name: *entity.DisplayName, Type: *entity.Type}
	}
	return nil
}

func (me *Environment) Export() (err error) {
	if err = me.InitialDownload(); err != nil {
		return err
	}

	if err = me.PostProcess(); err != nil {
		return err
	}

	if err = me.Finish(); err != nil {
		return err
	}
	return nil
}

func (me *Environment) InitialDownload() error {
	parallel := (os.Getenv("DYNATRACE_PARALLEL") != "false")
	resourceTypes := []string{}
	for resourceType := range me.ResArgs {
		resourceTypes = append(resourceTypes, string(resourceType))
	}
	sort.Strings(resourceTypes)

	if parallel {
		var wg sync.WaitGroup
		wg.Add(len(resourceTypes))
		for _, sResourceType := range resourceTypes {
			go func(sResourceType string) error {
				defer wg.Done()
				if shutdown.System.Stopped() {
					return nil
				}

				keys := me.ResArgs[sResourceType]
				module := me.Module(ResourceType(sResourceType))
				if err := module.Download(false, keys...); err != nil {
					return err
				}
				return nil
			}(sResourceType)
		}
		wg.Wait()
	} else {
		for _, sResourceType := range resourceTypes {
			if shutdown.System.Stopped() {
				return nil
			}

			keys := me.ResArgs[sResourceType]
			module := me.Module(ResourceType(sResourceType))
			if err := module.Download(true, keys...); err != nil {
				return err
			}
		}
	}

	return nil
}

func (me *Environment) PostProcess() error {
	fmt.Println("Post-Processing Resources ...")
	resources := me.GetNonPostProcessedResources()
	for len(resources) > 0 {
		if shutdown.System.Stopped() {
			return nil
		}
		m := map[ResourceType][]*Resource{}
		for _, resource := range resources {
			var reslist []*Resource
			if rl, found := m[resource.Type]; !found {
				reslist = []*Resource{}
			} else {
				reslist = rl
			}
			reslist = append(reslist, resource)
			m[resource.Type] = reslist
		}
		const ClearLine = "\033[2K"
		for k, reslist := range m {
			fmt.Printf("- %s (0 of %d)", k, len(reslist))
			for idx, resource := range reslist {
				if shutdown.System.Stopped() {
					return nil
				}
				if err := resource.PostProcess(); err != nil {
					return err
				}
				fmt.Print(ClearLine)
				fmt.Print("\r")
				fmt.Printf("- %s (%d of %d)", k, idx+1, len(reslist))

			}
			fmt.Print(ClearLine)
			fmt.Print("\r")
			fmt.Printf("- %s\n", k)
		}

		resources = me.GetNonPostProcessedResources()
	}

	for _, resource := range me.GetChildResources() {

		var parentBytes []byte
		var childBytes []byte
		var err error
		if parentBytes, err = resource.Parent.ReadFile(); err != nil {
			return err
		}
		if childBytes, err = resource.ReadFile(); err != nil {
			return err
		}
		var parentFile *os.File
		if parentFile, err = resource.Parent.CreateFile(); err != nil {
			return err
		}
		defer parentFile.Close()
		parentFile.Write(parentBytes)
		parentFile.Write([]byte("\n\n"))
		parentFile.Write(childBytes)
	}
	return nil
}

func (me *Environment) Finish() (err error) {
	if shutdown.System.Stopped() {
		return nil
	}

	if err = me.WriteResourceFiles(); err != nil {
		return err
	}
	if err = me.WriteDataSourceFiles(); err != nil {
		return err
	}
	if err = me.WriteMainFile(); err != nil {
		return err
	}
	if err = me.WriteVariablesFiles(); err != nil {
		return err
	}
	if err = me.WriteProviderFiles(); err != nil {
		return err
	}
	// for _, module := range me.Modules {
	// 	fmt.Println(module.Type, len(module.GetPostProcessedResources()))
	// }

	if err = me.RemoveNonReferencedModules(); err != nil {
		return err
	}
	return nil
}

func (me *Environment) Module(resType ResourceType) *Module {
	me.mu.Lock()
	defer me.mu.Unlock()
	if stored, found := me.Modules[resType]; found {
		return stored
	}
	module := &Module{
		Type:        resType,
		Resources:   map[string]*Resource{},
		DataSources: map[string]*DataSource{},
		namer:       NewUniqueNamer().Replace(ResourceName),
		Status:      ModuleStati.Untouched,
		Environment: me,
	}
	me.Modules[resType] = module
	return module
}

func (me *Environment) CreateFile(name string) (*os.File, error) {
	os.MkdirAll(me.GetFolder(), os.ModePerm)
	return os.Create(path.Join(me.GetFolder(), name))
}

func (me *Environment) GetFolder() string {
	return me.OutputFolder
}

func (me *Environment) GetAttentionFolder() string {
	return path.Join(me.OutputFolder, ".requires_attention")
}

func (me *Environment) GetFlawedFolder() string {
	return path.Join(me.OutputFolder, ".flawed")
}

func (me *Environment) RefersTo(resource *Resource) bool {
	if resource == nil {
		return false
	}
	for _, module := range me.Modules {
		if module.RefersTo(resource) {
			return true
		}
	}
	return false
}

func (me *Environment) GetNonPostProcessedResources() []*Resource {
	resources := []*Resource{}
	for _, module := range me.Modules {
		resources = append(resources, module.GetNonPostProcessedResources()...)
	}
	return resources
}

func (me *Environment) GetChildResources() []*Resource {
	resources := []*Resource{}
	for _, module := range me.Modules {
		if module.Descriptor.Parent != nil {
			resources = append(resources, module.GetChildResources()...)
		}
	}
	return resources
}

func (me *Environment) WriteDataSourceFiles() (err error) {
	fmt.Println("Writing ___datasources___.tf")

	if me.Flags.Flat {
		dataSources := map[string]*DataSource{}
		for _, module := range me.Modules {
			for k, v := range module.DataSources {
				dataSources[k] = v
			}
		}
		var datasourcesFile *os.File
		if datasourcesFile, err = me.CreateFile("___datasources___.tf"); err != nil {
			return err
		}
		defer func() {
			datasourcesFile.Close()
			format(datasourcesFile.Name(), true)
		}()

		for dataSourceID, dataSource := range dataSources {
			if _, err = datasourcesFile.WriteString(fmt.Sprintf(`data "dynatrace_entity" "%s" {
				type = "%s"
				name = "%s"				
			}
`, dataSourceID, dataSource.Type, dataSource.Name)); err != nil {
				return err
			}
		}

		return nil
	}
	for _, module := range me.Modules {
		if err = module.WriteDataSourcesFile(); err != nil {
			return err
		}
	}
	return nil
}

func (me *Environment) WriteResourceFiles() (err error) {
	if me.Flags.Flat {
		return nil
	}
	fmt.Println("Writing ___resources___.tf")
	for _, module := range me.Modules {
		if err = module.WriteResourcesFile(); err != nil {
			return err
		}
	}
	return nil
}

func (me *Environment) RemoveNonReferencedModules() (err error) {
	for _, module := range me.Modules {
		if module.IsReferencedAsDataSource() || module.Descriptor.Parent != nil {
			if err = module.PurgeFolder(); err != nil {
				return err
			}
		}
		if len(module.GetPostProcessedResources()) == 0 {
			if err = module.PurgeFolder(); err != nil {
				return err
			}
		}
	}
	return nil
}

func (me *Environment) WriteProviderFiles() (err error) {
	fmt.Println("Writing ___providers___.tf")

	var outputFile *os.File
	if outputFile, err = me.CreateFile("___providers___.tf"); err != nil {
		return err
	}
	defer func() {
		outputFile.Close()
		format(outputFile.Name(), true)
	}()
	providerSource := "dynatrace-oss/dynatrace"
	providerVersion := version.Current
	if value := os.Getenv(DYNATRACE_PROVIDER_SOURCE); len(value) != 0 {
		providerSource = value
	}
	if value := os.Getenv(DYNATRACE_PROVIDER_VERSION); len(value) != 0 {
		providerVersion = value
	}

	if _, err = outputFile.WriteString(fmt.Sprintf(`terraform {
	required_providers {
		dynatrace = {
		source = "%s"
		version = "%s"
		}
	}
	}

	provider "dynatrace" {
	}	  
`, providerSource, providerVersion)); err != nil {
		return err
	}
	if me.Flags.Flat {
		return nil
	}
	for _, module := range me.Modules {
		if err = module.WriteProviderFile(); err != nil {
			return err
		}
	}
	return nil
}

func (me *Environment) WriteVariablesFiles() (err error) {
	fmt.Println("Writing ___variables___.tf")
	for _, module := range me.Modules {
		if err = module.WriteVariablesFile(); err != nil {
			return err
		}
	}
	return nil
}

func (me *Environment) GetResourceTypesWithDownloads() []ResourceType {
	resourceTypesWithDownloads := map[ResourceType]ResourceType{}
	for _, module := range me.Modules {
		for _, resource := range module.Resources {
			if resource.Status == ResourceStati.PostProcessed {
				resourceTypesWithDownloads[resource.Type] = resource.Type
			}
		}
	}
	result := []ResourceType{}
	for resourceType := range resourceTypesWithDownloads {
		result = append(result, resourceType)
	}
	return result
}

func (me *Environment) WriteMainFile() error {
	if me.Flags.Flat {
		return nil
	}
	fmt.Println("Writing main.tf")

	var err error
	var mainFile *os.File
	if mainFile, err = os.Create(path.Join(me.OutputFolder, "main.tf")); err != nil {
		return err
	}
	defer func() {
		mainFile.Close()
		format(mainFile.Name(), true)
	}()
	resourceTypes := me.GetResourceTypesWithDownloads()
	sResourceTypes := []string{}
	for _, resourceType := range resourceTypes {
		sResourceTypes = append(sResourceTypes, string(resourceType))
	}
	sort.Strings(sResourceTypes)
	for _, sResourceType := range sResourceTypes {
		resourceType := ResourceType(sResourceType)
		if me.Module(resourceType).IsReferencedAsDataSource() {
			continue
		}
		if me.Module(resourceType).Descriptor.Parent != nil {
			continue
		}
		if len(me.Module(resourceType).GetPostProcessedResources()) == 0 {
			continue
		}
		mainFile.WriteString(fmt.Sprintf("module \"%s\" {\n", resourceType.Trim()))
		module := me.Module(resourceType)
		mainFile.WriteString(fmt.Sprintf("  source = \"./%s\"\n", module.GetFolder(true)))
		referencedResourceTypes := module.GetReferencedResourceTypes()
		if len(referencedResourceTypes) > 0 {
			for _, referencedResourceType := range referencedResourceTypes {
				if me.Module(referencedResourceType).IsReferencedAsDataSource() {
					continue
				}
				if referencedResourceType == resourceType {
					continue
				}
				mainFile.WriteString(fmt.Sprintf("  %s = module.%s.resources\n", referencedResourceType, referencedResourceType.Trim()))

			}
		}
		mainFile.WriteString("}\n\n")
	}
	return nil
}

func (me *Environment) ExecuteImport() error {
	if me.Flags.ImportState {
		return me.executeImportV1()
	}
	if me.Flags.ImportStateV2 {
		return me.executeImportV2()
	}

	return nil
}

func (me *Environment) executeImportV1() error {
	for _, module := range me.Modules {
		if shutdown.System.Stopped() {
			return errors.New("Import was stopped")
		}
		if err := module.ExecuteImportV1(); err != nil {
			return err
		}
	}
	return nil
}

func (me *Environment) executeImportV2() error {
	itemCount := len(me.Modules)
	fmt.Printf("Importing %d Modules", itemCount)
	channel := make(chan *Module, itemCount)
	mutex := sync.Mutex{}
	waitGroup := sync.WaitGroup{}
	maxThreads := 10
	if maxThreads > itemCount {
		maxThreads = itemCount
	}
	waitGroup.Add(maxThreads)
	errs := []error{}
	fs := afero.NewOsFs()
	var newStateObject interface{}

	processModule := func(module *Module) {
		stateObject, err := module.ExecuteImportV2(fs, "")
		if err != nil {
			mutex.Lock()
			errs = append(errs, err)
			mutex.Unlock()
			return
		}
		mutex.Lock()
		newStateObject = updateState(newStateObject, stateObject)
		mutex.Unlock()
	}

	for i := 0; i < maxThreads; i++ {

		go func() {

			for {
				module, ok := <-channel
				if shutdown.System.Stopped() {
					ok = false
				}
				if !ok {
					waitGroup.Done()
					return
				}
				processModule(module)
			}
		}()

	}

	for _, module := range me.Modules {
		channel <- module
	}

	close(channel)
	waitGroup.Wait()
	if shutdown.System.Stopped() {
		return errors.New("import was stopped")
	}

	if len(errs) >= 1 {
		return fmt.Errorf("Error during state import", errs)
	}

	bytes, err := json.MarshalIndent(newStateObject, "", "  ")
	if err != nil {
		return err
	}
	filename := fmt.Sprint(filepath.Join(me.OutputFolder, "terraform.tfstate"))
	err = afero.WriteFile(fs, filename, bytes, 0664)
	if err != nil {
		return err
	}

	return nil
}

func updateState(newStateObject interface{}, stateObject interface{}) interface{} {
	if stateObject == nil {
		return newStateObject
	}

	if newStateObject == nil {
		newStateObject = stateObject
	} else {
		serialValue := stateObject.(map[string]interface{})["serial"].(float64)
		newSerialValue := newStateObject.(map[string]interface{})["serial"].(float64)
		newStateObject.(map[string]interface{})["serial"] = (serialValue + newSerialValue)

		newStateObject.(map[string]interface{})["resources"] = append(
			newStateObject.(map[string]interface{})["resources"].([]interface{}),
			stateObject.(map[string]interface{})["resources"].([]interface{})...,
		)
	}
	return newStateObject
}
