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

package extension_config

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api"
	extension_config "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v2/hub/extension/config/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings"
)

func Service(credentials *settings.Credentials) settings.CRUDService[*extension_config.Settings] {
	return &service{credentials}
}

type service struct {
	credentials *settings.Credentials
}

func (me *service) Get(ctx context.Context, id string, v *extension_config.Settings) error {
	name, configurationID := splitID(id)
	var response GetMonitoringConfigurationResponse
	client := rest.DefaultClient(me.credentials.URL, me.credentials.Token)
	if err := client.Get(fmt.Sprintf("/api/v2/extensions/%s/monitoringConfigurations/%s", url.PathEscape(name), url.PathEscape(configurationID)), 200).Finish(&response); err != nil {
		return err
	}
	injectScope(response.Scope, v)
	v.Name = name
	v.Value = string(response.Value)
	return nil
}

func (me *service) List(ctx context.Context) (api.Stubs, error) {
	var stubs api.Stubs

	var extensionsList ExtensionsList
	client := rest.DefaultClient(me.credentials.URL, me.credentials.Token)

	if err := client.Get("/api/v2/extensions/info", 200).Finish(&extensionsList); err != nil {
		return stubs, err
	}
	for _, extension := range extensionsList.Extensions {
		nextPageKey := "first"
		for len(nextPageKey) > 0 {
			configList := new(MonitoringConfigurationList)
			query := "?pageSize=100"
			if len(nextPageKey) > 0 && nextPageKey != "first" {
				query = query + "&nextPageKey=" + url.QueryEscape(nextPageKey)
			}
			if err := client.Get(fmt.Sprintf("/api/v2/extensions/%s/monitoringConfigurations%s", url.PathEscape(extension.Name), query), 200).Finish(&configList); err != nil {
				return stubs, err
			}
			nextPageKey = configList.NextPageKey
			for _, config := range configList.Items {
				stubs = append(stubs, &api.Stub{ID: joinID(extension.Name, config.ObjectID), Name: joinID(extension.Name, config.ObjectID)})

			}
		}
	}

	return stubs, nil
}

func (me *service) Create(ctx context.Context, v *extension_config.Settings) (*api.Stub, error) {
	version, err := extractVersion(v)
	if err != nil {
		return nil, err
	}

	name := v.Name
	if err := me.ensureInstalled(name, version); err != nil {
		return nil, err
	}
	client := rest.DefaultClient(me.credentials.URL, me.credentials.Token)
	createResponse := []CreateMonitoringConfigResponse{}
	retry := 10
	for retry > 0 {
		payload := []MonitoringConfigCreateDto{{Scope: extractScope(v), Value: []byte(v.Value)}}
		if err := client.Post(fmt.Sprintf("/api/v2/extensions/%s/monitoringConfigurations", url.PathEscape(name)), &payload, 200).Finish(&createResponse); err != nil {
			if err.Error() == fmt.Sprintf("No schema with identifier 'ext:%s'", name) {
				time.Sleep(1 * time.Second)
				retry--
			} else {
				return nil, err
			}
		} else {
			retry = 0
		}
	}
	if len(createResponse) == 0 {
		return nil, fmt.Errorf("creating monitoring config didn't deliver ID")
	}
	return &api.Stub{ID: joinID(name, createResponse[0].ObjectID)}, nil
}

func (me *service) ensureInstalled(name string, version string) error {
	client := rest.DefaultClient(me.credentials.URL, me.credentials.Token)
	response := struct {
		Name    string `json:"extensionName"`
		Version string `json:"extensionVersion"`
	}{}
	if err := client.Post(fmt.Sprintf("/api/v2/extensions/%s?version=%s", url.PathEscape(name), url.QueryEscape(version)), nil, 200).Finish(&response); err != nil {
		if restErr, ok := err.(rest.Error); ok {
			if (restErr.Code == 400) && restErr.Message == fmt.Sprintf("Extension %s has already been added to environment", name) {
				return nil
			}
		}
		return err
	}
	return nil
}

func (me *service) Update(ctx context.Context, id string, v *extension_config.Settings) error {
	_, configID := splitID(id)
	version, err := extractVersion(v)
	if err != nil {
		return err
	}
	name := v.Name
	if err := me.ensureInstalled(name, version); err != nil {
		return err
	}
	client := rest.DefaultClient(me.credentials.URL, me.credentials.Token)
	createResponse := CreateMonitoringConfigResponse{}
	payload := MonitoringConfigCreateDto{Value: []byte(v.Value)}
	if err := client.Put(fmt.Sprintf("/api/v2/extensions/%s/monitoringConfigurations/%s", url.PathEscape(name), url.PathEscape(configID)), &payload, 200).Finish(&createResponse); err != nil {
		return err
	}
	if createResponse.ObjectID != configID {
		return fmt.Errorf("the ID of the configuration has changed during update")
	}
	return nil
}

func (me *service) Delete(ctx context.Context, id string) error {
	name, configID := splitID(id)
	client := rest.DefaultClient(me.credentials.URL, me.credentials.Token)
	if err := client.Delete(fmt.Sprintf("/api/v2/extensions/%s/monitoringConfigurations/%s", url.PathEscape(name), url.PathEscape(configID)), 200).Finish(nil); err != nil {
		if restErr, ok := err.(rest.Error); ok {
			if restErr.Code != 404 {
				return err
			}
		}
	}
	return nil
}

func (me *service) SchemaID() string {
	return "v2:extensions:twozero"
}

func (me *service) Validate(v *extension_config.Settings) error {
	return nil // no endpoint for that
}

func (me *service) New() *extension_config.Settings {
	return new(extension_config.Settings)
}

func injectScope(scope string, v *extension_config.Settings) {
	v.ActiveGateGroup = ""
	v.ManagementZone = ""
	v.Host = ""
	v.HostGroup = ""

	if strings.HasPrefix(scope, "ag_group-") {
		v.ActiveGateGroup = strings.TrimPrefix(scope, "ag_group-")
		return
	}
	if strings.HasPrefix(scope, "management_zone-") {
		v.ManagementZone = strings.TrimPrefix(scope, "management_zone-")
		return
	}
	if strings.HasPrefix(scope, "HOST-") {
		v.Host = scope
		return
	}
	if strings.HasPrefix(scope, "HOST_GROUP-") {
		v.HostGroup = scope
		return
	}
}

func extractScope(v *extension_config.Settings) string {
	if len(v.ActiveGateGroup) > 0 {
		return fmt.Sprintf("ag_group-%s", v.ActiveGateGroup)
	}
	if len(v.ManagementZone) > 0 {
		return fmt.Sprintf("management_zone-%s", v.ManagementZone)
	}
	if len(v.Host) > 0 {
		return v.Host
	}
	if len(v.HostGroup) > 0 {
		return v.HostGroup
	}
	return "environment"
}

func extractVersion(v *extension_config.Settings) (string, error) {
	valueMap := map[string]any{}
	if err := json.Unmarshal([]byte(v.Value), &valueMap); err != nil {
		return "", err
	}
	version := ""
	if untypedVersion, ok := valueMap["version"]; ok {
		if sVersion, ok := untypedVersion.(string); ok {
			version = sVersion
		}
	}
	if len(version) == 0 {
		return "", fmt.Errorf("the value of the monitoring configuration doesn't contain a 'version'")
	}
	return version, nil
}

func joinID(name string, configurationID string) string {
	return fmt.Sprintf("%s#-#%s", name, configurationID)
}

func splitID(id string) (name string, configurationID string) {
	parts := strings.Split(id, "#-#")
	if len(parts) > 0 {
		name = parts[0]
	}
	if len(parts) > 1 {
		configurationID = parts[1]
	}
	return
}

type GetMonitoringConfigurationResponse struct {
	ObjectID string          `json:"objectId"`
	Scope    string          `json:"scope"`
	Value    json.RawMessage `json:"value"`
}

type ExtensionsList struct {
	Extensions []struct {
		Name          string `json:"extensionName"`
		Version       string `json:"version"`
		ActiveVersion string `json:"activeVersion"`
	} `json:"extensions"`
}

type MonitoringConfigurationList struct {
	Items []struct {
		ObjectID string `json:"objectId"`
		Scope    string `json:"scope"`
		Value    struct {
			Version string `json:"version"`
		} `json:"value"`
	} `json:"items"`
	NextPageKey string `json:"nextPageKey"`
}

type MonitoringConfigCreateDto struct {
	Scope string          `json:"scope,omitempty"`
	Value json.RawMessage `json:"value"`
}

type CreateMonitoringConfigResponse struct {
	ObjectID string `json:"objectId"`
}
