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

// Key The key to identify the data we're matching.
// The actual set of fields and possible values vary, depending on the **type** of the key.
// Find the list of actual objects in the description of the **type** field.
type Key interface {
	MarshalHCL(hcl.Properties) error
	UnmarshalHCL(decoder hcl.Decoder) error
	MarshalJSON() ([]byte, error)
	UnmarshalJSON(data []byte) error
	GetAttribute() Attribute
	GetType() *ConditionKeyType
	Schema() map[string]*schema.Schema
}

// BaseConditionKey The key to identify the data we're matching.
// The actual set of fields and possible values vary, depending on the **type** of the key.
// Find the list of actual objects in the description of the **type** field.
type BaseConditionKey struct {
	Attribute Attribute                  `json:"attribute"`      // The attribute to be used for comparision.
	Type      *ConditionKeyType          `json:"type,omitempty"` // Defines the actual set of fields depending on the value. See one of the following objects:  * `PROCESS_CUSTOM_METADATA_KEY` -> CustomProcessMetadataConditionKey  * `HOST_CUSTOM_METADATA_KEY` -> CustomHostMetadataConditionKey  * `PROCESS_PREDEFINED_METADATA_KEY` -> ProcessMetadataConditionKey  * `STRING` -> StringConditionKey
	Unknowns  map[string]json.RawMessage `json:"-"`
}

func (bck *BaseConditionKey) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"attribute": {
			Type:        schema.TypeString,
			Description: "The attribute to be used for comparision",
			Required:    true,
		},
		"type": {
			Type:        schema.TypeString,
			Description: "Defines the actual set of fields depending on the value",
			Optional:    true,
		},
		"unknowns": {
			Type:        schema.TypeString,
			Description: "Any attributes that aren't yet supported by this provider but have meanwhile gotten introduced by a newer version of the Dynatrace REST API",
			Optional:    true,
		},
	}
}

func (bck *BaseConditionKey) MarshalHCL(properties hcl.Properties) error {
	if err := properties.Unknowns(bck.Unknowns); err != nil {
		return err
	}
	if err := properties.Encode("attribute", string(bck.Attribute)); err != nil {
		return err
	}
	if err := properties.Encode("type", bck.Type.String()); err != nil {
		return err
	}
	return nil
}

func (bck *BaseConditionKey) UnmarshalHCL(decoder hcl.Decoder) error {
	if value, ok := decoder.GetOk("unknowns"); ok {
		if err := json.Unmarshal([]byte(value.(string)), bck); err != nil {
			return err
		}
		if err := json.Unmarshal([]byte(value.(string)), &bck.Unknowns); err != nil {
			return err
		}
		delete(bck.Unknowns, "attribute")
		delete(bck.Unknowns, "type")
		if len(bck.Unknowns) == 0 {
			bck.Unknowns = nil
		}
	}
	if value, ok := decoder.GetOk("attribute"); ok {
		bck.Attribute = Attribute(value.(string))
	}
	if value, ok := decoder.GetOk("type"); ok {
		bck.Type = ConditionKeyType(value.(string)).Ref()
	}
	return nil
}

func (bck *BaseConditionKey) GetAttribute() Attribute {
	return bck.Attribute
}

func (bck *BaseConditionKey) GetType() *ConditionKeyType {
	return bck.Type
}

func (bck *BaseConditionKey) UnmarshalJSON(data []byte) error {
	m := map[string]json.RawMessage{}
	if err := json.Unmarshal(data, &m); err != nil {
		return err
	}
	if v, found := m["attribute"]; found {
		if err := json.Unmarshal(v, &bck.Attribute); err != nil {
			return err
		}
	}
	if v, found := m["type"]; found {
		if err := json.Unmarshal(v, &bck.Type); err != nil {
			return err
		}
	}
	delete(m, "attribute")
	delete(m, "type")
	if len(m) > 0 {
		bck.Unknowns = m
	}
	return nil
}

func (bck *BaseConditionKey) MarshalJSON() ([]byte, error) {
	m := map[string]json.RawMessage{}
	if len(bck.Unknowns) > 0 {
		for k, v := range bck.Unknowns {
			m[k] = v
		}
	}
	{
		rawMessage, err := json.Marshal(bck.Attribute)
		if err != nil {
			return nil, err
		}
		m["attribute"] = rawMessage
	}
	if bck.GetType() != nil {
		rawMessage, err := json.Marshal(*bck.Type)
		if err != nil {
			return nil, err
		}
		m["type"] = rawMessage
	}
	return json.Marshal(m)
}

// ConditionKeyType Defines the actual set of fields depending on the value. See one of the following objects:
// * `PROCESS_CUSTOM_METADATA_KEY` -> CustomProcessMetadataConditionKey
// * `HOST_CUSTOM_METADATA_KEY` -> CustomHostMetadataConditionKey
// * `PROCESS_PREDEFINED_METADATA_KEY` -> ProcessMetadataConditionKey
// * `STRING` -> StringConditionKey
type ConditionKeyType string

func (v ConditionKeyType) Ref() *ConditionKeyType {
	return &v
}

func (v *ConditionKeyType) String() string {
	return string(*v)
}

// ConditionKeyTypes offers the known enum values
var ConditionKeyTypes = struct {
	HostCustomMetadataKey        ConditionKeyType
	ProcessCustomMetadataKey     ConditionKeyType
	ProcessPredefinedMetadataKey ConditionKeyType
	String                       ConditionKeyType
}{
	"HOST_CUSTOM_METADATA_KEY",
	"PROCESS_CUSTOM_METADATA_KEY",
	"PROCESS_PREDEFINED_METADATA_KEY",
	"STRING",
}
