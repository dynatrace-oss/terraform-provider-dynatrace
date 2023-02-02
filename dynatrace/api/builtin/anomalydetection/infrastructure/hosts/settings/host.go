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

package hosts

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Host struct {
	ConnectionLostDetection    *ConnectionLostDetection    `json:"connectionLostDetection"`
	HighCpuSaturationDetection *HighCpuSaturationDetection `json:"highCpuSaturationDetection"`
	HighGcActivityDetection    *HighGcActivityDetection    `json:"highGcActivityDetection"`
	HighMemoryDetection        *HighMemoryDetection        `json:"highMemoryDetection"`
	OutOfMemoryDetection       *OutOfMemoryDetection       `json:"outOfMemoryDetection"`
	OutOfThreadsDetection      *OutOfThreadsDetection      `json:"outOfThreadsDetection"`
}

func (me *Host) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"connection_lost_detection": {
			Type:        schema.TypeList,
			Description: "no documentation available",
			Required:    true,
			Elem:        &schema.Resource{Schema: new(ConnectionLostDetection).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
		"high_cpu_saturation_detection": {
			Type:        schema.TypeList,
			Description: "no documentation available",
			Required:    true,
			Elem:        &schema.Resource{Schema: new(HighCpuSaturationDetection).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
		"high_gc_activity_detection": {
			Type:        schema.TypeList,
			Description: "no documentation available",
			Required:    true,
			Elem:        &schema.Resource{Schema: new(HighGcActivityDetection).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
		"high_memory_detection": {
			Type:        schema.TypeList,
			Description: "no documentation available",
			Required:    true,
			Elem:        &schema.Resource{Schema: new(HighMemoryDetection).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
		"out_of_memory_detection": {
			Type:        schema.TypeList,
			Description: "no documentation available",
			Required:    true,
			Elem:        &schema.Resource{Schema: new(OutOfMemoryDetection).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
		"out_of_threads_detection": {
			Type:        schema.TypeList,
			Description: "no documentation available",
			Required:    true,
			Elem:        &schema.Resource{Schema: new(OutOfThreadsDetection).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
	}
}

func (me *Host) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"connection_lost_detection":     me.ConnectionLostDetection,
		"high_cpu_saturation_detection": me.HighCpuSaturationDetection,
		"high_gc_activity_detection":    me.HighGcActivityDetection,
		"high_memory_detection":         me.HighMemoryDetection,
		"out_of_memory_detection":       me.OutOfMemoryDetection,
		"out_of_threads_detection":      me.OutOfThreadsDetection,
	})
}

func (me *Host) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"connection_lost_detection":     &me.ConnectionLostDetection,
		"high_cpu_saturation_detection": &me.HighCpuSaturationDetection,
		"high_gc_activity_detection":    &me.HighGcActivityDetection,
		"high_memory_detection":         &me.HighMemoryDetection,
		"out_of_memory_detection":       &me.OutOfMemoryDetection,
		"out_of_threads_detection":      &me.OutOfThreadsDetection,
	})
}
