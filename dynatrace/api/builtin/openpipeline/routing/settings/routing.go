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

package settings

import (
	"fmt"
	"strings"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/openpipeline"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

var AllowedKinds = []string{
	"logs", "events", "events.security", "security.events", "bizevents", "spans",
	"events.sdlc", "metrics", "usersessions", "davis.problems", "davis.events",
	"system.events", "azure.logs.forwarding", "user.events",
}

type Routing struct {
	Kind           string          `json:"-"`
	RoutingEntries []*RoutingEntry `json:"routingEntries,omitempty"`
}

func (r *Routing) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"kind": {
			Type:         schema.TypeString,
			Required:     true,
			Description:  "Indicates OpenPipeline data source",
			ForceNew:     true,
			ValidateFunc: validation.StringInSlice(AllowedKinds, true),
		},
		"routing_entry": {
			Type:        schema.TypeList,
			Description: "Groups all entries of the routing table together, mapping ingest sources to pipelines",
			Elem:        &schema.Resource{Schema: new(RoutingEntry).Schema()},
			Optional:    true,
			MaxItems:    3000,
		},
	}
}

func (r *Routing) MarshalHCL(properties hcl.Properties) error {
	err := properties.Encode("kind", r.Kind)
	if err != nil {
		return err
	}

	return properties.EncodeSlice("routing_entry", r.RoutingEntries)
}

func (r *Routing) UnmarshalHCL(decoder hcl.Decoder) error {
	err := decoder.Decode("kind", &r.Kind)
	if err != nil {
		return err
	}

	return decoder.DecodeSlice("routing_entry", &r.RoutingEntries)
}

func (r *Routing) Name() string {
	return "Routing for pipelines"
}

type RoutingEntry struct {
	Enabled           bool    `json:"enabled"`
	PipelineType      string  `json:"pipelineType"`
	BuiltinPipelineID *string `json:"builtinPipelineId,omitempty"`
	PipelineID        *string `json:"pipelineId,omitempty"`
	Matcher           string  `json:"matcher"`
	Description       string  `json:"description"`
}

const BuiltinPipelineIDMaxLength = 500
const DescriptionMaxLength = 512
const MatcherMaxLength = 1500

func (re *RoutingEntry) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"enabled": {
			Type:        schema.TypeBool,
			Description: "Indicates if the routing entry is active",
			Default:     true,
			Optional:    true,
		},
		"pipeline_id": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "ID of the pipeline. Only used if the pipeline type is \"custom\"",
		},
		"builtin_pipeline_id": {
			Type:         schema.TypeString,
			Optional:     true,
			Description:  "ID of the pipeline. Only used if the pipeline type is \"builtin\"",
			ValidateFunc: validation.StringLenBetween(1, BuiltinPipelineIDMaxLength),
		},
		"pipeline_type": {
			Type:         schema.TypeString,
			Required:     true,
			Description:  "Type of the pipeline. Must be \"custom\" or \"builtin\"",
			ValidateFunc: validation.StringInSlice([]string{"custom", "builtin"}, true),
		},
		"matcher": {
			Type:         schema.TypeString,
			Description:  "Query which determines whether the record should be routed to the target pipeline of this rule",
			Required:     true,
			ValidateFunc: validation.StringLenBetween(1, MatcherMaxLength),
		},
		"description": {
			Type:        schema.TypeString,
			Description: "Description of the routing table entry. Must not start with 'dt.' or 'dynatrace.'",
			Required:    true,
			ValidateFunc: validation.All(
				validation.StringLenBetween(1, DescriptionMaxLength),
				func(input interface{}, schema string) (warnings []string, errors []error) {
					id, ok := input.(string)
					if !ok {
						errors = append(errors, fmt.Errorf("expected type of %s to be string", schema))
						return warnings, errors
					}

					if strings.HasPrefix(id, "dt.") || strings.HasPrefix(id, "dynatrace.") {
						errors = append(errors,
							fmt.Errorf("%s must not start with 'dt.' or 'dynatrace.'", schema))
					}
					return warnings, errors
				}),
		},
	}
}

func (re *RoutingEntry) MarshalHCL(properties hcl.Properties) error {
	err := properties.EncodeAll(map[string]any{
		"enabled":             re.Enabled,
		"pipeline_id":         re.PipelineID,
		"builtin_pipeline_id": re.BuiltinPipelineID,
		"pipeline_type":       re.PipelineType,
		"matcher":             re.Matcher,
		"description":         re.Description,
	})
	openpipeline.RemoveNils(properties)

	return err
}

func (re *RoutingEntry) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"enabled":             &re.Enabled,
		"pipeline_id":         &re.PipelineID,
		"builtin_pipeline_id": &re.BuiltinPipelineID,
		"pipeline_type":       &re.PipelineType,
		"matcher":             &re.Matcher,
		"description":         &re.Description,
	})
}
