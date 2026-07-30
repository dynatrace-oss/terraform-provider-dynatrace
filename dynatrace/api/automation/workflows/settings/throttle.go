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

package workflows

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Throttle struct {
	IsLimitHit bool `json:"isLimitHit"` // Whether the execution limit is currently hit. Set `false` to clear an active throttle; the API rejects `true` unless the workflow is already throttled
}

func (me *Throttle) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"is_limit_hit": {
			Type: schema.TypeBool,
			Description: "Whether the workflow's execution limit is currently hit. This value is computed by the server from the current throttle state." +
				"\n  Set to `false` to reset (clear) an active throttle." +
				"\n  The API rejects `true` unless the workflow is already throttled." +
				"\n  When omitted, the value is read from the API",
			Optional: true,
			Computed: true,
		},
	}
}

func (me *Throttle) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"is_limit_hit": me.IsLimitHit,
	})
}

func (me *Throttle) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"is_limit_hit": &me.IsLimitHit,
	})
}
