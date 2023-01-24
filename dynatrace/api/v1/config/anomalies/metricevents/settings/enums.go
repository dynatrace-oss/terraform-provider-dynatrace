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

package metricevents

// AggregationType How the metric data points are aggregated for the evaluation.
//
//	The timeseries must support this aggregation.
type AggregationType string

func (me AggregationType) Ref() *AggregationType {
	return &me
}

// AggregationTypes offers the known enum values
var AggregationTypes = struct {
	Avg    AggregationType
	Count  AggregationType
	Max    AggregationType
	Median AggregationType
	Min    AggregationType
	P90    AggregationType
	Sum    AggregationType
	Value  AggregationType
}{
	"AVG",
	"COUNT",
	"MAX",
	"MEDIAN",
	"MIN",
	"P90",
	"SUM",
	"VALUE",
}

// WarningReason The reason of a warning set on the config.
// The `NONE` means config has no warnings.
type WarningReason string

func (me WarningReason) Ref() *WarningReason {
	return &me
}

// WarningReasons offers the known enum values
var WarningReasons = struct {
	None        WarningReason
	TooManyDims WarningReason
}{
	"NONE",
	"TOO_MANY_DIMS",
}

// DisabledReason The reason of automatic disabling.
// The `NONE` means config was not disabled automatically.
type DisabledReason string

func (me DisabledReason) Ref() *DisabledReason {
	return &me
}

// DisabledReasons offers the known enum values
var DisabledReasons = struct {
	MetricDefinitionInconsistency DisabledReason
	None                          DisabledReason
	TooManyDims                   DisabledReason
	TopxForciblyDeactivated       DisabledReason
}{
	"METRIC_DEFINITION_INCONSISTENCY",
	"NONE",
	"TOO_MANY_DIMS",
	"TOPX_FORCIBLY_DEACTIVATED",
}

// Severity The type of the event to trigger on the threshold violation.
// The `CUSTOM_ALERT` type is not correlated with other alerts.
// The `INFO` type does not open a problem.
type Severity string

func (me Severity) Ref() *Severity {
	return &me
}

// Severitys offers the known enum values
var Severitys = struct {
	Availability       Severity
	CustomAlert        Severity
	Error              Severity
	Info               Severity
	Performance        Severity
	ResourceContention Severity
}{
	"AVAILABILITY",
	"CUSTOM_ALERT",
	"ERROR",
	"INFO",
	"PERFORMANCE",
	"RESOURCE_CONTENTION",
}
