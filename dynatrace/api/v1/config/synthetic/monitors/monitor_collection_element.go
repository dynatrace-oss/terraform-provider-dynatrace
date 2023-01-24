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

// MonitorCollectionElement is the short representation of a synthetic monitor
type MonitorCollectionElement struct {
	// "required" : [ "enabled", "entityId", "name", "type" ],
	Name     string `json:"name"`     // The name of a synthetic object
	EntityID string `json:"entityId"` // The ID of a synthetic object
	Type     Type   `json:"type"`     // The type of a synthetic monitor
	Enabled  bool   `json:"enabled"`  // The state of a synthetic monitor
}

// Type The type of a synthetic monitor
// BROWSER: A Browser Monitor
// HTTP: A HTTP Monitor
type Type string

// Types offers the known enum values
var Types = struct {
	Browser Type
	HTTP    Type
}{
	"BROWSER",
	"HTTP",
}
