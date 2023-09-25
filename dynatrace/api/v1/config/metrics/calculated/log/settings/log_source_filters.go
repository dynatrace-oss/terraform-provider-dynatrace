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

package log

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type LogSourceFilters []*LogSourceFilter

func (me *LogSourceFilters) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"filter": {
			Type:        schema.TypeList,
			Optional:    true,
			MinItems:    1,
			Description: "A list of filters to define the logs to look into. If several criteria are specified, the AND logic applies.",
			Elem:        &schema.Resource{Schema: new(LogSourceFilter).Schema()},
		},
	}
}

func (me LogSourceFilters) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeSlice("filter", me)
}

func (me *LogSourceFilters) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeSlice("filter", me)
}

// LogSourceFilters Parameters of a filter of a calculated log metric.
type LogSourceFilter struct {
	PathDefinitions PathDefinitions `json:"pathDefinitions,omitempty"` // A list of filtering criteria for log path. If several criteria are specified, the OR logic applies.
	SourceEntities  *[]string       `json:"sourceEntities,omitempty"`  // A list of Dynatrace entities, where the log can originate from. Specify Dynatrace entity IDs here. Allowed types of entities are PROCESS_GROUP and PROCESS_GROUP_INSTANCE. You can't mix entity types. If several entities are specified, the OR logic applies. This field is mutually exclusive with the osTypes field.
	HostFilters     *[]string       `json:"hostFilters,omitempty"`     // A list of hosts, where the log can originate from. Specify Dynatrace entity IDs here.
	OSTypes         *[]OSType       `json:"osTypes,omitempty"`         // A list of operating system types, where the log can originate from. If set, only OS logs are included in the result. If several OS are specified, the OR logic applies. This field is mutually exclusive with the sourceEntities field.
}

func (me *LogSourceFilter) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"path_definitions": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "A list of filtering criteria for log path. If several criteria are specified, the OR logic applies.",
			Elem:        &schema.Resource{Schema: new(PathDefinitions).Schema()},
		},
		"source_entities": {
			Type:        schema.TypeSet,
			Optional:    true,
			Description: "A list of Dynatrace entities, where the log can originate from. Specify Dynatrace entity IDs here. Allowed types of entities are PROCESS_GROUP and PROCESS_GROUP_INSTANCE. You can't mix entity types. If several entities are specified, the OR logic applies. This field is mutually exclusive with the osTypes field.",
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
		"host_filters": {
			Type:        schema.TypeSet,
			Optional:    true,
			Description: "A list of hosts, where the log can originate from. Specify Dynatrace entity IDs here.",
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
		"os_types": {
			Type:        schema.TypeSet,
			Optional:    true,
			Description: "A list of operating system types, where the log can originate from. If set, only OS logs are included in the result. If several OS are specified, the OR logic applies. This field is mutually exclusive with the sourceEntities field.",
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
	}
}

func (me *LogSourceFilter) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"path_definitions": me.PathDefinitions,
		"source_entities":  me.SourceEntities,
		"host_filters":     me.HostFilters,
		"os_types":         me.OSTypes,
	})
}

func (me *LogSourceFilter) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"path_definitions": &me.PathDefinitions,
		"source_entities":  &me.SourceEntities,
		"host_filters":     &me.HostFilters,
		"os_types":         &me.OSTypes,
	})
}

type OSType string

var OSTypes = struct {
	AIX     OSType
	DARWIN  OSType
	HPUX    OSType
	LINUX   OSType
	SOLARIS OSType
	WINDOWS OSType
	ZOS     OSType
}{
	"AIX",
	"DARWIN",
	"HPUX",
	"LINUX",
	"SOLARIS",
	"WINDOWS",
	"ZOS",
}
