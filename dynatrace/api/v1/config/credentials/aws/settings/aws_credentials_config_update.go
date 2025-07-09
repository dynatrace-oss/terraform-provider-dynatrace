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

type AWSCredentialsConfigUpdate struct {
	Label                            string                        `json:"label"`
	TaggedOnly                       bool                          `json:"taggedOnly"`
	AuthenticationData               *AWSAuthenticationData        `json:"authenticationData"`
	PartitionType                    PartitionType                 `json:"partitionType"`
	TagsToMonitor                    []*AWSConfigTag               `json:"tagsToMonitor"`
	SupportingServicesToMonitor      []*AWSSupportingServiceConfig `json:"supportingServicesToMonitor"`
	CredentialsEnabled               bool                          `json:"credentialsEnabled"`               // Enable monitoring of specified AWS credentials
	RunningOnDynatraceInfrastructure bool                          `json:"runningOnDynatraceInfrastructure"` // Run credentials on Dynatrace infrastructure
}
