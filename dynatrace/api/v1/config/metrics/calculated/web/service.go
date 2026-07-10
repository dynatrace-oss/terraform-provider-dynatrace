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
	mysettings "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/metrics/calculated/web/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

const SchemaID = "v1:config:calculated-metrics-web"
const defaultTimeout = 30 * time.Second

// error messages (parts) that may be returned because of eventual consistency
const (
	overallErr = "Found more than one metric" // create => delete => create => leads to having the same "unique" identifier (metric key) twice
	deleteErr  = "Unable to delete"           // status 500
	getErr     = "does not exist"             // should exist but is not synced across server nodes.
)

var createErrs = []string{"Unable to create", "already exists"}

func Service(clientSet rest.ClientSet) settings.CRUDService[*mysettings.CalculatedWebMetric] {
	return &service{client: rest.APITokenClient(clientSet.Credentials())}
}

type service struct {
	client rest.Client
}

func (me *service) Get(ctx context.Context, id string, v *mysettings.CalculatedWebMetric) error {
	return retry.RetryContext(ctx, defaultTimeout, func() *retry.RetryError {
		err := me.get(ctx, id, v)
		return classifyRetryError(err, getErr)
	})
}

func (me *service) get(ctx context.Context, id string, v *mysettings.CalculatedWebMetric) error {
	return me.client.Get(ctx, fmt.Sprintf("/api/config/v1/calculatedMetrics/rum/%s", url.PathEscape(id)), 200).Finish(v)
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

	err := retry.RetryContext(ctx, defaultTimeout, func() *retry.RetryError {
		err := me.client.Post(ctx, "/api/config/v1/calculatedMetrics/rum", v, 201).Finish(&stub)
		return classifyRetryError(err, createErrs...)
	})

	return &stub, err
}

// classifyRetryError retries when there is a conflict during get/create/update/delete.
// Create after delete (e.g., ForceNew) => "Unable to create RumMetric my-metric in DemMetricsConfigPersistence or Metadata: Unable to create RumMetric my-metric in DemMetricsConfigPersistence"
// Delete after cleanup => "Found more than one metric with key <my-key>"
// GET after apply of an update => "Metric with key <my-key> does not exist"
func classifyRetryError(err error, expectedErrs ...string) *retry.RetryError {
	if err == nil {
		return nil
	}
	restError, ok := errors.AsType[rest.Error](err)

	if !ok || (restError.Code != http.StatusBadRequest && restError.Code != http.StatusInternalServerError) {
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
	return retry.RetryContext(ctx, defaultTimeout, func() *retry.RetryError {
		err := me.client.Put(ctx, fmt.Sprintf("/api/config/v1/calculatedMetrics/rum/%s", url.PathEscape(id)), v, 204).Finish()
		return classifyRetryError(err, getErr)
	})
}

func (me *service) delete(ctx context.Context, id string) error {
	return retry.RetryContext(ctx, defaultTimeout, func() *retry.RetryError {
		err := me.client.Delete(ctx, fmt.Sprintf("/api/config/v1/calculatedMetrics/rum/%s", url.PathEscape(id)), 204, 200).Finish()
		return classifyRetryError(err, deleteErr)
	})
}

func (me *service) Delete(ctx context.Context, id string) error {
	// There are several problems with this resource that are caused by server-node syncs:
	// 1. Consistency errors (server nodes issues) during delete => "500: Unable to delete" or "Found more than one metric with key <my-key>"
	//    Retrying on delete resolves that
	// 2. Creation errors after deletion (e.g., ForceNew) => "Found more than one metric with key <my-key>"
	//    Retrying inside create may not resolve it, even after minutes.
	//    Current solution: retry inside delete until GET returns 404 {desiredSuccesses} times
	//    Problem: GET after delete may return 404 but create fails again with "Found more than one metric with key <my-key>", even with 30s retry
	//    Assumption: Delete may be ignored
	err := me.delete(ctx, id)
	if err != nil {
		return err
	}

	desiredSuccesses := 5
	successes := 0
	consistencyRetryErr := errors.New("eventual consistency check")
	err = retry.RetryContext(ctx, 1*time.Minute /*just to be on the safe side*/, func() *retry.RetryError {
		err = me.get(ctx, id, new(mysettings.CalculatedWebMetric))
		if err == nil {
			// Metric still exists - delete and try again
			err = me.delete(ctx, id)
			// ignore 404 during delete
			if err != nil && !rest.IsNotFoundError(err) {
				return retry.NonRetryableError(err)
			}
			return retry.RetryableError(consistencyRetryErr)
		}
		if rest.IsNotFoundError(err) {
			successes++
			if successes >= desiredSuccesses {
				// Consider it a success if the metric is not found for several consecutive retries,
				// which indicates that the deletion has fully propagated across server nodes.
				return nil
			}
			// Not found, but might be eventual consistency - retry until desiredSuccesses is reached
			return retry.RetryableError(consistencyRetryErr)
		}
		// Unexpected error - stop retrying
		return retry.NonRetryableError(err)
	})

	if err != nil && errors.Is(err, consistencyRetryErr) {
		// ignore our retry check errors (returned by timeout)
		return nil
	}
	return err
}

func (me *service) SchemaID() string {
	return SchemaID
}
