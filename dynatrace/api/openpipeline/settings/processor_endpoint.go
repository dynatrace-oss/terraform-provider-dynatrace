package openpipeline

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/openpipeline/jsonmodel"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

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

func (ep *EndpointProcessors) FromJSON(endpointProcessors *[]jsonmodel.EndpointProcessor) error {
	ep.Processors = nil
	for _, ee := range *endpointProcessors {

		t, err := jsonmodel.ExtractEndpointProcessorType(ee)
		if err != nil {
			return err
		}
		switch t {
		case "DQLProcessor":
			dp, err := ee.AsDqlProcessor()
			if err != nil {
				return err
			}

			dqlProcessor := DqlProcessor{}
			if err := dqlProcessor.FromJSON(dp); err != nil {
				return err
			}
			ep.Processors = append(ep.Processors, EndpointProcessor{dqlProcessor: &dqlProcessor})

		case "FieldsAddProcessor":
			fap, err := ee.AsFieldsAddProcessor()
			if err != nil {
				return err
			}

			fieldsAddProcessor := FieldsAddProcessor{}
			if err := fieldsAddProcessor.FromJSON(fap); err != nil {
				return err
			}
			ep.Processors = append(ep.Processors, EndpointProcessor{fieldsAddProcessor: &fieldsAddProcessor})
		}

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
