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

package location

type DeploymentType string

var DeploymentTypes = struct {
	Kubernetes DeploymentType
	Openshift  DeploymentType
	Standard   DeploymentType
	Unknown    DeploymentType
}{
	DeploymentType("KUBERNETES"),
	DeploymentType("OPENSHIFT"),
	DeploymentType("STANDARD"),
	DeploymentType("UNKNOWN"),
}

func (me DeploymentType) Ref() *DeploymentType {
	return &me
}

type LocationType string

var LocationTypes = struct {
	Private LocationType
}{
	LocationType("PRIVATE"),
	// other types aren't supported for the location endpoint.
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

func (me Status) Ref() *Status {
	return &me
}

type FipsMode string

var FipsModes = struct {
	Disabled                  FipsMode
	Enabled                   FipsMode
	EnabledWithCorporateProxy FipsMode
}{
	FipsMode("DISABLED"),
	FipsMode("ENABLED"),
	FipsMode("ENABLED_WITH_CORPORATE_PROXY"),
}

func (me FipsMode) Ref() *FipsMode {
	return &me
}

type NodeSize string

var NodeSizes = struct {
	M           NodeSize
	S           NodeSize
	Unsupported NodeSize
	XS          NodeSize
}{
	NodeSize("M"),
	NodeSize("S"),
	NodeSize("UNSUPPORTED"),
	NodeSize("XS"),
}

func (me NodeSize) Ref() *NodeSize {
	return &me
}
