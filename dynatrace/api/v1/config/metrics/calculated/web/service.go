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

func Service(credentials *rest.Credentials) settings.CRUDService[*mysettings.CalculatedWebMetric] {
	return &service{client: rest.APITokenClient(credentials)}
}

type service struct {
	client rest.Client
}

func (me *service) Get(ctx context.Context, id string, v *mysettings.CalculatedWebMetric) error {
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
	client := me.client
	var stub api.Stub

	// retry max 10s
	err := retry.RetryContext(ctx, retrycommon.DurationUntilDeadlineOrDefault(ctx, 10*time.Second), func() *retry.RetryError {
		return classifyCreateRetryError(client.Post(ctx, "/api/config/v1/calculatedMetrics/rum", v, 201).Finish(&stub))
	})

	return &stub, err
}

// classifyCreateRetryError retries when there is a conflict during create.
// Create after delete doesn't immediately work, which is needed for re-create due to ForceNew.
// => "Unable to create RumMetric my-metric in DemMetricsConfigPersistence or Metadata: Unable to create RumMetric my-metric in DemMetricsConfigPersistence"
func classifyCreateRetryError(err error) *retry.RetryError {
	if err == nil {
		return nil
	}

	if restError, ok := errors.AsType[rest.Error](err); ok && restError.Code == http.StatusBadRequest && strings.Contains(restError.Message, "Unable to create RumMetric") {
		return retry.RetryableError(err)
	}

	return retry.NonRetryableError(err)
}

func (me *service) Update(ctx context.Context, id string, v *mysettings.CalculatedWebMetric) error {
	if err := me.client.Put(ctx, fmt.Sprintf("/api/config/v1/calculatedMetrics/rum/%s", url.PathEscape(id)), v, 204).Finish(); err != nil {
		return err
	}
	return nil
}

func (me *service) Delete(ctx context.Context, id string) error {
	var err error
	attempts := 30

	for range attempts {
		if err = me.client.Delete(ctx, fmt.Sprintf("/api/config/v1/calculatedMetrics/rum/%s", url.PathEscape(id)), 204, 200).Finish(); err != nil {
			if !strings.Contains(err.Error(), fmt.Sprintf("Unable to delete Rum metric with key: \"%s\" from DemMetricsConfigPersistence", id)) {
				return err
			}
		} else {
			if err = me.Get(ctx, id, &mysettings.CalculatedWebMetric{}); err != nil {
				if strings.Contains(err.Error(), fmt.Sprintf("Metric with key \"%s\" does not exist", id)) {
					break
				}
			}
		}
		time.Sleep(2 * time.Second)
	}

	if err != nil {
		return err
	}
	return nil
}

func (me *service) SchemaID() string {
	return SchemaID
}
