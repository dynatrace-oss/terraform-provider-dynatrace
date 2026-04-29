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

const SchemaVersion = "1.0.6"
const SchemaID = "builtin:failure-detection.environment.parameters"
const defaultTimeout = 10 * time.Second

// During "create" and potentially "update" there is a "delete"-error happening that is not related to the object that should be created.
// Seems like the delete-constraint is done at the wrong place, and that there is some eventual consistency.
// error: builtin:failure-detection.environment.parameters: Failure detection rule EnvironmentFailureDetectionRules {your-values-here} refers to this parameter ID: <uuid>.
const retryErr = "refers to this parameter"

type service struct {
	service settings.CRUDService[*parameters.Settings]
}

func Service(credentials *rest.Credentials) settings.CRUDService[*parameters.Settings] {
	return &service{settings20.Service[*parameters.Settings](credentials, SchemaID, SchemaVersion)}
}
func (me *service) Create(ctx context.Context, v *parameters.Settings) (*api.Stub, error) {
	var apiStub *api.Stub
	err := retry.RetryContext(ctx, defaultTimeout, func() *retry.RetryError {
		var err error
		apiStub, err = me.service.Create(ctx, v)
		return settings.ClassifyConstraintRetryError(err, retryErr)
	})

	if err != nil {
		return nil, err
	}
	return apiStub, nil
}

func (me *service) Update(ctx context.Context, id string, v *parameters.Settings) error {
	return retry.RetryContext(ctx, defaultTimeout, func() *retry.RetryError {
		err := me.service.Update(ctx, id, v)
		return settings.ClassifyConstraintRetryError(err, retryErr)
	})
}

func (me *service) Delete(ctx context.Context, id string) error {
	return me.service.Delete(ctx, id)
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
