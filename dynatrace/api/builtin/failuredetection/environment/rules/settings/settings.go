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

package rules

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Settings struct {
	Conditions  Conditions `json:"conditions"`            // Conditions
	Description *string    `json:"description,omitempty"` // Rule description
	Enabled     bool       `json:"enabled"`               // This setting is enabled (`true`) or disabled (`false`)
	RuleName    string     `json:"name"`                  // Rule name
	ParameterID string     `json:"parameterId"`           // Failure detection parameters
}

func (me *Settings) Name() string {
	if me.RuleName != "" {
		return me.RuleName
	}

	objID := settings.ObjectID{ID: me.ParameterID}
	err := objID.Decode()
	if err == nil && len(objID.Key) > 0 {
		return objID.Key
	}

	return uuid.NewString()
}

func (me *Settings) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"conditions": {
			Type:        schema.TypeList,
			Description: "Conditions",
			Required:    true,
			Elem:        &schema.Resource{Schema: new(Conditions).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
		"description": {
			Type:        schema.TypeString,
			Description: "Rule description",
			Optional:    true, // nullable
		},
		"enabled": {
			Type:        schema.TypeBool,
			Description: "This setting is enabled (`true`) or disabled (`false`)",
			Required:    true,
		},
		"name": {
			Type:        schema.TypeString,
			Description: "Rule name",
			Required:    true,
		},
		"parameter_id": {
			Type:        schema.TypeString,
			Description: "Failure detection parameters",
			Required:    true,
		},
	}
}

func (me *Settings) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"conditions":   me.Conditions,
		"description":  me.Description,
		"enabled":      me.Enabled,
		"name":         me.RuleName,
		"parameter_id": me.ParameterID,
	})
}

func (me *Settings) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"conditions":   &me.Conditions,
		"description":  &me.Description,
		"enabled":      &me.Enabled,
		"name":         &me.RuleName,
		"parameter_id": &me.ParameterID,
	})
}
