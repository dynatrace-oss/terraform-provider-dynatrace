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

package diskrules

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type SampleLimit struct {
	Samples          int `json:"samples"`          // .. within the last
	ViolatingSamples int `json:"violatingSamples"` // Minimum number of violating samples
}

func (me *SampleLimit) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"samples": {
			Type:        schema.TypeInt,
			Description: ".. within the last",
			Required:    true,
		},
		"violating_samples": {
			Type:        schema.TypeInt,
			Description: "Minimum number of violating samples",
			Required:    true,
		},
	}
}

func (me *SampleLimit) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"samples":           me.Samples,
		"violating_samples": me.ViolatingSamples,
	})
}

func (me *SampleLimit) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"samples":           &me.Samples,
		"violating_samples": &me.ViolatingSamples,
	})
}
