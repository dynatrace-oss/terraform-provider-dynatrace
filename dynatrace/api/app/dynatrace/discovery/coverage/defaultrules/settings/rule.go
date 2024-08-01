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

package defaultrules

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Rule struct {
	Actions          Actions `json:"actions,omitempty"`
	Category         string  `json:"category"`
	Description      string  `json:"description"`
	EnvironmentScope bool    `json:"environmentScope"` // Environment scope
	ID               string  `json:"id"`
	Priority         string  `json:"priority"`
	Query            string  `json:"query"` // Rule query
	Title            string  `json:"title"`
}

func (me *Rule) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"actions": {
			Type:        schema.TypeList,
			Description: "no documentation available",
			Optional:    true, // minobjects == 0
			Elem:        &schema.Resource{Schema: new(Actions).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
		"category": {
			Type:        schema.TypeString,
			Description: "no documentation available",
			Required:    true,
		},
		"description": {
			Type:        schema.TypeString,
			Description: "no documentation available",
			Required:    true,
		},
		"environment_scope": {
			Type:        schema.TypeBool,
			Description: "Environment scope",
			Required:    true,
		},
		"id": {
			Type:        schema.TypeString,
			Description: "no documentation available",
			Required:    true,
		},
		"priority": {
			Type:        schema.TypeString,
			Description: "no documentation available",
			Required:    true,
		},
		"query": {
			Type:        schema.TypeString,
			Description: "Rule query",
			Required:    true,
		},
		"title": {
			Type:        schema.TypeString,
			Description: "no documentation available",
			Required:    true,
		},
	}
}

func (me *Rule) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"actions":           me.Actions,
		"category":          me.Category,
		"description":       me.Description,
		"environment_scope": me.EnvironmentScope,
		"id":                me.ID,
		"priority":          me.Priority,
		"query":             me.Query,
		"title":             me.Title,
	})
}

func (me *Rule) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"actions":           &me.Actions,
		"category":          &me.Category,
		"description":       &me.Description,
		"environment_scope": &me.EnvironmentScope,
		"id":                &me.ID,
		"priority":          &me.Priority,
		"query":             &me.Query,
		"title":             &me.Title,
	})
}
