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

package quota

// HostUnits represents host units consumption and quota information on environment level. If skipped when editing via PUT method then already set quota will remain
type HostUnits struct {
	MaxLimit *int64 `json:"maxLimit"` // Concurrent environment quota. Not set if unlimited. When updating via PUT method, skipping this field will set quota unlimited
}

func (me *HostUnits) IsEmpty() bool {
	return me == nil || me.MaxLimit == nil
}
