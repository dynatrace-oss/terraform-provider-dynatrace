/**
* @license
* Copyright 2025 Dynatrace LLC
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

package generic

import (
	"encoding/json"
	"regexp"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Settings struct {
	Scope        string `json:"-" scope:"scope"` // The scope of this setting
	SchemaID     string `json:"schemaId"`
	Value        string `json:"value"`
	LocalStorage string `json:"-"`
}

func (me *Settings) Name() string {
	return me.SchemaID
}

type SecretMap map[string]any
type SecretSlice []any

func (ss SecretSlice) Merge(other SecretSlice) SecretSlice {
	if len(other) == 0 {
		return other
	}
	overlap := SecretSlice{}
	for i := 0; i < len(other); i++ {
		v := other[i]
		switch tv := v.(type) {
		case map[string]any:
			if len(ss) <= i {
				overlap = append(overlap, SecretMap(tv))
			} else {
				SecretMap(ss[i].(map[string]any)).Merge(SecretMap(tv))
			}
		case []any:
			if len(ss) <= i {
				overlap = append(overlap, SecretSlice(tv))
			} else {
				ss[i] = SecretSlice(ss[i].([]any)).Merge(SecretSlice(tv))
			}
		}
	}
	if len(ss) > len(other) {
		return ss[0:len(other)]
	} else if len(ss) < len(other) {
		return append(ss, overlap...)
	}
	return ss
}

var secretReg = regexp.MustCompile(`^\*\*\*\d\d\d\*\*\*$`)

func (ss SecretSlice) stripSecrets() {
	if len(ss) == 0 {
		return
	}
	for _, v := range ss {
		switch tv := v.(type) {
		case map[string]any:
			SecretMap(tv).stripSecrets()
		case []any:
			SecretSlice(tv).stripSecrets()
		default:
		}
	}
}

func (sm SecretMap) Merge(other SecretMap) {
	if other == nil {
		return
	}
	for k, v := range other {
		var value any
		if secretValue, found := sm["dynatrace_secret_"+k]; found {
			value = secretValue
		} else if nonSecretValue, found := sm[k]; found {
			value = nonSecretValue
		}
		if value != nil {
			switch tv := v.(type) {
			case map[string]any:
				SecretMap(value.(map[string]any)).Merge(SecretMap(tv))
			case []any:
				sm[k] = SecretSlice(value.([]any)).Merge(SecretSlice(tv))
			default:
				sm[k] = v
			}
		} else {
			sm[k] = v
		}

	}
}

func (sm SecretMap) stripSecrets() {
	if len(sm) == 0 {
		return
	}

	for k, v := range sm {
		switch tv := v.(type) {
		case map[string]any:
			SecretMap(tv).stripSecrets()
		case []any:
			SecretSlice(tv).stripSecrets()
		case string:
			if secretReg.MatchString(tv) {
				delete(sm, k)
			}
		default:
		}
	}
}

func (me *Settings) Merge(other *Settings) {
	if other == nil {
		return
	}
	base := SecretMap{}
	if err := json.Unmarshal([]byte(me.Value), &base); err != nil {
		return
	}
	otherm := SecretMap{}
	if err := json.Unmarshal([]byte(other.Value), &otherm); err != nil {
		return
	}
	otherm.stripSecrets()
	base.Merge(otherm)
	data, err := json.Marshal(base)
	if err != nil {
		return
	}
	me.Value = string(data)
}

func (me *Settings) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"value": {
			Type:             schema.TypeString,
			Required:         true,
			DiffSuppressFunc: hcl.SuppressJSONorEOT,
		},
		"schema": {
			Type:     schema.TypeString,
			Required: true,
		},
		"scope": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"local_storage": {
			Type:     schema.TypeString,
			Computed: true,
		},
	}
}

func (me *Settings) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"scope":         me.Scope,
		"value":         me.Value,
		"schema":        me.SchemaID,
		"local_storage": me.LocalStorage,
	})
}

func (me *Settings) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"scope":         &me.Scope,
		"value":         &me.Value,
		"schema":        &me.SchemaID,
		"local_storage": &me.LocalStorage,
	})
}
