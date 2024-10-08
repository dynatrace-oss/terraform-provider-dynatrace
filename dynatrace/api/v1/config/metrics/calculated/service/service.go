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

package service

import (
	"context"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api"
	mysettings "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/metrics/calculated/service/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings/services/httpcache"
)

const SchemaID = "v1:config:calculated-metrics-service"

func Service(credentials *settings.Credentials) settings.CRUDService[*mysettings.CalculatedServiceMetric] {
	return &service{client: httpcache.DefaultClient(credentials.URL, credentials.Token, SchemaID)}
}

type service struct {
	client rest.Client
}

func (me *service) Get(ctx context.Context, id string, v *mysettings.CalculatedServiceMetric) error {
	return me.client.Get(ctx, fmt.Sprintf("/api/config/v1/calculatedMetrics/service/%s", url.PathEscape(id)), 200).Finish(v)
}

func (me *service) List(ctx context.Context) (api.Stubs, error) {
	var err error

	req := me.client.Get(ctx, "/api/config/v1/calculatedMetrics/service", 200)
	var stubList api.StubList
	if err = req.Finish(&stubList); err != nil {
		return nil, err
	}
	return stubList.Values, nil
}

func (me *service) Validate(ctx context.Context, v *mysettings.CalculatedServiceMetric) error {
	var err error

	client := me.client

	retry := true
	maxAttempts := 64
	attempts := 0

	for retry {
		attempts = attempts + 1
		req := client.Post(ctx, "/api/config/v1/calculatedMetrics/service/validator", v, 204)
		if err = req.Finish(); err != nil {
			if !strings.Contains(err.Error(), "Metric definition must specify a known request attribute") {
				return err
			}
			// log.Println(".... request attribute is not fully known yet to cluster - retrying")
			if attempts < maxAttempts {
				time.Sleep(500 * time.Millisecond)
			} else {
				return err
			}
		} else {
			return nil
		}
	}
	return nil
}

func (me *service) Create(ctx context.Context, v *mysettings.CalculatedServiceMetric) (*api.Stub, error) {
	var err error

	client := me.client

	retry := true
	maxAttempts := 64
	attempts := 0
	var stub api.Stub

	for retry {
		attempts = attempts + 1
		req := client.Post(ctx, "/api/config/v1/calculatedMetrics/service", v, 201)
		if err = req.Finish(&stub); err != nil {
			if strings.Contains(err.Error(), "Metric definition must specify a known request attribute") {
				if attempts < maxAttempts {
					time.Sleep(500 * time.Millisecond)
				} else {
					return nil, err
				}
			}
			// if strings.Contains(err.Error(), "At least one condition of the following types must be used:") {
			// 	restErr := err.(rest.Error)
			// 	if len(restErr.ConstraintViolations) > 0 {
			// 		return &api.Stub{ID: v.TsmMetricKey + "---flawed----", Name: v.Name}, nil
			// 	}
			// }
			// if strings.Contains(err.Error(), `{"parameterLocation":"PAYLOAD_BODY","message":"must not be null","path":"metricDefinition.metric"}`) {
			// 	return &api.Stub{ID: v.TsmMetricKey + "---flawed----", Name: v.Name}, nil
			// }
			// if strings.Contains(err.Error(), `{"parameterLocation":"PAYLOAD_BODY","message":"Please check entityId: there is no such SERVICE in the system","path":"entityId"}`) {
			// 	return &api.Stub{ID: v.TsmMetricKey + "---flawed----", Name: v.Name}, nil
			// }
			return nil, err
		} else {
			return &stub, nil
		}
	}
	return &stub, nil
}

func (me *service) Update(ctx context.Context, id string, v *mysettings.CalculatedServiceMetric) error {
	if err := me.client.Put(ctx, fmt.Sprintf("/api/config/v1/calculatedMetrics/service/%s", url.PathEscape(id)), v, 204).Finish(); err != nil {
		return err
	}
	return nil
}

func (me *service) Delete(ctx context.Context, id string) error {
	var err error
	attempts := 30

	for i := 0; i < attempts; i++ {
		if err = me.client.Delete(ctx, fmt.Sprintf("/api/config/v1/calculatedMetrics/service/%s", url.PathEscape(id)), 204).Finish(); err != nil {
			if strings.Contains(err.Error(), fmt.Sprintf("Service metric with %s not found", id)) {
				return nil
			}
			if !strings.Contains(err.Error(), "Could not delete configuration") {
				return err
			}
		} else {
			break
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
