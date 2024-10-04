package openpipeline

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type RoutingTable struct {
	// CatchAllPipeline The default pipeline records are routed into if no dynamic routing entries apply.
	CatchAllPipeline *RoutingTableEntryTarget `json:"catchAllPipeline"`

	// Editable Indicates if the user is allowed to edit this object based on permissions and builtin property.
	Editable *bool `json:"editable,omitempty"`

	// Entries List of all dynamic routes.
	Entries []*RoutingTableEntry `json:"entries"`
}

func (t *RoutingTable) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"entry": {
			Type:        schema.TypeList,
			Description: "Dynamic routing entry",
			Elem:        &schema.Resource{Schema: new(RoutingTableEntry).Schema()},
			Optional:    true,
		},
	}
}

func (t *RoutingTable) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeSlice("entry", t.Entries)
}

func (t *RoutingTable) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeSlice("entry", &t.Entries)
}

func (t *RoutingTable) RemoveFixed() {
	filteredEntries := []*RoutingTableEntry{}
	for _, entry := range t.Entries {
		if !entry.IsFixed() {
			filteredEntries = append(filteredEntries, entry)
		}
	}
	t.Entries = filteredEntries
}

type RoutingTableEntryTarget struct {
	// Editable Indicates if the user is allowed to edit this object based on permissions and builtin property.
	Editable *bool `json:"editable,omitempty"`

	// PipelineId Identifier of the pipeline.
	PipelineId string `json:"pipelineId"`
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
		"enabled": {
			Type:        schema.TypeBool,
			Description: "Indicates if the object is active",
			Required:    true,
		},
		"matcher": {
			Type:        schema.TypeString,
			Description: "Matching condition to apply on incoming records",
			Required:    true,
		},
		"note": {
			Type:        schema.TypeString,
			Description: "Unique note describing the dynamic route",
			Required:    true,
		},
		"pipeline_id": {
			Type:        schema.TypeString,
			Description: "Identifier of the pipeline the record is routed into",
			Required:    true,
		},
	}
}

func (t *RoutingTableEntry) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"enabled":     t.Enabled,
		"matcher":     t.Matcher,
		"note":        t.Note,
		"pipeline_id": t.PipelineId,
	})
}

func (t *RoutingTableEntry) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"enabled":     &t.Enabled,
		"matcher":     &t.Matcher,
		"note":        &t.Note,
		"pipeline_id": &t.PipelineId,
	})
}

func (t *RoutingTableEntry) IsFixed() bool {
	return t.Builtin != nil && *t.Builtin
}

const (
	StaticRoutingType  = "static"
	DynamicRoutingType = "dynamic"
)

type Routing struct {
	Type       string  `json:"type"`
	PipelineId *string `json:"pipelineId,omitempty"`
}

func (ep *Routing) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"type": {
			Type:        schema.TypeString,
			Description: "Type of routing, static or dynamic",
			Required:    true,
		},
		"pipeline_id": {
			Type:        schema.TypeString,
			Description: "Pipeline ID of the static routing",
			Optional:    true,
		},
	}
}

func (ep *Routing) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"type":        ep.Type,
		"pipeline_id": ep.PipelineId,
	})
}

func (ep *Routing) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"type":        &ep.Type,
		"pipeline_id": &ep.PipelineId,
	})
}
