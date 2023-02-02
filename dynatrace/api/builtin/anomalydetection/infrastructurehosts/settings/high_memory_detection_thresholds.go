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

package infrastructurehosts

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type HighMemoryDetectionThresholds struct {
	EventThresholds                *EventThresholds `json:"eventThresholds"`
	PageFaultsPerSecondNonWindows  int              `json:"pageFaultsPerSecondNonWindows"`  // Alert if the memory page fault rate on Unix systems is higher than this threshold for the defined amount of samples
	PageFaultsPerSecondWindows     int              `json:"pageFaultsPerSecondWindows"`     // Alert if the memory page fault rate on Windows is higher than this threshold for the defined amount of samples
	UsedMemoryPercentageNonWindows int              `json:"usedMemoryPercentageNonWindows"` // Alert if the memory usage on Unix systems is higher than this threshold
	UsedMemoryPercentageWindows    int              `json:"usedMemoryPercentageWindows"`    // Alert if the memory usage on Windows is higher than this threshold
}

func (me *HighMemoryDetectionThresholds) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"event_thresholds": {
			Type:        schema.TypeList,
			Description: "no documentation available",
			Required:    true,
			Elem:        &schema.Resource{Schema: new(EventThresholds).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
		"page_faults_per_second_non_windows": {
			Type:        schema.TypeInt,
			Description: "Alert if the memory page fault rate on Unix systems is higher than this threshold for the defined amount of samples",
			Required:    true,
		},
		"page_faults_per_second_windows": {
			Type:        schema.TypeInt,
			Description: "Alert if the memory page fault rate on Windows is higher than this threshold for the defined amount of samples",
			Required:    true,
		},
		"used_memory_percentage_non_windows": {
			Type:        schema.TypeInt,
			Description: "Alert if the memory usage on Unix systems is higher than this threshold",
			Required:    true,
		},
		"used_memory_percentage_windows": {
			Type:        schema.TypeInt,
			Description: "Alert if the memory usage on Windows is higher than this threshold",
			Required:    true,
		},
	}
}

func (me *HighMemoryDetectionThresholds) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"event_thresholds":                   me.EventThresholds,
		"page_faults_per_second_non_windows": me.PageFaultsPerSecondNonWindows,
		"page_faults_per_second_windows":     me.PageFaultsPerSecondWindows,
		"used_memory_percentage_non_windows": me.UsedMemoryPercentageNonWindows,
		"used_memory_percentage_windows":     me.UsedMemoryPercentageWindows,
	})
}

func (me *HighMemoryDetectionThresholds) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"event_thresholds":                   &me.EventThresholds,
		"page_faults_per_second_non_windows": &me.PageFaultsPerSecondNonWindows,
		"page_faults_per_second_windows":     &me.PageFaultsPerSecondWindows,
		"used_memory_percentage_non_windows": &me.UsedMemoryPercentageNonWindows,
		"used_memory_percentage_windows":     &me.UsedMemoryPercentageWindows,
	})
}
