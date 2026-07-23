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

package slo

import (
	"context"
	"encoding/json"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings"
	sloclient "github.com/dynatrace/dynatrace-configuration-as-code-core/clients/slo"

	slo "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/slo/settings"
)

func Service(clientSet rest.ClientSet) (settings.CRUDService[*slo.SLO], error) {
	platformClient, err := clientSet.PlatformClient()
	if err != nil {
		return nil, err
	}

	return &service{client: sloclient.NewClient(platformClient)}, nil
}

type service struct {
	client *sloclient.Client
}

func (me *service) Get(ctx context.Context, id string, v *slo.SLO) error {
	response, err := me.client.Get(ctx, id)
	if err != nil {
		return err
	}

	return json.Unmarshal(response.Data, &v)
}

func (me *service) SchemaID() string {
	return "platform:slo"
}

func (me *service) List(ctx context.Context) (api.Stubs, error) {
	listResponse, err := me.client.List(ctx)
	if err != nil {
		return nil, err
	}

	var stubs api.Stubs
	for _, r := range listResponse.All() {
		var stub api.Stub
		if err := json.Unmarshal(r, &stub); err != nil {
			return nil, err
		}
		stubs = append(stubs, &api.Stub{ID: stub.ID, Name: stub.Name})
	}

	return stubs, nil
}

func (me *service) Validate(_ *slo.SLO) error {
	return nil // no endpoint for that
}

func (me *service) Create(ctx context.Context, v *slo.SLO) (*api.Stub, error) {
	data, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}

	response, err := me.client.Create(ctx, data)
	if err != nil {
		return nil, err
	}

	var stub api.Stub
	if err := json.Unmarshal(response.Data, &stub); err != nil {
		return nil, err
	}

	return &stub, nil
}

func (me *service) Update(ctx context.Context, id string, v *slo.SLO) error {
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}

	_, err = me.client.Update(ctx, id, data)
	if err != nil {
		return err
	}

	return nil
}

func (me *service) Delete(ctx context.Context, id string) error {
	_, err := me.client.Delete(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

func (me *service) New() *slo.SLO {
	return new(slo.SLO)
}
