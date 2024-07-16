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

package monitors

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type PerformanceThresholds struct {
	Enabled    bool       `json:"enabled"`              // Performance threshold is enabled (true) or disabled (false)
	Thresholds Thresholds `json:"thresholds,omitempty"` // The list of performance threshold rules
}

func (me *PerformanceThresholds) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"enabled": {
			Type:        schema.TypeBool,
			Description: "Performance threshold is enabled (true) or disabled (false)",
			Optional:    true,
			Default:     true,
		},
		"thresholds": {
			Type:        schema.TypeList,
			Description: "The list of performance threshold rules",
			Optional:    true,
			Elem:        &schema.Resource{Schema: new(Thresholds).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
	}
}

func (me *PerformanceThresholds) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"enabled":    me.Enabled,
		"thresholds": me.Thresholds,
	})
}

func (me *PerformanceThresholds) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"enabled":    &me.Enabled,
		"thresholds": &me.Thresholds,
	})
}
