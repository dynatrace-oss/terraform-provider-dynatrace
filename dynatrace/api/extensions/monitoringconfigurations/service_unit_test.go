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

package monitoringconfigurations_test

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/extensions/monitoringconfigurations"
	serviceSettings "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/extensions/monitoringconfigurations/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	coreapi "github.com/dynatrace/dynatrace-configuration-as-code-core/api"
)

// mockExtensionClient implements ExtensionClient for testing.
type mockExtensionClient struct {
	listExtensionsFn                func(ctx context.Context) (coreapi.PagedListResponse, error)
	listMonitoringConfigurationsFn  func(ctx context.Context, extensionName string) (coreapi.PagedListResponse, error)
	getMonitoringConfigurationFn    func(ctx context.Context, extensionName string, configurationID string) (coreapi.Response, error)
	createMonitoringConfigurationFn func(ctx context.Context, extensionName string, data []byte) (coreapi.Response, error)
	updateMonitoringConfigurationFn func(ctx context.Context, extensionName string, configurationID string, data []byte) (coreapi.Response, error)
	deleteMonitoringConfigurationFn func(ctx context.Context, extensionName string, configurationID string) error
}

func (m *mockExtensionClient) ListExtensions(ctx context.Context) (coreapi.PagedListResponse, error) {
	return m.listExtensionsFn(ctx)
}
func (m *mockExtensionClient) ListMonitoringConfigurations(ctx context.Context, extensionName string) (coreapi.PagedListResponse, error) {
	return m.listMonitoringConfigurationsFn(ctx, extensionName)
}
func (m *mockExtensionClient) GetMonitoringConfiguration(ctx context.Context, extensionName string, configurationID string) (coreapi.Response, error) {
	return m.getMonitoringConfigurationFn(ctx, extensionName, configurationID)
}
func (m *mockExtensionClient) CreateMonitoringConfiguration(ctx context.Context, extensionName string, data []byte) (coreapi.Response, error) {
	return m.createMonitoringConfigurationFn(ctx, extensionName, data)
}
func (m *mockExtensionClient) UpdateMonitoringConfiguration(ctx context.Context, extensionName string, configurationID string, data []byte) (coreapi.Response, error) {
	return m.updateMonitoringConfigurationFn(ctx, extensionName, configurationID, data)
}
func (m *mockExtensionClient) DeleteMonitoringConfiguration(ctx context.Context, extensionName string, configurationID string) error {
	return m.deleteMonitoringConfigurationFn(ctx, extensionName, configurationID)
}

func mockClientGetter(client *mockExtensionClient) func(ctx context.Context, credentials *rest.Credentials) (monitoringconfigurations.ExtensionClient, error) {
	return func(ctx context.Context, credentials *rest.Credentials) (monitoringconfigurations.ExtensionClient, error) {
		return client, nil
	}
}

func failingClientGetter(err error) func(ctx context.Context, credentials *rest.Credentials) (monitoringconfigurations.ExtensionClient, error) {
	return func(ctx context.Context, credentials *rest.Credentials) (monitoringconfigurations.ExtensionClient, error) {
		return nil, err
	}
}

func pagedResponse(objects ...[]byte) coreapi.PagedListResponse {
	return coreapi.PagedListResponse{
		{Objects: objects},
	}
}

func TestService_Get(t *testing.T) {
	t.Run("Returns error when client creation fails", func(t *testing.T) {
		svc := monitoringconfigurations.ServiceWithClientGetter(failingClientGetter(assert.AnError), &rest.Credentials{})
		err := svc.Get(t.Context(), "com.example.ext#-#cfg-1", &serviceSettings.Settings{})
		assert.ErrorIs(t, err, assert.AnError)
	})

	t.Run("Returns error when GetMonitoringConfiguration fails", func(t *testing.T) {
		mock := &mockExtensionClient{
			getMonitoringConfigurationFn: func(ctx context.Context, extensionName string, configurationID string) (coreapi.Response, error) {
				return coreapi.Response{}, assert.AnError
			},
		}
		svc := monitoringconfigurations.ServiceWithClientGetter(mockClientGetter(mock), &rest.Credentials{})
		err := svc.Get(t.Context(), "com.example.ext#-#cfg-1", &serviceSettings.Settings{})
		assert.ErrorIs(t, err, assert.AnError)
	})

	t.Run("Returns error on invalid response JSON", func(t *testing.T) {
		mock := &mockExtensionClient{
			getMonitoringConfigurationFn: func(ctx context.Context, extensionName string, configurationID string) (coreapi.Response, error) {
				return coreapi.Response{Data: []byte("not-json")}, nil
			},
		}
		svc := monitoringconfigurations.ServiceWithClientGetter(mockClientGetter(mock), &rest.Credentials{})
		err := svc.Get(t.Context(), "com.example.ext#-#cfg-1", &serviceSettings.Settings{})
		assert.Error(t, err)
	})

	t.Run("Splits compound ID and populates Settings on success", func(t *testing.T) {
		responseData, _ := json.Marshal(map[string]any{
			"scope": "HOST-ABC123",
			"value": map[string]any{"key": "val"},
		})
		var capturedExtName, capturedCfgID string
		mock := &mockExtensionClient{
			getMonitoringConfigurationFn: func(ctx context.Context, extensionName string, configurationID string) (coreapi.Response, error) {
				capturedExtName = extensionName
				capturedCfgID = configurationID
				return coreapi.Response{Data: responseData}, nil
			},
		}
		svc := monitoringconfigurations.ServiceWithClientGetter(mockClientGetter(mock), &rest.Credentials{})
		v := &serviceSettings.Settings{}
		err := svc.Get(t.Context(), "com.example.ext#-#cfg-1", v)
		require.NoError(t, err)
		assert.Equal(t, "com.example.ext", capturedExtName)
		assert.Equal(t, "cfg-1", capturedCfgID)
		assert.Equal(t, "com.example.ext", v.Name)
		assert.Equal(t, "HOST-ABC123", v.Scope)
	})
}

func TestService_List(t *testing.T) {
	t.Run("Returns error when client creation fails", func(t *testing.T) {
		svc := monitoringconfigurations.ServiceWithClientGetter(failingClientGetter(assert.AnError), &rest.Credentials{})
		_, err := svc.List(t.Context())
		assert.ErrorIs(t, err, assert.AnError)
	})

	t.Run("Returns error when ListExtensions fails", func(t *testing.T) {
		mock := &mockExtensionClient{
			listExtensionsFn: func(ctx context.Context) (coreapi.PagedListResponse, error) {
				return nil, assert.AnError
			},
		}
		svc := monitoringconfigurations.ServiceWithClientGetter(mockClientGetter(mock), &rest.Credentials{})
		_, err := svc.List(t.Context())
		assert.ErrorIs(t, err, assert.AnError)
	})

	t.Run("Returns error on invalid extension JSON", func(t *testing.T) {
		mock := &mockExtensionClient{
			listExtensionsFn: func(ctx context.Context) (coreapi.PagedListResponse, error) {
				return pagedResponse([]byte("bad-json")), nil
			},
		}
		svc := monitoringconfigurations.ServiceWithClientGetter(mockClientGetter(mock), &rest.Credentials{})
		_, err := svc.List(t.Context())
		assert.Error(t, err)
	})

	t.Run("Returns error when ListMonitoringConfigurations fails", func(t *testing.T) {
		extJSON, _ := json.Marshal(map[string]string{"extensionName": "com.example.ext"})
		mock := &mockExtensionClient{
			listExtensionsFn: func(ctx context.Context) (coreapi.PagedListResponse, error) {
				return pagedResponse(extJSON), nil
			},
			listMonitoringConfigurationsFn: func(ctx context.Context, extensionName string) (coreapi.PagedListResponse, error) {
				return nil, assert.AnError
			},
		}
		svc := monitoringconfigurations.ServiceWithClientGetter(mockClientGetter(mock), &rest.Credentials{})
		_, err := svc.List(t.Context())
		assert.ErrorIs(t, err, assert.AnError)
	})

	t.Run("Returns error on invalid configuration JSON", func(t *testing.T) {
		extJSON, _ := json.Marshal(map[string]string{"extensionName": "com.example.ext"})
		mock := &mockExtensionClient{
			listExtensionsFn: func(ctx context.Context) (coreapi.PagedListResponse, error) {
				return pagedResponse(extJSON), nil
			},
			listMonitoringConfigurationsFn: func(ctx context.Context, extensionName string) (coreapi.PagedListResponse, error) {
				return pagedResponse([]byte("bad-json")), nil
			},
		}
		svc := monitoringconfigurations.ServiceWithClientGetter(mockClientGetter(mock), &rest.Credentials{})
		_, err := svc.List(t.Context())
		assert.Error(t, err)
	})

	t.Run("Returns stubs with joined IDs on success", func(t *testing.T) {
		ext1JSON, _ := json.Marshal(map[string]string{"extensionName": "com.example.ext1"})
		ext2JSON, _ := json.Marshal(map[string]string{"extensionName": "com.example.ext2"})
		cfg1JSON, _ := json.Marshal(map[string]string{"objectId": "cfg-1"})
		cfg2JSON, _ := json.Marshal(map[string]string{"objectId": "cfg-2"})
		mock := &mockExtensionClient{
			listExtensionsFn: func(ctx context.Context) (coreapi.PagedListResponse, error) {
				return pagedResponse(ext1JSON, ext2JSON), nil
			},
			listMonitoringConfigurationsFn: func(ctx context.Context, extensionName string) (coreapi.PagedListResponse, error) {
				switch extensionName {
				case "com.example.ext1":
					return pagedResponse(cfg1JSON), nil
				case "com.example.ext2":
					return pagedResponse(cfg2JSON), nil
				}
				return pagedResponse(), nil
			},
		}
		svc := monitoringconfigurations.ServiceWithClientGetter(mockClientGetter(mock), &rest.Credentials{})
		stubs, err := svc.List(t.Context())
		require.NoError(t, err)
		require.Len(t, stubs, 2)
		assert.Equal(t, "com.example.ext1#-#cfg-1", stubs[0].ID)
		assert.Equal(t, "com.example.ext1#-#cfg-1", stubs[0].Name)
		assert.Equal(t, "com.example.ext2#-#cfg-2", stubs[1].ID)
		assert.Equal(t, "com.example.ext2#-#cfg-2", stubs[1].Name)
	})

	t.Run("Returns empty stubs when no extensions exist", func(t *testing.T) {
		mock := &mockExtensionClient{
			listExtensionsFn: func(ctx context.Context) (coreapi.PagedListResponse, error) {
				return pagedResponse(), nil
			},
		}
		svc := monitoringconfigurations.ServiceWithClientGetter(mockClientGetter(mock), &rest.Credentials{})
		stubs, err := svc.List(t.Context())
		require.NoError(t, err)
		assert.Empty(t, stubs)
	})

	t.Run("Returns empty stubs when extensions have no configurations", func(t *testing.T) {
		extJSON, _ := json.Marshal(map[string]string{"extensionName": "com.example.ext"})
		mock := &mockExtensionClient{
			listExtensionsFn: func(ctx context.Context) (coreapi.PagedListResponse, error) {
				return pagedResponse(extJSON), nil
			},
			listMonitoringConfigurationsFn: func(ctx context.Context, extensionName string) (coreapi.PagedListResponse, error) {
				return pagedResponse(), nil
			},
		}
		svc := monitoringconfigurations.ServiceWithClientGetter(mockClientGetter(mock), &rest.Credentials{})
		stubs, err := svc.List(t.Context())
		require.NoError(t, err)
		assert.Empty(t, stubs)
	})
}

func TestService_Create(t *testing.T) {
	t.Run("Returns error when client creation fails", func(t *testing.T) {
		svc := monitoringconfigurations.ServiceWithClientGetter(failingClientGetter(assert.AnError), &rest.Credentials{})
		_, err := svc.Create(t.Context(), &serviceSettings.Settings{Name: "com.example.ext"})
		assert.ErrorIs(t, err, assert.AnError)
	})

	t.Run("Returns error when CreateMonitoringConfiguration fails", func(t *testing.T) {
		mock := &mockExtensionClient{
			createMonitoringConfigurationFn: func(ctx context.Context, extensionName string, data []byte) (coreapi.Response, error) {
				return coreapi.Response{}, assert.AnError
			},
		}
		svc := monitoringconfigurations.ServiceWithClientGetter(mockClientGetter(mock), &rest.Credentials{})
		_, err := svc.Create(t.Context(), &serviceSettings.Settings{
			Name:  "com.example.ext",
			Scope: "HOST-ABC123",
		})
		assert.ErrorIs(t, err, assert.AnError)
	})

	t.Run("Returns error on invalid response JSON", func(t *testing.T) {
		mock := &mockExtensionClient{
			createMonitoringConfigurationFn: func(ctx context.Context, extensionName string, data []byte) (coreapi.Response, error) {
				return coreapi.Response{Data: []byte("bad-json")}, nil
			},
		}
		svc := monitoringconfigurations.ServiceWithClientGetter(mockClientGetter(mock), &rest.Credentials{})
		_, err := svc.Create(t.Context(), &serviceSettings.Settings{
			Name:  "com.example.ext",
			Scope: "HOST-ABC123",
		})
		assert.Error(t, err)
	})

	t.Run("Sends correct payload and returns stub with joined ID on success", func(t *testing.T) {
		var capturedExtName string
		var capturedPayload serviceSettings.Settings
		mock := &mockExtensionClient{
			createMonitoringConfigurationFn: func(ctx context.Context, extensionName string, data []byte) (coreapi.Response, error) {
				capturedExtName = extensionName
				err := json.Unmarshal(data, &capturedPayload)
				require.NoError(t, err)
				respJSON, _ := json.Marshal(map[string]string{"objectId": "new-cfg-id"})
				return coreapi.Response{Data: respJSON}, nil
			},
		}
		svc := monitoringconfigurations.ServiceWithClientGetter(mockClientGetter(mock), &rest.Credentials{})
		stub, err := svc.Create(t.Context(), &serviceSettings.Settings{
			Name:  "com.example.ext",
			Scope: "HOST-ABC123",
			Value: map[string]any{"key": "val"},
		})
		require.NoError(t, err)
		assert.Equal(t, "com.example.ext#-#new-cfg-id", stub.ID)
		assert.Equal(t, "com.example.ext", capturedExtName)
		assert.Equal(t, "HOST-ABC123", capturedPayload.Scope)
	})
}

func TestService_Update(t *testing.T) {
	t.Run("Returns error when client creation fails", func(t *testing.T) {
		svc := monitoringconfigurations.ServiceWithClientGetter(failingClientGetter(assert.AnError), &rest.Credentials{})
		err := svc.Update(t.Context(), "com.example.ext#-#cfg-1", &serviceSettings.Settings{})
		assert.ErrorIs(t, err, assert.AnError)
	})

	t.Run("Returns error when UpdateMonitoringConfiguration fails", func(t *testing.T) {
		mock := &mockExtensionClient{
			updateMonitoringConfigurationFn: func(ctx context.Context, extensionName string, configurationID string, data []byte) (coreapi.Response, error) {
				return coreapi.Response{}, assert.AnError
			},
		}
		svc := monitoringconfigurations.ServiceWithClientGetter(mockClientGetter(mock), &rest.Credentials{})
		err := svc.Update(t.Context(), "com.example.ext#-#cfg-1", &serviceSettings.Settings{})
		assert.ErrorIs(t, err, assert.AnError)
	})

	t.Run("Splits compound ID and sends correct payload on success", func(t *testing.T) {
		var capturedExtName, capturedCfgID string
		var capturedPayload serviceSettings.Settings
		mock := &mockExtensionClient{
			updateMonitoringConfigurationFn: func(ctx context.Context, extensionName string, configurationID string, data []byte) (coreapi.Response, error) {
				capturedExtName = extensionName
				capturedCfgID = configurationID
				err := json.Unmarshal(data, &capturedPayload)
				require.NoError(t, err)
				return coreapi.Response{}, nil
			},
		}
		svc := monitoringconfigurations.ServiceWithClientGetter(mockClientGetter(mock), &rest.Credentials{})
		err := svc.Update(t.Context(), "com.example.ext#-#cfg-1", &serviceSettings.Settings{
			Scope: "HOST-ABC123",
			Value: map[string]any{"key": "val"},
		})
		require.NoError(t, err)
		assert.Equal(t, "com.example.ext", capturedExtName)
		assert.Equal(t, "cfg-1", capturedCfgID)
		assert.Equal(t, "HOST-ABC123", capturedPayload.Scope)
	})
}

func TestService_Delete(t *testing.T) {
	t.Run("Returns error when client creation fails", func(t *testing.T) {
		svc := monitoringconfigurations.ServiceWithClientGetter(failingClientGetter(assert.AnError), &rest.Credentials{})
		err := svc.Delete(t.Context(), "com.example.ext#-#cfg-1")
		assert.ErrorIs(t, err, assert.AnError)
	})

	t.Run("Returns error when DeleteMonitoringConfiguration fails", func(t *testing.T) {
		mock := &mockExtensionClient{
			deleteMonitoringConfigurationFn: func(ctx context.Context, extensionName string, configurationID string) error {
				return assert.AnError
			},
		}
		svc := monitoringconfigurations.ServiceWithClientGetter(mockClientGetter(mock), &rest.Credentials{})
		err := svc.Delete(t.Context(), "com.example.ext#-#cfg-1")
		assert.ErrorIs(t, err, assert.AnError)
	})

	t.Run("Splits compound ID and returns nil on success", func(t *testing.T) {
		var capturedExtName, capturedCfgID string
		mock := &mockExtensionClient{
			deleteMonitoringConfigurationFn: func(ctx context.Context, extensionName string, configurationID string) error {
				capturedExtName = extensionName
				capturedCfgID = configurationID
				return nil
			},
		}
		svc := monitoringconfigurations.ServiceWithClientGetter(mockClientGetter(mock), &rest.Credentials{})
		err := svc.Delete(t.Context(), "com.example.ext#-#cfg-1")
		require.NoError(t, err)
		assert.Equal(t, "com.example.ext", capturedExtName)
		assert.Equal(t, "cfg-1", capturedCfgID)
	})
}
