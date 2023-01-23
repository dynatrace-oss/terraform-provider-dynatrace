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

package autotags

// PropagationType has no documentation
type PropagationType string

// PropagationTypes offers the known enum values
var PropagationTypes = struct {
	AzureToPg                  PropagationType
	AzureToService             PropagationType
	HostToProcessGroupInstance PropagationType
	ProcessGroupToHost         PropagationType
	ProcessGroupToService      PropagationType
	ServiceToHostLike          PropagationType
	ServiceToProcessGroupLike  PropagationType
}{
	"AZURE_TO_PG",
	"AZURE_TO_SERVICE",
	"HOST_TO_PROCESS_GROUP_INSTANCE",
	"PROCESS_GROUP_TO_HOST",
	"PROCESS_GROUP_TO_SERVICE",
	"SERVICE_TO_HOST_LIKE",
	"SERVICE_TO_PROCESS_GROUP_LIKE",
}
