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
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api"
	enable "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/networkzones"
	enablesettings "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/networkzones/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings"
	"github.com/google/uuid"

	networkzones "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v2/networkzones/settings"
)

const SchemaID = "v2:environment:network-zones"

func Service(credentials *rest.Credentials) settings.CRUDService[*networkzones.NetworkZone] {
	return &service{client: rest.APITokenClient(credentials), credentials: credentials}
}

type service struct {
	credentials *rest.Credentials
	client      rest.Client
}

func (me *service) Get(ctx context.Context, id string, v *networkzones.NetworkZone) error {
	return me.client.Get(ctx, fmt.Sprintf("/api/v2/networkZones/%s", url.PathEscape(id)), 200).Finish(v)
}

func (me *service) List(ctx context.Context) (api.Stubs, error) {
	var err error
	var stubList networkzones.NetworkZones

	if err = me.client.Get(ctx, "/api/v2/networkZones", 200).Finish(&stubList); err != nil {
		return nil, err
	}
	stubs := api.Stubs{}
	for _, zone := range stubList.Zones {
		stubs = append(stubs, &api.Stub{ID: *zone.ID, Name: *zone.ID})
	}
	return stubs, nil
}

func (me *service) Validate(v *networkzones.NetworkZone) error {
	return nil // no endpoint for that
}

func (me *service) Create(ctx context.Context, v *networkzones.NetworkZone) (*api.Stub, error) {
	var err error

	// id := *v.ID
	id := uuid.NewString()
	if v.NetworkZoneName != nil {
		id = *v.NetworkZoneName
	}

	req := me.client.Put(ctx, fmt.Sprintf("/api/v2/networkZones/%s", url.PathEscape(id)), v, 201, 204)

	if err = req.Finish(); err != nil {
		if strings.Contains(err.Error(), "Not allowed because network zones are disabled") {
			if _, err := enable.Service(me.credentials).Create(ctx, &enablesettings.NetworkZones{Enabled: true}); err != nil {
				return nil, err
			}
			return me.Create(ctx, v)
		}
		if strings.Contains(err.Error(), "Creation and modification of network zone is only possible via cluster API.") {
			return &api.Stub{ID: id + "---flawed----", Name: id}, nil
		}
		return nil, err
	}

	return &api.Stub{ID: id, Name: id}, nil
}

func (me *service) Update(ctx context.Context, id string, v *networkzones.NetworkZone) error {
	if v.NetworkZoneName != nil && id != strings.ToLower(*v.NetworkZoneName) {
		return fmt.Errorf("Network zone name cannot be modified, please destroy and create with the new name")
	}
	if err := me.client.Put(ctx, fmt.Sprintf("/api/v2/networkZones/%s", url.PathEscape(id)), v, 204).Finish(); err != nil {
		return err
	}
	return nil
}

func (me *service) Delete(ctx context.Context, id string) error {
	return me.client.Delete(ctx, fmt.Sprintf("/api/v2/networkZones/%s", url.PathEscape(id))).Expect(204).Finish()
}

func (me *service) SchemaID() string {
	return SchemaID
}
