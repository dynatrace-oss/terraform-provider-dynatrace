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

package entity

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Tags []*Tag // A set of tags assigned to the entity.

func (me Tags) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"tag": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "A tag assigned to the entity",
			Elem:        &schema.Resource{Schema: new(Tag).Schema()},
		},
	}
}

func (me Tags) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeSlice("tag", me)
}

func (me *Tags) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeSlice("tag", me)
}

// Tag A tag-based filter of monitored entities.
type Tag struct {
	Context              Context `json:"context"`                        // The origin of the tag, such as AWS or Cloud Foundry.  Custom tags use the `CONTEXTLESS` value.
	Key                  string  `json:"key"`                            // The key of the tag. Custom tags have the tag value here.
	Value                *string `json:"value,omitempty"`                // The value of the tag. Not applicable to custom tags.
	StringRepresentation *string `json:"stringRepresentation,omitempty"` // The string representation of the tag.
}

func (me *Tag) Schema() map[string]*schema.Schema {
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
		"string_representation": {
			Type:        schema.TypeString,
			Description: "The string representation of the tag",
			Optional:    true,
		},
	}
}

func (me *Tag) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"context":               me.Context,
		"key":                   me.Key,
		"value":                 me.Value,
		"string_representation": me.StringRepresentation,
	})
}

func (me *Tag) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"context":               &me.Context,
		"key":                   &me.Key,
		"value":                 &me.Value,
		"string_representation": &me.StringRepresentation,
	})
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
