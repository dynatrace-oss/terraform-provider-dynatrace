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

package logstoragesettings

type MatcherType string

var MatcherTypes = struct {
	ContainerName          MatcherType
	DtEntityContainerGroup MatcherType
	DtEntityProcessGroup   MatcherType
	HostTag                MatcherType
	K8SContainerName       MatcherType
	K8SDeploymentName      MatcherType
	K8SNamespaceName       MatcherType
	K8SPodAnnotation       MatcherType
	K8SPodLabel            MatcherType
	K8SWorkloadKind        MatcherType
	K8SWorkloadName        MatcherType
	LogContent             MatcherType
	LogSource              MatcherType
	LogSourceOrigin        MatcherType
	Loglevel               MatcherType
	ProcessTechnology      MatcherType
	WinlogEventid          MatcherType
	WinlogKeywords         MatcherType
	WinlogOpcode           MatcherType
	WinlogProvider         MatcherType
	WinlogTask             MatcherType
	WinlogUsername         MatcherType
}{
	"container.name",
	"dt.entity.container_group",
	"dt.entity.process_group",
	"host.tag",
	"k8s.container.name",
	"k8s.deployment.name",
	"k8s.namespace.name",
	"k8s.pod.annotation",
	"k8s.pod.label",
	"k8s.workload.kind",
	"k8s.workload.name",
	"log.content",
	"log.source",
	"log.source.origin",
	"loglevel",
	"process.technology",
	"winlog.eventid",
	"winlog.keywords",
	"winlog.opcode",
	"winlog.provider",
	"winlog.task",
	"winlog.username",
}

type Operator string

var Operators = struct {
	Matches Operator
}{
	"MATCHES",
}
