package openpipeline

import (
	"encoding/json"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Endpoints struct {
	Endpoints []*EndpointDefinition `json:"endpoints"`
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

func (d *Endpoints) MarshalJSON() ([]byte, error) {
	m := map[string]json.RawMessage{}
	return json.Marshal(m)
}

func (d *Endpoints) UnmarshalJSON(data []byte) error {
	if err := json.Unmarshal(data, &d.Endpoints); err != nil {
		return err
	}
	return nil
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
		"base_path": {
			Type:        schema.TypeString,
			Description: "The base path of the ingest source.",
			Required:    true,
		},
		"builtin": {
			Type:        schema.TypeString,
			Description: "Indicates if the object is provided by Dynatrace or customer defined.",
			Optional:    true,
		},
		"default_bucket": {
			Type:        schema.TypeString,
			Description: "The default bucket assigned to records for the ingest source.",
			Optional:    true,
		},
		"display_name": {
			Type:        schema.TypeString,
			Description: "Display name of the ingest source.",
			Optional:    true,
		},
		"editable": {
			Type:        schema.TypeBool,
			Description: "Indicates if the user is allowed to edit this object based on permissions and builtin property.",
			Optional:    true,
		},
		"enabled": {
			Type:        schema.TypeBool,
			Description: "Indicates if the object is active.",
			Required:    true,
		},
		"segment": {
			Type:        schema.TypeString,
			Description: "The segment of the ingest source, which is applied to the base path. Must be unique within a configuration.\"",
			Required:    true,
		},
		"routing": {
			Type:        schema.TypeList,
			Description: "Routing strategy, either dynamic or static.",
			MinItems:    1,
			MaxItems:    1,
			Elem:        &schema.Resource{Schema: new(Routing).Schema()},
			Optional:    true,
		},
		"processors": {
			Type:        schema.TypeList,
			Description: "The pre-processing done in the ingest source.",
			MinItems:    1,
			MaxItems:    1,
			Elem:        &schema.Resource{Schema: new(EndpointProcessors).Schema()},
			Optional:    true,
		},
	}
}

func (d *EndpointDefinition) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"base_path":      d.BasePath,
		"builtin":        d.Builtin,
		"default_bucket": d.DefaultBucket,
		"display_name":   d.DisplayName,
		"editable":       d.Editable,
		"enabled":        d.Enabled,
		"segment":        d.Segment,
		"routing":        d.Routing,
		"processors":     d.Processors,
	})
}

func (d *EndpointDefinition) UnmarshalHCL(decoder hcl.Decoder) error {

	return decoder.DecodeAll(map[string]any{
		"base_path":          &d.BasePath,
		"builtin":            &d.Builtin,
		"default_bucket":     &d.DefaultBucket,
		"display_name":       &d.DisplayName,
		"editable":           &d.Editable,
		"enabled":            d.Enabled,
		"segment":            d.Segment,
		"routing":            d.Routing,
		"endpoint_processor": &d.Processors,
	})
}
