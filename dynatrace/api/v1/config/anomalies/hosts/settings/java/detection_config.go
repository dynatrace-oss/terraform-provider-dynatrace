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

package java

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/anomalies/hosts/settings/java/oom"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/anomalies/hosts/settings/java/oot"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type DetectionConfig struct {
	OutOfMemoryDetection  *oom.DetectionConfig `json:"outOfMemoryDetection"`  // Configuration of Java out of memory problems detection.
	OutOfThreadsDetection *oot.DetectionConfig `json:"outOfThreadsDetection"` // Configuration of Java out of threads problems detection.
}

func (me *DetectionConfig) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"out_of_threads": {
			Type:        schema.TypeList,
			Optional:    true,
			MaxItems:    1,
			Description: "Configuration of Java out of threads problems detection",
			Elem:        &schema.Resource{Schema: new(oot.DetectionConfig).Schema()},
		},
		"out_of_memory": {
			Type:        schema.TypeList,
			Optional:    true,
			MaxItems:    1,
			Description: "Configuration of Java out of memory problems detection",
			Elem:        &schema.Resource{Schema: new(oom.DetectionConfig).Schema()},
		},
	}
}

func (me *DetectionConfig) IsConfigured() bool {
	if me.OutOfMemoryDetection != nil && me.OutOfMemoryDetection.Enabled {
		return true
	}
	if me.OutOfThreadsDetection != nil && me.OutOfThreadsDetection.Enabled {
		return true
	}
	return false
}

func (me *DetectionConfig) MarshalHCL(properties hcl.Properties) error {
	if !me.IsConfigured() {
		return nil
	}
	return properties.EncodeAll(map[string]any{
		"out_of_memory":  me.OutOfMemoryDetection,
		"out_of_threads": me.OutOfThreadsDetection,
	})
}

func (me *DetectionConfig) UnmarshalHCL(decoder hcl.Decoder) error {
	me.OutOfMemoryDetection = &oom.DetectionConfig{Enabled: false}
	me.OutOfThreadsDetection = &oot.DetectionConfig{Enabled: false}

	if _, ok := decoder.GetOk("out_of_threads.#"); ok {
		me.OutOfThreadsDetection = new(oot.DetectionConfig)
		if err := me.OutOfThreadsDetection.UnmarshalHCL(hcl.NewDecoder(decoder, "out_of_threads", 0)); err != nil {
			return err
		}
	}
	if _, ok := decoder.GetOk("out_of_memory.#"); ok {
		me.OutOfMemoryDetection = new(oom.DetectionConfig)
		if err := me.OutOfMemoryDetection.UnmarshalHCL(hcl.NewDecoder(decoder, "out_of_memory", 0)); err != nil {
			return err
		}
	}
	return nil
}
