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

package tag

import (
	"encoding/json"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/opt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Info Tag of a Dynatrace entity.
type Info struct {
	Context  Context                    `json:"context"`         // The origin of the tag, such as AWS or Cloud Foundry. Custom tags use the `CONTEXTLESS` value.
	Key      string                     `json:"key"`             // The key of the tag. Custom tags have the tag value here.
	Value    *string                    `json:"value,omitempty"` // The value of the tag. Not applicable to custom tags.
	Unknowns map[string]json.RawMessage `json:"-"`
}

func (ti *Info) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"context": {
			Type:        schema.TypeString,
			Description: "The origin of the tag, such as AWS or Cloud Foundry. Possible values are AWS, AWS_GENERIC, AZURE, CLOUD_FOUNDRY, CONTEXTLESS, ENVIRONMENT, GOOGLE_CLOUD and KUBERNETES. Custom tags use the `CONTEXTLESS` value",
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
			Description: "Any attributes that aren't yet supported by this provider but have meanwhile gotten introduced by a newer version of the Dynatrace REST API",
			Optional:    true,
		},
	}
}

func (ti *Info) MarshalHCL(properties hcl.Properties) error {
	if err := properties.Unknowns(ti.Unknowns); err != nil {
		return err
	}
	if err := properties.Encode("context", string(ti.Context)); err != nil {
		return err
	}
	if err := properties.Encode("key", ti.Key); err != nil {
		return err
	}
	if err := properties.Encode("value", ti.Value); err != nil {
		return err
	}
	return nil
}

func (ti *Info) UnmarshalHCL(decoder hcl.Decoder) error {
	if value, ok := decoder.GetOk("unknowns"); ok {
		if err := json.Unmarshal([]byte(value.(string)), ti); err != nil {
			return err
		}
		if err := json.Unmarshal([]byte(value.(string)), &ti.Unknowns); err != nil {
			return err
		}
		delete(ti.Unknowns, "context")
		delete(ti.Unknowns, "key")
		delete(ti.Unknowns, "value")
		if len(ti.Unknowns) == 0 {
			ti.Unknowns = nil
		}
	}
	if value, ok := decoder.GetOk("context"); ok {
		ti.Context = Context(value.(string))
	}
	if value, ok := decoder.GetOk("key"); ok {
		ti.Key = value.(string)
	}
	if value, ok := decoder.GetOk("value"); ok {
		ti.Value = opt.NewString(value.(string))
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
