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

package parameters

import (
	"context"
	"time"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api"
	parameters "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/failuredetection/environment/parameters/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings/services/settings20"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

const SchemaVersion = "1.0.8"
const SchemaID = "builtin:failure-detection.environment.parameters"
const defaultTimeout = 10 * time.Second

// There is a potential "delete"-error happening due to eventual consistency.
// error: builtin:failure-detection.environment.parameters: Failure detection rule EnvironmentFailureDetectionRules {your-values-here} refers to this parameter ID: <uuid>.
// The related object is deleted but not propagated to the server node.
const retryErr = "refers to this parameter"

type service struct {
	service settings.CRUDService[*parameters.Settings]
}

func Service(clientSet rest.ClientSet) settings.CRUDService[*parameters.Settings] {
	return &service{settings20.Service[*parameters.Settings](clientSet, SchemaID, SchemaVersion)}
}
func (me *service) Create(ctx context.Context, v *parameters.Settings) (*api.Stub, error) {
	return me.service.Create(ctx, v)
}

func (me *service) Update(ctx context.Context, id string, v *parameters.Settings) error {
	return me.service.Update(ctx, id, v)
}

func (me *service) Delete(ctx context.Context, id string) error {
	return retry.RetryContext(ctx, defaultTimeout, func() *retry.RetryError {
		err := me.service.Delete(ctx, id)
		return settings.ClassifyConstraintRetryError(err, retryErr)
	})
}

func (me *service) Get(ctx context.Context, id string, v *parameters.Settings) error {
	return me.service.Get(ctx, id, v)
}

func (me *service) List(ctx context.Context) (api.Stubs, error) {
	return me.service.List(ctx)
}

func (me *service) SchemaID() string {
	return SchemaID
}
