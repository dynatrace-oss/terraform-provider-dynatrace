package openpipeline

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Processor struct {
	// Description Name or description of the processor.
	Description string `json:"description"`

	// Editable Indicates if the user is allowed to edit this object based on permissions and builtin property.
	Editable *bool `json:"editable,omitempty"`

	// Enabled Indicates if the object is active.
	Enabled bool `json:"enabled"`

	// Id Identifier of the processor. Must be unique within a stage.
	Id string `json:"id"`

	// Matcher Matching condition to apply on incoming records.
	Matcher string `json:"matcher"`

	// SampleData Sample data related to the processor for documentation or testing.
	SampleData *string `json:"sampleData,omitempty"`
}

func (p *Processor) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"description": {
			Type:        schema.TypeString,
			Description: "Name or description of the processor",
			Required:    true,
		},
		"editable": {
			Type:        schema.TypeBool,
			Description: "",
			Optional:    true,
		},
		"enabled": {
			Type:        schema.TypeBool,
			Description: "",
			Required:    true,
		},
		"id": {
			Type:        schema.TypeString,
			Description: "",
			Required:    true,
		},
		"matcher": {
			Type:        schema.TypeString,
			Description: "",
			Required:    true,
		},
		"sample_data": {
			Type:        schema.TypeString,
			Description: "",
			Optional:    true,
		},
	}
}

func (p *Processor) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"description": p.Description,
		"editable":    p.Editable,
		"enabled":     p.Enabled,
		"id":          p.Id,
		"matcher":     p.Matcher,
		"sample_data": p.SampleData,
	})
}

func (p *Processor) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"description": &p.Description,
		"editable":    &p.Editable,
		"enabled":     &p.Enabled,
		"id":          &p.Id,
		"matcher":     &p.Matcher,
		"sample_data": &p.SampleData,
	})
}

// DqlProcessor Processor to apply a DQL script.
type DqlProcessor struct {
	Processor

	// DqlScript The DQL script to apply on the record.
	DqlScript string `json:"dqlScript"`
}

func (p *DqlProcessor) Schema() map[string]*schema.Schema {
	s := p.Processor.Schema()
	s["dql_script"] = &schema.Schema{
		Type:        schema.TypeString,
		Description: "The DQL script to apply on the record.",
		Required:    true,
	}

	return s
}

func (p *DqlProcessor) MarshalHCL(properties hcl.Properties) error {
	if err := p.Processor.MarshalHCL(properties); err != nil {
		return err
	}
	return properties.Encode("dql_script", p.DqlScript)
}

func (p *DqlProcessor) UnmarshalHCL(decoder hcl.Decoder) error {
	if err := p.Processor.UnmarshalHCL(decoder); err != nil {
		return err
	}
	return decoder.Decode("dql_script", &p.DqlScript)
}

type FieldsAddProcessor struct {
	Processor

	// Fields List of fields to add to the record.
	Fields []FieldsAddItem `json:"fields"`
}

func (p *FieldsAddProcessor) Schema() map[string]*schema.Schema {
	s := p.Processor.Schema()
	s["field"] = &schema.Schema{
		Type:        schema.TypeList,
		Elem:        &schema.Resource{Schema: new(FieldsAddItem).Schema()},
		Description: "Field to add to the record.",
		Required:    true,
	}

	return s
}

func (p *FieldsAddProcessor) MarshalHCL(properties hcl.Properties) error {
	if err := p.Processor.MarshalHCL(properties); err != nil {
		return err
	}
	return properties.Encode("field", p.Fields)
}

func (p *FieldsAddProcessor) UnmarshalHCL(decoder hcl.Decoder) error {
	if err := p.Processor.UnmarshalHCL(decoder); err != nil {
		return err
	}
	return decoder.Decode("field", &p.Fields)
}

// FieldsAddItem List of fields to add to the record.
type FieldsAddItem struct {
	// Name Name of the field.
	Name string `json:"name"`

	// Value Value to assign to the field.
	Value string `json:"value"`
}

func (f *FieldsAddItem) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:        schema.TypeString,
			Description: "Name of the field",
			Required:    true,
		},
		"value": {
			Type:        schema.TypeString,
			Description: " Value to assign to the field",
			Required:    true,
		},
	}
}

func (f *FieldsAddItem) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"name":  f.Name,
		"value": f.Value,
	})
}

func (f *FieldsAddItem) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"name":  &f.Name,
		"value": &f.Value,
	})
}

type FieldsRemoveProcessor struct {
	Processor

	// Fields List of fields to remove from the record.
	Fields []string `json:"fields"`
}

func (p *FieldsRemoveProcessor) Schema() map[string]*schema.Schema {
	s := p.Processor.Schema()
	s["field"] = &schema.Schema{
		Type:        schema.TypeList,
		Elem:        &schema.Schema{Type: schema.TypeString},
		Description: "Field to add to the record.",
		Required:    true,
	}

	return s
}

func (p *FieldsRemoveProcessor) MarshalHCL(properties hcl.Properties) error {
	if err := p.Processor.MarshalHCL(properties); err != nil {
		return err
	}
	return properties.Encode("field", p.Fields)
}

func (p *FieldsRemoveProcessor) UnmarshalHCL(decoder hcl.Decoder) error {
	if err := p.Processor.UnmarshalHCL(decoder); err != nil {
		return err
	}
	return decoder.Decode("field", &p.Fields)
}

// FieldsRenameProcessor Processor to rename fields.
type FieldsRenameProcessor struct {
	Processor

	// Fields List of fields to rename on the record.
	Fields []FieldsRenameItem `json:"fields"`
}

func (p *FieldsRenameProcessor) Schema() map[string]*schema.Schema {
	s := p.Processor.Schema()
	s["field"] = &schema.Schema{
		Type:        schema.TypeList,
		Elem:        &schema.Resource{Schema: new(FieldsRenameItem).Schema()},
		Description: "Field to rename on the record.",
		Required:    true,
	}

	return s
}

func (p *FieldsRenameProcessor) MarshalHCL(properties hcl.Properties) error {
	if err := p.Processor.MarshalHCL(properties); err != nil {
		return err
	}
	return properties.Encode("field", p.Fields)
}

func (p *FieldsRenameProcessor) UnmarshalHCL(decoder hcl.Decoder) error {
	if err := p.Processor.UnmarshalHCL(decoder); err != nil {
		return err
	}
	return decoder.Decode("field", &p.Fields)
}

// FieldsRenameItem List of fields to rename on the record.
type FieldsRenameItem struct {
	// FromName The field to rename.
	FromName string `json:"fromName"`

	// ToName The new field name.
	ToName string `json:"toName"`
}

func (f *FieldsRenameItem) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"from_name": {
			Type:        schema.TypeString,
			Description: "The field to rename",
			Required:    true,
		},
		"to_name": {
			Type:        schema.TypeString,
			Description: "The new field name",
			Required:    true,
		},
	}
}

func (f *FieldsRenameItem) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"from_name": f.FromName,
		"to_name":   f.ToName,
	})
}

func (f *FieldsRenameItem) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"from_name": &f.FromName,
		"to_name":   &f.ToName,
	})
}

type EndpointProcessor struct {
	dqlProcessor          *DqlProcessor
	fieldsAddProcessor    *FieldsAddProcessor
	fieldsRemoveProcessor *FieldsRemoveProcessor
	fieldsRenameProcessor *FieldsRenameProcessor
}

func (ep *EndpointProcessor) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"dql_processor": {
			Type:        schema.TypeList,
			Description: "Processor to apply a DQL script",
			MinItems:    1,
			MaxItems:    1,
			Elem:        &schema.Resource{Schema: new(DqlProcessor).Schema()},
			Optional:    true,
		},
		"fields_add_processor": {
			Type:        schema.TypeList,
			Description: "",
			MinItems:    1,
			MaxItems:    1,
			Elem:        &schema.Resource{Schema: new(FieldsAddProcessor).Schema()},
			Optional:    true,
		},
		"fields_remove_processor": {
			Type:        schema.TypeList,
			Description: "",
			MinItems:    1,
			MaxItems:    1,
			Elem:        &schema.Resource{Schema: new(FieldsRemoveProcessor).Schema()},
			Optional:    true,
		},
		"fields_rename_processor": {
			Type:        schema.TypeList,
			Description: "",
			MinItems:    1,
			MaxItems:    1,
			Elem:        &schema.Resource{Schema: new(FieldsRenameProcessor).Schema()},
			Optional:    true,
		},
	}
}

func (ep *EndpointProcessor) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"dql_processor":           ep.dqlProcessor,
		"fields_add_processor":    ep.fieldsAddProcessor,
		"fields_remove_processor": ep.fieldsRemoveProcessor,
		"fields_rename_processor": ep.fieldsRenameProcessor,
	})
}

func (ep *EndpointProcessor) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"dql_processor":           &ep.dqlProcessor,
		"fields_add_processor":    &ep.fieldsAddProcessor,
		"fields_remove_processor": &ep.fieldsRemoveProcessor,
		"fields_rename_processor": &ep.fieldsRenameProcessor,
	})
}
