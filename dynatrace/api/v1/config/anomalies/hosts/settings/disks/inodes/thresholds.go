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

package inodes

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Thresholds Custom thresholds for low disk inodes number. If not set, automatic mode is used.
type Thresholds struct {
	FreeInodesPercentage int32 `json:"freeInodesPercentage"` // Alert if percentage of available inodes is lower than *X*% in 3 out of 5 samples.
}

func (me *Thresholds) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"percentage": {
			Type:        schema.TypeInt,
			Required:    true,
			Description: "Alert if percentage of available inodes is lower than *X*% in 3 out of 5 samples",
		},
	}
}

func (me *Thresholds) MarshalHCL(properties hcl.Properties) error {
	return properties.Encode("percentage", me.FreeInodesPercentage)
}

func (me *Thresholds) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.Decode("percentage", &me.FreeInodesPercentage)
}
