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

package teams

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Settings struct {
	AdditionalInformation    AdditionalInformations   `json:"additionalInformation,omitempty"`    // Define key/value pairs that further describe this team — for example, cost center, solution type, or business unit assignments.
	ContactDetails           ContactDetailss          `json:"contactDetails,omitempty"`           // Define options for messaging integration or other means of contacting this team.
	Description              *string                  `json:"description,omitempty"`              // Description
	ExternalID               *string                  `json:"externalId,omitempty"`               // This field should only be used for the automation purpose when importing team information. Once the external ID is created it can’t be changed.
	Identifier               string                   `json:"identifier"`                         // The team identifier is used to reference the team from any entity in Dynatrace. This identifier can’t be edited once the team has been created.
	Links                    Links                    `json:"links,omitempty"`                    // Include links to online resources where information relevant to this team’s responsibilities can be found.
	Name                     string                   `json:"name"`                               // Team name
	Responsibilities         *Responsibilities        `json:"responsibilities"`                   // Turn on all responsibility assignments that apply to this team.
	SupplementaryIdentifiers SupplementaryIdentifiers `json:"supplementaryIdentifiers,omitempty"` // The supplementary team identifiers can be optionally used in addition to the main team identifier to reference this team from any entity in Dynatrace. Up to 3 supplementary identifiers are supported.
}

func (me *Settings) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"additional_information": {
			Type:        schema.TypeList,
			Description: "Define key/value pairs that further describe this team — for example, cost center, solution type, or business unit assignments.",
			Optional:    true, // minobjects == 0
			Elem:        &schema.Resource{Schema: new(AdditionalInformations).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
		"contact_details": {
			Type:        schema.TypeList,
			Description: "Define options for messaging integration or other means of contacting this team.",
			Optional:    true, // minobjects == 0
			Elem:        &schema.Resource{Schema: new(ContactDetailss).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
		"description": {
			Type:        schema.TypeString,
			Description: "Description",
			Optional:    true, // nullable
		},
		"external_id": {
			Type:        schema.TypeString,
			Description: "This field should only be used for the automation purpose when importing team information. Once the external ID is created it can’t be changed.",
			Optional:    true, // nullable
		},
		"identifier": {
			Type:        schema.TypeString,
			Description: "The team identifier is used to reference the team from any entity in Dynatrace. This identifier can’t be edited once the team has been created.",
			Required:    true,
		},
		"links": {
			Type:        schema.TypeList,
			Description: "Include links to online resources where information relevant to this team’s responsibilities can be found.",
			Optional:    true, // minobjects == 0
			Elem:        &schema.Resource{Schema: new(Links).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
		"name": {
			Type:        schema.TypeString,
			Description: "Team name",
			Required:    true,
		},
		"responsibilities": {
			Type:        schema.TypeList,
			Description: "Turn on all responsibility assignments that apply to this team.",
			Required:    true,
			Elem:        &schema.Resource{Schema: new(Responsibilities).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
		"supplementary_identifiers": {
			Type:        schema.TypeList,
			Description: "The supplementary team identifiers can be optionally used in addition to the main team identifier to reference this team from any entity in Dynatrace. Up to 3 supplementary identifiers are supported.",
			Optional:    true, // minobjects == 0
			Elem:        &schema.Resource{Schema: new(SupplementaryIdentifiers).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
	}
}

func (me *Settings) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"additional_information":    me.AdditionalInformation,
		"contact_details":           me.ContactDetails,
		"description":               me.Description,
		"external_id":               me.ExternalID,
		"identifier":                me.Identifier,
		"links":                     me.Links,
		"name":                      me.Name,
		"responsibilities":          me.Responsibilities,
		"supplementary_identifiers": me.SupplementaryIdentifiers,
	})
}

func (me *Settings) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"additional_information":    &me.AdditionalInformation,
		"contact_details":           &me.ContactDetails,
		"description":               &me.Description,
		"external_id":               &me.ExternalID,
		"identifier":                &me.Identifier,
		"links":                     &me.Links,
		"name":                      &me.Name,
		"responsibilities":          &me.Responsibilities,
		"supplementary_identifiers": &me.SupplementaryIdentifiers,
	})
}
