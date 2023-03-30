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
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/opt"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"golang.org/x/exp/slices"
)

type Predicate struct {
	CaseSensitive   *bool          `json:"caseSensitive,omitempty"`   // Case sensitive
	ManagementZones []string       `json:"managementZones,omitempty"` // Management zones
	PredicateType   string         `json:"predicateType"`             // Predicate type
	ServiceType     []ServiceTypes `json:"serviceType,omitempty"`     // Service types
	TagKeys         []string       `json:"tagKeys,omitempty"`         // Tag keys
	Tags            []string       `json:"tags,omitempty"`            // Tags (exact match)
	TextValues      []string       `json:"textValues,omitempty"`      // Names
}

func (me *Predicate) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"case_sensitive": {
			Type:        schema.TypeBool,
			Description: "Case sensitive",
			Optional:    true, // precondition
		},
		"management_zones": {
			Type:        schema.TypeSet,
			Description: "Management zones",
			Optional:    true, // precondition
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
		"predicate_type": {
			Type:        schema.TypeString,
			Description: "Predicate type",
			Required:    true,
		},
		"service_type": {
			Type:        schema.TypeSet,
			Description: "Service types",
			Optional:    true, // precondition & minobjects == 0
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
		"tag_keys": {
			Type:        schema.TypeSet,
			Description: "Tag keys",
			Optional:    true, // precondition
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
		"tags": {
			Type:        schema.TypeSet,
			Description: "Tags (exact match)",
			Optional:    true, // precondition
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
		"text_values": {
			Type:        schema.TypeSet,
			Description: "Names",
			Optional:    true, // precondition
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
	}
}

func (me *Predicate) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"case_sensitive":   me.CaseSensitive,
		"management_zones": me.ManagementZones,
		"predicate_type":   me.PredicateType,
		"service_type":     me.ServiceType,
		"tag_keys":         me.TagKeys,
		"tags":             me.Tags,
		"text_values":      me.TextValues,
	})
}

func (me *Predicate) HandlePreconditions() {
	if me.CaseSensitive == nil && slices.Contains([]string{"STRING_EQUALS", "STARTS_WITH", "ENDS_WITH", "CONTAINS"}, string(me.PredicateType)) {
		me.CaseSensitive = opt.NewBool(false)
	}
	// ---- ManagementZones []string -> {"expectedValues":["MANAGEMENT_ZONES_CONTAINS_ALL"],"property":"predicateType","type":"IN"}
	// ---- ServiceType []ServiceTypes -> {"expectedValues":["SERVICE_TYPE_EQUALS"],"property":"predicateType","type":"IN"}
	// ---- TagKeys []string -> {"expectedValues":["TAG_KEY_EQUALS"],"property":"predicateType","type":"IN"}
	// ---- Tags []string -> {"expectedValues":["TAG_EQUALS"],"property":"predicateType","type":"IN"}
	// ---- TextValues []string -> {"expectedValues":["STRING_EQUALS","STARTS_WITH","ENDS_WITH","CONTAINS"],"property":"predicateType","type":"IN"}
}

func (me *Predicate) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"case_sensitive":   &me.CaseSensitive,
		"management_zones": &me.ManagementZones,
		"predicate_type":   &me.PredicateType,
		"service_type":     &me.ServiceType,
		"tag_keys":         &me.TagKeys,
		"tags":             &me.Tags,
		"text_values":      &me.TextValues,
	})
}
