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
	"fmt"
	"net/url"
	"strings"
	"time"

	mysettings "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/metrics/calculated/service/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings"
)

const SchemaID = "v1:config:calculated-metrics:service"

func Service(credentials *settings.Credentials) settings.CRUDService[*mysettings.CalculatedServiceMetric] {
	return &service{client: rest.DefaultClient(credentials.URL, credentials.Token)}
}

type service struct {
	client rest.Client
}

func (me *service) Get(id string, v *mysettings.CalculatedServiceMetric) error {
	return me.client.Get(fmt.Sprintf("/api/config/v1/calculatedMetrics/service/%s", url.PathEscape(id)), 200).Finish(v)
}

func (me *service) List() (settings.Stubs, error) {
	var err error

	req := me.client.Get("/api/config/v1/calculatedMetrics/service", 200)
	var stubList settings.StubList
	if err = req.Finish(&stubList); err != nil {
		return nil, err
	}
	return stubList.Values, nil
}

func (me *service) Validate(v *mysettings.CalculatedServiceMetric) error {
	var err error

	client := me.client

	retry := true
	maxAttempts := 64
	attempts := 0

	for retry {
		attempts = attempts + 1
		req := client.Post("/api/config/v1/calculatedMetrics/service/validator", v, 204)
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

func (me *service) Create(v *mysettings.CalculatedServiceMetric) (*settings.Stub, error) {
	var err error

	client := me.client

	retry := true
	maxAttempts := 64
	attempts := 0
	var stub settings.Stub

	for retry {
		attempts = attempts + 1
		req := client.Post("/api/config/v1/calculatedMetrics/service", v, 201)
		if err = req.Finish(&stub); err != nil {
			if strings.Contains(err.Error(), "Metric definition must specify a known request attribute") {
				if attempts < maxAttempts {
					time.Sleep(500 * time.Millisecond)
				} else {
					return nil, err
				}
			}
			if strings.Contains(err.Error(), "At least one condition of the following types must be used:") {
				restErr := err.(rest.Error)
				// reasons := []string{}
				if len(restErr.ConstraintViolations) > 0 {
					// for _, violations := range restErr.ConstraintViolations {
					// 	reasons = append(reasons, violations.Message)
					// }
					return &settings.Stub{ID: v.TsmMetricKey + "---flawed----", Name: v.Name}, nil
				}
			}
			if strings.Contains(err.Error(), `{"parameterLocation":"PAYLOAD_BODY","message":"must not be null","path":"metricDefinition.metric"}`) {
				return &settings.Stub{ID: v.TsmMetricKey + "---flawed----", Name: v.Name}, nil
			}
			if strings.Contains(err.Error(), `{"parameterLocation":"PAYLOAD_BODY","message":"Please check entityId: there is no such SERVICE in the system","path":"entityId"}`) {
				return &settings.Stub{ID: v.TsmMetricKey + "---flawed----", Name: v.Name}, nil
			}
			return nil, err
		} else {
			return &stub, nil
		}
	}
	return &stub, nil
}

func (me *service) Update(id string, v *mysettings.CalculatedServiceMetric) error {
	if err := me.client.Put(fmt.Sprintf("/api/config/v1/calculatedMetrics/service/%s", url.PathEscape(id)), v, 204).Finish(); err != nil {
		return err
	}
	return nil
}

func (me *service) Delete(id string) error {
	for {
		if err := me.client.Delete(fmt.Sprintf("/api/config/v1/calculatedMetrics/service/%s", url.PathEscape(id)), 204).Finish(); err != nil {
			if strings.Contains(err.Error(), fmt.Sprintf("Service metric with %s not found", id)) {
				return nil
			}
			if !strings.Contains(err.Error(), "Could not delete configuration") {
				return err
			}
		} else {
			break
		}
	}
	return nil
}

func (me *service) SchemaID() string {
	return SchemaID
}
