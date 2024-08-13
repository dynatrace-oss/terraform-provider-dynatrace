package openpipeline

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Processor struct {
	Description string  `json:"description"`
	Editable    *bool   `json:"editable,omitempty"`
	Enabled     bool    `json:"enabled"`
	Id          string  `json:"id"`
	Matcher     string  `json:"matcher"`
	SampleData  *string `json:"sampleData,omitempty"`
}

func (p *Processor) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"description": {
			Type:        schema.TypeString,
			Description: "Name or description of the processor",
			Required:    true,
		},
		"editable": {
			Type:        schema.TypeBool,
			Description: "Indicates if the user is allowed to edit this object based on permissions and builtin property.",
			Optional:    true,
		},
		"enabled": {
			Type:        schema.TypeBool,
			Description: "Indicates if the object is active.",
			Required:    true,
		},
		"id": {
			Type:        schema.TypeString,
			Description: "Identifier of the processor. Must be unique within a stage.",
			Required:    true,
		},
		"matcher": {
			Type:        schema.TypeString,
			Description: "Matching condition to apply on incoming records.",
			Required:    true,
		},
		"sample_data": {
			Type:        schema.TypeString,
			Description: "Sample data related to the processor for documentation or testing.",
			Optional:    true,
		},
	}
}

func (p *Processor) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"description": p.Description,
		"editable":    p.Editable,
		"enabled":     p.Enabled,
		"id":          p.Id,
		"matcher":     p.Matcher,
		"sample_data": p.SampleData,
	})
}

func (p *Processor) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"description": &p.Description,
		"editable":    &p.Editable,
		"enabled":     &p.Enabled,
		"id":          &p.Id,
		"matcher":     &p.Matcher,
		"sample_data": &p.SampleData,
	})
}

type DqlProcessor struct {
	Processor
	DqlScript string `json:"dqlScript"`
}

func (p *DqlProcessor) Schema() map[string]*schema.Schema {
	s := p.Processor.Schema()
	s["dql_script"] = &schema.Schema{
		Type:        schema.TypeString,
		Description: "The DQL script to apply on the record.",
		Required:    true,
	}

	return s
}

func (p *DqlProcessor) MarshalHCL(properties hcl.Properties) error {
	if err := p.Processor.MarshalHCL(properties); err != nil {
		return err
	}
	return properties.Encode("dql_script", p.DqlScript)
}

func (p *DqlProcessor) UnmarshalHCL(decoder hcl.Decoder) error {
	if err := p.Processor.UnmarshalHCL(decoder); err != nil {
		return err
	}
	return decoder.Decode("dql_script", &p.DqlScript)
}

type FieldsAddProcessor struct {
	Processor
	Fields []FieldsAddItem `json:"fields"`
}

func (p *FieldsAddProcessor) Schema() map[string]*schema.Schema {
	s := p.Processor.Schema()
	s["field"] = &schema.Schema{
		Type:        schema.TypeList,
		Elem:        &schema.Resource{Schema: new(FieldsAddItem).Schema()},
		Description: "Field to add to the record.",
		Required:    true,
	}

	return s
}

func (p *FieldsAddProcessor) MarshalHCL(properties hcl.Properties) error {
	if err := p.Processor.MarshalHCL(properties); err != nil {
		return err
	}
	return properties.Encode("field", p.Fields)
}

func (p *FieldsAddProcessor) UnmarshalHCL(decoder hcl.Decoder) error {
	if err := p.Processor.UnmarshalHCL(decoder); err != nil {
		return err
	}
	return decoder.Decode("field", &p.Fields)
}

type FieldsAddItem struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

func (f *FieldsAddItem) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:        schema.TypeString,
			Description: "Name of the field",
			Required:    true,
		},
		"value": {
			Type:        schema.TypeString,
			Description: " Value to assign to the field",
			Required:    true,
		},
	}
}

func (f *FieldsAddItem) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"name":  f.Name,
		"value": f.Value,
	})
}

func (f *FieldsAddItem) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"name":  &f.Name,
		"value": &f.Value,
	})
}

type FieldsRemoveProcessor struct {
	Processor
	Fields []string `json:"fields"`
}

func (p *FieldsRemoveProcessor) Schema() map[string]*schema.Schema {
	s := p.Processor.Schema()
	s["field"] = &schema.Schema{
		Type:        schema.TypeList,
		Elem:        &schema.Schema{Type: schema.TypeString},
		Description: "Field to add to the record.",
		Required:    true,
	}

	return s
}

func (p *FieldsRemoveProcessor) MarshalHCL(properties hcl.Properties) error {
	if err := p.Processor.MarshalHCL(properties); err != nil {
		return err
	}
	return properties.Encode("field", p.Fields)
}

func (p *FieldsRemoveProcessor) UnmarshalHCL(decoder hcl.Decoder) error {
	if err := p.Processor.UnmarshalHCL(decoder); err != nil {
		return err
	}
	return decoder.Decode("field", &p.Fields)
}

type FieldsRenameProcessor struct {
	Processor
	Fields []FieldsRenameItem `json:"fields"`
}

func (p *FieldsRenameProcessor) Schema() map[string]*schema.Schema {
	s := p.Processor.Schema()
	s["field"] = &schema.Schema{
		Type:        schema.TypeList,
		Elem:        &schema.Resource{Schema: new(FieldsRenameItem).Schema()},
		Description: "Field to rename on the record.",
		Required:    true,
	}

	return s
}

func (p *FieldsRenameProcessor) MarshalHCL(properties hcl.Properties) error {
	if err := p.Processor.MarshalHCL(properties); err != nil {
		return err
	}
	return properties.Encode("field", p.Fields)
}

func (p *FieldsRenameProcessor) UnmarshalHCL(decoder hcl.Decoder) error {
	if err := p.Processor.UnmarshalHCL(decoder); err != nil {
		return err
	}
	return decoder.Decode("field", &p.Fields)
}

type FieldsRenameItem struct {
	FromName string `json:"fromName"`
	ToName   string `json:"toName"`
}

func (f *FieldsRenameItem) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"from_name": {
			Type:        schema.TypeString,
			Description: "The field to rename",
			Required:    true,
		},
		"to_name": {
			Type:        schema.TypeString,
			Description: "The new field name",
			Required:    true,
		},
	}
}

func (f *FieldsRenameItem) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"from_name": f.FromName,
		"to_name":   f.ToName,
	})
}

func (f *FieldsRenameItem) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"from_name": &f.FromName,
		"to_name":   &f.ToName,
	})
}

type BizEventExtractionProcessor struct {
	Processor
	EventProvider ValueAssignment `json:"eventProvider,omitempty"`
	EventType     ValueAssignment `json:"eventType,omitempty"`
}

func (ep *BizEventExtractionProcessor) Schema() map[string]*schema.Schema {
	s := ep.Processor.Schema()
	s["event_provider"] = &schema.Schema{
		Type:        schema.TypeList,
		MinItems:    1,
		MaxItems:    1,
		Elem:        &schema.Resource{Schema: new(ValueAssignment).Schema()},
		Description: "Strategy to assign a value.",
		Required:    true,
	}

	s["event_type"] = &schema.Schema{
		Type:        schema.TypeList,
		MinItems:    1,
		MaxItems:    1,
		Elem:        &schema.Resource{Schema: new(ValueAssignment).Schema()},
		Description: "Strategy to assign a value.",
		Required:    true,
	}
	return s
}

func (ep *BizEventExtractionProcessor) MarshalHCL(properties hcl.Properties) error {
	if err := ep.Processor.MarshalHCL(properties); err != nil {
		return err
	}

	return properties.EncodeAll(map[string]any{
		"event_provider": ep.EventProvider,
		"event_type":     ep.EventType,
	})
}

func (ep *BizEventExtractionProcessor) UnmarshalHCL(decoder hcl.Decoder) error {

	if err := ep.Processor.UnmarshalHCL(decoder); err != nil {
		return err
	}

	return decoder.DecodeAll(map[string]any{
		"event_provider": ep.EventProvider,
		"event_type":     ep.EventType,
	})
}

type DavisEventExtractionProcessor struct {
	Processor
	Properties []DavisEventProperty `json:"properties,omitempty"`
}

func (ep *DavisEventExtractionProcessor) Schema() map[string]*schema.Schema {

	s := ep.Processor.Schema()
	s["properties"] = &schema.Schema{
		Type:        schema.TypeList,
		Elem:        &schema.Resource{Schema: new(DavisEventProperty).Schema()},
		Description: "List of properties for the extracted davis event.",
		Required:    true,
	}
	return s
}

func (ep *DavisEventExtractionProcessor) MarshalHCL(properties hcl.Properties) error {
	if err := ep.Processor.MarshalHCL(properties); err != nil {
		return err
	}
	return properties.Encode("properties", ep.Properties)
}

func (ep *DavisEventExtractionProcessor) UnmarshalHCL(decoder hcl.Decoder) error {

	if err := ep.Processor.UnmarshalHCL(decoder); err != nil {
		return err
	}
	return decoder.Decode("field", &ep.Properties)
}

type DavisEventProperty struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func (ep *DavisEventProperty) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"key": {
			Type:        schema.TypeString,
			Description: "The key to set on the davis event.",
			Required:    true,
		},
		"value": {
			Type:        schema.TypeString,
			Description: "The value assigned to the key.",
			Required:    true,
		},
	}
}

func (ep *DavisEventProperty) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"key":   ep.Key,
		"value": ep.Value,
	})
}

func (ep *DavisEventProperty) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"key":   &ep.Key,
		"value": &ep.Value,
	})
}

type ValueAssignment struct {
	// Type Defines the actual set of fields depending on the value. See one of the following objects:
	Type string `json:"type"`
}

func (ep *ValueAssignment) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"type": {
			Type:        schema.TypeString,
			Description: "Strategy to assign a value.",
			Required:    true,
		},
	}

}

func (ep *ValueAssignment) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"type": ep.Type,
	})
}

func (ep *ValueAssignment) UnmarshalHCL(decoder hcl.Decoder) error {

	return decoder.DecodeAll(map[string]any{
		"type": ep.Type,
	})
}

type TechnologyProcessor struct {
	Processor
	TechnologyId string `json:"technologyId"`
}

func (p *TechnologyProcessor) Schema() map[string]*schema.Schema {
	s := p.Processor.Schema()
	s["technology_id"] = &schema.Schema{
		Type:        schema.TypeString,
		Description: "Identifier of the processor. Must be unique within a stage.",
		Required:    true,
	}

	return s
}

func (p *TechnologyProcessor) MarshalHCL(properties hcl.Properties) error {
	if err := p.Processor.MarshalHCL(properties); err != nil {
		return err
	}
	return properties.Encode("technology_id", p.TechnologyId)
}

func (p *TechnologyProcessor) UnmarshalHCL(decoder hcl.Decoder) error {
	if err := p.Processor.UnmarshalHCL(decoder); err != nil {
		return err
	}
	return decoder.Decode("technology_id", &p.TechnologyId)
}

type SqlxProcessor struct {
	Processor
	SqlxScript string `json:"sqlxScript"`
}

func (p *SqlxProcessor) Schema() map[string]*schema.Schema {
	s := p.Processor.Schema()
	s["sqlx_script"] = &schema.Schema{
		Type:        schema.TypeString,
		Description: "The SQLX script to apply on the record.",
		Required:    true,
	}

	return s
}

func (p *SqlxProcessor) MarshalHCL(properties hcl.Properties) error {
	if err := p.Processor.MarshalHCL(properties); err != nil {
		return err
	}
	return properties.Encode("sqlx_script", p.SqlxScript)
}

func (p *SqlxProcessor) UnmarshalHCL(decoder hcl.Decoder) error {
	if err := p.Processor.UnmarshalHCL(decoder); err != nil {
		return err
	}
	return decoder.Decode("sqlx_script", &p.SqlxScript)
}

type CounterMetricExtractionProcessor struct {
	Processor
	Dimensions *[]string `json:"dimensions,omitempty"`
	MetricKey  string    `json:"metricKey"`
}

func (p *CounterMetricExtractionProcessor) Schema() map[string]*schema.Schema {
	s := p.Processor.Schema()
	s["dimensions"] = &schema.Schema{
		Type:        schema.TypeList,
		Elem:        &schema.Schema{Type: schema.TypeString},
		Description: "List of dimensions to add to the metric.",
		Optional:    true,
	}
	s["metric_key"] = &schema.Schema{
		Type:        schema.TypeString,
		Description: "The key of the metric to write.",
		Required:    true,
	}

	return s
}

func (p *CounterMetricExtractionProcessor) MarshalHCL(properties hcl.Properties) error {
	if err := p.Processor.MarshalHCL(properties); err != nil {
		return err
	}
	return properties.EncodeAll(map[string]any{
		"dimensions": p.Dimensions,
		"metric_key": p.MetricKey,
	})
}

func (p *CounterMetricExtractionProcessor) UnmarshalHCL(decoder hcl.Decoder) error {
	if err := p.Processor.UnmarshalHCL(decoder); err != nil {
		return err
	}
	return decoder.DecodeAll(map[string]any{
		"dimensions": p.Dimensions,
		"metric_key": p.MetricKey,
	})
}

type ValueMetricExtractionProcessor struct {
	Processor
	Dimensions *[]string `json:"dimensions,omitempty"`
	Field      string    `json:"field"`
	MetricKey  string    `json:"metricKey"`
}

func (p *ValueMetricExtractionProcessor) Schema() map[string]*schema.Schema {
	s := p.Processor.Schema()
	s["dimensions"] = &schema.Schema{
		Type:        schema.TypeList,
		Elem:        &schema.Schema{Type: schema.TypeString},
		Description: "List of dimensions to add to the metric.",
		Optional:    true,
	}
	s["field"] = &schema.Schema{
		Type:        schema.TypeString,
		Description: "The field to extract the value for the metric.",
		Required:    true,
	}
	s["metric_key"] = &schema.Schema{
		Type:        schema.TypeString,
		Description: "The key of the metric to write.",
		Required:    true,
	}

	return s
}

func (p *ValueMetricExtractionProcessor) MarshalHCL(properties hcl.Properties) error {
	if err := p.Processor.MarshalHCL(properties); err != nil {
		return err
	}
	return properties.EncodeAll(map[string]any{
		"dimensions": p.Dimensions,
		"field":      p.Field,
		"metric_key": p.MetricKey,
	})
}

func (p *ValueMetricExtractionProcessor) UnmarshalHCL(decoder hcl.Decoder) error {
	if err := p.Processor.UnmarshalHCL(decoder); err != nil {
		return err
	}
	return decoder.DecodeAll(map[string]any{
		"dimensions": p.Dimensions,
		"field":      p.Field,
		"metric_key": p.MetricKey,
	})
}

type BucketAssignmentProcessor struct {
	Processor
	BucketName string `json:"bucketName"`
}

func (p *BucketAssignmentProcessor) Schema() map[string]*schema.Schema {
	s := p.Processor.Schema()
	s["bucket_name"] = &schema.Schema{
		Type:        schema.TypeString,
		Description: "Bucket that is assigned when the record is matched.",
		Required:    true,
	}

	return s
}

func (p *BucketAssignmentProcessor) MarshalHCL(properties hcl.Properties) error {
	if err := p.Processor.MarshalHCL(properties); err != nil {
		return err
	}
	return properties.EncodeAll(map[string]any{
		"bucket_name": p.BucketName,
	})
}

func (p *BucketAssignmentProcessor) UnmarshalHCL(decoder hcl.Decoder) error {
	if err := p.Processor.UnmarshalHCL(decoder); err != nil {
		return err
	}
	return decoder.DecodeAll(map[string]any{
		"bucket_name": p.BucketName,
	})
}

type NoStorageProcessor struct {
	Processor
}

func (p *NoStorageProcessor) Schema() map[string]*schema.Schema {
	return p.Processor.Schema()
}

func (p *NoStorageProcessor) MarshalHCL(properties hcl.Properties) error {
	return p.Processor.MarshalHCL(properties)

}

func (p *NoStorageProcessor) UnmarshalHCL(decoder hcl.Decoder) error {
	return p.Processor.UnmarshalHCL(decoder)
}

type SecurityContextProcessor struct {
	Processor
	Value ValueAssignment `json:"value"`
}

func (p *SecurityContextProcessor) Schema() map[string]*schema.Schema {
	s := p.Processor.Schema()
	s["value"] = &schema.Schema{
		Type:        schema.TypeList,
		MinItems:    1,
		MaxItems:    1,
		Elem:        &schema.Resource{Schema: new(ValueAssignment).Schema()},
		Description: "Strategy to assign a value.",
		Required:    true,
	}

	return s
}

func (p *SecurityContextProcessor) MarshalHCL(properties hcl.Properties) error {
	if err := p.Processor.MarshalHCL(properties); err != nil {
		return err
	}
	return properties.EncodeAll(map[string]any{
		"value": p.Value,
	})
}

func (p *SecurityContextProcessor) UnmarshalHCL(decoder hcl.Decoder) error {
	if err := p.Processor.UnmarshalHCL(decoder); err != nil {
		return err
	}
	return decoder.DecodeAll(map[string]any{
		"value": p.Value,
	})
}
