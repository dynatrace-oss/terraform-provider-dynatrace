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

package monitors

type MonitorType string

var MonitorTypes = struct {
	MultiProtocol MonitorType
}{
	"MULTI_PROTOCOL",
}

type TagSource string

var TagSources = struct {
	Auto      TagSource
	RuleBased TagSource
	User      TagSource
}{
	"AUTO",
	"RULE_BASED",
	"USER",
}

type Aggregation string

var Aggregations = struct {
	Avg Aggregation
	Max Aggregation
	Min Aggregation
}{
	"AVG",
	"MAX",
	"MIN",
}

type RequestType string

var RequestTypes = struct {
	ICMP RequestType
	TCP  RequestType
	DNS  RequestType
}{
	"ICMP",
	"TCP",
	"DNS",
}
