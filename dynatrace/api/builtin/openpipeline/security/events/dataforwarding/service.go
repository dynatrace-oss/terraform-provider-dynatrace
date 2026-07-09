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

package dataforwarding

import (
	service "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/openpipeline/security/events/dataforwarding/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings/services/settings20"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/provider/config"
)

const SchemaVersion = "1.68"
const SchemaID = "builtin:openpipeline.security.events.data-forwarding"

func Service(credentials *config.ProviderConfiguration) settings.CRUDService[*service.Settings] {
	return settings20.Service[*service.Settings](credentials, SchemaID, SchemaVersion)
}
