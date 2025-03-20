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
	"context"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings"

	dashboardsbase "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/dashboardsbase/settings"
)

const SchemaID = "v1:config:json-dashboards-base"

func Service(credentials *rest.Credentials) settings.CRUDService[*dashboardsbase.JSONDashboardBase] {
	return &service{settings.NewAPITokenService(
		credentials,
		SchemaID,
		settings.DefaultServiceOptions[*dashboardsbase.JSONDashboardBase]("/api/config/v1/dashboards").WithStubs(&dashboardsbase.DashboardList{}),
	)}
}

type service struct {
	service settings.CRUDService[*dashboardsbase.JSONDashboardBase]
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

func (me *service) Get(ctx context.Context, id string, v *dashboardsbase.JSONDashboardBase) error {
	return nil
}

func (me *service) Validate(ctx context.Context, v *dashboardsbase.JSONDashboardBase) error {
	return nil
}

func (me *service) Create(ctx context.Context, v *dashboardsbase.JSONDashboardBase) (*api.Stub, error) {
	return me.service.Create(ctx, v)
}

func (me *service) Update(ctx context.Context, id string, v *dashboardsbase.JSONDashboardBase) error {
	return nil
}

func (me *service) Delete(ctx context.Context, id string) error {
	return me.service.Delete(ctx, id)
}

func (me *service) SchemaID() string {
	return me.service.SchemaID()
}
