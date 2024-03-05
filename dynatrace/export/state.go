/**
* @license
* Copyright 2023 Dynatrace LLC
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
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/spf13/afero"
)

var PREV_STATE_ON = os.Getenv("DYNATRACE_PREV_STATE_ON") == "true"
var PREV_STATE_PATH_THIS = os.Getenv("DYNATRACE_PREV_STATE_PATH_THIS")
var PREV_STATE_PATH_LINKED = os.Getenv("DYNATRACE_PREV_STATE_PATH_LINKED")
var IMPORT_STATE_PATH = os.Getenv("DYNATRACE_IMPORT_STATE_PATH")

type StateMap struct {
	mutex     *sync.Mutex
	resources map[string]StateResource
}

type StateResource struct {
	Resource resource
	Used     bool
}

func NewStateMap() *StateMap {
	return &StateMap{
		mutex:     new(sync.Mutex),
		resources: map[string]StateResource{},
	}
}

func (sm *StateMap) AddToStateMapByName(res resource) {
	key := fmt.Sprintf("%s|||%s",
		res.Type,
		res.Name)

	sm.AddToStateMap(key, res)
}

func (sm *StateMap) AddToStateMapByID(res resource) {
	if len(res.Instances) <= 0 {
		return
	}
	key := fmt.Sprintf("%s|||%s",
		res.Type,
		res.Instances[0].Attributes.Id)

	sm.AddToStateMap(key, res)
}

func (sm *StateMap) AddToStateMap(key string, res resource) {
	sm.mutex.Lock()
	sm.resources[key] = StateResource{
		Resource: res,
		Used:     false,
	}
	sm.mutex.Unlock()
}

func (sm *StateMap) ExtractCommonStates(smLinked *StateMap) (*StateMap, map[string][]string) {
	commonStateMap := NewStateMap()
	namesByModule := map[string][]string{}

	for key, stateResource := range sm.resources {
		stateResourceLinked, found := smLinked.resources[key]
		if found {
			commonStateMap.AddToStateMapByID(stateResource.Resource)

			// Certain types of configs have a name and ID
			// But that ID is not only a config ID, it is also an entity ID
			// This means the ID will have changed on the second run of the Config Manager
			// This would then break matching
			// The terraform provider would show an Add action on new ID and Destroy on old ID, instead of Done or Change
			linkedResourceType := ResourceType(stateResourceLinked.Resource.Type)
			if linkedResourceType == ResourceTypes.WebApplication ||
				linkedResourceType == ResourceTypes.MobileApplication ||
				linkedResourceType == ResourceTypes.BrowserMonitor ||
				linkedResourceType == ResourceTypes.HTTPMonitor {

				commonStateMap.AddToStateMapByID(stateResourceLinked.Resource)
			}

			list, found := namesByModule[stateResource.Resource.Type]
			if found {
				// pass
			} else {
				list = []string{}
			}
			list = append(list, stateResource.Resource.Name)
			namesByModule[stateResource.Resource.Type] = list
		}
	}

	return commonStateMap, namesByModule
}

func (sm *StateMap) GetPrevUniqueName(res *Resource) string {

	typeOfReference := res.getTypeOfReference()
	name := ""

	if PREV_STATE_ON {
		// pass
	} else {
		return name
	}

	sm.mutex.Lock()

	resKey, resource, isResourceFound := getResourceByID(res, sm)
	resourceOfReference, isResourceOfReferenceFound := getResourceOfReferenceByID(res, typeOfReference, sm, resKey, isResourceFound, resource, name)

	if isResourceOfReferenceFound {

		if isResourceFound {
			resource.Used = true
			sm.resources[resKey] = resource
		}

		name = resourceOfReference.Resource.Name
	}
	sm.mutex.Unlock()

	return name
}

func getResourceOfReferenceByID(res *Resource, typeOfReference string, sm *StateMap, resKey string, isResourceFound bool, resource StateResource, name string) (StateResource, bool) {

	nameKey := fmt.Sprintf("%s|||%s",
		typeOfReference,
		res.ID)

	return getResourceOfReference(nameKey, res, typeOfReference, sm, resKey, isResourceFound, resource, name)
}

func getResourceOfReferenceByUniqueName(res *Resource, typeOfReference string, sm *StateMap, resKey string, isResourceFound bool, resource StateResource, name string) (StateResource, bool) {

	nameKey := fmt.Sprintf("%s|||%s",
		typeOfReference,
		res.UniqueName)

	return getResourceOfReference(nameKey, res, typeOfReference, sm, resKey, isResourceFound, resource, name)
}

func getResourceOfReference(nameKey string, res *Resource, typeOfReference string, sm *StateMap, resKey string, isResourceFound bool, resource StateResource, name string) (StateResource, bool) {

	resourceOfReference, isResourceOfReferenceFound := sm.resources[nameKey]

	if isResourceOfReferenceFound {
		// pass
	} else if nameKey != resKey {
		if isResourceFound {
			resourceOfReference = resource
			isResourceOfReferenceFound = isResourceFound
			fmt.Printf("\n??? No Parent but child: name: %s, ID: %s, Used: %v, typeOfId: %s, resType: %s", name, res.ID, sm.resources[resKey], typeOfReference, res.Type)
		}
	}
	return resourceOfReference, isResourceOfReferenceFound
}

func (sm *StateMap) GetResourceSplitId(res *Resource) (int, bool) {
	if IMPORT_STATE_PATH == "" {
		return -1, false
	}

	typeOfReference := res.getTypeOfReference()

	resKey, resource, isResourceFound := getResourceByUniqueName(res, sm)
	resourceOfReference, isResourceOfReferenceFound := getResourceOfReferenceByUniqueName(res, typeOfReference, sm, resKey, isResourceFound, resource, "")

	if isResourceOfReferenceFound {
		return extractSplitID(resourceOfReference), true
	}

	return -1, false
}

func extractSplitID(resource StateResource) int {
	resType := ResourceType(resource.Resource.Type).Trim()
	resModule := fmt.Sprintf("module.%s", resType)

	module := resource.Resource.Module

	if resModule == module {
		return 0
	}

	prefix := fmt.Sprintf("%s_", resModule)

	if strings.HasPrefix(module, prefix) {

		numberStr := strings.TrimPrefix(module, prefix)

		if num, err := strconv.Atoi(numberStr); err == nil {
			return num
		}
	}

	return -1
}

func getResourceByID(res *Resource, sm *StateMap) (string, StateResource, bool) {
	resType := string(res.Type)

	resKey := fmt.Sprintf("%s|||%s",
		resType,
		res.ID)

	resource, isResourceFound := sm.resources[resKey]
	return resKey, resource, isResourceFound
}

func getResourceByUniqueName(res *Resource, sm *StateMap) (string, StateResource, bool) {
	resType := string(res.Type)

	resKey := fmt.Sprintf("%s|||%s",
		resType,
		res.UniqueName)

	resource, isResourceFound := sm.resources[resKey]
	return resKey, resource, isResourceFound
}

func LoadStateThis() (*state, error) {
	return LoadStateFile(PREV_STATE_PATH_THIS)
}

func LoadStateLinked() (*state, error) {
	return LoadStateFile(PREV_STATE_PATH_LINKED)
}

func LoadStateFile(filePath string) (*state, error) {
	fs := afero.NewOsFs()

	stateLoaded := &state{}

	data, err := afero.ReadFile(fs, filePath)
	if err != nil {
		return stateLoaded, nil
	}

	err = json.Unmarshal(data, stateLoaded)
	if err != nil {
		return stateLoaded, nil
	}

	return stateLoaded, nil
}

func BuildStateMap(stateLoaded *state) *StateMap {
	sm := NewStateMap()

	for _, res := range stateLoaded.Resources {
		sm.AddToStateMapByName(res)
	}

	return sm
}
