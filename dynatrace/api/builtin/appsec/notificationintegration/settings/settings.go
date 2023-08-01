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

package notificationintegration

import (
	"fmt"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/opt"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Settings struct {
	AttackCandidateBasedAlertingProfile *string                             `json:"attackCandidateBasedAlertingProfile,omitempty"` // For attack candidate alerts, select an [alerting profile](/ui/settings/builtin:appsec.notification-attack-alerting-profile) to control the delivery of security notifications related to this integration.
	AttackCandidateBasedEmailPayload    *AttackCandidateBasedEmailPayload   `json:"attackCandidateBasedEmailPayload,omitempty"`    // Attack candidate based email payload, required when `trigger` equals `ATTACK_CANDIDATE` and `type` equals `EMAIL`
	AttackCandidateBasedJiraPayload     *AttackCandidateBasedJiraPayload    `json:"attackCandidateBasedJiraPayload,omitempty"`     // Attack candidate based Jira payload, required when `trigger` equals `ATTACK_CANDIDATE` and `type` equals `JIRA`
	AttackCandidateBasedWebhookPayload  *AttackCandidateBasedWebhookPayload `json:"attackCandidateBasedWebhookPayload,omitempty"`  // Attack candidate based webhook payload, required when `trigger` equals `ATTACK_CANDIDATE` and `type` equals `WEBHOOK`
	DisplayName                         string                              `json:"displayName"`                                   // Display name of the security notification
	EmailConfiguration                  *EmailConfiguration                 `json:"emailConfiguration,omitempty"`                  // Email configuration, required when `type` equals `EMAIL`
	Enabled                             bool                                `json:"enabled"`                                       // Enable/Disable the security notification, enabled (`true`) or disabled (`false`)
	JiraConfiguration                   *JiraConfiguration                  `json:"jiraConfiguration,omitempty"`                   // Jira configuration, required when `type` equals `JIRA`
	SecurityProblemBasedAlertingProfile *string                             `json:"securityProblemBasedAlertingProfile,omitempty"` // For security problem alerts, select an [alerting profile](/ui/settings/builtin:appsec.notification-alerting-profile) to control the delivery of security notifications related to this integration.
	SecurityProblemBasedEmailPayload    *SecurityProblemBasedEmailPayload   `json:"securityProblemBasedEmailPayload,omitempty"`    // Security problem based email payload, required when `trigger` equals `SECURITY_PROBLEM` and `type` equals `EMAIL`
	SecurityProblemBasedJiraPayload     *SecurityProblemBasedJiraPayload    `json:"securityProblemBasedJiraPayload,omitempty"`     // Security problem based Jira payload, required when `trigger` equals `SECURITY_PROBLEM` and `type` equals `JIRA`
	SecurityProblemBasedWebhookPayload  *SecurityProblemBasedWebhookPayload `json:"securityProblemBasedWebhookPayload,omitempty"`  // Security problem based webhook payload, required when `trigger` equals `SECURITY_PROBLEM` and `type` equals `WEBHOOK`
	Trigger                             NotificationTrigger                 `json:"trigger"`                                       // Security alert type, possible Values: `ATTACK_CANDIDATE`, `SECURITY_PROBLEM`
	Type                                NotificationType                    `json:"type"`                                          // Notification type, possible Values: `EMAIL`, `JIRA`, `WEBHOOK`
	WebhookConfiguration                *WebhookConfiguration               `json:"webhookConfiguration,omitempty"`                // Webhook configuration, required when `type` equals `WEBHOOK`
}

func (me *Settings) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"attack_candidate_based_alerting_profile": {
			Type:        schema.TypeString,
			Description: "For attack candidate alerts, select an [alerting profile](/ui/settings/builtin:appsec.notification-attack-alerting-profile) to control the delivery of security notifications related to this integration.",
			Optional:    true, // precondition
		},
		"attack_candidate_based_email_payload": {
			Type:        schema.TypeList,
			Description: "Attack candidate based email payload, required when `trigger` equals `ATTACK_CANDIDATE` and `type` equals `EMAIL`",
			Optional:    true, // precondition
			Elem:        &schema.Resource{Schema: new(AttackCandidateBasedEmailPayload).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
		"attack_candidate_based_jira_payload": {
			Type:        schema.TypeList,
			Description: "Attack candidate based Jira payload, required when `trigger` equals `ATTACK_CANDIDATE` and `type` equals `JIRA`",
			Optional:    true, // precondition
			Elem:        &schema.Resource{Schema: new(AttackCandidateBasedJiraPayload).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
		"attack_candidate_based_webhook_payload": {
			Type:        schema.TypeList,
			Description: "Attack candidate based webhook payload, required when `trigger` equals `ATTACK_CANDIDATE` and `type` equals `WEBHOOK`",
			Optional:    true, // precondition
			Elem:        &schema.Resource{Schema: new(AttackCandidateBasedWebhookPayload).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
		"display_name": {
			Type:        schema.TypeString,
			Description: "Display name of the security notification",
			Required:    true,
		},
		"email_configuration": {
			Type:        schema.TypeList,
			Description: "Email configuration, required when `type` equals `EMAIL`",
			Optional:    true, // precondition
			Elem:        &schema.Resource{Schema: new(EmailConfiguration).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
		"enabled": {
			Type:        schema.TypeBool,
			Description: "Enable/Disable the security notification, enabled (`true`) or disabled (`false`)",
			Required:    true,
		},
		"jira_configuration": {
			Type:        schema.TypeList,
			Description: "Jira configuration, required when `type` equals `JIRA`",
			Optional:    true, // precondition
			Elem:        &schema.Resource{Schema: new(JiraConfiguration).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
		"security_problem_based_alerting_profile": {
			Type:        schema.TypeString,
			Description: "For security problem alerts, select an [alerting profile](/ui/settings/builtin:appsec.notification-alerting-profile) to control the delivery of security notifications related to this integration.",
			Optional:    true, // precondition
		},
		"security_problem_based_email_payload": {
			Type:        schema.TypeList,
			Description: "Security problem based email payload, required when `trigger` equals `SECURITY_PROBLEM` and `type` equals `EMAIL`",
			Optional:    true, // precondition
			Elem:        &schema.Resource{Schema: new(SecurityProblemBasedEmailPayload).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
		"security_problem_based_jira_payload": {
			Type:        schema.TypeList,
			Description: "Security problem based Jira payload, required when `trigger` equals `SECURITY_PROBLEM` and `type` equals `JIRA`",
			Optional:    true, // precondition
			Elem:        &schema.Resource{Schema: new(SecurityProblemBasedJiraPayload).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
		"security_problem_based_webhook_payload": {
			Type:        schema.TypeList,
			Description: "Security problem based webhook payload, required when `trigger` equals `SECURITY_PROBLEM` and `type` equals `WEBHOOK`",
			Optional:    true, // precondition
			Elem:        &schema.Resource{Schema: new(SecurityProblemBasedWebhookPayload).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
		"trigger": {
			Type:        schema.TypeString,
			Description: "Security alert type, possible Values: `ATTACK_CANDIDATE`, `SECURITY_PROBLEM`",
			Required:    true,
		},
		"type": {
			Type:        schema.TypeString,
			Description: "Notification type, possible Values: `EMAIL`, `JIRA`, `WEBHOOK`",
			Required:    true,
		},
		"webhook_configuration": {
			Type:        schema.TypeList,
			Description: "Webhook configuration, required when `type` equals `WEBHOOK`",
			Optional:    true, // precondition
			Elem:        &schema.Resource{Schema: new(WebhookConfiguration).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
	}
}

func (me *Settings) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"attack_candidate_based_alerting_profile": me.AttackCandidateBasedAlertingProfile,
		"attack_candidate_based_email_payload":    me.AttackCandidateBasedEmailPayload,
		"attack_candidate_based_jira_payload":     me.AttackCandidateBasedJiraPayload,
		"attack_candidate_based_webhook_payload":  me.AttackCandidateBasedWebhookPayload,
		"display_name":                            me.DisplayName,
		"email_configuration":                     me.EmailConfiguration,
		"enabled":                                 me.Enabled,
		"jira_configuration":                      me.JiraConfiguration,
		"security_problem_based_alerting_profile": me.SecurityProblemBasedAlertingProfile,
		"security_problem_based_email_payload":    me.SecurityProblemBasedEmailPayload,
		"security_problem_based_jira_payload":     me.SecurityProblemBasedJiraPayload,
		"security_problem_based_webhook_payload":  me.SecurityProblemBasedWebhookPayload,
		"trigger":                                 me.Trigger,
		"type":                                    me.Type,
		"webhook_configuration":                   me.WebhookConfiguration,
	})
}

func (me *Settings) HandlePreconditions() error {
	if (me.AttackCandidateBasedAlertingProfile == nil) && (string(me.Trigger) == "ATTACK_CANDIDATE") {
		me.AttackCandidateBasedAlertingProfile = opt.NewString("")
	}
	if (me.SecurityProblemBasedAlertingProfile == nil) && (string(me.Trigger) == "SECURITY_PROBLEM") {
		me.SecurityProblemBasedAlertingProfile = opt.NewString("")
	}
	if (me.AttackCandidateBasedEmailPayload == nil) && ((string(me.Type) == "EMAIL") && (string(me.Trigger) == "ATTACK_CANDIDATE")) {
		return fmt.Errorf("'attack_candidate_based_email_payload' must be specified if 'type' is set to '%v' and 'trigger' is set to '%v'", me.Type, me.Trigger)
	}
	if (me.AttackCandidateBasedEmailPayload != nil) && ((string(me.Type) != "EMAIL") || (string(me.Trigger) != "ATTACK_CANDIDATE")) {
		return fmt.Errorf("'attack_candidate_based_email_payload' must not be specified if 'type' is set to '%v'", me.Type)
	}
	if (me.AttackCandidateBasedJiraPayload == nil) && ((string(me.Type) == "JIRA") && (string(me.Trigger) == "ATTACK_CANDIDATE")) {
		return fmt.Errorf("'attack_candidate_based_jira_payload' must be specified if 'type' is set to '%v' and 'trigger' is set to '%v'", me.Type, me.Trigger)
	}
	if (me.AttackCandidateBasedJiraPayload != nil) && ((string(me.Type) != "JIRA") || (string(me.Trigger) != "ATTACK_CANDIDATE")) {
		return fmt.Errorf("'attack_candidate_based_jira_payload' must not be specified if 'type' is set to '%v'", me.Type)
	}
	if (me.AttackCandidateBasedWebhookPayload == nil) && ((string(me.Type) == "WEBHOOK") && (string(me.Trigger) == "ATTACK_CANDIDATE")) {
		return fmt.Errorf("'attack_candidate_based_webhook_payload' must be specified if 'type' is set to '%v' and 'trigger' is set to '%v'", me.Type, me.Trigger)
	}
	if (me.AttackCandidateBasedWebhookPayload != nil) && ((string(me.Type) != "WEBHOOK") || (string(me.Trigger) != "ATTACK_CANDIDATE")) {
		return fmt.Errorf("'attack_candidate_based_webhook_payload' must not be specified if 'type' is set to '%v'", me.Type)
	}
	if (me.EmailConfiguration == nil) && (string(me.Type) == "EMAIL") {
		return fmt.Errorf("'email_configuration' must be specified if 'type' is set to '%v'", me.Type)
	}
	if (me.EmailConfiguration != nil) && (string(me.Type) != "EMAIL") {
		return fmt.Errorf("'email_configuration' must not be specified if 'type' is set to '%v'", me.Type)
	}
	if (me.JiraConfiguration == nil) && (string(me.Type) == "JIRA") {
		return fmt.Errorf("'jira_configuration' must be specified if 'type' is set to '%v'", me.Type)
	}
	if (me.JiraConfiguration != nil) && (string(me.Type) != "JIRA") {
		return fmt.Errorf("'jira_configuration' must not be specified if 'type' is set to '%v'", me.Type)
	}
	if (me.SecurityProblemBasedEmailPayload == nil) && ((string(me.Type) == "EMAIL") && (string(me.Trigger) == "SECURITY_PROBLEM")) {
		return fmt.Errorf("'security_problem_based_email_payload' must be specified if 'type' is set to '%v' and 'trigger' is set to '%v'", me.Type, me.Trigger)
	}
	if (me.SecurityProblemBasedEmailPayload != nil) && ((string(me.Type) != "EMAIL") || (string(me.Trigger) != "SECURITY_PROBLEM")) {
		return fmt.Errorf("'security_problem_based_email_payload' must not be specified if 'type' is set to '%v'", me.Type)
	}
	if (me.SecurityProblemBasedJiraPayload == nil) && ((string(me.Type) == "JIRA") && (string(me.Trigger) == "SECURITY_PROBLEM")) {
		return fmt.Errorf("'security_problem_based_jira_payload' must be specified if 'type' is set to '%v' and 'trigger' is set to '%v'", me.Type, me.Trigger)
	}
	if (me.SecurityProblemBasedJiraPayload != nil) && ((string(me.Type) != "JIRA") || (string(me.Trigger) != "SECURITY_PROBLEM")) {
		return fmt.Errorf("'security_problem_based_jira_payload' must not be specified if 'type' is set to '%v'", me.Type)
	}
	if (me.SecurityProblemBasedWebhookPayload == nil) && ((string(me.Type) == "WEBHOOK") && (string(me.Trigger) == "SECURITY_PROBLEM")) {
		return fmt.Errorf("'security_problem_based_webhook_payload' must be specified if 'type' is set to '%v' and 'trigger' is set to '%v'", me.Type, me.Trigger)
	}
	if (me.SecurityProblemBasedWebhookPayload != nil) && ((string(me.Type) != "WEBHOOK") || (string(me.Trigger) != "SECURITY_PROBLEM")) {
		return fmt.Errorf("'security_problem_based_webhook_payload' must not be specified if 'type' is set to '%v'", me.Type)
	}
	if (me.WebhookConfiguration == nil) && (string(me.Type) == "WEBHOOK") {
		return fmt.Errorf("'webhook_configuration' must be specified if 'type' is set to '%v'", me.Type)
	}
	if (me.WebhookConfiguration != nil) && (string(me.Type) != "WEBHOOK") {
		return fmt.Errorf("'webhook_configuration' must not be specified if 'type' is set to '%v'", me.Type)
	}
	return nil
}

func (me *Settings) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"attack_candidate_based_alerting_profile": &me.AttackCandidateBasedAlertingProfile,
		"attack_candidate_based_email_payload":    &me.AttackCandidateBasedEmailPayload,
		"attack_candidate_based_jira_payload":     &me.AttackCandidateBasedJiraPayload,
		"attack_candidate_based_webhook_payload":  &me.AttackCandidateBasedWebhookPayload,
		"display_name":                            &me.DisplayName,
		"email_configuration":                     &me.EmailConfiguration,
		"enabled":                                 &me.Enabled,
		"jira_configuration":                      &me.JiraConfiguration,
		"security_problem_based_alerting_profile": &me.SecurityProblemBasedAlertingProfile,
		"security_problem_based_email_payload":    &me.SecurityProblemBasedEmailPayload,
		"security_problem_based_jira_payload":     &me.SecurityProblemBasedJiraPayload,
		"security_problem_based_webhook_payload":  &me.SecurityProblemBasedWebhookPayload,
		"trigger":                                 &me.Trigger,
		"type":                                    &me.Type,
		"webhook_configuration":                   &me.WebhookConfiguration,
	})
}
