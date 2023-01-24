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

package alerting

import (
	"encoding/json"
	"sort"
	"strings"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/opt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Profile represents an Alerting Profile in Dynatrace
type Profile struct {
	ID             string        `json:"-"`
	Name           string        `json:"name"`                    // The name of the Alerting Profile
	ManagementZone *string       `json:"managementZone"`          // Define management zone filter for profile
	SeverityRules  SeverityRules `json:"severityRules,omitempty"` // Define severity rules for profile. A maximum of 100 severity rules is allowed.
	EventFilters   EventFilters  `json:"eventFilters,omitempty"`  // Define event filters for profile. A maximum of 100 event filters is allowed.
	LegacyID       *string       `json:"-"`
}

// Schema provides a map for terraform, containing all the current supported properties
func (me *Profile) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:        schema.TypeString,
			Description: "The name of the alerting profile, displayed in the UI",
			Required:    true,
		},
		"management_zone": {
			Type:        schema.TypeString,
			Description: "The ID of the management zone to which the alerting profile applies",
			Optional:    true,
		},
		"rules": {
			Type:        schema.TypeList,
			Description: "A list of rules for management zone usage.  Each rule is evaluated independently of all other rules",
			Optional:    true,
			MaxItems:    1,
			MinItems:    1,
			Elem:        &schema.Resource{Schema: new(SeverityRules).Schema()},
		},
		"filters": {
			Type:        schema.TypeList,
			Description: "The list of event filters.  For all filters that are *negated* inside of these event filters, that is all `Predefined` as well as `Custom` (Title and/or Description) ones the AND logic applies. For all *non-negated* ones the OR logic applies. Between these two groups, negated and non-negated, the AND logic applies.  If you specify both severity rule and event filter, the AND logic applies",
			Optional:    true,
			MaxItems:    1,
			MinItems:    1,
			Elem:        &schema.Resource{Schema: new(EventFilters).Schema()},
		},
		"legacy_id": {
			Type:        schema.TypeString,
			Description: "The ID of this setting when referred to by the Config REST API V1",
			Computed:    true,
			Optional:    true,
		},
	}
}

// EnsurePredictableOrder is currently an inconvenient necessity. The REST API does not guarantee to deliever the Severity Rules for an Alerting Profile in a predictable order.
// Terraform does however insist on the same order every time. Therefore the the Severity Rules are getting ordered based on their JSON representation, before any HCL code is getting produced.
func (me *Profile) EnsurePredictableOrder() {
	if len(me.SeverityRules) == 0 {
		return
	}
	conds := SeverityRules{}
	condStrings := sort.StringSlice{}
	for _, entry := range me.SeverityRules {
		condBytes, _ := json.Marshal(entry)
		condStrings = append(condStrings, string(condBytes))
	}
	condStrings.Sort()
	for _, condString := range condStrings {
		cond := SeverityRule{}
		json.Unmarshal([]byte(condString), &cond)
		conds = append(conds, &cond)
	}
	me.SeverityRules = conds
}

// MarshalHCL produces HCL structures for Terraform
func (me *Profile) MarshalHCL(properties hcl.Properties) error {
	if err := properties.Encode("name", me.Name); err != nil {
		return err
	}
	if err := properties.Encode("management_zone", me.ManagementZone); err != nil {
		return err
	}
	if err := properties.Encode("legacy_id", me.LegacyID); err != nil {
		return err
	}
	if len(me.SeverityRules) > 0 {
		me.EnsurePredictableOrder()
		marshalled := hcl.Properties{}
		err := me.SeverityRules.MarshalHCL(marshalled)
		if err != nil {
			return err
		}
		properties["rules"] = []any{marshalled}
	}
	if len(me.EventFilters) > 0 {
		filters := append(EventFilters{}, me.EventFilters...)
		sort.Slice(filters, func(i, j int) bool {
			d1, _ := json.Marshal(filters[i])
			d2, _ := json.Marshal(filters[j])
			cmp := strings.Compare(string(d1), string(d2))
			return (cmp == -1)
		})
		marshalled := hcl.Properties{}
		err := me.EventFilters.MarshalHCL(marshalled)
		if err != nil {
			return err
		}
		properties["filters"] = []any{marshalled}
	}
	return nil
}

// UnmarshalHCL decodes HCL code and fills this object
func (me *Profile) UnmarshalHCL(decoder hcl.Decoder) error {
	if value, ok := decoder.GetOk("name"); ok {
		me.Name = value.(string)
	}
	if value, ok := decoder.GetOk("management_zone"); ok {
		me.ManagementZone = opt.NewString(value.(string))
	}
	if value, ok := decoder.GetOk("legacy_id"); ok {
		me.LegacyID = opt.NewString(value.(string))
	}
	if _, ok := decoder.GetOk("rules.#"); ok {
		me.SeverityRules = SeverityRules{}
		if err := me.SeverityRules.UnmarshalHCL(hcl.NewDecoder(decoder, "rules", 0)); err != nil {
			return err
		}
	}
	if _, ok := decoder.GetOk("filters.#"); ok {
		me.EventFilters = EventFilters{}
		if err := me.EventFilters.UnmarshalHCL(hcl.NewDecoder(decoder, "filters", 0)); err != nil {
			return err
		}
	}
	return nil
}
