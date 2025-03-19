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

package rummobilecrashrateincrease

import (
	crashrate "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/anomalydetection/rum/mobile/crashrate/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings/services/settings20"
)

const SchemaVersion = "1.0.1"
const SchemaID = "builtin:anomaly-detection.rum-mobile-crash-rate-increase"

func Service(credentials *rest.Credentials) settings.CRUDService[*crashrate.Settings] {
	return settings20.Service[*crashrate.Settings](credentials, SchemaID, SchemaVersion)
}
