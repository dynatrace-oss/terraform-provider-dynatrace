/**
* @license
* Copyright 2025 Dynatrace LLC
*
* Licensed under the Apache License, Version 2.0 (the "License");
* you may not use this file except in compliance with the License.
* You may obtain a copy of the License at
*
*     http://www.apache.org/licenses/LICENSE-2.0
*
* Unless required by applicable law or agreed to in writing, software
* distributed under the License is distributed on an "AS IS" BASIS,
* WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
* See the License for the specific language governing permissions and
* limitations under the License.
 */

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
