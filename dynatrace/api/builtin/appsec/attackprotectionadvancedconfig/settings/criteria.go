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

package attackprotectionadvancedconfig

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Criteria struct {
	AttackType   AttackType `json:"attackType"`             // Possible Values: `ANY`, `CMD_INJECTION`, `JNDI_INJECTION`, `SQL_INJECTION`, `SSRF`
	ProcessGroup *string    `json:"processGroup,omitempty"` // Process group
}

func (me *Criteria) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"attack_type": {
			Type:        schema.TypeString,
			Description: "Possible Values: `ANY`, `CMD_INJECTION`, `JNDI_INJECTION`, `SQL_INJECTION`, `SSRF`",
			Required:    true,
		},
		"process_group": {
			Type:        schema.TypeString,
			Description: "Process group",
			Optional:    true, // nullable
			Deprecated:  "This field has been deprecated",
		},
	}
}

func (me *Criteria) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"attack_type":   me.AttackType,
		"process_group": me.ProcessGroup,
	})
}

func (me *Criteria) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"attack_type":   &me.AttackType,
		"process_group": &me.ProcessGroup,
	})
}
