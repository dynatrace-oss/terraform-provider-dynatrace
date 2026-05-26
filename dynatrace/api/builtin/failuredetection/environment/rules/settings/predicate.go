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
	"fmt"
	"slices"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// predicate. A predicate that tests a condition attribute value. The available fields depend on the predicate type: string predicates use `textValues` and `caseSensitive`, service type predicates use `serviceType`, management zone predicates use `managementZones`, and tag predicates use `tags` or `tagKeys`.
type Predicate struct {
	CaseSensitive   *bool          `json:"caseSensitive,omitempty"`   // If `true`, the string comparison is case-sensitive. Default: `false`.
	ManagementZones []string       `json:"managementZones,omitempty"` // A set of management zone references. The rule matches if the service belongs to all specified management zones. Only applicable for predicate type `MANAGEMENT_ZONES_CONTAINS_ALL`.
	PredicateType   string         `json:"predicateType"`             // The type of predicate to apply. Available types depend on the condition attribute:\n * `SERVICE_NAME` and `PG_NAME` support `STRING_EQUALS`, `STARTS_WITH`, `ENDS_WITH`, `CONTAINS`;\n * `SERVICE_TYPE` supports `SERVICE_TYPE_EQUALS`;\n * `SERVICE_MANAGEMENT_ZONE` supports `MANAGEMENT_ZONES_CONTAINS_ALL`;\n * `SERVICE_TAG` and `PG_TAG` support `TAG_EQUALS` and `TAG_KEY_EQUALS`.
	ServiceType     []ServiceTypes `json:"serviceType,omitempty"`     // A set of service types to match against. The rule matches if the service type is contained in this set. Only applicable for predicate type `SERVICE_TYPE_EQUALS`. Possible values: `CICS`, `CICSInteraction`, `CustomApplication`, `Database`, `EnterpriseServiceBus`, `External`, `IMS`, `IMSInteraction`, `Messaging`, `Method`, `Mobile`, `Process`, `QueueInteraction`, `QueueListener`, `RMI`, `RemoteCall`, `SaasVendor`, `WebRequest`, `WebService`, `WebSite`, `zOSConnect`
	TagKeys         []string       `json:"tagKeys,omitempty"`         // A set of tag keys to match. The rule matches if the entity has tags with all specified keys, regardless of tag value. Only applicable for predicate type `TAG_KEY_EQUALS`.
	Tags            []string       `json:"tags,omitempty"`            // A set of tags to match exactly. The rule matches if the entity has all specified tags (both key and value must match). Only applicable for predicate type `TAG_EQUALS`.
	TextValues      []string       `json:"textValues,omitempty"`      // A list of text values to match against. The rule matches if the attribute value matches any of these values according to the predicate type.
}

func (me *Predicate) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"case_sensitive": {
			Type:        schema.TypeBool,
			Description: "If `true`, the string comparison is case-sensitive. Default: `false`.",
			Optional:    true, // precondition
		},
		"management_zones": {
			Type:        schema.TypeSet,
			Description: "A set of management zone references. The rule matches if the service belongs to all specified management zones. Only applicable for predicate type `MANAGEMENT_ZONES_CONTAINS_ALL`.",
			Optional:    true, // precondition
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
		"predicate_type": {
			Type:        schema.TypeString,
			Description: "The type of predicate to apply. Available types depend on the condition attribute:\n * `SERVICE_NAME` and `PG_NAME` support `STRING_EQUALS`, `STARTS_WITH`, `ENDS_WITH`, `CONTAINS`;\n * `SERVICE_TYPE` supports `SERVICE_TYPE_EQUALS`;\n * `SERVICE_MANAGEMENT_ZONE` supports `MANAGEMENT_ZONES_CONTAINS_ALL`;\n * `SERVICE_TAG` and `PG_TAG` support `TAG_EQUALS` and `TAG_KEY_EQUALS`.",
			Required:    true,
		},
		"service_type": {
			Type:        schema.TypeSet,
			Description: "A set of service types to match against. The rule matches if the service type is contained in this set. Only applicable for predicate type `SERVICE_TYPE_EQUALS`. Possible values: `CICS`, `CICSInteraction`, `CustomApplication`, `Database`, `EnterpriseServiceBus`, `External`, `IMS`, `IMSInteraction`, `Messaging`, `Method`, `Mobile`, `Process`, `QueueInteraction`, `QueueListener`, `RMI`, `RemoteCall`, `SaasVendor`, `WebRequest`, `WebService`, `WebSite`, `zOSConnect`",
			Optional:    true, // precondition & minobjects == 0
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
		"tag_keys": {
			Type:        schema.TypeSet,
			Description: "A set of tag keys to match. The rule matches if the entity has tags with all specified keys, regardless of tag value. Only applicable for predicate type `TAG_KEY_EQUALS`.",
			Optional:    true, // precondition
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
		"tags": {
			Type:        schema.TypeSet,
			Description: "A set of tags to match exactly. The rule matches if the entity has all specified tags (both key and value must match). Only applicable for predicate type `TAG_EQUALS`.",
			Optional:    true, // precondition
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
		"text_values": {
			Type:        schema.TypeSet,
			Description: "A list of text values to match against. The rule matches if the attribute value matches any of these values according to the predicate type.",
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

func (me *Predicate) HandlePreconditions() error {
	if (me.CaseSensitive == nil) && (slices.Contains([]string{"STRING_EQUALS", "STARTS_WITH", "ENDS_WITH", "CONTAINS"}, string(me.PredicateType))) {
		me.CaseSensitive = new(false)
	}
	if (me.CaseSensitive != nil) && (!slices.Contains([]string{"STRING_EQUALS", "STARTS_WITH", "ENDS_WITH", "CONTAINS"}, string(me.PredicateType))) {
		return fmt.Errorf("'case_sensitive' must not be specified unless 'predicate_type' is one of ['STRING_EQUALS', 'STARTS_WITH', 'ENDS_WITH', 'CONTAINS']; got 'predicate_type'='%v'", me.PredicateType)
	}
	// ---- ManagementZones []string -> {"expectedValues":["MANAGEMENT_ZONES_CONTAINS_ALL"],"property":"predicateType","type":"IN"}
	// ---- ServiceType []ServiceTypes -> {"expectedValues":["SERVICE_TYPE_EQUALS"],"property":"predicateType","type":"IN"}
	// ---- TagKeys []string -> {"expectedValues":["TAG_KEY_EQUALS"],"property":"predicateType","type":"IN"}
	// ---- Tags []string -> {"expectedValues":["TAG_EQUALS"],"property":"predicateType","type":"IN"}
	// ---- TextValues []string -> {"expectedValues":["STRING_EQUALS","STARTS_WITH","ENDS_WITH","CONTAINS"],"property":"predicateType","type":"IN"}
	return nil
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
