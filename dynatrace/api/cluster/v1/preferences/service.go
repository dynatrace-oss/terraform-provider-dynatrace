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

package preferences

import (
	"context"

	preferences "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/cluster/v1/preferences/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest"
)

// ServiceClient TODO: documentation
type ServiceClient struct {
	client rest.Client
}

// NewService creates a new Service Client
// baseURL should look like this: "https://siz65484.live.dynatrace.com/api/config/v1"
// token is an API Token
func NewService(baseURL string, token string) *ServiceClient {
	return &ServiceClient{client: rest.DefaultClient(baseURL, token)}
}

// Create TODO: documentation
func (cs *ServiceClient) Create(ctx context.Context, config *preferences.Settings) error {
	return cs.Update(ctx, config)
}

// Update TODO: documentation
func (cs *ServiceClient) Update(ctx context.Context, config *preferences.Settings) error {
	return cs.client.Post(ctx, "/preferences", config, 200).Finish()
}

// Delete TODO: documentation
func (cs *ServiceClient) Delete(config *preferences.Settings) error {
	return nil
}

// Get TODO: documentation
func (cs *ServiceClient) Get(ctx context.Context) (*preferences.Settings, error) {
	var err error
	var config preferences.Settings
	if err = cs.client.Get(ctx, "/preferences", 200).Finish(&config); err != nil {
		return nil, err
	}
	return &config, nil
}
