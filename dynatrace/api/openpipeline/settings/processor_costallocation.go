package openpipeline

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type CostAllocationProcessor struct {
	Processor
	Value *ValueAssignment `json:"value,omitempty"`
}

func (p *CostAllocationProcessor) Schema() map[string]*schema.Schema {
	s := p.Processor.Schema()
	s["value"] = &schema.Schema{
		Type:        schema.TypeList,
		MinItems:    1,
		MaxItems:    1,
		Elem:        &schema.Resource{Schema: new(ValueAssignment).Schema()},
		Description: "Strategy to assign a value",
		Required:    true,
	}

	return s
}

func (p *CostAllocationProcessor) MarshalHCL(properties hcl.Properties) error {
	if err := p.Processor.MarshalHCL(properties); err != nil {
		return err
	}

	return properties.Encode("value", p.Value)
}

func (p *CostAllocationProcessor) UnmarshalHCL(decoder hcl.Decoder) error {
	if err := p.Processor.UnmarshalHCL(decoder); err != nil {
		return err
	}

	return decoder.Decode("value", &p.Value)
}

func (ep CostAllocationProcessor) MarshalJSON() ([]byte, error) {
	type costAllocationProcessor CostAllocationProcessor
	return MarshalAsJSONWithType((costAllocationProcessor)(ep), CostAllocationProcessorType)
}
