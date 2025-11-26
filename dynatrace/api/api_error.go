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

package api

// APIError wraps lower-level errors with an API domain context.
type APIError struct {
	Method string
	Err    error
}

func (e APIError) Error() string {
	if e.Err == nil {
		return e.Method
	}

	return "API Error: <" + e.Method + "> " + e.Err.Error()
}

func (e APIError) Unwrap() error {
	return e.Err
}

func NewAPIError(m string, err error) APIError {
	return APIError{Method: m, Err: err}
}
