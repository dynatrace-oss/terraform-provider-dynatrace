package generic

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Settings struct {
	Scope    string `json:"-" scope:"scope"` // The scope of this setting
	SchemaID string `json:"schemaId"`
	Value    string `json:"value"`
}

func (me *Settings) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"value": {
			Type:             schema.TypeString,
			Required:         true,
			DiffSuppressFunc: hcl.SuppressJSONorEOT,
		},
		"schema": {
			Type:     schema.TypeString,
			Required: true,
		},
		"scope": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
	}
}

func (me *Settings) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"scope":  me.Scope,
		"value":  me.Value,
		"schema": me.SchemaID,
	})
}

func (me *Settings) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"scope":  &me.Scope,
		"value":  &me.Value,
		"schema": &me.SchemaID,
	})
}
