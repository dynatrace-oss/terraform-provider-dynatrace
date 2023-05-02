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

package diskrules

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type DiskNameFilter struct {
	Operator DiskNameFilterOperator `json:"operator"`        // Possible Values: `CONTAINS`, `DOES_NOT_CONTAIN`, `DOES_NOT_EQUAL`, `DOES_NOT_START_WITH`, `EQUALS`, `STARTS_WITH`
	Value    *string                `json:"value,omitempty"` // Matching text
}

func (me *DiskNameFilter) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"operator": {
			Type:        schema.TypeString,
			Description: "Possible Values: `CONTAINS`, `DOES_NOT_CONTAIN`, `DOES_NOT_EQUAL`, `DOES_NOT_START_WITH`, `EQUALS`, `STARTS_WITH`",
			Required:    true,
		},
		"value": {
			Type:        schema.TypeString,
			Description: "Matching text",
			Optional:    true, // nullable
		},
	}
}

func (me *DiskNameFilter) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"operator": me.Operator,
		"value":    me.Value,
	})
}

func (me *DiskNameFilter) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"operator": &me.Operator,
		"value":    &me.Value,
	})
}
