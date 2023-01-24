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

package processgroups

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// AnomalyDetection Configuration of anomaly detection for the process group.
type AnomalyDetection struct {
	ProcessGroupId         string                  `json:"-"`                                // The ID of the process group
	AvailabilityMonitoring *AvailabilityMonitoring `json:"availabilityMonitoring,omitempty"` // Configuration of the availability monitoring for the process group.
}

func (me *AnomalyDetection) Name() string {
	return me.ProcessGroupId + "-anomalydetection"
}

func (me *AnomalyDetection) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"pg_id": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The ID of the process group",
		},
		"availability": {
			Type:        schema.TypeList,
			Optional:    true,
			MaxItems:    1,
			Description: "Configuration of the availability monitoring for the process group.",
			Elem:        &schema.Resource{Schema: new(AvailabilityMonitoring).Schema()},
		},
	}
}

func (me *AnomalyDetection) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"pg_id":        me.ProcessGroupId,
		"availability": me.AvailabilityMonitoring,
	})
}

func (me *AnomalyDetection) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"pg_id":        &me.ProcessGroupId,
		"availability": &me.AvailabilityMonitoring,
	})
}
