/**
* @license
* Copyright 2026 Dynatrace LLC
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

// Package awsmonitoring exposes a typed Terraform resource
// (dynatrace_aws_monitoring_configuration) backed by the generic Extensions
// 2.0 monitoring configuration API.
//
// The wire-level interactions mirror the existing generic
// `monitoringconfigurations` package — same endpoints, same client — but with
// the extension name hard-coded to `com.dynatrace.extension.da-aws` so the
// user does not have to know the extension at all.
package awsmonitoring

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/Masterminds/semver/v3"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api"
	serviceSettings "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/extensions/dac/awsmonitoring/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings"

	coreapi "github.com/dynatrace/dynatrace-configuration-as-code-core/api"
	coreextensions "github.com/dynatrace/dynatrace-configuration-as-code-core/clients/extensions"
)

// ExtensionClient is the subset of the upstream Extensions client we need.
// Mirrors the interface in the generic monitoringconfigurations package so
// tests can mock it.
type ExtensionClient interface {
	ListMonitoringConfigurations(ctx context.Context, extensionName string) (coreapi.PagedListResponse, error)
	GetMonitoringConfiguration(ctx context.Context, extensionName string, configurationID string) (coreapi.Response, error)
	CreateMonitoringConfiguration(ctx context.Context, extensionName string, data []byte) (coreapi.Response, error)
	UpdateMonitoringConfiguration(ctx context.Context, extensionName string, configurationID string, data []byte) (coreapi.Response, error)
	DeleteMonitoringConfiguration(ctx context.Context, extensionName string, configurationID string) error
	ListExtensionVersions(ctx context.Context, extensionName string) (coreapi.PagedListResponse, error)
}

type service struct {
	clientGetter func(ctx context.Context, credentials *rest.Credentials) (ExtensionClient, error)
	credentials  *rest.Credentials
}

// Service constructs the CRUD service registered in the resource descriptor
// for ResourceTypes.AWSMonitoringConfiguration.
func Service(credentials *rest.Credentials) settings.CRUDService[*serviceSettings.Settings] {
	return &service{clientGetter: createCoreClient, credentials: credentials}
}

// ServiceWithClientGetter is the test seam.
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
	return "extensions:v2:monitoringconfigurations:" + serviceSettings.AWSExtensionName
}

func (s *service) Get(ctx context.Context, id string, v *serviceSettings.Settings) error {
	client, err := s.clientGetter(ctx, s.credentials)
	if err != nil {
		return err
	}
	response, err := client.GetMonitoringConfiguration(ctx, serviceSettings.AWSExtensionName, id)
	if err != nil {
		return err
	}
	return json.Unmarshal(response.Data, v)
}

func (s *service) List(ctx context.Context) (api.Stubs, error) {
	client, err := s.clientGetter(ctx, s.credentials)
	if err != nil {
		return nil, err
	}
	resp, err := client.ListMonitoringConfigurations(ctx, serviceSettings.AWSExtensionName)
	if err != nil {
		return nil, err
	}
	type configEntry struct {
		ObjectID string `json:"objectId"`
		Value    struct {
			Description string `json:"description"`
		} `json:"value"`
	}
	stubs := api.Stubs{}
	for _, raw := range resp.All() {
		var entry configEntry
		if err := json.Unmarshal(raw, &entry); err != nil {
			return nil, err
		}
		name := entry.Value.Description
		if name == "" {
			name = entry.ObjectID
		}
		stubs = append(stubs, &api.Stub{ID: entry.ObjectID, Name: name})
	}
	return stubs, nil
}

func (s *service) Create(ctx context.Context, v *serviceSettings.Settings) (*api.Stub, error) {
	client, err := s.clientGetter(ctx, s.credentials)
	if err != nil {
		return nil, err
	}
	if v.ExtensionVersion == "" {
		resolved, err := resolveLatestExtensionVersion(ctx, client, serviceSettings.AWSExtensionName)
		if err != nil {
			return nil, err
		}
		v.ExtensionVersion = resolved
	}
	body, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	response, err := client.CreateMonitoringConfiguration(ctx, serviceSettings.AWSExtensionName, body)
	if err != nil {
		return nil, err
	}
	var created struct {
		ObjectID string `json:"objectId"`
	}
	if err := json.Unmarshal(response.Data, &created); err != nil {
		return nil, err
	}
	if created.ObjectID == "" {
		return nil, errors.New("create response did not contain objectId")
	}
	return &api.Stub{ID: created.ObjectID, Name: v.Name}, nil
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
	_, err = client.UpdateMonitoringConfiguration(ctx, serviceSettings.AWSExtensionName, id, body)
	return err
}

func (s *service) Delete(ctx context.Context, id string) error {
	client, err := s.clientGetter(ctx, s.credentials)
	if err != nil {
		return err
	}
	return client.DeleteMonitoringConfiguration(ctx, serviceSettings.AWSExtensionName, id)
}

// resolveLatestExtensionVersion mirrors dtctl's GetLatestVersion: list all
// installed versions of the extension and return the highest valid semver.
func resolveLatestExtensionVersion(ctx context.Context, client ExtensionClient, extensionName string) (string, error) {
	resp, err := client.ListExtensionVersions(ctx, extensionName)
	if err != nil {
		return "", fmt.Errorf("failed to list versions of extension %s: %w", extensionName, err)
	}
	type versionEntry struct {
		Version string `json:"version"`
	}
	var latest *semver.Version
	for _, raw := range resp.All() {
		var entry versionEntry
		if err := json.Unmarshal(raw, &entry); err != nil {
			return "", err
		}
		if entry.Version == "" {
			continue
		}
		parsed, err := semver.NewVersion(entry.Version)
		if err != nil {
			continue
		}
		if latest == nil || parsed.GreaterThan(latest) {
			latest = parsed
		}
	}
	if latest == nil {
		return "", fmt.Errorf("extension %s is not installed on this tenant — install it from the Hub first or set extension_version explicitly", extensionName)
	}
	return latest.String(), nil
}
