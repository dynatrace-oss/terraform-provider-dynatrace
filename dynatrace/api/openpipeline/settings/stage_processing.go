package openpipeline

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type ProcessingStage struct {
	Editable   *bool                       `json:"editable,omitempty"`
	Processors []*ProcessingStageProcessor `json:"processors"`
}

func (f *ProcessingStage) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"processor": {
			Type:        schema.TypeList,
			Description: "Groups all processors applicable for the ProcessingStage.\nApplicable processors are DqlProcessor, FieldsAddProcessor, FieldsRemoveProcessor, FieldsRenameProcessor, TechnologyProcessor and DropProcessor.",
			Elem:        &schema.Resource{Schema: new(ProcessingStageProcessor).Schema()},
			Optional:    true,
		},
	}
}

func (f *ProcessingStage) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeSlice("processor", f.Processors)
}

func (f *ProcessingStage) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeSlice("processor", &f.Processors)
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
			Description: "Processor to add fields",
			MinItems:    1,
			MaxItems:    1,
			Elem:        &schema.Resource{Schema: new(FieldsAddProcessor).Schema()},
			Optional:    true,
		},
		"fields_remove_processor": {
			Type:        schema.TypeList,
			Description: "Processor to remove fields",
			MinItems:    1,
			MaxItems:    1,
			Elem:        &schema.Resource{Schema: new(FieldsRemoveProcessor).Schema()},
			Optional:    true,
		},
		"fields_rename_processor": {
			Type:        schema.TypeList,
			Description: "Processor to rename fields",
			MinItems:    1,
			MaxItems:    1,
			Elem:        &schema.Resource{Schema: new(FieldsRenameProcessor).Schema()},
			Optional:    true,
		},
		"technology_processor": {
			Type:        schema.TypeList,
			Description: "Processor to apply a technology processors",
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
		"technology_processor":    &ep.technologyProcessor,
	})
}

func (ep ProcessingStageProcessor) MarshalJSON() ([]byte, error) {
	if ep.fieldsRenameProcessor != nil {
		return json.Marshal(ep.fieldsRenameProcessor)
	}
	if ep.dqlProcessor != nil {
		return json.Marshal(ep.dqlProcessor)
	}
	if ep.fieldsAddProcessor != nil {
		return json.Marshal(ep.fieldsAddProcessor)
	}
	if ep.fieldsRemoveProcessor != nil {
		return json.Marshal(ep.fieldsRemoveProcessor)
	}
	if ep.technologyProcessor != nil {
		return json.Marshal(ep.technologyProcessor)
	}

	return nil, errors.New("missing ProcessingStageProcessor value")
}

func (ep *ProcessingStageProcessor) UnmarshalJSON(b []byte) error {
	ttype, err := ExtractType(b)
	if err != nil {
		return err
	}

	switch ttype {
	case FieldsRenameProcessorType:
		fieldsRenameProcessor := FieldsRenameProcessor{}
		if err := json.Unmarshal(b, &fieldsRenameProcessor); err != nil {
			return err
		}
		ep.fieldsRenameProcessor = &fieldsRenameProcessor

	case FieldsAddProcessorType:
		fieldsAddProcessor := FieldsAddProcessor{}
		if err := json.Unmarshal(b, &fieldsAddProcessor); err != nil {
			return err
		}
		ep.fieldsAddProcessor = &fieldsAddProcessor

	case FieldsRemoveProcessorType:
		fieldsRemoveProcessor := FieldsRemoveProcessor{}
		if err := json.Unmarshal(b, &fieldsRemoveProcessor); err != nil {
			return err
		}
		ep.fieldsRemoveProcessor = &fieldsRemoveProcessor
	case DqlProcessorType:
		dqlProcessor := DqlProcessor{}
		if err := json.Unmarshal(b, &dqlProcessor); err != nil {
			return err
		}
		ep.dqlProcessor = &dqlProcessor
	case TechnologyProcessorType:
		technologyProcessor := TechnologyProcessor{}
		if err := json.Unmarshal(b, &technologyProcessor); err != nil {
			return err
		}
		ep.technologyProcessor = &technologyProcessor

	default:
		return fmt.Errorf("unknown ProcessingStageProcessor type: %s", ttype)
	}

	return nil
}
