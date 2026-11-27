/**
* @license
* Copyright 2026 Dynatrace LLC
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

package location

import (
	"context"
	"encoding/json"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api"
	location "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/synthetic/location/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings"

	"github.com/dynatrace/dynatrace-configuration-as-code-core/clients/synthetic"
)

type service struct {
	client *synthetic.Client
}

func Service(clientSet rest.ClientSet) (settings.CRUDService[*location.PrivateSyntheticLocation], error) {
	client, err := clientSet.PlatformClient()
	if err != nil {
		return nil, err
	}
	return &service{synthetic.NewClient(client)}, nil
}

func (me *service) Get(ctx context.Context, id string, v *location.PrivateSyntheticLocation) error {
	response, err := me.client.GetLocation(ctx, id)
	if err != nil {
		return err
	}

	return json.Unmarshal(response.Data, &v)
}

func (me *service) SchemaID() string {
	return "platform:synthetic:location"
}

func (me *service) List(ctx context.Context) (api.Stubs, error) {
	listResponse, err := me.client.ListLocations(ctx, synthetic.ListLocationsOptions{
		Type: "PRIVATE",
	})
	if err != nil {
		return nil, err
	}

	var stubs api.Stubs
	for _, r := range listResponse.Objects {
		var stub api.Stub
		if err := json.Unmarshal(r, &stub); err != nil {
			return nil, err
		}
		stubs = append(stubs, &api.Stub{ID: stub.ID, Name: stub.Name})
	}

	return stubs, nil
}

func (me *service) Create(ctx context.Context, v *location.PrivateSyntheticLocation) (*api.Stub, error) {
	data, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}

	response, err := me.client.CreateLocation(ctx, data)
	if err != nil {
		return nil, err
	}

	var stub api.Stub
	if err := json.Unmarshal(response.Data, &stub); err != nil {
		return nil, err
	}
	stub.Name = v.Name

	return &stub, nil
}

func (me *service) Update(ctx context.Context, id string, v *location.PrivateSyntheticLocation) error {
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}

	_, err = me.client.UpdateLocation(ctx, id, data)
	return err
}

func (me *service) Delete(ctx context.Context, id string) error {
	_, err := me.client.DeleteLocation(ctx, id)
	return err
}
