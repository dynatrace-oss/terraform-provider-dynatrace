package openpipeline

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type MetricExtractionStage struct {
	Editable   *bool                      `json:"editable,omitempty"`
	Processors MetricExtractionProcessors `json:"processors,omitempty"`
}

func (f *MetricExtractionStage) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"editable": {
			Type:        schema.TypeString,
			Description: "todo",
			Optional:    true,
		},
		"processors": {
			Type:        schema.TypeList,
			Description: "todo",
			MinItems:    1,
			MaxItems:    1,
			Elem:        &schema.Resource{Schema: new(MetricExtractionProcessors).Schema()},
			Optional:    true,
		},
	}
}

func (f *MetricExtractionStage) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"editable":   f.Editable,
		"processors": f.Processors,
	})
}

func (f *MetricExtractionStage) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"editable":   &f.Editable,
		"processors": &f.Processors,
	})
}

type DataExtractionStage struct {
	Editable   *bool                    `json:"editable,omitempty"`
	Processors DataExtractionProcessors `json:"processors"`
}

func (f *DataExtractionStage) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:        schema.TypeString,
			Description: "todo",
			Optional:    true,
		},
		"processors": {
			Type:        schema.TypeList,
			Description: "todo",
			MinItems:    1,
			MaxItems:    1,
			Elem:        &schema.Resource{Schema: new(DataExtractionProcessors).Schema()},
			Optional:    true,
		},
	}
}

func (f *DataExtractionStage) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"editable":   f.Editable,
		"processors": f.Processors,
	})
}

func (f *DataExtractionStage) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"editable":   &f.Editable,
		"processors": &f.Processors,
	})
}

type ProcessingStage struct {
	Editable   *bool                     `json:"editable,omitempty"`
	Processors ProcessingStageProcessors `json:"processors,omitempty"`
}

func (f *ProcessingStage) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"editable": {
			Type:        schema.TypeBool,
			Description: "todo",
			Optional:    true,
		},
		"processors": {
			Type:        schema.TypeList,
			Description: "todo",
			MinItems:    1,
			MaxItems:    1,
			Elem:        &schema.Resource{Schema: new(DataExtractionProcessors).Schema()},
			Optional:    true,
		},
	}
}

func (f *ProcessingStage) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"editable":   f.Editable,
		"processors": f.Processors,
	})
}

func (f *ProcessingStage) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"editable":   &f.Editable,
		"processors": &f.Processors,
	})
}

type StorageStage struct {
	Editable           *bool                  `json:"editable,omitempty"`
	CatchAllBucketName string                 `json:"catchAllBucketName"`
	Processors         StorageStageProcessors `json:"processors"`
}

func (f *StorageStage) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"editable": {
			Type:        schema.TypeString,
			Description: "todo",
			Optional:    true,
		},
		"catch_all_bucket_name": {
			Type:        schema.TypeString,
			Description: "todo",
			Optional:    true,
		},
		"processors": {
			Type:        schema.TypeList,
			Description: "todo",
			MinItems:    1,
			MaxItems:    1,
			Elem:        &schema.Resource{Schema: new(MetricExtractionProcessors).Schema()},
			Optional:    true,
		},
	}
}

func (f *StorageStage) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"editable":              f.Editable,
		"catch_all_bucket_name": f.CatchAllBucketName,
		"processors":            f.Processors,
	})
}

func (f *StorageStage) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"editable":              &f.Editable,
		"catch_all_bucket_name": f.CatchAllBucketName,
		"processors":            &f.Processors,
	})
}

type SecurityContextStage struct {
	Editable           *bool                     `json:"editable,omitempty"`
	CatchAllBucketName string                    `json:"catchAllBucketName"`
	Processors         SecurityContextProcessors `json:"processors"`
}

func (f *SecurityContextStage) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"editable": {
			Type:        schema.TypeString,
			Description: "todo",
			Optional:    true,
		},
		"processors": {
			Type:        schema.TypeList,
			Description: "todo",
			MinItems:    1,
			MaxItems:    1,
			Elem:        &schema.Resource{Schema: new(SecurityContextProcessors).Schema()},
			Optional:    true,
		},
	}
}

func (f *SecurityContextStage) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"editable":   f.Editable,
		"processors": f.Processors,
	})
}

func (f *SecurityContextStage) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"editable":   &f.Editable,
		"processors": &f.Processors,
	})
}

type ClassicProcessingStage struct {
	Editable   *bool                             `json:"editable,omitempty"`
	Processors []ClassicProcessingStageProcessor `json:"processors,omitempty"`
}

func (f *ClassicProcessingStage) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"editable": {
			Type:        schema.TypeBool,
			Description: "todo",
			Optional:    true,
		},
		"processors": {
			Type:        schema.TypeList,
			Description: "todo",
			MinItems:    1,
			MaxItems:    1,
			Elem:        &schema.Resource{Schema: new(ClassicProcessingStageProcessors).Schema()},
			Optional:    true,
		},
	}
}

func (f *ClassicProcessingStage) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"editable":   f.Editable,
		"processors": f.Processors,
	})
}

func (f *ClassicProcessingStage) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"editable":   &f.Editable,
		"processors": &f.Processors,
	})
}
