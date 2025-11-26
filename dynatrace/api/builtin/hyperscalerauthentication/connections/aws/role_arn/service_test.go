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

package role_arn

import (
	"context"
	"errors"
	"testing"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api"
	role_arn "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/hyperscalerauthentication/connections/aws/role_arn/settings"
	aws "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/hyperscalerauthentication/connections/aws/settings"
	"github.com/stretchr/testify/assert"
)

// mockCRUDService implements settings.CRUDService for testing
type mockCRUDService struct {
	getFunc    func(ctx context.Context, id string, v *aws.Settings) error
	updateFunc func(ctx context.Context, id string, v *aws.Settings) error
	listFunc   func(ctx context.Context) (api.Stubs, error)
}

func (m *mockCRUDService) Get(ctx context.Context, id string, v *aws.Settings) error {
	return m.getFunc(ctx, id, v)
}
func (m *mockCRUDService) Update(ctx context.Context, id string, v *aws.Settings) error {
	return m.updateFunc(ctx, id, v)
}
func (m *mockCRUDService) List(_ context.Context) (api.Stubs, error) {
	return api.Stubs{}, nil
}
func (m *mockCRUDService) Create(_ context.Context, _ *aws.Settings) (*api.Stub, error) {
	return nil, nil
}
func (m *mockCRUDService) Delete(_ context.Context, _ string) error { return nil }
func (m *mockCRUDService) SchemaID() string                         { return "" }

func TestService_Create(t *testing.T) {
	const connID = "conn-123"
	roleARN := "arn:aws:iam::123456789012:role/test"
	name := "TestConn"

	t.Run("success AWSRoleBasedAuthentication", func(t *testing.T) {
		mockConn := &aws.Settings{
			Type:                       aws.Types.AWSRoleBasedAuthentication,
			AWSRoleBasedAuthentication: &aws.AwsRoleBasedAuthenticationConfig{RoleARN: ""},
			Name:                       name,
		}
		mock := &mockCRUDService{
			getFunc: func(ctx context.Context, id string, v *aws.Settings) error {
				*v = *mockConn
				return nil
			},
			updateFunc: func(ctx context.Context, id string, v *aws.Settings) error {
				assert.Equal(t, roleARN, v.AWSRoleBasedAuthentication.RoleARN)
				return nil
			},
		}
		svc := &service{connService: mock}
		input := &role_arn.Settings{AWSConnectionID: connID, RoleARN: roleARN}

		stub, err := svc.Create(context.Background(), input)
		assert.NoError(t, err)
		assert.Equal(t, connID, stub.ID)
		assert.Equal(t, connID, stub.Name)
	})

	t.Run("success AWSWebIdentity", func(t *testing.T) {
		mockConn := &aws.Settings{
			Type:           aws.Types.AWSWebIdentity,
			AWSWebIdentity: &aws.AWSWebIdentity{RoleARN: ""},
			Name:           name,
		}
		mock := &mockCRUDService{
			getFunc: func(ctx context.Context, id string, v *aws.Settings) error {
				*v = *mockConn
				return nil
			},
			updateFunc: func(ctx context.Context, id string, v *aws.Settings) error {
				assert.Equal(t, roleARN, v.AWSWebIdentity.RoleARN)
				return nil
			},
		}
		svc := &service{connService: mock}
		input := &role_arn.Settings{AWSConnectionID: connID, RoleARN: roleARN}

		stub, err := svc.Create(context.Background(), input)

		assert.NoError(t, err)
		assert.Equal(t, connID, stub.ID)
		assert.Equal(t, connID, stub.Name)
	})

	t.Run("role ARN already set AWSRoleBasedAuthentication", func(t *testing.T) {
		mockConn := &aws.Settings{
			Type:                       aws.Types.AWSRoleBasedAuthentication,
			AWSRoleBasedAuthentication: &aws.AwsRoleBasedAuthenticationConfig{RoleARN: "other-arn"},
			Name:                       name,
		}
		mock := &mockCRUDService{
			getFunc: func(ctx context.Context, id string, v *aws.Settings) error {
				*v = *mockConn
				return nil
			},
		}
		svc := &service{connService: mock}
		input := &role_arn.Settings{AWSConnectionID: connID, RoleARN: roleARN}

		_, err := svc.Create(context.Background(), input)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "already set")
	})

	t.Run("role ARN already set but equal AWSRoleBasedAuthentication", func(t *testing.T) {
		mockConn := &aws.Settings{
			Type:                       aws.Types.AWSRoleBasedAuthentication,
			AWSRoleBasedAuthentication: &aws.AwsRoleBasedAuthenticationConfig{RoleARN: roleARN},
			Name:                       name,
		}
		mock := &mockCRUDService{
			getFunc: func(ctx context.Context, id string, v *aws.Settings) error {
				*v = *mockConn
				return nil
			},
			updateFunc: func(ctx context.Context, id string, v *aws.Settings) error {
				assert.Equal(t, roleARN, v.AWSRoleBasedAuthentication.RoleARN)
				return nil
			},
		}
		svc := &service{connService: mock}
		input := &role_arn.Settings{AWSConnectionID: connID, RoleARN: roleARN}

		stub, err := svc.Create(context.Background(), input)

		assert.NoError(t, err)
		assert.Equal(t, connID, stub.ID)
		assert.Equal(t, connID, stub.Name)
	})

	t.Run("role ARN already set AWSWebIdentity", func(t *testing.T) {
		mockConn := &aws.Settings{
			Type:           aws.Types.AWSWebIdentity,
			AWSWebIdentity: &aws.AWSWebIdentity{RoleARN: "other-arn"},
			Name:           name,
		}
		mock := &mockCRUDService{
			getFunc: func(ctx context.Context, id string, v *aws.Settings) error {
				*v = *mockConn
				return nil
			},
		}
		svc := &service{connService: mock}
		input := &role_arn.Settings{AWSConnectionID: connID, RoleARN: roleARN}

		_, err := svc.Create(context.Background(), input)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "already set")
	})

	t.Run("role ARN already set but equal AWSWebIdentity", func(t *testing.T) {
		mockConn := &aws.Settings{
			Type:           aws.Types.AWSWebIdentity,
			AWSWebIdentity: &aws.AWSWebIdentity{RoleARN: roleARN},
			Name:           name,
		}
		mock := &mockCRUDService{
			getFunc: func(ctx context.Context, id string, v *aws.Settings) error {
				*v = *mockConn
				return nil
			},
			updateFunc: func(ctx context.Context, id string, v *aws.Settings) error {
				assert.Equal(t, roleARN, v.AWSWebIdentity.RoleARN)
				return nil
			},
		}
		svc := &service{connService: mock}
		input := &role_arn.Settings{AWSConnectionID: connID, RoleARN: roleARN}

		stub, err := svc.Create(context.Background(), input)

		assert.NoError(t, err)
		assert.Equal(t, connID, stub.ID)
		assert.Equal(t, connID, stub.Name)
	})

	t.Run("connService.Get error", func(t *testing.T) {
		mock := &mockCRUDService{
			getFunc: func(ctx context.Context, id string, v *aws.Settings) error {
				return errors.New("get error")
			},
		}
		svc := &service{connService: mock}
		input := &role_arn.Settings{AWSConnectionID: connID, RoleARN: roleARN}

		_, err := svc.Create(context.Background(), input)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "get error")
	})

	t.Run("connService.Update error", func(t *testing.T) {
		mockConn := &aws.Settings{
			Type:                       aws.Types.AWSRoleBasedAuthentication,
			AWSRoleBasedAuthentication: &aws.AwsRoleBasedAuthenticationConfig{RoleARN: ""},
			Name:                       name,
		}
		mock := &mockCRUDService{
			getFunc: func(ctx context.Context, id string, v *aws.Settings) error {
				*v = *mockConn
				return nil
			},
			updateFunc: func(ctx context.Context, id string, v *aws.Settings) error {
				return errors.New("update error")
			},
		}
		svc := &service{connService: mock}
		input := &role_arn.Settings{AWSConnectionID: connID, RoleARN: roleARN}

		_, err := svc.Create(context.Background(), input)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "update error")
	})
}
