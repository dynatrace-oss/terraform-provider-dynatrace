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
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"sort"
	"strings"
	"sync"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/address"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v2/entity"
	entitysettings "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v2/entity/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings/services/cache"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/shutdown"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/provider/version"
	"github.com/google/uuid"
	"github.com/spf13/afero"
)

var NO_REFRESH_ON_IMPORT = os.Getenv("DYNATRACE_NO_REFRESH_ON_IMPORT") == "true"
var QUICK_INIT = os.Getenv("DYNATRACE_QUICK_INIT") == "true"
var ULTRA_PARALLEL = os.Getenv("DYNATRACE_ULTRA_PARALLEL") == "true"

const ENV_VAR_CUSTOM_PROVIDER_LOCATION = "DYNATRACE_CUSTOM_PROVIDER_LOCATION"

type Environment struct {
	mu                    sync.Mutex
	OutputFolder          string
	Credentials           *settings.Credentials
	Modules               map[ResourceType]*Module
	Flags                 Flags
	ResArgs               map[string][]string
	ChildResourceOverride bool
	PrevStateMapCommon    *StateMap
	PrevNamesByModule     map[string][]string
	ImportStateMap        *StateMap
	ChildParentGroups     map[ResourceType]ResourceType
	IsParentMap           map[ResourceType]bool
	HasDependenciesTo     map[ResourceType]bool
}

func (me *Environment) TenantID() string {
	var tenant string
	if strings.Contains(me.Credentials.URL, "/e/") {
		idx := strings.Index(me.Credentials.URL, "/e/")
		tenant = strings.TrimSuffix(strings.TrimPrefix(me.Credentials.URL[idx:], "/e/"), "/")
	} else if strings.HasPrefix(me.Credentials.URL, "http://") {
		tenant = strings.TrimPrefix(me.Credentials.URL, "http://")
		if idx := strings.Index(tenant, "."); idx != -1 {
			tenant = tenant[:idx]
		}
	} else if strings.HasPrefix(me.Credentials.URL, "https://") {
		tenant = strings.TrimPrefix(me.Credentials.URL, "https://")
		if idx := strings.Index(tenant, "."); idx != -1 {
			tenant = tenant[:idx]
		}
	}
	return tenant
}

func (me *Environment) DataSource(id string, kind DataSourceKind, excepts ...ResourceType) *DataSource {
	switch kind {
	case DataSourceKindTenant:
		return &DataSource{ID: id, Name: id, Type: "tenant", Kind: DataSourceKindTenant}
	case DataSourceKindEntity:
		return me.FetchEntity(id)
	}
	for _, module := range me.Modules {
		skipModule := false
		for _, except := range excepts {
			if module.Type == except {
				skipModule = true
			}
		}
		if skipModule {
			continue
		}
		if dataSource := module.DataSource(id, kind, excepts...); dataSource != nil {
			return dataSource
		}
	}
	return nil
}

func (me *Environment) FetchEntity(id string) *DataSource {
	service := cache.Read(entity.DataSourceService(me.Credentials))
	var entity entitysettings.Entity
	if err := service.Get(id, &entity); err == nil {
		return &DataSource{ID: *entity.EntityId, Name: *entity.DisplayName, Type: *entity.Type, Kind: DataSourceKindEntity}
	}
	return nil
}

func (me *Environment) Export() (err error) {

	if err = me.PreProcess(); err != nil {
		return err
	}

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

func (me *Environment) PreProcess() error {
	me.ProcessChildParentGroups()
	me.ProcessHasDependenciesTo()
	err := me.LoadImportState()
	if err != nil {
		return err
	}
	err = me.ProcessPrevState()
	if err != nil {
		return err
	}

	return nil
}

func (me *Environment) ProcessChildParentGroups() {
	me.ChildParentGroups = map[ResourceType]ResourceType{}
	me.IsParentMap = map[ResourceType]bool{}

	for resType, resource := range AllResources {
		if resource.Parent != nil {
			me.ChildParentGroups[resType] = *resource.Parent
			me.ChildParentGroups[*resource.Parent] = *resource.Parent
			me.IsParentMap[*resource.Parent] = true
		}
	}
}

func (me *Environment) ProcessHasDependenciesTo() {
	me.HasDependenciesTo = map[ResourceType]bool{}

	for _, resource := range AllResources {
		for _, dep := range resource.Dependencies {
			resSource := dep.ResourceType()
			if resSource == "" {
				continue
			}

			me.HasDependenciesTo[resSource] = true

		}
	}
}

func (me *Environment) LoadImportState() error {
	if IMPORT_STATE_PATH == "" {
		return nil
	}

	state, err := LoadStateFile(IMPORT_STATE_PATH)
	if err != nil {
		return err
	}

	stateMap := BuildStateMap(state)

	me.ImportStateMap = stateMap

	return nil
}

func (me *Environment) ProcessPrevState() error {
	if PREV_STATE_ON {
		// pass
	} else {
		return nil
	}

	stateThis, err := LoadStateThis()
	if err != nil {
		return err
	}
	stateLinked, err := LoadStateLinked()
	if err != nil {
		return err
	}

	stateMapThis := BuildStateMap(stateThis)
	stateMapLinked := BuildStateMap(stateLinked)

	me.PrevStateMapCommon, me.PrevNamesByModule = stateMapThis.ExtractCommonStates(stateMapLinked)

	return nil
}

func (me *Environment) InitialDownload() error {
	parallel := (os.Getenv("DYNATRACE_PARALLEL") != "false")
	resourceTypes := []string{}
	for resourceType := range me.ResArgs {
		resourceTypes = append(resourceTypes, string(resourceType))
		me.Module(ResourceType(resourceType)).blockPrevNames()
	}
	sort.Strings(resourceTypes)

	if parallel {
		var wg sync.WaitGroup
		itemCount := len(resourceTypes)
		channel := make(chan string, itemCount)
		maxThreads := 10
		if maxThreads > itemCount {
			maxThreads = itemCount
		}
		wg.Add(maxThreads)

		processItem := func(sResourceType string) error {
			keys := me.ResArgs[sResourceType]
			module := me.Module(ResourceType(sResourceType))
			if err := module.Download(parallel, keys...); err != nil {
				return err
			}
			return nil
		}

		for i := 0; i < maxThreads; i++ {

			go func() error {

				for {
					sResourceTypeLoop, ok := <-channel
					if !ok {
						wg.Done()
						return nil
					}
					if shutdown.System.Stopped() {
						wg.Done()
						return nil
					}

					err := processItem(sResourceTypeLoop)

					if err != nil {
						wg.Done()
						return err
					}
				}
			}()

		}

		for _, sResourceType := range resourceTypes {
			channel <- sResourceType
		}

		close(channel)
		wg.Wait()
	} else {
		for _, sResourceType := range resourceTypes {
			if shutdown.System.Stopped() {
				return nil
			}

			keys := me.ResArgs[sResourceType]
			module := me.Module(ResourceType(sResourceType))
			if err := module.Download(parallel, keys...); err != nil {
				return err
			}
		}
	}

	return nil
}

func (me *Environment) PostProcess() error {
	fmt.Println("Post-Processing Resources ...")
	parallel := (os.Getenv("DYNATRACE_PARALLEL") != "false")
	resources := me.GetNonPostProcessedResources()

	if parallel {
		for len(resources) > 0 {

			m := getResMap(resources)
			for _, reslist := range m {

				var wg sync.WaitGroup
				itemCount := len(reslist)
				channel := make(chan *Resource, itemCount)
				maxThreads := 50
				if maxThreads > itemCount {
					maxThreads = itemCount
				}
				wg.Add(maxThreads)

				processItem := func(resource *Resource) error {
					if shutdown.System.Stopped() {
						return nil
					}
					if err := resource.PostProcess(); err != nil {
						return err
					}
					fmt.Print("\r")
					fmt.Printf("- [POSTPROCESS] %s - %s", resource.Type, resource.UniqueName)

					return nil
				}

				for i := 0; i < maxThreads; i++ {

					go func() error {

						for {
							res, ok := <-channel
							if !ok {
								wg.Done()
								return nil
							}
							if shutdown.System.Stopped() {
								wg.Done()
								return nil
							}

							err := processItem(res)

							if err != nil {
								wg.Done()
								return err
							}
						}
					}()

				}

				for _, res := range reslist {
					channel <- res
				}

				close(channel)
				wg.Wait()
			}

			resources = me.GetNonPostProcessedResources()
		}

	} else {
		for len(resources) > 0 {
			if shutdown.System.Stopped() {
				return nil
			}
			m := getResMap(resources)
			const ClearLine = "\033[2K"
			for k, reslist := range m {
				fmt.Printf("- [POSTPROCESS] %s (0 of %d)", k, len(reslist))
				for idx, resource := range reslist {
					if shutdown.System.Stopped() {
						return nil
					}
					if err := resource.PostProcess(); err != nil {
						return err
					}
					fmt.Print(ClearLine)
					fmt.Print("\r")
					fmt.Printf("- [POSTPROCESS] %s (%d of %d)", k, idx+1, len(reslist))

				}
				fmt.Print(ClearLine)
				fmt.Print("\r")
				fmt.Printf("- [POSTPROCESS] %s\n", k)
			}

			resources = me.GetNonPostProcessedResources()
		}
	}

	fmt.Println("Post-Processing Resources - Group child configs with parent configs ...")
	for _, resource := range me.GetChildResources() {
		if resource.GetParent().Status == ResourceStati.Erronous {
			continue
		}
		var parentBytes []byte
		var childBytes []byte
		var err error
		if parentBytes, err = resource.GetParent().ReadFile(); err == nil {
			if childBytes, err = resource.ReadFile(); err == nil {
				resource.GetParent().Module.saveChildModule(resource.Module)
				var parentFile *os.File
				if parentFile, err = resource.GetParent().CreateFile(); err == nil {
					defer parentFile.Close()
					parentFile.Write(parentBytes)
					parentFile.Write([]byte("\n\n"))
					parentFile.Write(childBytes)
				}
			}
		}
	}
	return nil
}

func getResMap(resources []*Resource) map[ResourceType][]*Resource {
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
	return m
}

func (me *Environment) Finish() (err error) {
	fmt.Println("Finishing touches ...")
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
		Type:                 resType,
		Resources:            map[string]*Resource{},
		DataSources:          map[string]*DataSource{},
		namer:                NewUniqueNamer().Replace(ResourceName),
		Status:               ModuleStati.Untouched,
		Environment:          me,
		ChildParentIDNameMap: map[string]string{},
		ModuleMutex:          new(sync.Mutex),
		DataSourceLock:       new(sync.Mutex),
		ChildModules:         map[ResourceType]*Module{},
	}

	if resType == ResourceTypes.JSONDashboardBase {
		if module.Descriptor == nil {
			descriptor := AllResources[resType]
			module.Descriptor = &descriptor
		}
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

func (me *Environment) RefersTo(resource *Resource, parentType ResourceType) bool {
	if resource == nil {
		return false
	}
	for _, module := range me.Modules {
		if module.RefersTo(resource, parentType) {
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
	if me.ChildResourceOverride {
		return resources
	}
	for _, module := range me.Modules {
		if module != nil && module.Descriptor != nil && module.Descriptor.Parent != nil {
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
			module.GetDataSources(dataSources)
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
			if dataSource.ID == "tenant" {
				if _, err = datasourcesFile.WriteString(`data "dynatrace_tenant" "tenant" {
				}
	`); err != nil {
					return err
				}
			} else {
				if _, err = datasourcesFile.WriteString(fmt.Sprintf(`data "dynatrace_entity" "%s" {
					type = "%s"
					name = "%s"				
				}
	`, dataSourceID, dataSource.Type, dataSource.Name)); err != nil {
					return err
				}
			}
		}
		dsm := map[string]string{}
		for _, module := range me.Modules {
			mdsm, err := module.ProvideDataSources()
			if err != nil {
				return err
			}
			for k, v := range mdsm {
				dsm[k] = v
			}
		}
		for _, ds := range dsm {
			datasourcesFile.Write([]byte("\n" + ds))
		}

		return nil
	}
	parallel := (os.Getenv("DYNATRACE_PARALLEL") != "false")
	if parallel {
		var wg sync.WaitGroup
		wg.Add(len(me.Modules))
		for _, module := range me.Modules {
			go func(module *Module) error {
				defer wg.Done()
				if shutdown.System.Stopped() {
					return nil
				}
				if err = module.WriteDataSourcesFile(false); err != nil {
					return err
				}
				return nil
			}(module)
		}
		wg.Wait()
	} else {
		for _, module := range me.Modules {
			if err = module.WriteDataSourcesFile(true); err != nil {
				return err
			}
		}
	}
	return nil
}

func (me *Environment) WriteResourceFiles() (err error) {
	if me.Flags.Flat {
		return nil
	}
	fmt.Println("Writing ___resources___.tf")
	parallel := (os.Getenv("DYNATRACE_PARALLEL") != "false")
	if parallel {
		var wg sync.WaitGroup
		wg.Add(len(me.Modules))
		for _, module := range me.Modules {
			go func(module *Module) error {
				defer wg.Done()
				if shutdown.System.Stopped() {
					return nil
				}
				if err = module.WriteResourcesFile(); err != nil {
					return err
				}
				return nil
			}(module)
		}
		wg.Wait()
	} else {
		for _, module := range me.Modules {
			if err = module.WriteResourcesFile(); err != nil {
				return err
			}
		}
	}
	return nil
}

func (me *Environment) RemoveNonReferencedModules() (err error) {
	fmt.Println("Remove Non-Referenced Modules ...")
	m := map[ResourceType]*Module{}
	for k, module := range me.Modules {
		m[k] = module
	}
	for k, module := range m {
		if module.IsReferencedAsDataSource() {
			module.PurgeFolder()
			delete(me.Modules, k)
		} else if !module.Environment.ChildResourceOverride && module.Descriptor.Parent != nil {
			if me.Flags.FollowReferences {
				module.PurgeFolder()
			}
		} else if len(module.GetPostProcessedResources()) == 0 {
			module.PurgeFolder()
			delete(me.Modules, k)
		}
	}
	return nil
}

func (me *Environment) WriteProviderFiles() (err error) {

	if QUICK_INIT {
		// pass
	} else {
		err = me.WriteMainProviderFile()
		if err != nil {
			return err
		}
	}

	if me.Flags.Flat {
		return nil
	}

	fmt.Println("Writing modules ___providers___.tf")
	parallel := (os.Getenv("DYNATRACE_PARALLEL") != "false")
	if parallel {
		var wg sync.WaitGroup
		wg.Add(len(me.Modules))
		for _, module := range me.Modules {
			go func(module *Module) error {
				defer wg.Done()
				if shutdown.System.Stopped() {
					return nil
				}
				if err = module.WriteProviderFile(false); err != nil {
					return err
				}
				return nil
			}(module)
		}
		wg.Wait()
	} else {
		for _, module := range me.Modules {
			if err = module.WriteProviderFile(true); err != nil {
				return err
			}
		}
	}
	return nil
}

func (me *Environment) WriteMainProviderFile() error {
	fmt.Println("Writing main ___providers___.tf")
	var outputFile *os.File
	var err error = nil
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

	return nil
}

func (me *Environment) WriteVariablesFiles() (err error) {
	fmt.Println("Writing ___variables___.tf")
	parallel := (os.Getenv("DYNATRACE_PARALLEL") != "false")
	if parallel {
		var wg sync.WaitGroup

		for _, module := range me.Modules {
			if module.Descriptor == nil {
				continue
			}
			wg.Add(1)
			go func(module *Module) error {
				defer wg.Done()
				if shutdown.System.Stopped() {
					return nil
				}
				if err = module.WriteVariablesFile(false); err != nil {
					return err
				}
				return nil
			}(module)
		}
		wg.Wait()
	} else {
		for _, module := range me.Modules {
			if err = module.WriteVariablesFile(true); err != nil {
				return err
			}
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

	os.MkdirAll(me.OutputFolder, os.ModePerm)

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
		if !me.ChildResourceOverride && me.Module(resourceType).Descriptor.Parent != nil {
			continue
		}
		if len(me.Module(resourceType).GetPostProcessedResources()) == 0 {
			continue
		}

		module := me.Module(resourceType)
		me.writeOpeningMainSection(mainFile, resourceType.Trim(), module.GetFolder(true))

		if ATOMIC_DEPENDENCIES {

			uniqueNameExists := map[string]bool{}
			referencedResources := module.GetResourceReferences()
			if len(referencedResources) > 0 {
				for _, referencedResource := range referencedResources {
					if me.Module(referencedResource.Type).IsReferencedAsDataSource() {
						continue
					}
					if referencedResource.Type == resourceType || (referencedResource.XParent != nil && referencedResource.XParent.Type == resourceType) {
						continue
					}
					typeAndUniqueName := referencedResource.Type.Trim() + referencedResource.UniqueName
					if uniqueNameExists[typeAndUniqueName] {
						continue
					}
					uniqueNameExists[typeAndUniqueName] = true
					mainFile.WriteString(fmt.Sprintf("  %s_%s = module.%s.resources_%s\n", referencedResource.Type, referencedResource.UniqueName, referencedResource.Type.Trim(), referencedResource.UniqueName))
				}
			}
		} else {
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
		}
		writeClosingMainSection(mainFile)

		if module.SplitPathModuleNameMap != nil {
			for _, splitName := range module.SplitPathModuleNameMap {
				me.writeOpeningMainSection(mainFile, splitName, fmt.Sprintf("./modules/%s", splitName))
				writeClosingMainSection(mainFile)
			}
		}
	}
	return nil
}

func (me *Environment) writeOpeningMainSection(mainFile *os.File, trimmedResourceType string, resourceFolder string) {
	mainFile.WriteString(fmt.Sprintf("module \"%s\" {\n", trimmedResourceType))
	mainFile.WriteString(fmt.Sprintf("  source = \"./%s\"\n", resourceFolder))
}

func writeClosingMainSection(mainFile *os.File) {
	mainFile.WriteString("}\n\n")
}

func (me *Environment) ExecuteImport() error {
	if me.Flags.ImportStateV2 {
		fmt.Println("Importing Resources into Terraform State ...")
		err := me.executeImportV2()
		fmt.Println("Imported Resources into Terraform State ...")
		return err
	}

	return nil
}

type state struct {
	Version           int         `json:"version"`
	Terraform_version string      `json:"terraform_version"`
	Serial            int         `json:"serial"`
	Lineage           string      `json:"lineage"`
	Outputs           interface{} `json:"outputs"`
	Resources         resources   `json:"resources"`
	CheckResults      interface{} `json:"check_results"`
}

func (me *Environment) executeImportV2() error {
	fs := afero.NewOsFs()

	state := state{
		Version:           4,
		Terraform_version: "1.4.5",
		Serial:            0,
		Lineage:           uuid.NewString(),
		Outputs:           nil,
		Resources:         resources{},
		CheckResults:      nil,
	}

	for _, module := range me.Modules {
		resList, err := module.ExecuteImportV2(fs)
		if err != nil {
			return err
		}
		state.Resources = append(state.Resources, resList...)
	}

	me.importPrevResources(&state)

	bytes, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return err
	}
	filename := fmt.Sprint(filepath.Join(me.OutputFolder, "terraform.tfstate"))
	err = afero.WriteFile(fs, filename, bytes, 0664)
	if err != nil {
		return err
	}

	if NO_REFRESH_ON_IMPORT {
		fmt.Println("NO_REFRESH_ON_IMPORT set to true")
	} else {
		me.executeTF(genPlanCmd())
		me.executeTF(genApplyCmd())
	}

	return nil
}

func (me *Environment) importPrevResources(state *state) {
	if PREV_STATE_ON {
		for _, statePrev := range me.PrevStateMapCommon.resources {
			if statePrev.Used {
				continue
			}

			module, found := me.Modules[ResourceType(statePrev.Resource.Type)]
			if found {
				// pass
			} else {
				fmt.Println("ERROR: [importPrevResources] Could not find Module: ", statePrev.Resource.Type)
				continue
			}

			isWritten := module.namer.SetNameWritten(statePrev.Resource.Name)
			if isWritten {
				continue
			}
			fmt.Println("[importPrevResources] Post-Writing resource: ", statePrev.Resource.Name, statePrev.Resource.Type)
			state.Resources = append(state.Resources, statePrev.Resource)
		}
	}
}

func genPlanCmd() *exec.Cmd {
	fmt.Println("Refreshing State - Plan ...")
	exePath, _ := exec.LookPath("terraform")

	return exec.Command(
		exePath,
		"plan",
		"-lock=false",
		"-parallelism=50",
		"-refresh-only",
		"-no-color",
		"-out=terraform.plan",
	)
}

func genApplyCmd() *exec.Cmd {
	fmt.Println("Refreshing State - Apply ...")
	exePath, _ := exec.LookPath("terraform")

	return exec.Command(
		exePath,
		"apply",
		"-lock=false",
		"-parallelism=50",
		"-refresh-only",
		"-no-color",
		"-auto-approve",
		"terraform.plan",
	)
}

func (me *Environment) executeTF(cmd *exec.Cmd) (err error) {

	var outb, errb bytes.Buffer
	cmd.Stdout = &outb
	cmd.Stderr = &errb
	cmd.Dir = me.OutputFolder
	var cacheFolder string
	if cacheFolder, err = filepath.Abs(cache.GetCacheFolder()); err != nil {
		return err
	}
	cmd.Env = []string{
		// "TF_LOG_PROVIDER=INFO",
		"DYNATRACE_ENV_URL=" + me.Credentials.URL,
		"DYNATRACE_API_TOKEN=" + me.Credentials.Token,
		"DT_CACHE_FOLDER=" + cacheFolder,
		"CACHE_OFFLINE_MODE=true",
		"DT_CACHE_DELETE_ON_LAUNCH=false",
		"DT_NO_CACHE_CLEANUP=true",
		"DT_TERRAFORM_IMPORT=true",
	}
	cmd.Start()
	if err := cmd.Wait(); err != nil {
		fmt.Println("out:", outb.String())
		fmt.Println("err:", errb.String())
	}

	return nil
}

func (me *Environment) FinishExport() error {

	fmt.Println("Finish Export ...")
	address.SaveOriginalMap(me.OutputFolder)
	address.SaveCompletedMap(me.OutputFolder)

	if QUICK_INIT {
		err := me.WriteQuickModulesJSON()
		if err != nil {
			return err
		}
	} else if me.Flags.SkipTerraformInit {
		// pass
	} else {
		err := me.RunTerraformInit()
		if err != nil {
			return err
		}
	}

	if me.Flags.ImportStateV2 {
		if err := me.ExecuteImport(); err != nil {
			return err
		}
	}
	return nil
}

func (me *Environment) RunTerraformInit() error {
	exePath, err := exec.LookPath("terraform")
	fmt.Println("Terraform executable path: ", exePath)
	if err != nil {
		fmt.Println("Terraform executable path error: ", err)
	}
	cmdOptions := []string{"init", "-no-color"}

	customProviderLocation := os.Getenv(ENV_VAR_CUSTOM_PROVIDER_LOCATION)
	if len(customProviderLocation) != 0 && customProviderLocation != "" {
		cmdOptions = append(cmdOptions, fmt.Sprint("-plugin-dir=", customProviderLocation))
	}

	cmd := exec.Command(exePath, cmdOptions...)
	cmd.Dir = me.OutputFolder
	cmd.Env = os.Environ()
	outs, err := cmd.StderrPipe()
	if err != nil {
		return err
	}
	err = cmd.Start()
	if err != nil {
		fmt.Println("Terraform CLI not installed - skipping import, error: ", err)
		return nil
	} else {
		fmt.Println("Executing 'terraform init'")
		defer func() {
			cmd.Wait()
		}()

		go readStuff(bufio.NewScanner(outs))
	}
	return nil
}

type TerraformInitModuleList struct {
	Modules []TerraformInitModule `json:"Modules"`
}
type TerraformInitModule struct {
	Key    string `json:"Key"`
	Source string `json:"Source"`
	Dir    string `json:"Dir"`
}

func (me *Environment) RunQuickInit() error {

	if QUICK_INIT {
		// pass
	} else {
		return nil
	}

	fmt.Println("Executing Quick Init ...")
	err := me.WriteMainProviderFile()
	if err != nil {
		return err
	}

	err = me.RunTerraformInit()
	if err != nil {
		return err
	}

	return nil
}

func (me *Environment) WriteQuickModulesJSON() error {

	if QUICK_INIT {
		// pass
	} else {
		return nil
	}

	modules := []TerraformInitModule{}
	modules = append(modules, TerraformInitModule{Key: "", Source: "", Dir: ""})

	for _, module := range me.Modules {
		moduleNameTrimmed := module.Type.Trim()
		modules = appendModule(modules, moduleNameTrimmed)

		if module.SplitPathModuleNameMap != nil {
			for _, splitName := range module.SplitPathModuleNameMap {
				modules = appendModule(modules, splitName)
			}
		}

	}

	tfInitModules := TerraformInitModuleList{
		Modules: modules,
	}

	bytes, err := json.MarshalIndent(tfInitModules, "", "  ")
	if err != nil {
		return err
	}

	modulesDir := filepath.Join(me.OutputFolder, ".terraform", "modules")
	os.MkdirAll(modulesDir, os.ModePerm)
	fs := afero.NewOsFs()
	filename := fmt.Sprint(filepath.Join(modulesDir, "modules.json"))
	err = afero.WriteFile(fs, filename, bytes, 0664)
	if err != nil {
		return err
	}
	return nil
}

func readStuff(scanner *bufio.Scanner) {
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}
}

func appendModule(modules []TerraformInitModule, moduleNameTrimmed string) []TerraformInitModule {
	modules = append(modules, TerraformInitModule{
		Key:    moduleNameTrimmed,
		Source: fmt.Sprintf("./modules/%s", moduleNameTrimmed),
		Dir:    fmt.Sprintf("modules/%s", moduleNameTrimmed),
	})
	return modules
}
