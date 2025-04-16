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
	sloapi "github.com/dynatrace/dynatrace-configuration-as-code-core/api"
	sloclient "github.com/dynatrace/dynatrace-configuration-as-code-core/clients/slo"

	slo "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/slo/settings"
)

func Service(credentials *rest.Credentials) settings.CRUDService[*slo.SLO] {
	return &service{credentials}
}

type service struct {
	credentials *rest.Credentials
}

func (me *service) client(ctx context.Context) (*sloclient.Client, error) {
	platformClient, err := rest.CreatePlatformClient(ctx, me.credentials.OAuth.EnvironmentURL, me.credentials)
	if err != nil {
		return nil, err
	}
	return sloclient.NewClient(platformClient), nil
}

func (me *service) Get(ctx context.Context, id string, v *slo.SLO) (err error) {
	client, err := me.client(ctx)
	if err != nil {
		return err
	}
	response, err := client.Get(ctx, id)
	if err != nil {
		if apiError, ok := err.(sloapi.APIError); ok {
			return rest.Error{Code: apiError.StatusCode, Message: apiError.Error()}
		}
		return err
	}

	return json.Unmarshal(response.Data, &v)
}

func (me *service) SchemaID() string {
	return "platform:slo"
}

func (me *service) List(ctx context.Context) (api.Stubs, error) {
	client, err := me.client(ctx)
	if err != nil {
		return nil, err
	}
	listResponse, err := client.List(ctx)
	if err != nil {
		if apiError, ok := err.(sloapi.APIError); ok {
			return nil, rest.Error{Code: apiError.StatusCode, Message: apiError.Error()}
		}
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
	client, err := me.client(ctx)
	if err != nil {
		return nil, err
	}
	var data []byte
	if data, err = json.Marshal(v); err != nil {
		return nil, err
	}
	response, err := client.Create(ctx, data)
	if err != nil {
		if apiError, ok := err.(sloapi.APIError); ok {
			return nil, rest.Error{Code: apiError.StatusCode, Message: apiError.Error()}
		}
		return nil, err
	}

	var stub api.Stub
	if err := json.Unmarshal(response.Data, &stub); err != nil {
		return nil, err
	}

	return &stub, nil
}

func (me *service) Update(ctx context.Context, id string, v *slo.SLO) (err error) {
	client, err := me.client(ctx)
	if err != nil {
		return err
	}
	var data []byte
	if data, err = json.Marshal(v); err != nil {
		return err
	}

	_, err = client.Update(ctx, id, data)
	if err != nil {
		if apiError, ok := err.(sloapi.APIError); ok {
			return rest.Error{Code: apiError.StatusCode, Message: apiError.Error()}
		}
		return err
	}

	return nil
}

func (me *service) Delete(ctx context.Context, id string) error {
	client, err := me.client(ctx)
	if err != nil {
		return err
	}
	_, err = client.Delete(ctx, id)
	if err != nil {
		if apiError, ok := err.(sloapi.APIError); ok {
			return rest.Error{Code: apiError.StatusCode, Message: apiError.Error()}
		}
		return err
	}

	return nil
}

func (me *service) New() *slo.SLO {
	return new(slo.SLO)
}
