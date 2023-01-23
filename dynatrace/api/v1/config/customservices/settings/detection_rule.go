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

package customservices

import (
	"encoding/json"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/xjson"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// DetectionRule the defining rules for a CustomService
type DetectionRule struct {
	Enabled          bool                       `json:"enabled"`                   // Rule enabled/disabled
	FileName         *string                    `json:"fileName,omitempty"`        // The PHP file containing the class or methods to instrument. Required for PHP custom service. Not applicable to Java and .NET
	FileNameMatcher  *FileNameMatcher           `json:"fileNameMatcher,omitempty"` // Matcher applying to the file name. Default value is `ENDS_WITH` (if applicable)
	ClassName        *string                    `json:"className,omitempty"`       // The fully qualified class or interface to instrument. Required for Java and .NET custom services. Not applicable to PHP
	ClassNameMatcher *ClassNameMatcher          `json:"matcher,omitempty"`         // Matcher applying to the class name. `STARTS_WITH` can only be used if there is at least one annotation defined. Default value is `EQUALS`
	MethodRules      []*MethodRule              `json:"methodRules"`               // List of methods to instrument
	Annotations      []string                   `json:"annotations"`               // Additional annotations filter of the rule. Only classes where all listed annotations are available in the class itself or any of its superclasses are instrumented. nNot applicable to PHP
	Unknowns         map[string]json.RawMessage `json:"-"`
}

func (me *DetectionRule) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"enabled": {
			Type:        schema.TypeBool,
			Description: "Rule enabled/disabled",
			Required:    true,
		},
		"file": {
			Type:        schema.TypeList,
			Description: "The PHP file containing the class or methods to instrument. Required for PHP custom service. Not applicable to Java and .NET",
			Optional:    true,
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: new(FileSection).Schema(),
			},
		},
		"class": {
			Type:        schema.TypeList,
			Description: "The fully qualified class or interface to instrument (or a substring if matching to a string). Required for Java and .NET custom services. Not applicable to PHP",
			Optional:    true,
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: new(ClassSection).Schema(),
			},
		},
		"annotations": {
			Type:        schema.TypeSet,
			Description: "Additional annotations filter of the rule. Only classes where all listed annotations are available in the class itself or any of its superclasses are instrumented. Not applicable to PHP",
			Optional:    true,
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
		"method": {
			Type:        schema.TypeList,
			Description: "methods to instrument",
			Required:    true,
			Elem: &schema.Resource{
				Schema: new(MethodRule).Schema(),
			},
		},
		"unknowns": {
			Type:        schema.TypeString,
			Description: "allows for configuring properties that are not explicitly supported by the current version of this provider",
			Optional:    true,
		},
	}
}

func (me *DetectionRule) MarshalHCL(properties hcl.Properties) error {
	if len(me.Unknowns) > 0 {
		delete(me.Unknowns, "id")
		data, err := json.Marshal(me.Unknowns)
		if err != nil {
			return err
		}
		if err := properties.Encode("unknowns", string(data)); err != nil {
			return err
		}
	}
	// if me.ID != nil {
	// 	if err := properties.Encode("id", me.ID); err != nil { return err }
	// }
	if err := properties.Encode("enabled", me.Enabled); err != nil {
		return err
	}
	if me.FileName != nil || me.FileNameMatcher != nil {
		fileSection := &FileSection{
			Name:  me.FileName,
			Match: me.FileNameMatcher,
		}
		if !fileSection.IsEmpty() {
			if err := properties.Encode("file", fileSection); err != nil {
				return err
			}
		}
	}
	if me.ClassName != nil || me.ClassNameMatcher != nil {
		classSection := &ClassSection{
			Name:  me.ClassName,
			Match: me.ClassNameMatcher,
		}
		if err := properties.Encode("class", classSection); err != nil {
			return err
		}
	}
	if err := properties.Encode("method", me.MethodRules); err != nil {
		return err
	}
	if err := properties.Encode("annotations", me.Annotations); err != nil {
		return err
	}
	return nil
}

func (me *DetectionRule) UnmarshalHCL(decoder hcl.Decoder) error {
	if value, ok := decoder.GetOk("unknowns"); ok {
		if err := json.Unmarshal([]byte(value.(string)), me); err != nil {
			return err
		}
		if err := json.Unmarshal([]byte(value.(string)), &me.Unknowns); err != nil {
			return err
		}
		delete(me.Unknowns, "id")
		delete(me.Unknowns, "enabled")
		delete(me.Unknowns, "file")
		delete(me.Unknowns, "class")
		delete(me.Unknowns, "annotations")
		delete(me.Unknowns, "method")
		if len(me.Unknowns) == 0 {
			me.Unknowns = nil
		}
	}
	if value, ok := decoder.GetOk("enabled"); ok {
		me.Enabled = value.(bool)
	}
	if _, ok := decoder.GetOk("file.#"); ok {
		fileSection := new(FileSection)
		if err := fileSection.UnmarshalHCL(hcl.NewDecoder(decoder, "file", 0)); err != nil {
			return err
		}
		me.FileName = fileSection.Name
		me.FileNameMatcher = fileSection.Match
	}
	if _, ok := decoder.GetOk("class.#"); ok {
		classSection := new(ClassSection)
		if err := classSection.UnmarshalHCL(hcl.NewDecoder(decoder, "class", 0)); err != nil {
			return err
		}
		me.ClassName = classSection.Name
		me.ClassNameMatcher = classSection.Match
	}
	if result, ok := decoder.GetOk("method.#"); ok {
		me.MethodRules = []*MethodRule{}
		for idx := 0; idx < result.(int); idx++ {
			entry := new(MethodRule)
			if err := entry.UnmarshalHCL(hcl.NewDecoder(decoder, "method", idx)); err != nil {
				return err
			}
			me.MethodRules = append(me.MethodRules, entry)
		}
	}

	decoder.Decode("annotations", &me.Annotations)
	return nil
}

func (me *DetectionRule) MarshalJSON() ([]byte, error) {
	m := xjson.NewProperties(me.Unknowns)
	delete(m, "id")
	if err := m.Marshal("enabled", me.Enabled); err != nil {
		return nil, err
	}
	if err := m.Marshal("fileName", me.FileName); err != nil {
		return nil, err
	}
	if err := m.Marshal("fileNameMatcher", me.FileNameMatcher); err != nil {
		return nil, err
	}
	if err := m.Marshal("className", me.ClassName); err != nil {
		return nil, err
	}
	if err := m.Marshal("matcher", me.ClassNameMatcher); err != nil {
		return nil, err
	}
	if err := m.Marshal("methodRules", me.MethodRules); err != nil {
		return nil, err
	}
	if err := m.Marshal("annotations", me.Annotations); err != nil {
		return nil, err
	}
	return json.Marshal(m)
}

func (me *DetectionRule) UnmarshalJSON(data []byte) error {
	m := xjson.Properties{}
	if err := json.Unmarshal(data, &m); err != nil {
		return err
	}
	delete(m, "id")
	if err := m.Unmarshal("enabled", &me.Enabled); err != nil {
		return err
	}
	if err := m.Unmarshal("fileName", &me.FileName); err != nil {
		return err
	}
	if err := m.Unmarshal("fileNameMatcher", &me.FileNameMatcher); err != nil {
		return err
	}
	if err := m.Unmarshal("className", &me.ClassName); err != nil {
		return err
	}
	if err := m.Unmarshal("matcher", &me.ClassNameMatcher); err != nil {
		return err
	}
	if err := m.Unmarshal("methodRules", &me.MethodRules); err != nil {
		return err
	}
	if err := m.Unmarshal("annotations", &me.Annotations); err != nil {
		return err
	}

	if len(m) > 0 {
		me.Unknowns = m
	}
	return nil
}
