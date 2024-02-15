package permissions

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

type Permission struct {
	Name           string
	GroupID        string
	Environment    string
	ManagementZone string
	Account        string
}

func (me *Permission) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:         schema.TypeString,
			Required:     true,
			Description:  "Possible values: `account-company-info`, `account-user-management`, `account-viewer`, `account-saml-flexible-federation`, `tenant-viewer`, `tenant-manage-settings`, `tenant-agent-install`, `tenant-logviewer`, `tenant-view-sensitive-request-data`, `tenant-configure-request-capture-data`, `tenant-replay-sessions-with-masking`, `tenant-replay-sessions-without-masking`, `tenant-manage-security-problems`, `tenant-view-security-problems`, `tenant-manage-support-tickets`",
			ValidateFunc: validation.StringInSlice([]string{`account-company-info`, `account-user-management`, `account-viewer`, `account-saml-flexible-federation`, `tenant-viewer`, `tenant-manage-settings`, `tenant-agent-install`, `tenant-logviewer`, `tenant-view-sensitive-request-data`, `tenant-configure-request-capture-data`, `tenant-replay-sessions-with-masking`, `tenant-replay-sessions-without-masking`, `tenant-manage-security-problems`, `tenant-view-security-problems`, `tenant-manage-support-tickets`}, false),
		},
		"group": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The ID of the group this permission is valid for",
			ForceNew:    true,
		},
		"environment": {
			Type:          schema.TypeString,
			Optional:      true,
			Description:   "The environment this permission is valid (`https://<environmentid>.live.dynatrace.com`). Also required in when trying to specify a management zone permission.",
			ExactlyOneOf:  []string{`account`, `environment`},
			ConflictsWith: []string{"account"},
			ForceNew:      true,
		},
		"management_zone": {
			Type:          schema.TypeString,
			Optional:      true,
			Description:   "The management zone this permission is valid for. You need to use the attribute `legacy_id` when referring to a resource `dynatrace_management_zone_v2` or a data source `dynatrace_management_zone`. The attribute `environment` is required to get specified also in order to identify the management zone uniquely.",
			ConflictsWith: []string{"account"},
			RequiredWith:  []string{"environment"},
			ForceNew:      true,
		},
		"account": {
			Type:          schema.TypeString,
			Optional:      true,
			Description:   "The UUID of the account this permission is valid for",
			ExactlyOneOf:  []string{`account`, `environment`},
			ConflictsWith: []string{"management_zone", "environment"},
			ForceNew:      true,
		},
	}
}

func (me *Permission) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"name":            me.Name,
		"group":           me.GroupID,
		"environment":     me.Environment,
		"account":         me.Account,
		"management_zone": me.ManagementZone,
	})
}

func (me *Permission) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"name":            &me.Name,
		"group":           &me.GroupID,
		"environment":     &me.Environment,
		"account":         &me.Account,
		"management_zone": &me.ManagementZone,
	})
}
