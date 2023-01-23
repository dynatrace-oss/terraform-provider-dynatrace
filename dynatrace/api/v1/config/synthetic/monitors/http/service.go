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

package http

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/synthetic/monitors"
	http "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/synthetic/monitors/http/settings"
)

const SchemaID = "v1:synthetic:monitors:http"

func Service(credentials *settings.Credentials) settings.CRUDService[*http.SyntheticMonitor] {
	return settings.NewCRUDService(credentials, SchemaID, &settings.ServiceOptions[*http.SyntheticMonitor]{
		Get:            settings.Path("/api/v1/synthetic/monitors/%s"),
		List:           settings.Path("/api/v1/synthetic/monitors?type=HTTP"),
		CreateURL:      func(v *http.SyntheticMonitor) string { return "/api/v1/synthetic/monitors" },
		Stubs:          &monitors.Monitors{},
		HasNoValidator: true,
		CreateConfirm:  20,
	})
}
