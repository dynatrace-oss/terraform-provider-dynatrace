package openpipeline

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Configuration struct {
	CustomBasePath string `json:"customBasePath"`
	//Editable       *bool                `json:"editable,omitempty"`
	Endpoints []EndpointDefinition `json:"endpoints"`
	//Id             string               `json:"id"`
	//Pipelines      []Pipeline           `json:"pipelines"`
	//Routing        RoutingTable         `json:"routing"`
	//Version        string               `json:"version"`
}

func (d *Configuration) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"custom_base_path": {
			Type:        schema.TypeString,
			Description: "The custom base path of an openpipeline configuration",
			Required:    true,
		},

		"endpoint": {
			Type:        schema.TypeList,
			Description: "The endpoints of the openpipeline configuration",
			Elem:        &schema.Resource{Schema: new(EndpointDefinition).Schema()},
			Required:    true,
		},
	}
}

func (d *Configuration) MarshalHCL(properties hcl.Properties) error {

	if err := properties.EncodeAll(map[string]any{
		"custom_base_path": d.CustomBasePath,
		"endpoint":         d.Endpoints,
	}); err != nil {
		return err
	}
	return nil
}

func (d *Configuration) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"custom_base_path": &d.CustomBasePath,
		"endpoint":         &d.Endpoints,
	})

}
