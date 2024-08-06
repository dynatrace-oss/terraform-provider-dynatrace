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

package settings20

import (
	"encoding/json"
)

type SettingsObject struct {
	SchemaVersion string          `json:"schemaVersion"`
	SchemaID      string          `json:"schemaId"`
	Scope         string          `json:"scope"`
	Value         json.RawMessage `json:"value"`
}

type SettingsObjectList struct {
	Items       []*SettingsObjectListItem `json:"items"`
	NextPageKey *string                   `json:"nextPageKey,omitempty"`
}

type SettingsObjectListItem struct {
	ObjectID      string          `json:"objectId"`
	Scope         string          `json:"scope"`
	SchemaVersion string          `json:"schemaVersion"`
	SchemaID      string          `json:"schemaId"`
	Value         json.RawMessage `json:"value"`
}
