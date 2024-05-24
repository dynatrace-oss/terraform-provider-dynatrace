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

package business_calendars

import (
	"context"
	"encoding/json"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings"

	business_calendars "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/automation/business_calendars/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/monaco/pkg/client/auth"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/monaco/pkg/client/automation"
)

func Service(credentials *settings.Credentials) settings.CRUDService[*business_calendars.Settings] {
	return &service{credentials}
}

type service struct {
	credentials *settings.Credentials
}

func (me *service) client() *automation.Client {
	httpClient := auth.NewOAuthClient(context.TODO(), auth.OauthCredentials{
		ClientID:     me.credentials.Automation.ClientID,
		ClientSecret: me.credentials.Automation.ClientSecret,
		TokenURL:     me.credentials.Automation.TokenURL,
	})
	return automation.NewClient(me.credentials.Automation.EnvironmentURL, httpClient)
}

func (me *service) Get(id string, v *business_calendars.Settings) (err error) {
	var result *automation.Response
	if result, err = me.client().GET(automation.BusinessCalendars, id); err != nil {
		if responseErr, ok := err.(automation.ResponseErr); ok {
			return rest.Error{Code: responseErr.StatusCode, Message: string(responseErr.Data)}
		}
		return err
	}
	return json.Unmarshal(result.Data, &v)
}

func (me *service) SchemaID() string {
	return "automation:business.calendars"
}

func (me *service) List() (api.Stubs, error) {
	result, err := me.client().LIST(automation.BusinessCalendars)
	if err != nil {
		return nil, err
	}
	var stubs api.Stubs
	for _, r := range result {
		var setting business_calendars.Settings
		if err := json.Unmarshal(r.Data, &setting); err != nil {
			return nil, err
		}
		stubs = append(stubs, &api.Stub{ID: r.ID, Name: setting.Title})
	}
	return stubs, nil
}

func (me *service) Validate(v *business_calendars.Settings) error {
	return nil // no endpoint for that
}

func (me *service) Create(v *business_calendars.Settings) (stub *api.Stub, err error) {
	var data []byte
	var id string
	if data, err = json.Marshal(v); err != nil {
		return nil, err
	}
	if id, err = me.client().INSERT(automation.BusinessCalendars, data); err != nil {
		if responseErr, ok := err.(automation.ResponseErr); ok {
			return nil, rest.Error{Code: responseErr.StatusCode, Message: responseErr.Message}
		}
		return nil, err
	}
	return &api.Stub{Name: v.Title, ID: id}, nil
}

func (me *service) Update(id string, v *business_calendars.Settings) (err error) {
	var data []byte
	if data, err = json.Marshal(v); err != nil {
		return err
	}
	if err = me.client().UPDATE(automation.BusinessCalendars, id, data); err != nil {
		if responseErr, ok := err.(automation.ResponseErr); ok {
			return rest.Error{Code: responseErr.StatusCode, Message: string(responseErr.Data)}
		}
	}
	return err
}

func (me *service) Delete(id string) error {
	err := me.client().DELETE(automation.BusinessCalendars, id)
	if responseErr, ok := err.(automation.ResponseErr); ok {
		if responseErr.StatusCode == 404 {
			return nil
		}
		return rest.Error{Code: responseErr.StatusCode, Message: string(responseErr.Data)}
	}
	return err
}

func (me *service) New() *business_calendars.Settings {
	return new(business_calendars.Settings)
}

func (me *service) Name() string {
	return me.SchemaID()
}
