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

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api"
	customdevice "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v2/customdevice/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings"
	"github.com/google/uuid"
)

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
	req := client.Get(fmt.Sprintf("/api/v2/entities?from=now-3y&&entitySelector=%s", entitySelector)).Expect(200)
	var enitityList customdevice.CustomDeviceList
	if err = req.Finish(enitityList); err != nil {
		return err
	}

	if len(enitityList.Entities) == 0 {
		return rest.Error{Code: 404, Message: `Custom device with ID:` + id + " not found!"}
	}

	v.DisplayName = enitityList.Entities[0].DisplayName
	v.EntityId = enitityList.Entities[0].EntityId
	v.CustomDeviceID = id

	return nil
}

func (me *service) SchemaID() string {
	return "v2:environment:custom-device"
}

func (me *service) List() (api.Stubs, error) {
	return api.Stubs{}, nil
}

func (me *service) Validate(v *customdevice.CustomDevice) error {
	return nil // no endpoint for that
}

func (me *service) Create(v *customdevice.CustomDevice) (*api.Stub, error) {
	var err error
	if v.CustomDeviceID == "" {
		v.CustomDeviceID = uuid.NewString()
	}
	resultDevice := customdevice.CustomDevice{}
	client := rest.DefaultClient(me.credentials.URL, me.credentials.Token)
	if err = client.Post("/api/v2/entities/custom", v, 201, 204).Finish(&resultDevice); err != nil {
		return nil, err
	}
	resultDevice.CustomDeviceID = v.CustomDeviceID
	resultDevice.DisplayName = v.DisplayName

	return &api.Stub{ID: resultDevice.CustomDeviceID, Name: *resultDevice.DisplayName, Value: resultDevice}, nil
}

func (me *service) Update(id string, v *customdevice.CustomDevice) error {
	var err error
	v.CustomDeviceID = id
	v.EntityId = nil
	resultDevice := customdevice.CustomDevice{}
	client := rest.DefaultClient(me.credentials.URL, me.credentials.Token)
	if err = client.Post("/api/v2/entities/custom", v, 201, 204).Finish(&resultDevice); err != nil {
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
