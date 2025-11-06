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

package sitereliabilityguardian

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api"
	sitereliabilityguardian "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/app/dynatrace/sitereliabilityguardian/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings/services/settings20"
)

const SchemaVersion = "1.8"
const SchemaID = "app:dynatrace.site.reliability.guardian:guardians"

// func Service(credentials *rest.Credentials) settings.CRUDService[*sitereliabilityguardian.Settings] {
// 	return settings20.Service[*sitereliabilityguardian.Settings](credentials, SchemaID, SchemaVersion)
// }

func Service(credentials *rest.Credentials) settings.CRUDService[*sitereliabilityguardian.Settings] {
	return &service{credentials: credentials}
}

type service struct {
	credentials *rest.Credentials
}

func (me *service) Client(ctx context.Context, schemaIDs string) *settings20.Client {
	tokenClient, _ := rest.CreateClassicClient(me.credentials.URL, me.credentials.Token)
	oauthClient, _ := rest.CreateClassicOAuthBasedClient(ctx, me.credentials)
	return settings20.NewClient(tokenClient, oauthClient, schemaIDs)
}

func (me *service) Create(ctx context.Context, v *sitereliabilityguardian.Settings) (*api.Stub, error) {

	data, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}

	scope := "environment"
	response, err := me.Client(ctx, SchemaID).Create(ctx, scope, data)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != 200 {
		if err := rest.Envelope(response.Data, response.Request.URL, response.Request.Method); err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("status code %d (expected: %d): %s", response.StatusCode, 200, string(response.Data))
	}

	stub := &api.Stub{ID: response.ID, Name: response.ID}
	return stub, nil
}

func (me *service) Update(ctx context.Context, id string, v *sitereliabilityguardian.Settings) error {
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}

	response, err := me.Client(ctx, "").Update(ctx, id, data)
	if err != nil {
		return err
	}

	if response.StatusCode != 200 {
		if err := rest.Envelope(response.Data, response.Request.URL, response.Request.Method); err != nil {
			return err
		}
		return fmt.Errorf("status code %d (expected: %d): %s", response.StatusCode, 200, string(response.Data))
	}

	return nil
}

func (me *service) Validate(v *sitereliabilityguardian.Settings) error {
	return nil // Settings 2.0 doesn't offer validation
}

func (me *service) Delete(ctx context.Context, id string) error {
	response, err := me.Client(ctx, "").Delete(ctx, id)
	if err != nil {
		return err
	}

	if response.StatusCode != 204 {
		if err = rest.Envelope(response.Data, response.Request.URL, response.Request.Method); err != nil {
			return err
		}
		return fmt.Errorf("status code %d (expected: %d): %s", response.StatusCode, 204, string(response.Data))
	}

	return nil
}

type SettingsObject struct {
	SchemaVersion string          `json:"schemaVersion"`
	SchemaID      string          `json:"schemaId"`
	Scope         string          `json:"scope"`
	Value         json.RawMessage `json:"value"`
}

func (me *service) Get(ctx context.Context, id string, v *sitereliabilityguardian.Settings) error {

	response, err := me.Client(ctx, "").Get(ctx, id)
	if err != nil {
		return err
	}

	if response.StatusCode != 200 {
		if err := rest.Envelope(response.Data, response.Request.URL, response.Request.Method); err != nil {
			return err
		}
		return fmt.Errorf("status code %d (expected: %d): %s", response.StatusCode, 200, string(response.Data))
	}

	var settingsObject SettingsObject
	if err := json.Unmarshal(response.Data, &settingsObject); err != nil {
		return err
	}
	if err := json.Unmarshal(settingsObject.Value, &v); err != nil {
		return err
	}

	return nil
}

func (me *service) List(ctx context.Context) (api.Stubs, error) {
	response, err := me.Client(ctx, SchemaID).List(ctx)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != 200 {
		if err := rest.Envelope(response.Data, response.Request.URL, response.Request.Method); err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("status code %d (expected: %d): %s", response.StatusCode, 200, string(response.Data))
	}

	var stubs api.Stubs
	for _, item := range response.Items {
		stubs = append(stubs, &api.Stub{ID: item.ID, Name: item.ID})

	}
	return stubs, nil
}

func (me *service) SchemaID() string {
	return SchemaID
}
