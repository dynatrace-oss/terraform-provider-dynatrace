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

package attackprotectionallowlistconfig

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Criteria struct {
	AttackPattern *string `json:"attackPattern,omitempty"` // Only consider attacks matching the specified pattern.
	SourceIp      *string `json:"sourceIp,omitempty"`      // Source IP
}

func (me *Criteria) IsEmpty() bool {
	return me == nil || (me.AttackPattern == nil && me.SourceIp == nil)
}

func (me *Criteria) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"attack_pattern": {
			Type:        schema.TypeString,
			Description: "Only consider attacks matching the specified pattern.",
			Optional:    true, // nullable
		},
		"source_ip": {
			Type:        schema.TypeString,
			Description: "Source IP",
			Optional:    true, // nullable
		},
	}
}

func (me *Criteria) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"attack_pattern": me.AttackPattern,
		"source_ip":      me.SourceIp,
	})
}

func (me *Criteria) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"attack_pattern": &me.AttackPattern,
		"source_ip":      &me.SourceIp,
	})
}
