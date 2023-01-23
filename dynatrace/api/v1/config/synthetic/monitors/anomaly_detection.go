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

// AnomalyDetection The anomaly detection configuration
type AnomalyDetection struct {
	OutageHandling        *OutageHandlingPolicy        `json:"outageHandling,omitempty"`
	LoadingTimeThresholds *LoadingTimeThresholdsPolicy `json:"loadingTimeThresholds,omitempty"`
}

func (me *AnomalyDetection) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"outage_handling": {
			Type:        schema.TypeList,
			Description: "Outage handling configuration",
			Optional:    true,
			Elem:        &schema.Resource{Schema: new(OutageHandlingPolicy).Schema()},
		},
		"loading_time_thresholds": {
			Type:        schema.TypeList,
			Description: "Thresholds for loading times",
			Optional:    true,
			Elem: &schema.Resource{
				Schema: new(LoadingTimeThresholdsPolicy).Schema(),
			},
		},
	}
}

func (me *AnomalyDetection) MarshalHCL(properties hcl.Properties) error {
	if err := properties.Encode("outage_handling", me.OutageHandling); err != nil {
		return err
	}
	if err := properties.Encode("loading_time_thresholds", me.LoadingTimeThresholds); err != nil {
		return err
	}
	return nil
}

func (me *AnomalyDetection) UnmarshalHCL(decoder hcl.Decoder) error {
	if err := decoder.Decode("outage_handling", &me.OutageHandling); err != nil {
		return err
	}
	if err := decoder.Decode("loading_time_thresholds", &me.LoadingTimeThresholds); err != nil {
		return err
	}
	return nil
}
