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

package sharing

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings/services/cache"

	dashboards "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/dashboards/settings"
	sharing "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/dashboards/sharing/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/jsondashboards"
)

const SchemaID = "v1:config:dashboards:sharing"

func Service(credentials *settings.Credentials) settings.CRUDService[*sharing.DashboardSharing] {
	return &service{
		client:           rest.DefaultClient(credentials.URL, credentials.Token),
		dashboardService: cache.CRUD(jsondashboards.Service(credentials), true),
	}
}

type service struct {
	client           rest.Client
	dashboardService settings.CRUDService[*dashboards.JSONDashboard]
}

type NoValuesLister[T settings.Settings] interface {
	ListNoValues() (settings.Stubs, error)
}

func (me *service) Get(id string, v *sharing.DashboardSharing) error {
	id = strings.TrimSuffix(id, "-sharing")

	if err := me.client.Get(fmt.Sprintf("/api/config/v1/dashboards/%s/shareSettings", url.PathEscape(id)), 200).Finish(v); err != nil {
		return err
	}

	var dashboardName string
	var stubs settings.Stubs
	var err error

	if noValuesLister, ok := me.dashboardService.(NoValuesLister[*dashboards.Dashboard]); ok {
		if stubs, err = noValuesLister.ListNoValues(); err != nil {
			return err
		}
	} else {
		if stubs, err = me.dashboardService.List(); err != nil {
			return err
		}
	}
	for _, stub := range stubs {
		if stub.ID == id {
			dashboardName = stub.Name
			break
		}
	}
	if len(dashboardName) == 0 {
		dashboard := dashboards.JSONDashboard{}
		if err := me.dashboardService.Get(id, &dashboard); err != nil {
			return err
		}
		dashboardName = dashboard.Name()
	}

	v.Name = "ShareSettings for " + dashboardName
	return nil
}

func (me *service) Validate(v *sharing.DashboardSharing) error {
	id := v.DashboardID
	id = strings.TrimSuffix(id, "-sharing")
	if err := me.client.Post(fmt.Sprintf("/api/config/v1/dashboards/%s/shareSettings/validator", id), v, 204).Finish(); err != nil && !strings.HasPrefix(err.Error(), "No Content (PUT)") {
		return err
	}
	return nil
}

func (me *service) Update(id string, v *sharing.DashboardSharing) error {
	id = strings.TrimSuffix(id, "-sharing")
	if err := me.client.Put(fmt.Sprintf("/api/config/v1/dashboards/%s/shareSettings", id), v, 201).Finish(); err != nil && !strings.HasPrefix(err.Error(), "No Content (PUT)") {
		return err
	}
	return nil
}

func (me *service) Delete(id string) error {
	id = strings.TrimSuffix(id, "-sharing")
	settings := sharing.DashboardSharing{
		DashboardID: id,
		Enabled:     false,
		Preset:      false,
		Permissions: []*sharing.SharePermission{
			{
				Type:       sharing.PermissionTypes.All,
				Permission: sharing.Permissions.View,
			},
		},
		PublicAccess: &sharing.AnonymousAccess{
			ManagementZoneIDs: []string{},
			URLs:              map[string]string{},
		},
	}
	return me.Update(id, &settings)
}

func (me *service) Create(v *sharing.DashboardSharing) (*settings.Stub, error) {
	if err := me.Update(v.DashboardID, v); err != nil {
		return nil, err
	}
	return &settings.Stub{ID: v.DashboardID + "-sharing"}, nil
}

func (me *service) List() (settings.Stubs, error) {
	var err error

	var stubs settings.Stubs
	if stubs, err = me.dashboardService.List(); err != nil {
		return nil, err
	}
	for _, stub := range stubs {
		stub.Name = "ShareSettings for " + stub.Name
	}
	return stubs.ToStubs(), nil
}

func (me *service) SchemaID() string {
	return SchemaID
}
