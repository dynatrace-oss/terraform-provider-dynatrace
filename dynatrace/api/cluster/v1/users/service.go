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

package users

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api"
	users "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/cluster/v1/users/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings"
)

const SchemaID = "accounts:users"

func Service(credentials *rest.Credentials) settings.CRUDService[*users.UserConfig] {
	return &service{
		serviceClient: NewService(credentials),
	}
}

// ServiceClient TODO: documentation
type ServiceClient struct {
	client rest.Client
}

func (me *service) Create(ctx context.Context, v *users.UserConfig) (*api.Stub, error) {
	return me.serviceClient.Create(ctx, v)
}

func (me *service) Update(ctx context.Context, id string, v *users.UserConfig) error {
	return me.serviceClient.Update(ctx, v)
}

func (me *service) Delete(ctx context.Context, id string) error {
	return me.serviceClient.Delete(ctx, id)
}

func (me *service) List(ctx context.Context) (api.Stubs, error) {
	return me.serviceClient.List(ctx)
}

func (me *service) Get(ctx context.Context, id string, v *users.UserConfig) error {
	return me.serviceClient.Get(ctx, id, v)
}

func (me *service) SchemaID() string {
	return SchemaID
}

func (cs *ServiceClient) SchemaID() string {
	return SchemaID
}

// NewService creates a new Service Client
// baseURL should look like this: "https://siz65484.live.dynatrace.com/api/config/v1"
// token is an API Token
func NewService(credentials *rest.Credentials) *ServiceClient {
	return &ServiceClient{client: rest.ClusterV1Client(credentials)}
}

type service struct {
	serviceClient *ServiceClient
}

// Create TODO: documentation
func (cs *ServiceClient) Create(ctx context.Context, userConfig *users.UserConfig) (*api.Stub, error) {
	var err error

	var createdUserConfig users.UserConfig
	if err = cs.client.Post(ctx, "/users", userConfig, 200).Finish(&createdUserConfig); err != nil {
		return nil, err
	}
	return &api.Stub{ID: createdUserConfig.UserName, Name: createdUserConfig.UserName}, nil
}

// Update TODO: documentation
func (cs *ServiceClient) Update(ctx context.Context, userConfig *users.UserConfig) error {
	return cs.client.Put(ctx, "/users", userConfig, 200).Finish()
}

// Delete TODO: documentation
func (cs *ServiceClient) Delete(ctx context.Context, id string) error {
	if len(id) == 0 {
		return errors.New("empty ID provided for the user to delete")
	}
	return cs.client.Delete(ctx, fmt.Sprintf("/users/%s", id), 200).Finish()
}

// Get TODO: documentation
func (cs *ServiceClient) Get(ctx context.Context, id string, v *users.UserConfig) error {
	if len(id) == 0 {
		return errors.New("empty ID provided for the user to fetch")
	}

	var err error
	if err = cs.client.Get(ctx, fmt.Sprintf("/users/%s", id), 200).Finish(&v); err != nil {
		if strings.HasPrefix(err.Error(), "Not Found (GET) ") {
			return rest.Error{Code: 404, Message: fmt.Sprintf("user '%s' doesn't exist", id)}
		}
		return err
	}
	return nil
}

// ListAll TODO: documentation
func (cs *ServiceClient) List(ctx context.Context) (api.Stubs, error) {
	var err error
	var stubs api.Stubs
	var users []*users.UserConfig
	if err = cs.client.Get(ctx, "/users", 200).Finish(&users); err != nil {
		return nil, err
	}
	for _, user := range users {
		stubs = append(stubs, &api.Stub{ID: user.UserName, Name: user.UserName})
	}
	return stubs, nil
}
