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
	"encoding/json"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/xjson"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// CalculatedMetricDefinition The definition of a calculated service metric.
type CalculatedMetricDefinition struct {
	Metric           *Metric `json:"metric"`                     // The metric to be captured.
	RequestAttribute *string `json:"requestAttribute,omitempty"` // The request attribute to be captured. Only applicable when the **metric** parameter is set to `REQUEST_ATTRIBUTE`.
}

func (me *CalculatedMetricDefinition) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"metric": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The metric to be captured. Possible values are `CPU_TIME`, `DATABASE_CHILD_CALL_COUNT`, `DATABASE_CHILD_CALL_TIME`, `DISK_IO_TIME`, `EXCEPTION_COUNT`, `FAILED_REQUEST_COUNT`, `FAILED_REQUEST_COUNT_CLIENT`, `FAILURE_RATE`, `FAILURE_RATE_CLIENT`, `HTTP_4XX_ERROR_COUNT`, `HTTP_4XX_ERROR_COUNT_CLIENT`, `HTTP_5XX_ERROR_COUNT`, `HTTP_5XX_ERROR_COUNT_CLIENT`, `IO_TIME`, `LOCK_TIME`, `NETWORK_IO_TIME`, `NON_DATABASE_CHILD_CALL_COUNT`, `NON_DATABASE_CHILD_CALL_TIME`, `PROCESSING_TIME`, `REQUEST_ATTRIBUTE`, `REQUEST_COUNT`, `RESPONSE_TIME`, `RESPONSE_TIME_CLIENT`, `SUCCESSFUL_REQUEST_COUNT`, `SUCCESSFUL_REQUEST_COUNT_CLIENT` and `WAIT_TIME`",
		},
		"request_attribute": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The request attribute to be captured. Only applicable when the **metric** parameter is set to `REQUEST_ATTRIBUTE`",
		},
	}
}

func (me *CalculatedMetricDefinition) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"metric":            me.Metric,
		"request_attribute": me.RequestAttribute,
	})
}

func (me *CalculatedMetricDefinition) UnmarshalHCL(decoder hcl.Decoder) error {
	err := decoder.DecodeAll(map[string]any{
		"metric":            &me.Metric,
		"request_attribute": &me.RequestAttribute,
	})
	return err
}

func (me *CalculatedMetricDefinition) MarshalJSON() ([]byte, error) {
	properties := xjson.Properties{}
	if err := properties.MarshalAll(map[string]any{
		"metric":           me.Metric,
		"requestAttribute": me.RequestAttribute,
	}); err != nil {
		return nil, err
	}
	return json.Marshal(properties)
}

func (me *CalculatedMetricDefinition) UnmarshalJSON(data []byte) error {
	properties := xjson.Properties{}
	if err := json.Unmarshal(data, &properties); err != nil {
		return err
	}
	return properties.UnmarshalAll(map[string]any{
		"metric":           &me.Metric,
		"requestAttribute": &me.RequestAttribute,
	})
}

// Metric The metric to be captured.
type Metric string

// Metrics offers the known enum values
var Metrics = struct {
	CPUTime                      Metric
	DatabaseChildCallCount       Metric
	DatabaseChildCallTime        Metric
	DiskIoTime                   Metric
	ExceptionCount               Metric
	FailedRequestCount           Metric
	FailedRequestCountClient     Metric
	FailureRate                  Metric
	FailureRateClient            Metric
	Http4xxErrorCount            Metric
	Http4xxErrorCountClient      Metric
	Http5xxErrorCount            Metric
	Http5xxErrorCountClient      Metric
	IoTime                       Metric
	LockTime                     Metric
	NetworkIoTime                Metric
	NonDatabaseChildCallCount    Metric
	NonDatabaseChildCallTime     Metric
	ProcessingTime               Metric
	RequestAttribute             Metric
	RequestCount                 Metric
	ResponseTime                 Metric
	ResponseTimeClient           Metric
	SuccessfulRequestCount       Metric
	SuccessfulRequestCountClient Metric
	WaitTime                     Metric
}{
	"CPU_TIME",
	"DATABASE_CHILD_CALL_COUNT",
	"DATABASE_CHILD_CALL_TIME",
	"DISK_IO_TIME",
	"EXCEPTION_COUNT",
	"FAILED_REQUEST_COUNT",
	"FAILED_REQUEST_COUNT_CLIENT",
	"FAILURE_RATE",
	"FAILURE_RATE_CLIENT",
	"HTTP_4XX_ERROR_COUNT",
	"HTTP_4XX_ERROR_COUNT_CLIENT",
	"HTTP_5XX_ERROR_COUNT",
	"HTTP_5XX_ERROR_COUNT_CLIENT",
	"IO_TIME",
	"LOCK_TIME",
	"NETWORK_IO_TIME",
	"NON_DATABASE_CHILD_CALL_COUNT",
	"NON_DATABASE_CHILD_CALL_TIME",
	"PROCESSING_TIME",
	"REQUEST_ATTRIBUTE",
	"REQUEST_COUNT",
	"RESPONSE_TIME",
	"RESPONSE_TIME_CLIENT",
	"SUCCESSFUL_REQUEST_COUNT",
	"SUCCESSFUL_REQUEST_COUNT_CLIENT",
	"WAIT_TIME",
}
