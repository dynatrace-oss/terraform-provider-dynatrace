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

package urlbasedsampling

import (
	"fmt"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Settings struct {
	Enabled            bool                 `json:"enabled"`                      // This setting is enabled (`true`) or disabled (`false`)
	Factor             *SamplingScaleFactor `json:"factor,omitempty"`             // Select the scaling factor for the current sampling rate of the system. Possible values: `IncreaseCapturing128Times`, `IncreaseCapturing64Times`, `IncreaseCapturing32Times`, `IncreaseCapturing16Times`, `IncreaseCapturing8Times`, `IncreaseCapturing4Times`, `IncreaseCapturing2Times`, `ReduceCapturingByFactor2`, `ReduceCapturingByFactor4`, `ReduceCapturingByFactor8`, `ReduceCapturingByFactor16`, `ReduceCapturingByFactor32`, `ReduceCapturingByFactor64`, `ReduceCapturingByFactor128`
	HttpMethod         []HttpMethod         `json:"httpMethod,omitempty"`         // Possible values: `GET`, `POST`, `PUT`, `DELETE`, `HEAD`, `CONNECT`, `OPTIONS`, `TRACE`, `PATCH`
	HttpMethodAny      bool                 `json:"httpMethodAny"`                // The scaling factor for the defined URL will be applied to any HTTP method.
	Ignore             bool                 `json:"ignore"`                       // The matching URLs will always be ignored, also if Adaptive Traffic Management is not active.
	Path               *string              `json:"path,omitempty"`               // Path of the URL.
	PathComparisonType *PathComparisonType  `json:"pathComparisonType,omitempty"` // Path comparison condition. Possible values: `EQUALS`, `DOES_NOT_EQUAL`, `CONTAINS`, `DOES_NOT_CONTAIN`, `STARTS_WITH`, `DOES_NOT_START_WITH`, `ENDS_WITH`, `DOES_NOT_END_WITH`
	QueryParameters    QueryParameters      `json:"queryParameters"`              // Add URL parameters in any order. **All** specified parameters must be present in the query of an URL to get a match.
	Scope              *string              `json:"-" scope:"scope"`              // The scope of this setting (PROCESS_GROUP_INSTANCE, PROCESS_GROUP). Omit this property if you want to cover the whole environment.
	InsertAfter        string               `json:"-"`
}

func (me *Settings) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"enabled": {
			Type:        schema.TypeBool,
			Description: "This setting is enabled (`true`) or disabled (`false`)",
			Required:    true,
		},
		"factor": {
			Type:        schema.TypeString,
			Description: "Select the scaling factor for the current sampling rate of the system. Possible values: `IncreaseCapturing128Times`, `IncreaseCapturing64Times`, `IncreaseCapturing32Times`, `IncreaseCapturing16Times`, `IncreaseCapturing8Times`, `IncreaseCapturing4Times`, `IncreaseCapturing2Times`, `ReduceCapturingByFactor2`, `ReduceCapturingByFactor4`, `ReduceCapturingByFactor8`, `ReduceCapturingByFactor16`, `ReduceCapturingByFactor32`, `ReduceCapturingByFactor64`, `ReduceCapturingByFactor128`",
			Optional:    true,
		},
		"http_method": {
			Type:        schema.TypeSet,
			Description: "Possible values: `GET`, `POST`, `PUT`, `DELETE`, `HEAD`, `CONNECT`, `OPTIONS`, `TRACE`, `PATCH`",
			Optional:    true,
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
		"http_method_any": {
			Type:        schema.TypeBool,
			Description: "The scaling factor for the defined URL will be applied to any HTTP method.",
			Required:    true,
		},
		"ignore": {
			Type:        schema.TypeBool,
			Description: "The matching URLs will always be ignored, also if Adaptive Traffic Management is not active.",
			Required:    true,
		},
		"path": {
			Type:        schema.TypeString,
			Description: "Path of the URL.",
			Optional:    true,
		},
		"path_comparison_type": {
			Type:        schema.TypeString,
			Description: "Path comparison condition. Possible values: `EQUALS`, `DOES_NOT_EQUAL`, `CONTAINS`, `DOES_NOT_CONTAIN`, `STARTS_WITH`, `DOES_NOT_START_WITH`, `ENDS_WITH`, `DOES_NOT_END_WITH`",
			Optional:    true,
		},
		"query_parameters": {
			Type:        schema.TypeList,
			Description: "Add URL parameters in any order. **All** specified parameters must be present in the query of an URL to get a match.",
			Required:    true,
			Elem:        &schema.Resource{Schema: new(QueryParameters).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
		"scope": {
			Type:        schema.TypeString,
			Description: "The scope of this setting (PROCESS_GROUP_INSTANCE, PROCESS_GROUP). Omit this property if you want to cover the whole environment.",
			Optional:    true,
			Default:     "environment",
			ForceNew:    true,
		},
		"insert_after": {
			Type:        schema.TypeString,
			Description: "Because this resource allows for ordering you may specify the ID of the resource instance that comes before this instance regarding order. If not specified when creating the setting will be added to the end of the list. If not specified during update the order will remain untouched",
			Optional:    true,
			Computed:    true,
		},
	}
}

func (me *Settings) Name() string {
	if me.Scope == nil {
		return "environment"
	}
	return *me.Scope
}

func (me *Settings) MarshalHCL(properties hcl.Properties) error {
	err := properties.EncodeAll(map[string]any{
		"enabled":              me.Enabled,
		"http_method":          me.HttpMethod,
		"http_method_any":      me.HttpMethodAny,
		"ignore":               me.Ignore,
		"path":                 me.Path,
		"path_comparison_type": me.PathComparisonType,
		"query_parameters":     me.QueryParameters,
		"scope":                me.Scope,
		"insert_after":         me.InsertAfter,
	})
	if me.Factor != nil {
		for name, id := range SamplingScaleFactorLookup {
			if *me.Factor == id {
				if err := properties.Encode("factor", name); err != nil {
					return err
				}
			}
		}
	}
	return err
}

func (me *Settings) UnmarshalHCL(decoder hcl.Decoder) error {
	err := decoder.DecodeAll(map[string]any{
		"enabled":              &me.Enabled,
		"http_method":          &me.HttpMethod,
		"http_method_any":      &me.HttpMethodAny,
		"ignore":               &me.Ignore,
		"path":                 &me.Path,
		"path_comparison_type": &me.PathComparisonType,
		"query_parameters":     &me.QueryParameters,
		"scope":                &me.Scope,
		"insert_after":         &me.InsertAfter,
	})
	if factor, ok := decoder.GetOk("factor"); ok {
		if value, found := SamplingScaleFactorLookup[factor.(string)]; found {
			me.Factor = &value
		} else {
			return fmt.Errorf("invalid `factor` field, possible values: `IncreaseCapturing128Times`, `IncreaseCapturing64Times`, `IncreaseCapturing32Times`, `IncreaseCapturing16Times`, `IncreaseCapturing8Times`, `IncreaseCapturing4Times`, `IncreaseCapturing2Times`, `ReduceCapturingByFactor2`, `ReduceCapturingByFactor4`, `ReduceCapturingByFactor8`, `ReduceCapturingByFactor16`, `ReduceCapturingByFactor32`, `ReduceCapturingByFactor64`, `ReduceCapturingByFactor128`")
		}
	}
	return err
}
