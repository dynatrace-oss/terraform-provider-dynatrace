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

package aws

type ConsumersOfAwsRoleBasedAuthentication string
type ConsumersOfAwsWebIdentity string
type Type string

var Types = struct {
	AWSRoleBasedAuthentication Type
	AWSWebIdentity             Type
}{
	"awsRoleBasedAuthentication",
	"awsWebIdentity",
}

var ConsumersOfAwsRoleBasedAuthentications = struct {
	DataAcquisitionDeprecated ConsumersOfAwsRoleBasedAuthentication
	DataAcquisition           ConsumersOfAwsRoleBasedAuthentication
	CostAndCarbonOptimization ConsumersOfAwsRoleBasedAuthentication
	None                      ConsumersOfAwsRoleBasedAuthentication
}{
	"DA",
	"SVC:com.dynatrace.da",
	"APP:dynatrace.biz.carbon",
	"NONE",
}

var ConsumersOfAwsWebIdentities = struct {
	AWSConnector              ConsumersOfAwsWebIdentity
	CostAndCarbonOptimization ConsumersOfAwsWebIdentity
}{
	"APP:dynatrace.aws.connector",
	"APP:dynatrace.biz.carbon",
}
