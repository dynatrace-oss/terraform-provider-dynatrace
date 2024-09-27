package openpipeline

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	ClassicPipelineType = "classic"
	DefaultPipelineType = "default"
)

type Pipelines struct {
	Pipelines []*Pipeline
}

func (ep *Pipelines) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"pipeline": {
			Type:        schema.TypeList,
			Description: "todo",
			Elem:        &schema.Resource{Schema: new(Pipeline).Schema()},
			Optional:    true,
		},
	}
}

func (ep *Pipelines) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeSlice("pipeline", ep.Pipelines)
}

func (ep *Pipelines) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeSlice("pipeline", &ep.Pipelines)
}

func (d Pipelines) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.Pipelines)
}

func (d *Pipelines) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &d.Pipelines)
}

func (d *Pipelines) RemoveFixed() {
	filteredPipelines := []*Pipeline{}
	for _, pipeline := range d.Pipelines {
		if !pipeline.IsFixed() {
			filteredPipelines = append(filteredPipelines, pipeline)
		}
	}
	d.Pipelines = filteredPipelines
}

type Pipeline struct {
	classicPipeline *ClassicPipeline
	defaultPipeline *DefaultPipeline
}

func (ep *Pipeline) Schema() map[string]*schema.Schema {
	return new(DefaultPipeline).Schema()
}

func (ep *Pipeline) MarshalHCL(properties hcl.Properties) error {
	if ep.classicPipeline != nil {
		return errors.New("cannot marshal classic pipeline to HCL")
	}
	if ep.defaultPipeline == nil {
		return errors.New("missing default pipeline")
	}

	return ep.defaultPipeline.MarshalHCL(properties)
}

func (ep *Pipeline) UnmarshalHCL(decoder hcl.Decoder) error {
	ep.defaultPipeline = &DefaultPipeline{}
	return ep.defaultPipeline.UnmarshalHCL(decoder)
}

func (ep Pipeline) MarshalJSON() ([]byte, error) {
	if ep.classicPipeline != nil {
		return json.Marshal(ep.classicPipeline)
	}
	if ep.defaultPipeline != nil {
		return json.Marshal(ep.defaultPipeline)
	}

	return nil, errors.New("missing Pipeline value")
}

func (ep *Pipeline) UnmarshalJSON(b []byte) error {
	ttype, err := ExtractType(b)
	if err != nil {
		return err
	}

	switch ttype {
	case ClassicPipelineType:
		classicPipeline := ClassicPipeline{}
		if err := json.Unmarshal(b, &classicPipeline); err != nil {
			return err
		}
		ep.classicPipeline = &classicPipeline

	case DefaultPipelineType:
		defaultPipeline := DefaultPipeline{}
		if err := json.Unmarshal(b, &defaultPipeline); err != nil {
			return err
		}
		ep.defaultPipeline = &defaultPipeline

	default:
		return fmt.Errorf("unknown pipeline type: %s", ttype)
	}

	return nil
}

func (ep *Pipeline) IsFixed() bool {
	if ep.classicPipeline != nil {
		return ep.classicPipeline.IsFixed()
	}
	if ep.defaultPipeline != nil {
		return ep.defaultPipeline.IsFixed()
	}
	return false
}

type BasePipeline struct {
	Builtin          *bool                  `json:"builtin,omitempty"`
	DataExtraction   *DataExtractionStage   `json:"dataExtraction"`
	DisplayName      *string                `json:"displayName,omitempty"`
	Editable         *bool                  `json:"editable,omitempty"`
	Enabled          bool                   `json:"enabled"`
	Id               string                 `json:"id"`
	MetricExtraction *MetricExtractionStage `json:"metricExtraction,omitempty"`
	SecurityContext  *SecurityContextStage  `json:"securityContext"`
	Storage          *StorageStage          `json:"storage,omitempty"`
	Type             string                 `json:"type"`
}

func (ep *BasePipeline) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
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
			Required:    true,
		},
		"enabled": {
			Type:        schema.TypeBool,
			Description: "Indicates if the object is active.",
			Required:    true,
		},
		"id": {
			Type:        schema.TypeString,
			Description: "Identifier of the pipeline.",
			Required:    true,
		},
	}
}

func (ep *BasePipeline) MarshalHCL(properties hcl.Properties) error {
	if ep.DataExtraction != nil && len(ep.DataExtraction.Processors) > 0 {
		if err := properties.Encode("data_extraction", ep.DataExtraction); err != nil {
			return err
		}
	}

	if ep.MetricExtraction != nil && len(ep.MetricExtraction.Processors) > 0 {
		if err := properties.Encode("metric_extraction", ep.MetricExtraction); err != nil {
			return err
		}
	}

	if ep.SecurityContext != nil && len(ep.SecurityContext.Processors) > 0 {
		if err := properties.Encode("security_context", ep.SecurityContext); err != nil {
			return err
		}
	}

	if ep.Storage != nil && len(ep.Storage.Processors) > 0 {
		if err := properties.Encode("storage", ep.Storage); err != nil {
			return err
		}
	}

	return properties.EncodeAll(map[string]any{
		"display_name": ep.DisplayName,
		"enabled":      ep.Enabled,
		"id":           ep.Id,
	})
}

func (ep *BasePipeline) UnmarshalHCL(decoder hcl.Decoder) error {
	if err := decoder.DecodeAll(map[string]any{
		"data_extraction":   &ep.DataExtraction,
		"metric_extraction": &ep.MetricExtraction,
		"security_context":  &ep.SecurityContext,
		"storage":           &ep.Storage,
		"display_name":      &ep.DisplayName,
		"enabled":           &ep.Enabled,
		"id":                &ep.Id,
	}); err != nil {
		return err
	}

	if ep.DataExtraction == nil {
		ep.DataExtraction = &DataExtractionStage{}
	}

	if ep.MetricExtraction == nil {
		ep.MetricExtraction = &MetricExtractionStage{}
	}

	if ep.SecurityContext == nil {
		ep.SecurityContext = &SecurityContextStage{}
	}

	if ep.Storage == nil {
		ep.Storage = &StorageStage{}
	}

	return nil
}

func (ep *BasePipeline) IsFixed() bool {
	return (ep.Builtin != nil) && *ep.Builtin
}

type DefaultPipeline struct {
	BasePipeline
	Processing *ProcessingStage `json:"processing,omitempty"`
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

	if p.Processing != nil && len(p.Processing.Processors) > 0 {
		if err := properties.Encode("processing", p.Processing); err != nil {
			return err
		}
	}
	return nil
}

func (p *DefaultPipeline) UnmarshalHCL(decoder hcl.Decoder) error {
	if err := p.BasePipeline.UnmarshalHCL(decoder); err != nil {
		return err
	}

	if err := decoder.Decode("processing", &p.Processing); err != nil {
		return err
	}

	if p.Processing == nil {
		p.Processing = &ProcessingStage{}
	}

	return nil
}

func (p DefaultPipeline) MarshalJSON() ([]byte, error) {
	type defaultPipeline DefaultPipeline
	return MarshalAsJSONWithType((defaultPipeline)(p), DefaultPipelineType)
}

type ClassicPipeline struct {
	BasePipeline
	Processing     *ClassicProcessingStage `json:"processing,omitempty"`
	SettingsSchema string                  `json:"settingsSchema,omitempty"`
}

func (p ClassicPipeline) MarshalJSON() ([]byte, error) {
	type classicPipeline ClassicPipeline
	return MarshalAsJSONWithType((classicPipeline)(p), ClassicPipelineType)
}
