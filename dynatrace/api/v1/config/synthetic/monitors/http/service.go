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

package http

import (
	"context"
	"errors"
	"time"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/synthetic/monitors"
	http "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/synthetic/monitors/http/settings"
)

const SchemaID = "v1:synthetic:monitors:http"
const BasePath = "/api/v1/synthetic/monitors"

const DesiredReadySuccesses = 5

var ErrConsistencyRetry = errors.New("eventual consistency check")

func Service(credentials *rest.Credentials) settings.CRUDService[*http.SyntheticMonitor] {
	return &service{service: settings.NewAPITokenService(credentials, SchemaID, &settings.ServiceOptions[*http.SyntheticMonitor]{
		Get:            settings.Path("/api/v1/synthetic/monitors/%s"),
		List:           settings.Path("/api/v1/synthetic/monitors?type=HTTP"),
		CreateURL:      func(v *http.SyntheticMonitor) string { return "/api/v1/synthetic/monitors" },
		Stubs:          &monitors.Monitors{},
		HasNoValidator: true,
	})}
}

func WithPredefinedService(s settings.CRUDService[*http.SyntheticMonitor]) *service {
	return &service{service: s}
}

type service struct {
	service settings.CRUDService[*http.SyntheticMonitor]
}

func GetTempScript() *http.Script {
	return &http.Script{Requests: http.Requests{&http.Request{
		Description: new("--terraform-auto-generated--"),
		URL:         "http://localhost",
		Method:      "OPTIONS",
	}}}
}

func (me *service) Create(ctx context.Context, v *http.SyntheticMonitor) (*api.Stub, error) {
	if v.NoScript != nil && *v.NoScript && v.Script == nil {
		v.Script = GetTempScript()
	}
	stub, err := me.service.Create(ctx, v)
	if err != nil {
		return nil, err
	}
	err = me.validateReady(ctx, stub.ID, v)
	if err != nil {
		return nil, err
	}
	return stub, nil
}

func (me *service) Update(ctx context.Context, id string, v *http.SyntheticMonitor) error {
	if v.NoScript != nil && *v.NoScript && v.Script == nil {
		monitorSettings := new(http.SyntheticMonitor)
		if err := me.service.Get(ctx, id, monitorSettings); err != nil {
			return err
		}
		v.Script = monitorSettings.Script
	}
	err := me.service.Update(ctx, id, v)
	if err != nil {
		return err
	}
	return me.validateReady(ctx, id, v)
}

// validateReady check when an HTTP monitor is ready to use/get after create/update.
// Because of eventual consistency
// - `tags` may not be returned.
// - Get may return 404
// - Even if Get returns non 404 and the expected tags, we can't rely on it because server nodes may not be synced.
// That's why we need to
// 1. Retry on 404
// 2. Retry if tags are not present (if there should be any)
// 3. Retry several times, even if 1. and 2. succeeded. Server nodes may still be out of sync.
func (me *service) validateReady(ctx context.Context, id string, v *http.SyntheticMonitor) error {
	retryOnTags := len(v.Tags) > 0

	err := retry.RetryContext(ctx, 1*time.Minute, me.ReadyCheck(ctx, id, retryOnTags))
	if err != nil && errors.Is(err, ErrConsistencyRetry) {
		// ignore our retry check errors (returned by timeout)
		return nil
	}
	if err != nil {
		return err
	}
	return nil
}

// ReadyCheck returns a retry.RetryFunc that performs one readiness probe per
// invocation. The returned closure holds the success counter, so the same
// instance must be reused across retries (and can be invoked directly from
// tests to drive each branch synchronously).
func (me *service) ReadyCheck(ctx context.Context, id string, retryOnTags bool) retry.RetryFunc {
	successes := 0
	return func() *retry.RetryError {
		newVal := new(http.SyntheticMonitor)
		err := me.Get(ctx, id, newVal)

		// may not exist immediately after creation (not synced across server nodes)
		if rest.IsNotFoundError(err) {
			return retry.RetryableError(err)
		}

		// other errors not related to inconsistency and 404
		if err != nil {
			return retry.NonRetryableError(err)
		}

		// checks if tags are synced (consistent)
		if retryOnTags && len(newVal.Tags) == 0 {
			return retry.RetryableError(ErrConsistencyRetry)
		}

		successes++
		if successes < DesiredReadySuccesses {
			return retry.RetryableError(ErrConsistencyRetry)
		}
		return nil
	}
}

func (me *service) Delete(ctx context.Context, id string) error {
	return me.service.Delete(ctx, id)
}

func (me *service) Get(ctx context.Context, id string, v *http.SyntheticMonitor) error {
	if err := me.service.Get(ctx, id, v); err != nil {
		return err
	}
	if v.Script != nil && len(v.Script.Requests) == 1 && *v.Script.Requests[0].Description == *GetTempScript().Requests[0].Description {
		v.NoScript = new(true)
		v.Script = nil
	}
	return nil
}

func (me *service) List(ctx context.Context) (api.Stubs, error) {
	return me.service.List(ctx)
}

func (me *service) SchemaID() string {
	return me.service.SchemaID()
}
