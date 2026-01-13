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

package notifications

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/export/sensitive"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type ServiceNow struct {
	Enabled   bool   `json:"-"`
	Name      string `json:"-"`
	ProfileID string `json:"-"`

	FormatProblemDetailsAsText *bool   `json:"formatProblemDetailsAsText,omitempty"` // Use text format for problem details instead of HTML.
	InstanceName               *string `json:"instanceName"`                         // The ServiceNow instance identifier. It refers to the first part of your own ServiceNow URL. \n\n This field is mutually exclusive with the **url** field. You can only use one of them
	URL                        *string `json:"url"`                                  // The URL of the on-premise ServiceNow installation. \n\n This field is mutually exclusive with the **instanceName** field. You can only use one of them
	Username                   string  `json:"username"`                             // The username of the ServiceNow account. \n\n Make sure that your user account has the `web_service_admin` and `x_dynat_ruxit.Integration` roles
	Password                   string  `json:"password"`                             // The password to the ServiceNow account
	Message                    string  `json:"message"`                              // The content of the ServiceNow description. Type '{' for placeholder suggestions
	SendIncidents              bool    `json:"sendIncidents"`                        // Send incidents into ServiceNow ITSM
	SendEvents                 bool    `json:"sendEvents"`                           // Send events into ServiceNow ITOM
}

func (me *ServiceNow) PrepareMarshalHCL(decoder hcl.Decoder) error {
	if password, ok := decoder.GetOk("password"); ok && len(password.(string)) > 0 {
		me.Password = password.(string)
	}
	return nil
}

func (me *ServiceNow) FillDemoValues() []string {
	me.Password = "#######"
	return []string{"Please fill in the password"}
}

func (me *ServiceNow) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:        schema.TypeString,
			Description: "The name of the notification configuration",
			Required:    true,
		},
		"active": {
			Type:        schema.TypeBool,
			Description: "This setting is enabled (`true`) or disabled (`false`)",
			Required:    true,
		},
		"profile": {
			Type:        schema.TypeString,
			Description: "The ID of the associated alerting profile",
			Required:    true,
		},

		"format_problem_details_as_text": {
			Type:        schema.TypeBool,
			Description: "Use text format for problem details instead of HTML.",
			Optional:    true, // nullable
		},
		"instance": {
			Type:        schema.TypeString,
			Description: "The ServiceNow instance identifier. It refers to the first part of your own ServiceNow URL. \n\n This field is mutually exclusive with the **url** field. You can only use one of them.",
			Optional:    true, // precondition
		},
		"message": {
			Type:        schema.TypeString,
			Description: "The content of the ServiceNow description. Type '{' for placeholder suggestions.. #### Available placeholders\n**{ImpactedEntity}**: A short description of the problem and impacted entity (or multiple impacted entities).\n\n**{ImpactedEntityNames}**: The entity impacted by the problem.\n\n**{NamesOfImpactedEntities}**: The names of all entities that are impacted by the problem.\n\n**{PID}**: Unique system identifier of the reported problem.\n\n**{ProblemDetailsHTML}**: All problem event details including root cause as an HTML-formatted string.\n\n**{ProblemDetailsText}**: All problem event details including root cause as a text-formatted string.\n\n**{ProblemID}**: Display number of the reported problem.\n\n**{ProblemImpact}**: Impact level of the problem. Possible values are APPLICATION, SERVICE, or INFRASTRUCTURE.\n\n**{ProblemSeverity}**: Severity level of the problem. Possible values are AVAILABILITY, ERROR, PERFORMANCE, RESOURCE_CONTENTION, or CUSTOM_ALERT.\n\n**{ProblemTitle}**: Short description of the problem.\n\n**{State}**: Problem state. Possible values are OPEN or RESOLVED.\n\n**{Tags}**: Comma separated list of tags that are defined for all impacted entities. To refer to the value of a specific tag, specify the tag's key in square brackets: **{Tags[key]}**. If the tag does not have any assigned value, the placeholder will be replaced by an empty string. The placeholder will not be replaced if the tag key does not exist.",
			Required:    true,
		},
		"password": {
			Type:        schema.TypeString,
			Description: "The password to the ServiceNow account.",
			Optional:    true,
			Sensitive:   true,
		},
		"events": {
			Type:        schema.TypeBool,
			Description: "Send events into ServiceNow ITOM.",
			Optional:    true,
		},
		"incidents": {
			Type:        schema.TypeBool,
			Description: "Send incidents into ServiceNow ITSM.",
			Required:    true,
		},
		"url": {
			Type:        schema.TypeString,
			Description: "The URL of the on-premise ServiceNow installation. \n\n This field is mutually exclusive with the **instanceName** field. You can only use one of them.",
			Optional:    true, // nullable
		},
		"username": {
			Type:        schema.TypeString,
			Description: "The username of the ServiceNow account. \n\n Make sure that your user account has the `web_service_admin` and `x_dynat_ruxit.Integration` roles.",
			Required:    true,
		},
		"legacy_id": {
			Type:        schema.TypeString,
			Description: "The ID of these settings when referred to from resources requiring the REST API V1 keys",
			Computed:    true,
			Optional:    true,
		},
	}
}

func (me *ServiceNow) MarshalHCL(properties hcl.Properties) error { // The password field MUST NOT get serialized into HCL here
	// The Dynatrace Settings 2.0 API delivers a scrambled version of any previously stored password here
	// Evaluation at this point would lead to that scrambled version to make it into the Terraform State
	// As a result any plans would be non-empty
	return properties.EncodeAll(sensitive.ConditionalIgnoreChangesMap(
		me.Schema(),
		map[string]any{
			"name":    me.Name,
			"active":  me.Enabled,
			"profile": me.ProfileID,

			"format_problem_details_as_text": me.FormatProblemDetailsAsText,
			"events":                         me.SendEvents,
			"incidents":                      me.SendIncidents,
			"url":                            me.URL,
			"username":                       me.Username,
			"instance":                       me.InstanceName,
			"message":                        me.Message,
			"password":                       me.Password,
		},
	))
}

func (me *ServiceNow) HandlePreconditions() error {
	// ---- InstanceName *string -> {"preconditions":[{"property":"url","type":"NULL"},{"expectedValue":"","property":"url","type":"EQUALS"}],"type":"OR"}
	return nil
}

func (me *ServiceNow) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"name":    &me.Name,
		"active":  &me.Enabled,
		"profile": &me.ProfileID,

		"format_problem_details_as_text": &me.FormatProblemDetailsAsText,
		"events":                         &me.SendEvents,
		"incidents":                      &me.SendIncidents,
		"url":                            &me.URL,
		"username":                       &me.Username,
		"instance":                       &me.InstanceName,
		"message":                        &me.Message,
		"password":                       &me.Password,
	})
}
