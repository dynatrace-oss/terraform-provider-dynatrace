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

package configuration

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings"

	"github.com/dynatrace/dynatrace-configuration-as-code-core/clients/synthetic"

	configuration "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/platform/synthetic/configuration/settings"
)

// ID is the fixed ID of the singleton configuration.
const ID = "dynatrace_synthetic_v2_configuration"

func Service(clientSet rest.ClientSet) (settings.CRUDService[*configuration.Settings], error) {
	platformClient, err := clientSet.PlatformClient()
	if err != nil {
		return nil, fmt.Errorf("dynatrace_synthetic_v2_configuration requires OAuth client credentials (DT_CLIENT_ID, DT_CLIENT_SECRET, DT_ACCOUNT_ID) or a platform token (DYNATRACE_PLATFORM_TOKEN) with scopes synthetic:environment-configuration:read and synthetic:environment-configuration:write: %w", err)
	}

	return &service{client: synthetic.NewClient(platformClient)}, nil
}

type service struct {
	client *synthetic.Client
}

func (me *service) SchemaID() string {
	return "platform:synthetic:configuration"
}

func (me *service) List(ctx context.Context) (api.Stubs, error) {
	return api.Stubs{&api.Stub{ID: ID, Name: ID}}, nil
}

// Get ignores the ID, as there is only one configuration per environment.
func (me *service) Get(ctx context.Context, id string, v *configuration.Settings) error {
	response, err := me.client.GetConfiguration(ctx)
	if err != nil {
		return err
	}
	return json.Unmarshal(response.Data, v)
}

func (me *service) Create(ctx context.Context, v *configuration.Settings) (*api.Stub, error) {
	if err := me.put(ctx, v); err != nil {
		return nil, err
	}
	return &api.Stub{ID: ID, Name: ID}, nil
}

func (me *service) Update(ctx context.Context, id string, v *configuration.Settings) error {
	return me.put(ctx, v)
}

// put reads the current configuration and writes it back with the fields set in v overlaid,
// so that unset fields keep the value of the tenant.
func (me *service) put(ctx context.Context, v *configuration.Settings) error {
	response, err := me.client.GetConfiguration(ctx)
	if err != nil {
		return err
	}

	var current map[string]json.RawMessage
	if err := json.Unmarshal(response.Data, &current); err != nil {
		return err
	}

	// v only marshals the fields that are set (omitempty on pointers), so unset fields are not overlaid.
	desired, err := json.Marshal(v)
	if err != nil {
		return err
	}
	var overlay map[string]json.RawMessage
	if err := json.Unmarshal(desired, &overlay); err != nil {
		return err
	}
	for key, value := range overlay {
		current[key] = value
	}

	body, err := json.Marshal(current)
	if err != nil {
		return err
	}

	_, err = me.client.UpdateConfiguration(ctx, body)
	return err
}

// Delete only removes the resource from the state.
func (me *service) Delete(ctx context.Context, id string) error {
	return nil
}

func (me *service) New() *configuration.Settings {
	return new(configuration.Settings)
}
