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

package maintenance

import (
	"encoding/json"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/xjson"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/opt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type TagInfos []*TagInfo

func (me TagInfos) Equals(other any) bool {
	if tagInfos, ok := other.(TagInfos); ok {
		if len(me) != len(tagInfos) {
			return false
		}
		for _, tagInfo := range me {
			found := false
			for _, tagInfo2 := range tagInfos {
				if tagInfo.Equals(tagInfo2) {
					found = true
					break
				}
			}
			if !found {
				return false
			}
		}
		return true
	}
	return false
}

// TagInfo Tag of a Dynatrace entity.
type TagInfo struct {
	Context  Context                    `json:"context"`         // The origin of the tag, such as AWS or Cloud Foundry. Custom tags use the `CONTEXTLESS` value.
	Key      string                     `json:"key"`             // The key of the tag. Custom tags have the tag value here.
	Value    *string                    `json:"value,omitempty"` // The value of the tag. Not applicable to custom tags.
	Unknowns map[string]json.RawMessage `json:"-"`
}

func (me *TagInfo) Equals(other any) bool {
	if tagInfo, ok := other.(*TagInfo); ok {
		if me.Context != tagInfo.Context {
			return false
		}
		if me.Key != tagInfo.Key {
			return false
		}
		if me.Value == nil && tagInfo.Value != nil {
			return false
		}
		if me.Value != nil && tagInfo.Value == nil {
			return false
		}
		if me.Value != nil && tagInfo.Value != nil && *me.Value != *tagInfo.Value {
			return false
		}
		if len(me.Unknowns) != len(tagInfo.Unknowns) {
			return false
		}
		return true
	}

	return false
}

func (me *TagInfo) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"context": {
			Type:        schema.TypeString,
			Description: "The origin of the tag, such as AWS or Cloud Foundry. Custom tags use the `CONTEXTLESS` value",
			Required:    true,
		},
		"key": {
			Type:        schema.TypeString,
			Description: "The key of the tag. Custom tags have the tag value here",
			Required:    true,
		},
		"value": {
			Type:        schema.TypeString,
			Description: "The value of the tag. Not applicable to custom tags",
			Optional:    true,
		},
		"unknowns": {
			Type:        schema.TypeString,
			Description: "allows for configuring properties that are not explicitly supported by the current version of this provider",
			Optional:    true,
		},
	}
}

func (me *TagInfo) MarshalHCL(properties hcl.Properties) error {
	if err := properties.Unknowns(me.Unknowns); err != nil {
		return err
	}
	if err := properties.Encode("context", string(me.Context)); err != nil {
		return err
	}
	if err := properties.Encode("key", string(me.Key)); err != nil {
		return err
	}
	if err := properties.Encode("value", me.Value); err != nil {
		return err
	}

	return nil
}

func (me *TagInfo) UnmarshalHCL(decoder hcl.Decoder) error {
	if value, ok := decoder.GetOk("unknowns"); ok {
		if err := json.Unmarshal([]byte(value.(string)), me); err != nil {
			return err
		}
		if err := json.Unmarshal([]byte(value.(string)), &me.Unknowns); err != nil {
			return err
		}
		delete(me.Unknowns, "context")
		delete(me.Unknowns, "key")
		delete(me.Unknowns, "value")
		if len(me.Unknowns) == 0 {
			me.Unknowns = nil
		}
	}
	if value, ok := decoder.GetOk("context"); ok {
		me.Context = Context(value.(string))
	}
	if value, ok := decoder.GetOk("key"); ok {
		me.Key = value.(string)
	}
	if value, ok := decoder.GetOk("value"); ok {
		me.Value = opt.NewString(value.(string))
	}
	return nil
}

func (me *TagInfo) MarshalJSON() ([]byte, error) {
	m := xjson.NewProperties(me.Unknowns)
	if err := m.Marshal("context", me.Context); err != nil {
		return nil, err
	}
	if err := m.Marshal("key", me.Key); err != nil {
		return nil, err
	}
	if err := m.Marshal("value", me.Value); err != nil {
		return nil, err
	}
	return json.Marshal(m)
}

func (me *TagInfo) UnmarshalJSON(data []byte) error {
	m := xjson.Properties{}
	if err := json.Unmarshal(data, &m); err != nil {
		return err
	}
	if err := m.Unmarshal("context", &me.Context); err != nil {
		return err
	}
	if err := m.Unmarshal("key", &me.Key); err != nil {
		return err
	}
	if err := m.Unmarshal("value", &me.Value); err != nil {
		return err
	}

	if len(m) > 0 {
		me.Unknowns = m
	}
	return nil
}

// Context The origin of the tag, such as AWS or Cloud Foundry.
//
//	Custom tags use the `CONTEXTLESS` value.
type Context string

// Contexts offers the known enum values
var Contexts = struct {
	AWS          Context
	AWSGeneric   Context
	Azure        Context
	CloudFoundry Context
	Contextless  Context
	Environment  Context
	GoogleCloud  Context
	Kubernetes   Context
}{
	"AWS",
	"AWS_GENERIC",
	"AZURE",
	"CLOUD_FOUNDRY",
	"CONTEXTLESS",
	"ENVIRONMENT",
	"GOOGLE_CLOUD",
	"KUBERNETES",
}
