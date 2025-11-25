/**
* @license
* Copyright 2025 Dynatrace LLC
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

package policies

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Policy struct {
	Name           string   `json:"name"`
	Tags           []string `json:"tags"`
	Description    string   `json:"description,omitempty"`
	StatementQuery string   `json:"statementQuery"`
	Cluster        string   `json:"-"`
	Environment    string   `json:"-"`
}

func (me *Policy) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The name of the policy",
		},
		"description": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "An optional description text for the policy",
		},
		"tags": {
			Type:        schema.TypeSet,
			Optional:    true,
			MinItems:    1,
			Description: "Tags for this policy",
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
		"statement_query": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The Statement Query of the policy",
		},
		"cluster": {
			Type:          schema.TypeString,
			Optional:      true,
			ConflictsWith: []string{"environment"},
			AtLeastOneOf:  []string{"environment", "cluster"},
			ForceNew:      true,
			Description:   "The UUID of the cluster in case the policy should be applied to all environments of this cluster.",
		},
		"environment": {
			Type:          schema.TypeString,
			Optional:      true,
			ConflictsWith: []string{"cluster"},
			AtLeastOneOf:  []string{"environment", "cluster"},
			ForceNew:      true,
			Description:   "The ID of the environment if the policy should be applied to a specific environment",
		},
	}
}

func (me *Policy) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"name":            me.Name,
		"description":     me.Description,
		"statement_query": me.StatementQuery,
		"cluster":         me.Cluster,
		"environment":     me.Environment,
		"tags":            me.Tags,
	})
}

func (me *Policy) UnmarshalHCL(decoder hcl.Decoder) error {
	if err := decoder.DecodeAll(map[string]any{
		"name":            &me.Name,
		"description":     &me.Description,
		"statement_query": &me.StatementQuery,
		"cluster":         &me.Cluster,
		"environment":     &me.Environment,
		"tags":            &me.Tags,
	}); err != nil {
		return err
	}
	if me.Tags == nil {
		me.Tags = []string{}
	}
	return nil
}
