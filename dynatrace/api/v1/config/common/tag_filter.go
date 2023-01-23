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

package common

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/opt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// TagFilter A tag-based filter of monitored entities.
type TagFilter struct {
	Context Context `json:"context"`         // The origin of the tag, such as AWS or Cloud Foundry.  Custom tags use the `CONTEXTLESS` value.
	Key     string  `json:"key"`             // The key of the tag. Custom tags have the tag value here.
	Value   *string `json:"value,omitempty"` // The value of the tag. Not applicable to custom tags.
}

func (me *TagFilter) Equals(other *TagFilter) bool {
	if other == nil {
		return false
	}
	if me.Key != other.Key {
		return false
	}
	if me.Context != other.Context {
		return false
	}
	if me.Value == nil && other.Value != nil {
		return false
	}
	if me.Value != nil && other.Value == nil {
		return false
	}
	if me.Value != nil && *me.Value != *other.Value {
		return false
	}
	return true
}

func (me *TagFilter) Schema() map[string]*schema.Schema {
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
	}
}

func (me *TagFilter) MarshalHCL(properties hcl.Properties) error {
	if err := properties.Encode("context", string(me.Context)); err != nil {
		return err
	}
	if err := properties.Encode("key", me.Key); err != nil {
		return err
	}
	if err := properties.Encode("value", me.Value); err != nil {
		return err
	}
	return nil
}

func (me *TagFilter) UnmarshalHCL(decoder hcl.Decoder) error {
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

// Context The origin of the tag, such as AWS or Cloud Foundry.
// Custom tags use the `CONTEXTLESS` value.
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
