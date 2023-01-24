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

// DataSource has no documentation
type DataSource struct {
	CapturingAndStorageLocation *CapturingAndStorageLocation `json:"capturingAndStorageLocation,omitempty"` // Specifies the location where the values are captured and stored.  Required if the **source** is one of the following: `GET_PARAMETER`, `URI`, `REQUEST_HEADER`, `RESPONSE_HEADER`.   Not applicable in other cases.   If the **source** value is `REQUEST_HEADER` or `RESPONSE_HEADER`, the `CAPTURE_AND_STORE_ON_BOTH` location is not allowed.
	Scope                       *ScopeConditions             `json:"scope,omitempty"`                       // Conditions for data capturing.
	ParameterName               *string                      `json:"parameterName,omitempty"`               // The name of the web request parameter to capture.  Required if the **source** is one of the following: `POST_PARAMETER`, `GET_PARAMETER`, `REQUEST_HEADER`, `RESPONSE_HEADER`, `CUSTOM_ATTRIBUTE`.  Not applicable in other cases.
	IIBMethodNodeCondition      *ValueCondition              `json:"iibMethodNodeCondition,omitempty"`      // IBM integration bus label node name condition for which the value is captured.
	Methods                     []*CapturedMethod            `json:"methods,omitempty"`                     // The method specification if the **source** value is `METHOD_PARAM`.   Not applicable in other cases.
	SessionAttributeTechnology  *SessionAttributeTechnology  `json:"sessionAttributeTechnology,omitempty"`  // The technology of the session attribute to capture if the **source** value is `SESSION_ATTRIBUTE`. \n\n Not applicable in other cases.
	Technology                  *Technology                  `json:"technology,omitempty"`                  // The technology of the method to capture if the **source** value is `METHOD_PARAM`. \n\n Not applicable in other cases.
	ValueProcessing             *ValueProcessing             `json:"valueProcessing,omitempty"`             // Process values as specified.
	CICSSDKMethodNodeCondition  *ValueCondition              `json:"cicsSDKMethodNodeCondition,omitempty"`  // IBM integration bus label node name condition for which the value is captured.
	Enabled                     bool                         `json:"enabled"`                               // The data source is enabled (`true`) or disabled (`false`).
	Source                      Source                       `json:"source"`                                // The source of the attribute to capture. Works in conjunction with **parameterName** or **methods** and **technology**.
	IIBLabelMethodNodeCondition *ValueCondition              `json:"iibLabelMethodNodeCondition,omitempty"` // IBM integration bus label node name condition for which the value is captured.
	IIBNodeType                 *IIBNodeType                 `json:"iibNodeType,omitempty"`                 // The IBM integration bus node type for which the value is captured.  This or `iibMethodNodeCondition` is required if the **source** is: `IIB_NODE`.  Not applicable in other cases.
	Unknowns                    map[string]json.RawMessage   `json:"-"`
}

func (me *DataSource) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"capturing_and_storage_location": {
			Type:        schema.TypeString,
			Description: "Specifies the location where the values are captured and stored.  Required if the **source** is one of the following: `GET_PARAMETER`, `URI`, `REQUEST_HEADER`, `RESPONSE_HEADER`.   Not applicable in other cases.   If the **source** value is `REQUEST_HEADER` or `RESPONSE_HEADER`, the `CAPTURE_AND_STORE_ON_BOTH` location is not allowed",
			Optional:    true,
		},
		"scope": {
			Type:        schema.TypeList,
			Optional:    true,
			MaxItems:    1,
			Description: "Conditions for data capturing",
			Elem: &schema.Resource{
				Schema: new(ScopeConditions).Schema(),
			},
		},
		"parameter_name": {
			Type:        schema.TypeString,
			Description: "The name of the web request parameter to capture.  Required if the **source** is one of the following: `POST_PARAMETER`, `GET_PARAMETER`, `REQUEST_HEADER`, `RESPONSE_HEADER`, `CUSTOM_ATTRIBUTE`.  Not applicable in other cases",
			Optional:    true,
		},
		"iib_method_node_condition": {
			Type:        schema.TypeList,
			Optional:    true,
			MaxItems:    1,
			Description: "IBM integration bus label node name condition for which the value is captured",
			Elem: &schema.Resource{
				Schema: new(ValueCondition).Schema(),
			},
		},
		"methods": {
			Type:        schema.TypeList,
			Optional:    true,
			MinItems:    1,
			Description: "The method specification if the **source** value is `METHOD_PARAM`.   Not applicable in other cases",
			Elem: &schema.Resource{
				Schema: new(CapturedMethod).Schema(),
			},
		},
		"session_attribute_technology": {
			Type:        schema.TypeString,
			Description: "The technology of the session attribute to capture if the **source** value is `SESSION_ATTRIBUTE`. \n\n Not applicable in other cases",
			Optional:    true,
		},
		"technology": {
			Type:        schema.TypeString,
			Description: "The technology of the method to capture if the **source** value is `METHOD_PARAM`. \n\n Not applicable in other cases",
			Optional:    true,
		},
		"value_processing": {
			Type:        schema.TypeList,
			Optional:    true,
			MaxItems:    1,
			Description: "Process values as specified",
			Elem: &schema.Resource{
				Schema: new(ValueProcessing).Schema(),
			},
		},
		"cics_sdk_method_node_condition": {
			Type:        schema.TypeList,
			Optional:    true,
			MaxItems:    1,
			Description: "IBM integration bus label node name condition for which the value is captured",
			Elem: &schema.Resource{
				Schema: new(ValueCondition).Schema(),
			},
		},
		"enabled": {
			Type:        schema.TypeBool,
			Description: "The data source is enabled (`true`) or disabled (`false`)",
			Required:    true,
		},
		"source": {
			Type:        schema.TypeString,
			Description: "The source of the attribute to capture. Works in conjunction with **parameterName** or **methods** and **technology**",
			Required:    true,
		},
		"iib_label_method_node_condition": {
			Type:        schema.TypeList,
			Optional:    true,
			MaxItems:    1,
			Description: "IBM integration bus label node name condition for which the value is captured",
			Elem: &schema.Resource{
				Schema: new(ValueCondition).Schema(),
			},
		},
		"iib_node_type": {
			Type:        schema.TypeString,
			Description: "The IBM integration bus node type for which the value is captured.  This or `iibMethodNodeCondition` is required if the **source** is: `IIB_NODE`.  Not applicable in other cases",
			Optional:    true,
		},
		"unknowns": {
			Type:        schema.TypeString,
			Description: "allows for configuring properties that are not explicitly supported by the current version of this provider",
			Optional:    true,
		},
	}
}

func (me *DataSource) MarshalHCL(properties hcl.Properties) error {
	if err := properties.Unknowns(me.Unknowns); err != nil {
		return err
	}
	if err := properties.Encode("capturing_and_storage_location", me.CapturingAndStorageLocation); err != nil {
		return err
	}
	if err := properties.Encode("scope", me.Scope); err != nil {
		return err
	}
	if err := properties.Encode("parameter_name", me.ParameterName); err != nil {
		return err
	}
	if err := properties.Encode("iib_method_node_condition", me.IIBMethodNodeCondition); err != nil {
		return err
	}
	if err := properties.Encode("methods", me.Methods); err != nil {
		return err
	}
	if err := properties.Encode("session_attribute_technology", me.SessionAttributeTechnology); err != nil {
		return err
	}
	if err := properties.Encode("technology", me.Technology); err != nil {
		return err
	}
	if me.ValueProcessing != nil && !me.ValueProcessing.IsZero() {
		if err := properties.Encode("value_processing", me.ValueProcessing); err != nil {
			return err
		}
	}
	if err := properties.Encode("cics_sdk_method_node_condition", me.CICSSDKMethodNodeCondition); err != nil {
		return err
	}
	if err := properties.Encode("enabled", me.Enabled); err != nil {
		return err
	}
	if err := properties.Encode("source", string(me.Source)); err != nil {
		return err
	}
	if err := properties.Encode("iib_label_method_node_condition", me.IIBLabelMethodNodeCondition); err != nil {
		return err
	}
	if err := properties.Encode("iib_node_type", me.IIBNodeType); err != nil {
		return err
	}
	return nil
}

func (me *DataSource) UnmarshalHCL(decoder hcl.Decoder) error {
	if value, ok := decoder.GetOk("unknowns"); ok {
		if err := json.Unmarshal([]byte(value.(string)), me); err != nil {
			return err
		}
		if err := json.Unmarshal([]byte(value.(string)), &me.Unknowns); err != nil {
			return err
		}
		delete(me.Unknowns, "capturing_and_storage_location")
		delete(me.Unknowns, "scope")
		delete(me.Unknowns, "parameter_name")
		delete(me.Unknowns, "iib_method_node_condition")
		delete(me.Unknowns, "methods")
		delete(me.Unknowns, "session_attribute_technology")
		delete(me.Unknowns, "technology")
		delete(me.Unknowns, "value_processing")
		delete(me.Unknowns, "cics_sdk_method_node_condition")
		delete(me.Unknowns, "enabled")
		delete(me.Unknowns, "source")
		delete(me.Unknowns, "iib_label_method_node_condition")
		delete(me.Unknowns, "iib_node_type")
		if len(me.Unknowns) == 0 {
			me.Unknowns = nil
		}
	}
	if value, ok := decoder.GetOk("capturing_and_storage_location"); ok {
		me.CapturingAndStorageLocation = CapturingAndStorageLocation(value.(string)).Ref()
	}
	if _, ok := decoder.GetOk("scope.#"); ok {
		me.Scope = new(ScopeConditions)
		if err := me.Scope.UnmarshalHCL(hcl.NewDecoder(decoder, "scope", 0)); err != nil {
			return err
		}
	}
	if value, ok := decoder.GetOk("parameter_name"); ok {
		me.ParameterName = opt.NewString(value.(string))
	}
	if _, ok := decoder.GetOk("iib_method_node_condition.#"); ok {
		me.IIBMethodNodeCondition = new(ValueCondition)
		if err := me.IIBMethodNodeCondition.UnmarshalHCL(hcl.NewDecoder(decoder, "iib_method_node_condition", 0)); err != nil {
			return err
		}
	}
	if result, ok := decoder.GetOk("methods.#"); ok {
		me.Methods = []*CapturedMethod{}
		for idx := 0; idx < result.(int); idx++ {
			entry := new(CapturedMethod)
			if err := entry.UnmarshalHCL(hcl.NewDecoder(decoder, "methods", idx)); err != nil {
				return err
			}
			me.Methods = append(me.Methods, entry)
		}
	}
	if value, ok := decoder.GetOk("session_attribute_technology"); ok {
		me.SessionAttributeTechnology = SessionAttributeTechnology(value.(string)).Ref()
	}
	if value, ok := decoder.GetOk("technology"); ok {
		me.Technology = Technology(value.(string)).Ref()
	}
	if _, ok := decoder.GetOk("value_processing.#"); ok {
		me.ValueProcessing = new(ValueProcessing)
		if err := me.ValueProcessing.UnmarshalHCL(hcl.NewDecoder(decoder, "value_processing", 0)); err != nil {
			return err
		}
	}
	if _, ok := decoder.GetOk("cics_sdk_method_node_condition.#"); ok {
		me.CICSSDKMethodNodeCondition = new(ValueCondition)
		if err := me.CICSSDKMethodNodeCondition.UnmarshalHCL(hcl.NewDecoder(decoder, "cics_sdk_method_node_condition", 0)); err != nil {
			return err
		}
	}
	if value, ok := decoder.GetOk("enabled"); ok {
		me.Enabled = value.(bool)
	}
	if value, ok := decoder.GetOk("source"); ok {
		me.Source = Source(value.(string))
	}
	if _, ok := decoder.GetOk("iib_label_method_node_condition.#"); ok {
		me.IIBLabelMethodNodeCondition = new(ValueCondition)
		if err := me.IIBLabelMethodNodeCondition.UnmarshalHCL(hcl.NewDecoder(decoder, "iib_label_method_node_condition", 0)); err != nil {
			return err
		}
	}
	if value, ok := decoder.GetOk("iib_node_type"); ok {
		me.IIBNodeType = IIBNodeType(value.(string)).Ref()
	}
	return nil
}

func (me *DataSource) MarshalJSON() ([]byte, error) {
	m := xjson.NewProperties(me.Unknowns)
	if err := m.Marshal("capturingAndStorageLocation", me.CapturingAndStorageLocation); err != nil {
		return nil, err
	}
	if err := m.Marshal("scope", me.Scope); err != nil {
		return nil, err
	}
	if err := m.Marshal("parameterName", me.ParameterName); err != nil {
		return nil, err
	}
	if err := m.Marshal("iibMethodNodeCondition", me.IIBMethodNodeCondition); err != nil {
		return nil, err
	}
	if err := m.Marshal("methods", me.Methods); err != nil {
		return nil, err
	}
	if err := m.Marshal("sessionAttributeTechnology", me.SessionAttributeTechnology); err != nil {
		return nil, err
	}
	if err := m.Marshal("technology", me.Technology); err != nil {
		return nil, err
	}
	if err := m.Marshal("valueProcessing", me.ValueProcessing); err != nil {
		return nil, err
	}
	if err := m.Marshal("cicsSDKMethodNodeCondition", me.CICSSDKMethodNodeCondition); err != nil {
		return nil, err
	}
	if err := m.Marshal("enabled", me.Enabled); err != nil {
		return nil, err
	}
	if err := m.Marshal("source", me.Source); err != nil {
		return nil, err
	}
	if err := m.Marshal("iibLabelMethodNodeCondition", me.IIBLabelMethodNodeCondition); err != nil {
		return nil, err
	}
	if err := m.Marshal("iibNodeType", me.IIBNodeType); err != nil {
		return nil, err
	}
	return json.Marshal(m)
}

func (me *DataSource) UnmarshalJSON(data []byte) error {
	m := xjson.Properties{}
	if err := json.Unmarshal(data, &m); err != nil {
		return err
	}
	if err := m.Unmarshal("capturingAndStorageLocation", &me.CapturingAndStorageLocation); err != nil {
		return err
	}
	if err := m.Unmarshal("scope", &me.Scope); err != nil {
		return err
	}
	if err := m.Unmarshal("parameterName", &me.ParameterName); err != nil {
		return err
	}
	if err := m.Unmarshal("iibMethodNodeCondition", &me.IIBMethodNodeCondition); err != nil {
		return err
	}
	if err := m.Unmarshal("methods", &me.Methods); err != nil {
		return err
	}
	if err := m.Unmarshal("sessionAttributeTechnology", &me.SessionAttributeTechnology); err != nil {
		return err
	}
	if err := m.Unmarshal("technology", &me.Technology); err != nil {
		return err
	}
	if err := m.Unmarshal("enabled", &me.Enabled); err != nil {
		return err
	}
	if err := m.Unmarshal("valueProcessing", &me.ValueProcessing); err != nil {
		return err
	}
	if err := m.Unmarshal("cicsSDKMethodNodeCondition", &me.CICSSDKMethodNodeCondition); err != nil {
		return err
	}
	if err := m.Unmarshal("source", &me.Source); err != nil {
		return err
	}
	if err := m.Unmarshal("iibLabelMethodNodeCondition", &me.IIBLabelMethodNodeCondition); err != nil {
		return err
	}
	if err := m.Unmarshal("iibNodeType", &me.IIBNodeType); err != nil {
		return err
	}

	if len(m) > 0 {
		me.Unknowns = m
	}
	return nil
}
