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

package perdiskoverride

import (
	"encoding/json"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Settings struct {
	DiskId                              string                           `json:"-"`
	DiskLowInodesDetection              *DiskLowInodesDetection          `json:"diskLowInodesDetection"`
	DiskLowSpaceDetection               *DiskLowSpaceDetection           `json:"diskLowSpaceDetection"`
	DiskSlowWritesAndReadsDetection     *DiskSlowWritesAndReadsDetection `json:"diskSlowWritesAndReadsDetection"`
	OverrideDiskLowSpaceDetection       bool                             `json:"overrideDiskLowSpaceDetection"`       // Override low disk space detection settings
	OverrideLowInodesDetection          bool                             `json:"overrideLowInodesDetection"`          // Override low inodes detection settings
	OverrideSlowWritesAndReadsDetection bool                             `json:"overrideSlowWritesAndReadsDetection"` // Override slow writes and reads detection settings
}

func (me *Settings) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"disk_id": {
			Type:        schema.TypeString,
			Description: "The id for the disk anomaly detection",
			Required:    true,
		},
		"disk_low_inodes_detection": {
			Type:        schema.TypeList,
			Description: "no documentation available",
			Optional:    true,
			Elem:        &schema.Resource{Schema: new(DiskLowInodesDetection).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
		"disk_low_space_detection": {
			Type:        schema.TypeList,
			Description: "no documentation available",
			Optional:    true,
			Elem:        &schema.Resource{Schema: new(DiskLowSpaceDetection).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
		"disk_slow_writes_and_reads_detection": {
			Type:        schema.TypeList,
			Description: "no documentation available",
			Optional:    true,
			Elem:        &schema.Resource{Schema: new(DiskSlowWritesAndReadsDetection).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
		"override_disk_low_space_detection": {
			Type:        schema.TypeBool,
			Description: "Override low disk space detection settings",
			Required:    true,
		},
		"override_low_inodes_detection": {
			Type:        schema.TypeBool,
			Description: "Override low inodes detection settings",
			Required:    true,
		},
		"override_slow_writes_and_reads_detection": {
			Type:        schema.TypeBool,
			Description: "Override slow writes and reads detection settings",
			Required:    true,
		},
	}
}

func (me *Settings) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"disk_id":                                  me.DiskId,
		"disk_low_inodes_detection":                me.DiskLowInodesDetection,
		"disk_low_space_detection":                 me.DiskLowSpaceDetection,
		"disk_slow_writes_and_reads_detection":     me.DiskSlowWritesAndReadsDetection,
		"override_disk_low_space_detection":        me.OverrideDiskLowSpaceDetection,
		"override_low_inodes_detection":            me.OverrideLowInodesDetection,
		"override_slow_writes_and_reads_detection": me.OverrideSlowWritesAndReadsDetection,
	})
}

func (me *Settings) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"disk_id":                                  &me.DiskId,
		"disk_low_inodes_detection":                &me.DiskLowInodesDetection,
		"disk_low_space_detection":                 &me.DiskLowSpaceDetection,
		"disk_slow_writes_and_reads_detection":     &me.DiskSlowWritesAndReadsDetection,
		"override_disk_low_space_detection":        &me.OverrideDiskLowSpaceDetection,
		"override_low_inodes_detection":            &me.OverrideLowInodesDetection,
		"override_slow_writes_and_reads_detection": &me.OverrideSlowWritesAndReadsDetection,
	})
}

func (me *Settings) Name() string {
	return me.DiskId
}

func (me *Settings) SetScope(diskId string) {
	me.DiskId = diskId
}

func (me *Settings) GetScope() string {
	return me.DiskId
}

func (me *Settings) Store() ([]byte, error) {
	var data []byte
	var err error
	if data, err = json.Marshal(me); err != nil {
		return nil, err
	}
	m := map[string]json.RawMessage{}
	if err = json.Unmarshal(data, &m); err != nil {
		return nil, err
	}
	if data, err = json.Marshal(me.DiskId); err != nil {
		return nil, err
	}
	m["disk_id"] = data
	return json.MarshalIndent(m, "", "  ")
}

func (me *Settings) Load(data []byte) error {
	if err := json.Unmarshal(data, &me); err != nil {
		return err
	}

	c := struct {
		DiskId string `json:"disk_id"`
	}{}
	if err := json.Unmarshal(data, &c); err != nil {
		return err
	}
	me.DiskId = c.DiskId

	return nil
}
