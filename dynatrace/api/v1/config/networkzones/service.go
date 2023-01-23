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

package networkzones

import (
	"fmt"
	"net/url"
	"strings"

	enable "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v2/networkzones"
	enablesettings "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v2/networkzones/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings"
	"github.com/google/uuid"

	networkzones "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/networkzones/settings"
)

const SchemaID = "v2:environment:network-zones"

func Service(credentials *settings.Credentials) settings.CRUDService[*networkzones.NetworkZone] {
	return &service{client: rest.DefaultClient(credentials.URL, credentials.Token), credentials: credentials}
}

type service struct {
	credentials *settings.Credentials
	client      rest.Client
}

func (me *service) Get(id string, v *networkzones.NetworkZone) error {
	return me.client.Get(fmt.Sprintf("/api/v2/networkZones/%s", url.PathEscape(id)), 200).Finish(v)
}

func (me *service) List() (settings.Stubs, error) {
	var err error
	var stubList networkzones.NetworkZones

	if err = me.client.Get("/api/v2/networkZones", 200).Finish(&stubList); err != nil {
		return nil, err
	}
	stubs := settings.Stubs{}
	for _, zone := range stubList.Zones {
		stubs = append(stubs, &settings.Stub{ID: *zone.ID, Name: *zone.ID})
	}
	return stubs, nil
}

func (me *service) Validate(v *networkzones.NetworkZone) error {
	return nil // no endpoint for that
}

func (me *service) Create(v *networkzones.NetworkZone) (*settings.Stub, error) {
	var err error

	// id := *v.ID
	id := uuid.NewString()

	var stub settings.Stub
	req := me.client.Put(fmt.Sprintf("/api/v2/networkZones/%s", url.PathEscape(id)), v, 201)
	if err = req.Finish(&stub); err != nil {
		if strings.Contains(err.Error(), "Not allowed because network zones are disabled") {
			if _, err := enable.Service(me.credentials).Create(&enablesettings.NetworkZones{Enabled: true}); err != nil {
				return nil, err
			}
			return me.Create(v)
		}
		if strings.Contains(err.Error(), "Creation and modification of network zone is only possible via cluster API.") {
			name := uuid.NewString()
			return &settings.Stub{ID: name + "---flawed----", Name: name}, nil
		}
		return nil, err
	}
	return &stub, nil
}

func (me *service) Update(id string, v *networkzones.NetworkZone) error {
	if err := me.client.Put(fmt.Sprintf("/api/v2/networkZones/%s", url.PathEscape(id)), v, 204).Finish(); err != nil {
		return err
	}
	return nil
}

func (me *service) Delete(id string) error {
	return me.client.Delete(fmt.Sprintf("/api/v2/networkZones/%s", url.PathEscape(id))).Expect(204).Finish()
}

func (me *service) SchemaID() string {
	return SchemaID
}
