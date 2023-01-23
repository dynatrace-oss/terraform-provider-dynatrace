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
	"strconv"
	"strings"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/xjson"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/opt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// AnsibleTowerConfig Configuration of the Ansible Tower notification.
type AnsibleTowerConfig struct {
	BaseNotificationConfig
	AcceptAnyCertificate bool    `json:"acceptAnyCertificate"` // Accept any, including self-signed and invalid, SSL certificate (`true`) or only trusted (`false`) certificates.
	CustomMessage        string  `json:"customMessage"`        // The custom message of the notification.   This message will be displayed in the extra variables **Message** field of your job template.  You can use the following placeholders:  * `{ImpactedEntities}`: Details about the entities impacted by the problem in form of a JSON array.  * `{ImpactedEntity}`: The entity impacted by the problem or *X* impacted entities.  * `{PID}`: The ID of the reported problem.  * `{ProblemDetailsText}`: All problem event details, including root cause, as a text-formatted string.  * `{ProblemID}`: The display number of the reported problem.  * `{ProblemImpact}`: The [impact level](https://www.dynatrace.com/support/help/shortlink/impact-analysis) of the problem. Possible values are `APPLICATION`, `SERVICE`, and `INFRASTRUCTURE`.  * `{ProblemSeverity}`: The [severity level](https://www.dynatrace.com/support/help/shortlink/event-types) of the problem. Possible values are `AVAILABILITY`, `ERROR`, `PERFORMANCE`, `RESOURCE_CONTENTION`, and `CUSTOM_ALERT`.  * `{ProblemTitle}`: A short description of the problem.  * `{ProblemURL}`: The URL of the problem within Dynatrace.  * `{State}`: The state of the problem. Possible values are `OPEN` and `RESOLVED`.  * `{Tags}`: The list of tags that are defined for all impacted entities, separated by commas.
	JobTemplateID        int32   `json:"jobTemplateID"`        // The ID of the target Ansible Tower job template.
	JobTemplateURL       string  `json:"jobTemplateURL"`       // The URL of the target Ansible Tower job template.
	Password             *string `json:"password,omitempty"`   // The password for the Ansible Tower account.
	Username             string  `json:"username"`             // The username of the Ansible Tower account.
}

func (me *AnsibleTowerConfig) GetType() Type {
	return Types.Ansibletower
}

func (me *AnsibleTowerConfig) Schema() map[string]*schema.Schema {
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
		"accept_any_certificate": {
			Type:        schema.TypeBool,
			Description: "Accept any, including self-signed and invalid, SSL certificate (`true`) or only trusted (`false`) certificates",
			Required:    true,
		},
		"custom_message": {
			Type:        schema.TypeString,
			Description: "The custom message of the notification.   This message will be displayed in the extra variables **Message** field of your job template.  You can use the following placeholders:  * `{ImpactedEntities}`: Details about the entities impacted by the problem in form of a JSON array.  * `{ImpactedEntity}`: The entity impacted by the problem or *X* impacted entities.  * `{PID}`: The ID of the reported problem.  * `{ProblemDetailsText}`: All problem event details, including root cause, as a text-formatted string.  * `{ProblemID}`: The display number of the reported problem.  * `{ProblemImpact}`: The [impact level](https://www.dynatrace.com/support/help/shortlink/impact-analysis) of the problem. Possible values are `APPLICATION`, `SERVICE`, and `INFRASTRUCTURE`.  * `{ProblemSeverity}`: The [severity level](https://www.dynatrace.com/support/help/shortlink/event-types) of the problem. Possible values are `AVAILABILITY`, `ERROR`, `PERFORMANCE`, `RESOURCE_CONTENTION`, and `CUSTOM_ALERT`.  * `{ProblemTitle}`: A short description of the problem.  * `{ProblemURL}`: The URL of the problem within Dynatrace.  * `{State}`: The state of the problem. Possible values are `OPEN` and `RESOLVED`.  * `{Tags}`: The list of tags that are defined for all impacted entities, separated by commas",
			Required:    true,
		},
		"job_template_id": {
			Type:        schema.TypeInt,
			Description: "The ID of the target Ansible Tower job template",
			Required:    true,
		},
		"job_template_url": {
			Type:        schema.TypeString,
			Description: "The URL of the target Ansible Tower job template",
			Required:    true,
		},
		"password": {
			Type:        schema.TypeString,
			Description: "The password for the Ansible Tower account",
			Optional:    true,
		},
		"username": {
			Type:        schema.TypeString,
			Description: "The username of the Ansible Tower account",
			Required:    true,
		},
	}
}

func (me *AnsibleTowerConfig) FillDemoValues() []string {
	me.Password = opt.NewString("#######")
	if me.JobTemplateID == 0 && len(me.JobTemplateURL) > 0 {
		if idx := strings.LastIndex(me.JobTemplateURL, "/"); idx != -1 {
			if iJobTemplateID, err := strconv.Atoi(me.JobTemplateURL[idx+1:]); err == nil {
				me.JobTemplateID = int32(iJobTemplateID)
			}
		}
	}
	return []string{"The REST API didn't provide the credentials"}
}

func (me *AnsibleTowerConfig) PrepareMarshalHCL(decoder hcl.Decoder) error {
	if password, ok := decoder.GetOk("ansible_tower.0.password"); ok && len(password.(string)) > 0 {
		me.Password = opt.NewString(password.(string))
	}
	if job_template_id, ok := decoder.GetOk("ansible_tower.0.job_template_id"); ok && job_template_id.(int) != 0 {
		me.JobTemplateID = int32(job_template_id.(int))
	}
	if me.JobTemplateID == 0 && len(me.JobTemplateURL) > 0 {
		if idx := strings.LastIndex(me.JobTemplateURL, "/"); idx != -1 {
			sJobTemplateID := me.JobTemplateURL[idx+1:]
			if iJobTemplateID, err := strconv.Atoi(sJobTemplateID); err != nil {
				me.JobTemplateID = int32(iJobTemplateID)
			}
		}
	}
	return nil
}

func (me *AnsibleTowerConfig) MarshalHCL(properties hcl.Properties) error {
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
	if err := properties.Encode("accept_any_certificate", me.AcceptAnyCertificate); err != nil {
		return err
	}
	if err := properties.Encode("custom_message", me.CustomMessage); err != nil {
		return err
	}
	if err := properties.Encode("job_template_id", int(me.JobTemplateID)); err != nil {
		return err
	}
	if err := properties.Encode("job_template_url", me.JobTemplateURL); err != nil {
		return err
	}
	if err := properties.Encode("password", me.Password); err != nil {
		return err
	}
	if err := properties.Encode("username", me.Username); err != nil {
		return err
	}
	return nil
}

func (me *AnsibleTowerConfig) UnmarshalHCL(decoder hcl.Decoder) error {
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
	if value, ok := decoder.GetOk("accept_any_certificate"); ok {
		me.AcceptAnyCertificate = value.(bool)
	}
	if value, ok := decoder.GetOk("custom_message"); ok {
		me.CustomMessage = value.(string)
	}
	if value, ok := decoder.GetOk("job_template_id"); ok {
		me.JobTemplateID = int32(value.(int))
	}
	if value, ok := decoder.GetOk("job_template_url"); ok {
		me.JobTemplateURL = value.(string)
	}
	if value, ok := decoder.GetOk("password"); ok {
		me.Password = opt.NewString(value.(string))
	}
	if value, ok := decoder.GetOk("username"); ok {
		me.Username = value.(string)
	}
	return nil
}

func (me *AnsibleTowerConfig) MarshalJSON() ([]byte, error) {
	properties := xjson.NewProperties(me.Unknowns)
	if err := properties.MarshalAll(map[string]any{
		"id":                   me.ID,
		"name":                 me.Name,
		"type":                 me.GetType(),
		"active":               me.Active,
		"alertingProfile":      me.AlertingProfile,
		"acceptAnyCertificate": me.AcceptAnyCertificate,
		"customMessage":        me.CustomMessage,
		"jobTemplateID":        me.JobTemplateID,
		"jobTemplateURL":       me.JobTemplateURL,
		"password":             me.Password,
		"username":             me.Username,
	}); err != nil {
		return nil, err
	}
	return json.Marshal(properties)
}

func (me *AnsibleTowerConfig) UnmarshalJSON(data []byte) error {
	properties := xjson.Properties{}
	if err := json.Unmarshal(data, &properties); err != nil {
		return err
	}
	if err := properties.UnmarshalAll(map[string]any{
		"id":                   &me.ID,
		"name":                 &me.Name,
		"type":                 &me.Type,
		"active":               &me.Active,
		"alertingProfile":      &me.AlertingProfile,
		"acceptAnyCertificate": &me.AcceptAnyCertificate,
		"customMessage":        &me.CustomMessage,
		"jobTemplateID":        &me.JobTemplateID,
		"jobTemplateURL":       &me.JobTemplateURL,
		"password":             &me.Password,
		"username":             &me.Username,
	}); err != nil {
		return err
	}
	if len(properties) > 0 {
		me.Unknowns = properties
	}
	return nil
}
