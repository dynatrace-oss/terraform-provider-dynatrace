/*
 * @license
 * Copyright 2026 Dynatrace LLC
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package http

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type CustomProperties []*CustomProperty

func (me *CustomProperties) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"custom_property": {
			Type:        schema.TypeSet,
			Description: "Custom properties for the monitor",
			Required:    true,
			MinItems:    1,
			Elem:        &schema.Resource{Schema: new(CustomProperty).Schema()},
		},
	}
}

func (me CustomProperties) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeSlice("custom_property", me)
}

func (me *CustomProperties) UnmarshalHCL(decoder hcl.Decoder) error {
	if err := decoder.DecodeSlice("custom_property", me); err != nil {
		return err
	}

	*me = hcl.FilterEmpty(*me, CustomProperty{})
	return nil
}

type CustomProperty struct {
	Name  string `json:"name"`  // The name of the custom property
	Value string `json:"value"` // The value of the custom property
}

func (me *CustomProperty) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:        schema.TypeString,
			Description: "The name of the custom property. Possible values: `hmRequestTimeoutInMs`, `hmConnectTimeoutInMs`, `hmMaxHeaderSizeInBytes`, `hmMonitorExecutionTimeoutInMs`, `hmScriptExecutionTimeoutInMs`, `hmMaxRequestBodySizeInBytes`, `hmMaxCustomScriptSizeInBytes`, `hmMaxResponseBodySizeInBytes`, `hmMaxResponseBodySizeToCustomScriptInBytes`, `hmDnsQueryTimeoutInMs`",
			Required:    true,
		},
		"value": {
			Type:        schema.TypeString,
			Description: "The value of the custom property",
			Required:    true,
		},
	}
}

func (me *CustomProperty) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"name":  me.Name,
		"value": me.Value,
	})
}

func (me *CustomProperty) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"name":  &me.Name,
		"value": &me.Value,
	})
}
