package quota

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// UserSessions represents user sessions consumption and quota information on environment level. If skipped when editing via PUT method then already set quotas will remain
type UserSessions struct {
	TotalAnnualLimit  *int64 `json:"totalAnnualLimit"`  // Annual total User sessions environment quota. Not set if unlimited. When updating via PUT method, skipping this field will set quota unlimited
	TotalMonthlyLimit *int64 `json:"totalMonthlyLimit"` // Monthly total User sessions environment quota. Not set if unlimited. When updating via PUT method, skipping this field will set quota unlimited
}

func (me *UserSessions) IsEmpty() bool {
	return me == nil || (me.TotalAnnualLimit == nil && me.TotalMonthlyLimit == nil)
}

func (me *UserSessions) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"annual": {
			Type:        schema.TypeInt,
			Optional:    true,
			Description: "Annual total User sessions environment quota. Not set if unlimited",
		},
		"monthly": {
			Type:        schema.TypeInt,
			Optional:    true,
			Description: "Monthly total User sessions environment quota. Not set if unlimited",
		},
	}
}

func (me *UserSessions) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"monthly": me.TotalMonthlyLimit,
		"annual":  me.TotalAnnualLimit,
	})
}

func (me *UserSessions) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"monthly": &me.TotalMonthlyLimit,
		"annual":  &me.TotalAnnualLimit,
	})
}
