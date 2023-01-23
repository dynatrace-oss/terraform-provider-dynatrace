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

package notifications

// Stub The short representation of a notification.
type Stub struct {
	Description *string `json:"description,omitempty"` // A short description of the Dynatrace entity.
	ID          string  `json:"id"`                    // The ID of the Dynatrace entity.
	Name        *string `json:"name,omitempty"`        // The name of the Dynatrace entity.
	Type        *Type   `json:"type,omitempty"`        // The type of the notification.
}

// StubList has no documentation
type StubList struct {
	Values []Stub `json:"values,omitempty"` // has no documentation
}
