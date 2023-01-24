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

package browser

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// RestrictionSettings Settings for restricting certain browser type, version, platform and, comparator. It also restricts the mode
type RestrictionSettings struct {
	Mode                RestrictionMode `json:"mode"`                          // The mode of the list of browser restrictions. Possible values area `EXCLUDE` and `INCLUDE`.
	BrowserRestrictions Restrictions    `json:"browserRestrictions,omitempty"` // A list of browser restrictions
}

func (me *RestrictionSettings) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"mode": {
			Type:        schema.TypeString,
			Description: "The mode of the list of browser restrictions. Possible values area `EXCLUDE` and `INCLUDE`.",
			Required:    true,
		},
		"restrictions": {
			Type:        schema.TypeList,
			Description: "A list of browser restrictions",
			Optional:    true,
			MaxItems:    1,
			Elem:        &schema.Resource{Schema: new(Restrictions).Schema()},
		},
	}
}

func (me *RestrictionSettings) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"mode":         me.Mode,
		"restrictions": me.BrowserRestrictions,
	})
}

func (me *RestrictionSettings) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"mode":         &me.Mode,
		"restrictions": &me.BrowserRestrictions,
	})
}
