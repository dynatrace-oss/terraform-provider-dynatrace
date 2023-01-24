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

package kubernetes

// EndpointStatus The status of the configured endpoint.
// ASSIGNED: The credentials are assigned to an ActiveGate which is responsible for processing.
// UNASSIGNED: The credentials are not yet assigned to an ActiveGate so there is currently no processing.
// DISABLED: The credentials have been disabled by the user.
// FASTCHECK_AUTH_ERROR: The credentials are invalid.
// FASTCHECK_TLS_ERROR: The endpoint TLS certificate is invalid.
// FASTCHECK_NO_RESPONSE: The endpoint did not return a result until the timeout was reached.
// FASTCHECK_INVALID_ENDPOINT: The endpoint URL was invalid.
// FASTCHECK_AUTH_LOCKED: The credentials seem to be locked.
// UNKNOWN: An unknown error occured.
type EndpointStatus string

// EndpointStatuss offers the known enum values
var EndpointStates = struct {
	Assigned                 EndpointStatus
	Disabled                 EndpointStatus
	FastcheckAuthError       EndpointStatus
	FastcheckAuthLocked      EndpointStatus
	FastcheckInvalidEndpoint EndpointStatus
	FastcheckLowMemoryError  EndpointStatus
	FastcheckNoResponse      EndpointStatus
	FastcheckTlsError        EndpointStatus
	Unassigned               EndpointStatus
	Unknown                  EndpointStatus
}{
	"ASSIGNED",
	"DISABLED",
	"FASTCHECK_AUTH_ERROR",
	"FASTCHECK_AUTH_LOCKED",
	"FASTCHECK_INVALID_ENDPOINT",
	"FASTCHECK_LOW_MEMORY_ERROR",
	"FASTCHECK_NO_RESPONSE",
	"FASTCHECK_TLS_ERROR",
	"UNASSIGNED",
	"UNKNOWN",
}
