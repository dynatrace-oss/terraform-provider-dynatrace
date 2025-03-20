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

package cloudfoundry

import (
	"sync"

	cloudfoundry "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/credentials/cloudfoundry/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings"
)

const SchemaID = "v1:config:credentials:cloudfoundry"
const BasePath = "/api/config/v1/cloudFoundry/credentials"

var mu sync.Mutex

func Service(credentials *rest.Credentials) settings.CRUDService[*cloudfoundry.CloudFoundryCredentials] {
	return settings.NewAPITokenService(
		credentials,
		SchemaID,
		settings.DefaultServiceOptions[*cloudfoundry.CloudFoundryCredentials](BasePath).WithMutex(mu.Lock, mu.Unlock),
	)
}
