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

package useractioncustommetrics

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Settings struct {
	Dimensions []string     `json:"dimensions,omitempty"` // Defines the fields that are used as dimensions. A dimension is a collection of reference information about a metric data point that is of interest to your business. Dimensions are parameters like \"application\", \"type\", \"name\". For example, using \"type\" as a dimension allows you to split chart data based on the user action type.
	Enabled    bool         `json:"enabled"`              // This setting is enabled (`true`) or disabled (`false`)
	Filters    Filters      `json:"filters,omitempty"`    // Defines the filters for the user action. Filters apply at the moment of extracting the data and only sessions that satisfy the filtering criteria will be used to extract the custom metrics. You will not be able to modify these filters in the metric data explorer. For example, using \"type equals Xhr\" will give you only data from xhr actions, while forcing the rest of user actions of different types to be ignored.
	MetricKey  string       `json:"metricKey"`            // Metric key
	Value      *MetricValue `json:"value"`                // Defines the type of value to be extracted from the user action. When using **user action counter**, the number of user actions is counted (similar to count(*) when using USQL). When using **user action field value**, the value of a user action field is extracted.
}

func (me *Settings) Name() string {
	return me.MetricKey
}

func (me *Settings) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"dimensions": {
			Type:        schema.TypeList,
			Description: "Defines the fields that are used as dimensions. A dimension is a collection of reference information about a metric data point that is of interest to your business. Dimensions are parameters like \"application\", \"type\", \"name\". For example, using \"type\" as a dimension allows you to split chart data based on the user action type.",
			Optional:    true, // minobjects == 0
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
		"enabled": {
			Type:        schema.TypeBool,
			Description: "This setting is enabled (`true`) or disabled (`false`)",
			Required:    true,
		},
		"filters": {
			Type:        schema.TypeList,
			Description: "Defines the filters for the user action. Filters apply at the moment of extracting the data and only sessions that satisfy the filtering criteria will be used to extract the custom metrics. You will not be able to modify these filters in the metric data explorer. For example, using \"type equals Xhr\" will give you only data from xhr actions, while forcing the rest of user actions of different types to be ignored.",
			Optional:    true, // minobjects == 0
			Elem:        &schema.Resource{Schema: new(Filters).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
		"metric_key": {
			Type:        schema.TypeString,
			Description: "Metric key",
			Required:    true,
		},
		"value": {
			Type:        schema.TypeList,
			Description: "Defines the type of value to be extracted from the user action. When using **user action counter**, the number of user actions is counted (similar to count(*) when using USQL). When using **user action field value**, the value of a user action field is extracted.",
			Required:    true,
			Elem:        &schema.Resource{Schema: new(MetricValue).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
	}
}

func (me *Settings) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"dimensions": me.Dimensions,
		"enabled":    me.Enabled,
		"filters":    me.Filters,
		"metric_key": me.MetricKey,
		"value":      me.Value,
	})
}

func (me *Settings) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"dimensions": &me.Dimensions,
		"enabled":    &me.Enabled,
		"filters":    &me.Filters,
		"metric_key": &me.MetricKey,
		"value":      &me.Value,
	})
}
