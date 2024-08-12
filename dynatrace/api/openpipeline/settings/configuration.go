package openpipeline

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Configuration struct {
	CustomBasePath string `json:"customBasePath"`
	//Editable       *bool                `json:"editable,omitempty"`
	Endpoints Endpoints `json:"endpoints"`
	//Id             string               `json:"id"`
	//Pipelines      []Pipeline           `json:"pipelines"`
	Routing RoutingTable `json:"routing"`
	//Version        string               `json:"version"`
}

func (d *Configuration) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"custom_base_path": {
			Type:        schema.TypeString,
			Description: "The custom base path of an openpipeline configuration",
			Required:    true,
		},

		"endpoints": {
			Type:        schema.TypeList,
			Description: "The endpoints of the openpipeline configuration",
			MinItems:    1,
			MaxItems:    1,
			Elem:        &schema.Resource{Schema: new(Endpoints).Schema()},
			Required:    true,
		},
		"routing": {
			Type:        schema.TypeList,
			Description: "todo",
			MinItems:    1,
			MaxItems:    1,
			Elem:        &schema.Resource{Schema: new(RoutingTable).Schema()},
			Required:    true,
		},
	}
}

func (d *Configuration) MarshalHCL(properties hcl.Properties) error {

	if err := properties.EncodeAll(map[string]any{
		"custom_base_path": d.CustomBasePath,
		"endpoints":        d.Endpoints,
	}); err != nil {
		return err
	}
	return nil
}

func (d *Configuration) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"custom_base_path": &d.CustomBasePath,
		"endpoints":        &d.Endpoints,
	})

}
