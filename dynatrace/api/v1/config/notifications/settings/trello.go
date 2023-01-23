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

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/opt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// TrelloConfig Configuration of the Trello notification.
type TrelloConfig struct {
	BaseNotificationConfig
	ResolvedListID     string  `json:"resolvedListId"`               // The Trello list to which the card of the resolved problem should be assigned.
	Text               string  `json:"text"`                         // The text of the generated Trello card.  You can use the following placeholders:  * `{ImpactedEntity}`: The entity impacted by the problem or *X* impacted entities.  * `{PID}`: The ID of the reported problem.  * `{ProblemDetailsMarkdown}`: All problem event details, including root cause, as a [Markdown-formatted](https://www.markdownguide.org/cheat-sheet/) string.  * `{ProblemID}`: The display number of the reported problem.  * `{ProblemImpact}`: The [impact level](https://www.dynatrace.com/support/help/shortlink/impact-analysis) of the problem. Possible values are `APPLICATION`, `SERVICE`, and `INFRASTRUCTURE`.  * `{ProblemSeverity}`: The [severity level](https://www.dynatrace.com/support/help/shortlink/event-types) of the problem. Possible values are `AVAILABILITY`, `ERROR`, `PERFORMANCE`, `RESOURCE_CONTENTION`, and `CUSTOM_ALERT`.  * `{ProblemTitle}`: A short description of the problem.  * `{ProblemURL}`: The URL of the problem within Dynatrace.  * `{State}`: The state of the problem. Possible values are `OPEN` and `RESOLVED`.  * `{Tags}`: The list of tags that are defined for all impacted entities, separated by commas.
	ApplicationKey     string  `json:"applicationKey"`               // The application key for the Trello account.
	AuthorizationToken *string `json:"authorizationToken,omitempty"` // The application token for the Trello account.
	BoardID            string  `json:"boardId"`                      // The Trello board to which the card should be assigned.
	Description        string  `json:"description"`                  // The description of the Trello card.   You can use same placeholders as in card text.
	ListID             string  `json:"listId"`                       // The Trello list to which the card should be assigned.
}

func (me *TrelloConfig) GetType() Type {
	return Types.Trello
}

func (me *TrelloConfig) Schema() map[string]*schema.Schema {
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
		"resolved_list_id": {
			Type:        schema.TypeString,
			Description: "The Trello list to which the card of the resolved problem should be assigned",
			Required:    true,
		},
		"text": {
			Type:        schema.TypeString,
			Description: "The text of the generated Trello card.  You can use the following placeholders:  * `{ImpactedEntity}`: The entity impacted by the problem or *X* impacted entities.  * `{PID}`: The ID of the reported problem.  * `{ProblemDetailsMarkdown}`: All problem event details, including root cause, as a [Markdown-formatted](https://www.markdownguide.org/cheat-sheet/) string.  * `{ProblemID}`: The display number of the reported problem.  * `{ProblemImpact}`: The [impact level](https://www.dynatrace.com/support/help/shortlink/impact-analysis) of the problem. Possible values are `APPLICATION`, `SERVICE`, and `INFRASTRUCTURE`.  * `{ProblemSeverity}`: The [severity level](https://www.dynatrace.com/support/help/shortlink/event-types) of the problem. Possible values are `AVAILABILITY`, `ERROR`, `PERFORMANCE`, `RESOURCE_CONTENTION`, and `CUSTOM_ALERT`.  * `{ProblemTitle}`: A short description of the problem.  * `{ProblemURL}`: The URL of the problem within Dynatrace.  * `{State}`: The state of the problem. Possible values are `OPEN` and `RESOLVED`.  * `{Tags}`: The list of tags that are defined for all impacted entities, separated by commas",
			Required:    true,
		},
		"application_key": {
			Type:        schema.TypeString,
			Description: "The application key for the Trello account",
			Required:    true,
		},
		"authorization_token": {
			Type:        schema.TypeString,
			Description: "The application token for the Trello account",
			Optional:    true,
		},
		"board_id": {
			Type:        schema.TypeString,
			Description: "The Trello board to which the card should be assigned",
			Required:    true,
		},
		"description": {
			Type:        schema.TypeString,
			Description: "The description of the Trello card.   You can use same placeholders as in card text",
			Required:    true,
		},
		"list_id": {
			Type:        schema.TypeString,
			Description: "The Trello list to which the card should be assigned",
			Required:    true,
		},
	}
}

func (me *TrelloConfig) PrepareMarshalHCL(decoder hcl.Decoder) error {
	if value, ok := decoder.GetOk("trello.0.authorization_token"); ok && len(value.(string)) > 0 {
		me.AuthorizationToken = opt.NewString(value.(string))
	}
	return nil
}

func (me *TrelloConfig) FillDemoValues() []string {
	me.AuthorizationToken = opt.NewString("#######")
	return []string{"The REST API didn't provide the credentials"}
}

func (me *TrelloConfig) MarshalHCL(properties hcl.Properties) error {
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
	if err := properties.Encode("resolved_list_id", me.ResolvedListID); err != nil {
		return err
	}
	if err := properties.Encode("text", me.Text); err != nil {
		return err
	}
	if err := properties.Encode("application_key", me.ApplicationKey); err != nil {
		return err
	}
	if err := properties.Encode("authorization_token", me.AuthorizationToken); err != nil {
		return err
	}
	if err := properties.Encode("board_id", me.BoardID); err != nil {
		return err
	}
	if err := properties.Encode("description", me.Description); err != nil {
		return err
	}
	if err := properties.Encode("list_id", me.ListID); err != nil {
		return err
	}
	return nil
}

func (me *TrelloConfig) UnmarshalHCL(decoder hcl.Decoder) error {
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
		delete(me.Unknowns, "resolved_list_id")
		delete(me.Unknowns, "text")
		delete(me.Unknowns, "application_key")
		delete(me.Unknowns, "authorization_token")
		delete(me.Unknowns, "board_id")
		delete(me.Unknowns, "description")
		delete(me.Unknowns, "list_id")
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
	if value, ok := decoder.GetOk("resolved_list_id"); ok {
		me.ResolvedListID = value.(string)
	}
	if value, ok := decoder.GetOk("text"); ok {
		me.Text = value.(string)
	}
	if value, ok := decoder.GetOk("application_key"); ok {
		me.ApplicationKey = value.(string)
	}
	if value, ok := decoder.GetOk("authorization_token"); ok {
		me.AuthorizationToken = opt.NewString(value.(string))
	}
	if value, ok := decoder.GetOk("board_id"); ok {
		me.BoardID = value.(string)
	}
	if value, ok := decoder.GetOk("description"); ok {
		me.Description = value.(string)
	}
	if value, ok := decoder.GetOk("list_id"); ok {
		me.ListID = value.(string)
	}
	return nil
}

func (me *TrelloConfig) MarshalJSON() ([]byte, error) {
	properties := xjson.NewProperties(me.Unknowns)
	if err := properties.MarshalAll(map[string]any{
		"id":                 me.ID,
		"name":               me.Name,
		"type":               me.GetType(),
		"active":             me.Active,
		"alertingProfile":    me.AlertingProfile,
		"resolvedListId":     me.ResolvedListID,
		"text":               me.Text,
		"applicationKey":     me.ApplicationKey,
		"authorizationToken": me.AuthorizationToken,
		"boardId":            me.BoardID,
		"description":        me.Description,
		"listId":             me.ListID,
	}); err != nil {
		return nil, err
	}
	return json.Marshal(properties)
}

func (me *TrelloConfig) UnmarshalJSON(data []byte) error {
	properties := xjson.Properties{}
	if err := json.Unmarshal(data, &properties); err != nil {
		return err
	}
	if err := properties.UnmarshalAll(map[string]any{
		"id":                 &me.ID,
		"name":               &me.Name,
		"type":               &me.Type,
		"active":             &me.Active,
		"alertingProfile":    &me.AlertingProfile,
		"resolvedListId":     &me.ResolvedListID,
		"text":               &me.Text,
		"applicationKey":     &me.ApplicationKey,
		"authorizationToken": &me.AuthorizationToken,
		"boardId":            &me.BoardID,
		"description":        &me.Description,
		"listId":             &me.ListID,
	}); err != nil {
		return err
	}
	if len(properties) > 0 {
		me.Unknowns = properties
	}
	return nil
}
