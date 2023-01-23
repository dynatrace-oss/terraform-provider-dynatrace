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

package errors

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings/services/cache"

	webservice "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/applications/web"
	errors "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/applications/web/errors/settings"
	web "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/applications/web/settings"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/opt"
)

const SchemaID = "v1:config:applications:web:errors"

func Service(credentials *settings.Credentials) settings.CRUDService[*errors.Rules] {
	return &service{
		client:        rest.DefaultClient(credentials.URL, credentials.Token),
		webAppService: cache.CRUD(webservice.Service(credentials), true)}
}

type service struct {
	client        rest.Client
	webAppService settings.CRUDService[*web.Application]
}

func (me *service) Get(id string, v *errors.Rules) error {
	id = strings.TrimSuffix(id, "-error-rules")
	if err := me.client.Get(fmt.Sprintf("/api/config/v1/applications/web/%s/errorRules", url.PathEscape(id)), 200).Finish(v); err != nil {
		return err
	}

	// application := web.Application{}
	// if err := me.webAppService.Get(id, &application); err != nil {
	// 	return err
	// }
	stubs, err := me.webAppService.List()
	if err != nil {
		return err
	}
	for _, stub := range stubs {
		if stub.ID == id {
			v.Name = "Data Privacy Settings for " + stub.Name
			v.WebApplicationID = id
			break
		}
	}
	// v.Name = "Error Rules for " + application.Name
	// v.WebApplicationID = id
	return nil
}

func (me *service) Validate(v *errors.Rules) error {
	id := v.WebApplicationID
	id = strings.TrimSuffix(id, "-data-privacy")
	err := me.client.Post(fmt.Sprintf("/api/config/v1/applications/web/%s/errorRules/validator", id), v, 204).Finish()
	if err != nil && strings.HasPrefix(err.Error(), "No Content (PUT)") {
		return nil
	}
	return err
}

func (me *service) Update(id string, v *errors.Rules) error {
	id = strings.TrimSuffix(id, "-error-rules")
	err := me.client.Put(fmt.Sprintf("/api/config/v1/applications/web/%s/errorRules", id), v, 201, 204).Finish()
	if err != nil && strings.HasPrefix(err.Error(), "No Content (PUT)") {
		return nil
	}
	return err
}

func (me *service) Delete(id string) error {
	id = strings.TrimSuffix(id, "-error-rules")
	settings := errors.Rules{
		IgnoreJavaScriptErrorsInApdexCalculation: false,
		IgnoreHttpErrorsInApdexCalculation:       false,
		IgnoreCustomErrorsInApdexCalculation:     false,
		HTTPErrors: errors.HTTPErrorRules{
			{
				ConsiderUnknownErrorCode: false,
				ConsiderBlockedRequests:  false,
				ErrorCodes:               opt.NewString("401"),
				FilterByURL:              false,
				Filter:                   nil,
				URL:                      nil,
				Capture:                  false,
				ImpactApdex:              false,
				ConsiderForAI:            false,
			},
			{
				ConsiderUnknownErrorCode: false,
				ConsiderBlockedRequests:  false,
				ErrorCodes:               opt.NewString("404"),
				FilterByURL:              true,
				Filter:                   &errors.HTTPErrorRuleFilters.EndsWith,
				URL:                      opt.NewString("favicon.ico"),
				Capture:                  false,
				ImpactApdex:              false,
				ConsiderForAI:            false,
			},
			{
				ConsiderUnknownErrorCode: false,
				ConsiderBlockedRequests:  false,
				ErrorCodes:               opt.NewString("404"),
				FilterByURL:              true,
				Filter:                   &errors.HTTPErrorRuleFilters.EndsWith,
				URL:                      opt.NewString(".js"),
				Capture:                  true,
				ImpactApdex:              true,
				ConsiderForAI:            true,
			},
			{
				ConsiderUnknownErrorCode: false,
				ConsiderBlockedRequests:  false,
				ErrorCodes:               opt.NewString("404"),
				FilterByURL:              true,
				Filter:                   &errors.HTTPErrorRuleFilters.EndsWith,
				URL:                      opt.NewString(".css"),
				Capture:                  true,
				ImpactApdex:              true,
				ConsiderForAI:            true,
			},
			{
				ConsiderUnknownErrorCode: false,
				ConsiderBlockedRequests:  false,
				ErrorCodes:               opt.NewString("400-499"),
				FilterByURL:              false,
				Filter:                   nil,
				URL:                      nil,
				Capture:                  true,
				ImpactApdex:              true,
				ConsiderForAI:            false,
			},
			{
				ConsiderUnknownErrorCode: false,
				ConsiderBlockedRequests:  false,
				ErrorCodes:               opt.NewString("500-599"),
				FilterByURL:              false,
				Filter:                   nil,
				URL:                      nil,
				Capture:                  true,
				ImpactApdex:              true,
				ConsiderForAI:            true,
			},
			{
				ConsiderUnknownErrorCode: true,
				ConsiderBlockedRequests:  false,
				ErrorCodes:               nil,
				FilterByURL:              false,
				Filter:                   nil,
				URL:                      nil,
				Capture:                  true,
				ImpactApdex:              false,
				ConsiderForAI:            false,
			},
			{
				ConsiderUnknownErrorCode: false,
				ConsiderBlockedRequests:  true,
				ErrorCodes:               nil,
				FilterByURL:              false,
				Filter:                   nil,
				URL:                      nil,
				Capture:                  true,
				ImpactApdex:              false,
				ConsiderForAI:            false,
			},
		},
		CustomErrors: errors.CustomErrorRules{
			{
				KeyPattern:     nil,
				KeyMatcher:     nil,
				ValuePattern:   opt.NewString("example"),
				ValueMatcher:   &errors.CustomErrorRuleValueMatchers.BeginsWith,
				Capture:        true,
				ImpactApdex:    true,
				CustomAlerting: true,
			},
			{
				KeyPattern:     nil,
				KeyMatcher:     nil,
				ValuePattern:   nil,
				ValueMatcher:   nil,
				Capture:        true,
				ImpactApdex:    false,
				CustomAlerting: false,
			},
		},
	}
	return me.Update(id, &settings)
}

func (me *service) Create(v *errors.Rules) (*settings.Stub, error) {
	if err := me.Update(v.WebApplicationID, v); err != nil {
		return nil, err
	}
	return &settings.Stub{ID: v.WebApplicationID + "-error-rules"}, nil
}

func (me *service) List() (settings.Stubs, error) {
	var err error
	var stubs settings.Stubs
	if stubs, err = me.webAppService.List(); err != nil {
		return nil, err
	}
	for _, stub := range stubs {
		if strings.HasSuffix(stub.ID, "-data-privacy") {
			panic("")
		}
		if strings.HasSuffix(stub.ID, "-error-rules") {
			panic("")
		}
		stub.Name = "Error Rules for " + stub.Name
		stub.ID = stub.ID + "-error-rules"
	}
	return stubs.ToStubs(), nil
}

func (me *service) SchemaID() string {
	return SchemaID
}
