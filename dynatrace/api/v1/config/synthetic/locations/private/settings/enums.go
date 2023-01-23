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

package locations

type CloudPlatform string

var CloudPlatforms = struct {
	Alibaba        CloudPlatform
	AmazonEC2      CloudPlatform
	Azure          CloudPlatform
	DynatraceCloud CloudPlatform
	GoogleCloud    CloudPlatform
	Interoute      CloudPlatform
	Other          CloudPlatform
	Undefined      CloudPlatform
}{
	CloudPlatform("ALIBABA"),
	CloudPlatform("AMAZON_EC2"),
	CloudPlatform("AZURE"),
	CloudPlatform("DYNATRACE_CLOUD"),
	CloudPlatform("GOOGLE_CLOUD"),
	CloudPlatform("INTEROUTE"),
	CloudPlatform("OTHER"),
	CloudPlatform("UNDEFINED"),
}

type DeploymentType string

var DeploymentTypes = struct {
	Kubernetes DeploymentType
	Standard   DeploymentType
}{
	DeploymentType("KUBERNETES"),
	DeploymentType("STANDARD"),
}

func (me DeploymentType) Ref() *DeploymentType {
	return &me
}

type LocationType string

var LocationTypes = struct {
	Cluster LocationType
	Private LocationType
	Public  LocationType
}{
	LocationType("CLUSTER"),
	LocationType("PRIVATE"),
	LocationType("PUBLIC"),
}

type Stage string

var Stages = struct {
	Beta       Stage
	ComingSoon Stage
	Deleted    Stage
	GA         Stage
}{
	Stage("BETA"),
	Stage("COMING_SOON"),
	Stage("DELETED"),
	Stage("GA"),
}

type Status string

var Statuses = struct {
	Disabled Status
	Enabled  Status
	Hidden   Status
}{
	Status("DISABLED"),
	Status("ENABLED"),
	Status("HIDDEN"),
}
