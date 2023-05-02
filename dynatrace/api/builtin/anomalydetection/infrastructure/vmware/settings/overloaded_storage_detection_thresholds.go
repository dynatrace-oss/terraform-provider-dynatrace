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

package vmware

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// OverloadedStorageDetectionThresholds. Alert if the condition is met in 3 out of 5 samples
type OverloadedStorageDetectionThresholds struct {
	CommandAbortsNumber int `json:"commandAbortsNumber"` // Number of command aborts is higher than
}

func (me *OverloadedStorageDetectionThresholds) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"command_aborts_number": {
			Type:        schema.TypeInt,
			Description: "Number of command aborts is higher than",
			Required:    true,
		},
	}
}

func (me *OverloadedStorageDetectionThresholds) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"command_aborts_number": me.CommandAbortsNumber,
	})
}

func (me *OverloadedStorageDetectionThresholds) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"command_aborts_number": &me.CommandAbortsNumber,
	})
}
