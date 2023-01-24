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

package monitors

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// GlobalOutagePolicy Global outage handling configuration. \n\n Alert if **consecutiveRuns** times consecutively
type GlobalOutagePolicy struct {
	ConsecutiveRuns *int32 `json:"consecutiveRuns"` // The number of consecutive fails to trigger an alert
}

func (me *GlobalOutagePolicy) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"consecutive_runs": {
			Type:        schema.TypeInt,
			Description: "The number of consecutive fails to trigger an alert",
			Required:    true,
		},
	}
}

func (me *GlobalOutagePolicy) MarshalHCL(properties hcl.Properties) error {
	if me.ConsecutiveRuns == nil {
		return nil
	}
	if err := properties.Encode("consecutive_runs", me.ConsecutiveRuns); err != nil {
		return err
	}
	return nil
}

func (me *GlobalOutagePolicy) UnmarshalHCL(decoder hcl.Decoder) error {
	if err := decoder.Decode("consecutive_runs", &me.ConsecutiveRuns); err != nil {
		return err
	}
	return nil
}
