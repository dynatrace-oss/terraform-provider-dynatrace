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

package apitokens

import (
	"fmt"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings"

	apitokens "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v2/apitokens/settings"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/opt"
)

func Service(credentials *settings.Credentials) settings.CRUDService[*apitokens.APIToken] {
	return &service{credentials}
}

type service struct {
	credentials *settings.Credentials
}

func (me *service) Get(id string, v *apitokens.APIToken) error {
	var err error

	client := rest.DefaultClient(me.credentials.URL, me.credentials.Token)
	req := client.Get(fmt.Sprintf("/api/v2/apiTokens/%s", id)).Expect(200)
	if err = req.Finish(v); err != nil {
		return err
	}

	return nil
}

func (me *service) SchemaID() string {
	return "v2:environment:api-tokens"
}

func (me *service) List() (settings.Stubs, error) {
	var err error

	client := rest.DefaultClient(me.credentials.URL, me.credentials.Token)
	req := client.Get("/api/v2/apiTokens?pageSize=10000&sort=-creationDate").Expect(200)
	var tokenlist apitokens.TokenList
	if err = req.Finish(&tokenlist); err != nil {
		return nil, err
	}
	stubs := settings.Stubs{}
	for _, token := range tokenlist.APITokens {
		stubs = append(stubs, &settings.Stub{ID: *token.ID, Name: token.Name})
	}

	return stubs, nil
}

func (me *service) Validate(v *apitokens.APIToken) error {
	return nil // no endpoint for that
}

func (me *service) Create(v *apitokens.APIToken) (*settings.Stub, error) {
	var err error

	resultToken := struct {
		apitokens.APIToken
		ID *string `json:"id,omitempty"`
	}{}

	token := v

	client := rest.DefaultClient(me.credentials.URL, me.credentials.Token)
	if err = client.Post("/api/v2/apiTokens", v, 201).Finish(&resultToken); err != nil {
		return nil, err
	}
	item := v
	if item.Enabled == nil || !*item.Enabled {
		item.Enabled = opt.NewBool(false)
		if err = client.Put(fmt.Sprintf("/api/v2/apiTokens/%s", *resultToken.ID), item, 204).Finish(); err != nil {
			return nil, err
		}
		resultToken.Enabled = item.Enabled
	}
	resultToken.Name = item.Name
	resultToken.Scopes = item.Scopes

	return &settings.Stub{ID: *resultToken.ID, Name: token.Name, Value: resultToken}, nil
}

func (me *service) Update(id string, v *apitokens.APIToken) error {
	return rest.DefaultClient(me.credentials.URL, me.credentials.Token).Put(fmt.Sprintf("/api/v2/apiTokens/%s", id), v, 204).Finish()
}

func (me *service) Delete(id string) error {
	return rest.DefaultClient(me.credentials.URL, me.credentials.Token).Delete(fmt.Sprintf("/api/v2/apiTokens/%s", id), 204).Finish()
}

func (me *service) New() *apitokens.APIToken {
	return new(apitokens.APIToken)
}

func (me *service) Name() string {
	return me.SchemaID()
}
