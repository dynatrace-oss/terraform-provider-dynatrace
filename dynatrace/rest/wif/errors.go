/**
* @license
* Copyright 2026 Dynatrace LLC
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

package wif

// ConfigError marks a failure the user can only resolve by changing the provider configuration or
// the CI job, as opposed to a token service being temporarily unreachable. That distinction is what
// lets the provider report these while it is being configured rather than on the first request.
type ConfigError struct {
	Message string
}

func (err ConfigError) Error() string {
	return err.Message
}
