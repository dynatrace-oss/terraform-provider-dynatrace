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

package rumweb

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Settings struct {
	ErrorRate     *ErrorRate        `json:"errorRate"`       // Error rate
	ResponseTime  *ResponseTime     `json:"responseTime"`    // Response time
	Scope         *string           `json:"-" scope:"scope"` // The scope of this setting (APPLICATION_METHOD, APPLICATION). Omit this property if you want to cover the whole environment.
	TrafficDrops  *AppTrafficDrops  `json:"trafficDrops"`    // Detect traffic drops
	TrafficSpikes *AppTrafficSpikes `json:"trafficSpikes"`   // Detect traffic spikes
}

func (me *Settings) Name() string {
	return *me.Scope
}

func (me *Settings) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"error_rate": {
			Type:        schema.TypeList,
			Description: "Error rate",
			Required:    true,

			Elem:     &schema.Resource{Schema: new(ErrorRate).Schema()},
			MinItems: 1,
			MaxItems: 1,
		},
		"response_time": {
			Type:        schema.TypeList,
			Description: "Response time",
			Required:    true,

			Elem:     &schema.Resource{Schema: new(ResponseTime).Schema()},
			MinItems: 1,
			MaxItems: 1,
		},
		"scope": {
			Type:        schema.TypeString,
			Description: "The scope of this setting (APPLICATION_METHOD, APPLICATION). Omit this property if you want to cover the whole environment.",
			Optional:    true,
			Default:     "environment",
			ForceNew:    true,
		},
		"traffic_drops": {
			Type:        schema.TypeList,
			Description: "Detect traffic drops",
			Required:    true,

			Elem:     &schema.Resource{Schema: new(AppTrafficDrops).Schema()},
			MinItems: 1,
			MaxItems: 1,
		},
		"traffic_spikes": {
			Type:        schema.TypeList,
			Description: "Detect traffic spikes",
			Required:    true,

			Elem:     &schema.Resource{Schema: new(AppTrafficSpikes).Schema()},
			MinItems: 1,
			MaxItems: 1,
		},
	}
}

func (me *Settings) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"error_rate":     me.ErrorRate,
		"response_time":  me.ResponseTime,
		"scope":          me.Scope,
		"traffic_drops":  me.TrafficDrops,
		"traffic_spikes": me.TrafficSpikes,
	})
}

func (me *Settings) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"error_rate":     &me.ErrorRate,
		"response_time":  &me.ResponseTime,
		"scope":          &me.Scope,
		"traffic_drops":  &me.TrafficDrops,
		"traffic_spikes": &me.TrafficSpikes,
	})
}
