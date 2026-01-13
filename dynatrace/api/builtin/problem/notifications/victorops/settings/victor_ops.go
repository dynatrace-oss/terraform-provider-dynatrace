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

type VictorOps struct {
	Enabled   bool   `json:"-"`
	Name      string `json:"-"`
	ProfileID string `json:"-"`

	APIKey     string `json:"apiKey"`     // The API key for the target Splunk On-Call account.\n\nReceive your Splunk On-Call API key by navigating to: Settings -> Integrations -> Rest Endpoint -> Dynatrace within your VictorOps account
	RoutingKey string `json:"routingKey"` // The routing key, defining the group to be notified
	Message    string `json:"message"`    // The content of the message. Type '{' for placeholder suggestions
}

func (me *VictorOps) PrepareMarshalHCL(decoder hcl.Decoder) error {
	if apiKey, ok := decoder.GetOk("api_key"); ok && len(apiKey.(string)) > 0 {
		me.APIKey = apiKey.(string)
	}
	return nil
}

func (me *VictorOps) FillDemoValues() []string {
	me.APIKey = "#######"
	return []string{"Please fill in the API Key"}
}

func (me *VictorOps) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:        schema.TypeString,
			Description: "The name of the notification configuration.",
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

		"api_key": {
			Type:        schema.TypeString,
			Description: "The API key for the target Splunk On-Call account.\n\nReceive your Splunk On-Call API key by navigating to: Settings -> Integrations -> Rest Endpoint -> Dynatrace.",
			Optional:    true,
			Sensitive:   true,
		},
		"message": {
			Type:        schema.TypeString,
			Description: "The content of the message. Type '{' for placeholder suggestions.. #### Available placeholders\n**{ImpactedEntity}**: A short description of the problem and impacted entity (or multiple impacted entities).\n\n**{ImpactedEntityNames}**: The entity impacted by the problem.\n\n**{NamesOfImpactedEntities}**: The names of all entities that are impacted by the problem.\n\n**{ProblemDetailsText}**: All problem event details including root cause as a text-formatted string.\n\n**{ProblemID}**: Display number of the reported problem.\n\n**{ProblemImpact}**: Impact level of the problem. Possible values are APPLICATION, SERVICE, or INFRASTRUCTURE.\n\n**{ProblemSeverity}**: Severity level of the problem. Possible values are AVAILABILITY, ERROR, PERFORMANCE, RESOURCE_CONTENTION, or CUSTOM_ALERT.\n\n**{ProblemTitle}**: Short description of the problem.\n\n**{ProblemURL}**: URL of the problem within Dynatrace.\n\n**{State}**: Problem state. Possible values are OPEN or RESOLVED.",
			Required:    true,
		},
		"routing_key": {
			Type:        schema.TypeString,
			Description: "The routing key, defining the group to be notified.",
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

func (me *VictorOps) MarshalHCL(properties hcl.Properties) error { // The api_key field MUST NOT get serialized into HCL here
	// The Dynatrace Settings 2.0 API delivers a scrambled version of any previously stored api_key here
	// Evaluation at this point would lead to that scrambled version to make it into the Terraform State
	// As a result any plans would be non-empty
	return properties.EncodeAll(sensitive.ConditionalIgnoreChangesMap(
		me.Schema(),
		map[string]any{
			"name":    me.Name,
			"active":  me.Enabled,
			"profile": me.ProfileID,

			"message":     me.Message,
			"routing_key": me.RoutingKey,
			"api_key":     me.APIKey,
		},
	))
}

func (me *VictorOps) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"name":    &me.Name,
		"active":  &me.Enabled,
		"profile": &me.ProfileID,

		"message":     &me.Message,
		"routing_key": &me.RoutingKey,
		"api_key":     &me.APIKey,
	})
}
