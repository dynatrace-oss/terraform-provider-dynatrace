package openpipeline

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type EndpointDefinition struct {
	BasePath string `json:"basePath"`
	//Builtin       *bool   `json:"builtin,omitempty"`
	//DefaultBucket *string `json:"defaultBucket,omitempty"`
	//DisplayName   *string `json:"displayName,omitempty"`
	//Editable      *bool   `json:"editable,omitempty"`
	//Enabled       bool    `json:"enabled"`
	Processors []EndpointProcessor `json:"processors,omitempty"`
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
		"endpoint_processor": {
			Type:        schema.TypeList,
			Description: "todo",
			Optional:    true,
		},
	}
}

func (d *EndpointDefinition) MarshalHCL(properties hcl.Properties) error {

	if err := properties.EncodeAll(map[string]any{
		"base_path":          d.BasePath,
		"endpoint_processor": d.Processors,
	}); err != nil {
		return err
	}
	return nil
}

func (d *EndpointDefinition) UnmarshalHCL(decoder hcl.Decoder) error {

	return decoder.DecodeAll(map[string]any{
		"base_path":          &d.BasePath,
		"endpoint_processor": &d.Processors,
	})
}

type EndpointProcessor struct {
	Type string `json:"type"`
	//Description string  `json:"description"`
	//DqlScript   *string `json:"dqlScript"`
	//Editable    *bool   `json:"editable,omitempty"`
	//Enabled     bool    `json:"enabled"`
	//Id          string  `json:"id"`
	//Matcher     Matcher `json:"matcher"`
	//SampleData  *string `json:"sampleData,omitempty"`
	Fields *[]any `json:"fields"`
}

func (e *EndpointProcessor) Schema() map[string]*schema.Schema {

	var fieldsSchema *schema.Schema

	switch e.Type {
	case "fieldsRename":
		fieldsSchema = &schema.Schema{
			Type:        schema.TypeList,
			Description: "TODO",
			Elem:        &schema.Resource{Schema: new(FieldsRenameItem).Schema()},
			Required:    true,
		}
		//TODO: add other cases
	}

	return map[string]*schema.Schema{
		"type": {
			Type:        schema.TypeString,
			Description: "The custom base path of an openpipeline configuration",
			Required:    true,
		},
		"fields": fieldsSchema,
	}
}

func (e *EndpointProcessor) MarshalHCL(properties hcl.Properties) error {
	if err := properties.EncodeAll(map[string]any{
		"type":  e.Type,
		"field": e.Fields,
	}); err != nil {
		return err
	}
	return nil
}

func (e *EndpointProcessor) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"type":  &e.Type,
		"field": &e.Fields,
	})

}

type FieldsRenameItem struct {
	FromName *string `json:"fromName"`
	ToName   *string `json:"toName"`
}

func (f *FieldsRenameItem) MarshalHCL(properties hcl.Properties) error {
	if err := properties.EncodeAll(map[string]any{
		"from_name": f.FromName,
		"to_name":   f.ToName,
	}); err != nil {
		return err
	}
	return nil
}

func (f *FieldsRenameItem) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"from_name": &f.FromName,
		"to_name":   &f.ToName,
	})

}

func (f *FieldsRenameItem) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"from_name": {
			Type:        schema.TypeString,
			Description: "TODO",
			Required:    true,
		},
		"to_name": {
			Type:        schema.TypeString,
			Description: "TODO",
			Required:    true,
		},
	}
}

type Routing struct {
	PipelineId *string `json:"pipelineId"`
	Type       string  `json:"type"`
}
