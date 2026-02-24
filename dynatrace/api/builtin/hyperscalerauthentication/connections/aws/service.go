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

package aws

import (
	"context"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api"
	awsconnection "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/hyperscalerauthentication/connections/aws/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings/services/settings20"
)

const SchemaVersion = "0.0.24"
const SchemaID = "builtin:hyperscaler-authentication.connections.aws"

type service struct {
	service settings.CRUDService[*awsconnection.Settings]
}

// An update should only happen when the role ARN is added or the name is modified. Otherwise, all attributes and subresources
// are flagged as "forceNew", meaning instead of an update, the resource is destroyed and created from scratch
func Service(credentials *rest.Credentials) settings.CRUDService[*awsconnection.Settings] {
	return &service{service: settings20.Service[*awsconnection.Settings](credentials, SchemaID, SchemaVersion)}
}

func (me *service) SchemaID() string {
	return me.service.SchemaID()
}

func (me *service) Get(ctx context.Context, id string, v *awsconnection.Settings) error {
	return me.service.Get(ctx, id, v)
}

func (me *service) List(ctx context.Context) (api.Stubs, error) {
	return me.service.List(ctx)
}

func (me *service) Create(ctx context.Context, v *awsconnection.Settings) (*api.Stub, error) {
	return me.service.Create(ctx, v)
}

// Update is only used to update the name of the connection
// role_arn is not modifiable set in a different resource, and it's not in the state of this resource.
// Therefore, we need to fetch the current state and assign the values manually instead of overriding the role_arn with empty/null
func (me *service) Update(ctx context.Context, id string, v *awsconnection.Settings) error {
	var current awsconnection.Settings
	if err := me.Get(ctx, id, &current); err != nil {
		return err
	}
	current.Name = v.Name
	return me.service.Update(ctx, id, &current)
}

func (me *service) Delete(ctx context.Context, id string) error {
	return me.service.Delete(ctx, id)
}
