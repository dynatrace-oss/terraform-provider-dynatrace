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

package jsondashboardsbase

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings"

	dashboardsbase "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/dashboardsbase/settings"
)

const SchemaID = "v1:config:json-dashboards-base"

func Service(credentials *settings.Credentials) settings.CRUDService[*dashboardsbase.JSONDashboardBase] {
	return &service{settings.NewCRUDService(
		credentials,
		SchemaID,
		settings.DefaultServiceOptions[*dashboardsbase.JSONDashboardBase]("/api/config/v1/dashboards").WithStubs(&dashboardsbase.DashboardList{}),
	)}
}

type service struct {
	service settings.CRUDService[*dashboardsbase.JSONDashboardBase]
}

func (me *service) List() (api.Stubs, error) {
	stubs, err := me.service.List()
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

func (me *service) Get(id string, v *dashboardsbase.JSONDashboardBase) error {
	return nil
}

func (me *service) Validate(v *dashboardsbase.JSONDashboardBase) error {
	return nil
}

func (me *service) Create(v *dashboardsbase.JSONDashboardBase) (*api.Stub, error) {
	return me.service.Create(v)
}

func (me *service) Update(id string, v *dashboardsbase.JSONDashboardBase) error {
	return nil
}

func (me *service) Delete(id string) error {
	return me.service.Delete(id)
}

func (me *service) SchemaID() string {
	return me.service.SchemaID()
}

func (me *service) Name() string {
	return me.service.Name()
}
