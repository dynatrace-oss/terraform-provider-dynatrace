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

package disks

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type DiskLowInodesDetectionThresholds struct {
	FreeInodesPercentage int `json:"freeInodesPercentage"` // Alert if the percentage of available inodes is lower than this threshold in 3 out of 5 samples
}

func (me *DiskLowInodesDetectionThresholds) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"free_inodes_percentage": {
			Type:        schema.TypeInt,
			Description: "Alert if the percentage of available inodes is lower than this threshold in 3 out of 5 samples",
			Required:    true,
		},
	}
}

func (me *DiskLowInodesDetectionThresholds) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"free_inodes_percentage": me.FreeInodesPercentage,
	})
}

func (me *DiskLowInodesDetectionThresholds) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"free_inodes_percentage": &me.FreeInodesPercentage,
	})
}
