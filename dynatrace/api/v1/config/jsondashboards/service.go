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

package jsondashboards

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings"

	dashboards "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/dashboards/settings"
)

const SchemaID = "v1:config:json-dashboards"

func Service(credentials *settings.Credentials) settings.CRUDService[*dashboards.JSONDashboard] {
	return &service{settings.NewCRUDService(
		credentials,
		SchemaID,
		settings.DefaultServiceOptions[*dashboards.JSONDashboard]("/api/config/v1/dashboards").WithStubs(&dashboards.DashboardList{}).WithAfterCreate(AfterCreate),
	)}
}

type DashboardPreset struct {
	DashboardPreset string `json:"DashboardPreset"`
	UserGroup       string `json:"UserGroup"`
}

type DashboardPresetPayloadValue struct {
	EnableDashboardPresets bool               `json:"enableDashboardPresets"`
	DashboardPresetsList   []*DashboardPreset `json:"dashboardPresetsList"`
}

type DashboardPresetPayload struct {
	SchemaID string                       `json:"schemaId"`
	Scope    string                       `json:"scope"`
	Value    *DashboardPresetPayloadValue `json:"value"`
}

const dummUserGroupID = "0d80b6f2-64c8-430e-8c10-83b8beda3fd3"
const userGroupViolationMessagePrefix = "Given property 'UserGroup' with value: '0d80b6f2-64c8-430e-8c10-83b8beda3fd3' is not a valid value in datasource. Value must be one of ["

func AfterCreate(ctx context.Context, client rest.Client, stub *api.Stub) (stubs *api.Stub, err error) {
	// This function just ATTEMPTS to validate whether the Dashboard is ready for use via Dashboard Presets
	// If an error happens (e.g. because of missing / wrong credentials) it doesn't error out
	//
	// What happens here:
	// * The Settings 2.0 API receives payload for VALIDATION of Dashboard Presets
	// * The first request is bound to fail because we don't have a User Group to specify available at this point
	//   * We can extract a valid User Group ID from the first error message
	// * After that sending the payload is required to succeed FIVE times
	//   * At most 50 attempts (with 500 sec sleep time in between) are being made
	//   * If less than 5 successes are counted the function doesn't error out
	//   * But even then it may have stalled Terraform long enough, so that subsequent
	//     dynatrace_dashboard_presets resources won't fail anymore
	payload := DashboardPresetPayload{
		SchemaID: "builtin:dashboards.presets",
		Scope:    "environment",
		Value: &DashboardPresetPayloadValue{
			EnableDashboardPresets: true,
			DashboardPresetsList: []*DashboardPreset{
				{
					DashboardPreset: stub.ID,
					UserGroup:       dummUserGroupID, // first request sends dummy user groupID
				},
			},
		},
	}
	numRetries := 50
	numSuccesses := 0
	for numRetries > 0 && numSuccesses < 5 {
		numRetries--
		validated, err := ValidatePreset(ctx, client, &payload)
		if err != nil {
			// some other error has happened - we will silently ignore that
			// the dashboard HAS been created, the sanity check just couldn't be done
			return stub, nil
		}
		if validated {
			numSuccesses++
		}
		time.Sleep(500)
	}
	return stub, nil
}

func ValidatePreset(ctx context.Context, client rest.Client, payload *DashboardPresetPayload) (validated bool, err error) {
	p := []*DashboardPresetPayload{payload}
	err = client.Post(ctx, "/api/v2/settings/objects?repairInput=true&validateOnly=true", &p, 200).Finish(nil)
	if err != nil {
		if restErr, ok := err.(rest.Error); ok {
			for _, violation := range restErr.ConstraintViolations {
				violationMessage := violation.Message
				if strings.HasPrefix(violationMessage, userGroupViolationMessagePrefix) {
					violationMessage = violationMessage[len(userGroupViolationMessagePrefix):]
					idx := strings.Index(violationMessage, ",")
					if idx < 0 {
						return false, errors.New("no groupID found")
					}
					parts := strings.Split(violationMessage, ", ")
					if len(parts) == 0 {
						return false, errors.New("no groupID found")
					}
					payload.Value.DashboardPresetsList[0].UserGroup = parts[0]
					return false, nil
				}
			}

			for _, violation := range restErr.ConstraintViolations {
				violationMessage := violation.Message
				if strings.HasPrefix(violationMessage, fmt.Sprintf("Given property 'DashboardPreset' with value: '%s' is not a valid value in datasource.", payload.Value.DashboardPresetsList[0].DashboardPreset)) {
					return false, nil
				}
			}
		} else {
			return false, err
		}
	}
	return true, nil
}

type service struct {
	service settings.CRUDService[*dashboards.JSONDashboard]
}

func (me *service) List(ctx context.Context) (api.Stubs, error) {
	stubs, err := me.service.List(ctx)
	if err != nil {
		return stubs, err
	}
	var filteredStubs api.Stubs
	for _, stub := range stubs {
		if stub.Name != "Config owned by " {
			filteredStubs = append(filteredStubs, stub)
		}
	}
	return filteredStubs, nil
}

func (me *service) Get(ctx context.Context, id string, v *dashboards.JSONDashboard) error {
	if err := me.service.Get(ctx, id, v); err != nil {
		return err
	}
	v.DeNull()
	return nil
}

func (me *service) Validate(ctx context.Context, v *dashboards.JSONDashboard) error {
	if validator, ok := me.service.(settings.Validator[*dashboards.JSONDashboard]); ok {
		return validator.Validate(ctx, v)
	}
	return nil
}

func (me *service) Create(ctx context.Context, v *dashboards.JSONDashboard) (*api.Stub, error) {
	doCreateService := true
	var stub *api.Stub = nil
	var err error = nil

	if len(v.LinkID) > 0 {
		doCreateService = false

		err = me.Update(ctx, v.LinkID, v)
		stub = &api.Stub{ID: v.LinkID}

		if restError, ok := err.(rest.Error); ok {
			if restError.Code == 404 {
				doCreateService = true
			}
		}
	}
	if doCreateService {
		stub, err = me.service.Create(ctx, v.EnrichRequireds())
	}

	return stub, err
}

func (me *service) Update(ctx context.Context, id string, v *dashboards.JSONDashboard) error {

	if len(v.LinkID) > 0 {
		if id != strings.ToLower(v.LinkID) {
			return fmt.Errorf("dashboard ID cannot be modified, please destroy and create with the new ID")
		}
	}

	jsonDashboard := v
	oldContents := jsonDashboard.Contents
	jsonDashboard.Contents = strings.Replace(oldContents, "{", `{ "id": "`+id+`", `, 1)
	err := me.service.Update(ctx, id, v.EnrichRequireds())
	jsonDashboard.Contents = oldContents
	return err
}

func (me *service) Delete(ctx context.Context, id string) error {
	return me.service.Delete(ctx, id)
}

func (me *service) SchemaID() string {
	return me.service.SchemaID()
}
