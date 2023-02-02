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

package disks

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Disk struct {
	DiskLowInodesDetection          *DiskLowInodesDetection          `json:"diskLowInodesDetection"`
	DiskLowSpaceDetection           *DiskLowSpaceDetection           `json:"diskLowSpaceDetection"`
	DiskSlowWritesAndReadsDetection *DiskSlowWritesAndReadsDetection `json:"diskSlowWritesAndReadsDetection"`
}

func (me *Disk) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"disk_low_inodes_detection": {
			Type:        schema.TypeList,
			Description: "no documentation available",
			Required:    true,
			Elem:        &schema.Resource{Schema: new(DiskLowInodesDetection).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
		"disk_low_space_detection": {
			Type:        schema.TypeList,
			Description: "no documentation available",
			Required:    true,
			Elem:        &schema.Resource{Schema: new(DiskLowSpaceDetection).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
		"disk_slow_writes_and_reads_detection": {
			Type:        schema.TypeList,
			Description: "no documentation available",
			Required:    true,
			Elem:        &schema.Resource{Schema: new(DiskSlowWritesAndReadsDetection).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
	}
}

func (me *Disk) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"disk_low_inodes_detection":            me.DiskLowInodesDetection,
		"disk_low_space_detection":             me.DiskLowSpaceDetection,
		"disk_slow_writes_and_reads_detection": me.DiskSlowWritesAndReadsDetection,
	})
}

func (me *Disk) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"disk_low_inodes_detection":            &me.DiskLowInodesDetection,
		"disk_low_space_detection":             &me.DiskLowSpaceDetection,
		"disk_slow_writes_and_reads_detection": &me.DiskSlowWritesAndReadsDetection,
	})
}
