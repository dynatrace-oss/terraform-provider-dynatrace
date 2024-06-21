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

package activegatetokens

import (
	"context"
	"fmt"
	"net/url"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings"

	activegatetokens "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v2/activegatetokens/settings"
)

func Service(credentials *settings.Credentials) settings.CRUDService[*activegatetokens.Settings] {
	return &service{credentials}
}

type TokenCreateResponse struct {
	ID             string  `json:"id"`
	ExpirationDate *string `json:"expirationDate,omitempty"`
	Token          *string `json:"token,omitempty"`
}

type TenantTokenResponse struct {
	TenantToken string `json:"tenantToken"`
}

type service struct {
	credentials *settings.Credentials
}

func (me *service) Get(ctx context.Context, id string, v *activegatetokens.Settings) error {
	var err error

	client := rest.DefaultClient(me.credentials.URL, me.credentials.Token)
	req := client.Get(fmt.Sprintf("/api/v2/activeGateTokens/%s", url.PathEscape(id))).Expect(200)
	if err = req.Finish(v); err != nil {
		return err
	}
	var ttr TenantTokenResponse
	req = client.Get("/api/v1/deployment/installer/agent/connectioninfo").Expect(200)
	if err = req.Finish(&ttr); err == nil {
		v.TenantToken = &ttr.TenantToken
	}

	return nil
}

func (me *service) SchemaID() string {
	return "v2:environment:activegate-tokens"
}

func (me *service) List(ctx context.Context) (api.Stubs, error) {
	return api.Stubs{}, nil
}

func (me *service) Validate(v *activegatetokens.Settings) error {
	return nil // no endpoint for that
}

func (me *service) Create(ctx context.Context, v *activegatetokens.Settings) (*api.Stub, error) {
	var err error

	response := TokenCreateResponse{}
	client := rest.DefaultClient(me.credentials.URL, me.credentials.Token)
	if err = client.Post("/api/v2/activeGateTokens", v, 201).Finish(&response); err != nil {
		return nil, err
	}
	v.ExpirationDate = response.ExpirationDate
	v.Token = response.Token

	var ttr TenantTokenResponse
	req := client.Get("/api/v1/deployment/installer/agent/connectioninfo").Expect(200)
	if err = req.Finish(&ttr); err == nil {
		v.TenantToken = &ttr.TenantToken
	}

	return &api.Stub{ID: response.ID, Name: v.Name, Value: v}, nil
}

func (me *service) Update(ctx context.Context, id string, v *activegatetokens.Settings) error {
	return nil
}

func (me *service) Delete(ctx context.Context, id string) error {
	return rest.DefaultClient(me.credentials.URL, me.credentials.Token).Delete(fmt.Sprintf("/api/v2/activeGateTokens/%s", url.PathEscape(id)), 204).Finish()
}

func (me *service) New() *activegatetokens.Settings {
	return new(activegatetokens.Settings)
}
