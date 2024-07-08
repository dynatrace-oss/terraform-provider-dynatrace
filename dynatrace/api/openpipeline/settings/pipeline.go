package openpipeline

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Pipeline struct {
	Builtin          *bool                   `json:"builtin,omitempty"`
	DataExtraction   *DataExtractionStage    `json:"dataExtraction,omitempty"`
	DisplayName      *string                 `json:"displayName,omitempty"`
	Editable         *bool                   `json:"editable,omitempty"`
	Enabled          bool                    `json:"enabled"`
	Id               string                  `json:"id"`
	MetricExtraction *MetricExtractionStage  `json:"metricExtraction,omitempty"`
	Processing       *ClassicProcessingStage `json:"processing,omitempty"`
	SecurityContext  *SecurityContextStage   `json:"securityContext,omitempty"`
	SettingsSchema   *string                 `json:"settingsSchema,omitempty"` // classic
	Storage          *StorageStage           `json:"storage,omitempty"`
	Type             string                  `json:"type"`
}

func (d *Pipeline) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"builtin": {
			Type:        schema.TypeBool,
			Description: "Defines whether this is a builtin pipeline",
			Optional:    true,
		},

		"display_name": {
			Type:        schema.TypeString,
			Description: "Defines whether this openpipeline configuration is editable",
			Optional:    true,
		},
		"editable": {
			Type:        schema.TypeBool,
			Description: "",
			Optional:    true,
		},
		"enabled": {
			Type:        schema.TypeBool,
			Description: "",
			Optional:    true,
		},
		"id": {
			Type:        schema.TypeString,
			Description: "",
			Optional:    true,
		},
		"data_extraction": {
			Type:        schema.TypeList,
			Description: "",
			Optional:    true,
			MinItems:    1,
			MaxItems:    1,
			Elem:        &schema.Resource{Schema: new(DataExtractionStage).Schema()},
		},
		"metric_extraction": {
			Type:        schema.TypeList,
			Description: "",
			Optional:    true,
			MinItems:    1,
			MaxItems:    1,
			Elem:        &schema.Resource{Schema: new(MetricExtractionStage).Schema()},
		},
		"processing": {
			Type:        schema.TypeList,
			Description: "",
			Optional:    true,
			MinItems:    1,
			MaxItems:    1,
			Elem:        &schema.Resource{Schema: new(ClassicProcessingStage).Schema()},
		},
		"security_context": {
			Type:        schema.TypeList,
			Description: "",
			Optional:    true,
			MinItems:    1,
			MaxItems:    1,
			Elem:        &schema.Resource{Schema: new(SecurityContextStage).Schema()},
		},

		"settings_schema": {
			Type:        schema.TypeString,
			Description: "",
			Optional:    true,
		},
		"storage": {
			Type:        schema.TypeList,
			Description: "",
			Optional:    true,
			MinItems:    1,
			MaxItems:    1,
			Elem:        &schema.Resource{Schema: new(StorageStage).Schema()},
		},
		"type": {
			Type:        schema.TypeString,
			Description: "",
			Optional:    true,
		},
	}
}

func (d *Pipeline) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{})
}

func (d *Pipeline) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{})
}

type DataExtractionStage struct {
	Editable   *bool                     `json:"editable,omitempty"`
	Processors []DataExtractionProcessor `json:"processors"`
}

func (d *DataExtractionStage) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"editable": {
			Type:        schema.TypeBool,
			Description: "",
			Optional:    true,
		},
	}
}

func (d *DataExtractionStage) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{})
}

func (d *DataExtractionStage) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{})
}

type MetricExtractionStage struct {
	Editable   *bool                       `json:"editable,omitempty"`
	Processors []MetricExtractionProcessor `json:"processors"`
}

func (m *MetricExtractionStage) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"editable": {
			Type:        schema.TypeBool,
			Description: "",
			Optional:    true,
		},
	}
}

func (m *MetricExtractionStage) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{})
}

func (m *MetricExtractionStage) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{})
}

type ClassicProcessingStage struct {
	Editable   *bool           `json:"editable,omitempty"`
	Processors []SxqlProcessor `json:"processors"`
}

func (c *ClassicProcessingStage) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"editable": {
			Type:        schema.TypeBool,
			Description: "",
			Optional:    true,
		},
	}
}

func (c *ClassicProcessingStage) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{})
}

func (c *ClassicProcessingStage) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{})
}

type SecurityContextStage struct {
	Editable   *bool                      `json:"editable,omitempty"`
	Processors []SecurityContextProcessor `json:"processors"`
}

func (s *SecurityContextStage) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"editable": {
			Type:        schema.TypeBool,
			Description: "",
			Optional:    true,
		},
	}
}

func (s *SecurityContextStage) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{})
}

func (s *SecurityContextStage) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{})
}

type StorageStage struct {
	CatchAllBucketName string                  `json:"catchAllBucketName"`
	Editable           *bool                   `json:"editable,omitempty"`
	Processors         []StorageStageProcessor `json:"processors"`
}

func (s *StorageStage) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"catch_all_bucket_name": {
			Type:        schema.TypeString,
			Description: "",
			Optional:    true,
		},
		"editable": {
			Type:        schema.TypeBool,
			Description: "",
			Optional:    true,
		},
	}
}

func (s *StorageStage) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{})
}

func (s *StorageStage) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{})
}
