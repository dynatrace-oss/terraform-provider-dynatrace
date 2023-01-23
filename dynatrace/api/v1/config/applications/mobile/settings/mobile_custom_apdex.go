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

package mobile

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// MobileCustomApdex represents Apdex configuration of a mobile or custom application. \n\nA duration less than the **tolerable** threshold is considered satisfied
type MobileCustomApdex struct {
	ToleratedThreshold   int32 `json:"toleratedThreshold"`   // Apdex **tolerable** threshold, in milliseconds: a duration greater than or equal to this value is considered tolerable
	FrustratingThreshold int32 `json:"frustratingThreshold"` // Apdex **frustrated** threshold, in milliseconds: a duration greater than or equal to this value is considered frustrated
	FrustratedOnError    bool  `json:"frustratedOnError"`    // Apdex error condition: if `true` the user session is considered frustrated when an error is reported
}

func (me *MobileCustomApdex) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"tolerable": {
			Type:        schema.TypeInt,
			Description: "Apdex **tolerable** threshold, in milliseconds: a duration greater than or equal to this value is considered tolerable",
			Required:    true,
		},
		"frustrated": {
			Type:        schema.TypeInt,
			Description: "Apdex **frustrated** threshold, in milliseconds: a duration greater than or equal to this value is considered frustrated",
			Required:    true,
		},
		"frustrated_on_error": {
			Type:        schema.TypeBool,
			Description: "Apdex error condition: if `true` the user session is considered frustrated when an error is reported",
			Optional:    true,
		},
	}
}

func (me *MobileCustomApdex) UnmarshalHCL(decoder hcl.Decoder) error {
	if err := decoder.Decode("tolerable", &me.ToleratedThreshold); err != nil {
		return err
	}
	if err := decoder.Decode("frustrated", &me.FrustratingThreshold); err != nil {
		return err
	}
	if err := decoder.Decode("frustrated_on_error", &me.FrustratedOnError); err != nil {
		return err
	}
	return nil
}

func (me *MobileCustomApdex) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"tolerable":           me.ToleratedThreshold,
		"frustrated":          me.FrustratingThreshold,
		"frustrated_on_error": me.FrustratedOnError,
	})
}
