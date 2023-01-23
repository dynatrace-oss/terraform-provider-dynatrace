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

// LocalOutagePolicy Local outage handling configuration. \n\n Alert if **affectedLocations** of locations are unable to access the web application **consecutiveRuns** times consecutively
type LocalOutagePolicy struct {
	AffectedLocations *int32 `json:"affectedLocations"` // The number of affected locations to trigger an alert
	ConsecutiveRuns   *int32 `json:"consecutiveRuns"`   // The number of consecutive fails to trigger an alert
}

func (me *LocalOutagePolicy) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"affected_locations": {
			Type:        schema.TypeInt,
			Description: "The number of affected locations to trigger an alert",
			Required:    true,
		},
		"consecutive_runs": {
			Type:        schema.TypeInt,
			Description: "The number of consecutive fails to trigger an alert",
			Required:    true,
		},
	}
}

func (me *LocalOutagePolicy) MarshalHCL(properties hcl.Properties) error {
	if me.AffectedLocations == nil && me.ConsecutiveRuns == nil {
		return nil
	}
	if err := properties.Encode("affected_locations", me.AffectedLocations); err != nil {
		return err
	}
	if err := properties.Encode("consecutive_runs", me.ConsecutiveRuns); err != nil {
		return err
	}
	return nil
}

func (me *LocalOutagePolicy) UnmarshalHCL(decoder hcl.Decoder) error {
	if err := decoder.Decode("affected_locations", &me.AffectedLocations); err != nil {
		return err
	}
	if err := decoder.Decode("consecutive_runs", &me.ConsecutiveRuns); err != nil {
		return err
	}
	return nil
}
