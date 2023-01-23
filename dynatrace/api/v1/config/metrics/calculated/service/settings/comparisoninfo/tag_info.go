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

package comparisoninfo

import (
	"encoding/json"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/xjson"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// TagInfo Tag of a Dynatrace entity.
type TagInfo struct {
	Context TagInfoContext `json:"context"`         // The origin of the tag, such as AWS or Cloud Foundry.   Custom tags use the `CONTEXTLESS` value.
	Key     string         `json:"key"`             // The key of the tag.   Custom tags have the tag value here.
	Value   *string        `json:"value,omitempty"` // The value of the tag.   Not applicable to custom tags.
}

func (me *TagInfo) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"key": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The key of the tag. Custom tags have the tag value here",
		},
		"value": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The value of the tag. Not applicable to custom tags",
		},
		"context": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The origin of the tag, such as AWS or Cloud Foundry. Custom tags use the `CONTEXTLESS` value. Possible values are `AWS`, `AWS_GENERIC`, `AZURE`, `CLOUD_FOUNDRY`, `CONTEXTLESS`, `ENVIRONMENT`, `GOOGLE_CLOUD` and `KUBERNETES`",
		},
		"unknowns": {
			Type:        schema.TypeString,
			Description: "allows for configuring properties that are not explicitly supported by the current version of this provider",
			Optional:    true,
		},
	}
}

func (me *TagInfo) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"key":     me.Key,
		"value":   me.Value,
		"context": me.Context,
	})
}

func (me *TagInfo) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"key":     &me.Key,
		"value":   &me.Value,
		"context": &me.Context,
	})
}

func (me *TagInfo) MarshalJSON() ([]byte, error) {
	properties := xjson.Properties{}
	if err := properties.MarshalAll(map[string]any{
		"key":     me.Key,
		"value":   me.Value,
		"context": me.Context,
	}); err != nil {
		return nil, err
	}
	return json.Marshal(properties)
}

func (me *TagInfo) UnmarshalJSON(data []byte) error {
	properties := xjson.Properties{}
	if err := json.Unmarshal(data, &properties); err != nil {
		return err
	}
	return properties.UnmarshalAll(map[string]any{
		"key":     &me.Key,
		"value":   &me.Value,
		"context": &me.Context,
	})
}

// TagInfoContext The origin of the tag, such as AWS or Cloud Foundry.
//
//	Custom tags use the `CONTEXTLESS` value.
type TagInfoContext string

// TagInfoContexts offers the known enum values
var TagInfoContexts = struct {
	AWS          TagInfoContext
	AWSGeneric   TagInfoContext
	Azure        TagInfoContext
	CloudFoundry TagInfoContext
	Contextless  TagInfoContext
	Environment  TagInfoContext
	GoogleCloud  TagInfoContext
	Kubernetes   TagInfoContext
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

type TagInfos []*TagInfo

func (me TagInfos) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"value": {
			Type:        schema.TypeList,
			MinItems:    1,
			Optional:    true,
			Description: "The values to compare to",
			Elem:        &schema.Resource{Schema: new(TagInfo).Schema()},
		},
	}
}

func (me TagInfos) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeSlice("value", me)
}

func (me *TagInfos) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeSlice("value", me)
}
