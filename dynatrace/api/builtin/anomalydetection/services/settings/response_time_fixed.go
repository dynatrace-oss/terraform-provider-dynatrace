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

package services

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type ResponseTimeFixed struct {
	OverAlertingProtection *OverAlertingProtection   `json:"overAlertingProtection"` // Avoid over-alerting
	ResponseTimeAll        *ResponseTimeFixedAll     `json:"responseTimeAll"`        // Alert if the median response time of all requests degrades beyond this threshold:
	ResponseTimeSlowest    *ResponseTimeFixedSlowest `json:"responseTimeSlowest"`    // Alert if the response time of the slowest 10% of requests degrades beyond this threshold:
	Sensitivity            Sensitivity               `json:"sensitivity"`            // Possible Values: `High`, `Low`, `Medium`
}

func (me *ResponseTimeFixed) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"over_alerting_protection": {
			Type:        schema.TypeList,
			Description: "Avoid over-alerting",
			Required:    true,
			Elem:        &schema.Resource{Schema: new(OverAlertingProtection).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
		"response_time_all": {
			Type:        schema.TypeList,
			Description: "Alert if the median response time of all requests degrades beyond this threshold:",
			Required:    true,
			Elem:        &schema.Resource{Schema: new(ResponseTimeFixedAll).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
		"response_time_slowest": {
			Type:        schema.TypeList,
			Description: "Alert if the response time of the slowest 10% of requests degrades beyond this threshold:",
			Required:    true,
			Elem:        &schema.Resource{Schema: new(ResponseTimeFixedSlowest).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
		"sensitivity": {
			Type:        schema.TypeString,
			Description: "Possible Values: `High`, `Low`, `Medium`",
			Required:    true,
		},
	}
}

func (me *ResponseTimeFixed) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"over_alerting_protection": me.OverAlertingProtection,
		"response_time_all":        me.ResponseTimeAll,
		"response_time_slowest":    me.ResponseTimeSlowest,
		"sensitivity":              me.Sensitivity,
	})
}

func (me *ResponseTimeFixed) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"over_alerting_protection": &me.OverAlertingProtection,
		"response_time_all":        &me.ResponseTimeAll,
		"response_time_slowest":    &me.ResponseTimeSlowest,
		"sensitivity":              &me.Sensitivity,
	})
}
