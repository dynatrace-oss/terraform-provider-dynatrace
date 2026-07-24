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
	"fmt"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// CustomNamespace models an entry of `aws.namespaces` on the wire. Lets users
// ingest CloudWatch metrics for namespaces not covered by the built-in feature
// sets, including non-standard custom namespaces.
type CustomNamespace struct {
	Namespace            string
	AutoDiscoveryEnabled bool
	Metrics              []*CustomMetric
}

type CustomNamespaces []*CustomNamespace

func (me *CustomNamespace) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"namespace": {
			Type:        schema.TypeString,
			Description: "CloudWatch namespace, e.g. `AWS/GroundStation` for a standard service or `MyApp/Metrics` for a custom namespace.",
			Required:    true,
		},
		"auto_discovery_enabled": {
			Type:        schema.TypeBool,
			Description: "Whether the extension auto-discovers entities in this namespace. Default false.",
			Optional:    true,
			Default:     false,
		},
		"metric": {
			Type:        schema.TypeList,
			Description: "CloudWatch metrics to ingest from this namespace.",
			Required:    true,
			MinItems:    1,
			Elem:        &schema.Resource{Schema: new(CustomMetric).Schema()},
		},
	}
}

func (me *CustomNamespace) MarshalHCL(properties hcl.Properties) error {
	if err := properties.EncodeAll(map[string]any{
		"namespace":              me.Namespace,
		"auto_discovery_enabled": me.AutoDiscoveryEnabled,
	}); err != nil {
		return err
	}
	return properties.EncodeSlice("metric", me.Metrics)
}

func (me *CustomNamespace) UnmarshalHCL(decoder hcl.Decoder) error {
	if err := decoder.DecodeAll(map[string]any{
		"namespace":              &me.Namespace,
		"auto_discovery_enabled": &me.AutoDiscoveryEnabled,
	}); err != nil {
		return err
	}
	return decoder.DecodeSlice("metric", &me.Metrics)
}

// CustomMetric is one CloudWatch metric inside a CustomNamespace.
type CustomMetric struct {
	Name         string
	Unit         string
	Dimensions   []string
	Aggregations []string
	// Type is either CUSTOM_AWS (standard CloudWatch namespace, e.g. AWS/*) or
	// CUSTOM (non-AWS namespace published by your own application).
	Type string
}

func (me *CustomMetric) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:        schema.TypeString,
			Description: "CloudWatch metric name.",
			Required:    true,
		},
		"unit": {
			Type:        schema.TypeString,
			Description: "CloudWatch metric unit (`Count`, `Bytes`, `Seconds`, `Percent`, …).",
			Required:    true,
		},
		"dimensions": {
			Type:        schema.TypeList,
			Description: "CloudWatch dimensions to collect for this metric.",
			Optional:    true,
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
		"aggregations": {
			Type:        schema.TypeList,
			Description: "Statistics to retrieve (`Sum`, `Average`, `Maximum`, `Minimum`, `SampleCount`).",
			Required:    true,
			MinItems:    1,
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
		"type": {
			Type:        schema.TypeString,
			Description: "`CUSTOM_AWS` for standard AWS/* namespaces or `CUSTOM` for non-AWS namespaces published by your own app.",
			Required:    true,
			ValidateFunc: func(i any, k string) (warnings []string, errs []error) {
				v, _ := i.(string)
				if v != "CUSTOM_AWS" && v != "CUSTOM" {
					errs = append(errs, fmt.Errorf("%s must be CUSTOM_AWS or CUSTOM, got %q", k, v))
				}
				return
			},
		},
	}
}

func (me *CustomMetric) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"name":         me.Name,
		"unit":         me.Unit,
		"dimensions":   me.Dimensions,
		"aggregations": me.Aggregations,
		"type":         me.Type,
	})
}

func (me *CustomMetric) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"name":         &me.Name,
		"unit":         &me.Unit,
		"dimensions":   &me.Dimensions,
		"aggregations": &me.Aggregations,
		"type":         &me.Type,
	})
}
