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

package memory

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Thresholds Custom thresholds for high memory usage. If not set then the automatic mode is used.
//
//	**Both** conditions must be met to trigger an alert.
type Thresholds struct {
	UsedMemoryPercentageWindows    int32 `json:"usedMemoryPercentageWindows"`    // Memory usage is higher than *X*% on Windows.
	UsedMemoryPercentageNonWindows int32 `json:"usedMemoryPercentageNonWindows"` // Memory usage is higher than *X*% on Linux.
	PageFaultsPerSecondNonWindows  int32 `json:"pageFaultsPerSecondNonWindows"`  // Memory page fault rate is higher than *X* faults per second on Linux.
	PageFaultsPerSecondWindows     int32 `json:"pageFaultsPerSecondWindows"`     // Memory page fault rate is higher than *X* faults per second on Windows.
}

type osThresholds struct {
	UsedMemoryPercentage int32 // Memory usage is higher than *X*%.
	PageFaultsPerSecond  int32 // Memory page fault rate is higher than *X* faults per second.
}

func (me *osThresholds) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"usage": {
			Type:        schema.TypeInt,
			Required:    true,
			Description: "Memory usage is higher than *X*%",
		},
		"page_faults": {
			Type:        schema.TypeInt,
			Required:    true,
			Description: "Memory page fault rate is higher than *X* faults per second",
		},
	}
}

func (me *osThresholds) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"usage":       me.UsedMemoryPercentage,
		"page_faults": me.PageFaultsPerSecond,
	})
}

func (me *osThresholds) UnmarshalHCL(decoder hcl.Decoder) error {
	if value, ok := decoder.GetOk("usage"); ok {
		me.UsedMemoryPercentage = int32(value.(int))
	}
	if value, ok := decoder.GetOk("page_faults"); ok {
		me.PageFaultsPerSecond = int32(value.(int))
	}
	return nil
}

func (me *Thresholds) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"windows": {
			Type:        schema.TypeList,
			Required:    true,
			MaxItems:    1,
			MinItems:    1,
			Description: "Custom thresholds for Windows",
			Elem:        &schema.Resource{Schema: new(osThresholds).Schema()},
		},
		"linux": {
			Type:        schema.TypeList,
			Required:    true,
			MaxItems:    1,
			MinItems:    1,
			Description: "Custom thresholds for Linux",
			Elem:        &schema.Resource{Schema: new(osThresholds).Schema()},
		},
	}
}

func (me *Thresholds) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"linux":   &osThresholds{UsedMemoryPercentage: me.UsedMemoryPercentageNonWindows, PageFaultsPerSecond: me.PageFaultsPerSecondNonWindows},
		"windows": &osThresholds{UsedMemoryPercentage: me.UsedMemoryPercentageWindows, PageFaultsPerSecond: me.PageFaultsPerSecondWindows},
	})
}

func (me *Thresholds) UnmarshalHCL(decoder hcl.Decoder) error {
	if _, ok := decoder.GetOk("windows.#"); ok {
		cfg := new(osThresholds)
		if err := cfg.UnmarshalHCL(hcl.NewDecoder(decoder, "windows", 0)); err != nil {
			return err
		}
		me.PageFaultsPerSecondWindows = cfg.PageFaultsPerSecond
		me.UsedMemoryPercentageWindows = cfg.UsedMemoryPercentage
	}
	if _, ok := decoder.GetOk("linux.#"); ok {
		cfg := new(osThresholds)
		if err := cfg.UnmarshalHCL(hcl.NewDecoder(decoder, "linux", 0)); err != nil {
			return err
		}
		me.PageFaultsPerSecondNonWindows = cfg.PageFaultsPerSecond
		me.UsedMemoryPercentageNonWindows = cfg.UsedMemoryPercentage
	}
	return nil
}
