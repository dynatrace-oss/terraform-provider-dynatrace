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

package vault

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api"
	vault "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/platform/credential-vault/v1/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings"

	"github.com/dynatrace/dynatrace-configuration-as-code-core/clients/credentialvault"
)

func Service(clientSet rest.ClientSet) (settings.CRUDService[*vault.Credentials], error) {
	platformClient, err := clientSet.PlatformClient()
	if err != nil {
		return nil, err
	}

	return &service{client: credentialvault.NewClient(platformClient), envURL: clientSet.Credentials().Platform.EnvironmentURL}, nil
}

type service struct {
	client *credentialvault.Client
	envURL string
}

type entryReference struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func (me *service) SchemaID() string {
	return "platform:credential-vault"
}

func (me *service) Create(ctx context.Context, v *vault.Credentials) (*api.Stub, error) {
	if v.ExternalVault != nil && v.ExternalVault.Certificate != nil {
		v.ExternalVault.Certificate = new(strings.Replace(*v.ExternalVault.Certificate, "CREDENTIAL_VAULT_ENTRY", "CREDENTIALS_VAULT", -1))
	}

	data, err := json.Marshal(v)

	if err != nil {
		return nil, err
	}

	raw, err := me.client.Create(ctx, data)
	if err != nil {
		return nil, err
	}

	var response entryReference
	if err = json.Unmarshal(raw.Data, &response); err != nil {
		return nil, err
	}

	return &api.Stub{Name: v.Name, ID: response.ID}, nil
}

func (me *service) Update(ctx context.Context, id string, v *vault.Credentials) error {
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}

	if _, err = me.client.Update(ctx, id, data); err != nil {
		return err
	}

	return nil
}

func (me *service) Get(ctx context.Context, id string, v *vault.Credentials) error {
	response, err := me.client.Get(ctx, id)
	if err != nil {
		return err
	}
	if err = json.Unmarshal(response.Data, v); err != nil {
		return err
	}

	return nil
}

func (me *service) List(ctx context.Context) (api.Stubs, error) {
	result, err := me.client.List(ctx)
	if err != nil {
		return nil, err
	}
	var stubs api.Stubs
	for _, raw := range result.All() {
		var entry entryReference
		if err := json.Unmarshal(raw, &entry); err != nil {
			return nil, err
		}
		stubs = append(stubs, &api.Stub{ID: entry.ID, Name: entry.Name})
	}
	return stubs, nil
}

func (me *service) Delete(ctx context.Context, id string) error {
	_, err := me.client.Delete(ctx, id)

	// add synth monitor cleanup

	return err
}
