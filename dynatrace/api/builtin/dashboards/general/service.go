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

package general

import (
	"context"
	"errors"
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api"
	general "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/dashboards/general/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings/services/settings20"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

const SchemaVersion = "1.0.18"
const SchemaID = "builtin:dashboards.general"
const defaultTimeout = 20 * time.Second

type service struct {
	service settings.CRUDService[*general.Settings]
}

func Service(credentials *rest.Credentials) settings.CRUDService[*general.Settings] {
	return &service{settings20.Service[*general.Settings](credentials, SchemaID, SchemaVersion)}
}
func (me *service) Create(ctx context.Context, v *general.Settings) (*api.Stub, error) {
	var apiStub *api.Stub
	err := retry.RetryContext(ctx, defaultTimeout, func() *retry.RetryError {
		var err error
		apiStub, err = me.service.Create(ctx, v)
		return classifyRetryError(err)
	})

	if err != nil {
		return nil, err
	}
	return apiStub, nil
}

// classifyRetryErrors retries on certain conflicts
// During "create" the data source of the dashboard may not be up to date.
func classifyRetryError(err error) *retry.RetryError {
	if err == nil {
		return nil
	}
	if restError, ok := errors.AsType[rest.Error](err); ok && restError.Code == http.StatusBadRequest {
		containsConflict := slices.ContainsFunc(restError.ConstraintViolations, func(cv rest.ConstraintViolation) bool {
			return strings.Contains(cv.Message, "Invalid value in datasource")
		})
		if containsConflict {
			return retry.RetryableError(err)
		}
	}

	return retry.NonRetryableError(err)
}

func (me *service) Update(ctx context.Context, id string, v *general.Settings) error {
	return retry.RetryContext(ctx, defaultTimeout, func() *retry.RetryError {
		err := me.service.Update(ctx, id, v)
		return classifyRetryError(err)
	})
}

func (me *service) Delete(ctx context.Context, id string) error {
	return me.service.Delete(ctx, id)
}

func (me *service) Get(ctx context.Context, id string, v *general.Settings) error {
	return me.service.Get(ctx, id, v)
}

func (me *service) List(ctx context.Context) (api.Stubs, error) {
	return me.service.List(ctx)
}

func (me *service) SchemaID() string {
	return SchemaID
}
