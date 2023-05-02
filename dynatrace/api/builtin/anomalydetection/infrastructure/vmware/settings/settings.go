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

package vmware

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Settings struct {
	DroppedPacketsDetection      *DroppedPacketsDetectionConfig      `json:"droppedPacketsDetection"`
	EsxiHighCpuDetection         *EsxiHighCpuDetectionConfig         `json:"esxiHighCpuDetection"`
	EsxiHighMemoryDetection      *EsxiHighMemoryDetectionConfig      `json:"esxiHighMemoryDetection"`
	GuestCpuLimitDetection       *GuestCPULimitDetectionConfig       `json:"guestCpuLimitDetection"`
	LowDatastoreSpaceDetection   *LowDatastoreSpaceDetectionConfig   `json:"lowDatastoreSpaceDetection"`
	OverloadedStorageDetection   *OverloadedStorageDetectionConfig   `json:"overloadedStorageDetection"`
	SlowPhysicalStorageDetection *SlowPhysicalStorageDetectionConfig `json:"slowPhysicalStorageDetection"`
	UndersizedStorageDetection   *UndersizedStorageDetectionConfig   `json:"undersizedStorageDetection"`
}

func (me *Settings) Name() string {
	return "vmware_anomalies"
}

func (me *Settings) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"dropped_packets_detection": {
			Type:        schema.TypeList,
			Description: "no documentation available",
			Required:    true,
			Elem:        &schema.Resource{Schema: new(DroppedPacketsDetectionConfig).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
		"esxi_high_cpu_detection": {
			Type:        schema.TypeList,
			Description: "no documentation available",
			Required:    true,
			Elem:        &schema.Resource{Schema: new(EsxiHighCpuDetectionConfig).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
		"esxi_high_memory_detection": {
			Type:        schema.TypeList,
			Description: "no documentation available",
			Required:    true,
			Elem:        &schema.Resource{Schema: new(EsxiHighMemoryDetectionConfig).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
		"guest_cpu_limit_detection": {
			Type:        schema.TypeList,
			Description: "no documentation available",
			Required:    true,
			Elem:        &schema.Resource{Schema: new(GuestCPULimitDetectionConfig).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
		"low_datastore_space_detection": {
			Type:        schema.TypeList,
			Description: "no documentation available",
			Required:    true,
			Elem:        &schema.Resource{Schema: new(LowDatastoreSpaceDetectionConfig).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
		"overloaded_storage_detection": {
			Type:        schema.TypeList,
			Description: "no documentation available",
			Required:    true,
			Elem:        &schema.Resource{Schema: new(OverloadedStorageDetectionConfig).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
		"slow_physical_storage_detection": {
			Type:        schema.TypeList,
			Description: "no documentation available",
			Required:    true,
			Elem:        &schema.Resource{Schema: new(SlowPhysicalStorageDetectionConfig).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
		"undersized_storage_detection": {
			Type:        schema.TypeList,
			Description: "no documentation available",
			Required:    true,
			Elem:        &schema.Resource{Schema: new(UndersizedStorageDetectionConfig).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
	}
}

func (me *Settings) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"dropped_packets_detection":       me.DroppedPacketsDetection,
		"esxi_high_cpu_detection":         me.EsxiHighCpuDetection,
		"esxi_high_memory_detection":      me.EsxiHighMemoryDetection,
		"guest_cpu_limit_detection":       me.GuestCpuLimitDetection,
		"low_datastore_space_detection":   me.LowDatastoreSpaceDetection,
		"overloaded_storage_detection":    me.OverloadedStorageDetection,
		"slow_physical_storage_detection": me.SlowPhysicalStorageDetection,
		"undersized_storage_detection":    me.UndersizedStorageDetection,
	})
}

func (me *Settings) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"dropped_packets_detection":       &me.DroppedPacketsDetection,
		"esxi_high_cpu_detection":         &me.EsxiHighCpuDetection,
		"esxi_high_memory_detection":      &me.EsxiHighMemoryDetection,
		"guest_cpu_limit_detection":       &me.GuestCpuLimitDetection,
		"low_datastore_space_detection":   &me.LowDatastoreSpaceDetection,
		"overloaded_storage_detection":    &me.OverloadedStorageDetection,
		"slow_physical_storage_detection": &me.SlowPhysicalStorageDetection,
		"undersized_storage_detection":    &me.UndersizedStorageDetection,
	})
}
