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
	"context"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings/services/cache"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings/services/httpcache"

	dashboards "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/dashboards/settings"
	sharing "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/dashboards/sharing/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/jsondashboards"
)

const SchemaID = "v1:config:dashboards:sharing"

func Service(credentials *settings.Credentials) settings.CRUDService[*sharing.DashboardSharing] {
	return &service{
		client:           httpcache.DefaultClient(credentials.URL, credentials.Token, SchemaID),
		dashboardService: cache.CRUD(jsondashboards.Service(credentials), true),
	}
}

type service struct {
	client           rest.Client
	dashboardService settings.CRUDService[*dashboards.JSONDashboard]
}

type NoValuesLister[T settings.Settings] interface {
	ListNoValues() (api.Stubs, error)
}

type DashbordMeta struct {
	DashboardMetaData struct {
		Preset bool `json:"preset"`
	} `json:"dashboardMetadata"`
}

func (me *service) Get(ctx context.Context, id string, v *sharing.DashboardSharing) error {
	id = strings.TrimSuffix(id, "-sharing")

	var dbm DashbordMeta
	if err := me.client.Get(fmt.Sprintf("/api/config/v1/dashboards/%s", url.PathEscape(id)), 200).Finish(&dbm); err != nil {
		return err
	}

	if err := me.client.Get(fmt.Sprintf("/api/config/v1/dashboards/%s/shareSettings", url.PathEscape(id)), 200).Finish(v); err != nil {
		return err
	}

	v.Muted = dbm.DashboardMetaData.Preset

	var dashboardName string
	var stubs api.Stubs
	var err error

	if noValuesLister, ok := me.dashboardService.(NoValuesLister[*dashboards.Dashboard]); ok {
		if stubs, err = noValuesLister.ListNoValues(); err != nil {
			return err
		}
	} else {
		if stubs, err = me.dashboardService.List(ctx); err != nil {
			return err
		}
	}
	for _, stub := range stubs {

		if stub.ID == id {
			dashboardName = stub.Name
			break
		}
	}
	if len(dashboardName) == 0 || dashboardName == id {
		dashboard := dashboards.JSONDashboard{}
		if err := me.dashboardService.Get(ctx, id, &dashboard); err != nil {
			return err
		}
		dashboardName = dashboard.Name()
	}

	v.DashboardName = dashboardName
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

func (me *service) Update(ctx context.Context, id string, v *sharing.DashboardSharing) error {
	return me.update(ctx, id, v, 0)
}

const max_retries = 10

func (me *service) update(ctx context.Context, id string, v *sharing.DashboardSharing, retry int) error {
	id = strings.TrimSuffix(id, "-sharing")

	var dbm DashbordMeta
	if err := me.client.Get(fmt.Sprintf("/api/config/v1/dashboards/%s", url.PathEscape(id)), 200).Finish(&dbm); err != nil {
		return err
	}

	if err := me.client.Put(fmt.Sprintf("/api/config/v1/dashboards/%s/shareSettings", id), v, 201, 204).Finish(); err != nil && !strings.HasPrefix(err.Error(), "No Content (PUT)") {
		// newly created dashboards are sometimes not yet known cluster wide
		// this retry functionality tries to make up for that
		if strings.Contains(err.Error(), "Dashboard does not exist") {
			if retry > max_retries {
				return err
			}
			time.Sleep(10 * time.Second)
			return me.update(ctx, id, v, retry+1)
		}
		if strings.Contains(err.Error(), "Sharing settings of a preset can't be updated. It's shared by default") {
			return nil
		}
		return err
	}
	return nil
}

func (me *service) Delete(ctx context.Context, id string) error {
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
	return me.Update(ctx, id, &settings)
}

func (me *service) Create(ctx context.Context, v *sharing.DashboardSharing) (*api.Stub, error) {
	if err := me.Update(ctx, v.DashboardID, v); err != nil {
		return nil, err
	}
	return &api.Stub{ID: v.DashboardID + "-sharing"}, nil
}

func (me *service) List(ctx context.Context) (api.Stubs, error) {
	var err error

	var stubs api.Stubs
	if stubs, err = me.dashboardService.List(ctx); err != nil {
		return nil, err
	}
	return stubs.ToStubs(), nil
}

func (me *service) SchemaID() string {
	return SchemaID
}
