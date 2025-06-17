package openpipeline

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	DqlProcessorType          = "dql"
	FieldsAddProcessorType    = "fieldsAdd"
	FieldsRemoveProcessorType = "fieldsRemove"
	FieldsRenameProcessorType = "fieldsRename"
	DropProcessorType         = "drop"

	CounterMetricProcessorType = "counterMetric"
	ValueMetricProcessorType   = "valueMetric"

	DavisEventExtractionProcessorType = "davis"
	BizEventExtractionProcessorType   = "bizevent"

	SecurityContextProcessorType = "securityContext"

	NoStorageStageProcessorType        = "noStorage"
	BucketAssignmentStageProcessorType = "bucketAssignment"

	TechnologyProcessorType = "technology"
	SqlxProcessorType       = "sqlx"
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
		"enabled": {
			Type:        schema.TypeBool,
			Description: "Indicates if the object is active",
			Required:    true,
		},
		"id": {
			Type:        schema.TypeString,
			Description: "Identifier of the processor. Must be unique within a stage.",
			Required:    true,
		},
		"matcher": {
			Type:        schema.TypeString,
			Description: "Matching condition to apply on incoming records",
			Required:    true,
		},
		"sample_data": {
			Type:        schema.TypeString,
			Description: "Sample data related to the processor for documentation or testing",
			Optional:    true,
		},
	}
}

func (p *Processor) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"description": p.Description,
		"enabled":     p.Enabled,
		"id":          p.Id,
		"matcher":     p.Matcher,
		"sample_data": p.SampleData,
	})
}

func (p *Processor) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"description": &p.Description,
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
		Description: "The DQL script to apply on the record",
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

func (ep DqlProcessor) MarshalJSON() ([]byte, error) {
	type dqlProcessor DqlProcessor
	return MarshalAsJSONWithType((dqlProcessor)(ep), DqlProcessorType)
}

type FieldsAddProcessor struct {
	Processor
	Fields []*FieldsAddItem `json:"fields"`
}

func (p *FieldsAddProcessor) Schema() map[string]*schema.Schema {
	s := p.Processor.Schema()
	s["field"] = &schema.Schema{
		Type:        schema.TypeList,
		Elem:        &schema.Resource{Schema: new(FieldsAddItem).Schema()},
		Description: "Field to add to the record",
		Required:    true,
	}

	return s
}

func (p *FieldsAddProcessor) MarshalHCL(properties hcl.Properties) error {
	if err := p.Processor.MarshalHCL(properties); err != nil {
		return err
	}
	return properties.EncodeSlice("field", p.Fields)
}

func (p *FieldsAddProcessor) UnmarshalHCL(decoder hcl.Decoder) error {
	if err := p.Processor.UnmarshalHCL(decoder); err != nil {
		return err
	}
	return decoder.DecodeSlice("field", &p.Fields)
}

func (ep FieldsAddProcessor) MarshalJSON() ([]byte, error) {
	type fieldsAddProcessor FieldsAddProcessor
	return MarshalAsJSONWithType((fieldsAddProcessor)(ep), FieldsAddProcessorType)
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
	s["fields"] = &schema.Schema{
		Type:        schema.TypeList,
		Elem:        &schema.Schema{Type: schema.TypeString},
		Description: "Field to add to the record",
		Required:    true,
	}

	return s
}

func (p *FieldsRemoveProcessor) MarshalHCL(properties hcl.Properties) error {
	if err := p.Processor.MarshalHCL(properties); err != nil {
		return err
	}
	return properties.Encode("fields", p.Fields)
}

func (p *FieldsRemoveProcessor) UnmarshalHCL(decoder hcl.Decoder) error {
	if err := p.Processor.UnmarshalHCL(decoder); err != nil {
		return err
	}
	return decoder.Decode("fields", &p.Fields)
}

func (ep FieldsRemoveProcessor) MarshalJSON() ([]byte, error) {
	type fieldsRemoveProcessor FieldsRemoveProcessor
	return MarshalAsJSONWithType((fieldsRemoveProcessor)(ep), FieldsRemoveProcessorType)
}

type FieldsRenameProcessor struct {
	Processor
	Fields []*FieldsRenameItem `json:"fields"`
}

func (p *FieldsRenameProcessor) Schema() map[string]*schema.Schema {
	s := p.Processor.Schema()
	s["field"] = &schema.Schema{
		Type:        schema.TypeList,
		Elem:        &schema.Resource{Schema: new(FieldsRenameItem).Schema()},
		Description: "Field to rename on the record",
		Required:    true,
	}

	return s
}

func (p *FieldsRenameProcessor) MarshalHCL(properties hcl.Properties) error {
	if err := p.Processor.MarshalHCL(properties); err != nil {
		return err
	}
	return properties.EncodeSlice("field", p.Fields)
}

func (p *FieldsRenameProcessor) UnmarshalHCL(decoder hcl.Decoder) error {
	if err := p.Processor.UnmarshalHCL(decoder); err != nil {
		return err
	}
	return decoder.DecodeSlice("field", &p.Fields)
}

func (ep FieldsRenameProcessor) MarshalJSON() ([]byte, error) {
	type fieldsRenameProcessor FieldsRenameProcessor
	return MarshalAsJSONWithType((fieldsRenameProcessor)(ep), FieldsRenameProcessorType)
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

type DropProcessor struct {
	Processor
}

func (p *DropProcessor) Schema() map[string]*schema.Schema {
	return p.Processor.Schema()
}

func (p *DropProcessor) MarshalHCL(properties hcl.Properties) error {
	return p.Processor.MarshalHCL(properties)
}

func (p *DropProcessor) UnmarshalHCL(decoder hcl.Decoder) error {
	return p.Processor.UnmarshalHCL(decoder)
}

func (ep DropProcessor) MarshalJSON() ([]byte, error) {
	type dropProcessor DropProcessor
	return MarshalAsJSONWithType((dropProcessor)(ep), DropProcessorType)
}

type FieldExtraction struct {
	Fields   []string `json:"fields"`
	Semantic string   `json:"semantic"`
}

func (ep *FieldExtraction) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"fields": {
			Type:        schema.TypeList,
			Elem:        &schema.Schema{Type: schema.TypeString},
			Description: "Unique fields to include/exclude in the extracted record",
			Optional:    true,
		},
		"semantic": {
			Type:        schema.TypeString,
			Description: "Defines how the fields of the source record should be extracted",
			Required:    true,
		},
	}

}

func (ep *FieldExtraction) MarshalHCL(properties hcl.Properties) error {
	if err := properties.EncodeSlice("fields", ep.Fields); err != nil {
		return err
	}

	return properties.Encode("semantic", ep.Semantic)
}

func (ep *FieldExtraction) UnmarshalHCL(decoder hcl.Decoder) error {
	if err := decoder.DecodeSlice("fields", &ep.Fields); err != nil {
		return err
	}

	return decoder.Decode("semantic", &ep.Semantic)
}

type BizEventExtractionProcessor struct {
	Processor
	EventProvider   *ValueAssignment `json:"eventProvider,omitempty"`
	EventType       *ValueAssignment `json:"eventType,omitempty"`
	FieldExtraction *FieldExtraction `json:"fieldExtraction,omitempty"`
}

func (ep *BizEventExtractionProcessor) Schema() map[string]*schema.Schema {
	s := ep.Processor.Schema()
	s["event_provider"] = &schema.Schema{
		Type:        schema.TypeList,
		MinItems:    1,
		MaxItems:    1,
		Elem:        &schema.Resource{Schema: new(ValueAssignment).Schema()},
		Description: "Strategy to assign a value",
		Required:    true,
	}

	s["event_type"] = &schema.Schema{
		Type:        schema.TypeList,
		MinItems:    1,
		MaxItems:    1,
		Elem:        &schema.Resource{Schema: new(ValueAssignment).Schema()},
		Description: "Strategy to assign a value",
		Required:    true,
	}

	s["field_extraction"] = &schema.Schema{
		Type:        schema.TypeList,
		MinItems:    1,
		MaxItems:    1,
		Elem:        &schema.Resource{Schema: new(FieldExtraction).Schema()},
		Description: "Definition of the field extraction",
		Optional:    true,
	}

	return s
}

func (ep *BizEventExtractionProcessor) MarshalHCL(properties hcl.Properties) error {
	if err := ep.Processor.MarshalHCL(properties); err != nil {
		return err
	}

	return properties.EncodeAll(map[string]any{
		"event_provider":   ep.EventProvider,
		"event_type":       ep.EventType,
		"field_extraction": ep.FieldExtraction,
	})
}

func (ep *BizEventExtractionProcessor) UnmarshalHCL(decoder hcl.Decoder) error {
	if err := ep.Processor.UnmarshalHCL(decoder); err != nil {
		return err
	}

	return decoder.DecodeAll(map[string]any{
		"event_provider":   &ep.EventProvider,
		"event_type":       &ep.EventType,
		"field_extraction": &ep.FieldExtraction,
	})
}

func (ep BizEventExtractionProcessor) MarshalJSON() ([]byte, error) {
	type bizEventExtractionProcessor BizEventExtractionProcessor
	return MarshalAsJSONWithType((bizEventExtractionProcessor)(ep), BizEventExtractionProcessorType)
}

type DavisEventExtractionProcessor struct {
	Processor
	Properties []*DavisEventProperty `json:"properties,omitempty"`
}

func (ep *DavisEventExtractionProcessor) Schema() map[string]*schema.Schema {

	s := ep.Processor.Schema()
	s["properties"] = &schema.Schema{
		Type:        schema.TypeList,
		Elem:        &schema.Resource{Schema: new(DavisEventProperty).Schema()},
		Description: "List of properties for the extracted davis event",
		Required:    true,
	}
	return s
}

func (ep *DavisEventExtractionProcessor) MarshalHCL(properties hcl.Properties) error {
	if err := ep.Processor.MarshalHCL(properties); err != nil {
		return err
	}
	return properties.EncodeSlice("properties", ep.Properties)
}

func (ep *DavisEventExtractionProcessor) UnmarshalHCL(decoder hcl.Decoder) error {

	if err := ep.Processor.UnmarshalHCL(decoder); err != nil {
		return err
	}
	return decoder.DecodeSlice("properties", &ep.Properties)
}

func (p DavisEventExtractionProcessor) MarshalJSON() ([]byte, error) {
	type davisEventExtractionProcessor DavisEventExtractionProcessor
	return MarshalAsJSONWithType((davisEventExtractionProcessor)(p), DavisEventExtractionProcessorType)
}

type DavisEventProperty struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func (ep *DavisEventProperty) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"key": {
			Type:        schema.TypeString,
			Description: "The key to set on the davis event",
			Required:    true,
		},
		"value": {
			Type:        schema.TypeString,
			Description: "The value assigned to the key",
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
	Type     string  `json:"type"`
	Field    *string `json:"field"`
	Constant *string `json:"constant,omitempty"`
}

func (ep *ValueAssignment) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"type": {
			Type:        schema.TypeString,
			Description: "Strategy to assign a value",
			Required:    true,
		},
		"field": {
			Type:        schema.TypeString,
			Description: "Strategy to assign a value",
			Optional:    true,
		},
		"constant": {
			Type:        schema.TypeString,
			Description: "Strategy to assign a value",
			Optional:    true,
		},
	}

}

func (ep *ValueAssignment) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"type":     ep.Type,
		"field":    ep.Field,
		"constant": ep.Constant,
	})
}

func (ep *ValueAssignment) UnmarshalHCL(decoder hcl.Decoder) error {

	return decoder.DecodeAll(map[string]any{
		"type":     &ep.Type,
		"field":    &ep.Field,
		"constant": &ep.Constant,
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
		Description: "The reference identifier to a specific technology. This technology is applied on the record.",
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

func (ep TechnologyProcessor) MarshalJSON() ([]byte, error) {
	type technologyProcessor TechnologyProcessor
	return MarshalAsJSONWithType((technologyProcessor)(ep), TechnologyProcessorType)
}

type SqlxProcessor struct {
	Processor
	SqlxScript string `json:"sqlxScript"`
}

func (p *SqlxProcessor) Schema() map[string]*schema.Schema {
	s := p.Processor.Schema()
	s["sqlx_script"] = &schema.Schema{
		Type:        schema.TypeString,
		Description: "The SQLX script to apply on the record",
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
	Dimensions []string `json:"dimensions,omitempty"`
	MetricKey  string   `json:"metricKey"`
}

func (p *CounterMetricExtractionProcessor) Schema() map[string]*schema.Schema {
	s := p.Processor.Schema()
	s["dimensions"] = &schema.Schema{
		Type:        schema.TypeList,
		Elem:        &schema.Schema{Type: schema.TypeString},
		Description: "List of dimensions to add to the metric",
		Optional:    true,
	}
	s["metric_key"] = &schema.Schema{
		Type:        schema.TypeString,
		Description: "The key of the metric to write",
		Required:    true,
	}

	return s
}

func (p *CounterMetricExtractionProcessor) MarshalHCL(properties hcl.Properties) error {
	if err := p.Processor.MarshalHCL(properties); err != nil {
		return err
	}

	if err := properties.Encode("dimensions", p.Dimensions); err != nil {
		return err
	}

	return properties.Encode("metric_key", p.MetricKey)
}

func (p *CounterMetricExtractionProcessor) UnmarshalHCL(decoder hcl.Decoder) error {
	if err := p.Processor.UnmarshalHCL(decoder); err != nil {
		return err
	}

	if err := decoder.Decode("dimensions", &p.Dimensions); err != nil {
		return err
	}

	return decoder.Decode("metric_key", &p.MetricKey)
}

func (ep CounterMetricExtractionProcessor) MarshalJSON() ([]byte, error) {
	type counterMetricExtractionProcessor CounterMetricExtractionProcessor
	return MarshalAsJSONWithType((counterMetricExtractionProcessor)(ep), CounterMetricProcessorType)
}

type ValueMetricExtractionProcessor struct {
	Processor
	Dimensions []string `json:"dimensions,omitempty"`
	Field      string   `json:"field"`
	MetricKey  string   `json:"metricKey"`
}

func (p *ValueMetricExtractionProcessor) Schema() map[string]*schema.Schema {
	s := p.Processor.Schema()
	s["dimensions"] = &schema.Schema{
		Type:        schema.TypeList,
		Elem:        &schema.Schema{Type: schema.TypeString},
		Description: "List of dimensions to add to the metric",
		Optional:    true,
	}
	s["field"] = &schema.Schema{
		Type:        schema.TypeString,
		Description: "The field to extract the value for the metric",
		Required:    true,
	}
	s["metric_key"] = &schema.Schema{
		Type:        schema.TypeString,
		Description: "The key of the metric to write",
		Required:    true,
	}

	return s
}

func (p *ValueMetricExtractionProcessor) MarshalHCL(properties hcl.Properties) error {
	if err := p.Processor.MarshalHCL(properties); err != nil {
		return err
	}

	if err := properties.Encode("dimensions", p.Dimensions); err != nil {
		return err
	}

	return properties.EncodeAll(map[string]any{
		"field":      p.Field,
		"metric_key": p.MetricKey,
	})
}

func (p *ValueMetricExtractionProcessor) UnmarshalHCL(decoder hcl.Decoder) error {
	if err := p.Processor.UnmarshalHCL(decoder); err != nil {
		return err
	}

	if err := decoder.Decode("dimensions", &p.Dimensions); err != nil {
		return err
	}

	return decoder.DecodeAll(map[string]any{
		"field":      &p.Field,
		"metric_key": &p.MetricKey,
	})
}

func (ep ValueMetricExtractionProcessor) MarshalJSON() ([]byte, error) {
	type valueMetricExtractionProcessor ValueMetricExtractionProcessor
	return MarshalAsJSONWithType((valueMetricExtractionProcessor)(ep), ValueMetricProcessorType)
}

type BucketAssignmentProcessor struct {
	Processor
	BucketName string `json:"bucketName"`
}

func (p *BucketAssignmentProcessor) Schema() map[string]*schema.Schema {
	s := p.Processor.Schema()
	s["bucket_name"] = &schema.Schema{
		Type:        schema.TypeString,
		Description: "Bucket that is assigned when the record is matched",
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
		"bucket_name": &p.BucketName,
	})
}

func (ep BucketAssignmentProcessor) MarshalJSON() ([]byte, error) {
	type bucketAssignmentProcessor BucketAssignmentProcessor
	return MarshalAsJSONWithType((bucketAssignmentProcessor)(ep), BucketAssignmentStageProcessorType)
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

func (ep NoStorageProcessor) MarshalJSON() ([]byte, error) {
	type noStorageProcessor NoStorageProcessor
	return MarshalAsJSONWithType((noStorageProcessor)(ep), NoStorageStageProcessorType)
}

type SecurityContextProcessor struct {
	Processor
	Value *ValueAssignment `json:"value"`
}

func (p *SecurityContextProcessor) Schema() map[string]*schema.Schema {
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

func (p *SecurityContextProcessor) MarshalHCL(properties hcl.Properties) error {
	if err := p.Processor.MarshalHCL(properties); err != nil {
		return err
	}
	return properties.Encode("value", p.Value)
}

func (p *SecurityContextProcessor) UnmarshalHCL(decoder hcl.Decoder) error {
	if err := p.Processor.UnmarshalHCL(decoder); err != nil {
		return err
	}
	return decoder.Decode("value", &p.Value)
}

func (ep SecurityContextProcessor) MarshalJSON() ([]byte, error) {
	type securityContextProcessor SecurityContextProcessor
	return MarshalAsJSONWithType((securityContextProcessor)(ep), SecurityContextProcessorType)
}
