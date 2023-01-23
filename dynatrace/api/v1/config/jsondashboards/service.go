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
	"strings"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings"

	dashboards "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/dashboards/settings"
)

const SchemaID = "v1:config:json-dashboards"

func Service(credentials *settings.Credentials) settings.CRUDService[*dashboards.JSONDashboard] {
	return &service{settings.NewCRUDService(
		credentials,
		SchemaID,
		settings.DefaultServiceOptions[*dashboards.JSONDashboard]("/api/config/v1/dashboards").WithStubs(&dashboards.DashboardList{}),
	)}
}

type service struct {
	service settings.CRUDService[*dashboards.JSONDashboard]
}

func (me *service) List() (settings.Stubs, error) {
	stubs, err := me.service.List()
	if err != nil {
		return stubs, err
	}
	var filteredStubs settings.Stubs
	for _, stub := range stubs {
		if stub.Name != "Config owned by " {
			filteredStubs = append(filteredStubs, stub)
		}
	}
	return filteredStubs, nil
}

func (me *service) Get(id string, v *dashboards.JSONDashboard) error {
	return me.service.Get(id, v)
}

func (me *service) Validate(v *dashboards.JSONDashboard) error {
	if validator, ok := me.service.(settings.Validator[*dashboards.JSONDashboard]); ok {
		return validator.Validate(v)
	}
	return nil
}

func (me *service) Create(v *dashboards.JSONDashboard) (*settings.Stub, error) {
	return me.service.Create(v)
}

func (me *service) Update(id string, v *dashboards.JSONDashboard) error {
	jsonDashboard := v
	oldContents := jsonDashboard.Contents
	jsonDashboard.Contents = strings.Replace(oldContents, "{", `{ "id": "`+id+`", `, 1)
	err := me.service.Update(id, v)
	jsonDashboard.Contents = oldContents
	return err
}

func (me *service) Delete(id string) error {
	return me.service.Delete(id)
}

func (me *service) SchemaID() string {
	return me.service.SchemaID()
}
