package iam

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Policy struct {
	// UUID           string   `json:"uuid,omitempty"`
	Name           string   `json:"name"`
	Tags           []string `json:"tags,omitempty"`
	Description    string   `json:"description,omitempty"`
	StatementQuery string   `json:"statementQuery"`
}

func (me *Policy) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:     schema.TypeString,
			Required: true,
		},
		"description": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"statement_query": {
			Type:     schema.TypeString,
			Required: true,
		},
	}
}

func (me *Policy) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"name":            me.Name,
		"description":     me.Description,
		"statement_query": me.StatementQuery,
	})
}

func (me *Policy) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"name":            &me.Name,
		"description":     &me.Description,
		"statement_query": &me.StatementQuery,
	})
}
