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
	"path/filepath"
	"strings"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/address"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/export/multiuse"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/shutdown"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/provider/logging"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hclgen"
)

var SHORTER_NAMES = os.Getenv("DYNATRACE_SHORTER_NAMES") == "true"

type Resource struct {
	ID                              string
	LegacyID                        string
	Name                            string
	UniqueName                      string
	Type                            ResourceType
	Module                          *Module
	Status                          ResourceStatus
	Error                           error
	ResourceReferences              []*Resource
	DataSourceReferences            []*DataSource
	OutputFileAbs                   string
	Flawed                          bool
	XParent                         *Resource
	ParentID                        *string
	SplitId                         int
	BundleFilePath                  string
	ExtractedIdsPerDependencyModule map[string]map[string]bool
}

func (me *Resource) GetParent() *Resource {
	if me.Module.Environment.ChildResourceOverride {
		return nil
	}
	return me.XParent
}

func (me *Resource) IsReferencedAsDataSource() bool {
	return me.Module.IsReferencedAsDataSource()
}

func (me *Resource) SetName(name string) *Resource {
	if me.Name == name {
		return me
	}
	me.Name = name

	var parentUniqueNameFound bool = false
	var parentUniqueName string = ""

	parentType, parentFound := me.Module.Environment.ChildParentGroups[me.Module.Type]
	nameModule := me.Module

	if parentFound {
		nameModule = me.Module.Environment.Module(parentType)
		nameModule.ModuleMutex.Lock()
		defer nameModule.ModuleMutex.Unlock()
		parentUniqueName, parentUniqueNameFound = nameModule.ChildParentIDNameMap[me.ID]
	}

	if parentUniqueNameFound {
		me.UniqueName = parentUniqueName
	} else {
		prevUniqueName := me.Module.Environment.PrevStateMapCommon.GetPrevUniqueName(me)
		if prevUniqueName == "" {
			terraformName := toTerraformName(name)
			me.UniqueName = nameModule.namer.Name(terraformName)
		} else {
			me.UniqueName = prevUniqueName
		}
	}
	if parentFound && !parentUniqueNameFound {
		nameModule.ChildParentIDNameMap[me.ID] = me.UniqueName
	}

	me.Status = ResourceStati.Discovered
	return me
}

func (me *Resource) getTypeOfReference() string {
	parentType, parentFound := me.Module.Environment.ChildParentGroups[me.Module.Type]

	typeOfId := string(me.Type)
	if parentFound {
		typeOfId = string(parentType)
	}
	return typeOfId
}

func (me *Resource) GetResourceReferences() []*Resource {
	resources := map[string]*Resource{}

	resources = me.getResourceReferences(resources)

	result := []*Resource{}
	for _, resource := range resources {
		result = append(result, resource)
	}
	return result
}

func (me *Resource) getResourceReferences(resources map[string]*Resource) map[string]*Resource {
	if len(me.ResourceReferences) == 0 {
		return resources
	}
	for _, resource := range me.ResourceReferences {
		if !resource.Status.IsOneOf(ResourceStati.PostProcessed, ResourceStati.Downloaded) {
			continue
		}
		key := fmt.Sprintf("%s.%s", resource.ID, resource.Type)
		_, foundReference := resources[key]
		if foundReference {
			continue
		}
		resources[key] = resource

		resources = resource.getResourceReferences(resources)
	}

	return resources
}

func (me *Resource) RefersTo(other *Resource) bool {
	if other == nil {
		return false
	}
	if len(me.ResourceReferences) == 0 {
		return false
	}
	for _, reference := range me.ResourceReferences {
		if reference.Type == other.Type && reference.ID == other.ID {
			return true
		}
	}
	return false
}

func (me *Resource) CreateFile() (*os.File, error) {
	return os.Create(me.GetFile())
}

func (me *Resource) ReadFile() ([]byte, error) {
	return os.ReadFile(me.GetFile())
}

const MAX_PATH_LENGTH_FILENAME_SHORTER = 240

func (me *Resource) GetFileName() string {
	filename := fileSystemName(fmt.Sprintf("%s.%s.tf", strings.TrimSpace(me.UniqueName), me.Type.Trim()))

	if SHORTER_NAMES {
		filename = me.getShorterFileName(filename)
	}

	return filename
}

func (me *Resource) getShorterFileName(filename string) string {
	folderPath, err := filepath.Abs(me.Module.GetFolder())
	if err != nil {
		folderPath = me.Module.GetFolder()
	}

	if (len(folderPath) + len(filename)) > MAX_PATH_LENGTH_FILENAME_SHORTER {
		filename = fileSystemName(fmt.Sprintf("%s.%s.tf", GetHashName(strings.TrimSpace(me.UniqueName)), me.Type.Trim()))
	}
	return filename
}

func (me *Resource) GetFile() string {
	if me.BundleFilePath == "" {
		return path.Join(me.Module.GetFolder(), me.GetFileName())
	}

	return me.BundleFilePath
}

func (me *Resource) GetAttentionFile() string {
	return path.Join(me.Module.GetAttentionFolder(false), me.GetFileName())
}

func (me *Resource) GetFlawedFile() string {
	return path.Join(me.Module.GetFlawedFolder(false), me.GetFileName())
}

func (me *Resource) Download() error {
	if shutdown.System.Stopped() {
		return nil
	}
	if me.Status.IsOneOf(ResourceStati.Erronous, ResourceStati.Excluded, ResourceStati.Downloaded, ResourceStati.PostProcessed) {
		return nil
	}

	var err error

	if me.Module.Status == ModuleStati.Erronous {
		me.Status = ResourceStati.Erronous
	}

	if me.Module.Status == ModuleStati.Untouched {
		if err = me.Module.Discover(); err != nil {
			return err
		}
	}

	if me.Module.Descriptor.except != nil {
		if me.Module.Descriptor.except(me.ID, me.Name) {
			me.Status = ResourceStati.Excluded
			return nil
		}
	}

	var service = me.Module.Service

	settngs := me.Module.Descriptor.NewSettings()

	getID := multiuse.EncodeIDParent(me.ID, me.ParentID)

	if err = service.Get(getID, settngs); err != nil {
		if restError, ok := err.(rest.Error); ok {
			if strings.HasPrefix(restError.Message, "Editing or deleting a non user specific dashboard preset is not allowed.") {
				me.Status = ResourceStati.Erronous
				me.Error = err
				return nil
			}
			if restError.Code == 404 {
				logging.Debug.Info.Printf("[DOWNLOAD] [%s] [FAILED] [%s] 404 not found", me.Type, me.ID)
				logging.Debug.Warn.Printf("[DOWNLOAD] [%s] [FAILED] [%s] 404 not found", me.Type, me.ID)
				me.Status = ResourceStati.Erronous
				me.Error = err
				return nil
			}
			if restError.Code == 400 {
				logging.Debug.Info.Printf("[DOWNLOAD] [%s] [FAILED] [%s] 400 %s", me.Type, me.ID, err.Error())
				logging.Debug.Warn.Printf("[DOWNLOAD] [%s] [FAILED] [%s] 400 %s", me.Type, me.ID, err.Error())
				me.Status = ResourceStati.Erronous
				me.Error = err
				return nil
			}
			if restError.Code == 500 {
				logging.Debug.Info.Printf("[DOWNLOAD] [%s] [FAILED] [%s] 500 Internal Server error", me.Type, me.ID)
				logging.Debug.Warn.Printf("[DOWNLOAD] [%s] [FAILED] [%s] 500 Internal Server error", me.Type, me.ID)
				me.Status = ResourceStati.Erronous
				me.Error = err
				return nil
			}
			if strings.HasPrefix(restError.Message, "Token is missing required scope") {
				logging.Debug.Info.Printf("[DOWNLOAD] [%s] [FAILED] [%s] Token is missing required scope", me.Type, me.ID)
				logging.Debug.Warn.Printf("[DOWNLOAD] [%s] [FAILED] [%s] Token is missing required scope", me.Type, me.ID)
				me.Status = ResourceStati.Erronous
				me.Error = err
				return nil
			}
		}
		return err
	}
	name := settings.Name(settngs, me.ID)
	me.SetName(name)

	legacyID := settings.GetLegacyID(settngs)
	if legacyID != nil {
		me.LegacyID = *legacyID
	}

	idOnly, _, _ := settings.SplitID(me.ID)
	address.AddToComplete(address.AddressComplete{
		AddressOriginal: address.AddressOriginal{
			TerraformSchemaID: service.SchemaID(),
			OriginalID:        idOnly,
		},
		UniqueName:  me.UniqueName,
		Type:        string(me.Type),
		TrimmedType: me.Type.Trim(),
	})
	comments := settings.FillDemoValues(settngs)
	comments = append(comments, settings.Validate(settngs)...)

	if len(comments) > 0 {
		for _, comment := range comments {
			if strings.HasPrefix(comment, "FLAWED SETTINGS") {
				me.Flawed = true
			}
		}
	}

	finalComments := []string{}
	if me.Module.Environment.Flags.PersistIDs {
		finalComments = []string{"ID " + me.ID}
		if legacyID := settings.ClearLegacyID(settngs); legacyID != nil {
			finalComments = append(finalComments, "LEGACY_ID "+*legacyID)
		}
	}
	if len(comments) > 0 {
		for _, comment := range comments {
			if len(finalComments) > 0 {
				finalComments = append(finalComments, "")
			}
			finalComments = append(finalComments, "ATTENTION "+comment)
		}
	}

	me.Module.MkdirAll(me.Flawed)

	splitFolder, err := me.Module.getLockBundleFile(me)
	if err != nil {
		return err
	}

	var outputFile *os.File
	if splitFolder != nil {
		me.SplitId = splitFolder.splitId
		outputFile = splitFolder.currentBundleFile
	}

	if outputFile != nil {
		me.BundleFilePath = outputFile.Name()
		outputFile.Write([]byte("# BUNDLE-ITEM\n"))
		defer splitFolder.releaseUnlockBundleFile()
	} else {
		me.BundleFilePath = ""
		me.SplitId = 0
		if outputFile, err = me.CreateFile(); err != nil {
			return err
		}
		defer outputFile.Close()
	}

	if err = hclgen.ExportResource(settngs, outputFile, string(me.Type), me.UniqueName, finalComments...); err != nil {
		return err
	}

	if !me.Flawed && me.Status != ResourceStati.Erronous && len(comments) > 0 {
		orig, _ := filepath.Abs(me.GetFile())
		att, _ := filepath.Abs(me.GetAttentionFile())
		absdir, _ := filepath.Abs(path.Dir(me.GetAttentionFile()))
		os.MkdirAll(absdir, os.ModePerm)
		os.Link(orig, att)
	}
	if me.Flawed && me.Status != ResourceStati.Erronous {
		orig, _ := filepath.Abs(me.GetFile())
		att, _ := filepath.Abs(me.GetFlawedFile())
		absdir, _ := filepath.Abs(path.Dir(me.GetFlawedFile()))
		os.MkdirAll(absdir, os.ModePerm)
		os.Link(orig, att)
	}
	if me.Status != ResourceStati.Erronous {
		me.Status = ResourceStati.Downloaded
	}
	SetOptimizedRegexResource(me)
	return nil
}

func (me *Resource) PostProcess() error {
	if shutdown.System.Stopped() {
		return nil
	}

	if me.Status.IsOneOf(ResourceStati.Erronous, ResourceStati.Excluded, ResourceStati.PostProcessed) {
		return nil
	}
	var err error
	if me.Status == ResourceStati.Discovered {
		if err = me.Download(); err != nil {
			return err
		}
	}
	me.Status = ResourceStati.PostProcessed
	if me.IsReferencedAsDataSource() {
		return nil
	}
	if !me.Module.Environment.Flags.FollowReferences {
		return nil
	}
	descriptor := me.Module.Descriptor
	if len(descriptor.Dependencies) == 0 {
		return nil
	}

	var data []byte
	var foundItemsInFileContents []any
	if data, err = me.ReadFile(); err != nil {
		return err
	}
	fileContents := string(data)

	idx := strings.Index(fileContents, "\" {")
	fileHeader := fileContents[:idx]
	fileBody := fileContents[idx:]

	isModifiedFile := false

	for _, dependency := range descriptor.Dependencies {
		resourceType := dependency.ResourceType()
		if len(resourceType) > 0 {
			module := me.Module.Environment.Module(resourceType)
			if module.Status == ModuleStati.Erronous {
				continue
			}
			if !module.Status.IsOneOf(ModuleStati.Downloaded, ModuleStati.Discovered, ModuleStati.Erronous) {
				if err = module.Discover(); err != nil {
					return err
				}
			}
		}

		if fileBody, foundItemsInFileContents = dependency.Replace(me.Module.Environment, fileBody, me.Type, me.ID); len(foundItemsInFileContents) > 0 {
			isModifiedFile = true

			for _, item := range foundItemsInFileContents {
				switch typedItem := item.(type) {
				case *Resource:
					if err = typedItem.Download(); err != nil {
						return err
					}
					if dependency.IsParent() {
						me.XParent = typedItem
					}
					me.ResourceReferences = append(me.ResourceReferences, typedItem)
				case *DataSource:
					// me.DataSourceReferences = append(me.DataSourceReferences, typedItem)
				}
			}
		}
	}

	if isModifiedFile {
		var outputFile *os.File
		if outputFile, err = me.CreateFile(); err != nil {
			return err
		}
		fileContents = fileHeader + fileBody
		defer func() {
			outputFile.Close()
			format(outputFile.Name(), false)
		}()
		if _, err = outputFile.Write([]byte(fileContents)); err != nil {
			return err
		}
	} else {
		format(me.GetFile(), false)
	}

	return nil
}

func (me *Resource) GetExtractedIdsPerRegexType(idRegexType string, tfFileContent string, optimizers map[string]optimizedIdDep) map[string]bool {
	idMap, exists := me.ExtractedIdsPerDependencyModule[idRegexType]

	if exists {
		return idMap
	}

	me.ExtractedIdsPerDependencyModule[idRegexType] = map[string]bool{}

	if idRegexType == NONE {
		return me.ExtractedIdsPerDependencyModule[idRegexType]
	}

	optimizedIdDep := optimizers[idRegexType]

	optimizedMatchList := optimizedIdDep.regex.FindAll([]byte(tfFileContent), -1)
	for _, match := range optimizedMatchList {
		me.ExtractedIdsPerDependencyModule[idRegexType][string(match)] = true
	}

	return me.ExtractedIdsPerDependencyModule[idRegexType]
}
