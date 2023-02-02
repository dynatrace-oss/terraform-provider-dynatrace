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

package infrastructurehosts

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type OutOfMemoryDetectionThresholds struct {
	EventThresholds             *StrictEventThresholds `json:"eventThresholds"`
	OutOfMemoryExceptionsNumber int                    `json:"outOfMemoryExceptionsNumber"` // Alert if the number of Java out-of-memory exceptions is at least this value
}

func (me *OutOfMemoryDetectionThresholds) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"event_thresholds": {
			Type:        schema.TypeList,
			Description: "no documentation available",
			Required:    true,
			Elem:        &schema.Resource{Schema: new(StrictEventThresholds).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
		"out_of_memory_exceptions_number": {
			Type:        schema.TypeInt,
			Description: "Alert if the number of Java out-of-memory exceptions is at least this value",
			Required:    true,
		},
	}
}

func (me *OutOfMemoryDetectionThresholds) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"event_thresholds":                me.EventThresholds,
		"out_of_memory_exceptions_number": me.OutOfMemoryExceptionsNumber,
	})
}

func (me *OutOfMemoryDetectionThresholds) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"event_thresholds":                &me.EventThresholds,
		"out_of_memory_exceptions_number": &me.OutOfMemoryExceptionsNumber,
	})
}
