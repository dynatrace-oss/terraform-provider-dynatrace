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

type Aggregation string

var Aggregations = struct {
	Avg          Aggregation
	Count        Aggregation
	Max          Aggregation
	Median       Aggregation
	Min          Aggregation
	Percentile90 Aggregation
	Sum          Aggregation
	Value        Aggregation
}{
	"AVG",
	"COUNT",
	"MAX",
	"MEDIAN",
	"MIN",
	"PERCENTILE90",
	"SUM",
	"VALUE",
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

type DimensionFilterOperator string

var DimensionFilterOperators = struct {
	ContainsCaseSensitive       DimensionFilterOperator
	DoesNotContainCaseSensitive DimensionFilterOperator
	DoesNotEqual                DimensionFilterOperator
	DoesNotStartWith            DimensionFilterOperator
	Equals                      DimensionFilterOperator
	StartsWith                  DimensionFilterOperator
}{
	"CONTAINS_CASE_SENSITIVE",
	"DOES_NOT_CONTAIN_CASE_SENSITIVE",
	"DOES_NOT_EQUAL",
	"DOES_NOT_START_WITH",
	"EQUALS",
	"STARTS_WITH",
}

type EntityFilterOperator string

var EntityFilterOperators = struct {
	ContainsCaseInsensitive       EntityFilterOperator
	ContainsCaseSensitive         EntityFilterOperator
	DoesNotContainCaseInsensitive EntityFilterOperator
	DoesNotContainCaseSensitive   EntityFilterOperator
	DoesNotEqual                  EntityFilterOperator
	DoesNotStartWith              EntityFilterOperator
	Equals                        EntityFilterOperator
	StartsWith                    EntityFilterOperator
}{
	"CONTAINS_CASE_INSENSITIVE",
	"CONTAINS_CASE_SENSITIVE",
	"DOES_NOT_CONTAIN_CASE_INSENSITIVE",
	"DOES_NOT_CONTAIN_CASE_SENSITIVE",
	"DOES_NOT_EQUAL",
	"DOES_NOT_START_WITH",
	"EQUALS",
	"STARTS_WITH",
}

type EntityFilterType string

var EntityFilterTypes = struct {
	CustomDeviceGroupName EntityFilterType
	EntityId              EntityFilterType
	HostGroupName         EntityFilterType
	HostName              EntityFilterType
	ManagementZone        EntityFilterType
	Name                  EntityFilterType
	ProcessGroupId        EntityFilterType
	ProcessGroupName      EntityFilterType
	Tag                   EntityFilterType
}{
	"CUSTOM_DEVICE_GROUP_NAME",
	"ENTITY_ID",
	"HOST_GROUP_NAME",
	"HOST_NAME",
	"MANAGEMENT_ZONE",
	"NAME",
	"PROCESS_GROUP_ID",
	"PROCESS_GROUP_NAME",
	"TAG",
}

type EventTypeEnum string

var EventTypeEnums = struct {
	Availability         EventTypeEnum
	CustomAlert          EventTypeEnum
	CustomAnnotation     EventTypeEnum
	CustomConfiguration  EventTypeEnum
	CustomDeployment     EventTypeEnum
	Error                EventTypeEnum
	Info                 EventTypeEnum
	MarkedForTermination EventTypeEnum
	Resource             EventTypeEnum
	Slowdown             EventTypeEnum
}{
	"AVAILABILITY",
	"CUSTOM_ALERT",
	"CUSTOM_ANNOTATION",
	"CUSTOM_CONFIGURATION",
	"CUSTOM_DEPLOYMENT",
	"ERROR",
	"INFO",
	"MARKED_FOR_TERMINATION",
	"RESOURCE",
	"SLOWDOWN",
}

type ModelType string

var ModelTypes = struct {
	AutoAdaptiveThreshold ModelType
	SeasonalBaseline      ModelType
	StaticThreshold       ModelType
}{
	"AUTO_ADAPTIVE_THRESHOLD",
	"SEASONAL_BASELINE",
	"STATIC_THRESHOLD",
}

type Type string

var Types = struct {
	MetricKey      Type
	MetricSelector Type
}{
	"METRIC_KEY",
	"METRIC_SELECTOR",
}
