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

package ingestsources

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Settings struct {
	DefaultBucket *string        `json:"defaultBucket,omitempty"` // Default Bucket
	DisplayName   string         `json:"displayName"`             // Endpoint display name
	Enabled       bool           `json:"enabled"`                 // This setting is enabled (`true`) or disabled (`false`)
	PathSegment   string         `json:"pathSegment"`             // Endpoint segment
	Processing    *Stage         `json:"processing"`              // Processing stage
	StaticRouting *StaticRouting `json:"staticRouting,omitempty"` // Static routing of endpoint
}

func (me *Settings) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"default_bucket": {
			Type:        schema.TypeString,
			Description: "Default Bucket",
			Optional:    true, // nullable
		},
		"display_name": {
			Type:        schema.TypeString,
			Description: "Endpoint display name",
			Required:    true,
		},
		"enabled": {
			Type:        schema.TypeBool,
			Description: "This setting is enabled (`true`) or disabled (`false`)",
			Required:    true,
		},
		"path_segment": {
			Type:        schema.TypeString,
			Description: "Endpoint segment",
			Required:    true,
		},
		"processing": {
			Type:        schema.TypeList,
			Description: "Processing stage",
			Required:    true,
			Elem:        &schema.Resource{Schema: new(Stage).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
		"static_routing": {
			Type:        schema.TypeList,
			Description: "Static routing of endpoint",
			Optional:    true, // nullable
			Elem:        &schema.Resource{Schema: new(StaticRouting).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
	}
}

func (me *Settings) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"default_bucket": me.DefaultBucket,
		"display_name":   me.DisplayName,
		"enabled":        me.Enabled,
		"path_segment":   me.PathSegment,
		"processing":     me.Processing,
		"static_routing": me.StaticRouting,
	})
}

func (me *Settings) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"default_bucket": &me.DefaultBucket,
		"display_name":   &me.DisplayName,
		"enabled":        &me.Enabled,
		"path_segment":   &me.PathSegment,
		"processing":     &me.Processing,
		"static_routing": &me.StaticRouting,
	})
}
