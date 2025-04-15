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

package objectivetemplates

type TemplatesResponse struct {
	TotalCount  int                 `json:"totalCount"`
	NextPageKey *string             `json:"nextPageKey,omitempty"`
	Items       []ObjectiveTemplate `json:"items"`
}

type ObjectiveTemplate struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description *string `json:"description,omitempty"`
	Indicator   string  `json:"indicator"`
	Variables   []struct {
		Name  string `json:"name"`
		Scope string `json:"scope"`
	} `json:"variables"`
	Tags            []string `json:"tags,omitempty"`
	ApplicableScope string   `json:"applicableScope"`
	Version         *string  `json:"version,omitempty"`
	BuiltIn         bool     `json:"builtIn"`
}
