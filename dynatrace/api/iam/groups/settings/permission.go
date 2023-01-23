package groups

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Permissions []*Permission

func (me *Permissions) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"permission": {
			Type:        schema.TypeSet,
			Required:    true,
			MinItems:    1,
			Description: "A Permission",
			Elem:        &schema.Resource{Schema: new(Permission).Schema()},
		},
	}
}

func (me Permissions) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeSlice("permission", me)
}

func (me *Permissions) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeSlice("permission", me)
}

type Permission struct {
	Name      string `json:"permissionName"` // The name of the permission. Possible values are account-company-info, account-user-management, account-viewer, tenant-viewer, tenant-manage-settings, tenant-agent-install, tenant-logviewer, tenant-view-sensitive-request-data, tenant-configure-request-capture-data, tenant-replay-sessions-with-masking, tenant-replay-sessions-without-masking, tenant-manage-security-problems, tenant-manage-support-tickets.
	Scope     string `json:"scope"`          // The scope of the permission. Depending on the scope type, it is defined by the UUID of the account (scopeType = `account`), the ID of the environment (scopeType = `tenant`) or the ID of the management zone from an environment in `{environment-id}-{management-zone-id}` format (scopeType = `management-zone`)
	ScopeType string `json:"scopeType"`      // The type of the permission scope. Possible values are `account`, `tenant` and `management-zone`
}

func (me *Permission) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:     schema.TypeString,
			Required: true,
		},
		"scope": {
			Type:     schema.TypeString,
			Required: true,
		},
		"type": {
			Type:     schema.TypeString,
			Required: true,
		},
	}
}

func (me *Permission) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"name":  me.Name,
		"scope": me.Scope,
		"type":  string(me.ScopeType),
	})
}

func (me *Permission) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"name":  &me.Name,
		"scope": &me.Scope,
		"type":  &me.ScopeType,
	})
}

type ScopeType string

var ScopeTypes = struct {
	Account        ScopeType
	Tenant         ScopeType
	ManagementZone ScopeType
}{
	ScopeType("account"),
	ScopeType("tenant"),
	ScopeType("management-zone"),
}
