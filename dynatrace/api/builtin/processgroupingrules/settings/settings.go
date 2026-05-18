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

package processgroupingrules

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Settings struct {
	CustomTechnologyName *string                 `json:"customTechnologyName,omitempty"` // Note: Reported only in full-stack, infrastructure and discovery modes.
	Enabled              bool                    `json:"enabled"`                        // This setting is enabled (`true`) or disabled (`false`)
	PgExtraction         ProcessGroupExtractions `json:"pgExtraction"`                   // Define process groups and processes.
	Scope                *string                 `json:"-" scope:"scope"`                // The scope of this setting (HOST, KUBERNETES_CLUSTER, HOST_GROUP). Omit this property if you want to cover the whole environment.
	InsertAfter          string                  `json:"-"`
}

func (me *Settings) Name() string {
	return "process_grouping_rules"
}

func (me *Settings) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"custom_technology_name": {
			Type:        schema.TypeString,
			Description: "Note: Reported only in full-stack, infrastructure and discovery modes.",
			Optional:    true, // nullable
		},
		"enabled": {
			Type:        schema.TypeBool,
			Description: "This setting is enabled (`true`) or disabled (`false`)",
			Required:    true,
		},
		"pg_extraction": {
			Type:        schema.TypeList,
			Description: "Define process groups and processes.",
			Required:    true,
			Elem:        &schema.Resource{Schema: new(ProcessGroupExtractions).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
		"scope": {
			Type:        schema.TypeString,
			Description: "The scope of this setting (HOST, KUBERNETES_CLUSTER, HOST_GROUP). Omit this property if you want to cover the whole environment.",
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
		"custom_technology_name": me.CustomTechnologyName,
		"enabled":                me.Enabled,
		"pg_extraction":          me.PgExtraction,
		"scope":                  me.Scope,
		"insert_after":           me.InsertAfter,
	})
}

func (me *Settings) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"custom_technology_name": &me.CustomTechnologyName,
		"enabled":                &me.Enabled,
		"pg_extraction":          &me.PgExtraction,
		"scope":                  &me.Scope,
		"insert_after":           &me.InsertAfter,
	})
}
