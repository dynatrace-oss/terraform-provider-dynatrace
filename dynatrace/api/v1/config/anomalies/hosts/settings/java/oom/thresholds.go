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

package oom

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Thresholds Custom thresholds for Java out of memory. If not set, automatic mode is used.
type Thresholds struct {
	ExceptionCount int32 `json:"outOfMemoryExceptionsNumber"` // Alert if the number of Java out of memory exceptions is *X* per minute or higher.
}

func (me *Thresholds) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"exception_count": {
			Type:        schema.TypeInt,
			Required:    true,
			Description: "Alert if the number of Java out of memory exceptions is *X* per minute or higher",
		},
	}
}

func (me *Thresholds) MarshalHCL(properties hcl.Properties) error {
	return properties.Encode("exception_count", me.ExceptionCount)
}

func (me *Thresholds) UnmarshalHCL(decoder hcl.Decoder) error {
	if value, ok := decoder.GetOk("exception_count"); ok {
		me.ExceptionCount = int32(value.(int))
	}
	return nil
}
