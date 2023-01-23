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
	"fmt"
	"os"
	"path"
	"sort"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v2/entity"
	entitysettings "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v2/entity/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/shutdown"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/provider/version"
)

type Environment struct {
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
	service := entity.Service(me.Credentials)
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
	resourceTypes := []string{}
	for resourceType := range me.ResArgs {
		resourceTypes = append(resourceTypes, string(resourceType))
	}
	sort.Strings(resourceTypes)

	for _, sResourceType := range resourceTypes {
		if shutdown.System.Stopped() {
			return nil
		}

		keys := me.ResArgs[sResourceType]
		module := me.Module(ResourceType(sResourceType))
		if err := module.Download(keys...); err != nil {
			return err
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
		for k, reslist := range m {
			fmt.Println("- " + k)
			for _, resource := range reslist {
				if shutdown.System.Stopped() {
					return nil
				}
				if err := resource.PostProcess(); err != nil {
					return err
				}
			}
		}

		resources = me.GetNonPostProcessedResources()
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
	return nil
}

func (me *Environment) Module(resType ResourceType) *Module {
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

func (me *Environment) WriteDataSourceFiles() (err error) {
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
	for _, module := range me.Modules {
		if err = module.WriteResourcesFile(); err != nil {
			return err
		}
	}
	return nil
}

func (me *Environment) WriteProviderFiles() (err error) {
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
		mainFile.WriteString(fmt.Sprintf("module \"%s\" {\n", resourceType.Trim()))
		module := me.Module(resourceType)
		mainFile.WriteString(fmt.Sprintf("  source = \"./%s\"\n", module.GetFolder(true)))
		referencedResourceTypes := module.GetReferencedResourceTypes()
		if len(referencedResourceTypes) > 0 {
			for _, referencedResourceType := range referencedResourceTypes {
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
	if !me.Flags.ImportState {
		return nil
	}
	for _, module := range me.Modules {
		if err := module.ExecuteImport(); err != nil {
			return err
		}
	}
	return nil
}
