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

package locations

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings"

	locations "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/synthetic/locations/settings"
)

const SchemaID = "v1:synthetic:locations:all"

func Service(credentials *settings.Credentials) settings.RService[*locations.SyntheticLocation] {
	return &service{client: rest.DefaultClient(credentials.URL, credentials.Token)}
}

type service struct {
	client rest.Client
}

func (me *service) List() (stubs settings.Stubs, err error) {
	var stubList locations.SyntheticLocations
	if err = me.client.Get("/api/v1/synthetic/locations", 200).Finish(&stubList); err != nil {
		return nil, err
	}
	for _, location := range stubList.Locations {
		localLocation := location
		stubs = append(stubs, &settings.Stub{ID: location.ID, Name: location.Name, Value: localLocation})
	}
	return stubs, nil
}

func (me *service) Get(id string, v *locations.SyntheticLocation) (err error) {
	var stubs settings.Stubs
	if stubs, err = me.List(); err != nil {
		return err
	}
	for _, stub := range stubs {
		if stub.ID == id {
			value := stub.Value.(*locations.SyntheticLocation)
			v.ID = value.ID
			v.CloudPlatform = value.CloudPlatform
			v.IPs = value.IPs
			v.Name = value.Name
			v.Type = value.Type
			v.Status = value.Status
			v.Stage = value.Stage
			return nil
		}
	}
	return rest.Error{Code: 404, Message: "No synthetic location with id " + id + "found"}
}

func (me *service) SchemaID() string {
	return SchemaID
}
