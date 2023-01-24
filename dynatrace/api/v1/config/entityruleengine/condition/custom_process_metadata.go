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

package condition

import (
	"encoding/json"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// CustomProcessMetadata The key for dynamic attributes of the `PROCESS_CUSTOM_METADATA_KEY` type.
type CustomProcessMetadata struct {
	BaseConditionKey
	DynamicKey *CustomProcessMetadataKey  `json:"dynamicKey"` // The key of the attribute, which need dynamic keys.  Not applicable otherwise, as the attibute itself acts as a key.
	Unknowns   map[string]json.RawMessage `json:"-"`
}

func (cpmck *CustomProcessMetadata) GetType() *ConditionKeyType {
	return &ConditionKeyTypes.ProcessCustomMetadataKey
}

func (cpmck *CustomProcessMetadata) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"attribute": {
			Type:        schema.TypeString,
			Description: "The attribute to be used for comparision",
			Required:    true,
		},
		"type": {
			Type:        schema.TypeString,
			Description: "if specified, needs to be PROCESS_CUSTOM_METADATA_KEY",
			Optional:    true,
			Deprecated:  "The value of the attribute type is implicit, therefore shouldn't get specified",
		},
		"dynamic_key": {
			Type:        schema.TypeList,
			MaxItems:    1,
			Description: "The key of the attribute, which need dynamic keys. Not applicable otherwise, as the attibute itself acts as a key",
			Required:    true,
			Elem: &schema.Resource{
				Schema: new(CustomProcessMetadataKey).Schema(),
			},
		},
		"unknowns": {
			Type:        schema.TypeString,
			Description: "Any attributes that aren't yet supported by this provider but have meanwhile gotten introduced by a newer version of the Dynatrace REST API",
			Optional:    true,
		},
	}
}

func (cpmck *CustomProcessMetadata) MarshalHCL(properties hcl.Properties) error {
	if err := properties.Unknowns(cpmck.Unknowns); err != nil {
		return err
	}
	if err := properties.Encode("attribute", string(cpmck.Attribute)); err != nil {
		return err
	}
	if err := properties.Encode("dynamic_key", cpmck.DynamicKey); err != nil {
		return err
	}
	return nil
}

func (cpmck *CustomProcessMetadata) UnmarshalHCL(decoder hcl.Decoder) error {
	if value, ok := decoder.GetOk("unknowns"); ok {
		if err := json.Unmarshal([]byte(value.(string)), cpmck); err != nil {
			return err
		}
		if err := json.Unmarshal([]byte(value.(string)), &cpmck.Unknowns); err != nil {
			return err
		}
		delete(cpmck.Unknowns, "attribute")
		delete(cpmck.Unknowns, "dynamic_key")
		delete(cpmck.Unknowns, "type")
		if len(cpmck.Unknowns) == 0 {
			cpmck.Unknowns = nil
		}
	}
	if value, ok := decoder.GetOk("attribute"); ok {
		cpmck.Attribute = Attribute(value.(string))
	}
	if _, ok := decoder.GetOk("dynamic_key.#"); ok {
		cpmck.DynamicKey = new(CustomProcessMetadataKey)
		if err := cpmck.DynamicKey.UnmarshalHCL(hcl.NewDecoder(decoder, "dynamic_key", 0)); err != nil {
			return err
		}
	}
	return nil
}

func (cpmck *CustomProcessMetadata) MarshalJSON() ([]byte, error) {
	m := map[string]json.RawMessage{}
	if len(cpmck.Unknowns) > 0 {
		for k, v := range cpmck.Unknowns {
			m[k] = v
		}
	}
	{
		rawMessage, err := json.Marshal(cpmck.Attribute)
		if err != nil {
			return nil, err
		}
		m["attribute"] = rawMessage
	}
	if cpmck.GetType() != nil {
		rawMessage, err := json.Marshal(ConditionKeyTypes.ProcessCustomMetadataKey)
		if err != nil {
			return nil, err
		}
		m["type"] = rawMessage
	}
	if cpmck.DynamicKey != nil {
		rawMessage, err := json.Marshal(cpmck.DynamicKey)
		if err != nil {
			return nil, err
		}
		m["dynamicKey"] = rawMessage
	}
	return json.Marshal(m)
}

func (cpmck *CustomProcessMetadata) UnmarshalJSON(data []byte) error {
	m := map[string]json.RawMessage{}
	if err := json.Unmarshal(data, &m); err != nil {
		return err
	}
	cpmck.Type = cpmck.GetType()
	if v, found := m["attribute"]; found {
		if err := json.Unmarshal(v, &cpmck.Attribute); err != nil {
			return err
		}
	}
	if v, found := m["dynamicKey"]; found {
		if err := json.Unmarshal(v, &cpmck.DynamicKey); err != nil {
			return err
		}
	}
	delete(m, "attribute")
	delete(m, "dynamicKey")
	delete(m, "type")
	if len(m) > 0 {
		cpmck.Unknowns = m
	}
	return nil
}
