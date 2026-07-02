/**
* @license
* Copyright 2026 Dynatrace LLC
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

package settings

import (
	"fmt"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// DTLabelEnrichment adds a Dynatrace label (`dt.*`) to all monitored Azure
// entities. Either a static `literal` value or an Azure `tag_key` lookup must
// be provided — never both.
type DTLabelEnrichment struct {
	Label   string
	Literal string
	TagKey  string
}

type DTLabelEnrichments []*DTLabelEnrichment

func (me *DTLabelEnrichment) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"label": {
			Type:        schema.TypeString,
			Description: "Dynatrace label key, e.g. `dt.security_context` or `dt.cost.product`.",
			Required:    true,
		},
		"literal": {
			Type:        schema.TypeString,
			Description: "Static value applied to every monitored entity. Mutually exclusive with `tag_key`.",
			Optional:    true,
		},
		"tag_key": {
			Type:        schema.TypeString,
			Description: "Azure tag key whose value will be copied into the Dynatrace label. Mutually exclusive with `literal`.",
			Optional:    true,
		},
	}
}

func (me *DTLabelEnrichment) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"label":   me.Label,
		"literal": me.Literal,
		"tag_key": me.TagKey,
	})
}

func (me *DTLabelEnrichment) UnmarshalHCL(decoder hcl.Decoder) error {
	if err := decoder.DecodeAll(map[string]any{
		"label":   &me.Label,
		"literal": &me.Literal,
		"tag_key": &me.TagKey,
	}); err != nil {
		return err
	}
	if me.Literal == "" && me.TagKey == "" {
		return fmt.Errorf("dt_label_enrichment %q: exactly one of `literal` or `tag_key` must be set", me.Label)
	}
	if me.Literal != "" && me.TagKey != "" {
		return fmt.Errorf("dt_label_enrichment %q: `literal` and `tag_key` are mutually exclusive", me.Label)
	}
	return nil
}
