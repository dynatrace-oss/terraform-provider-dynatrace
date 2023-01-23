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

package strategy

import (
	"encoding/json"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/xjson"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Wrapper struct {
	Strategy MonitoringStrategy
}

func (me *Wrapper) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"auto": {
			Type:        schema.TypeList,
			Optional:    true,
			MaxItems:    1,
			Description: "An auto-adaptive baseline strategy to detect anomalies within metrics that show a regular change over time, as the baseline is also updated automatically. An example is to detect an anomaly in the number of received network packets or within the number of user actions over time",
			Elem:        &schema.Resource{Schema: new(Auto).Schema()},
		},
		"static": {
			Type:        schema.TypeList,
			Optional:    true,
			MaxItems:    1,
			Description: "A static threshold monitoring strategy to alert on hard limits within a given metric. An example is the violation of a critical memory limit",
			Elem:        &schema.Resource{Schema: new(Static).Schema()},
		},
		"generic": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "A generic monitoring strategy",
			Elem:        &schema.Resource{Schema: new(BaseMonitoringStrategy).Schema()},
		},
	}
}

func (me *Wrapper) MarshalHCL(properties hcl.Properties) error {
	if me.Strategy != nil {
		switch strategy := me.Strategy.(type) {
		case *Auto:
			if err := properties.Encode("auto", strategy); err != nil {
				return err
			}
		case *Static:
			if err := properties.Encode("static", strategy); err != nil {
				return err
			}
		case *BaseMonitoringStrategy:
			if err := properties.Encode("generic", strategy); err != nil {
				return err
			}
		default:
		}
	}
	return nil
}

func (me *Wrapper) UnmarshalHCL(decoder hcl.Decoder) error {
	if _, ok := decoder.GetOk("auto.#"); ok {
		cfg := new(Auto)
		if err := cfg.UnmarshalHCL(hcl.NewDecoder(decoder, "auto", 0)); err != nil {
			return err
		}
		me.Strategy = cfg
	}
	if _, ok := decoder.GetOk("static.#"); ok {
		cfg := new(Static)
		if err := cfg.UnmarshalHCL(hcl.NewDecoder(decoder, "static", 0)); err != nil {
			return err
		}
		me.Strategy = cfg
	}
	if _, ok := decoder.GetOk("generic.#"); ok {
		cfg := new(BaseMonitoringStrategy)
		if err := cfg.UnmarshalHCL(hcl.NewDecoder(decoder, "generic", 0)); err != nil {
			return err
		}
		me.Strategy = cfg
	}
	return nil
}

func (me *Wrapper) UnmarshalJSON(data []byte) error {
	properties := xjson.Properties{}
	if err := json.Unmarshal(data, &properties); err != nil {
		return err
	}
	if rawType, found := properties["type"]; found {
		var sType string
		if err := json.Unmarshal(rawType, &sType); err != nil {
			return err
		}
		switch sType {
		case string(Types.AutoAdaptiveBaseline):
			cfg := new(Auto)
			if err := json.Unmarshal(data, &cfg); err != nil {
				return err
			}
			me.Strategy = cfg
		case string(Types.StaticThreshold):
			cfg := new(Static)
			if err := json.Unmarshal(data, &cfg); err != nil {
				return err
			}
			me.Strategy = cfg
		default:
			cfg := new(BaseMonitoringStrategy)
			if err := json.Unmarshal(data, &cfg); err != nil {
				return err
			}
			me.Strategy = cfg
		}
	}
	return nil
}
