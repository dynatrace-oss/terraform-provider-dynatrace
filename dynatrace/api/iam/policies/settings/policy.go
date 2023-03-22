package policies

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Policy struct {
	Name           string   `json:"name"`
	Tags           []string `json:"tags"`
	Description    string   `json:"description,omitempty"`
	StatementQuery string   `json:"statementQuery"`
	Account        string   `json:"-"`
	Environment    string   `json:"-"`
}

func (me *Policy) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The name of the policy",
		},
		"description": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "An optional description text for the policy",
		},
		"tags": {
			Type:        schema.TypeSet,
			Optional:    true,
			MinItems:    1,
			Description: "Tags for this policy",
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
		"statement_query": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The Statement Query of the policy",
		},
		"account": {
			Type:          schema.TypeString,
			Optional:      true,
			ConflictsWith: []string{"environment"},
			AtLeastOneOf:  []string{"environment", "account"},
			ForceNew:      true,
			Description:   "The UUID of the account (`urn:dtaccount:<account-uuid>`) in case the policy should be applied to all environments of this account. The prefix `urn:dtaccount:` MUST be omitted here.",
		},
		"environment": {
			Type:          schema.TypeString,
			Optional:      true,
			ConflictsWith: []string{"account"},
			AtLeastOneOf:  []string{"environment", "account"},
			ForceNew:      true,
			Description:   "The ID of the environment (https://<environmentid>.live.dynatrace.com) if the policy should be applied to a specific environment",
		},
	}
}

func (me *Policy) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"name":            me.Name,
		"description":     me.Description,
		"statement_query": me.StatementQuery,
		"account":         me.Account,
		"environment":     me.Environment,
		"tags":            me.Tags,
	})
}

func (me *Policy) UnmarshalHCL(decoder hcl.Decoder) error {
	if err := decoder.DecodeAll(map[string]any{
		"name":            &me.Name,
		"description":     &me.Description,
		"statement_query": &me.StatementQuery,
		"account":         &me.Account,
		"environment":     &me.Environment,
		"tags":            &me.Tags,
	}); err != nil {
		return err
	}
	if me.Tags == nil {
		me.Tags = []string{}
	}
	return nil
}
