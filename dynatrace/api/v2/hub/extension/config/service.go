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
	"regexp"
	"strings"
	"time"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api"
	extension_config "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v2/hub/extension/config/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/provider/logging"
)

func Service(credentials *rest.Credentials) settings.CRUDService[*extension_config.Settings] {
	return &service{credentials}
}

type service struct {
	credentials *rest.Credentials
}

var pattern = regexp.MustCompile(`^\*{3}\d{3}\*{3}$`)

func replaceCredPlaceholders(a, b any) any {
	switch aVal := a.(type) {
	case map[string]any:
		bMap, ok := b.(map[string]any)
		if !ok {
			return b
		}
		for key, bVal := range bMap {
			if aInnerVal, exists := aVal[key]; exists {
				bMap[key] = replaceCredPlaceholders(aInnerVal, bVal)
			}
		}
	case []any:
		bSlice, ok := b.([]any)
		if !ok || len(aVal) != len(bSlice) {
			return b
		}
		for i := range bSlice {
			bSlice[i] = replaceCredPlaceholders(aVal[i], bSlice[i])
		}
	case string:
		bStr, ok := b.(string)
		if ok && pattern.MatchString(bStr) {
			return aVal
		}
	}
	return b
}

func (me *service) Get(ctx context.Context, id string, v *extension_config.Settings) error {
	cfg := ctx.Value(settings.ContextKeyStateConfig)
	stateConfig, _ := cfg.(*extension_config.Settings)

	name, configurationID := splitID(id)

	var response GetMonitoringConfigurationResponse
	client := rest.APITokenClient(me.credentials)

	urlPath := fmt.Sprintf(
		"/api/v2/extensions/%s/monitoringConfigurations/%s",
		url.PathEscape(name),
		url.PathEscape(configurationID),
	)

	if err := client.Get(ctx, urlPath, 200).Finish(&response); err != nil {
		return err
	}

	injectScope(response.Scope, v)
	v.Name = name

	// Try to replace placeholders only if stateConfig and response.Value are valid
	if stateConfig != nil && response.Value != nil {
		var (
			receivedValue map[string]any
			stateValue    map[string]any
		)

		if err := json.Unmarshal(response.Value, &receivedValue); err != nil {
			logging.File.Printf("Error unmarshalling received value: %v\n", err)
		} else if err := json.Unmarshal([]byte(stateConfig.Value), &stateValue); err != nil {
			logging.File.Printf("Error unmarshalling state config: %v\n", err)
		} else {
			merged := replaceCredPlaceholders(stateValue, receivedValue)
			if marshalled, err := json.MarshalIndent(merged, "", "  "); err != nil {
				logging.File.Printf("Error marshalling merged config: %v\n", err)
			} else {
				v.Value = string(marshalled)
				return nil
			}
		}
	}

	// Fallback: use response value directly
	v.Value = string(response.Value)
	return nil
}

func (me *service) List(ctx context.Context) (api.Stubs, error) {
	var stubs api.Stubs

	var extensionsList ExtensionsList
	client := rest.APITokenClient(me.credentials)

	nextPageKey := "first"
	for len(nextPageKey) > 0 {
		extensionsList = ExtensionsList{}
		query := ""
		if len(nextPageKey) > 0 && nextPageKey != "first" {
			query = "?nextPageKey=" + url.QueryEscape(nextPageKey)
		} else {
			query = "?pageSize=100"
		}
		if err := client.Get(ctx, "/api/v2/extensions/info"+query, 200).Finish(&extensionsList); err != nil {
			return stubs, err
		}
		nextPageKey = extensionsList.NextPageKey
	}
	for _, extension := range extensionsList.Extensions {
		nextPageKey := "first"
		for len(nextPageKey) > 0 {
			configList := new(MonitoringConfigurationList)
			query := "?pageSize=100"
			if len(nextPageKey) > 0 && nextPageKey != "first" {
				query = "?nextPageKey=" + url.QueryEscape(nextPageKey)
			}
			if err := client.Get(ctx, fmt.Sprintf("/api/v2/extensions/%s/monitoringConfigurations%s", url.PathEscape(extension.Name), query), 200).Finish(&configList); err != nil {
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
	if err := me.ensureInstalled(ctx, name, version); err != nil {
		return nil, err
	}
	client := rest.APITokenClient(me.credentials)
	createResponse := []CreateMonitoringConfigResponse{}
	retry := 10
	for retry > 0 {
		payload := []MonitoringConfigCreateDto{{Scope: extractScope(v), Value: []byte(v.Value)}}
		if err := client.Post(ctx, fmt.Sprintf("/api/v2/extensions/%s/monitoringConfigurations", url.PathEscape(name)), &payload, 200).Finish(&createResponse); err != nil {
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

func (me *service) ensureInstalled(ctx context.Context, name string, version string) error {
	client := rest.APITokenClient(me.credentials)
	if strings.HasPrefix(name, "custom:") {
		request := client.Get(ctx, fmt.Sprintf("/api/v2/extensions/%s/%s", url.PathEscape(name), url.QueryEscape(version)), 200)
		request.SetHeader("Accept", "application/json; charset=utf-8")
		if err := request.Finish(); err != nil {
			return err
		}
	}
	response := struct {
		Name    string `json:"extensionName"`
		Version string `json:"extensionVersion"`
	}{}
	if err := client.Post(ctx, fmt.Sprintf("/api/v2/extensions/%s?version=%s", url.PathEscape(name), url.QueryEscape(version)), nil, 200).Finish(&response); err != nil {
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
	if err := me.ensureInstalled(ctx, name, version); err != nil {
		return err
	}
	client := rest.APITokenClient(me.credentials)
	createResponse := CreateMonitoringConfigResponse{}
	payload := MonitoringConfigCreateDto{Value: []byte(v.Value)}
	if err := client.Put(ctx, fmt.Sprintf("/api/v2/extensions/%s/monitoringConfigurations/%s", url.PathEscape(name), url.PathEscape(configID)), &payload, 200).Finish(&createResponse); err != nil {
		return err
	}
	if createResponse.ObjectID != configID {
		return fmt.Errorf("the ID of the configuration has changed during update")
	}
	return nil
}

func (me *service) Delete(ctx context.Context, id string) error {
	name, configID := splitID(id)
	client := rest.APITokenClient(me.credentials)
	if err := client.Delete(ctx, fmt.Sprintf("/api/v2/extensions/%s/monitoringConfigurations/%s", url.PathEscape(name), url.PathEscape(configID)), 200).Finish(nil); err != nil {
		// Potential response when the configuration contains
		//    import {
		//     ...
		//    }
		// with an invalid ID
		if strings.Contains(err.Error(), "Version property invalid. Required format is: 1.0.0") {
			return nil
		}
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
	NextPageKey string `json:"nextPageKey"`
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
