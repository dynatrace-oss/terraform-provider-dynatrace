//go:build unit

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

package testing

import (
	"context"

	coreapi "github.com/dynatrace/dynatrace-configuration-as-code-core/api"
)

// MockExtensionClient implements ExtensionClient for testing.
type MockExtensionClient struct {
	ListExtensionsFn                func(ctx context.Context) (coreapi.PagedListResponse, error)
	ListMonitoringConfigurationsFn  func(ctx context.Context, extensionName string, filter string) (coreapi.PagedListResponse, error)
	GetMonitoringConfigurationFn    func(ctx context.Context, extensionName string, configurationID string) (coreapi.Response, error)
	CreateMonitoringConfigurationFn func(ctx context.Context, extensionName string, data []byte) (coreapi.Response, error)
	UpdateMonitoringConfigurationFn func(ctx context.Context, extensionName string, configurationID string, data []byte) (coreapi.Response, error)
	DeleteMonitoringConfigurationFn func(ctx context.Context, extensionName string, configurationID string) error
}

func (m *MockExtensionClient) ListExtensions(ctx context.Context) (coreapi.PagedListResponse, error) {
	return m.ListExtensionsFn(ctx)
}
func (m *MockExtensionClient) ListMonitoringConfigurations(ctx context.Context, extensionName string, filter string) (coreapi.PagedListResponse, error) {
	return m.ListMonitoringConfigurationsFn(ctx, extensionName, filter)
}
func (m *MockExtensionClient) GetMonitoringConfiguration(ctx context.Context, extensionName string, configurationID string) (coreapi.Response, error) {
	return m.GetMonitoringConfigurationFn(ctx, extensionName, configurationID)
}
func (m *MockExtensionClient) CreateMonitoringConfiguration(ctx context.Context, extensionName string, data []byte) (coreapi.Response, error) {
	return m.CreateMonitoringConfigurationFn(ctx, extensionName, data)
}
func (m *MockExtensionClient) UpdateMonitoringConfiguration(ctx context.Context, extensionName string, configurationID string, data []byte) (coreapi.Response, error) {
	return m.UpdateMonitoringConfigurationFn(ctx, extensionName, configurationID, data)
}
func (m *MockExtensionClient) DeleteMonitoringConfiguration(ctx context.Context, extensionName string, configurationID string) error {
	return m.DeleteMonitoringConfigurationFn(ctx, extensionName, configurationID)
}
