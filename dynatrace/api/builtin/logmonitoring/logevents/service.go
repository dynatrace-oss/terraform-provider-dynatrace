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

package logevents

import (
	"fmt"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api"
	logevents "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/logmonitoring/logevents/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings/services/settings20"
)

const SchemaVersion = "3.2.3"
const SchemaID = "builtin:logmonitoring.log-events"

func Service(credentials *settings.Credentials) settings.CRUDService[*logevents.Settings] {
	return settings20.Service[*logevents.Settings](credentials, SchemaID, SchemaVersion, &settings20.ServiceOptions[*logevents.Settings]{Duplicates: Duplicates})
}

func Duplicates(service settings.RService[*logevents.Settings], v *logevents.Settings) (*api.Stub, error) {
	if settings.RejectDuplicate("dynatrace_log_events") {
		var err error
		var stubs api.Stubs
		if stubs, err = service.List(); err != nil {
			return nil, err
		}
		for _, stub := range stubs {
			if v.Name() == stub.Name {
				return nil, fmt.Errorf("Log event summary with value `%s` already exists", v.Name())
			}
		}
	} else if settings.HijackDuplicate("dynatrace_log_events") {
		var err error
		var stubs api.Stubs
		if stubs, err = service.List(); err != nil {
			return nil, err
		}
		for _, stub := range stubs {
			if v.Name() == stub.Name {
				return stub, nil
			}
		}
	}
	return nil, nil
}
