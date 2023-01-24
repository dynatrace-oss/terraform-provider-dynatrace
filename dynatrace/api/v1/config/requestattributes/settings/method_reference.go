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

// MethodReference Configuration of a method to be captured.
type MethodReference struct {
	ReturnType      string                     `json:"returnType"`                // The return type.
	Visibility      Visibility                 `json:"visibility"`                // The visibility of the method to capture.
	ArgumentTypes   []string                   `json:"argumentTypes"`             // The list of argument types.
	ClassName       *string                    `json:"className,omitempty"`       // The class name where the method to capture resides.   Either this or the **fileName** must be set.
	FileName        *string                    `json:"fileName,omitempty"`        // The file name where the method to capture resides.   Either this or **className** must be set.
	FileNameMatcher *FileNameMatcher           `json:"fileNameMatcher,omitempty"` // The operator of the comparison.   If not set, `EQUALS` is used.
	MethodName      string                     `json:"methodName"`                // The name of the method to capture.
	Modifiers       []Modifier                 `json:"modifiers"`                 // The modifiers of the method to capture.
	Unknowns        map[string]json.RawMessage `json:"-"`
}

func (me *MethodReference) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"return_type": {
			Type:        schema.TypeString,
			Description: "The return type",
			Required:    true,
		},
		"visibility": {
			Type:        schema.TypeString,
			Description: "The visibility of the method to capture",
			Required:    true,
		},
		"argument_types": {
			Type:        schema.TypeList,
			Optional:    true,
			MinItems:    1,
			Description: "Configuration of a method to be captured",
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
		"class_name": {
			Type:        schema.TypeString,
			Description: "The class name where the method to capture resides.   Either this or the **fileName** must be set",
			Optional:    true,
		},
		"file_name": {
			Type:        schema.TypeString,
			Description: "The file name where the method to capture resides.   Either this or **className** must be set",
			Optional:    true,
		},
		"file_name_matcher": {
			Type:        schema.TypeString,
			Description: "The operator of the comparison. If not set, `EQUALS` is used",
			Optional:    true,
		},
		"method_name": {
			Type:        schema.TypeString,
			Description: "The name of the method to capture",
			Required:    true,
		},
		"modifiers": {
			Type:        schema.TypeSet,
			Optional:    true,
			MinItems:    1,
			Description: "The modifiers of the method to capture",
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
		"unknowns": {
			Type:        schema.TypeString,
			Description: "allows for configuring properties that are not explicitly supported by the current version of this provider",
			Optional:    true,
		},
	}
}

func (me *MethodReference) MarshalHCL(properties hcl.Properties) error {
	if err := properties.Unknowns(me.Unknowns); err != nil {
		return err
	}
	if err := properties.Encode("return_type", me.ReturnType); err != nil {
		return err
	}
	if err := properties.Encode("visibility", Visibility(me.Visibility)); err != nil {
		return err
	}
	if err := properties.Encode("argument_types", me.ArgumentTypes); err != nil {
		return err
	}
	if err := properties.Encode("class_name", me.ClassName); err != nil {
		return err
	}
	if err := properties.Encode("file_name", me.FileName); err != nil {
		return err
	}
	if err := properties.Encode("file_name_matcher", me.FileNameMatcher); err != nil {
		return err
	}
	if err := properties.Encode("method_name", me.MethodName); err != nil {
		return err
	}
	if len(me.Modifiers) > 0 {
		mods := []string{}
		for _, mod := range me.Modifiers {
			mods = append(mods, string(mod))
		}
		if err := properties.Encode("modifiers", mods); err != nil {
			return err
		}
	}
	return nil
}

func (me *MethodReference) UnmarshalHCL(decoder hcl.Decoder) error {
	if value, ok := decoder.GetOk("unknowns"); ok {
		if err := json.Unmarshal([]byte(value.(string)), me); err != nil {
			return err
		}
		if err := json.Unmarshal([]byte(value.(string)), &me.Unknowns); err != nil {
			return err
		}
		delete(me.Unknowns, "return_type")
		delete(me.Unknowns, "visibility")
		delete(me.Unknowns, "argument_types")
		delete(me.Unknowns, "class_name")
		delete(me.Unknowns, "file_name")
		delete(me.Unknowns, "file_name_matcher")
		delete(me.Unknowns, "method_name")
		delete(me.Unknowns, "modifiers")
		if len(me.Unknowns) == 0 {
			me.Unknowns = nil
		}
	}
	if value, ok := decoder.GetOk("return_type"); ok {
		me.ReturnType = value.(string)
	}
	if value, ok := decoder.GetOk("visibility"); ok {
		me.Visibility = Visibility(value.(string))
	}
	if value, ok := decoder.GetOk("argument_types"); ok {
		me.ArgumentTypes = []string{}
		for _, e := range value.([]any) {
			me.ArgumentTypes = append(me.ArgumentTypes, e.(string))
		}
	}
	if value, ok := decoder.GetOk("class_name"); ok {
		me.ClassName = opt.NewString(value.(string))
	}
	if value, ok := decoder.GetOk("file_name"); ok {
		me.FileName = opt.NewString(value.(string))
	}
	if value, ok := decoder.GetOk("file_name_matcher"); ok {
		me.FileNameMatcher = FileNameMatcher(value.(string)).Ref()
	}
	if value, ok := decoder.GetOk("method_name"); ok {
		me.MethodName = value.(string)
	}
	if value, ok := decoder.GetOk("modifiers"); ok {
		me.Modifiers = []Modifier{}
		for _, e := range value.(*schema.Set).List() {
			me.Modifiers = append(me.Modifiers, Modifier(e.(string)))
		}
	}
	return nil
}

func (me *MethodReference) MarshalJSON() ([]byte, error) {
	m := xjson.NewProperties(me.Unknowns)
	if err := m.Marshal("returnType", me.ReturnType); err != nil {
		return nil, err
	}
	if err := m.Marshal("visibility", me.Visibility); err != nil {
		return nil, err
	}
	if err := m.Marshal("argumentTypes", me.ArgumentTypes); err != nil {
		return nil, err
	}
	if err := m.Marshal("className", me.ClassName); err != nil {
		return nil, err
	}
	if err := m.Marshal("fileName", me.FileName); err != nil {
		return nil, err
	}
	if err := m.Marshal("fileNameMatcher", me.FileNameMatcher); err != nil {
		return nil, err
	}
	if err := m.Marshal("methodName", me.MethodName); err != nil {
		return nil, err
	}
	if err := m.Marshal("modifiers", me.Modifiers); err != nil {
		return nil, err
	}
	return json.Marshal(m)
}

func (me *MethodReference) UnmarshalJSON(data []byte) error {
	m := xjson.Properties{}
	if err := json.Unmarshal(data, &m); err != nil {
		return err
	}
	if err := m.Unmarshal("returnType", &me.ReturnType); err != nil {
		return err
	}
	if err := m.Unmarshal("visibility", &me.Visibility); err != nil {
		return err
	}
	if err := m.Unmarshal("argumentTypes", &me.ArgumentTypes); err != nil {
		return err
	}
	if err := m.Unmarshal("className", &me.ClassName); err != nil {
		return err
	}
	if err := m.Unmarshal("fileName", &me.FileName); err != nil {
		return err
	}
	if err := m.Unmarshal("fileNameMatcher", &me.FileNameMatcher); err != nil {
		return err
	}
	if err := m.Unmarshal("methodName", &me.MethodName); err != nil {
		return err
	}
	if err := m.Unmarshal("modifiers", &me.Modifiers); err != nil {
		return err
	}

	if len(m) > 0 {
		me.Unknowns = m
	}
	return nil
}
