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

package databaseservices

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings"

	databaseservices "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/anomalies/databaseservices/settings"
)

const SchemaID = "v1:config:anomaly-detection:database-services"
const EndpointURL = "/api/config/v1/anomalyDetection/databaseServices"
const StaticID = "70276f24-2efe-4203-89f9-4fd3b7d7e5d9"
const StaticName = "database_anomalies"

func Service(credentials *settings.Credentials) settings.CRUDService[*databaseservices.AnomalyDetection] {
	return settings.StaticService[*databaseservices.AnomalyDetection](
		credentials,
		SchemaID,
		EndpointURL,
		settings.Stub{ID: StaticID, Name: StaticName},
	)
}
