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

package traffic

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Detection struct {
	Drops  *DropDetection  // The configuration of traffic drops detection.
	Spikes *SpikeDetection // The configuration of traffic spikes detection.
}

func (me *Detection) IsEmpty() bool {
	if me.Drops != nil {
		if me.Drops.Enabled {
			return false
		}
		me.Drops = nil
	}
	if me.Spikes != nil {
		if me.Spikes.Enabled {
			return false
		}
		me.Spikes = nil
	}
	return true
}

func (me *Detection) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"drops": {
			Type:        schema.TypeList,
			Optional:    true,
			MaxItems:    1,
			Description: "The configuration of traffic drops detection",
			Elem:        &schema.Resource{Schema: new(DropDetection).Schema()},
		},
		"spikes": {
			Type:        schema.TypeList,
			Optional:    true,
			MaxItems:    1,
			Description: "The configuration of traffic spikes detection",
			Elem:        &schema.Resource{Schema: new(SpikeDetection).Schema()},
		},
	}
}

func (me *Detection) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"drops":  me.Drops,
		"spikes": me.Spikes,
	})
}

func (me *Detection) UnmarshalHCL(decoder hcl.Decoder) error {
	if _, ok := decoder.GetOk("drops.#"); ok {
		me.Drops = new(DropDetection)
		if err := me.Drops.UnmarshalHCL(hcl.NewDecoder(decoder, "drops", 0)); err != nil {
			return err
		}
		if !me.Drops.Enabled {
			me.Drops = nil
		}
	}
	if _, ok := decoder.GetOk("spikes.#"); ok {
		me.Spikes = new(SpikeDetection)
		if err := me.Spikes.UnmarshalHCL(hcl.NewDecoder(decoder, "spikes", 0)); err != nil {
			return err
		}
		if !me.Spikes.Enabled {
			me.Spikes = nil
		}
	}

	return nil
}
