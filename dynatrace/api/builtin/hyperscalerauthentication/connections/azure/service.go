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

package azure

import (
	"context"
	"time"

	serviceSettings "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/hyperscalerauthentication/connections/azure/settings"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api"
	retrycommon "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/hyperscalerauthentication/connections/retry"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings/services/settings20"
)

const SchemaVersion = "0.0.10"
const SchemaID = "builtin:hyperscaler-authentication.connections.azure"
const DefaultTimeout = 2 * time.Minute

func Service(credentials *rest.Credentials) settings.CRUDService[*serviceSettings.Settings] {
	return &service{
		service: settings20.Service[*serviceSettings.Settings](credentials, SchemaID, SchemaVersion),
	}
}

type service struct {
	service settings.CRUDService[*serviceSettings.Settings]
}

func (me *service) List(ctx context.Context) (api.Stubs, error) {
	return me.service.List(ctx)
}

func (me *service) Get(ctx context.Context, id string, v *serviceSettings.Settings) error {
	return me.service.Get(ctx, id, v)
}

func (me *service) SchemaID() string {
	return me.service.SchemaID()
}

func (me *service) Create(ctx context.Context, v *serviceSettings.Settings) (*api.Stub, error) {
	var stub *api.Stub
	err := retry.RetryContext(ctx, retrycommon.DurationUntilDeadlineOrDefault(ctx, serviceSettings.DefaultTimeout), func() *retry.RetryError {
		var err error
		stub, err = me.service.Create(ctx, v)
		return retrycommon.ClassifyRetryError(err)
	})
	if err != nil {
		return nil, err
	}
	return stub, nil
}

func (me *service) Update(ctx context.Context, id string, v *serviceSettings.Settings) error {
	return retry.RetryContext(ctx, retrycommon.DurationUntilDeadlineOrDefault(ctx, serviceSettings.DefaultTimeout), func() *retry.RetryError {
		return retrycommon.ClassifyRetryError(me.service.Update(ctx, id, v))
	})
}

func (me *service) Delete(ctx context.Context, id string) error {
	return me.service.Delete(ctx, id)
}
