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

package attackprotectionallowlistconfig

import (
	attackprotectionallowlistconfig "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/appsec/attackprotectionallowlistconfig/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings/services/settings20"
)

const SchemaVersion = "1.0.1"
const SchemaID = "builtin:appsec.attack-protection-allowlist-config"

func Service(credentials *settings.Credentials) settings.CRUDService[*attackprotectionallowlistconfig.Settings] {
	return settings20.Service[*attackprotectionallowlistconfig.Settings](credentials, SchemaID, SchemaVersion)
}
