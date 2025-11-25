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

package storage

// LogMonitoringStorage represents log monitoring storage usage and limit information on environment level. Not editable when Log monitoring is not allowed by license or not configured on cluster level. If skipped when editing via PUT method then already set limit will remain
type LogMonitoringStorage struct {
	MaxLimit *int64 `json:"maxLimit"` // Maximum storage limit [bytes]
}

func (me *LogMonitoringStorage) IsEmpty() bool {
	return me == nil || me.MaxLimit == nil || *me.MaxLimit == 0
}
