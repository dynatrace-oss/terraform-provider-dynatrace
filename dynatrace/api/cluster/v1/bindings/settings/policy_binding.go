package bindings

import (
	"fmt"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type PolicyBinding struct {
	GroupID     string   `json:"-"`
	Cluster     string   `json:"-"`
	Environment string   `json:"-"`
	PolicyIDs   []string `json:"policyUuids"`
}

func (me *PolicyBinding) Name() string {
	if len(me.Cluster) > 0 {
		return fmt.Sprintf("%s#-#%s#-#%s", me.GroupID, "cluster", me.Cluster)
	}
	if len(me.Environment) > 0 {
		return fmt.Sprintf("%s#-#%s#-#%s", me.GroupID, "environment", me.Environment)
	}
	return fmt.Sprintf("%s#-#%s#-#%s", me.GroupID, "global", "global")
}

func (me *PolicyBinding) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"group": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The name of the policy",
			ForceNew:    true,
		},
		"cluster": {
			Type:          schema.TypeString,
			Optional:      true,
			ConflictsWith: []string{"environment"},
			AtLeastOneOf:  []string{"environment", "cluster"},
			ForceNew:      true,
			Description:   "The UUID of the cluster. The attribute `policies` must contain ONLY policies defined for that cluster.",
		},
		"environment": {
			Type:          schema.TypeString,
			Optional:      true,
			ConflictsWith: []string{"cluster"},
			AtLeastOneOf:  []string{"environment", "cluster"},
			ForceNew:      true,
			Description:   "The ID of the environment (https://<environmentid>.live.dynatrace.com). The attribute `policies` must contain ONLY policies defined for that environment.",
		},
		"policies": {
			Type:        schema.TypeSet,
			Required:    true,
			MinItems:    1,
			Description: "A list of IDs referring to policies bound to that group. It's not possible to mix policies here that are defined for different scopes (different clusters or environments) than specified via attributes `cluster` or `environment`.",
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
	}
}

func (me *PolicyBinding) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"group":       me.GroupID,
		"cluster":     me.Cluster,
		"environment": me.Environment,
		"policies":    me.PolicyIDs,
	})
}

func (me *PolicyBinding) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"group":       &me.GroupID,
		"cluster":     &me.Cluster,
		"environment": &me.Environment,
		"policies":    &me.PolicyIDs,
	})
}
