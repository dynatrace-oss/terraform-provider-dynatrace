package openpipeline

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type EndpointProcessors struct {
	Processors []*EndpointProcessor
}

func (ep *EndpointProcessors) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"processor": {
			Type:        schema.TypeList,
			Description: "Groups all processors applicable for processing in the EndpointDefinition.\nApplicable processors are DqlProcessor, FieldsAddProcessor, FieldsRemoveProcessor, FieldsRenameProcessor and DropProcessor.",
			Elem:        &schema.Resource{Schema: new(EndpointProcessor).Schema()},
			Optional:    true,
		},
	}
}

func (ep *EndpointProcessors) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeSlice("processor", ep.Processors)
}

func (ep *EndpointProcessors) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeSlice("processor", &ep.Processors)
}

func (ep EndpointProcessors) MarshalJSON() ([]byte, error) {
	rawProcessors := []json.RawMessage{}
	for _, processor := range ep.Processors {
		rawProcessor, err := processor.MarshalJSON()
		if err != nil {
			return nil, err
		}

		rawProcessors = append(rawProcessors, rawProcessor)
	}

	return json.Marshal(rawProcessors)
}

func (ep *EndpointProcessors) UnmarshalJSON(b []byte) error {
	rawProcessors := []json.RawMessage{}
	if err := json.Unmarshal(b, &rawProcessors); err != nil {
		return err
	}

	ep.Processors = nil
	for _, rawProcessor := range rawProcessors {
		processor := EndpointProcessor{}
		if err := json.Unmarshal(rawProcessor, &processor); err != nil {
			return err
		}

		ep.Processors = append(ep.Processors, &processor)
	}
	return nil
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

func (ep EndpointProcessor) MarshalJSON() ([]byte, error) {
	if ep.dqlProcessor != nil {
		return json.Marshal(ep.dqlProcessor)
	}
	if ep.fieldsAddProcessor != nil {
		return json.Marshal(ep.fieldsAddProcessor)
	}
	if ep.fieldsRemoveProcessor != nil {
		return json.Marshal(ep.fieldsRemoveProcessor)
	}
	if ep.fieldsRenameProcessor != nil {
		return json.Marshal(ep.fieldsRenameProcessor)
	}

	return nil, errors.New("missing EndpointProcessor value")
}

func (ep *EndpointProcessor) UnmarshalJSON(b []byte) error {
	ttype, err := ExtractType(b)
	if err != nil {
		return err
	}

	switch ttype {
	case DqlProcessorType:
		dqlProcessor := DqlProcessor{}
		if err := json.Unmarshal(b, &dqlProcessor); err != nil {
			return err
		}
		ep.dqlProcessor = &dqlProcessor

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

	case FieldsRenameProcessorType:
		fieldsRenameProcessor := FieldsRenameProcessor{}
		if err := json.Unmarshal(b, &fieldsRenameProcessor); err != nil {
			return err
		}
		ep.fieldsRenameProcessor = &fieldsRenameProcessor

	default:
		return fmt.Errorf("unknown EndpointProcessor type: %s", ttype)
	}

	return nil
}
