package mgmz

import (
	"fmt"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

type Permission struct {
	EnvironmentID    string
	GroupID          string
	ManagementZoneID string
	Permissions      []string
}

func (me *Permission) Name() string {
	return fmt.Sprintf("%s#-#%s#-#%s", me.EnvironmentID, me.GroupID, me.ManagementZoneID)
}

func (me *Permission) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"environment": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The UUID of the environment",
		},
		"group": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The ID of the group the permissions are valid for. You may refer to the id of a resource `dynatrace_user_group` here",
		},
		"management_zone": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The ID of the management zone the permissions are valid for. When referring to resource `dynatrace_management_zone_v2` or data source `dynatrace_management_zone` you need to refer to the attribute `legacy_id`.",
		},
		"permissions": {
			Type:        schema.TypeSet,
			Required:    true,
			MinItems:    1,
			Description: "The permissions to assign for that management zone. Allowed values are `DEMO_USER`, `LOG_VIEWER`, `MANAGE_SECURITY_PROBLEMS`, `MANAGE_SETTINGS`, `REPLAY_SESSION_DATA`, `REPLAY_SESSION_DATA_WITHOUT_MASKING`, `VIEWER`, `VIEW_SENSITIVE_REQUEST_DATA`",
			Elem: &schema.Schema{
				Type:         schema.TypeString,
				ValidateFunc: validation.StringInSlice([]string{"DEMO_USER", "LOG_VIEWER", "MANAGE_SECURITY_PROBLEMS", "MANAGE_SETTINGS", "REPLAY_SESSION_DATA", "REPLAY_SESSION_DATA_WITHOUT_MASKING", "VIEWER", "VIEW_SENSITIVE_REQUEST_DATA"}, false),
			},
		},
	}
}

func (me *Permission) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"environment":     me.EnvironmentID,
		"group":           me.GroupID,
		"management_zone": me.ManagementZoneID,
		"permissions":     me.Permissions,
	})
}

func (me *Permission) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"environment":     &me.EnvironmentID,
		"group":           &me.GroupID,
		"management_zone": &me.ManagementZoneID,
		"permissions":     &me.Permissions,
	})
}
