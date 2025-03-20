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
	"net/url"
	"time"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api"
	tfrest "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings"

	automationerr "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/automation"
	business_calendars "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/automation/business_calendars/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/automation/httplog"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/monaco/pkg/client/auth"
	apiClient "github.com/dynatrace/dynatrace-configuration-as-code-core/api/clients/automation"
	"github.com/dynatrace/dynatrace-configuration-as-code-core/api/rest"
	"github.com/dynatrace/dynatrace-configuration-as-code-core/clients/automation"
)

func Service(credentials *tfrest.Credentials) settings.CRUDService[*business_calendars.Settings] {
	return &service{credentials}
}

type service struct {
	credentials *tfrest.Credentials
}

func (me *service) client(ctx context.Context) *automation.Client {
	httplog.InstallRoundTripper()
	httpClient := auth.NewOAuthClient(ctx, auth.OauthCredentials{
		ClientID:     me.credentials.OAuth.ClientID,
		ClientSecret: me.credentials.OAuth.ClientSecret,
		TokenURL:     me.credentials.OAuth.TokenURL,
	})
	u, _ := url.Parse(me.credentials.OAuth.EnvironmentURL)

	restClient := rest.NewClient(u, httpClient, rest.WithRateLimiter(), rest.WithRetryOptions(&rest.RetryOptions{MaxRetries: 30, DelayAfterRetry: 10 * time.Second, ShouldRetryFunc: rest.RetryIfTooManyRequests}))
	restClient.SetHeader("User-Agent", "Dynatrace Terraform Provider")
	return automation.NewClient(restClient)
}

func (me *service) Get(ctx context.Context, id string, v *business_calendars.Settings) (err error) {
	var response automation.Response
	if response, err = me.client(ctx).Get(ctx, apiClient.BusinessCalendars, id); err != nil {
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
	return "automation:business.calendars"
}

type BusinessCalendarStub struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}

func (me *service) List(ctx context.Context) (api.Stubs, error) {
	listResponse, err := me.client(ctx).List(ctx, apiClient.BusinessCalendars)
	if err != nil {
		return nil, err
	}
	if apiErr, ok := listResponse.AsAPIError(); ok {
		return nil, tfrest.Error{Code: apiErr.StatusCode, Message: string(apiErr.Body)}
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
	var response automation.Response
	if response, err = me.client(ctx).Create(ctx, apiClient.BusinessCalendars, data); err != nil {
		return nil, err
	}
	if response.StatusCode == 201 {
		var stub BusinessCalendarStub
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

func (me *service) Update(ctx context.Context, id string, v *business_calendars.Settings) (err error) {
	var data []byte
	if data, err = json.Marshal(v); err != nil {
		return err
	}
	var response automation.Response
	if response, err = me.client(ctx).Update(ctx, apiClient.BusinessCalendars, id, data); err != nil {
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
	response, err := me.client(ctx).Delete(ctx, apiClient.BusinessCalendars, id)
	if response.StatusCode == 204 || response.StatusCode == 404 {
		return nil
	}
	if response.StatusCode != 0 {
		return tfrest.Error{Code: response.StatusCode, Message: string(response.Data)}
	}
	return err
}

func (me *service) New() *business_calendars.Settings {
	return new(business_calendars.Settings)
}
