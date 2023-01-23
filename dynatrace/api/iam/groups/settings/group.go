package groups

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Group struct {
	Name                     string      `json:"name"`
	Description              string      `json:"description"`
	FederatedAttributeValues []string    `json:"federatedAttributeValues"`
	Permissions              Permissions `json:"-"`
}

func (me *Group) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:     schema.TypeString,
			Required: true,
		},
		"description": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"federated_attribute_values": {
			Type:     schema.TypeSet,
			Optional: true,
			MinItems: 1,
			Elem:     &schema.Schema{Type: schema.TypeString},
		},
		"permissions": {
			Type:     schema.TypeList,
			Optional: true,
			MinItems: 1,
			MaxItems: 1,
			Elem:     &schema.Resource{Schema: new(Permissions).Schema()},
		},
	}
}

func (me *Group) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"name":                       me.Name,
		"description":                me.Description,
		"federated_attribute_values": me.FederatedAttributeValues,
		"permissions":                me.Permissions,
	})
}

func (me *Group) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"name":                       &me.Name,
		"description":                &me.Description,
		"federated_attribute_values": &me.FederatedAttributeValues,
		"permissions":                &me.Permissions,
	})
}
