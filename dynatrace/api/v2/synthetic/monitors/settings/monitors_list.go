/**
* @license
* Copyright 2020 Dynatrace LLC
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

package monitors

import "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api"

type MonitorsList struct {
	Monitors []MonitorsResponseElement `json:"monitors"` // A list of monitors
}

func (me MonitorsList) ToStubs() api.Stubs {
	stubs := api.Stubs{}
	for _, elem := range me.Monitors {
		stubs = append(stubs, &api.Stub{ID: elem.EntityId, Name: elem.Name})
	}
	return stubs
}

type MonitorsResponseElement struct {
	EntityId string      `json:"entityId"` // Entity ID of the monitor
	Name     string      `json:"name"`     // Name of the monitor
	Type     MonitorType `json:"type"`     // Type of the monitor
	Enabled  bool        `json:"enabled"`  // If true, the monitor is enabled
}
