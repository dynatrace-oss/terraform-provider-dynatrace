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

package logdpprules

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type RuleTesting struct {
	SampleLog string `json:"sampleLog"` // Sample log in JSON format.
}

func (me *RuleTesting) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"sample_log": {
			Type:        schema.TypeString,
			Description: "Sample log in JSON format.",
			Required:    true,
		},
	}
}

func (me *RuleTesting) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"sample_log": me.SampleLog,
	})
}

func (me *RuleTesting) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"sample_log": &me.SampleLog,
	})
}
