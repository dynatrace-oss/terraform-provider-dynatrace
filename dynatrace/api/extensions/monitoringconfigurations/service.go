/*
 * @license
 * Copyright 2026 Dynatrace LLC
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package monitoringconfigurations

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api"
	serviceSettings "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/extensions/monitoringconfigurations/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	coreapi "github.com/dynatrace/dynatrace-configuration-as-code-core/api"
	coreextensions "github.com/dynatrace/dynatrace-configuration-as-code-core/clients/extensions"
)

var invalidIdErr = errors.New("invalid id. Expected format: `{extensionName}#-#{configurationId}`")

// ExtensionClient defines the interface for interacting with the Extensions API.
type ExtensionClient interface {
	// ListExtensions returns all extensions.
	ListExtensions(ctx context.Context) (coreapi.PagedListResponse, error)

	// ListMonitoringConfigurations returns all monitoring configurations for a given extension.
	ListMonitoringConfigurations(ctx context.Context, extensionName string) (coreapi.PagedListResponse, error)

	// GetMonitoringConfiguration returns a specific monitoring configuration by extension name and configuration ID.
	GetMonitoringConfiguration(ctx context.Context, extensionName string, configurationID string) (coreapi.Response, error)

	// CreateMonitoringConfiguration creates a new monitoring configuration for a given extension.
	CreateMonitoringConfiguration(ctx context.Context, extensionName string, data []byte) (coreapi.Response, error)

	// UpdateMonitoringConfiguration updates an existing monitoring configuration.
	UpdateMonitoringConfiguration(ctx context.Context, extensionName string, configurationID string, data []byte) (coreapi.Response, error)

	// DeleteMonitoringConfiguration deletes a specific monitoring configuration.
	DeleteMonitoringConfiguration(ctx context.Context, extensionName string, configurationID string) error
}

type service struct {
	clientGetter func(ctx context.Context, credentials *rest.Credentials) (ExtensionClient, error)
	credentials  *rest.Credentials
}

func Service(credentials *rest.Credentials) settings.CRUDService[*serviceSettings.Settings] {
	return &service{clientGetter: createCoreClient, credentials: credentials}
}

func ServiceWithClientGetter(clientGetter func(ctx context.Context, credentials *rest.Credentials) (ExtensionClient, error), credentials *rest.Credentials) settings.CRUDService[*serviceSettings.Settings] {
	return &service{clientGetter: clientGetter, credentials: credentials}
}

func createCoreClient(ctx context.Context, credentials *rest.Credentials) (ExtensionClient, error) {
	platformClient, err := rest.CreatePlatformClient(ctx, credentials.OAuth.EnvironmentURL, credentials)
	if err != nil {
		return nil, err
	}
	return coreextensions.NewClient(platformClient), nil
}

func (s *service) SchemaID() string {
	return "extensions:v2:monitoringconfigurations"
}

func (s *service) Get(ctx context.Context, id string, v *serviceSettings.Settings) error {
	client, err := s.clientGetter(ctx, s.credentials)
	if err != nil {
		return err
	}

	extensionName, configID, err := splitID(id)
	if err != nil {
		return err
	}
	response, err := client.GetMonitoringConfiguration(ctx, extensionName, configID)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(response.Data, &v); err != nil {
		return err
	}
	v.Name = extensionName
	cfg := ctx.Value(settings.ContextKeyStateConfig)
	stateConfig, _ := cfg.(*serviceSettings.Settings)

	if stateConfig != nil && v.Value != nil {
		hcl.ReplaceCredPlaceholders(stateConfig.Value, v.Value)
	}
	return nil
}

func (s *service) List(ctx context.Context) (api.Stubs, error) {
	client, err := s.clientGetter(ctx, s.credentials)
	if err != nil {
		return nil, err
	}

	extensionsResponse, err := client.ListExtensions(ctx)
	if err != nil {
		return nil, err
	}

	stubs := api.Stubs{}
	type extensionResponse struct {
		ExtensionName string `json:"extensionName"`
	}
	for _, extensionBody := range extensionsResponse.All() {
		var extension extensionResponse
		if err := json.Unmarshal(extensionBody, &extension); err != nil {
			return nil, err
		}

		configurationsResponse, err := client.ListMonitoringConfigurations(ctx, extension.ExtensionName)
		if err != nil {
			return nil, err
		}

		type monitoringConfigurationResponse struct {
			ConfigurationID string `json:"objectId"`
		}
		for _, config := range configurationsResponse.All() {
			var configuration monitoringConfigurationResponse
			if err := json.Unmarshal(config, &configuration); err != nil {
				return nil, err
			}
			joinedID := joinID(extension.ExtensionName, configuration.ConfigurationID)
			stubs = append(stubs, &api.Stub{
				ID:   joinedID,
				Name: joinedID,
			})
		}
	}

	return stubs, nil
}

func (s *service) Create(ctx context.Context, v *serviceSettings.Settings) (stub *api.Stub, err error) {
	client, err := s.clientGetter(ctx, s.credentials)
	if err != nil {
		return nil, err
	}
	body, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	type createResponse struct {
		ID string `json:"objectId"`
	}

	response, err := client.CreateMonitoringConfiguration(ctx, v.Name, body)
	if err != nil {
		return nil, err
	}
	var responseData createResponse
	if err := json.Unmarshal(response.Data, &responseData); err != nil {
		return nil, err
	}
	return &api.Stub{ID: joinID(v.Name, responseData.ID)}, nil
}

func (s *service) Update(ctx context.Context, id string, v *serviceSettings.Settings) error {
	client, err := s.clientGetter(ctx, s.credentials)
	if err != nil {
		return err
	}
	body, err := json.Marshal(v)
	if err != nil {
		return err
	}
	extensionName, configID, err := splitID(id)
	if err != nil {
		return err
	}
	_, err = client.UpdateMonitoringConfiguration(ctx, extensionName, configID, body)
	return err
}

func (s *service) Delete(ctx context.Context, id string) error {
	client, err := s.clientGetter(ctx, s.credentials)
	if err != nil {
		return err
	}
	extensionName, configID, err := splitID(id)
	if err != nil {
		return err
	}
	return client.DeleteMonitoringConfiguration(ctx, extensionName, configID)
}

func joinID(extensionName string, configurationID string) string {
	return fmt.Sprintf("%s#-#%s", extensionName, configurationID)
}

func splitID(id string) (extensionName string, configurationID string, err error) {
	var found bool
	extensionName, configurationID, found = strings.Cut(id, "#-#")
	if !found {
		return "", "", invalidIdErr
	}
	return
}
