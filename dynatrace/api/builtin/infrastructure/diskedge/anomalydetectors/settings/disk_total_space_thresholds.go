/**
* @license
* Copyright 2026 Dynatrace LLC
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

package anomalydetectors

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type DiskTotalSpaceThresholds struct {
	ThresholdAbove *int `json:"thresholdAbove,omitempty"` // If this field is empty then there is no lower limit. Minimum total disk space in GiB
	ThresholdBelow *int `json:"thresholdBelow,omitempty"` // If this field is empty then there is no upper limit. Maximum total disk space in GiB
}

func (me *DiskTotalSpaceThresholds) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"threshold_above": {
			Type:        schema.TypeInt,
			Description: "If this field is empty then there is no lower limit. Minimum total disk space in GiB",
			Optional:    true, // nullable
		},
		"threshold_below": {
			Type:        schema.TypeInt,
			Description: "If this field is empty then there is no upper limit. Maximum total disk space in GiB",
			Optional:    true, // nullable
		},
	}
}

func (me *DiskTotalSpaceThresholds) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"threshold_above": me.ThresholdAbove,
		"threshold_below": me.ThresholdBelow,
	})
}

func (me *DiskTotalSpaceThresholds) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"threshold_above": &me.ThresholdAbove,
		"threshold_below": &me.ThresholdBelow,
	})
}
