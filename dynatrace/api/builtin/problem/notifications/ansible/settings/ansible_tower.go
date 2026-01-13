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

type AnsibleTower struct {
	Enabled   bool   `json:"-" hcl:"active,omitempty"` // The notification is active (`true`) or inactive (`false`). Default is `false`
	Name      string `json:"-" hcl:"name"`             // The display name within the Dynatrace WebUI
	ProfileID string `json:"-" hcl:"alerting_profile"` // The ID of the associated alerting profile

	JobTemplateURL string `json:"jobTemplateURL" hcl:"job_template_url"`             // The URL of the target Ansible Tower job template.\n\nFor example, https://<Ansible Tower server name>/#/templates/job_template/<JobTemplateID>\n\n**Note:** Be sure to select the **Prompt on Launch** option in the Extra Variables section of your job template configuration
	Insecure       bool   `json:"acceptAnyCertificate" hcl:"accept_any_certificate"` // Accept any SSL certificate (including self-signed and invalid certificates)
	Username       string `json:"username" hcl:"username"`                           // The username of the Ansible Tower account
	Password       string `json:"password" hcl:"password,secret"`                    // The password for the Ansible Tower account
	CustomMessage  string `json:"customMessage" hcl:"custom_message"`                // This message will be displayed in the Extra Variables **Message** field of your job template. Type '{' for placeholder suggestions
}

func (me *AnsibleTower) PrepareMarshalHCL(decoder hcl.Decoder) error {
	if password, ok := decoder.GetOk("password"); ok && len(password.(string)) > 0 {
		me.Password = password.(string)
	}
	return nil
}

func (me *AnsibleTower) FillDemoValues() []string {
	me.Password = "#######"
	return []string{"Please fill in the correct password"}
}

func (me *AnsibleTower) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:        schema.TypeString,
			Description: "The name of the notification configuration.",
			Required:    true,
		},
		"active": {
			Type:        schema.TypeBool,
			Description: "This setting is enabled (`true`) or disabled (`false`)",
			Optional:    true,
		},
		"profile": {
			Type:        schema.TypeString,
			Description: "The ID of the associated alerting profile.",
			Required:    true,
		},
		"insecure": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "Accept any SSL certificate (including self-signed and invalid certificates)",
		},
		"custom_message": {
			Type:        schema.TypeString,
			Description: "This message will be displayed in the Extra Variables **Message** field of your job template. Type '{' for placeholder suggestions.. #### Available placeholders\n**{ImpactedEntities}**: Details about the entities impacted by the problem in form of a json array.\n\n**{ImpactedEntity}**: A short description of the problem and impacted entity (or multiple impacted entities).\n\n**{ImpactedEntityNames}**: The entity impacted by the problem.\n\n**{NamesOfImpactedEntities}**: The names of all entities that are impacted by the problem.\n\n**{PID}**: Unique system identifier of the reported problem.\n\n**{ProblemDetailsText}**: All problem event details including root cause as a text-formatted string.\n\n**{ProblemID}**: Display number of the reported problem.\n\n**{ProblemImpact}**: Impact level of the problem. Possible values are APPLICATION, SERVICE, or INFRASTRUCTURE.\n\n**{ProblemSeverity}**: Severity level of the problem. Possible values are AVAILABILITY, ERROR, PERFORMANCE, RESOURCE_CONTENTION, or CUSTOM_ALERT.\n\n**{ProblemTitle}**: Short description of the problem.\n\n**{ProblemURL}**: URL of the problem within Dynatrace.\n\n**{State}**: Problem state. Possible values are OPEN or RESOLVED.\n\n**{Tags}**: Comma separated list of tags that are defined for all impacted entities. To refer to the value of a specific tag, specify the tag's key in square brackets: **{Tags[key]}**. If the tag does not have any assigned value, the placeholder will be replaced by an empty string. The placeholder will not be replaced if the tag key does not exist.",
			Required:    true,
		},
		"job_template_url": {
			Type:        schema.TypeString,
			Description: "The URL of the target job template.\n\nFor example, https://<Ansible server name>/#/templates/job_template/<JobTemplateID>\n\n**Note:** Be sure to select the **Prompt on Launch** option in the Extra Variables section of your job template configuration.",
			Required:    true,
		},
		"password": {
			Type:        schema.TypeString,
			Description: "Account password.",
			Sensitive:   true,
			Optional:    true,
		},
		"username": {
			Type:        schema.TypeString,
			Description: "Account username.",
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

func (me *AnsibleTower) MarshalHCL(properties hcl.Properties) error { // The password field MUST NOT get serialized into HCL here
	// The Dynatrace Settings 2.0 API delivers a scrambled version of any previously stored password here
	// Evaluation at this point would lead to that scrambled version to make it into the Terraform State
	// As a result any plans would be non-empty
	return properties.EncodeAll(sensitive.ConditionalIgnoreChangesMap(
		me.Schema(),
		map[string]any{
			"name":    me.Name,
			"active":  me.Enabled,
			"profile": me.ProfileID,

			"insecure":         me.Insecure,
			"custom_message":   me.CustomMessage,
			"job_template_url": me.JobTemplateURL,
			"password":         me.Password,
			"username":         me.Username,
		},
	))
}

func (me *AnsibleTower) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"name":    &me.Name,
		"active":  &me.Enabled,
		"profile": &me.ProfileID,

		"insecure":         &me.Insecure,
		"custom_message":   &me.CustomMessage,
		"job_template_url": &me.JobTemplateURL,
		"password":         &me.Password,
		"username":         &me.Username,
	})
}
