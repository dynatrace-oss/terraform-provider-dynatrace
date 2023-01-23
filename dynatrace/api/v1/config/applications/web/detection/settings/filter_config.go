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

package detection

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type FilterConfig struct {
	Pattern                string `json:"pattern"`
	ApplicationMatchType   string `json:"applicationMatchType"`
	ApplicationMatchTarget string `json:"applicationMatchTarget"`
}

func (me *FilterConfig) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"pattern": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The value to look for with the application detection rule",
		},
		"application_match_type": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The operator used for matching the application detection rule, possible values are `BEGINS_WITH`, `CONTAINS`, `ENDS_WITH`, `EQUALS`, `MATCHES`",
		},
		"application_match_target": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Where to look for the pattern value, possible values are `DOMAIN` or `URL`",
		},
	}
}

func (me *FilterConfig) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"pattern":                  me.Pattern,
		"application_match_type":   me.ApplicationMatchType,
		"application_match_target": me.ApplicationMatchTarget,
	})
}

func (me *FilterConfig) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"pattern":                  &me.Pattern,
		"application_match_type":   &me.ApplicationMatchType,
		"application_match_target": &me.ApplicationMatchTarget,
	})
}
