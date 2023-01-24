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

package kubernetes

import (
	"encoding/json"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// KubernetesEventPattern Represents a single Kubernetes events field selector (=event filter based on the K8s field selector).
type KubernetesEventPattern struct {
	Active        bool                       `json:"active"`        // Whether subscription to this events field selector is enabled (value set to `true`). If disabled (value set to `false`), Dynatrace will stop fetching events from the Kubernetes API for this events field selector
	FieldSelector string                     `json:"fieldSelector"` // The field selector string (url decoding is applied) when storing it.
	Label         string                     `json:"label"`         // A label of the events field selector.
	Unknowns      map[string]json.RawMessage `json:"-"`
}

func (kep *KubernetesEventPattern) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"unknowns": {
			Type:        schema.TypeString,
			Description: "Any attributes that aren't yet supported by this provider",
			Optional:    true,
		},
		"active": {
			Type:        schema.TypeBool,
			Description: "Whether subscription to this events field selector is enabled (value set to `true`). If disabled (value set to `false`), Dynatrace will stop fetching events from the Kubernetes API for this events field selector",
			Required:    true,
		},
		"field_selector": {
			Type:        schema.TypeString,
			Description: "The field selector string (url decoding is applied) when storing it.",
			Required:    true,
		},
		"label": {
			Type:        schema.TypeString,
			Description: "A label of the events field selector.",
			Required:    true,
		},
	}
}

func (kep *KubernetesEventPattern) MarshalHCL(properties hcl.Properties) error {
	if err := properties.Unknowns(kep.Unknowns); err != nil {
		return err
	}

	if err := properties.Encode("active", kep.Active); err != nil {
		return err
	}
	if err := properties.Encode("field_selector", kep.FieldSelector); err != nil {
		return err
	}
	if err := properties.Encode("label", kep.Label); err != nil {
		return err
	}

	return nil
}

func (kep *KubernetesEventPattern) UnmarshalHCL(decoder hcl.Decoder) error {
	if value, ok := decoder.GetOk("unknowns"); ok {
		if err := json.Unmarshal([]byte(value.(string)), kep); err != nil {
			return err
		}
		if err := json.Unmarshal([]byte(value.(string)), &kep.Unknowns); err != nil {
			return err
		}
		delete(kep.Unknowns, "active")
		delete(kep.Unknowns, "field_selector")
		delete(kep.Unknowns, "label")
		if len(kep.Unknowns) == 0 {
			kep.Unknowns = nil
		}
	}
	if value, ok := decoder.GetOk("active"); ok {
		kep.Active = value.(bool)
	}
	if value, ok := decoder.GetOk("field_selector"); ok {
		kep.FieldSelector = value.(string)
	}
	if value, ok := decoder.GetOk("label"); ok {
		kep.Label = value.(string)
	}
	return nil
}

func (kep *KubernetesEventPattern) UnmarshalJSON(data []byte) error {
	m := map[string]json.RawMessage{}
	if err := json.Unmarshal(data, &m); err != nil {
		return err
	}
	if v, found := m["label"]; found {
		if err := json.Unmarshal(v, &kep.Label); err != nil {
			return err
		}
	}
	if v, found := m["active"]; found {
		if err := json.Unmarshal(v, &kep.Active); err != nil {
			return err
		}
	}
	if v, found := m["fieldSelector"]; found {
		if err := json.Unmarshal(v, &kep.FieldSelector); err != nil {
			return err
		}
	}
	delete(m, "active")
	delete(m, "label")
	delete(m, "fieldSelector")
	if len(m) > 0 {
		kep.Unknowns = m
	}
	return nil
}

func (kep *KubernetesEventPattern) MarshalJSON() ([]byte, error) {
	m := map[string]json.RawMessage{}
	if len(kep.Unknowns) > 0 {
		for k, v := range kep.Unknowns {
			m[k] = v
		}
	}
	{
		rawMessage, err := json.Marshal(kep.Label)
		if err != nil {
			return nil, err
		}
		m["label"] = rawMessage
	}
	{
		rawMessage, err := json.Marshal(kep.Active)
		if err != nil {
			return nil, err
		}
		m["active"] = rawMessage
	}
	{
		rawMessage, err := json.Marshal(kep.FieldSelector)
		if err != nil {
			return nil, err
		}
		m["fieldSelector"] = rawMessage
	}
	return json.Marshal(m)
}
