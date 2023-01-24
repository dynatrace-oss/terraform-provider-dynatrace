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

package diskevents

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// DiskNameFilter Narrows the rule usage down to disks, matching the specified criteria.
type DiskNameFilter struct {
	Operator Operator `json:"operator"` // Comparison operator.
	Value    string   `json:"value"`    // Value to compare to.
}

func (me *DiskNameFilter) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"operator": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Possible values are: `CONTAINS`, `DOES_NOT_CONTAIN`, `DOES_NOT_EQUAL`, `DOES_NOT_START_WITH`, `EQUALS` and `STARTS_WITH`",
		},
		"value": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Value to compare to",
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
	if value, ok := decoder.GetOk("operator"); ok {
		me.Operator = Operator(value.(string))
	}
	if value, ok := decoder.GetOk("value"); ok {
		me.Value = value.(string)
	}
	return nil
}

// Operator Comparison operator.
type Operator string

// Operators offers the known enum values
var Operators = struct {
	Contains         Operator
	DoesNotContain   Operator
	DoesNotEqual     Operator
	DoesNotStartWith Operator
	Equals           Operator
	StartsWith       Operator
}{
	"CONTAINS",
	"DOES_NOT_CONTAIN",
	"DOES_NOT_EQUAL",
	"DOES_NOT_START_WITH",
	"EQUALS",
	"STARTS_WITH",
}
