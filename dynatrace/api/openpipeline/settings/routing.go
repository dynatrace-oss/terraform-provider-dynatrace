package openpipeline

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Routing struct {
	// Type Defines the actual set of fields depending on the value. See one of the following objects:
	//
	// * `static` -> StaticRouting
	// * `dynamic` -> DynamicRouting
	Type RoutingType `json:"type"`
}

type RoutingType string

type RoutingTable struct {
	// CatchAllPipeline The default pipeline records are routed into if no dynamic routing entries apply.
	CatchAllPipeline RoutingTableEntryTarget `json:"catchAllPipeline"`

	// Editable Indicates if the user is allowed to edit this object based on permissions and builtin property.
	Editable *bool `json:"editable,omitempty"`

	// Entries List of all dynamic routes.
	Entries RoutingTableEntries `json:"entries"`
}

func (t *RoutingTable) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"catch_all_pipeline": {
			Type:        schema.TypeList,
			Description: "todo",
			MinItems:    1,
			MaxItems:    1,
			Elem:        &schema.Resource{Schema: new(RoutingTableEntryTarget).Schema()},
			Optional:    true,
		},
		"editable": {
			Type:        schema.TypeBool,
			Description: "todo",
			Optional:    true,
		},
		"entries": {
			Type:        schema.TypeList,
			Description: "todo",
			MinItems:    1,
			MaxItems:    1,
			Elem:        &schema.Resource{Schema: new(RoutingTableEntries).Schema()},
			Optional:    true,
		},
	}
}

func (t *RoutingTable) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"catch_all_pipeline": t.CatchAllPipeline,
		"editable":           t.Editable,
		"entries":            t.Entries,
	})
}

func (t *RoutingTable) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"catch_all_pipeline": t.CatchAllPipeline,
		"editable":           t.Editable,
		"entries":            t.Entries,
	})
}

type RoutingTableEntryTarget struct {
	// Editable Indicates if the user is allowed to edit this object based on permissions and builtin property.
	Editable *bool `json:"editable,omitempty"`

	// PipelineId Identifier of the pipeline.
	PipelineId string `json:"pipelineId"`
}

func (t *RoutingTableEntryTarget) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"editable": {
			Type:        schema.TypeBool,
			Description: "todo",
			Required:    true,
		},
		"pipeline_id": {
			Type:        schema.TypeString,
			Description: "todo",
			Required:    true,
		},
	}
}

func (t *RoutingTableEntryTarget) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"editable":    t.Editable,
		"pipeline_id": t.PipelineId,
	})
}

func (t *RoutingTableEntryTarget) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"editable":    t.Editable,
		"pipeline_id": t.PipelineId,
	})
}

type RoutingTableEntries struct {
	Entries []RoutingTableEntry `json:"entries,omitempty"`
}

func (e *RoutingTableEntries) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"entry": {
			Type:        schema.TypeSet,
			Description: "todo",
			Elem:        &schema.Resource{Schema: new(RoutingTableEntry).Schema()},
			Optional:    true,
		},
	}
}

func (e *RoutingTableEntries) MarshalHCL(properties hcl.Properties) error {
	return properties.Encode("entries", e.Entries)
}

func (e *RoutingTableEntries) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.Decode("entries", &e.Entries)
}

type RoutingTableEntry struct {
	// Builtin Indicates if the object is provided by Dynatrace or customer defined.
	Builtin *bool `json:"builtin,omitempty"`

	// Editable Indicates if the user is allowed to edit this object based on permissions and builtin property.
	Editable *bool `json:"editable,omitempty"`

	// Enabled Indicates if the object is active.
	Enabled bool `json:"enabled"`

	// Matcher Matching condition to apply on incoming records.
	Matcher string `json:"matcher"`

	// Note Unique note describing the dynamic route.
	Note string `json:"note"`

	// PipelineId Identifier of the pipeline the record is routed into.
	PipelineId string `json:"pipelineId"`
}

func (e *RoutingTableEntry) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"builtin": {
			Type:        schema.TypeBool,
			Description: "todo",
			Required:    true,
		},
		"editable": {
			Type:        schema.TypeBool,
			Description: "todo",
			Required:    true,
		},
		"enabled": {
			Type:        schema.TypeBool,
			Description: "todo",
			Required:    true,
		},
		"matcher": {
			Type:        schema.TypeString,
			Description: "todo",
			Required:    true,
		},
		"note": {
			Type:        schema.TypeString,
			Description: "todo",
			Required:    true,
		},
		"pipeline_id": {
			Type:        schema.TypeString,
			Description: "todo",
			Required:    true,
		},
	}
}

func (t *RoutingTableEntry) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"builtin":     t.Builtin,
		"editable":    t.Editable,
		"enabled":     t.Enabled,
		"matcher":     t.Matcher,
		"note":        t.Note,
		"pipeline_id": t.PipelineId,
	})
}

func (t *RoutingTableEntry) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"builtin":     t.Builtin,
		"editable":    t.Editable,
		"enabled":     t.Enabled,
		"matcher":     t.Matcher,
		"note":        t.Note,
		"pipeline_id": t.PipelineId,
	})
}
