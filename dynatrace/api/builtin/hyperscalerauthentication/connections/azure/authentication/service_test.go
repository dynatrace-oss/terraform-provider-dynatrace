//go:build unit

/**
* @license
* Copyright 2025 Dynatrace LLC
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

package connectionauthentication

import (
	"context"
	"errors"
	"testing"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api"

	connectionauthentication_settings "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/hyperscalerauthentication/connections/azure/authentication/settings"
	azure "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/hyperscalerauthentication/connections/azure/settings"
	"github.com/stretchr/testify/assert"
)

// mockCRUDService implements settings.CRUDService for testing
type mockCRUDService struct {
	getFunc    func(ctx context.Context, id string, v *azure.Settings) error
	updateFunc func(ctx context.Context, id string, v *azure.Settings) error
	listFunc   func(ctx context.Context) (api.Stubs, error)
}

func (m *mockCRUDService) Get(ctx context.Context, id string, v *azure.Settings) error {
	return m.getFunc(ctx, id, v)
}
func (m *mockCRUDService) Update(ctx context.Context, id string, v *azure.Settings) error {
	return m.updateFunc(ctx, id, v)
}
func (m *mockCRUDService) List(_ context.Context) (api.Stubs, error) {
	return api.Stubs{}, nil
}
func (m *mockCRUDService) Create(_ context.Context, _ *azure.Settings) (*api.Stub, error) {
	return nil, nil
}
func (m *mockCRUDService) Delete(_ context.Context, _ string) error { return nil }
func (m *mockCRUDService) SchemaID() string                         { return "" }

// TestService_Create tests the Create method of the dynatrace_azure_connection_authentication resource service
func TestService_Create(t *testing.T) {
	const connID = "conn-123"
	name := "TestConn"
	applicationID := "app-123"
	directoryID := "dir-123"

	t.Run("success configuring federated identity credentials", func(t *testing.T) {
		mockConn := &azure.Settings{
			Type:                        azure.Types.Federatedidentitycredential,
			FederatedIdentityCredential: &azure.FederatedIdentityCredential{Consumers: []azure.ConsumersOfFederatedIdentityCredential{"DA"}},
			Name:                        name,
		}
		mock := &mockCRUDService{
			getFunc: func(ctx context.Context, id string, v *azure.Settings) error {
				*v = *mockConn
				return nil
			},
			updateFunc: func(ctx context.Context, id string, v *azure.Settings) error {
				assert.NotNil(t, v.FederatedIdentityCredential)
				assert.Equal(t, applicationID, *v.FederatedIdentityCredential.ApplicationID)
				assert.Equal(t, directoryID, *v.FederatedIdentityCredential.DirectoryID)
				return nil
			},
		}
		svc := &service{connService: mock}
		input := &connectionauthentication_settings.Settings{AzureConnectionID: connID, ApplicationID: &applicationID, DirectoryID: &directoryID}

		stub, err := svc.Create(context.Background(), input)

		assert.NoError(t, err)
		assert.Equal(t, connID, stub.ID)
		assert.Equal(t, connID, stub.Name)
	})

	t.Run("failure configuring client secret", func(t *testing.T) {
		mockConn := &azure.Settings{
			Type:         azure.Types.Clientsecret,
			ClientSecret: &azure.ClientSecretConfig{},
			Name:         name,
		}

		mock := &mockCRUDService{
			getFunc: func(ctx context.Context, id string, v *azure.Settings) error {
				*v = *mockConn
				return nil
			},
		}
		svc := &service{connService: mock}
		input := &connectionauthentication_settings.Settings{AzureConnectionID: connID, ApplicationID: &applicationID, DirectoryID: &directoryID}

		_, err := svc.Create(context.Background(), input)

		assert.ErrorContains(t, err, "not configured")
	})

	t.Run("federated identity credentials already set", func(t *testing.T) {
		mockConn := &azure.Settings{
			Type:                        azure.Types.Federatedidentitycredential,
			FederatedIdentityCredential: &azure.FederatedIdentityCredential{Consumers: []azure.ConsumersOfFederatedIdentityCredential{"DA"}, ApplicationID: &applicationID, DirectoryID: &directoryID},
			Name:                        name,
		}
		mock := &mockCRUDService{
			getFunc: func(ctx context.Context, id string, v *azure.Settings) error {
				*v = *mockConn
				return nil
			},
		}

		applicationID2 := "app-456"

		svc := &service{connService: mock}
		input := &connectionauthentication_settings.Settings{AzureConnectionID: connID, ApplicationID: &applicationID2, DirectoryID: &directoryID}

		_, err := svc.Create(context.Background(), input)

		assert.ErrorContains(t, err, "immutable after creation")
	})

	t.Run("federated identity credentials already set but equal", func(t *testing.T) {
		mockConn := &azure.Settings{
			Type:                        azure.Types.Federatedidentitycredential,
			FederatedIdentityCredential: &azure.FederatedIdentityCredential{Consumers: []azure.ConsumersOfFederatedIdentityCredential{"DA"}, ApplicationID: &applicationID, DirectoryID: &directoryID},
			Name:                        name,
		}
		mock := &mockCRUDService{
			getFunc: func(ctx context.Context, id string, v *azure.Settings) error {
				*v = *mockConn
				return nil
			},
		}
		svc := &service{connService: mock}
		input := &connectionauthentication_settings.Settings{AzureConnectionID: connID, ApplicationID: &applicationID, DirectoryID: &directoryID}

		stub, err := svc.Create(context.Background(), input)

		assert.NoError(t, err)
		assert.Equal(t, connID, stub.ID)
		assert.Equal(t, connID, stub.Name)
	})

	t.Run("connService.Get error", func(t *testing.T) {
		mock := &mockCRUDService{
			getFunc: func(ctx context.Context, id string, v *azure.Settings) error {
				return errors.New("get error")
			},
		}
		svc := &service{connService: mock}
		input := &connectionauthentication_settings.Settings{AzureConnectionID: connID, ApplicationID: &applicationID, DirectoryID: &directoryID}

		_, err := svc.Create(context.Background(), input)

		assert.ErrorContains(t, err, "get error")
	})

	t.Run("connService.Update error", func(t *testing.T) {
		mockConn := &azure.Settings{
			Type:                        azure.Types.Federatedidentitycredential,
			FederatedIdentityCredential: &azure.FederatedIdentityCredential{Consumers: []azure.ConsumersOfFederatedIdentityCredential{"DA"}},
			Name:                        name,
		}
		mock := &mockCRUDService{
			getFunc: func(ctx context.Context, id string, v *azure.Settings) error {
				*v = *mockConn
				return nil
			},
			updateFunc: func(ctx context.Context, id string, v *azure.Settings) error {
				return errors.New("update error")
			},
		}
		svc := &service{connService: mock}
		input := &connectionauthentication_settings.Settings{AzureConnectionID: connID, ApplicationID: &applicationID, DirectoryID: &directoryID}

		_, err := svc.Create(context.Background(), input)

		assert.ErrorContains(t, err, "update error")
	})
}
