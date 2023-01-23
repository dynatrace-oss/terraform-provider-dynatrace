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

package managementzones

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings"

	managementzones "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/managementzones/settings"
)

const SchemaID = "v1:config:management-zones"
const BasePath = "/api/config/v1/managementZones"

func Service(credentials *settings.Credentials) settings.CRUDService[*managementzones.ManagementZone] {
	return settings.NewCRUDService(
		credentials,
		SchemaID,
		settings.DefaultServiceOptions[*managementzones.ManagementZone](BasePath).Hijack(Hijack),
	)
}

func Hijack(err error, service settings.RService[*managementzones.ManagementZone], v *managementzones.ManagementZone) (*settings.Stub, error) {
	if rest.ContainsViolation(err, "Management zone must have a unique name.") {
		return settings.FindByName(service, settings.Name(v))
	}
	return nil, nil
}
