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

package remoteaccess

import (
	"context"
	"fmt"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api"
	remoteaccess "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/cluster/v2/remoteaccess/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest"
)

// ServiceClient TODO: documentation
type ServiceClient struct {
	client rest.Client
}

// NewService creates a new Service Client
// baseURL should look like this: "https://siz65484.live.dynatrace.com/api/config/v1"
// token is an API Token
func NewService(credentials *rest.Credentials) *ServiceClient {
	return &ServiceClient{client: rest.ClusterV2Client(credentials)}
}

// Create TODO: documentation
func (cs *ServiceClient) Create(ctx context.Context, config *remoteaccess.Settings) (*api.Stub, error) {
	config.State = nil

	response := remoteaccess.Settings{}
	if err := cs.client.Post(ctx, "/remoteaccess/requests", config, 201).Finish(&response); err != nil {
		return nil, err
	}

	return &api.Stub{ID: *response.ID, Name: response.UserId, Value: response}, nil
}

// Update TODO: documentation
func (cs *ServiceClient) Update(ctx context.Context, id string, config *remoteaccess.UpdateSettings) error {
	return cs.client.Put(ctx, fmt.Sprintf("/remoteaccess/requests/%s/state", id), config, 200).Finish()
}

// Delete TODO: documentation
func (cs *ServiceClient) Delete() error {
	return nil
}

// Get TODO: documentation
func (cs *ServiceClient) Get(ctx context.Context, id string) (*remoteaccess.Settings, error) {
	var err error
	var config remoteaccess.Settings
	if err = cs.client.Get(ctx, fmt.Sprintf("/remoteaccess/requests/%s", id), 200).Finish(&config); err != nil {
		return nil, err
	}
	return &config, nil
}
