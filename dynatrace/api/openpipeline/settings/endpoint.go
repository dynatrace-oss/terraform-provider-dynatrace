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

package openpipeline

import (
	"encoding/json"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Endpoints struct {
	Endpoints []*EndpointDefinition
}

func (ep *Endpoints) RemoveFixed() {
	filteredEndpointDefinitions := []*EndpointDefinition{}
	for _, endpoint := range ep.Endpoints {
		if !endpoint.IsFixed() {
			filteredEndpointDefinitions = append(filteredEndpointDefinitions, endpoint)
		}
	}
	ep.Endpoints = filteredEndpointDefinitions
}

func (ep *Endpoints) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"endpoint": {
			Type:        schema.TypeList,
			Description: "Definition of a single ingest source",
			Elem:        &schema.Resource{Schema: new(EndpointDefinition).Schema()},
			Optional:    true,
		},
	}
}

func (ep *Endpoints) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeSlice("endpoint", ep.Endpoints)
}

func (ep *Endpoints) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeSlice("endpoint", &ep.Endpoints)
}

func (d Endpoints) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.Endpoints)
}

func (d *Endpoints) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &d.Endpoints)
}

type EndpointDefinition struct {
	BasePath      string              `json:"basePath"`
	Builtin       *bool               `json:"builtin,omitempty"`
	DefaultBucket *string             `json:"defaultBucket,omitempty"`
	DisplayName   *string             `json:"displayName,omitempty"`
	Editable      *bool               `json:"editable,omitempty"`
	Enabled       bool                `json:"enabled"`
	Segment       string              `json:"segment"`
	Routing       *Routing            `json:"routing"`
	Processors    *EndpointProcessors `json:"-"`
}

func (d *EndpointDefinition) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"default_bucket": {
			Type:        schema.TypeString,
			Description: "The default bucket assigned to records for the ingest source",
			Optional:    true,
		},
		"display_name": {
			Type:        schema.TypeString,
			Description: "Display name of the ingest source",
			Required:    true,
		},
		"enabled": {
			Type:        schema.TypeBool,
			Description: "Indicates if the object is active",
			Required:    true,
		},
		"segment": {
			Type:        schema.TypeString,
			Description: "The segment of the ingest source, which is applied to the base path. Must be unique within a configuration.\"",
			Required:    true,
		},
		"routing": {
			Type:        schema.TypeList,
			Description: "Routing strategy, either dynamic or static",
			MinItems:    1,
			MaxItems:    1,
			Elem:        &schema.Resource{Schema: new(Routing).Schema()},
			Required:    true,
		},
		"processors": {
			Type:        schema.TypeList,
			Description: "The pre-processing done in the ingest source",
			MinItems:    1,
			MaxItems:    1,
			Elem:        &schema.Resource{Schema: new(EndpointProcessors).Schema()},
			Optional:    true,
		},
	}
}

func (d *EndpointDefinition) MarshalHCL(properties hcl.Properties) error {
	if err := properties.EncodeAll(map[string]any{
		"default_bucket": d.DefaultBucket,
		"display_name":   d.DisplayName,
		"enabled":        d.Enabled,
		"segment":        d.Segment,
		"routing":        d.Routing,
	}); err != nil {
		return err
	}

	if d.Processors != nil && len(d.Processors.Processors) > 0 {
		if err := properties.Encode("processors", d.Processors); err != nil {
			return err
		}
	}

	return nil
}

func (d *EndpointDefinition) UnmarshalHCL(decoder hcl.Decoder) error {

	return decoder.DecodeAll(map[string]any{
		"default_bucket": &d.DefaultBucket,
		"display_name":   &d.DisplayName,
		"enabled":        &d.Enabled,
		"segment":        &d.Segment,
		"routing":        &d.Routing,
		"processors":     &d.Processors,
	})
}

func (d EndpointDefinition) MarshalJSON() ([]byte, error) {
	rawProcessors, err := json.Marshal(d.Processors)
	if err != nil {
		return nil, err
	}

	type endpointDefinition EndpointDefinition
	endpointDef := struct {
		endpointDefinition
		RawProcessors json.RawMessage `json:"processors"`
	}{
		endpointDefinition: (endpointDefinition)(d),
		RawProcessors:      rawProcessors,
	}

	return json.Marshal(endpointDef)
}

func (d *EndpointDefinition) UnmarshalJSON(b []byte) error {
	type endpointDefinition EndpointDefinition

	endpointDef := struct {
		endpointDefinition
		RawProcessors json.RawMessage `json:"processors"`
	}{}
	if err := json.Unmarshal(b, &endpointDef); err != nil {
		return err
	}

	*d = EndpointDefinition(endpointDef.endpointDefinition)

	d.Processors = &EndpointProcessors{}
	return json.Unmarshal(endpointDef.RawProcessors, d.Processors)
}

func (d *EndpointDefinition) IsFixed() bool {
	return (d.Builtin != nil) && *d.Builtin
}
