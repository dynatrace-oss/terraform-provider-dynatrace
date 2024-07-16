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

package monitors

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Steps []*Step

func (me *Steps) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"step": {
			Type:        schema.TypeList,
			Description: "The step of a network availability monitor",
			Required:    true,
			Elem:        &schema.Resource{Schema: new(Step).Schema()},
		},
	}
}

func (me Steps) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeSlice("step", me)
}

func (me *Steps) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeSlice("step", me)
}

type Step struct {
	Name                  string                `json:"name"`                   // Step name
	RequestType           RequestType           `json:"requestType"`            // Request type, possible values `ICMP`, `TCP`, `DNS`
	TargetList            []string              `json:"targetList"`             // Target list
	TargetFilter          *string               `json:"targetFilter,omitempty"` // Target filter
	Properties            map[string]string     `json:"properties"`             // Key/value pairs of properties which apply to all requests in the step
	Constraints           Constraints           `json:"constraints"`            // The list of constraints which apply to all requests in the step
	RequestConfigurations RequestConfigurations `json:"requestConfigurations"`  // Request configurations
}

func (me *Step) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:        schema.TypeString,
			Description: "Step name",
			Required:    true,
		},
		"request_type": {
			Type:        schema.TypeString,
			Description: "Request type, possible values `ICMP`, `TCP`, `DNS`",
			Required:    true,
		},
		"target_list": {
			Type:        schema.TypeSet,
			Description: "Target list",
			Required:    true,
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
		"target_filter": {
			Type:        schema.TypeString,
			Description: "Target filter",
			Optional:    true,
		},
		"properties": {
			Type:        schema.TypeMap,
			Description: "Key/value pairs of properties which apply to all requests in the step",
			Required:    true,
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
		"constraints": {
			Type:        schema.TypeList,
			Description: "The list of constraints which apply to all requests in the step",
			Required:    true,
			Elem:        &schema.Resource{Schema: new(Constraints).Schema()},
		},
		"request_configurations": {
			Type:        schema.TypeList,
			Description: "Request configurations",
			Optional:    true,
			Elem:        &schema.Resource{Schema: new(RequestConfigurations).Schema()},
		},
	}
}

func (me Step) MarshalHCL(properties hcl.Properties) error {
	if me.RequestConfigurations != nil && len(me.RequestConfigurations) == 1 && len(me.RequestConfigurations[0].Constraints) == 0 {
		me.RequestConfigurations = RequestConfigurations{}
	}

	return properties.EncodeAll(map[string]any{
		"name":                   me.Name,
		"request_type":           me.RequestType,
		"target_list":            me.TargetList,
		"target_filter":          me.TargetFilter,
		"properties":             me.Properties,
		"constraints":            me.Constraints,
		"request_configurations": me.RequestConfigurations,
	})
}

func (me *Step) UnmarshalHCL(decoder hcl.Decoder) error {
	err := decoder.DecodeAll(map[string]any{
		"name":                   &me.Name,
		"request_type":           &me.RequestType,
		"target_list":            &me.TargetList,
		"target_filter":          &me.TargetFilter,
		"properties":             &me.Properties,
		"constraints":            &me.Constraints,
		"request_configurations": &me.RequestConfigurations,
	})

	if len(me.RequestConfigurations) == 0 {
		me.RequestConfigurations = RequestConfigurations{&RequestConfiguration{}}
	}

	return err
}
