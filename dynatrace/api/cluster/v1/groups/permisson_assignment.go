package groups

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type PermissionAssignments []*PermissionAssignment

func (me *PermissionAssignments) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"grant": {
			Type:        schema.TypeList,
			Description: "A permission granted to one or multiple environments",
			Required:    true,
			MinItems:    1,
			Elem:        &schema.Resource{Schema: new(PermissionAssignment).Schema()},
		},
	}
}

func (me PermissionAssignments) MarshalHCL(properties hcl.Properties) error {
	if len(me) > 0 {
		if err := properties.EncodeSlice("grant", me); err != nil {
			return err
		}
	}
	return nil
}

func (me *PermissionAssignments) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeSlice("grant", me)
}

type PermissionAssignment struct {
	Permission   Permission
	Environments []string
}

func (me *PermissionAssignment) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"permission": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The permission. Possible values are `VIEWER`, `MANAGE_SETTINGS`, `AGENT_INSTALL`, `LOG_VIEWER`, `VIEW_SENSITIVE_REQUEST_DATA`, `CONFIGURE_REQUEST_CAPTURE_DATA`, `REPLAY_SESSION_DATA`, `REPLAY_SESSION_DATA_WITHOUT_MASKING`, `MANAGE_SECURITY_PROBLEMS` and `MANAGE_SUPPORT_TICKETS`.",
		},
		"environments": {
			Type:        schema.TypeSet,
			Elem:        &schema.Schema{Type: schema.TypeString},
			Optional:    true,
			Description: "The ids of the environments this permission grants the user access to.",
		},
	}
}

func (me *PermissionAssignment) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"permission":   me.Permission,
		"environments": me.Environments,
	})
}

func (me *PermissionAssignment) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"permission":   &me.Permission,
		"environments": &me.Environments,
	})
}

type Permission string

var Permissions = struct {
	Viewer                          Permission
	ManageSettings                  Permission
	AgentInstall                    Permission
	LogViewer                       Permission
	ViewSensitiveRequestData        Permission
	ConfigureRequestCaptureData     Permission
	ReplaySessionData               Permission
	ReplaySessionDataWithoutMasking Permission
	ManageSecurityProblems          Permission
	ManageSupportTickets            Permission
}{
	"VIEWER",
	"MANAGE_SETTINGS",
	"AGENT_INSTALL",
	"LOG_VIEWER",
	"VIEW_SENSITIVE_REQUEST_DATA",
	"CONFIGURE_REQUEST_CAPTURE_DATA",
	"REPLAY_SESSION_DATA",
	"REPLAY_SESSION_DATA_WITHOUT_MASKING",
	"MANAGE_SECURITY_PROBLEMS",
	"MANAGE_SUPPORT_TICKETS",
}
