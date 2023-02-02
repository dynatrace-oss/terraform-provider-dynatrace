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

package perdiskoverride

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type DiskSlowWritesAndReadsDetection struct {
	CustomThresholds *DiskSlowWritesAndReadsDetectionThresholds `json:"customThresholds,omitempty"`
	DetectionMode    *DetectionMode                             `json:"detectionMode,omitempty"` // Detection mode for slow running disks
	Enabled          bool                                       `json:"enabled"`                 // Detect slow-running disks
}

func (me *DiskSlowWritesAndReadsDetection) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"custom_thresholds": {
			Type:        schema.TypeList,
			Description: "no documentation available",
			Optional:    true,
			Elem:        &schema.Resource{Schema: new(DiskSlowWritesAndReadsDetectionThresholds).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
		"detection_mode": {
			Type:        schema.TypeString,
			Description: "Detection mode for slow running disks",
			Optional:    true,
		},
		"enabled": {
			Type:        schema.TypeBool,
			Description: "Detect slow-running disks",
			Required:    true,
		},
	}
}

func (me *DiskSlowWritesAndReadsDetection) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"custom_thresholds": me.CustomThresholds,
		"detection_mode":    me.DetectionMode,
		"enabled":           me.Enabled,
	})
}

func (me *DiskSlowWritesAndReadsDetection) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"custom_thresholds": &me.CustomThresholds,
		"detection_mode":    &me.DetectionMode,
		"enabled":           &me.Enabled,
	})
}
