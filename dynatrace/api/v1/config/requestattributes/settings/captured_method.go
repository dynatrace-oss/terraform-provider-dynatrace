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

package requestattributes

import (
	"encoding/json"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/xjson"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/opt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// CapturedMethod has no documentation
type CapturedMethod struct {
	ArgumentIndex    *int32                     `json:"argumentIndex,omitempty"`    // The index of the argument to capture. Set `0` to capture the return value, `1` or higher to capture a mehtod argument.   Required if the **capture** is set to `ARGUMENT`.  Not applicable in other cases.
	Capture          Capture                    `json:"capture"`                    // What to capture from the method.
	DeepObjectAccess *string                    `json:"deepObjectAccess,omitempty"` // The getter chain to apply to the captured object. It is required in one of the following cases:  The **capture** is set to `THIS`.    The **capture** is set to `ARGUMENT`, and the argument is not a primitive, a primitive wrapper class, a string, or an array.   Not applicable in other cases.
	Method           *MethodReference           `json:"method"`                     // Configuration of a method to be captured.
	Unknowns         map[string]json.RawMessage `json:"-"`
}

func (me *CapturedMethod) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"argument_index": {
			Type:        schema.TypeInt,
			Description: "The index of the argument to capture. Set `0` to capture the return value, `1` or higher to capture a mehtod argument.   Required if the **capture** is set to `ARGUMENT`.  Not applicable in other cases",
			Optional:    true,
		},
		"capture": {
			Type:        schema.TypeString,
			Description: "What to capture from the method",
			Required:    true,
		},
		"deep_object_access": {
			Type:        schema.TypeString,
			Description: "The getter chain to apply to the captured object. It is required in one of the following cases:  The **capture** is set to `THIS`.    The **capture** is set to `ARGUMENT`, and the argument is not a primitive, a primitive wrapper class, a string, or an array.   Not applicable in other cases",
			Optional:    true,
		},
		"method": {
			Type:        schema.TypeList,
			Optional:    true,
			MaxItems:    1,
			Description: "Configuration of a method to be captured",
			Elem: &schema.Resource{
				Schema: new(MethodReference).Schema(),
			},
		},
		"unknowns": {
			Type:        schema.TypeString,
			Description: "allows for configuring properties that are not explicitly supported by the current version of this provider",
			Optional:    true,
		},
	}
}

func (me *CapturedMethod) MarshalHCL(properties hcl.Properties) error {
	if err := properties.Unknowns(me.Unknowns); err != nil {
		return err
	}
	if err := properties.Encode("argument_index", int(opt.Int32(me.ArgumentIndex))); err != nil {
		return err
	}
	if err := properties.Encode("capture", string(me.Capture)); err != nil {
		return err
	}
	if err := properties.Encode("deep_object_access", me.DeepObjectAccess); err != nil {
		return err
	}
	if err := properties.Encode("method", me.Method); err != nil {
		return err
	}
	return nil
}

func (me *CapturedMethod) UnmarshalHCL(decoder hcl.Decoder) error {
	if value, ok := decoder.GetOk("unknowns"); ok {
		if err := json.Unmarshal([]byte(value.(string)), me); err != nil {
			return err
		}
		if err := json.Unmarshal([]byte(value.(string)), &me.Unknowns); err != nil {
			return err
		}
		delete(me.Unknowns, "argument_index")
		delete(me.Unknowns, "capture")
		delete(me.Unknowns, "deep_object_access")
		delete(me.Unknowns, "method")
		if len(me.Unknowns) == 0 {
			me.Unknowns = nil
		}
	}
	if value, ok := decoder.GetOk("argument_index"); ok {
		me.ArgumentIndex = opt.NewInt32(int32(value.(int)))
	}
	if value, ok := decoder.GetOk("capture"); ok {
		me.Capture = Capture(value.(string))
	}
	if me.Capture == Captures.Argument {
		if me.ArgumentIndex == nil {
			me.ArgumentIndex = opt.NewInt32(0)
		}
	}
	if value, ok := decoder.GetOk("deep_object_access"); ok {
		me.DeepObjectAccess = opt.NewString(value.(string))
	}
	if _, ok := decoder.GetOk("method.#"); ok {
		me.Method = new(MethodReference)
		if err := me.Method.UnmarshalHCL(hcl.NewDecoder(decoder, "method", 0)); err != nil {
			return err
		}
	}
	return nil
}

func (me *CapturedMethod) MarshalJSON() ([]byte, error) {
	m := xjson.NewProperties(me.Unknowns)
	if err := m.Marshal("argumentIndex", me.ArgumentIndex); err != nil {
		return nil, err
	}
	if err := m.Marshal("capture", me.Capture); err != nil {
		return nil, err
	}
	if err := m.Marshal("deepObjectAccess", me.DeepObjectAccess); err != nil {
		return nil, err
	}
	if err := m.Marshal("method", me.Method); err != nil {
		return nil, err
	}
	return json.Marshal(m)
}

func (me *CapturedMethod) UnmarshalJSON(data []byte) error {
	m := xjson.Properties{}
	if err := json.Unmarshal(data, &m); err != nil {
		return err
	}
	if err := m.Unmarshal("argumentIndex", &me.ArgumentIndex); err != nil {
		return err
	}
	if err := m.Unmarshal("capture", &me.Capture); err != nil {
		return err
	}
	if err := m.Unmarshal("deepObjectAccess", &me.DeepObjectAccess); err != nil {
		return err
	}
	if err := m.Unmarshal("method", &me.Method); err != nil {
		return err
	}

	if len(m) > 0 {
		me.Unknowns = m
	}
	return nil
}
