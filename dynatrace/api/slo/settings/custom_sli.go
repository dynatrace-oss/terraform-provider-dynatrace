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

package slo

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type CustomSli struct {
	Indicator      string         `json:"indicator" minlength:"1" maxlength:"2000"`
	FilterSegments FilterSegments `json:"filterSegments,omitempty"`
}

func (me *CustomSli) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"indicator": {
			Type:             schema.TypeString,
			Description:      "Indicator of the custom SLI. Example: `timeseries sli=avg(dt.host.cpu.idle)`",
			Required:         true,
			ValidateDiagFunc: Validate(ValidateMinLength(1), ValidateMaxLength(2000)),
		},
		"filter_segments": {
			Type:        schema.TypeList,
			Description: "A filter segment is identified by an ID. Each segment includes a list of variable definitions.",
			MinItems:    1,
			MaxItems:    1,
			Optional:    true,
			Elem:        &schema.Resource{Schema: new(FilterSegments).Schema()},
		},
	}
}

func (me *CustomSli) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"indicator":       me.Indicator,
		"filter_segments": me.FilterSegments,
	})
}

func (me *CustomSli) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"indicator":       &me.Indicator,
		"filter_segments": &me.FilterSegments,
	})
}
