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

// RequestAttribute has no documentation
type RequestAttribute struct {
	Name                    string                     `json:"name"`                    // The name of the request attribute.
	SkipPersonalDataMasking *bool                      `json:"skipPersonalDataMasking"` // Personal data masking flag. Set `true` to skip masking.   Warning: This will potentially access personalized data.
	Confidential            *bool                      `json:"confidential"`            // Confidential data flag. Set `true` to treat the captured data as confidential.
	DataSources             []*DataSource              `json:"dataSources"`             // The list of data sources.
	DataType                DataType                   `json:"dataType"`                // The data type of the request attribute.
	Normalization           Normalization              `json:"normalization"`           // String values transformation.   If the **dataType** is not `string`, set the `Original` here.
	Enabled                 *bool                      `json:"enabled"`                 // The request attribute is enabled (`true`) or disabled (`false`).
	Aggregation             Aggregation                `json:"aggregation"`             // Aggregation type for the request values.
	Unknowns                map[string]json.RawMessage `json:"-"`
}

func (me *RequestAttribute) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:        schema.TypeString,
			Description: "The name of the request attribute",
			Required:    true,
		},
		"skip_personal_data_masking": {
			Type:        schema.TypeBool,
			Description: "Personal data masking flag. Set `true` to skip masking.   Warning: This will potentially access personalized data",
			Optional:    true,
		},
		"confidential": {
			Type:        schema.TypeBool,
			Description: "Confidential data flag. Set `true` to treat the captured data as confidential",
			Optional:    true,
		},
		"data_sources": {
			Type:        schema.TypeList,
			Optional:    true,
			MinItems:    1,
			Description: "The list of data sources",
			Elem: &schema.Resource{
				Schema: new(DataSource).Schema(),
			},
		},
		"data_type": {
			Type:        schema.TypeString,
			Description: "The data type of the request attribute",
			Required:    true,
		},
		"normalization": {
			Type:        schema.TypeString,
			Description: "String values transformation.   If the **dataType** is not `string`, set the `Original` here",
			Required:    true,
		},
		"enabled": {
			Type:        schema.TypeBool,
			Description: "The request attribute is enabled (`true`) or disabled (`false`)",
			Optional:    true,
		},
		"aggregation": {
			Type:        schema.TypeString,
			Description: "Aggregation type for the request values",
			Required:    true,
		},
		"unknowns": {
			Type:        schema.TypeString,
			Description: "allows for configuring properties that are not explicitly supported by the current version of this provider",
			Optional:    true,
		},
	}
}

func (me *RequestAttribute) MarshalHCL(properties hcl.Properties) error {
	if me.Unknowns != nil {
		delete(me.Unknowns, "metadata")
		delete(me.Unknowns, "id")
	}

	if err := properties.Unknowns(me.Unknowns); err != nil {
		return err
	}
	if err := properties.Encode("name", me.Name); err != nil {
		return err
	}
	if err := properties.Encode("skip_personal_data_masking", opt.Bool(me.SkipPersonalDataMasking)); err != nil {
		return err
	}
	if err := properties.Encode("confidential", opt.Bool(me.Confidential)); err != nil {
		return err
	}
	if err := properties.Encode("data_sources", me.DataSources); err != nil {
		return err
	}
	if err := properties.Encode("data_type", string(me.DataType)); err != nil {
		return err
	}
	if err := properties.Encode("normalization", string(me.Normalization)); err != nil {
		return err
	}
	if err := properties.Encode("enabled", opt.Bool(me.Enabled)); err != nil {
		return err
	}
	if err := properties.Encode("aggregation", string(me.Aggregation)); err != nil {
		return err
	}
	return nil
}

func (me *RequestAttribute) UnmarshalHCL(decoder hcl.Decoder) error {
	if value, ok := decoder.GetOk("unknowns"); ok {
		if err := json.Unmarshal([]byte(value.(string)), me); err != nil {
			return err
		}
		if err := json.Unmarshal([]byte(value.(string)), &me.Unknowns); err != nil {
			return err
		}
		delete(me.Unknowns, "name")
		delete(me.Unknowns, "skip_personal_data_masking")
		delete(me.Unknowns, "confidential")
		delete(me.Unknowns, "data_sources")
		delete(me.Unknowns, "data_type")
		delete(me.Unknowns, "normalization")
		delete(me.Unknowns, "enabled")
		delete(me.Unknowns, "aggregation")
		delete(me.Unknowns, "metadata")
		delete(me.Unknowns, "id")
		if len(me.Unknowns) == 0 {
			me.Unknowns = nil
		}
	}
	if value, ok := decoder.GetOk("name"); ok {
		me.Name = value.(string)
	}
	if value, ok := decoder.GetOk("skip_personal_data_masking"); ok {
		me.SkipPersonalDataMasking = opt.NewBool(value.(bool))
	}
	if value, ok := decoder.GetOk("confidential"); ok {
		me.Confidential = opt.NewBool(value.(bool))
	}
	if result, ok := decoder.GetOk("data_sources.#"); ok {
		me.DataSources = []*DataSource{}
		for idx := 0; idx < result.(int); idx++ {
			entry := new(DataSource)
			if err := entry.UnmarshalHCL(hcl.NewDecoder(decoder, "data_sources", idx)); err != nil {
				return err
			}
			me.DataSources = append(me.DataSources, entry)
		}
	}
	if value, ok := decoder.GetOk("data_type"); ok {
		me.DataType = DataType(value.(string))
	}
	if value, ok := decoder.GetOk("normalization"); ok {
		me.Normalization = Normalization(value.(string))
	}
	if value, ok := decoder.GetOk("enabled"); ok {
		me.Enabled = opt.NewBool(value.(bool))
	}
	if value, ok := decoder.GetOk("aggregation"); ok {
		me.Aggregation = Aggregation(value.(string))
	}
	return nil
}

func (me *RequestAttribute) MarshalJSON() ([]byte, error) {
	m := xjson.NewProperties(me.Unknowns)
	if err := m.Marshal("name", me.Name); err != nil {
		return nil, err
	}
	if err := m.Marshal("skipPersonalDataMasking", opt.Bool(me.SkipPersonalDataMasking)); err != nil {
		return nil, err
	}
	if err := m.Marshal("confidential", opt.Bool(me.Confidential)); err != nil {
		return nil, err
	}
	if err := m.Marshal("dataSources", me.DataSources); err != nil {
		return nil, err
	}
	if err := m.Marshal("dataType", me.DataType); err != nil {
		return nil, err
	}
	if err := m.Marshal("normalization", me.Normalization); err != nil {
		return nil, err
	}
	if err := m.Marshal("enabled", me.Enabled); err != nil {
		return nil, err
	}
	if err := m.Marshal("aggregation", me.Aggregation); err != nil {
		return nil, err
	}
	delete(m, "id")
	delete(m, "metadata")

	return json.Marshal(m)
}

func (me *RequestAttribute) UnmarshalJSON(data []byte) error {
	m := xjson.Properties{}
	if err := json.Unmarshal(data, &m); err != nil {
		return err
	}
	if err := m.Unmarshal("name", &me.Name); err != nil {
		return err
	}
	if err := m.Unmarshal("skipPersonalDataMasking", &me.SkipPersonalDataMasking); err != nil {
		return err
	}
	if err := m.Unmarshal("confidential", &me.Confidential); err != nil {
		return err
	}
	if err := m.Unmarshal("dataSources", &me.DataSources); err != nil {
		return err
	}
	if err := m.Unmarshal("dataType", &me.DataType); err != nil {
		return err
	}
	if err := m.Unmarshal("normalization", &me.Normalization); err != nil {
		return err
	}
	if err := m.Unmarshal("enabled", &me.Enabled); err != nil {
		return err
	}
	if err := m.Unmarshal("aggregation", &me.Aggregation); err != nil {
		return err
	}
	delete(m, "id")
	delete(m, "metadata")

	if len(m) > 0 {
		me.Unknowns = m
	}
	return nil
}
