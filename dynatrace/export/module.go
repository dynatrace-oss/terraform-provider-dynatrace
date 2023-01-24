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
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings/services/cache"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/shutdown"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/provider/version"
)

type Module struct {
	Environment *Environment
	Type        ResourceType
	Resources   map[string]*Resource
	DataSources map[string]*DataSource
	namer       UniqueNamer
	Status      ModuleStatus
	Error       error
	Descriptor  *ResourceDescriptor
	Service     settings.CRUDService[settings.Settings]
}

func (me *Module) DataSource(id string) *DataSource {
	if dataSource, found := me.DataSources[id]; found {
		return dataSource
	}
	dataSource := me.Environment.DataSource(id)
	if dataSource != nil {
		me.DataSources[id] = dataSource
	}
	return dataSource
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
	for _, resource := range me.Resources {
		if me.Environment.RefersTo(resource) {
			resources = append(resources, resource)
		}
	}
	return resources
}

func (me *Module) GetReferencedResourceTypes() []ResourceType {
	resourceTypes := map[ResourceType]ResourceType{}
	for _, referencedResource := range me.GetResourceReferences() {
		if referencedResource.Type == me.Type {
			continue
		}
		resourceTypes[referencedResource.Type] = referencedResource.Type
	}
	result := []ResourceType{}
	for resourceType := range resourceTypes {
		result = append(result, resourceType)
	}
	return result
}

func (me *Module) GetResourceReferences() []*Resource {
	resources := map[string]*Resource{}
	if len(me.Resources) == 0 {
		return []*Resource{}
	}
	for _, resource := range me.Resources {
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
		result = append(result, resource)
	}
	return result
}

func (me *Module) Resource(id string) *Resource {
	if stored, found := me.Resources[id]; found {
		return stored
	}
	res := &Resource{ID: id, Type: me.Type, Module: me, Status: ResourceStati.Discovered}
	me.Resources[id] = res
	return res
}

func (me *Module) MkdirAll(flawed bool) error {
	if flawed {
		return os.MkdirAll(me.GetFlawedFolder(), os.ModePerm)
	}
	return os.MkdirAll(me.GetFolder(), os.ModePerm)
}

func (me *Module) GetFolder(relative ...bool) string {
	if me.Environment.Flags.Flat {
		return me.Environment.GetFolder()
	}
	if len(relative) == 0 || !relative[0] {
		return path.Join(me.Environment.GetFolder(), path.Join("modules", me.Type.Trim()))
	}
	return path.Join("modules", me.Type.Trim())
}

func (me *Module) GetAttentionFolder(relative ...bool) string {
	if me.Environment.Flags.Flat {
		return me.Environment.GetAttentionFolder()
	}
	if len(relative) == 0 || !relative[0] {
		return path.Join(me.Environment.GetAttentionFolder(), path.Join(me.Type.Trim()))
	}
	return path.Join(me.Type.Trim())
}

func (me *Module) GetFlawedFolder(relative ...bool) string {
	if me.Environment.Flags.Flat {
		return me.Environment.GetFlawedFolder()
	}
	if len(relative) == 0 || !relative[0] {
		return path.Join(me.Environment.GetFlawedFolder(), path.Join(me.Type.Trim()))
	}
	return path.Join(me.Type.Trim())
}

func (me *Module) GetFile(name string) string {
	return path.Join(me.GetFolder(), name)
}

func (me *Module) OpenFile(name string) (file *os.File, err error) {
	return os.OpenFile(me.GetFile(name), os.O_APPEND|os.O_CREATE, 0666)
}

func (me *Module) CreateFile(name string) (*os.File, error) {
	fileName := me.GetFile(name)
	os.MkdirAll(filepath.Dir(fileName), os.ModePerm)
	return os.Create(fileName)
}

func (me *Module) WriteProviderFile() (err error) {
	if me.Environment.Flags.Flat {
		return nil
	}
	if !me.ContainsPostProcessedResources() {
		return
	}
	if err = me.MkdirAll(false); err != nil {
		return err
	}
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

func (me *Module) WriteVariablesFile() (err error) {
	if me.Environment.Flags.Flat {
		return nil
	}
	if !me.ContainsPostProcessedResources() {
		return
	}
	referencedResourceTypes := me.GetReferencedResourceTypes()
	if len(referencedResourceTypes) == 0 {
		return nil
	}

	var variablesFile *os.File
	if variablesFile, err = me.CreateFile("___variables___.tf"); err != nil {
		return err
	}
	defer func() {
		variablesFile.Close()

		exePath, _ := exec.LookPath("terraform.exe")
		cmd := exec.Command(exePath, "fmt", variablesFile.Name())
		cmd.Start()
		cmd.Wait()
	}()

	for _, resourceType := range referencedResourceTypes {
		if _, err = variablesFile.WriteString(fmt.Sprintf(`variable "%s" {
			type = any
		}
		
		`, resourceType)); err != nil {
			return err
		}
	}
	return nil
}

func (me *Module) WriteDataSourcesFile() (err error) {
	if me.Environment.Flags.Flat {
		return nil
	}

	var datasourcesFile *os.File
	if datasourcesFile, err = me.CreateFile("___datasources___.tf"); err != nil {
		return err
	}
	defer func() {
		datasourcesFile.Close()
		format(datasourcesFile.Name(), true)
	}()
	for dataSourceID, dataSource := range me.DataSources {
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

func (me *Module) WriteResourcesFile() (err error) {
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
	return nil
}

func (me *Module) RefersTo(resource *Resource) bool {
	if resource == nil {
		return false
	}
	if me.Type == resource.Type {
		return false
	}
	for _, res := range me.Resources {
		if res.RefersTo(resource) {
			return true
		}
	}
	return false
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

func (me *Module) Download(keys ...string) (err error) {
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

	// 	fmt.Print(ClearLine)
	// 	fmt.Print("\r")
	// 	fmt.Printf("  - %s (%d of %d)", me.Type, idx, length)
	// }
	// fmt.Print(ClearLine)
	// fmt.Print("\r")
	// fmt.Printf("  - %s\n", me.Type)

	const ClearLine = "\033[2K"
	if len(keys) == 0 {
		length := len(me.Resources)
		fmt.Printf("Downloading \"%s\" (0 of %d)", me.Type, length)
		idx := 0
		for _, resource := range me.Resources {
			if err := resource.Download(); err != nil {
				return err
			}
			idx++
			fmt.Print(ClearLine)
			fmt.Print("\r")
			fmt.Printf("Downloading \"%s\" (%d of %d)", me.Type, idx, length)
		}
		fmt.Print(ClearLine)
		fmt.Print("\r")
		fmt.Printf("Downloading \"%s\"\n", me.Type)
		return nil
	}
	resourcesToDownload := []*Resource{}
	for _, key := range keys {
		for _, resource := range me.Resources {
			if resource.ID == key {
				resourcesToDownload = append(resourcesToDownload, resource)
			}
		}
	}
	length := len(me.Resources)
	fmt.Printf("Downloading \"%s\" (0 of %d)", me.Type, length)
	idx := 0
	for _, resource := range resourcesToDownload {
		if err := resource.Download(); err != nil {
			return err
		}
		idx++
		fmt.Print(ClearLine)
		fmt.Print("\r")
		fmt.Printf("Downloading \"%s\" (%d of %d)", me.Type, idx, length)
	}
	fmt.Print(ClearLine)
	fmt.Print("\r")
	fmt.Printf("Downloading \"%s\"\n", me.Type)
	return nil
}

func (me *Module) Discover() error {
	if shutdown.System.Stopped() {
		return nil
	}

	if me.Status.IsOneOf(ModuleStati.Downloaded, ModuleStati.Discovered, ModuleStati.Erronous) {
		return nil
	}

	if me.Descriptor == nil {
		descriptor := AllResources[me.Type]
		me.Descriptor = &descriptor
	}

	if me.Service == nil {
		me.Service = me.Descriptor.Service(me.Environment.Credentials)
	}

	var err error

	var stubs []*settings.Stub
	// log.Println("Discovering \"" + me.Type + "\" ...")
	if stubs, err = me.Service.List(); err != nil {
		if strings.Contains(err.Error(), "Token is missing required scope") {
			me.Status = ModuleStati.Erronous
			me.Error = err
			return nil
		}
		if strings.Contains(err.Error(), "No schema with topic identifier") {
			me.Status = ModuleStati.Erronous
			me.Error = err
			return nil
		}
		return err
	}
	for _, stub := range stubs {
		res := me.Resource(stub.ID).SetName(stub.Name)
		if stub.LegacyID != nil {
			res.LegacyID = *stub.LegacyID
		}
	}
	me.Status = ModuleStati.Discovered
	hide(stubs)
	// log.Println("   ", fmt.Sprintf("%d items found", len(stubs)))
	return nil
}

func (me *Module) ExecuteImport() (err error) {
	if !me.Environment.Flags.ImportState {
		return nil
	}
	if me.Status.IsOneOf(ModuleStati.Imported, ModuleStati.Erronous, ModuleStati.Untouched) {
		return nil
	}
	referencedResourceTypes := me.GetReferencedResourceTypes()
	if len(referencedResourceTypes) > 0 {
		for _, resourceType := range referencedResourceTypes {
			if err := me.Environment.Module(resourceType).ExecuteImport(); err != nil {
				return err
			}
		}
	}
	length := 0
	for _, resource := range me.Resources {
		if !resource.Status.IsOneOf(ResourceStati.PostProcessed) {
			continue
		}
		length++
	}
	fmt.Printf("  - %s (0 of %d)", me.Type, length)
	exePath, _ := exec.LookPath("terraform")
	const ClearLine = "\033[2K"
	idx := 0
	for _, resource := range me.Resources {
		if !resource.Status.IsOneOf(ResourceStati.PostProcessed) {
			continue
		}
		statement := fmt.Sprintf("module.%s.%s.%s", me.Type.Trim(), me.Type, resource.UniqueName)
		if me.Environment.Flags.Flat {
			statement = fmt.Sprintf("%s.%s", me.Type, resource.UniqueName)
		}
		// log.Println("terraform", "import", statement, resource.ID)
		cmd := exec.Command(
			exePath,
			"import",
			"-lock=false",
			"-input=false",
			"-no-color",
			statement,
			resource.ID,
		)
		var outb, errb bytes.Buffer
		cmd.Stdout = &outb
		cmd.Stderr = &errb
		cmd.Dir = me.Environment.OutputFolder
		var cacheFolder string
		if cacheFolder, err = filepath.Abs(cache.GetCacheFolder()); err != nil {
			return err
		}
		cmd.Env = []string{
			// "TF_LOG_PROVIDER=INFO",
			"DYNATRACE_ENV_URL=" + me.Environment.Credentials.URL,
			"DYNATRACE_API_TOKEN=" + me.Environment.Credentials.Token,
			"DT_CACHE_FOLDER=" + cacheFolder,
			"CACHE_OFFLINE_MODE=true",
			"DELETE_CACHE_ON_LAUNCH=false",
			"DT_TERRAFORM_IMPORT=true",
			"DT_REST_DEBUG_REQUESTS=terraform-provider-dynatrace.http.log",
		}
		cmd.Start()
		if err := cmd.Wait(); err != nil {
			fmt.Println("out:", outb.String())
			fmt.Println("err:", errb.String())
		}
		idx++
		fmt.Print(ClearLine)
		fmt.Print("\r")
		fmt.Printf("  - %s (%d of %d)", me.Type, idx, length)
	}
	fmt.Print(ClearLine)
	fmt.Print("\r")
	fmt.Printf("  - %s\n", me.Type)
	me.Status = ModuleStati.Imported
	return nil
}

func hide(v any) {}
