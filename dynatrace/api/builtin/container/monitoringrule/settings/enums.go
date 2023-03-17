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

package monitoringrule

type ConditionOperator string

var ConditionOperators = struct {
	Contains    ConditionOperator
	Ends        ConditionOperator
	Equals      ConditionOperator
	Exists      ConditionOperator
	NotContains ConditionOperator
	NotEnds     ConditionOperator
	NotEquals   ConditionOperator
	NotExists   ConditionOperator
	NotStarts   ConditionOperator
	Starts      ConditionOperator
}{
	"CONTAINS",
	"ENDS",
	"EQUALS",
	"EXISTS",
	"NOT_CONTAINS",
	"NOT_ENDS",
	"NOT_EQUALS",
	"NOT_EXISTS",
	"NOT_STARTS",
	"STARTS",
}

type ContainerItem string

var ContainerItems = struct {
	ContainerName           ContainerItem
	ImageName               ContainerItem
	KubernetesBasepodname   ContainerItem
	KubernetesContainername ContainerItem
	KubernetesFullpodname   ContainerItem
	KubernetesNamespace     ContainerItem
	KubernetesPoduid        ContainerItem
}{
	"CONTAINER_NAME",
	"IMAGE_NAME",
	"KUBERNETES_BASEPODNAME",
	"KUBERNETES_CONTAINERNAME",
	"KUBERNETES_FULLPODNAME",
	"KUBERNETES_NAMESPACE",
	"KUBERNETES_PODUID",
}

type MonitoringMode string

var MonitoringModes = struct {
	MonitoringOff MonitoringMode
	MonitoringOn  MonitoringMode
}{
	"MONITORING_OFF",
	"MONITORING_ON",
}
