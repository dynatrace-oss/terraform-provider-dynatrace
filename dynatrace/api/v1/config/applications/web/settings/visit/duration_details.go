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

package visit

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// DurationDetails Configuration of a visit duration-based conversion goal
type DurationDetails struct {
	DurationInMillis int64 `json:"durationInMillis"` // The duration of session to hit the conversion goal, in milliseconds
}

func (me *DurationDetails) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"duration": {
			Type:        schema.TypeInt,
			Description: "The duration of session to hit the conversion goal, in milliseconds",
			Required:    true,
		},
	}
}

func (me *DurationDetails) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"duration": me.DurationInMillis,
	})
}

func (me *DurationDetails) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"duration": &me.DurationInMillis,
	})
}
