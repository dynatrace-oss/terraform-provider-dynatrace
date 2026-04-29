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

package web

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"slices"
	"strings"
	"time"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api"
	retrycommon "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/hyperscalerauthentication/connections/retry"
	mysettings "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/metrics/calculated/web/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

const SchemaID = "v1:config:calculated-metrics-web"
const defaultTimeout = 20 * time.Second

// error messages (parts) that may be returned because of eventual consistency
const (
	overallErr = "Found more than one metric" // create => delete => create => leads to having the same "unique" identifier (metric key) twice
	createErr  = "Unable to create"
	deleteErr  = "Unable to delete"
	getErr     = "does not exist" // should exist but is not synced across server nodes.
)

func Service(credentials *rest.Credentials) settings.CRUDService[*mysettings.CalculatedWebMetric] {
	return &service{client: rest.APITokenClient(credentials)}
}

type service struct {
	client rest.Client
}

func (me *service) Get(ctx context.Context, id string, v *mysettings.CalculatedWebMetric) error {
	return retry.RetryContext(ctx, retrycommon.DurationUntilDeadlineOrDefault(ctx, defaultTimeout), func() *retry.RetryError {
		err := me.client.Get(ctx, fmt.Sprintf("/api/config/v1/calculatedMetrics/rum/%s", url.PathEscape(id)), 200).Finish(v)
		return classifyRetryError(err, getErr)
	})
}

func (me *service) List(ctx context.Context) (api.Stubs, error) {
	var err error

	req := me.client.Get(ctx, "/api/config/v1/calculatedMetrics/rum", 200)
	var stubList api.StubList
	if err = req.Finish(&stubList); err != nil {
		return nil, err
	}
	return stubList.Values, nil
}

func (me *service) Validate(ctx context.Context, v *mysettings.CalculatedWebMetric) error {
	var err error
	client := me.client

	req := client.Post(ctx, "/api/config/v1/calculatedMetrics/rum/validator", v, 204)
	if err = req.Finish(); err != nil {
		return err
	}

	return nil
}

func (me *service) Create(ctx context.Context, v *mysettings.CalculatedWebMetric) (*api.Stub, error) {
	var stub api.Stub

	err := retry.RetryContext(ctx, retrycommon.DurationUntilDeadlineOrDefault(ctx, defaultTimeout), func() *retry.RetryError {
		err := me.client.Post(ctx, "/api/config/v1/calculatedMetrics/rum", v, 201).Finish(&stub)
		return classifyRetryError(err, createErr)
	})

	return &stub, err
}

// classifyRetryError retries when there is a conflict during create/update/delete.
// Create after delete (e.g., ForceNew) => "Unable to create RumMetric my-metric in DemMetricsConfigPersistence or Metadata: Unable to create RumMetric my-metric in DemMetricsConfigPersistence"
// Delete after cleanup => "Found more than one metric with key <my-key>"
// GET after apply of an update => "Metric with key <my-key> does not exist"
func classifyRetryError(err error, expectedErrs ...string) *retry.RetryError {
	if err == nil {
		return nil
	}
	restError, ok := errors.AsType[rest.Error](err)

	if !ok || restError.Code != http.StatusBadRequest {
		return retry.NonRetryableError(err)
	}

	expectedErrs = append(expectedErrs, overallErr)

	containsEventualConsistencyErr := slices.ContainsFunc(expectedErrs, func(s string) bool {
		return strings.Contains(restError.Message, s)
	})
	if containsEventualConsistencyErr {
		return retry.RetryableError(err)
	}

	return retry.NonRetryableError(err)
}

func (me *service) Update(ctx context.Context, id string, v *mysettings.CalculatedWebMetric) error {
	return retry.RetryContext(ctx, retrycommon.DurationUntilDeadlineOrDefault(ctx, defaultTimeout), func() *retry.RetryError {
		err := me.client.Put(ctx, fmt.Sprintf("/api/config/v1/calculatedMetrics/rum/%s", url.PathEscape(id)), v, 204).Finish()
		return classifyRetryError(err, getErr)
	})
}

func (me *service) Delete(ctx context.Context, id string) error {
	return retry.RetryContext(ctx, retrycommon.DurationUntilDeadlineOrDefault(ctx, defaultTimeout), func() *retry.RetryError {
		err := me.client.Delete(ctx, fmt.Sprintf("/api/config/v1/calculatedMetrics/rum/%s", url.PathEscape(id)), 204, 200).Finish()
		return classifyRetryError(err, deleteErr)
	})
}

func (me *service) SchemaID() string {
	return SchemaID
}
