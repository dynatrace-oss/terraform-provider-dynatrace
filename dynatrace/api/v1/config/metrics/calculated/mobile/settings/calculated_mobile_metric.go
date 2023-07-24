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

package mobile

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// CalculatedMobileMetric Descriptor of a calculated mobile/custom app metric.
type CalculatedMobileMetric struct {
	Description      *string           `json:"description,omitempty"`      // Descriptor of a calculated mobile/custom app metric.
	AppIdentifier    string            `json:"applicationIdentifier"`      // The Dynatrace entity ID of the application to which the metric belongs.
	Name             string            `json:"name"`                       // The displayed name of the metric.
	MetricKey        string            `json:"metricKey"`                  // The unique key of the calculated mobile/custom app metric.
	Enabled          bool              `json:"enabled"`                    // The metric is enabled (`true`) or disabled (`false`).
	MetricType       MetricType        `json:"metricType"`                 // The type of the metric. Possible values: [ REPORTED_ERROR_COUNT, USER_ACTION_DURATION, WEB_REQUEST_COUNT, WEB_REQUEST_ERROR_COUNT ]
	Dimensions       Dimensions        `json:"dimensions,omitempty"`       // Parameters of a definition of a calculated mobile/custom app metric.
	UserActionFilter *UserActionFilter `json:"userActionFilter,omitempty"` // User actions filter of the calculated mobile/custom app metric. Only user actions matching the provided criteria are used for metric calculation. A user action must match all the criteria.
}

func (me *CalculatedMobileMetric) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"description": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Descriptor of a calculated mobile/custom app metric.",
		},
		"app_identifier": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The Dynatrace entity ID of the application to which the metric belongs.",
		},
		"name": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The displayed name of the metric.",
		},
		"metric_key": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The unique key of the calculated mobile/custom app metric.",
			ForceNew:    true,
		},
		"enabled": {
			Type:        schema.TypeBool,
			Required:    true,
			Description: "The metric is enabled (`true`) or disabled (`false`)",
		},
		"metric_type": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The type of the metric. Possible values: [ REPORTED_ERROR_COUNT, USER_ACTION_DURATION, WEB_REQUEST_COUNT, WEB_REQUEST_ERROR_COUNT ]",
		},
		"dimensions": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "Parameters of a definition of a calculated mobile/custom app metric.",
			Elem:        &schema.Resource{Schema: new(Dimensions).Schema()},
		},
		"user_action_filter": {
			Type:        schema.TypeList,
			Optional:    true,
			MaxItems:    1,
			Description: "Parameters of a definition of a calculated mobile/custom app metric.",
			Elem:        &schema.Resource{Schema: new(UserActionFilter).Schema()},
		},
	}
}

func (me *CalculatedMobileMetric) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"description":        me.Description,
		"app_identifier":     me.AppIdentifier,
		"name":               me.Name,
		"metric_key":         me.MetricKey,
		"enabled":            me.Enabled,
		"metric_type":        me.MetricType,
		"dimensions":         me.Dimensions,
		"user_action_filter": me.UserActionFilter,
	})
}

func (me *CalculatedMobileMetric) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"description":        &me.Description,
		"app_identifier":     &me.AppIdentifier,
		"name":               &me.Name,
		"metric_key":         &me.MetricKey,
		"enabled":            &me.Enabled,
		"metric_type":        &me.MetricType,
		"dimensions":         &me.Dimensions,
		"user_action_filter": &me.UserActionFilter,
	})
}

type MetricType string

var MetricTypes = struct {
	ReportedErrorCount   MetricType
	UserActionDuration   MetricType
	WebRequestCount      MetricType
	WebRequestErrorCount MetricType
}{
	"REPORTED_ERROR_COUNT",
	"USER_ACTION_DURATION",
	"WEB_REQUEST_COUNT",
	"WEB_REQUEST_ERROR_COUNT",
}
