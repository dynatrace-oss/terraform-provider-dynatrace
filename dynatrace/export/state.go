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
	"sync"

	"github.com/spf13/afero"
)

var PREV_STATE_ON = os.Getenv("DYNATRACE_PREV_STATE_ON") == "true"
var PREV_STATE_PATH_THIS = os.Getenv("DYNATRACE_PREV_STATE_PATH_THIS")
var PREV_STATE_PATH_LINKED = os.Getenv("DYNATRACE_PREV_STATE_PATH_LINKED")

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
		_, found := smLinked.resources[key]
		if found {
			commonStateMap.AddToStateMapByID(stateResource.Resource)

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
	name := ""

	if PREV_STATE_ON {
		// pass
	} else {
		return name
	}

	key := fmt.Sprintf("%s|||%s",
		string(res.Type),
		res.ID)

	sm.mutex.Lock()
	resFound, found := sm.resources[key]

	if found {
		resFound.Used = true
		sm.resources[key] = resFound
		name = resFound.Resource.Name
	}
	sm.mutex.Unlock()

	return name
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
