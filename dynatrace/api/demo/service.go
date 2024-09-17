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

package demo

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api"
	demo "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/demo/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings"
	"github.com/google/uuid"
)

const SchemaID = "demo:demo"

var mu sync.Mutex

func Service(credentials *settings.Credentials) settings.CRUDService[*demo.Demo] {
	return &service{credentials}
}

type service struct {
	credentials *settings.Credentials
}

func getRecords() ([]*demo.Demo, error) {
	jsonFile, err := os.Open("records.json")
	if err != nil {
		if strings.Contains(err.Error(), "The system cannot find the file specified") {
			return []*demo.Demo{}, nil
		}
		return nil, err
	}
	defer jsonFile.Close()

	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		return nil, err
	}

	records := []*demo.Demo{}
	if err := json.Unmarshal(byteValue, &records); err != nil {
		return nil, err
	}
	return records, nil
}

func storeRecords(records []*demo.Demo) error {
	data, err := json.MarshalIndent(records, "", "  ")
	if err != nil {
		return err
	}
	os.Remove("records.json")
	err = os.WriteFile("records.json", data, 0777)
	if err != nil {
		return err
	}
	return nil
}

func (me *service) List(ctx context.Context) (api.Stubs, error) {
	mu.Lock()
	defer mu.Unlock()
	stubs := api.Stubs{}
	records, err := getRecords()
	if err != nil {
		return stubs, err
	}
	for _, record := range records {
		stubs = append(stubs, &api.Stub{ID: record.ID, Name: record.Name})
	}
	return stubs, nil
}

func (me *service) Get(ctx context.Context, id string, v *demo.Demo) error {
	mu.Lock()
	defer mu.Unlock()
	records, err := getRecords()
	if err != nil {
		return err
	}
	for _, record := range records {
		if record.ID == id {
			v.ID = record.ID
			v.Name = record.Name
			v.Value = record.Value
			return nil
		}
	}
	return fmt.Errorf("'%s' not found", id)
}

func (me *service) SchemaID() string {
	return SchemaID
}

func (me *service) Create(ctx context.Context, v *demo.Demo) (*api.Stub, error) {
	mu.Lock()
	defer mu.Unlock()
	records, err := getRecords()
	if err != nil {
		return nil, err
	}
	newRecord := demo.Demo{
		ID:    uuid.NewString(),
		Name:  v.Name,
		Value: v.Value,
	}
	records = append(records, &newRecord)
	err = storeRecords(records)
	if err != nil {
		return nil, err
	}
	return &api.Stub{ID: newRecord.ID, Name: newRecord.Name}, nil
}

func (me *service) Update(ctx context.Context, id string, v *demo.Demo) error {
	mu.Lock()
	defer mu.Unlock()
	records, err := getRecords()
	if err != nil {
		return err
	}
	for _, record := range records {
		if record.ID == id {
			record.Name = v.Name
			record.Value = v.Value
			if err := storeRecords(records); err != nil {
				return err
			}
			return nil
		}
	}
	return fmt.Errorf("'%s' not found", id)
}

func (me *service) Delete(ctx context.Context, id string) error {
	mu.Lock()
	defer mu.Unlock()
	newRecords := []*demo.Demo{}
	records, err := getRecords()
	if err != nil {
		return err
	}
	for _, record := range records {
		if record.ID != id {
			newRecords = append(newRecords, record)
		}
	}

	return storeRecords(newRecords)
}
