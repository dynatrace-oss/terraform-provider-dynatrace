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

package customdevice

import (
	"fmt"
	"net/url"
	"sync"
	"time"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api"
	customdevice "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v2/customdevice/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings"
	"github.com/google/uuid"
)

var mutex sync.Mutex

func Service(credentials *settings.Credentials) settings.CRUDService[*customdevice.CustomDevice] {
	return &service{credentials}
}

type service struct {
	credentials *settings.Credentials
}

func (me *service) Get(id string, v *customdevice.CustomDevice) error {
	var err error
	client := rest.DefaultClient(me.credentials.URL, me.credentials.Token)
	entitySelector := `detectedName("` + id + `"),type("CUSTOM_DEVICE")`
	var CustomDeviceGetResponse customdevice.CustomDeviceGetResponse

	// The result from the GET API enpoint is not very stable, so attepting to get the custom device once is not enough.
	// 20 is an arbitraty number (it takes 40s before the method gives up) that should be long enough for the endpoint to return a value.
	for i := 0; i < 20; i++ {
		req := client.Get(fmt.Sprintf("/api/v2/entities?from=now-3y&entitySelector=%s", url.QueryEscape(entitySelector))).Expect(200)
		err = req.Finish(&CustomDeviceGetResponse)
		if len(CustomDeviceGetResponse.Entities) != 0 {
			break
		}
		time.Sleep(2 * time.Second)
	}

	if len(CustomDeviceGetResponse.Entities) == 0 {
		// We only throw this error if the Finish method failed for the last attempt because sometimes random calls fail.
		// This way if all calls fail, the last will fail as well, and we only get a false positive if the last call happens to be the only one to fail.
		if err != nil {
			return err
		}
		return rest.Error{Code: 404, Message: `Custom device with ID:` + id + " not found!"}
	}

	v.DisplayName = CustomDeviceGetResponse.Entities[0].DisplayName
	v.EntityId = CustomDeviceGetResponse.Entities[0].EntityId
	v.CustomDeviceID = id

	return nil
}

func (me *service) CheckGet(id string, v *customdevice.CustomDevice) error {
	var err error
	client := rest.DefaultClient(me.credentials.URL, me.credentials.Token)
	entitySelector := `detectedName("` + id + `"),type("CUSTOM_DEVICE")`
	req := client.Get(fmt.Sprintf("/api/v2/entities?from=now-3y&entitySelector=%s", url.QueryEscape(entitySelector))).Expect(200)
	var CustomDeviceGetResponse customdevice.CustomDeviceGetResponse
	if err = req.Finish(&CustomDeviceGetResponse); err != nil {
		return err
	}
	if len(CustomDeviceGetResponse.Entities) == 0 {
		return nil
	}
	v.EntityId = CustomDeviceGetResponse.Entities[0].EntityId
	return nil
}

func (me *service) SchemaID() string {
	return "v2:environment:custom-device"
}

func (me *service) List() (api.Stubs, error) {
	var err error
	client := rest.DefaultClient(me.credentials.URL, me.credentials.Token)
	entitySelector := `type("CUSTOM_DEVICE")`
	req := client.Get(fmt.Sprintf("/api/v2/entities?from=now-3y&entitySelector=%s&fields=properties.detectedName&pageSize=500", url.QueryEscape(entitySelector))).Expect(200)
	listResponse := struct {
		Entities []struct {
			EntityId    string `json:"entityId"`
			Type        string `json:"type"`
			DisplayName string `json:"displayName"`
			Properties  struct {
				DetectedName string `json:"detectedName"`
			} `json:"properties"`
		} `json:"entities"`
	}{}
	if err = req.Finish(&listResponse); err != nil {
		return nil, err
	}
	var stubs api.Stubs
	if len(listResponse.Entities) == 0 {
		return stubs, nil
	}
	for _, entity := range listResponse.Entities {
		if entity.Type != "CUSTOM_DEVICE" {
			continue
		}
		if len(entity.Properties.DetectedName) == 0 {
			continue
		}
		if len(entity.DisplayName) == 0 {
			continue
		}
		stubs = append(stubs, &api.Stub{ID: entity.Properties.DetectedName, Name: entity.DisplayName})
	}

	return stubs, nil
}

func (me *service) Validate(v *customdevice.CustomDevice) error {
	return nil // no endpoint for that
}

func (me *service) Create(v *customdevice.CustomDevice) (*api.Stub, error) {
	mutex.Lock()
	defer mutex.Unlock()
	var err error
	if v.CustomDeviceID == "" {
		v.CustomDeviceID = uuid.NewString()
	}
	client := rest.DefaultClient(me.credentials.URL, me.credentials.Token)
	if err = client.Post("/api/v2/entities/custom", v, 201, 204).Finish(); err != nil {
		return nil, err
	}

	// Check the custom device was indeed created before finishing up
	for i := 0; i < 50; i++ {
		me.CheckGet(v.CustomDeviceID, v)
		time.Sleep(2 * time.Second)
		if v.EntityId != "" {
			break
		}
	}
	return &api.Stub{ID: v.CustomDeviceID, Name: *v.DisplayName}, nil
}

func (me *service) Update(id string, v *customdevice.CustomDevice) error {
	var err error
	v.CustomDeviceID = id
	v.EntityId = ""
	client := rest.DefaultClient(me.credentials.URL, me.credentials.Token)
	if err = client.Post("/api/v2/entities/custom", v, 204).Finish(); err != nil {
		return err
	}
	return nil
}

func (me *service) Delete(id string) error {
	return nil // no endpoint for that
}

func (me *service) New() *customdevice.CustomDevice {
	return new(customdevice.CustomDevice)
}

func (me *service) Name() string {
	return me.SchemaID()
}
