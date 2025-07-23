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

package scheduling_rules

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api"
	automationerr "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/automation"
	scheduling_rules "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/automation/scheduling_rules/settings"
	tfrest "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings"

	cacapi "github.com/dynatrace/dynatrace-configuration-as-code-core/api"
	"github.com/dynatrace/dynatrace-configuration-as-code-core/clients/automation"
)

func Service(credentials *tfrest.Credentials) settings.CRUDService[*scheduling_rules.Settings] {
	return &service{credentials}
}

type service struct {
	credentials *tfrest.Credentials
}

func (me *service) client(ctx context.Context) (*automation.Client, error) {
	platformClient, err := tfrest.CreatePlatformClient(ctx, me.credentials.OAuth.EnvironmentURL, me.credentials)
	if err != nil {
		return nil, err
	}
	return automation.NewClient(platformClient), nil
}

func (me *service) Get(ctx context.Context, id string, v *scheduling_rules.Settings) (err error) {
	client, err := me.client(ctx)
	if err != nil {
		return err
	}
	var response cacapi.Response
	if response, err = client.Get(ctx, automation.SchedulingRules, id); err != nil {
		return err
	}
	if response.StatusCode != 200 {
		var e automationerr.ErrorEnvelope
		if e.Unmarshal(response.Data) {
			return e.Err.ToRESTError()
		}
		return tfrest.Error{Code: response.StatusCode, Message: string(response.Data)}
	}
	return json.Unmarshal(response.Data, &v)
}

func (me *service) SchemaID() string {
	return "automation:scheduling.rules"
}

type SchedulingRuleStub struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}

func (me *service) List(ctx context.Context) (api.Stubs, error) {
	client, err := me.client(ctx)
	if err != nil {
		return nil, err
	}
	listResponse, err := client.List(ctx, automation.SchedulingRules)
	if err != nil {
		apiErr := cacapi.APIError{}
		if errors.As(err, &apiErr) {
			return nil, tfrest.Error{Code: apiErr.StatusCode, Message: string(apiErr.Body)}
		}

		return nil, err
	}
	var stubs api.Stubs
	for _, r := range listResponse.All() {
		var stub SchedulingRuleStub
		if err := json.Unmarshal(r, &stub); err != nil {
			return nil, err
		}
		stubs = append(stubs, &api.Stub{ID: stub.ID, Name: stub.Title})
	}
	return stubs, nil
}

func (me *service) Validate(v *scheduling_rules.Settings) error {
	return nil // no endpoint for that
}

func (me *service) Create(ctx context.Context, v *scheduling_rules.Settings) (stub *api.Stub, err error) {
	client, err := me.client(ctx)
	if err != nil {
		return nil, err
	}
	var data []byte
	if data, err = json.Marshal(v); err != nil {
		return nil, err
	}
	var response cacapi.Response
	if response, err = client.Create(ctx, automation.SchedulingRules, data); err != nil {
		return nil, err
	}
	if response.StatusCode == 201 {
		var stub SchedulingRuleStub
		if err := json.Unmarshal(response.Data, &stub); err != nil {
			return nil, err
		}
		return &api.Stub{Name: v.Title, ID: stub.ID}, nil
	}
	var e automationerr.ErrorEnvelope
	if e.Unmarshal(response.Data) {
		return nil, e.Err.ToRESTError()
	}
	return nil, tfrest.Error{Code: response.StatusCode, Message: string(response.Data)}
}

func (me *service) Update(ctx context.Context, id string, v *scheduling_rules.Settings) (err error) {
	client, err := me.client(ctx)
	if err != nil {
		return err
	}
	var data []byte
	if data, err = json.Marshal(v); err != nil {
		return err
	}
	var response cacapi.Response
	if response, err = client.Update(ctx, automation.SchedulingRules, id, data); err != nil {
		return err
	}
	if response.StatusCode == 200 {
		return nil
	}
	var e automationerr.ErrorEnvelope
	if e.Unmarshal(response.Data) {
		return e.Err.ToRESTError()
	}
	return tfrest.Error{Code: response.StatusCode, Message: string(response.Data)}
}

func (me *service) Delete(ctx context.Context, id string) error {
	client, err := me.client(ctx)
	if err != nil {
		return err
	}
	response, err := client.Delete(ctx, automation.SchedulingRules, id)
	if response.StatusCode == 204 || response.StatusCode == 404 {
		return nil
	}
	if response.StatusCode != 0 {
		return tfrest.Error{Code: response.StatusCode, Message: string(response.Data)}
	}
	return err
}

func (me *service) New() *scheduling_rules.Settings {
	return new(scheduling_rules.Settings)
}
