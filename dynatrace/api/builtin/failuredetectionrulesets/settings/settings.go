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

package failuredetectionrulesets

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Settings struct {
	Enabled     bool     `json:"enabled"`         // This setting is enabled (`true`) or disabled (`false`)
	Ruleset     *Ruleset `json:"ruleset"`         // Ruleset
	Scope       *string  `json:"-" scope:"scope"` // The scope of this setting (CLOUD_APPLICATION_NAMESPACE, KUBERNETES_CLUSTER, HOST_GROUP). Omit this property if you want to cover the whole environment.
	InsertAfter string   `json:"-"`
}

func (me *Settings) Name() string {
	if me.Scope == nil {
		return "Environment_" + me.Ruleset.RulesetName
	}
	return *me.Scope + "_" + me.Ruleset.RulesetName
}

func (me *Settings) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"enabled": {
			Type:        schema.TypeBool,
			Description: "This setting is enabled (`true`) or disabled (`false`)",
			Required:    true,
		},
		"ruleset": {
			Type:        schema.TypeList,
			Description: "Ruleset",
			Required:    true,
			Elem:        &schema.Resource{Schema: new(Ruleset).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
		"scope": {
			Type:        schema.TypeString,
			Description: "The scope of this setting (CLOUD_APPLICATION_NAMESPACE, KUBERNETES_CLUSTER, HOST_GROUP). Omit this property if you want to cover the whole environment.",
			Optional:    true,
			Default:     "environment",
		},
		"insert_after": {
			Type:        schema.TypeString,
			Description: "Because this resource allows for ordering you may specify the ID of the resource instance that comes before this instance regarding order. If not specified when creating the setting will be added to the end of the list. If not specified during update the order will remain untouched",
			Optional:    true,
			Computed:    true,
		},
	}
}

func (me *Settings) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"enabled":      me.Enabled,
		"ruleset":      me.Ruleset,
		"scope":        me.Scope,
		"insert_after": me.InsertAfter,
	})
}

func (me *Settings) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"enabled":      &me.Enabled,
		"ruleset":      &me.Ruleset,
		"scope":        &me.Scope,
		"insert_after": &me.InsertAfter,
	})
}
