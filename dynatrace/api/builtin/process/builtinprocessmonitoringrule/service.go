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

package builtinprocessmonitoringrule

import (
	service "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/process/builtinprocessmonitoringrule/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings/services/settings20"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/provider/config"
)

const SchemaVersion = "0.4.7"
const SchemaID = "builtin:process.built-in-process-monitoring-rule"

func Service(credentials *config.ProviderConfiguration) settings.CRUDService[*service.Settings] {
	return settings20.Service[*service.Settings](credentials, SchemaID, SchemaVersion)
}
