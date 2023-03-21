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

package requesterrors

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type FilterSettings struct {
	Filter *UrlFilter `json:"filter,omitempty"` // Possible Values: `BEGINS_WITH`, `CONTAINS`, `ENDS_WITH`, `EQUALS`
	Url    *string    `json:"url,omitempty"`
}

func (me *FilterSettings) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"filter": {
			Type:        schema.TypeString,
			Description: "Possible Values: `BEGINS_WITH`, `CONTAINS`, `ENDS_WITH`, `EQUALS`",
			Optional:    true,
		},
		"url": {
			Type:        schema.TypeString,
			Description: "no documentation available",
			Optional:    true,
		},
	}
}

func (me *FilterSettings) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"filter": me.Filter,
		"url":    me.Url,
	})
}

func (me *FilterSettings) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"filter": &me.Filter,
		"url":    &me.Url,
	})
}
