package boundaries

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type PolicyBoundary struct {
	// LevelType string `json:"levelType"`
	// LevelID   string `json:"levelId"`
	Name  string `json:"name"`
	Query string `json:"boundaryQuery"`
}

func (me *PolicyBoundary) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The name of the policy",
		},
		"query": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The boundary query",
		},
	}
}

func (me *PolicyBoundary) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"name":  me.Name,
		"query": me.Query,
	})
}

func (me *PolicyBoundary) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"name":  &me.Name,
		"query": &me.Query,
	})
}
