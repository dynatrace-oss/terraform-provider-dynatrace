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

package web

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type MetaDataCaptureSettings []*MetaDataCapturing

func (me *MetaDataCaptureSettings) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"capture": {
			Type:        schema.TypeList,
			Description: "Java script agent meta data capture settings",
			Optional:    true,
			Elem:        &schema.Resource{Schema: new(MetaDataCapturing).Schema()},
		},
	}
}

func (me MetaDataCaptureSettings) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeSlice("capture", me)
}

func (me *MetaDataCaptureSettings) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeSlice("capture", me)
}

type MetaDataCapturing struct {
	Type           MetaDataCapturingType `json:"type"`               // The type of the meta data to capture. Possible values are `COOKIE`, `CSS_SELECTOR`, `JAVA_SCRIPT_FUNCTION`, `JAVA_SCRIPT_VARIABLE`, `META_TAG` and `QUERY_STRING`.
	CapturingName  string                `json:"capturingName"`      // The name of the meta data to capture
	Name           string                `json:"name"`               // Name for displaying the captured values in Dynatrace
	UniqueID       *int32                `json:"uniqueId,omitempty"` // The unique ID of the meta data to capture
	PublicMetadata bool                  `json:"publicMetadata"`     // `true` if this metadata should be captured regardless of the privacy settings, `false` otherwise
	UseLastValue   bool                  `json:"useLastValue"`       // `true` if the last captured value should be used for this metadata. By default the first value will be used.
}

func (me *MetaDataCapturing) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"type": {
			Type:        schema.TypeString,
			Description: "The type of the meta data to capture. Possible values are `COOKIE`, `CSS_SELECTOR`, `JAVA_SCRIPT_FUNCTION`, `JAVA_SCRIPT_VARIABLE`, `META_TAG` and `QUERY_STRING`.",
			Required:    true,
		},
		"capturing_name": {
			Type:        schema.TypeString,
			Description: "The name of the meta data to capture",
			Required:    true,
		},
		"name": {
			Type:        schema.TypeString,
			Description: "Name for displaying the captured values in Dynatrace",
			Required:    true,
		},
		"unique_id": {
			Type:        schema.TypeInt,
			Description: "The unique ID of the meta data to capture",
			Optional:    true,
		},
		"public_metadata": {
			Type:        schema.TypeBool,
			Description: "`true` if this metadata should be captured regardless of the privacy settings, `false` otherwise",
			Optional:    true,
		},
		"use_last_value": {
			Type:        schema.TypeBool,
			Description: "`true` if the last captured value should be used for this metadata. By default the first value will be used.",
			Optional:    true,
		},
	}
}

func (me *MetaDataCapturing) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"type":            me.Type,
		"capturing_name":  me.CapturingName,
		"name":            me.Name,
		"unique_id":       me.UniqueID,
		"public_metadata": me.PublicMetadata,
		"use_last_value":  me.UseLastValue,
	})
}

func (me *MetaDataCapturing) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"type":            &me.Type,
		"capturing_name":  &me.CapturingName,
		"name":            &me.Name,
		"unique_id":       &me.UniqueID,
		"public_metadata": &me.PublicMetadata,
		"use_last_value":  &me.UseLastValue,
	})
}
