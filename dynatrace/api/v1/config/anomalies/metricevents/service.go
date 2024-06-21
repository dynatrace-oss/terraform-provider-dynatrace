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

package metricevents

import (
	"context"
	"fmt"
	"strings"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings"

	metricevents "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/anomalies/metricevents/settings"
)

const SchemaID = "v1:config:anomaly-detection:metric-events"
const BasePath = "/api/config/v1/anomalyDetection/metricEvents"

func Service(credentials *settings.Credentials) settings.CRUDService[*metricevents.MetricEvent] {
	return settings.NewCRUDService(
		credentials,
		SchemaID,
		settings.DefaultServiceOptions[*metricevents.MetricEvent](BasePath).WithCreateRetry(RetryOnCreate).WithDuplicates(Duplicates),
	)
}

func Duplicates(ctx context.Context, service settings.RService[*metricevents.MetricEvent], v *metricevents.MetricEvent) (*api.Stub, error) {
	if settings.RejectDuplicate("dynatrace_custom_anomalies") {
		var err error
		var stubs api.Stubs
		if stubs, err = service.List(ctx); err != nil {
			return nil, err
		}
		for _, stub := range stubs {
			if v.Name == stub.Name {
				return nil, fmt.Errorf("Custom Anomalies named '%s' already exists", v.Name)
			}
		}
	} else if settings.HijackDuplicate("dynatrace_custom_anomalies") {
		var err error
		var stubs api.Stubs
		if stubs, err = service.List(ctx); err != nil {
			return nil, err
		}
		for _, stub := range stubs {
			if v.Name == stub.Name {
				return stub, nil
			}
		}
	}
	return nil, nil
}

func RetryOnCreate(v *metricevents.MetricEvent, err error) *metricevents.MetricEvent {
	var ok bool
	var restErr *rest.Error
	if restErr, ok = err.(*rest.Error); !ok {
		return nil
	}
	if len(restErr.ConstraintViolations) == 0 {
		return nil
	}
	for _, violation := range restErr.ConstraintViolations {
		violationMessage := violation.Message
		if strings.Contains(violationMessage, "Metric selectors could only be parsed in backward compatibility mode") {
			if strings.Contains(violationMessage, "Consider querying `") && strings.Contains(violationMessage, "` instead..") {
				metricSelector := violationMessage[strings.Index(violationMessage, "Consider querying `")+len("Consider querying `") : strings.Index(violationMessage, "` instead..")]
				v.MetricSelector = &metricSelector
				return v
			}
		}
	}
	return nil
}
