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
	"encoding/json"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// OutageHandlingPolicy Outage handling configuration
type OutageHandlingPolicy struct {
	GlobalOutage       bool                `json:"globalOutage"`       // When enabled (`true`), generate a problem and send an alert when the monitor is unavailable at all configured locations
	GlobalOutagePolicy *GlobalOutagePolicy `json:"globalOutagePolicy"` // Global outage handling configuration. \n\n Alert if **consecutiveRuns** times consecutively
	LocalOutage        bool                `json:"localOutage"`        // When enabled (`true`), generate a problem and send an alert when the monitor is unavailable for one or more consecutive runs at any location
	LocalOutagePolicy  *LocalOutagePolicy  `json:"localOutagePolicy"`  // Local outage handling configuration. \n\n Alert if **affectedLocations** of locations are unable to access the web application **consecutiveRuns** times consecutively
	RetryOnError       bool                `json:"retryOnError"`       // Schedule retry if browser monitor execution results in a fail. For HTTP monitors this property is ignored
}

func (me *OutageHandlingPolicy) MarshalJSON() ([]byte, error) {
	lop := me.LocalOutagePolicy
	if lop == nil {
		lop = new(LocalOutagePolicy)
	}
	gop := me.GlobalOutagePolicy
	if gop == nil {
		gop = new(GlobalOutagePolicy)
	}
	return json.Marshal(struct {
		GlobalOutage       bool                `json:"globalOutage"`
		GlobalOutagePolicy *GlobalOutagePolicy `json:"globalOutagePolicy"`
		LocalOutage        bool                `json:"localOutage"`
		LocalOutagePolicy  *LocalOutagePolicy  `json:"localOutagePolicy"`
		RetryOnError       bool                `json:"retryOnError"`
	}{
		me.GlobalOutage,
		gop,
		me.LocalOutage,
		lop,
		me.RetryOnError,
	})
}

func (me *OutageHandlingPolicy) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"global_outage": {
			Type:        schema.TypeBool,
			Description: "When enabled (`true`), generate a problem and send an alert when the monitor is unavailable at all configured locations",
			Optional:    true,
		},
		"local_outage": {
			Type:        schema.TypeBool,
			Description: "When enabled (`true`), generate a problem and send an alert when the monitor is unavailable for one or more consecutive runs at any location",
			Optional:    true,
		},
		"retry_on_error": {
			Type:        schema.TypeBool,
			Description: "Schedule retry if browser monitor execution results in a fail. For HTTP monitors this property is ignored",
			Optional:    true,
		},
		"local_outage_policy": {
			Type:        schema.TypeList,
			Description: "Local outage handling configuration. \n\n Alert if **affectedLocations** of locations are unable to access the web application **consecutiveRuns** times consecutively",
			Optional:    true,
			Elem:        &schema.Resource{Schema: new(LocalOutagePolicy).Schema()},
		},
		"global_outage_policy": {
			Type:        schema.TypeList,
			Description: "Global outage handling configuration. \n\n Alert if **consecutiveRuns** times consecutively",
			Optional:    true,
			Elem:        &schema.Resource{Schema: new(GlobalOutagePolicy).Schema()},
		},
	}
}

func (me *OutageHandlingPolicy) MarshalHCL(properties hcl.Properties) error {
	if err := properties.Encode("global_outage", me.GlobalOutage); err != nil {
		return err
	}
	if err := properties.Encode("local_outage", me.LocalOutage); err != nil {
		return err
	}
	if err := properties.Encode("retry_on_error", me.RetryOnError); err != nil {
		return err
	}
	if err := properties.Encode("local_outage_policy", me.LocalOutagePolicy); err != nil {
		return err
	}
	if err := properties.Encode("global_outage_policy", me.GlobalOutagePolicy); err != nil {
		return err
	}
	return nil
}

func (me *OutageHandlingPolicy) UnmarshalHCL(decoder hcl.Decoder) error {
	if err := decoder.Decode("global_outage", &me.GlobalOutage); err != nil {
		return err
	}
	if err := decoder.Decode("local_outage", &me.LocalOutage); err != nil {
		return err
	}
	if err := decoder.Decode("retry_on_error", &me.RetryOnError); err != nil {
		return err
	}
	if err := decoder.Decode("local_outage_policy", &me.LocalOutagePolicy); err != nil {
		return err
	}
	if err := decoder.Decode("global_outage_policy", &me.GlobalOutagePolicy); err != nil {
		return err
	}
	return nil
}
