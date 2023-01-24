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

type ModelType string

var ModelTypes = struct {
	Static       ModelType
	AutoAdaptive ModelType
	Seasonal     ModelType
}{
	"STATIC_THRESHOLD",
	"AUTO_ADAPTIVE_THRESHOLD",
	"SEASONAL_BASELINE",
}

type AlertCondition string

var AlertConditions = struct {
	Above   AlertCondition
	Below   AlertCondition
	Outside AlertCondition
}{
	"ABOVE",
	"BELOW",
	"OUTSIDE",
}

type Type string

var Types = struct {
	MetricKey      Type
	MetricSelector Type
}{
	"METRIC_KEY",
	"METRIC_SELECTOR",
}

type Aggregation string

var Aggregations = struct {
	Minimum    Aggregation
	Maximum    Aggregation
	Summary    Aggregation
	Count      Aggregation
	Average    Aggregation
	Value      Aggregation
	Median     Aggregation
	Percentile Aggregation
}{
	"MIN",
	"MAX",
	"SUM",
	"COUNT",
	"AVG",
	"VALUE",
	"MEDIAN",
	"PERCENTILE90",
}

type EntityFilterType string

var EntityFilterTypes = struct {
	Name             EntityFilterType
	Entity           EntityFilterType
	MgmtZone         EntityFilterType
	Tag              EntityFilterType
	Host             EntityFilterType
	HostGrpName      EntityFilterType
	ProcGrpName      EntityFilterType
	ProcGrpID        EntityFilterType
	CustomDevGrpName EntityFilterType
}{
	"NAME",
	"ENTITY_ID",
	"MANAGEMENT_ZONE",
	"TAG",
	"HOST_NAME",
	"HOST_GROUP_NAME",
	"PROCESS_GROUP_NAME",
	"PROCESS_GROUP_ID",
	"CUSTOM_DEVICE_GROUP_NAME",
}

type EntityFilterOperator string

var EntityFilterOperators = struct {
	ContainsInsensitive EntityFilterOperator
	ContainsSensitive   EntityFilterOperator
	Equals              EntityFilterOperator
}{
	"CONTAINS_CASE_INSENSITIVE",
	"CONTAINS_CASE_SENSITIVE",
	"EQUALS",
}

type EventTypeEnum string

var EventTypeEnums = struct {
	Info              EventTypeEnum
	Error             EventTypeEnum
	Availability      EventTypeEnum
	Slowdown          EventTypeEnum
	Resource          EventTypeEnum
	CustomAlert       EventTypeEnum
	CustomAnnotation  EventTypeEnum
	CustomConfig      EventTypeEnum
	CustomDeployment  EventTypeEnum
	MaskedTermination EventTypeEnum
}{
	"INFO",
	"ERROR",
	"AVAILABILITY",
	"SLOWDOWN",
	"RESOURCE",
	"CUSTOM_ALERT",
	"CUSTOM_ANNOTATION",
	"CUSTOM_CONFIGURATION",
	"CUSTOM_DEPLOYMENT",
	"MARKED_FOR_TERMINATION",
}
