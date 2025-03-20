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

package services

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api"
	services "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/anomalies/services/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings"
)

const SchemaID = "v1:config:anomaly-detection:services"
const EndpointURL = "/api/config/v1/anomalyDetection/services"
const StaticID = "aa622c87-3a39-42b0-85fb-62ba48758482"
const StaticName = "service_anomalies"

func Service(credentials *rest.Credentials) settings.CRUDService[*services.AnomalyDetection] {
	return settings.APITokenStaticService[*services.AnomalyDetection](
		credentials,
		SchemaID,
		EndpointURL,
		api.Stub{ID: StaticID, Name: StaticName},
	)
}
