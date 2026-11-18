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
	"fmt"
	"strconv"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type OutageHandling struct {
	GlobalConsecutiveOutageCountThreshold *int `json:"globalConsecutiveOutageCountThreshold,omitempty"` // Number of consecutive failures for all locations
	GlobalOutages                         bool `json:"globalOutages"`                                   // Generate a problem and send an alert when the monitor is unavailable at all configured locations
	LocalConsecutiveOutageCountThreshold  *int `json:"localConsecutiveOutageCountThreshold,omitempty"`  // Number of consecutive failures
	LocalLocationOutageCountThreshold     *int `json:"localLocationOutageCountThreshold,omitempty"`     // Number of failing locations
	LocalOutages                          bool `json:"localOutages"`                                    // Generate a problem and send an alert when the monitor is unavailable for one or more consecutive runs at any location
}

func (me *OutageHandling) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"global_consecutive_outage_count_threshold": {
			Type:             schema.TypeInt,
			Description:      "Number of consecutive failures for all locations",
			Optional:         true,
			DiffSuppressFunc: suppressDefaultGlobalConsecutiveOutageCountThreshold,
		},
		"global_outages": {
			Type:        schema.TypeBool,
			Description: "Generate a problem and send an alert when the monitor is unavailable at all configured locations",
			Optional:    true,
			Default:     true,
		},
		"local_consecutive_outage_count_threshold": {
			Type:        schema.TypeInt,
			Description: "Number of consecutive failures",
			Optional:    true,
		},
		"local_location_outage_count_threshold": {
			Type:        schema.TypeInt,
			Description: "Number of failing locations",
			Optional:    true,
		},
		"local_outages": {
			Type:        schema.TypeBool,
			Description: "Generate a problem and send an alert when the monitor is unavailable for one or more consecutive runs at any location",
			Optional:    true,
			Default:     false,
		},
	}
}

const defaultGlobalConsecutiveOutageCountThreshold = 1

// suppressDefaultGlobalConsecutiveOutageCountThreshold hides the diff between an unset threshold and the default sent by HandlePreconditions.
func suppressDefaultGlobalConsecutiveOutageCountThreshold(_, oldValue, newValue string, _ *schema.ResourceData) bool {
	return (newValue == "" || newValue == "0") && oldValue == strconv.Itoa(defaultGlobalConsecutiveOutageCountThreshold)
}

func (me *OutageHandling) HandlePreconditions() error {
	if me.GlobalOutages && me.GlobalConsecutiveOutageCountThreshold == nil {
		me.GlobalConsecutiveOutageCountThreshold = new(defaultGlobalConsecutiveOutageCountThreshold) // 0 is invalid
	}

	if !me.LocalOutages && me.LocalLocationOutageCountThreshold != nil {
		return fmt.Errorf("'local_location_outage_count_threshold' must not be specified unless 'local_outages' is set to 'true'; got 'local_outages'='%v'", me.LocalOutages)
	}

	if !me.LocalOutages && me.LocalConsecutiveOutageCountThreshold != nil {
		return fmt.Errorf("'local_consecutive_outage_count_threshold' must not be specified unless 'local_outages' is set to 'true'; got 'local_outages'='%v'", me.LocalOutages)
	}

	if !me.GlobalOutages && me.GlobalConsecutiveOutageCountThreshold != nil {
		return fmt.Errorf("'global_consecutive_outage_count_threshold' must not be specified unless 'global_outages' is set to 'true'; got 'global_outages'='%v'", me.GlobalOutages)
	}

	return nil
}

func (me *OutageHandling) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"global_consecutive_outage_count_threshold": me.GlobalConsecutiveOutageCountThreshold,
		"global_outages": me.GlobalOutages,
		"local_consecutive_outage_count_threshold": me.LocalConsecutiveOutageCountThreshold,
		"local_location_outage_count_threshold":    me.LocalLocationOutageCountThreshold,
		"local_outages":                            me.LocalOutages,
	})
}

func (me *OutageHandling) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"global_consecutive_outage_count_threshold": &me.GlobalConsecutiveOutageCountThreshold,
		"global_outages": &me.GlobalOutages,
		"local_consecutive_outage_count_threshold": &me.LocalConsecutiveOutageCountThreshold,
		"local_location_outage_count_threshold":    &me.LocalLocationOutageCountThreshold,
		"local_outages":                            &me.LocalOutages,
	})
}
