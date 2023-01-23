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

package comparison

import (
	"encoding/json"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/entityruleengine/comparison/dcrum_decoder"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// DCRumDecoder Comparison for `DCRUM_DECODER_TYPE` attributes.
type DCRumDecoder struct {
	BaseComparison
	Operator dcrum_decoder.Operator `json:"operator"`        // Operator of the comparison. You can reverse it by setting **negate** to `true`.  Possible values depend on the **type** of the comparison. Find the list of actual models in the description of the **type** field and check the description of the model you need.
	Value    *dcrum_decoder.Value   `json:"value,omitempty"` // The value to compare to.
}

func (ddc *DCRumDecoder) GetType() ComparisonBasicType {
	return ComparisonBasicTypes.DCRumDecoderType
}

func (ddc *DCRumDecoder) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"type": {
			Type:        schema.TypeString,
			Description: "if specified, needs to be DCRUM_DECODER_TYPE",
			Optional:    true,
			Deprecated:  "The value of the attribute type is implicit, therefore shouldn't get specified",
		},
		"negate": {
			Type:        schema.TypeBool,
			Description: "Reverses the operator. For example it turns EQUALS into DOES NOT EQUAL",
			Optional:    true,
		},
		"operator": {
			Type:        schema.TypeString,
			Description: "Either EQUALS or EXISTS. You can reverse it by setting **negate** to `true`",
			Required:    true,
		},
		"value": {
			Type:        schema.TypeString,
			Description: "The value to compare to. Possible values are ALL_OTHER, CITRIX_APPFLOW, CITRIX_ICA, CITRIX_ICA_OVER_SSL, DB2_DRDA, HTTP, HTTPS, HTTP_EXPRESS, INFORMIX, MYSQL, ORACLE, SAP_GUI, SAP_GUI_OVER_HTTP, SAP_GUI_OVER_HTTPS, SAP_HANA_DB, SAP_RFC, SSL and TDS.",
			Optional:    true,
		},
		"unknowns": {
			Type:        schema.TypeString,
			Description: "Any attributes that aren't yet supported by this provider but have meanwhile gotten introduced by a newer version of the Dynatrace REST API",
			Optional:    true,
		},
	}
}

func (ddc *DCRumDecoder) MarshalHCL(properties hcl.Properties) error {
	if err := properties.Unknowns(ddc.Unknowns); err != nil {
		return err
	}
	if err := properties.Encode("negate", ddc.Negate); err != nil {
		return err
	}
	if err := properties.Encode("operator", string(ddc.Operator)); err != nil {
		return err
	}
	if err := properties.Encode("value", ddc.Value.String()); err != nil {
		return err
	}
	return nil
}

func (ddc *DCRumDecoder) UnmarshalHCL(decoder hcl.Decoder) error {
	if value, ok := decoder.GetOk("unknowns"); ok {
		if err := json.Unmarshal([]byte(value.(string)), ddc); err != nil {
			return err
		}
		if err := json.Unmarshal([]byte(value.(string)), &ddc.Unknowns); err != nil {
			return err
		}
		delete(ddc.Unknowns, "type")
		delete(ddc.Unknowns, "negate")
		delete(ddc.Unknowns, "operator")
		delete(ddc.Unknowns, "value")
		if len(ddc.Unknowns) == 0 {
			ddc.Unknowns = nil
		}
	}
	if value, ok := decoder.GetOk("type"); ok {
		ddc.Type = ComparisonBasicType(value.(string))
	}
	if value, ok := decoder.GetOk("negate"); ok {
		ddc.Negate = value.(bool)
	}
	if value, ok := decoder.GetOk("operator"); ok {
		ddc.Operator = dcrum_decoder.Operator(value.(string))
	}
	if value, ok := decoder.GetOk("value"); ok {
		ddc.Value = dcrum_decoder.Value(value.(string)).Ref()
	}
	return nil
}

func (ddc *DCRumDecoder) MarshalJSON() ([]byte, error) {
	m := map[string]json.RawMessage{}
	if len(ddc.Unknowns) > 0 {
		for k, v := range ddc.Unknowns {
			m[k] = v
		}
	}
	{
		rawMessage, err := json.Marshal(ddc.Negate)
		if err != nil {
			return nil, err
		}
		m["negate"] = rawMessage
	}
	{
		rawMessage, err := json.Marshal(ComparisonBasicTypes.DCRumDecoderType)
		if err != nil {
			return nil, err
		}
		m["type"] = rawMessage
	}
	{
		rawMessage, err := json.Marshal(&ddc.Operator)
		if err != nil {
			return nil, err
		}
		m["operator"] = rawMessage
	}
	if ddc.Value != nil {
		rawMessage, err := json.Marshal(ddc.Value)
		if err != nil {
			return nil, err
		}
		m["value"] = rawMessage
	}
	return json.Marshal(m)
}

func (ddc *DCRumDecoder) UnmarshalJSON(data []byte) error {
	m := map[string]json.RawMessage{}
	if err := json.Unmarshal(data, &m); err != nil {
		return err
	}
	ddc.Type = ddc.GetType()
	if v, found := m["negate"]; found {
		if err := json.Unmarshal(v, &ddc.Negate); err != nil {
			return err
		}
	}
	if v, found := m["operator"]; found {
		if err := json.Unmarshal(v, &ddc.Operator); err != nil {
			return err
		}
	}
	if v, found := m["value"]; found {
		if err := json.Unmarshal(v, &ddc.Value); err != nil {
			return err
		}
	}
	delete(m, "negate")
	delete(m, "operator")
	delete(m, "value")
	delete(m, "type")
	if len(m) > 0 {
		ddc.Unknowns = m
	}
	return nil
}
