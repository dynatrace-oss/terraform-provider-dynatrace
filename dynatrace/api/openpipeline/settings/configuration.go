package openpipeline

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Configuration struct {
	CustomBasePath string       `json:"customBasePath"`
	Editable       *bool        `json:"editable,omitempty"`
	Endpoints      Endpoints    `json:"endpoints"`
	Kind           string       `json:"id"`
	Pipelines      Pipelines    `json:"pipelines"`
	Routing        RoutingTable `json:"routing"`
	Version        string       `json:"version"`
}

func (d *Configuration) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"custom_base_path": {
			Type:        schema.TypeString,
			Description: "The base path for custom ingest endpoints.",
			Required:    true,
		},
		"editable": {
			Type:        schema.TypeBool,
			Description: "Indicates if the user is allowed to edit this object based on permissions and builtin property.",
			Optional:    true,
		},

		"endpoints": {
			Type:        schema.TypeList,
			Description: "List of all ingest sources of the configuration.",
			MinItems:    1,
			MaxItems:    1,
			Elem:        &schema.Resource{Schema: new(Endpoints).Schema()},
			Required:    true,
		},
		"kind": {
			Type:        schema.TypeString,
			Description: "Identifier of the configuration.",
			Required:    true,
		},
		"routing": {
			Type:        schema.TypeList,
			Description: "Dynamic routing definition.",
			MinItems:    1,
			MaxItems:    1,
			Elem:        &schema.Resource{Schema: new(RoutingTable).Schema()},
			Required:    true,
		},
		"pipelines": {
			Type:        schema.TypeList,
			Description: "List of all pipelines of the configuration.",
			MinItems:    1,
			MaxItems:    1,
			Elem:        &schema.Resource{Schema: new(Pipelines).Schema()},
			Required:    true,
		},
		"version": {
			Type:        schema.TypeString,
			Description: "The current version of the configuration.",
			Optional:    true,
		},
	}
}

func (d *Configuration) MarshalHCL(properties hcl.Properties) error {

	if err := properties.EncodeAll(map[string]any{
		"custom_base_path": d.CustomBasePath,
		"endpoints":        d.Endpoints,
		"editable":         d.Editable,
		"kind":             d.Kind,
		"version":          d.Version,
	}); err != nil {
		return err
	}
	return nil
}

func (d *Configuration) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"custom_base_path": &d.CustomBasePath,
		"endpoints":        &d.Endpoints,
		"editable":         &d.Editable,
		"kind":             &d.Kind,
		"version":          &d.Version,
	})

}
