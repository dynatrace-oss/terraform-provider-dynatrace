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

package hosts

import (
	hosts "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/anomalies/hosts/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings"
)

const SchemaID = "v1:config:anomaly-detection:hosts"
const EndpointURL = "/api/config/v1/anomalyDetection/hosts"
const StaticID = "7100e39c-a80c-4ebc-a431-211a938cb3ee"
const StaticName = "host_anomalies"

func Service(credentials *settings.Credentials) settings.CRUDService[*hosts.AnomalyDetection] {
	return settings.StaticService[*hosts.AnomalyDetection](
		credentials,
		SchemaID,
		EndpointURL,
		settings.Stub{ID: StaticID, Name: StaticName},
	)
}
