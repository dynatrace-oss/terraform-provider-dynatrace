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

package outgoing

import (
	"fmt"
	"slices"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type EventDataAttributeComplex struct {
	Path       *string                        `json:"path,omitempty"`   // [See our documentation](https://dt-url.net/ei034bx)
	Source     *string                        `json:"source,omitempty"` // Fixed value
	SourceType DataSourceWithStaticStringEnum `json:"sourceType"`       // Data source. Possible Values: `constant.string`, `request.body`, `request.headers`, `request.method`, `request.parameters`, `request.path`, `request.url`, `response.body`, `response.headers`, `response.statusCode`
}

func (me *EventDataAttributeComplex) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"path": {
			Type:        schema.TypeString,
			Description: "[See our documentation](https://dt-url.net/ei034bx)",
			Optional:    true, // precondition
		},
		"source": {
			Type:        schema.TypeString,
			Description: "Fixed value",
			Optional:    true, // precondition
		},
		"source_type": {
			Type:        schema.TypeString,
			Description: "Data source. Possible Values: `constant.string`, `request.body`, `request.headers`, `request.method`, `request.parameters`, `request.path`, `request.url`, `response.body`, `response.headers`, `response.statusCode`",
			Required:    true,
		},
	}
}

func (me *EventDataAttributeComplex) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"path":        me.Path,
		"source":      me.Source,
		"source_type": me.SourceType,
	})
}

func (me *EventDataAttributeComplex) HandlePreconditions() error {
	if (me.Path == nil) && (slices.Contains([]string{"request.body", "request.headers", "request.parameters", "response.body", "response.headers"}, string(me.SourceType))) {
		return fmt.Errorf("'path' must be specified if 'source_type' is set to '%v'", me.SourceType)
	}
	if (me.Path != nil) && (!slices.Contains([]string{"request.body", "request.headers", "request.parameters", "response.body", "response.headers"}, string(me.SourceType))) {
		return fmt.Errorf("'path' must not be specified if 'source_type' is set to '%v'", me.SourceType)
	}
	if (me.Source == nil) && (slices.Contains([]string{"constant.string"}, string(me.SourceType))) {
		return fmt.Errorf("'source' must be specified if 'source_type' is set to '%v'", me.SourceType)
	}
	if (me.Source != nil) && (!slices.Contains([]string{"constant.string"}, string(me.SourceType))) {
		return fmt.Errorf("'source' must not be specified if 'source_type' is set to '%v'", me.SourceType)
	}
	return nil
}

func (me *EventDataAttributeComplex) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"path":        &me.Path,
		"source":      &me.Source,
		"source_type": &me.SourceType,
	})
}
