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

package diskrules

type DiskMetric string

var DiskMetrics = struct {
	LowDiskSpace       DiskMetric
	LowInodes          DiskMetric
	ReadTimeExceeding  DiskMetric
	WriteTimeExceeding DiskMetric
}{
	"LOW_DISK_SPACE",
	"LOW_INODES",
	"READ_TIME_EXCEEDING",
	"WRITE_TIME_EXCEEDING",
}

type DiskNameFilterOperator string

var DiskNameFilterOperators = struct {
	Contains         DiskNameFilterOperator
	DoesNotContain   DiskNameFilterOperator
	DoesNotEqual     DiskNameFilterOperator
	DoesNotStartWith DiskNameFilterOperator
	Equals           DiskNameFilterOperator
	StartsWith       DiskNameFilterOperator
}{
	"CONTAINS",
	"DOES_NOT_CONTAIN",
	"DOES_NOT_EQUAL",
	"DOES_NOT_START_WITH",
	"EQUALS",
	"STARTS_WITH",
}
