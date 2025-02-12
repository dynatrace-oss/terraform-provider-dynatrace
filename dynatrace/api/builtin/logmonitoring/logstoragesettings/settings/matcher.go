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

package logstoragesettings

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Matchers []*Matcher

func (me *Matchers) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"matcher": {
			Type:        schema.TypeSet,
			Required:    true,
			MinItems:    1,
			Description: "",
			Elem:        &schema.Resource{Schema: new(Matcher).Schema()},
		},
	}
}

func (me Matchers) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeSlice("matcher", me)
}

func (me *Matchers) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeSlice("matcher", me)
}

type Matcher struct {
	Attribute MatcherType `json:"attribute"` // Possible Values: `Container_name`, `Dt_entity_container_group`, `Dt_entity_process_group`, `Host_tag`, `Journald_unit`, `K8s_container_name`, `K8s_deployment_name`, `K8s_namespace_name`, `K8s_pod_annotation`, `K8s_pod_label`, `K8s_workload_kind`, `K8s_workload_name`, `Log_content`, `Log_source`, `Log_source_origin`, `Loglevel`, `Process_technology`, `Winlog_eventid`, `Winlog_keywords`, `Winlog_opcode`, `Winlog_provider`, `Winlog_task`, `Winlog_username`
	Operator  Operator    `json:"operator"`  // Possible Values: `MATCHES`
	Values    []string    `json:"values"`
}

func (me *Matcher) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"attribute": {
			Type:        schema.TypeString,
			Description: "Possible Values: `Container_name`, `Dt_entity_container_group`, `Dt_entity_process_group`, `Host_tag`, `Journald_unit`, `K8s_container_name`, `K8s_deployment_name`, `K8s_namespace_name`, `K8s_pod_annotation`, `K8s_pod_label`, `K8s_workload_kind`, `K8s_workload_name`, `Log_content`, `Log_source`, `Log_source_origin`, `Loglevel`, `Process_technology`, `Winlog_eventid`, `Winlog_keywords`, `Winlog_opcode`, `Winlog_provider`, `Winlog_task`, `Winlog_username`",
			Required:    true,
		},
		"operator": {
			Type:        schema.TypeString,
			Description: "Possible Values: `MATCHES`",
			Required:    true,
		},
		"values": {
			Type:        schema.TypeSet,
			Description: "no documentation available",
			Required:    true,
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
	}
}

func (me *Matcher) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"attribute": me.Attribute,
		"operator":  me.Operator,
		"values":    me.Values,
	})
}

func (me *Matcher) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"attribute": &me.Attribute,
		"operator":  &me.Operator,
		"values":    &me.Values,
	})
}
