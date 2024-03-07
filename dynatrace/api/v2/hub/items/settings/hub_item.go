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

package items

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type HubItem struct {
	ActivationLink       string   `json:"activationLink"`       // The activation link for a technology
	ArtifactId           string   `json:"artifactId"`           // The unique ID used by the artifacts contained in releases
	AuthorLogo           string   `json:"authorLogo"`           // Url for the author's logo
	AuthorName           string   `json:"authorName"`           // Name of the author of the item
	ClusterCompatible    bool     `json:"clusterCompatible"`    // Checks if the item is compatible with the cluster version
	ComingSoon           bool     `json:"comingSoon"`           // Whether or not the item is planned to be released soon
	Description          string   `json:"description"`          // Description of the item
	DocumentationLink    string   `json:"documentationLink"`    // An absolute link to the documentation page of this item
	HasDescriptionBlocks bool     `json:"hasDescriptionBlocks"` // Whether or not the details call will contain description blocks
	ItemId               string   `json:"itemId"`               // Unique Id of the item
	Logo                 string   `json:"logo"`                 // The logo of the item. Can be a URL or Base64 encoded. Intended for HTML tags
	MarketingLink        string   `json:"marketingLink"`        // An absolute link to the marketing page of this item
	Name                 string   `json:"name"`                 // Name of the item
	NotCompatibleReason  string   `json:"notCompatibleReason"`  // The reason why the item is not compatible with the cluster version
	Tags                 []string `json:"tags"`                 // Grouping of items with keywords
	Type                 string   `json:"type"`                 // Represents the type of item. It can be `TECHNOLOGY`, `EXTENSION1` or `EXTENSION2`
}

func (me *HubItem) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"activation_link": {
			Type:        schema.TypeString,
			Description: "The activation link for a technology",
			Computed:    true,
			Optional:    true,
		},
		"artifact_id": {
			Type:        schema.TypeString,
			Description: "The unique ID used by the artifacts contained in releases",
			Computed:    true,
		},
		"author_logo": {
			Type:        schema.TypeString,
			Description: "URL for the author's logo",
			Computed:    true,
			Optional:    true,
		},
		"author_name": {
			Type:        schema.TypeString,
			Description: "Name of the author of the item",
			Computed:    true,
			Optional:    true,
		},
		"cluster_compatible": {
			Type:        schema.TypeBool,
			Description: "Checks if the item is compatible with the cluster version",
			Computed:    true,
			Optional:    true,
		},
		"coming_soon": {
			Type:        schema.TypeBool,
			Description: "Whether or not the item is planned to be released soon",
			Computed:    true,
			Optional:    true,
		},
		"description": {
			Type:        schema.TypeString,
			Description: "Description of the item",
			Computed:    true,
			Optional:    true,
		},
		"documentation_link": {
			Type:        schema.TypeString,
			Description: "An absolute link to the documentation page of this item",
			Computed:    true,
			Optional:    true,
		},
		"has_description_blocks": {
			Type:        schema.TypeBool,
			Description: "Whether or not the details call will contain description blocks",
			Computed:    true,
			Optional:    true,
		},
		"item_id": {
			Type:        schema.TypeString,
			Description: "Unique Id of the item",
			Computed:    true,
		},
		"logo": {
			Type:        schema.TypeString,
			Description: "The logo of the item. Can be a URL or Base64 encoded. Intended for HTML tags",
			Computed:    true,
			Optional:    true,
		},
		"marketing_link": {
			Type:        schema.TypeString,
			Description: "An absolute link to the marketing page of this item",
			Computed:    true,
			Optional:    true,
		},
		"name": {
			Type:        schema.TypeString,
			Description: "Name of the item",
			Computed:    true,
		},
		"not_compatible_reason": {
			Type:        schema.TypeString,
			Description: "The reason why the item is not compatible with the cluster version",
			Computed:    true,
			Optional:    true,
		},
		"tags": {
			Type:        schema.TypeSet,
			Description: "Grouping of items with keywords",
			Computed:    true,
			Optional:    true,
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
		"type": {
			Type:        schema.TypeString,
			Description: "Represents the type of item. It can be `TECHNOLOGY`, `EXTENSION1` or `EXTENSION2`",
			Computed:    true,
		},
	}
}

func (me *HubItem) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"activation_link":        me.ActivationLink,
		"artifact_id":            me.ArtifactId,
		"author_logo":            me.AuthorLogo,
		"author_name":            me.AuthorName,
		"cluster_compatible":     me.ClusterCompatible,
		"coming_soon":            me.ComingSoon,
		"description":            me.Description,
		"documentation_link":     me.DocumentationLink,
		"has_description_blocks": me.HasDescriptionBlocks,
		"item_id":                me.ItemId,
		"logo":                   me.Logo,
		"marketing_link":         me.MarketingLink,
		"name":                   me.Name,
		"not_compatible_reason":  me.NotCompatibleReason,
		"tags":                   me.Tags,
		"type":                   me.Type,
	})
}

func (me *HubItem) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"activation_link":        &me.ActivationLink,
		"artifact_id":            &me.ArtifactId,
		"author_logo":            &me.AuthorLogo,
		"author_name":            &me.AuthorName,
		"cluster_compatible":     &me.ClusterCompatible,
		"coming_soon":            &me.ComingSoon,
		"description":            &me.Description,
		"documentation_link":     &me.DocumentationLink,
		"has_description_blocks": &me.HasDescriptionBlocks,
		"item_id":                &me.ItemId,
		"logo":                   &me.Logo,
		"marketing_link":         &me.MarketingLink,
		"name":                   &me.Name,
		"not_compatible_reason":  &me.NotCompatibleReason,
		"tags":                   &me.Tags,
		"type":                   &me.Type,
	})
}
