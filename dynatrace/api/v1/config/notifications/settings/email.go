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
	"encoding/json"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/xjson"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// EmailConfig Configuration of the email notification.
type EmailConfig struct {
	BaseNotificationConfig
	BccReceivers []string `json:"bccReceivers,omitempty"` // The list of the email BCC-recipients.
	Body         string   `json:"body"`                   // The template of the email notification.  You can use the following placeholders:  * `{ImpactedEntities}`: Details about the entities impacted by the problem in form of a JSON array.  * `{ImpactedEntity}`: The entity impacted by the problem or *X* impacted entities.  * `{PID}`: The ID of the reported problem.  * `{ProblemDetailsHTML}`: All problem event details, including root cause, as an HTML-formatted string.  * `{ProblemDetailsJSON}`: All problem event details, including root cause, as a JSON object.  * `{ProblemDetailsMarkdown}`: All problem event details, including root cause, as a [Markdown-formatted](https://www.markdownguide.org/cheat-sheet/) string.  * `{ProblemDetailsText}`: All problem event details, including root cause, as a text-formatted string.  * `{ProblemID}`: The display number of the reported problem.  * `{ProblemImpact}`: The [impact level](https://www.dynatrace.com/support/help/shortlink/impact-analysis) of the problem. Possible values are `APPLICATION`, `SERVICE`, and `INFRASTRUCTURE`.  * `{ProblemSeverity}`: The [severity level](https://www.dynatrace.com/support/help/shortlink/event-types) of the problem. Possible values are `AVAILABILITY`, `ERROR`, `PERFORMANCE`, `RESOURCE_CONTENTION`, and `CUSTOM_ALERT`.  * `{ProblemTitle}`: A short description of the problem.  * `{ProblemURL}`: The URL of the problem within Dynatrace.  * `{State}`: The state of the problem. Possible values are `OPEN` and `RESOLVED`.  * `{Tags}`: The list of tags that are defined for all impacted entities, separated by commas.
	CcReceivers  []string `json:"ccReceivers,omitempty"`  // The list of the email CC-recipients.
	Receivers    []string `json:"receivers"`              // The list of the email recipients.
	Subject      string   `json:"subject"`                // The subject of the email notifications.
}

func (me *EmailConfig) GetType() Type {
	return Types.Email
}

func (me *EmailConfig) Schema() map[string]*schema.Schema {
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
		"alerting_profile": {
			Type:        schema.TypeString,
			Description: "The ID of the associated alerting profile",
			Required:    true,
		},
		"unknowns": {
			Type:        schema.TypeString,
			Description: "allows for configuring properties that are not explicitly supported by the current version of this provider",
			Optional:    true,
		},
		"bcc_receivers": {
			Type:        schema.TypeSet,
			Optional:    true,
			MinItems:    1,
			Description: "The list of the email BCC-recipients",
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
		"cc_receivers": {
			Type:        schema.TypeSet,
			Optional:    true,
			MinItems:    1,
			Description: "The list of the email CC-recipients",
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
		"receivers": {
			Type:        schema.TypeSet,
			Optional:    true,
			MinItems:    1,
			Description: "The list of the email recipients",
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
		"body": {
			Type:        schema.TypeString,
			Description: "The template of the email notification.  You can use the following placeholders:  * `{ImpactedEntities}`: Details about the entities impacted by the problem in form of a JSON array.  * `{ImpactedEntity}`: The entity impacted by the problem or *X* impacted entities.  * `{PID}`: The ID of the reported problem.  * `{ProblemDetailsHTML}`: All problem event details, including root cause, as an HTML-formatted string.  * `{ProblemDetailsJSON}`: All problem event details, including root cause, as a JSON object.  * `{ProblemDetailsMarkdown}`: All problem event details, including root cause, as a [Markdown-formatted](https://www.markdownguide.org/cheat-sheet/) string.  * `{ProblemDetailsText}`: All problem event details, including root cause, as a text-formatted string.  * `{ProblemID}`: The display number of the reported problem.  * `{ProblemImpact}`: The [impact level](https://www.dynatrace.com/support/help/shortlink/impact-analysis) of the problem. Possible values are `APPLICATION`, `SERVICE`, and `INFRASTRUCTURE`.  * `{ProblemSeverity}`: The [severity level](https://www.dynatrace.com/support/help/shortlink/event-types) of the problem. Possible values are `AVAILABILITY`, `ERROR`, `PERFORMANCE`, `RESOURCE_CONTENTION`, and `CUSTOM_ALERT`.  * `{ProblemTitle}`: A short description of the problem.  * `{ProblemURL}`: The URL of the problem within Dynatrace.  * `{State}`: The state of the problem. Possible values are `OPEN` and `RESOLVED`.  * `{Tags}`: The list of tags that are defined for all impacted entities, separated by commas",
			Required:    true,
		},
		"subject": {
			Type:        schema.TypeString,
			Description: "The subject of the email notifications",
			Required:    true,
		},
	}
}

func (me *EmailConfig) MarshalHCL(properties hcl.Properties) error {
	if err := properties.Unknowns(me.Unknowns); err != nil {
		return err
	}
	if err := properties.Encode("name", me.Name); err != nil {
		return err
	}
	if err := properties.Encode("active", me.Active); err != nil {
		return err
	}
	if err := properties.Encode("alerting_profile", me.AlertingProfile); err != nil {
		return err
	}
	if err := properties.Encode("bcc_receivers", me.BccReceivers); err != nil {
		return err
	}
	if err := properties.Encode("cc_receivers", me.CcReceivers); err != nil {
		return err
	}
	if err := properties.Encode("receivers", me.Receivers); err != nil {
		return err
	}
	if err := properties.Encode("body", me.Body); err != nil {
		return err
	}
	if err := properties.Encode("subject", me.Subject); err != nil {
		return err
	}
	return nil
}

func (me *EmailConfig) UnmarshalHCL(decoder hcl.Decoder) error {
	if value, ok := decoder.GetOk("unknowns"); ok {
		if err := json.Unmarshal([]byte(value.(string)), me); err != nil {
			return err
		}
		if err := json.Unmarshal([]byte(value.(string)), &me.Unknowns); err != nil {
			return err
		}
		delete(me.Unknowns, "name")
		delete(me.Unknowns, "active")
		delete(me.Unknowns, "alerting_profile")
		delete(me.Unknowns, "accept_any_certificate")
		delete(me.Unknowns, "custom_message")
		delete(me.Unknowns, "job_template_id")
		delete(me.Unknowns, "job_template_url")
		delete(me.Unknowns, "password")
		delete(me.Unknowns, "username")
		if len(me.Unknowns) == 0 {
			me.Unknowns = nil
		}
	}
	if value, ok := decoder.GetOk("name"); ok {
		me.Name = value.(string)
	}
	if value, ok := decoder.GetOk("active"); ok {
		me.Active = value.(bool)
	}
	if value, ok := decoder.GetOk("alerting_profile"); ok {
		me.AlertingProfile = value.(string)
	}
	if err := decoder.Decode("bcc_receivers", &me.BccReceivers); err != nil {
		return err
	}
	if err := decoder.Decode("cc_receivers", &me.CcReceivers); err != nil {
		return err
	}
	if err := decoder.Decode("receivers", &me.Receivers); err != nil {
		return err
	}
	if value, ok := decoder.GetOk("body"); ok {
		me.Body = value.(string)
	}
	if value, ok := decoder.GetOk("subject"); ok {
		me.Subject = value.(string)
	}
	return nil
}

func (me *EmailConfig) MarshalJSON() ([]byte, error) {
	properties := xjson.NewProperties(me.Unknowns)
	if err := properties.MarshalAll(map[string]any{
		"id":              me.ID,
		"name":            me.Name,
		"type":            me.GetType(),
		"active":          me.Active,
		"alertingProfile": me.AlertingProfile,
		"bccReceivers":    me.BccReceivers,
		"ccReceivers":     me.CcReceivers,
		"receivers":       me.Receivers,
		"body":            me.Body,
		"subject":         me.Subject,
	}); err != nil {
		return nil, err
	}
	return json.Marshal(properties)
}

func (me *EmailConfig) UnmarshalJSON(data []byte) error {
	properties := xjson.Properties{}
	if err := json.Unmarshal(data, &properties); err != nil {
		return err
	}
	if err := properties.UnmarshalAll(map[string]any{
		"id":              &me.ID,
		"name":            &me.Name,
		"type":            &me.Type,
		"active":          &me.Active,
		"alertingProfile": &me.AlertingProfile,
		"bccReceivers":    &me.BccReceivers,
		"ccReceivers":     &me.CcReceivers,
		"receivers":       &me.Receivers,
		"body":            &me.Body,
		"subject":         &me.Subject,
	}); err != nil {
		return err
	}
	if len(properties) > 0 {
		me.Unknowns = properties
	}
	return nil
}
