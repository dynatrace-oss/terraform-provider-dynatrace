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

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api"

	azureconnection "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/hyperscalerauthentication/connections/azure"
	connectionauthentication_settings "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/hyperscalerauthentication/connections/azure/authentication/settings"
	azureconnection_settings "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/hyperscalerauthentication/connections/azure/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings"
)

func Service(credentials *rest.Credentials) settings.CRUDService[*connectionauthentication_settings.Settings] {
	return &service{
		connService: azureconnection.Service(credentials),
	}
}

type service struct {
	connService settings.CRUDService[*azureconnection_settings.Settings]
}

func (me *service) List(ctx context.Context) (api.Stubs, error) {
	stubs, err := me.connService.List(ctx)
	if err != nil {
		return nil, err
	}

	var filteredStubs api.Stubs
	for _, stub := range stubs {
		connValue := azureconnection_settings.Settings{}
		if err := me.connService.Get(ctx, stub.ID, &connValue); err != nil {
			return nil, err
		}
		if connValue.Type == azureconnection_settings.Types.Federatedidentitycredential {
			filteredStubs = append(filteredStubs, stub)
		}
	}

	return filteredStubs, nil
}

func (me *service) Get(ctx context.Context, id string, v *connectionauthentication_settings.Settings) error {
	connValue := azureconnection_settings.Settings{}
	if err := me.connService.Get(ctx, id, &connValue); err != nil {
		return err
	}

	if connValue.Type != azureconnection_settings.Types.Federatedidentitycredential {
		return errors.New("the associated Azure connection is not configured for federated identity credential authentication")
	}

	if connValue.FederatedIdentityCredential != nil {
		v.ApplicationID = connValue.FederatedIdentityCredential.ApplicationID
		v.DirectoryID = connValue.FederatedIdentityCredential.DirectoryID
	}
	v.AzureConnectionID = id
	v.Name = connValue.Name
	return nil
}

func (me *service) SchemaID() string {
	return me.connService.SchemaID()
}

func (me *service) Create(ctx context.Context, v *connectionauthentication_settings.Settings) (*api.Stub, error) {
	connValue := azureconnection_settings.Settings{}
	if err := me.connService.Get(ctx, v.AzureConnectionID, &connValue); err != nil {
		return nil, err
	}

	if (connValue.Type != azureconnection_settings.Types.Federatedidentitycredential) || (connValue.FederatedIdentityCredential == nil) {
		return nil, errors.New("the associated Azure Connection is not configured for Federated Identity Credential authentication")
	}

	// no update needed if there are no changes
	if connValue.FederatedIdentityCredential.DirectoryID == v.DirectoryID && connValue.FederatedIdentityCredential.ApplicationID == v.ApplicationID {
		return &api.Stub{ID: v.AzureConnectionID, Name: v.AzureConnectionID}, nil
	}

	// it is not possible to change directory or application ID after they have been set
	if connValue.FederatedIdentityCredential.DirectoryID != nil || connValue.FederatedIdentityCredential.ApplicationID != nil {
		return nil, errors.New("update not supported: This resource is immutable after creation. Changes require `dynatrace_azure_connection` to be recreated")
	}

	connValue.FederatedIdentityCredential.DirectoryID = v.DirectoryID
	connValue.FederatedIdentityCredential.ApplicationID = v.ApplicationID

	if err := me.connService.Update(ctx, v.AzureConnectionID, &connValue); err != nil {
		return nil, err
	}

	return &api.Stub{ID: v.AzureConnectionID, Name: v.AzureConnectionID}, nil
}

func (me *service) Update(_ context.Context, _ string, _ *connectionauthentication_settings.Settings) error {
	return errors.New("update not supported: This resource is immutable after creation. Changes require replacement")
}

func (me *service) Delete(_ context.Context, _ string) error {
	return nil
}
