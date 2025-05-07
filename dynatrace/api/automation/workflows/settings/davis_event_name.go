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

package workflows

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

type DavisEventNames []*DavisEventName

func (me *DavisEventNames) Schema(prefix string) map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:        schema.TypeList,
			Description: "A combination of name and match",
			MinItems:    1,
			Optional:    true,
			Elem:        &schema.Resource{Schema: new(DavisEventName).Schema(prefix + ".0.name")},
		},
	}
}

func (me DavisEventNames) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeSlice("name", me)
}

func (me *DavisEventNames) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeSlice("name", me)
}

type DavisEventName struct {
	Name  string `json:"name"`
	Match string `json:"match"`
}

func (me *DavisEventName) Schema(prefix string) map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:        schema.TypeString,
			Description: "The event name",
			Required:    true,
		},
		"match": {
			Type:         schema.TypeString,
			Description:  "Possible values: `equals` and `contains`. The Davis event name must equal or contain the string provided in attribute `name`",
			Required:     true,
			ValidateFunc: validation.StringInSlice([]string{"equals", "contains"}, false),
		},
	}
}

func (me *DavisEventName) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"name":  me.Name,
		"match": me.Match,
	})
}

func (me *DavisEventName) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"name":  &me.Name,
		"match": &me.Match,
	})
}
