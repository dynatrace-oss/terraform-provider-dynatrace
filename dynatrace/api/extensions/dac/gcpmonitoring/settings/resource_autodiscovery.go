/**
* @license
* Copyright 2026 Dynatrace LLC
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
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ResourceAutodiscovery is a per-resource-type override of the extension's
// default autodiscovery behaviour.
type ResourceAutodiscovery struct {
	ResourceType         string
	AutoDiscoveryEnabled bool
	ExcludeMetricType    []string
}

type ResourceAutodiscoveries []*ResourceAutodiscovery

func (me *ResourceAutodiscovery) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"resource_type": {
			Type:        schema.TypeString,
			Description: "GCP monitored resource type in the form `<service>.googleapis.com/<Kind>`, e.g. `compute.googleapis.com/Instance`.",
			Required:    true,
		},
		"auto_discovery_enabled": {
			Type:        schema.TypeBool,
			Description: "Whether autodiscovery is enabled for this resource type.",
			Required:    true,
		},
		"exclude_metric_type": {
			Type:        schema.TypeSet,
			Description: "Metric types to exclude from autodiscovery.",
			Optional:    true,
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
	}
}

func (me *ResourceAutodiscovery) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"resource_type":          me.ResourceType,
		"auto_discovery_enabled": me.AutoDiscoveryEnabled,
		"exclude_metric_type":    me.ExcludeMetricType,
	})
}

func (me *ResourceAutodiscovery) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"resource_type":          &me.ResourceType,
		"auto_discovery_enabled": &me.AutoDiscoveryEnabled,
		"exclude_metric_type":    &me.ExcludeMetricType,
	})
}
