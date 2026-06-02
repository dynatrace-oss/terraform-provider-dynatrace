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
	"fmt"
	"slices"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type ContactDetailss []*ContactDetails

func (me *ContactDetailss) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"contact_detail": {
			Type:        schema.TypeList,
			Required:    true,
			MinItems:    1,
			Description: "",
			Elem:        &schema.Resource{Schema: new(ContactDetails).Schema()},
		},
	}
}

func (me ContactDetailss) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeSlice("contact_detail", me)
}

func (me *ContactDetailss) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeSlice("contact_detail", me)
}

type ContactDetails struct {
	Email           *string         `json:"email,omitempty"`
	IntegrationType IntegrationType `json:"integrationType"` // Integration type. Possible values: `EMAIL`, `JIRA`, `MS_TEAMS`, `SLACK`
	Jira            *JiraConnection `json:"jira,omitempty"`
	MsTeams         *string         `json:"msTeams,omitempty"`      // Team
	SlackChannel    *string         `json:"slackChannel,omitempty"` // Channel
	Url             *string         `json:"url,omitempty"`
}

func (me *ContactDetails) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"email": {
			Type:        schema.TypeString,
			Description: "No documentation available",
			Optional:    true, // precondition
		},
		"integration_type": {
			Type:        schema.TypeString,
			Description: "Integration type. Possible values: `EMAIL`, `JIRA`, `MS_TEAMS`, `SLACK`",
			Required:    true,
		},
		"jira": {
			Type:        schema.TypeList,
			Description: "No documentation available",
			Optional:    true, // precondition
			Elem:        &schema.Resource{Schema: new(JiraConnection).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
		"ms_teams": {
			Type:        schema.TypeString,
			Description: "Team",
			Optional:    true, // precondition
		},
		"slack_channel": {
			Type:        schema.TypeString,
			Description: "Channel",
			Optional:    true, // precondition
		},
		"url": {
			Type:        schema.TypeString,
			Description: "No documentation available",
			Optional:    true, // nullable & precondition
		},
	}
}

func (me *ContactDetails) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"email":            me.Email,
		"integration_type": me.IntegrationType,
		"jira":             me.Jira,
		"ms_teams":         me.MsTeams,
		"slack_channel":    me.SlackChannel,
		"url":              me.Url,
	})
}

func (me *ContactDetails) HandlePreconditions() error {
	if (me.Email != nil) && (string(me.IntegrationType) != "EMAIL") {
		return fmt.Errorf("'email' must not be specified unless 'integration_type' is set to 'EMAIL'; got 'integration_type'='%v'", me.IntegrationType)
	}
	if (me.Email == nil) && (string(me.IntegrationType) == "EMAIL") {
		return fmt.Errorf("'email' must be specified when 'integration_type' is set to 'EMAIL'; got 'integration_type'='%v'", me.IntegrationType)
	}
	if (me.Jira != nil) && (string(me.IntegrationType) != "JIRA") {
		return fmt.Errorf("'jira' must not be specified unless 'integration_type' is set to 'JIRA'; got 'integration_type'='%v'", me.IntegrationType)
	}
	if (me.Jira == nil) && (string(me.IntegrationType) == "JIRA") {
		return fmt.Errorf("'jira' must be specified when 'integration_type' is set to 'JIRA'; got 'integration_type'='%v'", me.IntegrationType)
	}
	if (me.MsTeams != nil) && (string(me.IntegrationType) != "MS_TEAMS") {
		return fmt.Errorf("'ms_teams' must not be specified unless 'integration_type' is set to 'MS_TEAMS'; got 'integration_type'='%v'", me.IntegrationType)
	}
	if (me.MsTeams == nil) && (string(me.IntegrationType) == "MS_TEAMS") {
		return fmt.Errorf("'ms_teams' must be specified when 'integration_type' is set to 'MS_TEAMS'; got 'integration_type'='%v'", me.IntegrationType)
	}
	if (me.SlackChannel != nil) && (string(me.IntegrationType) != "SLACK") {
		return fmt.Errorf("'slack_channel' must not be specified unless 'integration_type' is set to 'SLACK'; got 'integration_type'='%v'", me.IntegrationType)
	}
	if (me.SlackChannel == nil) && (string(me.IntegrationType) == "SLACK") {
		return fmt.Errorf("'slack_channel' must be specified when 'integration_type' is set to 'SLACK'; got 'integration_type'='%v'", me.IntegrationType)
	}
	if (me.Url != nil) && (!slices.Contains([]string{"SLACK", "JIRA", "MS_TEAMS"}, string(me.IntegrationType))) {
		return fmt.Errorf("'url' must not be specified unless 'integration_type' is one of ['SLACK', 'JIRA', 'MS_TEAMS']; got 'integration_type'='%v'", me.IntegrationType)
	}
	return nil
}

func (me *ContactDetails) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"email":            &me.Email,
		"integration_type": &me.IntegrationType,
		"jira":             &me.Jira,
		"ms_teams":         &me.MsTeams,
		"slack_channel":    &me.SlackChannel,
		"url":              &me.Url,
	})
}
