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

package mode

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api"
	mode "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/host/monitoring/mode/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings/services/settings20"
)

const SchemaVersion = "1.3"
const SchemaID = "builtin:host.monitoring.mode"

var WarnOnAgentOffline = os.Getenv("DYNATRACE_HOST_MONITORING_WARNINGS") == "true"
var ExportOfflineHosts = os.Getenv("DYNATRACE_HOST_MONITORING_OFFLINE") == "true"
var StrictUpdateRetries = os.Getenv("DYNATRACE_HOST_MONITORING_STRICT_UPDATE_RETRIES")

func Service(credentials *settings.Credentials) settings.CRUDService[*mode.Settings] {
	return &service{
		service:     settings20.Service[*mode.Settings](credentials, SchemaID, SchemaVersion),
		credentials: credentials,
		client:      rest.DefaultClient(credentials.URL, credentials.Token),
	}
}

type service struct {
	service     settings.CRUDService[*mode.Settings]
	credentials *settings.Credentials
	client      rest.Client
}

type ListResponse struct {
	NextPageKey string `json:"nextPageKey"`
	Entities    []struct {
		EntityID    string `json:"entityId"`
		DisplayName string `json:"displayName"`
		Properties  struct {
			MonitoringMode string `json:"monitoringMode"`
			State          string `json:"state"`
		} `json:"properties"`
	} `json:"entities"`
}

func (me *service) List(ctx context.Context) (api.Stubs, error) {
	collectStubs := func(stubs api.Stubs, listResponse *ListResponse) api.Stubs {
		for _, entity := range listResponse.Entities {
			setting := new(mode.Settings)
			setting.HostID = entity.EntityID
			if len(entity.Properties.MonitoringMode) == 0 {
				continue
			}
			if len(entity.Properties.State) == 0 {
				continue
			}
			if !ExportOfflineHosts && (entity.Properties.State == "OFFLINE" || entity.Properties.State == "SHUTDOWN" || entity.Properties.State == "MONITORING_DISABLED") {
				continue
			}
			if entity.Properties.MonitoringMode == "INFRASTRUCTURE" {
				entity.Properties.MonitoringMode = "INFRA_ONLY"
			}
			setting.MonitoringMode = mode.OamonitoringMode(entity.Properties.MonitoringMode)
			stub := new(api.Stub)
			stub.ID = entity.EntityID
			stub.Name = entity.DisplayName
			stub.Value = setting
			stubs = append(stubs, stub)
		}
		return stubs
	}
	var stubs api.Stubs
	listResponse := new(ListResponse)
	from := "now-3y"
	fields := "properties.monitoringMode,properties.state"
	fullURL := fmt.Sprintf("/api/v2/entities?pageSize=500&entitySelector=%s&from=%s&fields=%s", url.QueryEscape("type(HOST)"), url.QueryEscape(from), url.QueryEscape(fields))
	err := me.client.Get(ctx, fullURL).Expect(200).Finish(listResponse)
	if err != nil {
		return stubs, err
	}
	stubs = collectStubs(stubs, listResponse)
	nextPageKey := listResponse.NextPageKey
	for len(nextPageKey) > 0 {
		listResponse := new(ListResponse)
		err := me.client.Get(ctx, fmt.Sprintf("/api/v2/entities?nextPageKey=%s", url.QueryEscape(nextPageKey))).Expect(200).Finish(listResponse)
		if err != nil {
			return stubs, err
		}
		nextPageKey = listResponse.NextPageKey
		stubs = collectStubs(stubs, listResponse)
	}
	return stubs, nil
}

func (me *service) Get(ctx context.Context, id string, v *mode.Settings) error {
	listResponse := new(ListResponse)
	from := "now-3y"
	fields := "properties.monitoringMode,properties.state"
	fullURL := fmt.Sprintf("/api/v2/entities?pageSize=1&entitySelector=%s&from=%s&fields=%s", url.QueryEscape(fmt.Sprintf("type(HOST),entityId(%s)", id)), url.QueryEscape(from), url.QueryEscape(fields))
	err := me.client.Get(ctx, fullURL).Expect(200).Finish(listResponse)
	if err != nil {
		return err
	}
	if len(listResponse.Entities) == 0 {
		return rest.Error{Code: 404, Message: fmt.Sprintf("monitoringMode for host '%s' not found", id)}
	}
	if len(listResponse.Entities[0].Properties.MonitoringMode) == 0 {
		return rest.Error{Code: 404, Message: fmt.Sprintf("monitoringMode for host '%s' not found", id)}
	}
	v.HostID = id
	monitoringMode := listResponse.Entities[0].Properties.MonitoringMode
	if monitoringMode == "INFRASTRUCTURE" {
		monitoringMode = "INFRA_ONLY"
	}
	v.MonitoringMode = mode.OamonitoringMode(monitoringMode)

	return nil
}

func (me *service) SchemaID() string {
	return SchemaID
}

func (me *service) Create(ctx context.Context, v *mode.Settings) (*api.Stub, error) {
	err := me.update(ctx, v.HostID, v)
	// it's ok to use HostID as the name here - not being used during creation
	return &api.Stub{ID: v.HostID, Name: v.HostID}, err
}

func (me *service) Update(ctx context.Context, id string, v *mode.Settings) error {
	return me.update(ctx, id, v)
}

func (me *service) update(ctx context.Context, id string, v *mode.Settings) error {
	_, err := me.service.Create(ctx, v)
	if err == nil {
		remainingRetries := 0
		if len(StrictUpdateRetries) > 0 {
			if iStrictUpdateRetries, err := strconv.Atoi(StrictUpdateRetries); err == nil && iStrictUpdateRetries >= 0 {
				remainingRetries = iStrictUpdateRetries
			}
		}
		readMode := ""
		for remainingRetries > 0 && readMode != string(v.MonitoringMode) {
			getSettings := new(mode.Settings)
			me.Get(ctx, id, getSettings)
			readMode = string(getSettings.MonitoringMode)
			if readMode == string(v.MonitoringMode) {
				break
			}
			time.Sleep(20 * time.Second)
			remainingRetries--
		}

		return nil
	}
	// we don't want to error out just because the Agent isn't online
	if strings.Contains(err.Error(), "Monitoring mode can't be changed while OneAgent is offline") {
		if WarnOnAgentOffline {
			return rest.Warning{Message: fmt.Sprintf("The host '%s' is currently offline - changing the monitoring mode was not possible", id)}
		}
		return nil
	}
	return err
}

func (me *service) Delete(ctx context.Context, id string) error {
	// Deleting the monitoring mode for a host doesn't make sense
	// it always exists
	return nil
}
