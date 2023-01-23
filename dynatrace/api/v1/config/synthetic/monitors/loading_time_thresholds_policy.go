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

type LoadingTimeThresholdsPolicy struct {
	Enabled    bool                  `json:"enabled"`    // Performance threshold is enabled (`true`) or disabled (`false`)
	Thresholds LoadingTimeThresholds `json:"thresholds"` // The list of performance threshold rules
}

func (me *LoadingTimeThresholdsPolicy) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"enabled": {
			Type:        schema.TypeBool,
			Description: "Performance threshold is enabled (`true`) or disabled (`false`)",
			Optional:    true,
		},
		"thresholds": {
			Type:        schema.TypeList,
			Description: "The list of performance threshold rules",
			Optional:    true,
			Elem: &schema.Resource{
				Schema: new(LoadingTimeThresholds).Schema(),
			},
		},
	}
}

func (me *LoadingTimeThresholdsPolicy) MarshalHCL(properties hcl.Properties) error {
	if err := properties.Encode("enabled", me.Enabled); err != nil {
		return err
	}
	if err := properties.Encode("thresholds", me.Thresholds); err != nil {
		return err
	}
	return nil
}

func (me *LoadingTimeThresholdsPolicy) UnmarshalHCL(decoder hcl.Decoder) error {
	if err := decoder.Decode("enabled", &me.Enabled); err != nil {
		return err
	}
	if err := decoder.Decode("thresholds", &me.Thresholds); err != nil {
		return err
	}
	if me.Thresholds == nil {
		me.Thresholds = LoadingTimeThresholds{}
	}
	return nil
}
