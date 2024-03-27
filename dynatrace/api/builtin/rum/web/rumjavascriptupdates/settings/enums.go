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

package rumjavascriptupdates

type JavascriptVersion string

var JavascriptVersions = struct {
	Custom               JavascriptVersion
	LatestIe11Supported  JavascriptVersion
	LatestIe710Supported JavascriptVersion
	LatestStable         JavascriptVersion
	PreviousStable       JavascriptVersion
}{
	"CUSTOM",
	"LATEST_IE11_SUPPORTED",
	"LATEST_IE7_10_SUPPORTED",
	"LATEST_STABLE",
	"PREVIOUS_STABLE",
}
