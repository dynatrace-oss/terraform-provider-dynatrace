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

package web

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// WaterfallSettings These settings influence the monitoring data you receive for 3rd party, CDN, and 1st party resources
type WaterfallSettings struct {
	UncompressedResourcesThreshold           int32 `json:"uncompressedResourcesThreshold"`           // Warn about uncompressed resources larger than *X* bytes. Values between 0 and 99999 are allowed.
	ResourcesThreshold                       int32 `json:"resourcesThreshold"`                       // Warn about resources larger than *X* bytes. Values between 0 and 99999000 are allowed.
	ResourceBrowserCachingThreshold          int32 `json:"resourceBrowserCachingThreshold"`          // Warn about resources with a lower browser cache rate above *X*%. Values between 1 and 100 are allowed.
	SlowFirstPartyResourcesThreshold         int32 `json:"slowFirstPartyResourcesThreshold"`         // Warn about slow 1st party resources with a response time above *X* ms. Values between 0 and 99999000 are allowed.
	SlowThirdPartyResourcesThreshold         int32 `json:"slowThirdPartyResourcesThreshold"`         // Warn about slow 3rd party resources with a response time above *X* ms. Values between 0 and 99999000 are allowed.
	SlowCdnResourcesThreshold                int32 `json:"slowCdnResourcesThreshold"`                // Warn about slow CDN resources with a response time above *X* ms. Values between 0 and 99999000 are allowed.
	SpeedIndexVisuallyCompleteRatioThreshold int32 `json:"speedIndexVisuallyCompleteRatioThreshold"` // Warn if Speed index exceeds *X* % of Visually complete. Values between 1 and 99 are allowed.
}

func (me *WaterfallSettings) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"uncompressed_resources_threshold": {
			Type:        schema.TypeInt,
			Description: "Warn about uncompressed resources larger than *X* bytes. Values between 0 and 99999 are allowed.",
			Required:    true,
		},
		"resources_threshold": {
			Type:        schema.TypeInt,
			Description: "Warn about resources larger than *X* bytes. Values between 0 and 99999000 are allowed.",
			Required:    true,
		},
		"resource_browser_caching_threshold": {
			Type:        schema.TypeInt,
			Description: "Warn about resources with a lower browser cache rate above *X*%. Values between 1 and 100 are allowed.",
			Required:    true,
		},
		"slow_first_party_resources_threshold": {
			Type:        schema.TypeInt,
			Description: "Warn about slow 1st party resources with a response time above *X* ms. Values between 0 and 99999000 are allowed.",
			Required:    true,
		},
		"slow_third_party_resources_threshold": {
			Type:        schema.TypeInt,
			Description: "Warn about slow 3rd party resources with a response time above *X* ms. Values between 0 and 99999000 are allowed.",
			Required:    true,
		},
		"slow_cnd_resources_threshold": {
			Type:        schema.TypeInt,
			Description: "Warn about slow CDN resources with a response time above *X* ms. Values between 0 and 99999000 are allowed.",
			Required:    true,
		},
		"speed_index_visually_complete_ratio_threshold": {
			Type:        schema.TypeInt,
			Description: "Warn if Speed index exceeds *X* % of Visually complete. Values between 1 and 99 are allowed.",
			Required:    true,
		},
	}
}

func (me *WaterfallSettings) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"uncompressed_resources_threshold":              me.UncompressedResourcesThreshold,
		"resources_threshold":                           me.ResourcesThreshold,
		"resource_browser_caching_threshold":            me.ResourceBrowserCachingThreshold,
		"slow_first_party_resources_threshold":          me.SlowFirstPartyResourcesThreshold,
		"slow_third_party_resources_threshold":          me.SlowThirdPartyResourcesThreshold,
		"slow_cnd_resources_threshold":                  me.SlowCdnResourcesThreshold,
		"speed_index_visually_complete_ratio_threshold": me.SpeedIndexVisuallyCompleteRatioThreshold,
	})
}

func (me *WaterfallSettings) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"uncompressed_resources_threshold":              &me.UncompressedResourcesThreshold,
		"resources_threshold":                           &me.ResourcesThreshold,
		"resource_browser_caching_threshold":            &me.ResourceBrowserCachingThreshold,
		"slow_first_party_resources_threshold":          &me.SlowFirstPartyResourcesThreshold,
		"slow_third_party_resources_threshold":          &me.SlowThirdPartyResourcesThreshold,
		"slow_cnd_resources_threshold":                  &me.SlowCdnResourcesThreshold,
		"speed_index_visually_complete_ratio_threshold": &me.SpeedIndexVisuallyCompleteRatioThreshold,
	})
}
