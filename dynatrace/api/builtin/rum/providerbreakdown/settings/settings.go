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

package providerbreakdown

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Settings struct {
	DomainNamePatternList   DomainNamePatternListObjects `json:"domainNamePatternList"`   // Domain name pattern
	IconUrl                 *string                      `json:"iconUrl,omitempty"`       // Specify an URL for the provider's brand icon
	ReportPublicImprovement bool                         `json:"reportPublicImprovement"` // Send the patterns of this provider to Dynatrace to help us improve 3rd-party detection.
	ResourceName            string                       `json:"resourceName"`            // Resource name
	ResourceType            ResourceType                 `json:"resourceType"`            // Possible Values: `FirstParty`, `ThirdParty`, `Cdn`
}

func (me *Settings) Name() string {
	return me.ResourceName
}

func (me *Settings) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"domain_name_pattern_list": {
			Type:        schema.TypeList,
			Description: "Domain name pattern",
			Required:    true,

			Elem:     &schema.Resource{Schema: new(DomainNamePatternListObjects).Schema()},
			MinItems: 1,
			MaxItems: 1,
		},
		"icon_url": {
			Type:        schema.TypeString,
			Description: "Specify an URL for the provider's brand icon",
			Optional:    true,
		},
		"report_public_improvement": {
			Type:        schema.TypeBool,
			Description: "Send the patterns of this provider to Dynatrace to help us improve 3rd-party detection.",
			Required:    true,
		},
		"resource_name": {
			Type:        schema.TypeString,
			Description: "Resource name",
			Required:    true,
		},
		"resource_type": {
			Type:        schema.TypeString,
			Description: "Possible Values: `FirstParty`, `ThirdParty`, `Cdn`",
			Required:    true,
		},
	}
}

func (me *Settings) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"domain_name_pattern_list":  me.DomainNamePatternList,
		"icon_url":                  me.IconUrl,
		"report_public_improvement": me.ReportPublicImprovement,
		"resource_name":             me.ResourceName,
		"resource_type":             me.ResourceType,
	})
}

func (me *Settings) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"domain_name_pattern_list":  &me.DomainNamePatternList,
		"icon_url":                  &me.IconUrl,
		"report_public_improvement": &me.ReportPublicImprovement,
		"resource_name":             &me.ResourceName,
		"resource_type":             &me.ResourceType,
	})
}
