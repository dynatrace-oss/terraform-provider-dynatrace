/**
* @license
* Copyright 2026 Dynatrace LLC
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

package config

type Type string

var Types = struct {
	AwsTag                 Type
	AzureTag               Type
	Custom                 Type
	GcpLabel               Type
	GcpTag                 Type
	HostProcessProperty    Type
	K8SNamespaceAnnotation Type
	K8SNamespaceLabel      Type
	K8SPodAnnotation       Type
	K8SPodLabel            Type
	K8SWorkloadAnnotation  Type
	K8SWorkloadLabel       Type
}{
	"AWS_TAG",
	"AZURE_TAG",
	"CUSTOM",
	"GCP_LABEL",
	"GCP_TAG",
	"HOST_PROCESS_PROPERTY",
	"K8S_NAMESPACE_ANNOTATION",
	"K8S_NAMESPACE_LABEL",
	"K8S_POD_ANNOTATION",
	"K8S_POD_LABEL",
	"K8S_WORKLOAD_ANNOTATION",
	"K8S_WORKLOAD_LABEL",
}
