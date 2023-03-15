package permissions

import (
	"fmt"
)

type PermissionDTO struct {
	Name      string `json:"permissionName"`
	GroupID   string `json:"-"`
	Scope     string `json:"scope"`
	ScopeType string `json:"scopeType"`
}

func (me *PermissionDTO) ToID(groupID string) string {
	return fmt.Sprintf("%s#-#%s#-#%s#-#%s", groupID, me.Name, me.Scope, me.ScopeType)
}
