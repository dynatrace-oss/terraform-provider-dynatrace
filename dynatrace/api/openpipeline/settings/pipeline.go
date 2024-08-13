package openpipeline

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Pipelines struct {
	Pipelines []Pipeline `json:"pipelines"`
}

func (ep *Pipelines) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"pipeline": {
			Type:        schema.TypeSet,
			Description: "todo",
			Elem:        &schema.Resource{Schema: new(Pipeline).Schema()},
			Optional:    true,
		},
	}
}

func (ep *Pipelines) MarshalHCL(properties hcl.Properties) error {
	return properties.Encode("pipelines", ep.Pipelines)
}

func (ep *Pipelines) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.Decode("pipelines", &ep.Pipelines)
}

type Pipeline struct {
	classicPipeline *ClassicPipeline
	defaultPipeline *DefaultPipeline
}

func (ep *Pipeline) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"classic_pipeline": {
			Type:        schema.TypeList,
			Description: "Processor to apply a DQL script",
			MinItems:    1,
			MaxItems:    1,
			Elem:        &schema.Resource{Schema: new(ClassicPipeline).Schema()},
			Optional:    true,
		},
		"default_pipeline": {
			Type:        schema.TypeList,
			Description: "",
			MinItems:    1,
			MaxItems:    1,
			Elem:        &schema.Resource{Schema: new(DefaultPipeline).Schema()},
			Optional:    true,
		},
	}
}

func (ep *Pipeline) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"classic_pipeline": ep.classicPipeline,
		"default_pipeline": ep.defaultPipeline,
	})
}

func (ep *Pipeline) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"classic_pipeline": &ep.classicPipeline,
		"default_pipeline": &ep.defaultPipeline,
	})
}

type BasePipeline struct {
	Builtin *bool `json:"builtin,omitempty"`

	DataExtraction *DataExtractionStage `json:"dataExtraction,omitempty"`

	DisplayName *string `json:"displayName,omitempty"`

	Editable *bool `json:"editable,omitempty"`

	Enabled bool `json:"enabled"`

	Id string `json:"id"`

	MetricExtraction *MetricExtractionStage `json:"metricExtraction,omitempty"`

	SecurityContext *SecurityContextStage `json:"securityContext,omitempty"`

	Storage *StorageStage `json:"storage,omitempty"`

	Type string `json:"type"`
}

func (ep *BasePipeline) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"builtin": {
			Type:        schema.TypeBool,
			Description: "todo.",
			Required:    true,
		},
		"data_extraction": {
			Type:        schema.TypeList,
			Description: "todo",
			MinItems:    1,
			MaxItems:    1,
			Elem:        &schema.Resource{Schema: new(DataExtractionStage).Schema()},
			Optional:    true,
		},
		"metric_extraction": {
			Type:        schema.TypeList,
			Description: "todo",
			MinItems:    1,
			MaxItems:    1,
			Elem:        &schema.Resource{Schema: new(MetricExtractionStage).Schema()},
			Optional:    true,
		},
		"security_context": {
			Type:        schema.TypeList,
			Description: "todo",
			MinItems:    1,
			MaxItems:    1,
			Elem:        &schema.Resource{Schema: new(SecurityContextStage).Schema()},
			Optional:    true,
		},
		"storage": {
			Type:        schema.TypeList,
			Description: "Data extraction stage configuration of the pipeline.",
			MinItems:    1,
			MaxItems:    1,
			Elem:        &schema.Resource{Schema: new(StorageStage).Schema()},
			Optional:    true,
		},
		"display_name": {
			Type:        schema.TypeString,
			Description: "Display name of the pipeline.",
			Optional:    true,
		},
		"editable": {
			Type:        schema.TypeBool,
			Description: "Indicates if the user is allowed to edit this object based on permissions and builtin property.",
			Optional:    true,
		},
		"enabled": {
			Type:        schema.TypeBool,
			Description: "Indicates if the object is active.",
			Optional:    true,
		},
		"id": {
			Type:        schema.TypeString,
			Description: "Identifier of the pipeline.",
			Optional:    true,
		},
	}
}

func (ep *BasePipeline) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"builtin":         ep.Builtin,
		"data_extraction": ep.DataExtraction,
		"display_name":    ep.DisplayName,
		"editable":        ep.Editable,
		"enabled":         ep.Enabled,
		"id":              ep.Id,
	})
}

func (ep *BasePipeline) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"builtin":         ep.Builtin,
		"data_extraction": ep.DataExtraction,
		"display_name":    ep.DisplayName,
		"editable":        ep.Editable,
		"enabled":         ep.Enabled,
		"id":              ep.Id,
	})
}

type DefaultPipeline struct {
	BasePipeline
	Processing ProcessingStage `json:"processing,omitempty"`
}

func (p *DefaultPipeline) Schema() map[string]*schema.Schema {
	s := p.BasePipeline.Schema()
	s["processing"] = &schema.Schema{
		Type:        schema.TypeList,
		Description: "todo",
		MinItems:    1,
		MaxItems:    1,
		Elem:        &schema.Resource{Schema: new(ProcessingStage).Schema()},
		Optional:    true,
	}
	return s
}

func (p *DefaultPipeline) MarshalHCL(properties hcl.Properties) error {
	if err := p.BasePipeline.MarshalHCL(properties); err != nil {
		return err
	}
	return nil
}

func (p *DefaultPipeline) UnmarshalHCL(decoder hcl.Decoder) error {
	if err := p.BasePipeline.UnmarshalHCL(decoder); err != nil {
		return err
	}
	return nil
}

type ClassicPipeline struct {
	BasePipeline
	Processing     ClassicProcessingStage `json:"processing,omitempty"`
	SettingsSchema string                 `json:"settingsSchema,omitempty"`
}

func (p *ClassicPipeline) Schema() map[string]*schema.Schema {
	s := p.BasePipeline.Schema()
	s["settings_schema"] = &schema.Schema{
		Type:        schema.TypeString,
		Description: "The DQL script to apply on the record.",
		Required:    true,
	}
	s["processing"] = &schema.Schema{
		Type:        schema.TypeList,
		Description: "Processor to apply a DQL script",
		MinItems:    1,
		MaxItems:    1,
		Elem:        &schema.Resource{Schema: new(ClassicProcessingStage).Schema()},
		Optional:    true,
	}

	return s
}

func (p *ClassicPipeline) MarshalHCL(properties hcl.Properties) error {
	if err := p.BasePipeline.MarshalHCL(properties); err != nil {
		return err
	}
	return properties.Encode("settings_schema", p.SettingsSchema)
}

func (p *ClassicPipeline) UnmarshalHCL(decoder hcl.Decoder) error {
	if err := p.BasePipeline.UnmarshalHCL(decoder); err != nil {
		return err
	}
	return decoder.Decode("settings_schema", &p.SettingsSchema)
}
