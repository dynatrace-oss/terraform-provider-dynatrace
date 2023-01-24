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

package propagation

import (
	"encoding/json"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/xjson"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// UniversalTagKey has no documentation
type UniversalTagKey struct {
	Key     *string                 `json:"key,omitempty"`     // has no documentation
	Context *UniversalTagKeyContext `json:"context,omitempty"` // has no documentation
}

func (me *UniversalTagKey) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"key": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "has no documentation",
		},
		"context": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "has no documentation",
		},
	}
}

func (me *UniversalTagKey) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"key":     me.Key,
		"context": me.Context,
	})
}

func (me *UniversalTagKey) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"key":     &me.Key,
		"context": &me.Context,
	})
}

func (me *UniversalTagKey) MarshalJSON() ([]byte, error) {
	properties := xjson.Properties{}
	if err := properties.MarshalAll(map[string]any{
		"key":     me.Key,
		"context": me.Context,
	}); err != nil {
		return nil, err
	}
	return json.Marshal(properties)
}

func (me *UniversalTagKey) UnmarshalJSON(data []byte) error {
	properties := xjson.Properties{}
	if err := json.Unmarshal(data, &properties); err != nil {
		return err
	}
	return properties.UnmarshalAll(map[string]any{
		"key":     me.Key,
		"context": me.Context,
	})
}

// UniversalTagKeyContext has no documentation
type UniversalTagKeyContext string

// UniversalTagKeyContexts offers the known enum values
var UniversalTagKeyContexts = struct {
	AWS                 UniversalTagKeyContext
	AWSGeneric          UniversalTagKeyContext
	Azure               UniversalTagKeyContext
	CloudFoundry        UniversalTagKeyContext
	Contextless         UniversalTagKeyContext
	Environment         UniversalTagKeyContext
	GoogleComputeEngine UniversalTagKeyContext
	Kubernetes          UniversalTagKeyContext
}{
	"AWS",
	"AWS_GENERIC",
	"AZURE",
	"CLOUD_FOUNDRY",
	"CONTEXTLESS",
	"ENVIRONMENT",
	"GOOGLE_COMPUTE_ENGINE",
	"KUBERNETES",
}
