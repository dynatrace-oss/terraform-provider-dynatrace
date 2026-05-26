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

package custominjectionrules

import (
	"fmt"
	"slices"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Settings struct {
	ApplicationID string   `json:"-" scope:"applicationId"` // The scope of this settings. If the settings should cover the whole environment, just don't specify any scope.
	Enabled       bool     `json:"enabled"`                 // This setting is enabled (`true`) or disabled (`false`)
	HtmlPattern   *string  `json:"htmlPattern,omitempty"`
	Operator      Operator `json:"operator"`             // **Example**: \n\n   **For the URL:**  \n  `http://www.example.com:8080/lorem/ipsum.jsp?mode=desktop` \n\n   A rule can be specified on the URL pattern:  \n  `/lorem/ipsum.jsp` \n\n   Using the operator:  \n  `URL ends with` \n\n   **Result:**  \n  If URL ends with .jsp do not inject the JavaScript library. Possible values: `AllPages`, `Contains`, `Ends`, `Equals`, `Starts`
	Rule          Rule     `json:"rule"`                 // Rule. Possible values: `AfterSpecificHtml`, `Automatic`, `BeforeSpecificHtml`, `DoNotInject`
	UrlPattern    *string  `json:"urlPattern,omitempty"` // URL pattern
	InsertAfter   string   `json:"-"`
}

func (me *Settings) Name() string {
	return me.ApplicationID
}

func (me *Settings) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"application_id": {
			Type:        schema.TypeString,
			Description: "The scope of this settings. If the settings should cover the whole environment, just don't specify any scope.",
			Required:    true,
		},
		"enabled": {
			Type:        schema.TypeBool,
			Description: "This setting is enabled (`true`) or disabled (`false`)",
			Required:    true,
		},
		"html_pattern": {
			Type:        schema.TypeString,
			Description: "No documentation available",
			Optional:    true, // precondition
		},
		"operator": {
			Type:        schema.TypeString,
			Description: "**Example**: \n\n   **For the URL:**  \n  `http://www.example.com:8080/lorem/ipsum.jsp?mode=desktop` \n\n   A rule can be specified on the URL pattern:  \n  `/lorem/ipsum.jsp` \n\n   Using the operator:  \n  `URL ends with` \n\n   **Result:**  \n  If URL ends with .jsp do not inject the JavaScript library. Possible values: `AllPages`, `Contains`, `Ends`, `Equals`, `Starts`",
			Required:    true,
		},
		"rule": {
			Type:        schema.TypeString,
			Description: "Rule. Possible values: `AfterSpecificHtml`, `Automatic`, `BeforeSpecificHtml`, `DoNotInject`",
			Required:    true,
		},
		"url_pattern": {
			Type:        schema.TypeString,
			Description: "URL pattern",
			Optional:    true, // precondition
		},
		"insert_after": {
			Type:        schema.TypeString,
			Description: "Because this resource allows for ordering you may specify the ID of the resource instance that comes before this instance regarding order. If not specified when creating the setting will be added to the end of the list. If not specified during update the order will remain untouched",
			Optional:    true,
			Computed:    true,
		},
	}
}

func (me *Settings) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"application_id": me.ApplicationID,
		"enabled":        me.Enabled,
		"html_pattern":   me.HtmlPattern,
		"operator":       me.Operator,
		"rule":           me.Rule,
		"url_pattern":    me.UrlPattern,
		"insert_after":   me.InsertAfter,
	})
}

func (me *Settings) HandlePreconditions() error {
	if (me.HtmlPattern == nil) && (slices.Contains([]string{"BeforeSpecificHtml", "AfterSpecificHtml"}, string(me.Rule))) {
		me.HtmlPattern = new("")
	}
	if (me.UrlPattern == nil) && (slices.Contains([]string{"Equals", "Starts", "Ends", "Contains"}, string(me.Operator))) {
		me.UrlPattern = new("")
	}
	if (me.HtmlPattern != nil) && (!slices.Contains([]string{"BeforeSpecificHtml", "AfterSpecificHtml"}, string(me.Rule))) {
		return fmt.Errorf("'html_pattern' must not be specified unless 'rule' is one of ['BeforeSpecificHtml', 'AfterSpecificHtml']; got 'rule'='%v'", me.Rule)
	}
	if (me.UrlPattern != nil) && (!slices.Contains([]string{"Equals", "Starts", "Ends", "Contains"}, string(me.Operator))) {
		return fmt.Errorf("'url_pattern' must not be specified unless 'operator' is one of ['Equals', 'Starts', 'Ends', 'Contains']; got 'operator'='%v'", me.Operator)
	}
	return nil
}

func (me *Settings) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"application_id": &me.ApplicationID,
		"enabled":        &me.Enabled,
		"html_pattern":   &me.HtmlPattern,
		"operator":       &me.Operator,
		"rule":           &me.Rule,
		"url_pattern":    &me.UrlPattern,
		"insert_after":   &me.InsertAfter,
	})
}
