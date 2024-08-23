package openpipeline

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type SecurityContextProcessors struct {
	Processors []SecContextProcessor
}

func (ep *SecurityContextProcessors) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"processor": {
			Type:        schema.TypeSet,
			Description: "todo",
			Elem:        &schema.Resource{Schema: new(SecContextProcessor).Schema()},
			Optional:    true,
		},
	}
}

func (ep *SecurityContextProcessors) MarshalHCL(properties hcl.Properties) error {
	return properties.Encode("processors", ep.Processors)
}

func (ep *SecurityContextProcessors) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.Decode("processors", &ep.Processors)
}

type SecContextProcessor struct {
	securityContextProcessor *SecurityContextProcessor
}

func (ep *SecContextProcessor) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"security_context_processor": {
			Type:        schema.TypeList,
			Description: "Processor to set the security context field.",
			MinItems:    1,
			MaxItems:    1,
			Elem:        &schema.Resource{Schema: new(SecurityContextProcessor).Schema()},
			Optional:    true,
		},
	}
}

func (ep *SecContextProcessor) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"security_context_processor": ep.securityContextProcessor,
	})
}

func (ep *SecContextProcessor) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"security_context_processor": ep.securityContextProcessor,
	})
}
