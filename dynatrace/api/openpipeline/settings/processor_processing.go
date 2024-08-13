package openpipeline

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type ProcessingStageProcessors struct {
	Processors []ProcessingStageProcessor
}

func (ep *ProcessingStageProcessors) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"processor": {
			Type:        schema.TypeSet,
			Description: "todo",
			Elem:        &schema.Resource{Schema: new(ProcessingStageProcessor).Schema()},
			Optional:    true,
		},
	}
}

func (ep *ProcessingStageProcessors) MarshalHCL(properties hcl.Properties) error {
	return properties.Encode("processors", ep.Processors)
}

func (ep *ProcessingStageProcessors) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.Decode("processors", &ep.Processors)
}

type ProcessingStageProcessor struct {
	dqlProcessor          *DqlProcessor
	fieldsAddProcessor    *FieldsAddProcessor
	fieldsRemoveProcessor *FieldsRemoveProcessor
	fieldsRenameProcessor *FieldsRenameProcessor
	technologyProcessor   *TechnologyProcessor
}

func (ep *ProcessingStageProcessor) Schema() map[string]*schema.Schema {
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
			Description: "Processor to add fields.",
			MinItems:    1,
			MaxItems:    1,
			Elem:        &schema.Resource{Schema: new(FieldsAddProcessor).Schema()},
			Optional:    true,
		},
		"fields_remove_processor": {
			Type:        schema.TypeList,
			Description: "Processor to remove fields.",
			MinItems:    1,
			MaxItems:    1,
			Elem:        &schema.Resource{Schema: new(FieldsRemoveProcessor).Schema()},
			Optional:    true,
		},
		"fields_rename_processor": {
			Type:        schema.TypeList,
			Description: "Processor to rename fields.",
			MinItems:    1,
			MaxItems:    1,
			Elem:        &schema.Resource{Schema: new(FieldsRenameProcessor).Schema()},
			Optional:    true,
		},
		"technology_processor": {
			Type:        schema.TypeList,
			Description: "Processor to apply a technology processors.",
			MinItems:    1,
			MaxItems:    1,
			Elem:        &schema.Resource{Schema: new(TechnologyProcessor).Schema()},
			Optional:    true,
		},
	}
}

func (ep *ProcessingStageProcessor) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"dql_processor":           ep.dqlProcessor,
		"fields_add_processor":    ep.fieldsAddProcessor,
		"fields_remove_processor": ep.fieldsRemoveProcessor,
		"fields_rename_processor": ep.fieldsRenameProcessor,
		"technology_processor":    ep.technologyProcessor,
	})
}

func (ep *ProcessingStageProcessor) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"dql_processor":           &ep.dqlProcessor,
		"fields_add_processor":    &ep.fieldsAddProcessor,
		"fields_remove_processor": &ep.fieldsRemoveProcessor,
		"fields_rename_processor": &ep.fieldsRenameProcessor,
		"technology_processor":    ep.technologyProcessor,
	})
}

type ClassicProcessingStageProcessors struct {
	Processors []ClassicProcessingStageProcessor
}

func (ep *ClassicProcessingStageProcessors) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"processor": {
			Type:        schema.TypeSet,
			Description: "todo",
			Elem:        &schema.Resource{Schema: new(ClassicProcessingStageProcessor).Schema()},
			Optional:    true,
		},
	}
}

func (ep *ClassicProcessingStageProcessors) MarshalHCL(properties hcl.Properties) error {
	return properties.Encode("processors", ep.Processors)
}

func (ep *ClassicProcessingStageProcessors) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.Decode("processors", &ep.Processors)
}

type ClassicProcessingStageProcessor struct {
	sqlxProcessor *SqlxProcessor
}

func (ep *ClassicProcessingStageProcessor) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"sqlx_processor": {
			Type:        schema.TypeList,
			Description: "",
			MinItems:    1,
			MaxItems:    1,
			Elem:        &schema.Resource{Schema: new(SqlxProcessor).Schema()},
			Optional:    true,
		},
	}
}

func (ep *ClassicProcessingStageProcessor) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"sqlx_processor": ep.sqlxProcessor,
	})
}

func (ep *ClassicProcessingStageProcessor) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"sqlx_processor": ep.sqlxProcessor,
	})
}
