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

package anomalydetectors

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type DavisEventTemplate struct {
	Properties EventProperties `json:"properties,omitempty"` // Set of additional key-value properties to be attached to the triggered event.
}

func (me *DavisEventTemplate) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"properties": {
			Type:        schema.TypeList,
			Description: "Set of additional key-value properties to be attached to the triggered event.",
			Optional:    true, // minobjects == 0
			Elem:        &schema.Resource{Schema: new(EventProperties).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
	}
}

func (me *DavisEventTemplate) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"properties": me.Properties,
	})
}

func (me *DavisEventTemplate) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"properties": &me.Properties,
	})
}
