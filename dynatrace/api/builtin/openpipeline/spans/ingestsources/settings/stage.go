/**
* @license
* Copyright 2025 Dynatrace LLC
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

package ingestsources

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Stage struct {
	Processors Processors `json:"processors,omitempty"` // Processors of stage
}

func (me *Stage) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"processors": {
			Type:        schema.TypeList,
			Description: "Processors of stage",
			Optional:    true, // minobjects == 0
			Elem:        &schema.Resource{Schema: new(Processors).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
	}
}

func (me *Stage) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"processors": me.Processors,
	})
}

func (me *Stage) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"processors": &me.Processors,
	})
}
