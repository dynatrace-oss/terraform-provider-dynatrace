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

package relation

type IngestDataSource string

var IngestDataSources = struct {
	Entities IngestDataSource
	Events   IngestDataSource
	Logs     IngestDataSource
	Metrics  IngestDataSource
	Spans    IngestDataSource
	Topology IngestDataSource
}{
	"Entities",
	"Events",
	"Logs",
	"Metrics",
	"Spans",
	"Topology",
}

type Normalization string

var Normalizations = struct {
	LeavetextasIs Normalization
	Tolowercase   Normalization
	Touppercase   Normalization
}{
	"Leavetextas_is",
	"Tolowercase",
	"Touppercase",
}

type RelationType string

var RelationTypes = struct {
	Calls      RelationType
	ChildOf    RelationType
	InstanceOf RelationType
	PartOf     RelationType
	RunsOn     RelationType
	SameAs     RelationType
}{
	"CALLS",
	"CHILD_OF",
	"INSTANCE_OF",
	"PART_OF",
	"RUNS_ON",
	"SAME_AS",
}
