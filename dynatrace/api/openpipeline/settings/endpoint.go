package openpipeline

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Endpoints struct {
	Endpoints []EndpointDefinition `json:"endpoints"`
}

func (ep *Endpoints) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"endpoint": {
			Type:        schema.TypeSet,
			Description: "todo",
			Elem:        &schema.Resource{Schema: new(EndpointDefinition).Schema()},
			Optional:    true,
		},
	}
}

func (ep *Endpoints) MarshalHCL(properties hcl.Properties) error {
	return properties.Encode("endpoints", ep.Endpoints)
}

func (ep *Endpoints) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.Decode("endpoints", &ep.Endpoints)
}

type EndpointDefinition struct {
	BasePath string `json:"basePath"`
	//Builtin       *bool   `json:"builtin,omitempty"`
	//DefaultBucket *string `json:"defaultBucket,omitempty"`
	//DisplayName   *string `json:"displayName,omitempty"`
	//Editable      *bool   `json:"editable,omitempty"`
	//Enabled       bool    `json:"enabled"`
	Processors EndpointProcessors `json:"processors,omitempty"`
	// Routing Routing `json:"routing"`
	//Segment string  `json:"segment"`
}

func (d *EndpointDefinition) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"base_path": {
			Type:        schema.TypeString,
			Description: "todo",
			Required:    true,
		},
		"processors": {
			Type:        schema.TypeList,
			Description: "todo",
			MinItems:    1,
			MaxItems:    1,
			Elem:        &schema.Resource{Schema: new(EndpointProcessors).Schema()},
			Optional:    true,
		},
	}
}

func (d *EndpointDefinition) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"base_path":  d.BasePath,
		"processors": d.Processors,
	})
}

func (d *EndpointDefinition) UnmarshalHCL(decoder hcl.Decoder) error {

	return decoder.DecodeAll(map[string]any{
		"base_path":          &d.BasePath,
		"endpoint_processor": &d.Processors,
	})
}

type EndpointProcessors struct {
	Processors []EndpointProcessor
}

func (ep *EndpointProcessors) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"processor": {
			Type:        schema.TypeSet,
			Description: "todo",
			Elem:        &schema.Resource{Schema: new(EndpointProcessor).Schema()},
			Optional:    true,
		},
	}
}

func (ep *EndpointProcessors) MarshalHCL(properties hcl.Properties) error {
	return properties.Encode("processors", ep.Processors)
}

func (ep *EndpointProcessors) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.Decode("processors", &ep.Processors)
}
