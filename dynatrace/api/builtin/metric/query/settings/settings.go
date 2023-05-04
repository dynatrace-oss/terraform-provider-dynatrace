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

package query

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Settings struct {
	MetricSelector string `json:"metricSelector"`     // Query
	MetricID       string `json:"-" scope:"metricId"` // The scope of this setting (metric)
}

func (me *Settings) Name() string {
	return me.MetricID
}

func (me *Settings) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"metric_selector": {
			Type:             schema.TypeString,
			Description:      "Query",
			Required:         true,
			DiffSuppressFunc: hcl.SuppressJSONorEOT,
		},
		"metric_id": {
			Type:        schema.TypeString,
			Description: "The scope of this setting (metric)",
			Required:    true,
		},
	}
}

func (me *Settings) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"metric_selector": me.MetricSelector,
		"metric_id":       me.MetricID,
	})
}

func (me *Settings) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"metric_selector": &me.MetricSelector,
		"metric_id":       &me.MetricID,
	})
}
