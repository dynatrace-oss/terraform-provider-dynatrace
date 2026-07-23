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
	business_calendars "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/automation/business_calendars/settings"
	tfrest "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings"

	cacapi "github.com/dynatrace/dynatrace-configuration-as-code-core/api"
	"github.com/dynatrace/dynatrace-configuration-as-code-core/clients/automation"
)

func Service(clientSet tfrest.ClientSet) (settings.CRUDService[*business_calendars.Settings], error) {
	platformClient, err := clientSet.PlatformClient()
	if err != nil {
		return nil, err
	}

	return &service{
		client: automation.NewClient(platformClient),
	}, nil
}

type service struct {
	client *automation.Client
}

func (me *service) Get(ctx context.Context, id string, v *business_calendars.Settings) (err error) {
	response, err := me.client.Get(ctx, automation.BusinessCalendars, id)
	if err != nil {
		return err
	}
	return json.Unmarshal(response.Data, &v)
}

func (me *service) SchemaID() string {
	return "automation:business.calendars"
}

type BusinessCalendarStub struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}

func (me *service) List(ctx context.Context) (api.Stubs, error) {
	listResponse, err := me.client.List(ctx, automation.BusinessCalendars)
	if err != nil {
		return nil, err
	}

	var stubs api.Stubs
	for _, r := range listResponse.All() {
		var stub BusinessCalendarStub
		if err := json.Unmarshal(r, &stub); err != nil {
			return nil, err
		}
		stubs = append(stubs, &api.Stub{ID: stub.ID, Name: stub.Title})
	}
	return stubs, nil
}

func (me *service) Validate(v *business_calendars.Settings) error {
	return nil // no endpoint for that
}

func (me *service) Create(ctx context.Context, v *business_calendars.Settings) (stub *api.Stub, err error) {
	var data []byte
	if data, err = json.Marshal(v); err != nil {
		return nil, err
	}

	response, err := me.client.Create(ctx, automation.BusinessCalendars, data)
	if err != nil {
		return nil, err
	}
	var bStub BusinessCalendarStub
	if err := json.Unmarshal(response.Data, &bStub); err != nil {
		return nil, err
	}
	return &api.Stub{Name: v.Title, ID: bStub.ID}, nil
}

func (me *service) Update(ctx context.Context, id string, v *business_calendars.Settings) (err error) {
	var data []byte
	if data, err = json.Marshal(v); err != nil {
		return err
	}

	_, err = me.client.Update(ctx, automation.BusinessCalendars, id, data)
	return err
}

func (me *service) Delete(ctx context.Context, id string) error {
	if _, err := me.client.Delete(ctx, automation.BusinessCalendars, id); err != nil && !cacapi.IsNotFoundError(err) {
		return err
	}
	return nil
}

func (me *service) New() *business_calendars.Settings {
	return new(business_calendars.Settings)
}
