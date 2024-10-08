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

package azure

import (
	"context"
	"fmt"
	"sync"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api"
	azure "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/credentials/azure/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings/services/httpcache"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/provider/logging"
)

const SchemaID = "v1:config:credentials:azure"
const BasePath = "/api/config/v1/azure/credentials"

var mu sync.Mutex

func Service(credentials *settings.Credentials) settings.CRUDService[*azure.AzureCredentials] {
	return &service{service: settings.NewCRUDService(
		credentials,
		SchemaID,
		settings.DefaultServiceOptions[*azure.AzureCredentials](BasePath).
			WithMutex(mu.Lock, mu.Unlock),
	), client: httpcache.DefaultClient(credentials.URL, credentials.Token, SchemaID)}
}

type service struct {
	service settings.CRUDService[*azure.AzureCredentials]
	client  rest.Client
}

func (me *service) List(ctx context.Context) (api.Stubs, error) {
	return me.service.List(ctx)
}

func (me *service) Get(ctx context.Context, id string, v *azure.AzureCredentials) error {
	return me.service.Get(ctx, id, v)
}

func (me *service) SchemaID() string {
	return me.service.SchemaID()
}

func (me *service) Create(ctx context.Context, v *azure.AzureCredentials) (*api.Stub, error) {
	stub, err := me.service.Create(ctx, v)
	if err != nil {
		return nil, err
	}
	logging.File.Println("RemoveDefault:", v.RemoveDefaults)
	if v.RemoveDefaults {
		if err := me.client.Put(ctx, fmt.Sprintf("%s/%s/services", BasePath, stub.ID), struct {
			Services []string `json:"services"`
		}{Services: []string{}}, 204).Finish(); err != nil {
			me.Delete(ctx, stub.ID)
			return nil, err
		}
	}
	return stub, err
}

func (me *service) Update(ctx context.Context, id string, v *azure.AzureCredentials) error {
	// by not sending supported services they won't get changed
	v.SupportingServices = nil
	return me.service.Update(ctx, id, v)
}

func (me *service) Delete(ctx context.Context, id string) error {
	return me.service.Delete(ctx, id)
}
