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
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/anomalies/hosts/settings/disks/inodes"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/anomalies/hosts/settings/disks/slow"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/anomalies/hosts/settings/disks/space"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type DetectionConfig struct {
	Speed  *slow.DetectionConfig   `json:"diskSlowWritesAndReadsDetection"` // Configuration of slow running disks detection.
	Space  *space.DetectionConfig  `json:"diskLowSpaceDetection"`           // Configuration of low disk space detection.
	Inodes *inodes.DetectionConfig `json:"diskLowInodesDetection"`          // Configuration of low disk inodes number detection.
}

func (me *DetectionConfig) IsConfigured() bool {
	if me.Speed != nil && me.Speed.Enabled {
		return true
	}
	if me.Space != nil && me.Space.Enabled {
		return true
	}
	if me.Inodes != nil && me.Inodes.Enabled {
		return true
	}
	return false
}

func (me *DetectionConfig) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"space": {
			Type:        schema.TypeList,
			Optional:    true,
			MaxItems:    1,
			Description: "Configuration of low disk space detection",
			Elem:        &schema.Resource{Schema: new(space.DetectionConfig).Schema()},
		},
		"speed": {
			Type:        schema.TypeList,
			Optional:    true,
			MaxItems:    1,
			Description: "Configuration of slow running disks detection",
			Elem:        &schema.Resource{Schema: new(slow.DetectionConfig).Schema()},
		},
		"inodes": {
			Type:        schema.TypeList,
			Optional:    true,
			MaxItems:    1,
			Description: "Configuration of low disk inodes number detection",
			Elem:        &schema.Resource{Schema: new(inodes.DetectionConfig).Schema()},
		},
	}
}

func (me *DetectionConfig) MarshalHCL(properties hcl.Properties) error {
	if !me.IsConfigured() {
		return nil
	}
	return properties.EncodeAll(map[string]any{
		"space":  me.Space,
		"speed":  me.Speed,
		"inodes": me.Inodes,
	})
}

func (me *DetectionConfig) UnmarshalHCL(decoder hcl.Decoder) error {
	me.Speed = &slow.DetectionConfig{Enabled: false}
	me.Space = &space.DetectionConfig{Enabled: false}
	me.Inodes = &inodes.DetectionConfig{Enabled: false}

	if _, ok := decoder.GetOk("space.#"); ok {
		me.Space = new(space.DetectionConfig)
		if err := me.Space.UnmarshalHCL(hcl.NewDecoder(decoder, "space", 0)); err != nil {
			return err
		}
	}
	if _, ok := decoder.GetOk("speed.#"); ok {
		me.Speed = new(slow.DetectionConfig)
		if err := me.Speed.UnmarshalHCL(hcl.NewDecoder(decoder, "speed", 0)); err != nil {
			return err
		}
	}
	if _, ok := decoder.GetOk("inodes.#"); ok {
		me.Inodes = new(inodes.DetectionConfig)
		if err := me.Inodes.UnmarshalHCL(hcl.NewDecoder(decoder, "inodes", 0)); err != nil {
			return err
		}
	}
	return nil
}
