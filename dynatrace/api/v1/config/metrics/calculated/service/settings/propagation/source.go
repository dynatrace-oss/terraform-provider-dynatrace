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

package propagation

import (
	"encoding/json"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/xjson"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Source Defines valid sources of request attributes for conditions or placeholders.
type Source struct {
	ManagementZone *string                    `json:"managementZone,omitempty"` // Use only request attributes from services that belong to this management zone.. Use either this or `serviceTag`
	ServiceTag     *UniversalTag              `json:"serviceTag,omitempty"`     // Use only request attributes from services that have this tag. Use either this or `managementZone`
	Unknowns       map[string]json.RawMessage `json:"-"`
}

func (me *Source) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"management_zone": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Use only request attributes from services that belong to this management zone.. Use either this or `serviceTag`",
		},
		"service_tag": {
			Type:        schema.TypeList,
			Optional:    true,
			MaxItems:    1,
			MinItems:    1,
			Description: "Use only request attributes from services that have this tag. Use either this or `managementZone`",
			Elem:        &schema.Resource{Schema: new(UniversalTag).Schema()},
		},
		"unknowns": {
			Type:        schema.TypeString,
			Description: "allows for configuring properties that are not explicitly supported by the current version of this provider",
			Optional:    true,
		},
	}
}

func (me *Source) MarshalHCL(properties hcl.Properties) error {
	if err := properties.Unknowns(me.Unknowns); err != nil {
		return err
	}
	return properties.EncodeAll(map[string]any{
		"management_zone": me.ManagementZone,
		"service_tag":     me.ServiceTag,
		"unknowns":        me.Unknowns,
	})
}

func (me *Source) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"management_zone": &me.ManagementZone,
		"service_tag":     &me.ServiceTag,
		"unknowns":        &me.Unknowns,
	})
}

func (me *Source) MarshalJSON() ([]byte, error) {
	properties := xjson.NewProperties(me.Unknowns)
	if err := properties.MarshalAll(map[string]any{
		"managementZone": me.ManagementZone,
		"serviceTag":     me.ServiceTag,
	}); err != nil {
		return nil, err
	}
	return json.Marshal(properties)
}

func (me *Source) UnmarshalJSON(data []byte) error {
	properties := xjson.NewProperties(me.Unknowns)
	if err := json.Unmarshal(data, &properties); err != nil {
		return err
	}
	return properties.UnmarshalAll(map[string]any{
		"managementZone": &me.ManagementZone,
		"serviceTag":     &me.ServiceTag,
	})
}
