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

package pipelinegroups

type StageType string

var StageTypes = struct {
	Costallocation           StageType
	Dataextraction           StageType
	Davis                    StageType
	Metricextraction         StageType
	Processing               StageType
	Productallocation        StageType
	Securitycontext          StageType
	Smartscapeedgeextraction StageType
	Smartscapenodeextraction StageType
	Storage                  StageType
}{
	"costAllocation",
	"dataExtraction",
	"davis",
	"metricExtraction",
	"processing",
	"productAllocation",
	"securityContext",
	"smartscapeEdgeExtraction",
	"smartscapeNodeExtraction",
	"storage",
}
