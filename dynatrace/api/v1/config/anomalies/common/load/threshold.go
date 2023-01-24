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

package load

// Threshold Minimal service load to detect response time degradation.
//  Response time degradation of services with smaller load won't trigger alerts.
type Threshold string

// Thresholds offers the known enum values
var Thresholds = struct {
	FifteenRequestsPerMinute Threshold
	FiveRequestsPerMinute    Threshold
	OneRequestPerMinute      Threshold
	TenRequestsPerMinute     Threshold
}{
	"FIFTEEN_REQUESTS_PER_MINUTE",
	"FIVE_REQUESTS_PER_MINUTE",
	"ONE_REQUEST_PER_MINUTE",
	"TEN_REQUESTS_PER_MINUTE",
}
