package bindings

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type PolicyBinding struct {
	GroupID     string   `json:"-"`
	Account     string   `json:"-"`
	Environment string   `json:"-"`
	PolicyIDs   []string `json:"policyUuids"`
}

func (me *PolicyBinding) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"group": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The name of the policy",
			ForceNew:    true,
		},
		"account": {
			Type:          schema.TypeString,
			Optional:      true,
			ConflictsWith: []string{"environment"},
			AtLeastOneOf:  []string{"environment", "account"},
			ForceNew:      true,
			Description:   "The UUID of the account (`urn:dtaccount:<account-uuid>`). The attribute `policies` must contain ONLY policies defined for that account. The prefix `urn:dtaccount:` MUST be omitted here.",
		},
		"environment": {
			Type:          schema.TypeString,
			Optional:      true,
			ConflictsWith: []string{"account"},
			AtLeastOneOf:  []string{"environment", "account"},
			ForceNew:      true,
			Description:   "The ID of the environment (https://<environmentid>.live.dynatrace.com). The attribute `policies` must contain ONLY policies defined for that environment.",
		},
		"policies": {
			Type:        schema.TypeSet,
			Required:    true,
			MinItems:    1,
			Description: "A list of IDs referring to policies bound to that group. It's not possible to mix policies here that are defined for different scopes (different accounts or environments) than specified via attributes `account` or `environment`.",
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
	}
}

func (me *PolicyBinding) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"group":       me.GroupID,
		"account":     me.Account,
		"environment": me.Environment,
		"policies":    me.PolicyIDs,
	})
}

func (me *PolicyBinding) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"group":       &me.GroupID,
		"account":     &me.Account,
		"environment": &me.Environment,
		"policies":    &me.PolicyIDs,
	})
}
