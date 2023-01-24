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
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type ServiceNow struct {
	Enabled   bool   `json:"-"`
	Name      string `json:"-"`
	ProfileID string `json:"-"`

	InstanceName  *string `json:"instanceName"`  // The ServiceNow instance identifier. It refers to the first part of your own ServiceNow URL. \n\n This field is mutually exclusive with the **url** field. You can only use one of them
	URL           *string `json:"url"`           // The URL of the on-premise ServiceNow installation. \n\n This field is mutually exclusive with the **instanceName** field. You can only use one of them
	Username      string  `json:"username"`      // The username of the ServiceNow account. \n\n Make sure that your user account has the `web_service_admin` and `x_dynat_ruxit.Integration` roles
	Password      string  `json:"password"`      // The password to the ServiceNow account
	Message       string  `json:"message"`       // The content of the ServiceNow description. Type '{' for placeholder suggestions
	SendIncidents bool    `json:"sendIncidents"` // Send incidents into ServiceNow ITSM
	SendEvents    bool    `json:"sendEvents"`    // Send events into ServiceNow ITOM
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
			Description: "The configuration is enabled (`true`) or disabled (`false`)",
			Required:    true,
		},
		"profile": {
			Type:        schema.TypeString,
			Description: "The ID of the associated alerting profile",
			Required:    true,
		},

		"events": {
			Type:        schema.TypeBool,
			Description: "Send events into ServiceNow ITOM",
			Optional:    true,
		},
		"incidents": {
			Type:        schema.TypeBool,
			Description: "Send incidents into ServiceNow ITSM",
			Required:    true,
		},
		"url": {
			Type:        schema.TypeString,
			Description: "The URL of the on-premise ServiceNow installation. This field is mutually exclusive with the **instance** field. You can only use one of them",
			Optional:    true,
		},
		"username": {
			Type:        schema.TypeString,
			Description: "The username of the ServiceNow account.   Make sure that your user account has the `rest_service`, `web_request_admin`, and `x_dynat_ruxit.Integration` roles",
			Required:    true,
		},
		"instance": {
			Type:        schema.TypeString,
			Description: "The ServiceNow instance identifier. It refers to the first part of your own ServiceNow URL. This field is mutually exclusive with the **url** field. You can only use one of them",
			Optional:    true,
		},
		"message": {
			Type:        schema.TypeString,
			Description: "The content of the ServiceNow description. You can use the following placeholders:  * `{ImpactedEntity}`: The entity impacted by the problem or *X* impacted entities.  * `{PID}`: The ID of the reported problem.  * `{ProblemDetailsHTML}`: All problem event details, including root cause, as an HTML-formatted string.  * `{ProblemID}`: The display number of the reported problem.  * `{ProblemImpact}`: The [impact level](https://www.dynatrace.com/support/help/shortlink/impact-analysis) of the problem. Possible values are `APPLICATION`, `SERVICE`, and `INFRASTRUCTURE`.  * `{ProblemSeverity}`: The [severity level](https://www.dynatrace.com/support/help/shortlink/event-types) of the problem. Possible values are `AVAILABILITY`, `ERROR`, `PERFORMANCE`, `RESOURCE_CONTENTION`, and `CUSTOM_ALERT`.  * `{ProblemTitle}`: A short description of the problem.  * `{ProblemURL}`: The URL of the problem within Dynatrace.  * `{State}`: The state of the problem. Possible values are `OPEN` and `RESOLVED`.  * `{Tags}`: The list of tags that are defined for all impacted entities, separated by commas",
			Required:    true,
		},
		"password": {
			Type:        schema.TypeString,
			Description: "The password to the ServiceNow account",
			Sensitive:   true,
			Optional:    true,
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
	return properties.EncodeAll(map[string]any{
		"name":    me.Name,
		"active":  me.Enabled,
		"profile": me.ProfileID,

		"events":    me.SendEvents,
		"incidents": me.SendIncidents,
		"url":       me.URL,
		"username":  me.Username,
		"instance":  me.InstanceName,
		"message":   me.Message,
		"password":  me.Password,
	})
}

func (me *ServiceNow) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"name":    &me.Name,
		"active":  &me.Enabled,
		"profile": &me.ProfileID,

		"events":    &me.SendEvents,
		"incidents": &me.SendIncidents,
		"url":       &me.URL,
		"username":  &me.Username,
		"instance":  &me.InstanceName,
		"message":   &me.Message,
		"password":  &me.Password,
	})
}
