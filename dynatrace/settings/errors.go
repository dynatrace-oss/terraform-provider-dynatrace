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

package settings

import "fmt"

// ErrValidatorNotImplemented reports that the service for the given schema was expected to implement
// Validator but does not. This is an internal invariant violation (a provider bug) rather than a
// user configuration error, so it directs the reader to reach out to support. The schema ID is
// included to pinpoint which service is affected without resorting to a stack trace.
func ErrValidatorNotImplemented(schemaID string) error {
	return fmt.Errorf("internal provider bug: the service for schema %q was expected to support validation but does not. Please reach out to support", schemaID)
}
