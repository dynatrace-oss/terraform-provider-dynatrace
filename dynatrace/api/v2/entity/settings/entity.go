package entity

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Entity struct {
	EntityId    *string `json:"entityId,omitempty"`    // The ID of the entity.
	Type        *string `json:"type,omitempty"`        // The type of the entity.
	DisplayName *string `json:"displayName,omitempty"` // The name of the entity, displayed in the UI.
	Tags        Tags    `json:"tags,omitempty"`        // A set of tags assigned to the entity.
}

func (me *Entity) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"entity_id": {
			Type:        schema.TypeString,
			Description: "The ID of the entity.",
			Optional:    true,
		},
		"type": {
			Type:        schema.TypeString,
			Description: "The type of the entity.",
			Optional:    true,
		},
		"display_name": {
			Type:        schema.TypeString,
			Description: "The name of the entity, displayed in the UI.",
			Optional:    true,
		},
		"tags": {
			Type:        schema.TypeList,
			Description: "A set of tags assigned to the entity.",
			Optional:    true,
			Elem: &schema.Resource{
				Schema: new(Tags).Schema(),
			},
		},
	}
}

func (me *Entity) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"entity_id":    me.EntityId,
		"type":         me.Type,
		"display_name": me.DisplayName,
		"tags":         me.Tags,
	})
}

func (me *Entity) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"entity_id":    &me.EntityId,
		"type":         &me.Type,
		"display_name": &me.DisplayName,
		"tags":         &me.Tags,
	})
}
