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
package address

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/spf13/afero"
)

// To speed things up when using Dynatrace Config Manager
var BUILD_ADDRESS_FILES = os.Getenv("DYNATRACE_BUILD_ADDRESS_FILES") == "true"

type AddressOriginal struct {
	TerraformSchemaID string
	OriginalID        string
	OriginalSchemaID  string
}

func (ao *AddressOriginal) getKey() string {
	return ao.TerraformSchemaID + ao.OriginalID
}

type AddressComplete struct {
	AddressOriginal
	UniqueName  string
	Type        string
	TrimmedType string
}

type Address interface {
	getKey() string
}

type AddressMap struct {
	mutex     *sync.Mutex
	addresses map[string]Address
}

func NewAddressMap() *AddressMap {
	return &AddressMap{
		mutex:     new(sync.Mutex),
		addresses: map[string]Address{},
	}
}

func (al *AddressMap) AddToAddressMap(a Address) {
	if !BUILD_ADDRESS_FILES {
		return
	}

	key := a.getKey()
	al.mutex.Lock()
	value, found := al.addresses[key]
	if found {
		fmt.Printf("ERROR: Duplicate for key: %v, value: %v \n", key, value)
	}
	al.addresses[key] = a
	al.mutex.Unlock()
}

func SaveAddressMap(addresses interface{}, OutputFolder string, filename string) error {
	if !BUILD_ADDRESS_FILES {
		return nil
	}

	fs := afero.NewOsFs()
	bytes, err := json.MarshalIndent(addresses, "", "  ")
	if err != nil {
		return err
	}
	filename = fmt.Sprint(filepath.Join(OutputFolder, filename))
	err = afero.WriteFile(fs, filename, bytes, 0664)
	if err != nil {
		return err
	}

	return nil
}

var originalMap = NewAddressMap()
var completedMap = NewAddressMap()

func AddToOriginal(a AddressOriginal) {
	fmt.Println("ORIGINAL: ", a)
	originalMap.AddToAddressMap(&a)
}

func AddToComplete(a AddressComplete) {
	completedMap.AddToAddressMap(&a)
}

func SaveOriginalMap(OutputFolder string) error {
	return SaveAddressMap(originalMap.addresses, OutputFolder, "original.json")
}

func SaveCompletedMap(OutputFolder string) error {
	for key, item := range originalMap.addresses {
		_, exists := completedMap.addresses[key]
		if exists {
			completedMap.addresses[key].(*AddressComplete).AddressOriginal.OriginalSchemaID = item.(*AddressOriginal).OriginalSchemaID
		} else {
			fmt.Printf("ERROR: Missing key: %v", key)
		}
	}
	addressMap := map[string]map[string]Address{}
	for _, item := range completedMap.addresses {
		keyL1 := item.(*AddressComplete).TerraformSchemaID
		keyL2 := item.(*AddressComplete).OriginalID

		if item.(*AddressComplete).OriginalSchemaID != "" {
			keyL1 = item.(*AddressComplete).OriginalSchemaID
		}

		_, ok := addressMap[keyL1]
		if !ok {
			addressMap[keyL1] = map[string]Address{}
		}

		addressMap[keyL1][keyL2] = item
	}
	return SaveAddressMap(addressMap, OutputFolder, "address.json")
}
