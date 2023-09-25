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
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Settings struct {
	AttackHandling *AttackHandling `json:"attackHandling"` // Step 2: Define attack control for chosen criteria
	Criteria       *Criteria       `json:"criteria"`       // Step 1: Define criteria. Please specify at least one of source IP or attack pattern.
	Enabled        bool            `json:"enabled"`        // This setting is enabled (`true`) or disabled (`false`)
	Metadata       *Metadata       `json:"metadata"`       // Step 3: Leave comment
}

func (me *Settings) Name() string {
	return uuid.NewString()
}

func (me *Settings) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"attack_handling": {
			Type:        schema.TypeList,
			Description: "Step 2: Define attack control for chosen criteria",
			Required:    true,
			Elem:        &schema.Resource{Schema: new(AttackHandling).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
		"criteria": {
			Type:        schema.TypeList,
			Description: "Step 1: Define criteria. Please specify at least one of source IP or attack pattern.",
			Required:    true,
			Elem:        &schema.Resource{Schema: new(Criteria).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
		"enabled": {
			Type:        schema.TypeBool,
			Description: "This setting is enabled (`true`) or disabled (`false`)",
			Required:    true,
		},
		"metadata": {
			Type:        schema.TypeList,
			Description: "Step 3: Leave comment",
			Required:    true,
			Elem:        &schema.Resource{Schema: new(Metadata).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
	}
}

func (me *Settings) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"attack_handling": me.AttackHandling,
		"criteria":        me.Criteria,
		"enabled":         me.Enabled,
		"metadata":        me.Metadata,
	})
}

func (me *Settings) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"attack_handling": &me.AttackHandling,
		"criteria":        &me.Criteria,
		"enabled":         &me.Enabled,
		"metadata":        &me.Metadata,
	})
}
