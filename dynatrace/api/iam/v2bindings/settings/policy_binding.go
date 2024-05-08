package v2bindings

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type PolicyBinding struct {
	GroupID     string    `json:"-"`
	Account     string    `json:"-"`
	Environment string    `json:"-"`
	Policies    []*Policy `json:"policies"`
}

func (me *PolicyBinding) Name() string {
	return me.GroupID
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
		"policy": {
			Type:        schema.TypeSet,
			Description: "A list of Policies (ID and parameters/metadata) referring to policies bound to that group. It's not possible to mix policies here that are defined for different scopes (different accounts or environments) than specified via attributes `account` or `environment`.",
			Optional:    true,
			Elem:        &schema.Resource{Schema: new(Policy).Schema()},
		},
	}
}

func (me *PolicyBinding) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"group":       me.GroupID,
		"account":     me.Account,
		"environment": me.Environment,
		"policy":      me.Policies,
	})
}

func (me *PolicyBinding) UnmarshalHCL(decoder hcl.Decoder) error {
	if err := decoder.DecodeAll(map[string]any{
		"group":       &me.GroupID,
		"account":     &me.Account,
		"environment": &me.Environment,
	}); err != nil {
		return err
	}
	if err := decoder.DecodeSlice("policy", &me.Policies); err != nil {
		return err
	}
	// Terraform has a bug when it comes to TypeSet
	// Sometimes empty entries are getting produced
	if len(me.Policies) > 0 {
		finalPolicies := []*Policy{}
		for _, policy := range me.Policies {
			if policy == nil {
				continue
			}
			if len(policy.ID) == 0 {
				continue
			}
			finalPolicies = append(finalPolicies, policy)
		}
		me.Policies = finalPolicies
	}

	return nil
}
