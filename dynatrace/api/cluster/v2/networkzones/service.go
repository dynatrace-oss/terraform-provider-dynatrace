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

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api"
	networkzones "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/cluster/v2/networkzones/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings"
	"github.com/google/uuid"
)

const SchemaID = "cluster:networkzones"

func Service(credentials *settings.Credentials) settings.CRUDService[*networkzones.NetworkZone] {
	return &service{
		serviceClient: NewService(fmt.Sprintf("%s%s", credentials.Cluster.URL, "/api/cluster/v2"), credentials.Cluster.Token),
	}
}

type ServiceClient struct {
	client rest.Client
}

func (me *service) Create(v *networkzones.NetworkZone) (*api.Stub, error) {
	return me.serviceClient.Create(v)
}

func (me *service) Update(id string, v *networkzones.NetworkZone) error {
	return me.serviceClient.Update(id, v)
}

func (me *service) Delete(id string) error {
	return me.serviceClient.Delete(id)
}

func (me *service) List() (api.Stubs, error) {
	return me.serviceClient.List()
}

func (me *service) Get(id string, v *networkzones.NetworkZone) error {
	return me.serviceClient.Get(id, v)
}

func (me *service) SchemaID() string {
	return SchemaID
}

func (cs *ServiceClient) SchemaID() string {
	return SchemaID
}

func NewService(baseURL string, token string) *ServiceClient {
	return &ServiceClient{client: rest.DefaultClient(baseURL, token)}
}

type service struct {
	serviceClient *ServiceClient
}

func (cs *ServiceClient) Create(v *networkzones.NetworkZone) (*api.Stub, error) {
	var id string
	if v.NetworkZoneName != nil {
		id = *v.NetworkZoneName
	} else {
		id = uuid.NewString()
	}

	response := networkzones.NetworkZone{}
	if err := cs.client.Put(fmt.Sprintf("/networkZones/%s", url.PathEscape(id)), v, 201, 204).Finish(&response); err != nil {
		return nil, err
	}

	return &api.Stub{ID: id, Name: id}, nil
}

func (cs *ServiceClient) Update(id string, v *networkzones.NetworkZone) error {
	return cs.client.Put(fmt.Sprintf("/networkZones/%s", id), v, 204).Finish()
}

func (cs *ServiceClient) Delete(id string) error {
	return cs.client.Delete(fmt.Sprintf("/networkZones/%s", url.PathEscape(id))).Expect(204).Finish()
}

func (cs *ServiceClient) Get(id string, v *networkzones.NetworkZone) error {
	return cs.client.Get(fmt.Sprintf("/networkZones/%s", url.PathEscape(id)), 200).Finish(v)
}

func (cs *ServiceClient) List() (api.Stubs, error) {
	var err error
	var stubList networkzones.NetworkZones

	if err = cs.client.Get("/networkZones", 200).Finish(&stubList); err != nil {
		return nil, err
	}
	stubs := api.Stubs{}
	for _, zone := range stubList.Zones {
		stubs = append(stubs, &api.Stub{ID: *zone.ID, Name: *zone.ID})
	}
	return stubs, nil
}
