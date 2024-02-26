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

package maintenancewindow

import (
	"fmt"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api"
	maintenancewindow "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/alerting/maintenancewindow/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings/services/settings20"
)

const SchemaID = "builtin:alerting.maintenance-window"
const SchemaVersion = "2.14.5"

func Service(credentials *settings.Credentials) settings.CRUDService[*maintenancewindow.Settings] {
	return settings20.Service(credentials, SchemaID, SchemaVersion, &settings20.ServiceOptions[*maintenancewindow.Settings]{LegacyID: settings.LegacyObjIDDecode, Duplicates: Duplicates})
}

func Duplicates(service settings.RService[*maintenancewindow.Settings], v *maintenancewindow.Settings) (*api.Stub, error) {
	if settings.RejectDuplicate("dynatrace_maintenance") {
		var err error
		var stubs api.Stubs
		if stubs, err = service.List(); err != nil {
			return nil, err
		}
		for _, stub := range stubs {
			if v.GeneralProperties != nil && v.GeneralProperties.Name == stub.Name {
				return nil, fmt.Errorf("Maintenance Window named '%s' already exists", v.GeneralProperties.Name)
			}
		}
	} else if settings.HijackDuplicate("dynatrace_maintenance") {
		var err error
		var stubs api.Stubs
		if stubs, err = service.List(); err != nil {
			return nil, err
		}
		for _, stub := range stubs {
			if v.GeneralProperties != nil && v.GeneralProperties.Name == stub.Name {
				return stub, nil
			}
		}
	}
	return nil, nil
}
