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

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/shutdown"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/provider/logging"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/provider/version"
	"github.com/spf13/afero"
)

type Module struct {
	Environment            *Environment
	Type                   ResourceType
	Resources              map[string]*Resource
	DataSources            map[string]*DataSource
	namer                  UniqueNamer
	Status                 ModuleStatus
	Error                  error
	PrivDescriptor         *ResourceDescriptor
	Service                settings.CRUDService[settings.Settings]
	ChildParentIDNameMap   map[string]string
	ModuleMutex            *sync.Mutex
	SplitPathModuleNameMap map[string]string
	SplitList              *splitList
	ChildModules           map[ResourceType]*Module
	IdRegexType            string
	LegacyIdMap            map[string]*Resource
	DataSourceLock         *sync.Mutex
	DescriptorLock         sync.Mutex
}

func (me *Module) GetDescriptor() *ResourceDescriptor {
	me.DescriptorLock.Lock()
	defer me.DescriptorLock.Unlock()
	if me.PrivDescriptor == nil {
		if descriptor, found := AllResources[me.Type]; found {
			me.PrivDescriptor = &descriptor
		} else {
			panic(fmt.Sprintf("Tried to resolve a Resource Descriptor for resource type `%s` - that key doesn't exist in AllResource. Please contact Dynatrace.", me.Type))
		}
	}
	return me.PrivDescriptor
}

func (me *Module) GetReferringResources(resource *Resource) []*Resource {
	var resources []*Resource
	me.ModuleMutex.Lock()
	for _, res := range me.Resources {
		resources = append(resources, res)
	}
	me.ModuleMutex.Unlock()
	var referringResources []*Resource
	for _, res := range resources {
		if res.RefersTo(resource) {
			referringResources = append(referringResources, res)
		}
	}
	return referringResources
}

func (me *Module) IsReferencedAsDataSource() bool {
	if !me.Environment.Flags.DataSources {
		return false
	}
	if _, found := me.Environment.ResArgs[string(me.Type)]; found {
		return false
	}
	return me.Type == ResourceTypes.ManagementZoneV2 ||
		me.Type == ResourceTypes.Alerting ||
		me.Type == ResourceTypes.RequestAttribute ||
		me.Type == ResourceTypes.WebApplication ||
		me.Type == ResourceTypes.RequestNaming ||
		me.Type == ResourceTypes.JSONDashboard ||
		me.Type == ResourceTypes.SLO ||
		me.Type == ResourceTypes.CalculatedServiceMetric ||
		me.Type == ResourceTypes.MobileApplication ||
		me.Type == ResourceTypes.BrowserMonitor ||
		me.Type == ResourceTypes.HTTPMonitor ||
		me.Type == ResourceTypes.Credentials ||
		me.Type == ResourceTypes.SyntheticLocation ||
		me.Type == ResourceTypes.FailureDetectionParameters ||
		me.Type == ResourceTypes.UpdateWindows ||
		me.Type == ResourceTypes.AWSCredentials ||
		me.Type == ResourceTypes.AzureCredentials ||
		me.Type == ResourceTypes.IAMGroup ||
		me.Type == ResourceTypes.AppSecVulnerabilityAlerting ||
		me.Type == ResourceTypes.AppSecAttackAlerting
}

func (me *Module) DataSource(id string, kind DataSourceKind, excepts ...ResourceType) *DataSource {
	me.DataSourceLock.Lock()
	defer me.DataSourceLock.Unlock()
	if dataSource, found := me.DataSources[id]; found {
		return dataSource
	}
	dataSource := me.Environment.DataSource(id, kind, append(excepts, me.Type)...)
	if dataSource != nil {
		me.DataSources[id] = dataSource
	}
	return dataSource
}

func (me *Module) GetDataSources(dataSources map[string]*DataSource) {
	me.DataSourceLock.Lock()
	defer me.DataSourceLock.Unlock()
	for k, v := range me.DataSources {
		dataSources[k] = v
	}
}

func (me *Module) SortedDataSources() (result []*DataSource) {
	me.DataSourceLock.Lock()
	defer me.DataSourceLock.Unlock()
	dataSourceIDs := []string{}
	for dataSourceID := range me.DataSources {
		dataSourceIDs = append(dataSourceIDs, dataSourceID)
	}
	sort.Strings(dataSourceIDs)
	for _, dataSourceID := range dataSourceIDs {
		result = append(result, me.DataSources[dataSourceID])
	}
	return result
}

func (me *Module) ContainsPostProcessedResources() bool {
	for _, resource := range me.Resources {
		if resource.Status.IsOneOf(ResourceStati.PostProcessed) {
			return true
		}
	}
	return false
}

func (me *Module) GetResourcesReferencedFromOtherModules() []*Resource {
	resources := []*Resource{}

	myResources := me.GetResourcesAndChildOfResources()

	for _, resource := range myResources {
		if me.Environment.RefersTo(resource, me.Type) {
			resources = append(resources, resource)
		}
	}
	sort.Slice(resources, func(i, j int) bool {
		return resources[i].ID < resources[j].ID
	})
	return resources
}

func (me *Module) GetReferencedResourceTypes() []ResourceType {
	resourceTypes := map[ResourceType]ResourceType{}
	for _, referencedResource := range me.GetResourceReferences() {
		if referencedResource.Type == me.Type || (referencedResource.XParent != nil && referencedResource.XParent.Type == me.Type) {
			continue
		}
		resourceTypes[referencedResource.Type] = referencedResource.Type
	}
	result := []ResourceType{}
	for resourceType := range resourceTypes {
		result = append(result, resourceType)
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i] < result[j]
	})
	return result
}

func (me *Module) GetResourcesAndChildOfResources() []*Resource {
	resources := make([]*Resource, 0, len(me.Resources))

	for _, resource := range me.Resources {
		resources = append(resources, resource)
	}

	childResourcesList := me.GetChildOfResources()
	resources = append(resources, childResourcesList...)

	return resources

}

func (me *Module) GetResourceReferences() []*Resource {
	resources := map[string]*Resource{}
	myResources := me.GetResourcesAndChildOfResources()

	if len(myResources) == 0 {
		return []*Resource{}
	}
	for _, resource := range myResources {
		if !resource.Status.IsOneOf(ResourceStati.PostProcessed, ResourceStati.Downloaded) {
			continue
		}
		key := fmt.Sprintf("%s.%s", resource.ID, resource.Type)
		resources[key] = resource
		for _, resource := range resource.GetResourceReferences() {
			if !resource.Status.IsOneOf(ResourceStati.PostProcessed, ResourceStati.Downloaded) {
				continue
			}
			key := fmt.Sprintf("%s.%s", resource.ID, resource.Type)
			resources[key] = resource
		}
	}
	result := []*Resource{}
	for _, resource := range resources {
		isReferedTo := false

		if me.RefersTo(resource, resource.Type) {
			isReferedTo = true
		}

		for _, childModule := range me.ChildModules {
			if childModule.RefersTo(resource, resource.Type) {
				isReferedTo = true
			}
		}

		if isReferedTo {
			result = append(result, resource)
		}

	}
	return result
}

func (me *Module) Resource(id string) *Resource {
	me.ModuleMutex.Lock()
	defer me.ModuleMutex.Unlock()
	if stored, found := me.Resources[id]; found {
		return stored
	}
	res := &Resource{ID: id, Type: me.Type, Module: me, Status: ResourceStati.Discovered, ExtractedIdsPerDependencyModule: map[string]map[string]bool{}}
	me.Resources[id] = res
	return res
}

var mkdirMutex = new(sync.Mutex)

func (me *Module) MkdirAll(flawed bool) error {
	mkdirMutex.Lock()
	defer mkdirMutex.Unlock()
	if flawed {
		if err := os.MkdirAll(me.GetFlawedFolder(), os.ModePerm); err != nil {
			return err
		}
	}
	return os.MkdirAll(me.GetFolder(), os.ModePerm)
}

func (me *Module) FolderNameOverride() string {
	if descriptor := me.GetDescriptor(); descriptor != nil {
		return descriptor.FolderName
	}
	return ""
}

func (me *Module) GetFolder(relative ...bool) string {
	if me.Environment.Flags.Flat {
		return me.Environment.GetFolder()
	}
	if len(relative) == 0 || !relative[0] {
		return path.Join(me.Environment.GetFolder(), path.Join("modules", me.Type.GetFolderName(me.FolderNameOverride())))
	}
	return path.Join("modules", me.Type.GetFolderName(me.FolderNameOverride()))
}

func (me *Module) GetAttentionFolder(relative ...bool) string {
	if me.Environment.Flags.Flat {
		return me.Environment.GetAttentionFolder()
	}
	if len(relative) == 0 || !relative[0] {
		return path.Join(me.Environment.GetAttentionFolder(), path.Join(me.Type.GetFolderName(me.FolderNameOverride())))
	}
	return path.Join(me.Type.GetFolderName(me.FolderNameOverride()))
}

func (me *Module) GetFlawedFolder(relative ...bool) string {
	if me.Environment.Flags.Flat {
		return me.Environment.GetFlawedFolder()
	}
	if len(relative) == 0 || !relative[0] {
		return path.Join(me.Environment.GetFlawedFolder(), path.Join(me.Type.GetFolderName(me.FolderNameOverride())))
	}
	return path.Join(me.Type.GetFolderName(me.FolderNameOverride()))
}

func (me *Module) GetFile(name string) string {
	return path.Join(me.GetFolder(), name)
}

func (me *Module) GetFileSpecificPath(name string, specificPath string) string {
	return path.Join(specificPath, name)
}

func (me *Module) OpenFile(name string) (file *os.File, err error) {
	return os.OpenFile(me.GetFile(name), os.O_APPEND|os.O_CREATE, 0666)
}

func (me *Module) CreateFile(name string) (*os.File, error) {
	fileName := me.GetFile(name)
	os.MkdirAll(filepath.Dir(fileName), os.ModePerm)
	return os.Create(fileName)
}

func (me *Module) CreateFileSpecificPath(name string, specificPath string) (*os.File, error) {
	fileName := me.GetFileSpecificPath(name, specificPath)
	os.MkdirAll(filepath.Dir(fileName), os.ModePerm)
	return os.Create(fileName)
}

func (me *Module) WriteProviderFile(logToScreen bool) (err error) {
	if me.IsReferencedAsDataSource() {
		return nil
	}
	if me.Environment.Flags.Flat {
		return nil
	}
	if !me.ContainsPostProcessedResources() {
		return
	}
	if logToScreen {
		fmt.Println("- " + me.Type)
	}
	err = me.writeProviderFile(me.GetFolder())
	if err != nil {
		return err
	}
	if me.SplitPathModuleNameMap != nil {
		for specificPath := range me.SplitPathModuleNameMap {
			err = me.writeProviderFile(specificPath)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (me *Module) writeProviderFile(specificPath string) error {
	var err error

	if err = me.MkdirAll(false); err != nil {
		return err
	}
	var outputFile *os.File
	if outputFile, err = me.CreateFileSpecificPath("___providers___.tf", specificPath); err != nil {
		return err
	}
	defer func() {
		outputFile.Close()
		format(outputFile.Name(), true)
	}()
	providerSource := "dynatrace-oss/dynatrace"
	providerVersion := version.Current
	if value := os.Getenv("DYNATRACE_PROVIDER_SOURCE"); len(value) != 0 {
		providerSource = value
	}
	if value := os.Getenv("DYNATRACE_PROVIDER_VERSION"); len(value) != 0 {
		providerVersion = value
	}

	if _, err = outputFile.WriteString(fmt.Sprintf(`
	terraform {
		required_providers {
		  dynatrace = {
			source = "%s"
			version = "%s"
		  }
		}
	  } 
`, providerSource, providerVersion)); err != nil {
		return err
	}

	return nil
}

func (me *Module) WriteVariablesFile(logToScreen bool) (err error) {
	if me.IsReferencedAsDataSource() {
		return nil
	}
	if me.GetDescriptor().Parent != nil {
		return nil
	}
	if me.Environment.Flags.Flat {
		return nil
	}
	if !me.ContainsPostProcessedResources() {
		return
	}
	var variablesFile *os.File
	if variablesFile, err = me.CreateFile("___variables___.tf"); err != nil {
		return err
	}

	defer func() {
		variablesFile.Close()

		if HCL_NO_FORMAT {
			// pass
		} else {
			exePath, _ := exec.LookPath("terraform.exe")
			cmd := exec.Command(exePath, "fmt", variablesFile.Name())
			cmd.Start()
			cmd.Wait()
		}
	}()

	if ATOMIC_DEPENDENCIES {
		uniqueNameExists := map[string]bool{}
		referencedResources := me.GetResourceReferences()
		if len(referencedResources) == 0 {
			return nil
		}
		if logToScreen {
			fmt.Println("- " + me.Type)
		}

		sort.Slice(referencedResources, func(i, j int) bool {
			if referencedResources[i].UniqueName == referencedResources[j].UniqueName {
				return referencedResources[i].UniqueName < referencedResources[j].UniqueName
			}
			return referencedResources[i].UniqueName < referencedResources[j].UniqueName
		})
		for _, resource := range referencedResources {
			if resource.Type == me.Type || (resource.XParent != nil && resource.XParent.Type == me.Type) {
				continue
			}

			typeAndUniqueName := resource.Type.Trim() + resource.UniqueName
			if uniqueNameExists[typeAndUniqueName] {
				continue
			}
			uniqueNameExists[typeAndUniqueName] = true

			if !me.Environment.Module(resource.Type).IsReferencedAsDataSource() {
				if _, err = variablesFile.WriteString(fmt.Sprintf(`variable "%s_%s" {
	type = any
}
				
`, resource.Type, resource.UniqueName)); err != nil {
					return err
				}
			}
		}
	} else {
		referencedResourceTypes := me.GetReferencedResourceTypes()
		if len(referencedResourceTypes) == 0 {
			return nil
		}

		if logToScreen {
			fmt.Println("- " + me.Type)
		}

		sort.Slice(referencedResourceTypes, func(i, j int) bool {
			return referencedResourceTypes[i] < referencedResourceTypes[j]
		})
		for _, resourceType := range referencedResourceTypes {
			if !me.Environment.Module(resourceType).IsReferencedAsDataSource() {
				if _, err = variablesFile.WriteString(fmt.Sprintf(`variable "%s" {
	type = any
}
		
`, resourceType)); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (me *Module) WriteDataSourcesFile(logToScreen bool) (err error) {
	if me.IsReferencedAsDataSource() {
		return nil
	}
	if !me.Environment.ChildResourceOverride && me.GetDescriptor().Parent != nil {
		return nil
	}
	if me.Environment.Flags.Flat {
		return nil
	}
	if logToScreen {
		fmt.Println("- " + me.Type)
	}
	buf := new(bytes.Buffer)
	dsm := map[string]string{}
	for _, v := range me.Resources {
		for _, referencedResource := range v.ResourceReferences {
			if !me.Environment.Module(referencedResource.Type).IsReferencedAsDataSource() {
				continue
			}
			if asDS := AsDataSource(referencedResource); len(asDS) > 0 {
				dsm[asDS] = asDS
			}
		}
	}
	for ds := range dsm {
		buf.Write([]byte("\n" + ds))
	}
	for _, dataSource := range me.SortedDataSources() {
		dataSourceName := dataSource.Name
		dataSourceID := dataSource.ID
		dd, _ := json.Marshal(dataSourceName)
		if dataSourceID == "tenant" {
			if _, err = buf.WriteString(`
			data "dynatrace_tenant" "tenant" {
			}`); err != nil {
				return err
			}
		} else {
			if _, err = buf.WriteString(fmt.Sprintf(`
			data "dynatrace_entity" "%s" {
				type = "%s"
				name = %s
			}`, dataSourceID, dataSource.Type, string(dd))); err != nil {
				return err
			}
		}
	}
	data := buf.Bytes()
	if len(data) > 0 {
		var datasourcesFile *os.File
		if datasourcesFile, err = me.CreateFile("___datasources___.tf"); err != nil {
			return err
		}
		if _, err := datasourcesFile.Write(data); err != nil {
			return err
		}
		defer func() {
			datasourcesFile.Close()
			format(datasourcesFile.Name(), true)
		}()
	}

	return nil
}

func (me *Module) ProvideDataSources() (dsm map[string]string, err error) {
	if me.IsReferencedAsDataSource() {
		return map[string]string{}, nil
	}
	if !me.Environment.ChildResourceOverride && me.GetDescriptor().Parent != nil {
		return map[string]string{}, nil
	}
	dsm = map[string]string{}
	for _, v := range me.Resources {
		for _, referencedResource := range v.ResourceReferences {
			if asDS := AsDataSource(referencedResource); len(asDS) > 0 {
				dsm[string(referencedResource.Type)+"."+referencedResource.ID] = asDS
			}
		}
	}
	for _, dataSource := range me.SortedDataSources() {
		dataSourceName := dataSource.Name
		dataSourceID := dataSource.ID
		dd, _ := json.Marshal(dataSourceName)
		dsm["dynatrace_entity."+dataSource.Type+"."+string(dd)] = fmt.Sprintf(`data "dynatrace_entity" "%s" {
			type = "%s"
			name = %s
		}`, dataSourceID, dataSource.Type, string(dd))
	}
	return dsm, nil
}

func (me *Module) PurgeFolder() (err error) {
	if me.Environment.Flags.Flat {
		for _, resource := range me.Resources {
			os.Remove(resource.GetFile())
		}
	} else {
		if err = os.RemoveAll(me.GetFolder()); err != nil {
			return err
		}
		if err = os.RemoveAll(me.GetAttentionFolder()); err != nil {
			return err
		}
	}
	return nil
}

func (me *Module) WriteResourcesFile() (err error) {
	if me.IsReferencedAsDataSource() {
		return nil
	}
	if me.Environment == nil {
		return nil
	}
	if !me.Environment.ChildResourceOverride && me.GetDescriptor().Parent != nil {
		return nil
	}
	if me.Environment.Flags.Flat {
		return nil
	}
	if !me.ContainsPostProcessedResources() {
		return
	}
	referencedResources := me.GetResourcesReferencedFromOtherModules()
	if len(referencedResources) == 0 {
		return nil
	}

	var resourcesFile *os.File
	if resourcesFile, err = me.CreateFile("___resources___.tf"); err != nil {
		return err
	}
	defer func() {
		resourcesFile.Close()
		format(resourcesFile.Name(), true)
	}()

	if ATOMIC_DEPENDENCIES {

		uniqueNameExists := map[string]bool{}
		for _, resource := range referencedResources {

			if uniqueNameExists[resource.UniqueName] {
				continue
			}
			uniqueNameExists[resource.UniqueName] = true

			if _, err = resourcesFile.WriteString(fmt.Sprintf(`output "resources_%s" {
  value = {
  `, resource.UniqueName)); err != nil {
				return err
			}

			if _, err = resourcesFile.WriteString(fmt.Sprintf(`  value = %s.%s
	  `, resource.Type, resource.UniqueName)); err != nil {
				return err
			}

			if _, err = resourcesFile.WriteString(`}
  }
`); err != nil {
				return err
			}
		}
	} else {
		if _, err = resourcesFile.WriteString(`output "resources" {
  value = {
  `); err != nil {
			return err
		}

		for _, resource := range referencedResources {
			if _, err = resourcesFile.WriteString(fmt.Sprintf(`  %s = %s.%s
			  `, resource.UniqueName, resource.Type, resource.UniqueName)); err != nil {
				return err
			}
		}
		if _, err = resourcesFile.WriteString(`}
  }
`); err != nil {

			return err
		}

	}
	return nil
}

func (me *Module) RefersTo(resource *Resource, parentType ResourceType) bool {
	if resource == nil {
		return false
	}
	if me.Type == resource.Type || (me.GetDescriptor().Parent != nil && *me.GetDescriptor().Parent == parentType) {
		return false
	}
	for _, res := range me.Resources {
		if res.RefersTo(resource) {
			return true
		}
	}
	return false
}

func (me *Module) IsReferenced(resource *Resource) bool {
	me.ModuleMutex.Lock()
	var resources []*Resource
	for _, res := range me.Resources {
		resources = append(resources, res)
	}
	me.ModuleMutex.Unlock()
	for _, res := range resources {
		if res.RefersTo(resource) {
			return true
		}
	}
	return false
}

func (me *Module) GetChildOfResources() []*Resource {
	resources := []*Resource{}
	if me.Environment.ChildResourceOverride {
		return resources
	}

	for _, module := range me.Environment.Modules {
		childDescriptor := module.GetDescriptor()
		isParent := !me.Environment.ChildResourceOverride && childDescriptor != nil && childDescriptor.Parent != nil && string(*childDescriptor.Parent) == string(me.Type)
		if isParent {
			for _, resource := range module.Resources {
				if resource.Status == ResourceStati.PostProcessed && resource.GetParent() != nil {
					resources = append(resources, resource)
				}
			}
		}
	}

	return resources
}

func (me *Module) GetChildResources() []*Resource {
	resources := []*Resource{}
	if me.Environment.ChildResourceOverride {
		return resources
	}
	for _, resource := range me.Resources {
		if resource.Status == ResourceStati.PostProcessed && resource.GetParent() != nil {
			resources = append(resources, resource)
		}
	}
	return resources
}

func (me *Module) GetNonPostProcessedResources() []*Resource {
	resources := []*Resource{}
	for _, resource := range me.Resources {
		if resource.Status == ResourceStati.Downloaded {
			resources = append(resources, resource)
		}
	}
	return resources
}

func (me *Module) GetPostProcessedResources() []*Resource {
	resources := []*Resource{}
	for _, resource := range me.Resources {
		if resource.Status == ResourceStati.PostProcessed {
			resources = append(resources, resource)
		}
	}
	return resources
}

const ClearLine = "\033[2K"

func (me *Module) Download(multiThreaded bool, keys ...string) (err error) {
	if shutdown.System.Stopped() {
		return nil
	}

	if me.Status.IsOneOf(ModuleStati.Erronous) {
		return nil
	}
	if !me.Status.IsOneOf(ModuleStati.Discovered) {
		if err := me.Discover(); err != nil {
			return err
		}
	}

	resourcesToDownload := []*Resource{}

	if len(keys) == 0 {
		for _, resource := range me.Resources {
			resourcesToDownload = append(resourcesToDownload, resource)
		}
	} else {
		for _, key := range keys {
			for _, resource := range me.Resources {
				if resource.ID == key {
					resourcesToDownload = append(resourcesToDownload, resource)
				}
			}
		}

	}

	err = me.downloadResources(resourcesToDownload, multiThreaded)
	if err != nil {
		return err
	}

	return nil
}

func (me *Module) blockPrevNames() {
	if PREV_STATE_ON {
		names, found := me.Environment.PrevNamesByModule[string(me.Type)]
		if found {
			for _, name := range names {
				me.namer.BlockName(name)
			}
		}
	}
}

func (me *Module) CreateBundleFile(bundle bundleToDownload) (*os.File, error) {
	me.MkdirAll(false)
	file, err := me.GetBundleFile(bundle)
	if err != nil {
		return nil, err
	}

	return os.Create(file)
}

func (me *Module) GetBundleFileName(bundle bundleToDownload) string {
	return fileSystemName(fmt.Sprintf("bundle_%d.%s.tf", bundle.bundleId, me.Type.Trim()))
}

func (me *Module) GetBundleFile(bundle bundleToDownload) (string, error) {
	folderPath := me.GetFolder()

	if bundle.splitId > 0 {
		folderPath = fmt.Sprintf("%s_%d", folderPath, bundle.splitId)
		err := os.MkdirAll(folderPath, os.ModePerm)
		if err != nil {
			return "", err
		}
		splitName := fmt.Sprintf("%s_%d", me.Type.Trim(), bundle.splitId)
		me.saveSplitPath(folderPath, splitName)
	}

	return path.Join(folderPath, me.GetBundleFileName(bundle)), nil
}

func (me *Module) saveSplitPath(folderPath string, splitName string) {
	me.ModuleMutex.Lock()
	if me.SplitPathModuleNameMap == nil {
		me.SplitPathModuleNameMap = map[string]string{}
	}
	me.SplitPathModuleNameMap[folderPath] = splitName
	me.ModuleMutex.Unlock()
}

func (me *Module) saveChildModule(childModule *Module) {
	_, exist := me.ChildModules[childModule.Type]

	if exist {
		return
	}

	me.ChildModules[childModule.Type] = childModule
}

type bundleToDownload struct {
	bundleId int
	splitId  int
}

func (me *Module) downloadResources(resourcesToDownload []*Resource, multiThreaded bool) error {
	length := len(me.Resources)
	resourcesCount := len(resourcesToDownload)
	maxThreads := 50

	bundleConfig := me.getResourcesPerBundle(resourcesCount, maxThreads)
	bundleText := ""
	if bundleConfig.resourcesPerBundle > 1 {
		bundleText = fmt.Sprintf(" in %d bundles of %d resources", bundleConfig.bundlesCount, bundleConfig.resourcesPerBundle)
	}
	splitText := ""
	if bundleConfig.splitCount > 1 {
		splitText = fmt.Sprintf(" in %d splits of %d bundles ", bundleConfig.splitCount, bundleConfig.bundlesPerSplitModule)
	}

	if me.SplitList == nil {
		me.SplitList = &splitList{
			folders:      make(map[int]*splitFolder, bundleConfig.splitCount),
			folderMutex:  new(sync.Mutex),
			bundleConfig: bundleConfig,
		}
		defer me.closeLastBundleFiles()
	}

	if multiThreaded {
		fmt.Printf("Downloading \"%s\" Count:  %d %s %s\n", me.Type, resourcesCount, bundleText, splitText)
	} else {
		fmt.Printf("Downloading \"%s\" (0 of %d) %s %s", me.Type, resourcesCount, bundleText, splitText)
	}

	idx := 0
	var wg sync.WaitGroup
	itemCount := len(resourcesToDownload)
	channel := make(chan *Resource, itemCount)
	mutex := sync.Mutex{}
	if !multiThreaded {
		maxThreads = 1
	}
	if maxThreads > itemCount {
		maxThreads = itemCount
	}
	wg.Add(maxThreads)

	processItem := func(resource *Resource) error {
		if err := resource.Download(); err != nil {
			return err
		}
		mutex.Lock()
		idx++
		if !multiThreaded {
			fmt.Print(ClearLine)
			fmt.Print("\r")
			fmt.Printf("Downloading \"%s\" (%d of %d)", me.Type, idx, length)
		}
		mutex.Unlock()
		return nil
	}

	for i := 0; i < maxThreads; i++ {

		go func() error {

			for {
				resourceLoop, ok := <-channel
				if !ok {
					wg.Done()
					return nil
				}
				if shutdown.System.Stopped() {
					wg.Done()
					return nil
				}

				err := processItem(resourceLoop)

				if err != nil {
					wg.Done()
					return err
				}
			}
		}()

	}

	if len(resourcesToDownload) > 0 {
		logging.Debug.Info.Printf("[DOWNLOAD] [%s]", me.Type)
	}
	for _, resource := range resourcesToDownload {
		channel <- resource
	}

	close(channel)
	wg.Wait()
	if !multiThreaded {
		fmt.Print(ClearLine)
		fmt.Print("\r")
		fmt.Printf("Downloading \"%s\"\n", me.Type)
	}
	return nil
}

type splitList struct {
	folders      map[int]*splitFolder
	folderMutex  *sync.Mutex
	bundleConfig bundlingConfig
}

type splitFolder struct {
	splitId               int
	currentBundleId       int
	currentBundleFile     *os.File
	currentBundleResCount int
	splitList             *splitList
	splitMutex            *sync.Mutex
}

func (me *splitList) getSplitFolder(splitIdRes int) *splitFolder {
	(*me).folderMutex.Lock()
	defer (*me).folderMutex.Unlock()

	folder, found := (*me).folders[splitIdRes]

	if found {
		// pass
	} else {
		folder = &splitFolder{
			splitId:               splitIdRes,
			currentBundleId:       -1,
			currentBundleFile:     nil,
			currentBundleResCount: 0,
			splitList:             me,
			splitMutex:            new(sync.Mutex),
		}
		(*me).folders[splitIdRes] = folder
	}

	return folder
}

func (me *Module) getLockBundleFile(resource *Resource) (*splitFolder, error) {

	splitId, importSplitFound := me.Environment.ImportStateMap.GetResourceSplitId(resource)

	hasSplitOrBundles := false
	if me.SplitList != nil {
		hasSplitOrBundles = me.SplitList.bundleConfig.splitCount > 1 || me.SplitList.bundleConfig.resourcesPerBundle > 1
	}

	if importSplitFound {
		if splitId == 0 {
			if hasSplitOrBundles {
				// pass
			} else {
				return nil, nil
			}
		}
	} else {
		if hasSplitOrBundles {
			splitId = getSplitID(resource.UniqueName, me.SplitList.bundleConfig.splitCount)
		} else {
			return nil, nil
		}
	}

	splitFolder := me.SplitList.getSplitFolder(splitId)

	splitFolder.splitMutex.Lock()

	if splitFolder.currentBundleFile == nil {
		var err error
		splitFolder.currentBundleId += 1
		splitFolder.currentBundleFile, err = me.CreateBundleFile(bundleToDownload{
			splitId:  splitFolder.splitId,
			bundleId: splitFolder.currentBundleId,
		})
		splitFolder.currentBundleResCount = 0

		if err != nil {
			return nil, err
		}
	}

	return splitFolder, nil
}

func (me *splitFolder) releaseUnlockBundleFile() error {
	defer me.splitMutex.Unlock()

	me.currentBundleResCount++

	if me.currentBundleResCount >= me.splitList.bundleConfig.resourcesPerBundle {
		var err error
		me.currentBundleResCount = 0

		err = me.currentBundleFile.Close()
		if err != nil {
			return err
		}

		me.currentBundleFile = nil
	}

	return nil
}

func (me *Module) closeLastBundleFiles() error {
	me.SplitList.folderMutex.Lock()
	defer me.SplitList.folderMutex.Unlock()

	for _, splitFolder := range me.SplitList.folders {
		err := splitFolder.closeLastFile()
		if err != nil {
			return err
		}
	}

	return nil
}

func (me *splitFolder) closeLastFile() error {
	me.splitMutex.Lock()
	defer me.splitMutex.Unlock()

	if me.currentBundleFile != nil {
		err := me.currentBundleFile.Close()
		me.currentBundleFile = nil
		return err
	}

	return nil
}

func getSplitID(input string, splitCount int) int {
	hashValue := GetHash32(input)
	splitID := int(hashValue % uint32(splitCount))
	return splitID
}

type bundlingConfig struct {
	resourcesPerBundle    int
	bundlesPerSplitModule int
	bundlesCount          int
	splitCount            int
}

func (me *Module) getResourcesPerBundle(resourcesCount int, maxThreads int) bundlingConfig {
	resourcesPerBundle := 1
	bundlesPerSplitModule := resourcesCount

	if me.IsBundleImpossible() {
		return bundlingConfig{
			resourcesPerBundle:    resourcesPerBundle,
			bundlesPerSplitModule: bundlesPerSplitModule,
			bundlesCount:          1,
			splitCount:            1,
		}
	}

	resourcesForSplitModuleMin := 100
	splitCountMax := maxThreads
	bundlingDivisor := 10

	resourcesForBundlingMin := 30
	resourcesPerBundlingMax := 100

	if resourcesCount > resourcesForBundlingMin {
		resourcesPerBundle = CeilDivide(resourcesCount, bundlingDivisor)
		if resourcesPerBundle > resourcesPerBundlingMax {
			resourcesPerBundle = resourcesPerBundlingMax
		}
	}

	bundlesCount := CeilDivide(resourcesCount, resourcesPerBundle)
	splitCount := CeilDivide(resourcesCount, resourcesForSplitModuleMin)
	if splitCount > splitCountMax {
		splitCount = splitCountMax
	}

	if resourcesCount > resourcesForSplitModuleMin {
		bundlesPerSplitModule = CeilDivide(bundlesCount, splitCount)
	} else {
		bundlesPerSplitModule = bundlesCount
	}

	return bundlingConfig{
		resourcesPerBundle:    resourcesPerBundle,
		bundlesPerSplitModule: bundlesPerSplitModule,
		bundlesCount:          bundlesCount,
		splitCount:            splitCount,
	}
}

func (me *Module) IsBundleImpossible() bool {
	resourceDefinition := AllResources[me.Type]

	if ULTRA_PARALLEL {
		// pass
	} else {
		return true
	}

	if me.Environment.HasDependenciesTo[me.Type] {
		return true
	}

	if len(resourceDefinition.Dependencies) > 0 {

		hasEntityIdsDepsOnly := true
		for _, d := range resourceDefinition.Dependencies {
			switch d.(type) {
			case *entityds:
				continue
			}
			hasEntityIdsDepsOnly = false
		}

		if hasEntityIdsDepsOnly {
			// pass
		} else {
			return true
		}
	}

	if resourceDefinition.Parent != nil {
		return true
	}

	if me.Environment.IsParentMap[me.Type] {
		return true
	}

	return false
}

func (me *Module) Discover() error {
	if shutdown.System.Stopped() {
		return nil
	}

	if me.Status.IsOneOf(ModuleStati.Downloaded, ModuleStati.Discovered, ModuleStati.Erronous) {
		return nil
	}

	descriptor := me.GetDescriptor()

	if me.Service == nil {
		me.Service = descriptor.Service(me.Environment.Credentials)
	}

	var err error

	var stubs api.Stubs
	if stubs, err = me.Service.List(); err != nil {
		if strings.Contains(err.Error(), "Token is missing required scope") {
			logging.Debug.Info.Printf("[DISCOVER] [%s] Module will not get exported. Token is missing required scope.", me.Type)
			me.Status = ModuleStati.Erronous
			me.Error = err
			return nil
		}
		if strings.Contains(err.Error(), "No schema with topic identifier") {
			logging.Debug.Info.Printf("[DISCOVER] [%s] Module will not get exported. The schema doesn't exist on that environment.", me.Type)
			me.Status = ModuleStati.Erronous
			me.Error = err
			return nil
		}
		return err
	}
	stubs = stubs.Sort()
	for _, stub := range stubs {
		if stub.Name == "" {
			panic(me.Type)
		}
		res := me.Resource(stub.ID).SetName(stub.Name)
		if stub.LegacyID != nil {
			res.LegacyID = *stub.LegacyID
		}
		if stub.ParentID != nil {
			res.ParentID = stub.ParentID
		}
	}
	me.Status = ModuleStati.Discovered
	hide(stubs)

	SetOptimizedRegexModule(me)

	logging.Debug.Info.Printf("[DISCOVER] [%s] %d items found.", me.Type, len(stubs))
	return nil
}

func (me *Module) GetDependencyOptimizationInfo() (string, bool) {
	idRegexType := me.GetRegexType()

	if idRegexType == NONE {
		return idRegexType, false
	}
	return idRegexType, true
}

func (me *Module) GetRegexType() string {
	me.ModuleMutex.Lock()
	defer me.ModuleMutex.Unlock()

	if me.IdRegexType != "" {
		return me.IdRegexType
	}

	me.IdRegexType = NONE

	if len(me.Resources) == 0 {
		return me.IdRegexType
	}

	for id := range me.Resources {

		if len(id) == 0 {
			continue
		}

		for idRegexType, optimizedIdDep := range OptimizedKeyRegexId {
			matchesValidation := optimizedIdDep.regex.FindAll([]byte(id), -1)
			if len(matchesValidation) > 0 {
				if bytes.Equal(matchesValidation[0], []byte(id)) {
					me.IdRegexType = idRegexType
					break
				}
			}
		}
		break
	}

	return me.IdRegexType
}

func (me *Module) GetLegacyIdMap() map[string]*Resource {
	me.ModuleMutex.Lock()
	defer me.ModuleMutex.Unlock()

	if me.LegacyIdMap != nil {
		return me.LegacyIdMap
	}

	me.LegacyIdMap = map[string]*Resource{}

	if len(me.Resources) == 0 {
		return me.LegacyIdMap
	}

	for _, resource := range me.Resources {
		if len(resource.LegacyID) == 0 {
			continue
		}
		me.LegacyIdMap[resource.LegacyID] = resource
	}

	return me.LegacyIdMap
}

type resources []resource

type resource struct {
	Module    string     `json:"module"`
	Mode      string     `json:"mode"`
	Type      string     `json:"type"`
	Name      string     `json:"name"`
	Provider  string     `json:"provider"`
	Instances []instance `json:"instances"`
}

type instance struct {
	Attributes          attrs         `json:"attributes"`
	SchemaVersion       int           `json:"schema_version"`
	SensitiveAttributes []interface{} `json:"sensitive_attributes"`
	Private             string        `json:"private"`
}
type attrs struct {
	Id string `json:"id"`
}

func (me *Module) ExecuteImportV2(fs afero.Fs) (resList resources, err error) {
	if !me.Environment.Flags.ImportStateV2 {
		return nil, nil
	}
	if me.Status.IsOneOf(ModuleStati.Imported, ModuleStati.Erronous, ModuleStati.Untouched) {
		return nil, nil
	}

	resList = make(resources, 0, len(me.Resources))

	for _, res := range me.Resources {
		if !res.Status.IsOneOf(ResourceStati.PostProcessed) {
			continue
		}

		isWrittenAlready := me.namer.SetNameWritten(res.UniqueName)

		if isWrittenAlready {
			fmt.Println("ERROR: [ExecuteImportV2] Duplicate UniqueName for ", string(me.Type), res.UniqueName, res.ID)
			continue

		}

		providerSource := os.Getenv("DYNATRACE_PROVIDER_SOURCE")
		if len(providerSource) == 0 {
			providerSource = `provider["registry.terraform.io/dynatrace-oss/dynatrace"]`
		} else {
			providerSource = fmt.Sprintf(`provider["%s"]`, providerSource)
		}

		moduleValue := fmt.Sprintf("module.%s", me.Type.Trim())
		if me.GetDescriptor().Parent != nil {
			moduleValue = fmt.Sprintf("module.%s", me.GetDescriptor().Parent.Trim())
		}
		if res.SplitId > 0 {
			moduleValue = fmt.Sprintf("%s_%d", moduleValue, res.SplitId)
		}

		resList = append(resList, resource{
			Module: moduleValue,
			Mode:   "managed",
			Type:   string(me.Type),
			Name:   res.UniqueName,
			// Provider: `provider["dynatrace.com/com/dynatrace"]`,
			Provider: providerSource,
			Instances: []instance{
				{
					Attributes: attrs{
						Id: res.ID,
					},
					SchemaVersion:       0,
					SensitiveAttributes: make([]any, 0),
					Private:             "eyJzY2hlbWFfdmVyc2lvbiI6IjAifQ==",
				},
			},
		})
	}

	me.Status = ModuleStati.Imported

	return resList, nil
}

func hide(v any) {}

func CeilDivide(a, b int) int {
	quotient := a / b

	remainder := a % b
	if remainder > 0 {
		quotient++
	}

	return quotient
}
