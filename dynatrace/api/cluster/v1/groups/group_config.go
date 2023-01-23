package groups

import (
	"sort"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// GroupConfig represents the configuration of the group
type GroupConfig struct {
	ID                  *string             `json:"id"`                       // Group ID. Leave empty if creating group. Set if updating group
	Name                string              `json:"name"`                     // Group name
	IsClusterAdminGroup bool                `json:"isClusterAdminGroup"`      // If true, then the group has the cluster administrator rights
	LDAPGroupNames      []string            `json:"ldapGroupNames,omitempty"` // LDAP group names
	SSOGroupNames       []string            `json:"ssoGroupNames,omitempty"`  // SSO group names. If defined it's used to map SSO group name to Dynatrace group name, otherwise mapping is done by group name
	AccessRight         map[string][]string `json:"accessRight,omitempty"`    // Access rights
	AccessAccount       bool                `json:"accessAccount,omitempty"`  // write-only - no documentation available
	ManageAccount       bool                `json:"manageAccount,omitempty"`  // write-only - no documentation available
	// HasAccessAccountRole                    bool                `json:"hasAccessAccountRole,omitempty"`                    // If true, then the group has the access account rights
	// HasManageAccountAndViewProductUsageRole bool                `json:"hasManageAccountAndViewProductUsageRole,omitempty"` // If true, then the group has the manage account rights
}

func (me *GroupConfig) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The name of the user group",
		},
		"cluster_admin": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "If `true`, then the group has the cluster administrator rights",
		},
		"access_account": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "If `true`, then the group has the access account rights",
		},
		"manage_account": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "If `true`, then the group has the manage account rights",
		},
		"ldap_groups": {
			Type:        schema.TypeSet,
			Elem:        &schema.Schema{Type: schema.TypeString},
			Optional:    true,
			Description: "LDAP group names",
		},
		"sso_groups": {
			Type:        schema.TypeSet,
			Elem:        &schema.Schema{Type: schema.TypeString},
			Optional:    true,
			Description: "SSO group names. If defined it's used to map SSO group name to Dynatrace group name, otherwise mapping is done by group name",
		},
		"permissions": {
			Type:        schema.TypeList,
			Description: "Permissions for environments",
			Optional:    true,
			MaxItems:    1,
			Elem:        &schema.Resource{Schema: new(PermissionAssignments).Schema()},
		},
	}
}

func (me *GroupConfig) MarshalHCL(properties hcl.Properties) error {
	if err := properties.EncodeAll(map[string]any{
		"name":           me.Name,
		"cluster_admin":  me.IsClusterAdminGroup,
		"access_account": me.AccessAccount,
		"manage_account": me.ManageAccount,
		"ldap_groups":    me.LDAPGroupNames,
		"sso_groups":     me.SSOGroupNames,
	}); err != nil {
		return err
	}
	if len(me.AccessRight) > 0 {
		perms := PermissionAssignments{}
		for k, v := range me.AccessRight {
			sort.Strings(v)
			perms = append(perms, &PermissionAssignment{
				Permission:   Permission(k),
				Environments: v,
			})
		}
		sort.SliceStable(perms, func(i, j int) bool {
			return string(perms[i].Permission) < string(perms[j].Permission)
		})
		properties.Encode("permissions", perms)
	}
	return nil
}

func (me *GroupConfig) UnmarshalHCL(decoder hcl.Decoder) error {
	err := decoder.DecodeAll(map[string]any{
		"name":           &me.Name,
		"cluster_admin":  &me.IsClusterAdminGroup,
		"access_account": &me.AccessAccount,
		"manage_account": &me.ManageAccount,
		"ldap_groups":    &me.LDAPGroupNames,
		"sso_groups":     &me.SSOGroupNames,
	})
	if err != nil {
		return err
	}
	perms := PermissionAssignments{}
	err = decoder.Decode("permissions", &perms)
	if err != nil {
		return err
	}
	if len(perms) > 0 {
		me.AccessRight = map[string][]string{}
		for _, perm := range perms {
			me.AccessRight[string(perm.Permission)] = perm.Environments
		}
	}
	return nil
}

func (me *GroupConfig) Anonymize() *GroupConfig {
	return &GroupConfig{
		ID:                  nil,
		Name:                me.Name,
		IsClusterAdminGroup: me.IsClusterAdminGroup,
		AccessAccount:       me.AccessAccount,
		ManageAccount:       me.ManageAccount,
		LDAPGroupNames:      me.LDAPGroupNames,
		SSOGroupNames:       me.SSOGroupNames,
		AccessRight:         me.AccessRight,
	}
}

func (me *GroupConfig) Slim() *GroupConfig {
	result := GroupConfig{
		ID:                  me.ID,
		Name:                me.Name,
		IsClusterAdminGroup: me.IsClusterAdminGroup,
		AccessAccount:       me.AccessAccount,
		ManageAccount:       me.ManageAccount,
		LDAPGroupNames:      strSlice(me.LDAPGroupNames).Nil(),
		SSOGroupNames:       strSlice(me.SSOGroupNames).Nil(),
		AccessRight:         strSliceMap(me.AccessRight).Nil(),
	}
	return &result
}

func (me *GroupConfig) Equals(other *GroupConfig) bool {
	if other == nil {
		return false
	}
	if other.ID != me.ID {
		return false
	}
	if other.Name != me.Name {
		return false
	}
	if other.IsClusterAdminGroup != me.IsClusterAdminGroup {
		return false
	}
	if other.AccessAccount != me.AccessAccount {
		return false
	}
	if other.ManageAccount != me.ManageAccount {
		return false
	}
	if !EqualStringSlice(other.LDAPGroupNames, me.LDAPGroupNames) {
		return false
	}
	if !EqualStringSlice(other.SSOGroupNames, me.SSOGroupNames) {
		return false
	}
	if !EqualStringSliceMap(other.AccessRight, me.AccessRight) {
		return false
	}
	return true
}

func EqualStringSliceMap(a map[string][]string, b map[string][]string) bool {
	if a == nil && b != nil {
		return false
	}
	if a != nil && b == nil {
		return false
	}
	if a == nil && b == nil {
		return true
	}
	if len(a) != len(b) {
		return false
	}

	for k, va := range a {
		vb := b[k]
		if !EqualStringSlice(va, vb) {
			return false
		}
	}

	return true
}

func EqualStringSlice(a []string, b []string) bool {
	if a == nil && b != nil {
		return false
	}
	if a != nil && b == nil {
		return false
	}
	if a == nil && b == nil {
		return true
	}
	if len(a) != len(b) {
		return false
	}
	aLen := len(a)
	bLen := len(b)

	visited := make([]bool, bLen)

	for i := 0; i < aLen; i++ {
		found := false
		element := a[i]
		for j := 0; j < bLen; j++ {
			if visited[j] {
				continue
			}
			if element == b[j] {
				visited[j] = true
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}
