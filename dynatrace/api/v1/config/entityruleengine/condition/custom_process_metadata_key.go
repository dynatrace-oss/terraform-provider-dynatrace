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

// CustomProcessMetadataKey The key of the attribute, which need dynamic keys.
// Not applicable otherwise, as the attibute itself acts as a key.
type CustomProcessMetadataKey struct {
	Source   CustomProcessMetadataKeySource `json:"source"` // The source of the custom metadata.
	Key      string                         `json:"key"`    // The actual key of the custom metadata.
	Unknowns map[string]json.RawMessage     `json:"-"`
}

func (cpmk *CustomProcessMetadataKey) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"source": {
			Type:        schema.TypeString,
			Description: "The source of the custom metadata. Possible values are CLOUD_FOUNDRY, ENVIRONMENT, GOOGLE_CLOUD, KUBERNETES and PLUGIN",
			Required:    true,
		},
		"key": {
			Type:        schema.TypeString,
			Description: " The actual key of the custom metadata",
			Required:    true,
		},
		"unknowns": {
			Type:        schema.TypeString,
			Description: "Any attributes that aren't yet supported by this provider but have meanwhile gotten introduced by a newer version of the Dynatrace REST API",
			Optional:    true,
		},
	}
}

func (cpmk *CustomProcessMetadataKey) MarshalHCL(properties hcl.Properties) error {
	if err := properties.Unknowns(cpmk.Unknowns); err != nil {
		return err
	}
	if err := properties.Encode("key", cpmk.Key); err != nil {
		return err
	}
	if err := properties.Encode("source", string(cpmk.Source)); err != nil {
		return err
	}
	return nil
}

func (cpmk *CustomProcessMetadataKey) UnmarshalHCL(decoder hcl.Decoder) error {
	if value, ok := decoder.GetOk("unknowns"); ok {
		if err := json.Unmarshal([]byte(value.(string)), cpmk); err != nil {
			return err
		}
		if err := json.Unmarshal([]byte(value.(string)), &cpmk.Unknowns); err != nil {
			return err
		}
		delete(cpmk.Unknowns, "key")
		delete(cpmk.Unknowns, "source")
		if len(cpmk.Unknowns) == 0 {
			cpmk.Unknowns = nil
		}
	}
	if value, ok := decoder.GetOk("key"); ok {
		cpmk.Key = value.(string)
	}
	if value, ok := decoder.GetOk("source"); ok {
		cpmk.Source = CustomProcessMetadataKeySource(value.(string))
	}
	return nil
}

func (cpmk *CustomProcessMetadataKey) MarshalJSON() ([]byte, error) {
	m := map[string]json.RawMessage{}
	if len(cpmk.Unknowns) > 0 {
		for k, v := range cpmk.Unknowns {
			m[k] = v
		}
	}
	{
		rawMessage, err := json.Marshal(cpmk.Key)
		if err != nil {
			return nil, err
		}
		m["key"] = rawMessage
	}
	{
		rawMessage, err := json.Marshal(cpmk.Source)
		if err != nil {
			return nil, err
		}
		m["source"] = rawMessage
	}
	return json.Marshal(m)
}

func (cpmk *CustomProcessMetadataKey) UnmarshalJSON(data []byte) error {
	m := map[string]json.RawMessage{}
	if err := json.Unmarshal(data, &m); err != nil {
		return err
	}
	if v, found := m["key"]; found {
		if err := json.Unmarshal(v, &cpmk.Key); err != nil {
			return err
		}
	}
	if v, found := m["source"]; found {
		if err := json.Unmarshal(v, &cpmk.Source); err != nil {
			return err
		}
	}
	delete(m, "key")
	delete(m, "source")
	if len(m) > 0 {
		cpmk.Unknowns = m
	}
	return nil
}

// CustomProcessMetadataKeySource The source of the custom metadata.
type CustomProcessMetadataKeySource string

// CustomProcessMetadataKeySources offers the known enum values
var CustomProcessMetadataKeySources = struct {
	CloudFoundry CustomProcessMetadataKeySource
	Environment  CustomProcessMetadataKeySource
	GoogleCloud  CustomProcessMetadataKeySource
	Kubernetes   CustomProcessMetadataKeySource
	Plugin       CustomProcessMetadataKeySource
}{
	"CLOUD_FOUNDRY",
	"ENVIRONMENT",
	"GOOGLE_CLOUD",
	"KUBERNETES",
	"PLUGIN",
}
