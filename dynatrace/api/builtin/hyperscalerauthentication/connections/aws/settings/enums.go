/**
* @license
* Copyright 2025 Dynatrace LLC
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

var ConsumersOfAwsRoleBasedAuthentications = struct {
	AppDynatraceBizCarbon       ConsumersOfAwsRoleBasedAuthentication
	Da                          ConsumersOfAwsRoleBasedAuthentication
	None                        ConsumersOfAwsRoleBasedAuthentication
	SvcComDynatraceBo           ConsumersOfAwsRoleBasedAuthentication
	SvcComDynatraceDa           ConsumersOfAwsRoleBasedAuthentication
	SvcComDynatraceGrail        ConsumersOfAwsRoleBasedAuthentication
	SvcComDynatraceOpenpipeline ConsumersOfAwsRoleBasedAuthentication
}{
	"APP:dynatrace.biz.carbon",
	"DA",
	"NONE",
	"SVC:com.dynatrace.bo",
	"SVC:com.dynatrace.da",
	"SVC:com.dynatrace.grail",
	"SVC:com.dynatrace.openpipeline",
}

type ConsumersOfAwsWebIdentity string

var ConsumersOfAwsWebIdentities = struct {
	AppDynatraceAwsConnector ConsumersOfAwsWebIdentity
	AppDynatraceBizCarbon    ConsumersOfAwsWebIdentity
}{
	"APP:dynatrace.aws.connector",
	"APP:dynatrace.biz.carbon",
}

type Type string

var Types = struct {
	Awsrolebasedauthentication Type
	Awswebidentity             Type
}{
	"awsRoleBasedAuthentication",
	"awsWebIdentity",
}
