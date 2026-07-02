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

package settings

import (
	"fmt"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type TagFilter struct {
	Key       string
	Value     string
	Condition string // INCLUDE | EXCLUDE
}

type TagFilters []*TagFilter

func (me *TagFilter) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"key": {
			Type:        schema.TypeString,
			Description: "AWS tag key.",
			Required:    true,
		},
		"value": {
			Type:        schema.TypeString,
			Description: "AWS tag value to match.",
			Required:    true,
		},
		"condition": {
			Type:        schema.TypeString,
			Description: "`INCLUDE` to only monitor matching resources, `EXCLUDE` to skip them.",
			Required:    true,
			ValidateFunc: func(i any, k string) (warnings []string, errs []error) {
				v, _ := i.(string)
				if v != "INCLUDE" && v != "EXCLUDE" {
					errs = append(errs, fmt.Errorf("%s must be INCLUDE or EXCLUDE, got %q", k, v))
				}
				return
			},
		},
	}
}

func (me *TagFilter) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"key":       me.Key,
		"value":     me.Value,
		"condition": me.Condition,
	})
}

func (me *TagFilter) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"key":       &me.Key,
		"value":     &me.Value,
		"condition": &me.Condition,
	})
}
