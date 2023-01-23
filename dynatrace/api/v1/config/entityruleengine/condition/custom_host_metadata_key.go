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

// CustomHostMetadataKey The key of the attribute, which need dynamic keys.
// Not applicable otherwise, as the attibute itself acts as a key.
type CustomHostMetadataKey struct {
	Key      string                      `json:"key"`    // The actual key of the custom metadata.
	Source   CustomHostMetadataKeySource `json:"source"` // The source of the custom metadata.
	Unknowns map[string]json.RawMessage  `json:"-"`
}

func (chmk *CustomHostMetadataKey) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"key": {
			Type:        schema.TypeString,
			Description: "The actual key of the custom metadata",
			Required:    true,
		},
		"source": {
			Type:        schema.TypeString,
			Description: "The source of the custom metadata. Possible values are ENVIRONMENT, GOOGLE_COMPUTE_ENGINE and PLUGIN",
			Required:    true,
		},
		"unknowns": {
			Type:        schema.TypeString,
			Description: "Any attributes that aren't yet supported by this provider but have meanwhile gotten introduced by a newer version of the Dynatrace REST API",
			Optional:    true,
		},
	}
}

func (chmk *CustomHostMetadataKey) MarshalHCL(properties hcl.Properties) error {
	if err := properties.Unknowns(chmk.Unknowns); err != nil {
		return err
	}
	if err := properties.Encode("key", chmk.Key); err != nil {
		return err
	}
	if err := properties.Encode("source", string(chmk.Source)); err != nil {
		return err
	}
	return nil
}

func (chmk *CustomHostMetadataKey) UnmarshalHCL(decoder hcl.Decoder) error {
	if value, ok := decoder.GetOk("unknowns"); ok {
		if err := json.Unmarshal([]byte(value.(string)), chmk); err != nil {
			return err
		}
		if err := json.Unmarshal([]byte(value.(string)), &chmk.Unknowns); err != nil {
			return err
		}
		delete(chmk.Unknowns, "key")
		delete(chmk.Unknowns, "source")
		if len(chmk.Unknowns) == 0 {
			chmk.Unknowns = nil
		}
	}
	if value, ok := decoder.GetOk("key"); ok {
		chmk.Key = value.(string)
	}
	if value, ok := decoder.GetOk("source"); ok {
		chmk.Source = CustomHostMetadataKeySource(value.(string))
	}
	return nil
}

func (chmk *CustomHostMetadataKey) MarshalJSON() ([]byte, error) {
	m := map[string]json.RawMessage{}
	if len(chmk.Unknowns) > 0 {
		for k, v := range chmk.Unknowns {
			m[k] = v
		}
	}
	{
		rawMessage, err := json.Marshal(chmk.Key)
		if err != nil {
			return nil, err
		}
		m["key"] = rawMessage
	}
	{
		rawMessage, err := json.Marshal(chmk.Source)
		if err != nil {
			return nil, err
		}
		m["source"] = rawMessage
	}
	return json.Marshal(m)
}

func (chmk *CustomHostMetadataKey) UnmarshalJSON(data []byte) error {
	m := map[string]json.RawMessage{}
	if err := json.Unmarshal(data, &m); err != nil {
		return err
	}
	if v, found := m["key"]; found {
		if err := json.Unmarshal(v, &chmk.Key); err != nil {
			return err
		}
	}
	if v, found := m["source"]; found {
		if err := json.Unmarshal(v, &chmk.Source); err != nil {
			return err
		}
	}
	delete(m, "key")
	delete(m, "source")
	if len(m) > 0 {
		chmk.Unknowns = m
	}
	return nil
}

// CustomHostMetadataKeySource The source of the custom metadata.
type CustomHostMetadataKeySource string

// CustomHostMetadataKeySources offers the known enum values
var CustomHostMetadataKeySources = struct {
	Environment         CustomHostMetadataKeySource
	GoogleComputeEngine CustomHostMetadataKeySource
	Plugin              CustomHostMetadataKeySource
}{
	"ENVIRONMENT",
	"GOOGLE_COMPUTE_ENGINE",
	"PLUGIN",
}
