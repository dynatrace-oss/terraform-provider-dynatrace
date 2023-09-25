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

package appdetection

import (
	"fmt"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api"
	appdetection "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/rum/web/appdetection/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings/services/settings20"
)

const SchemaVersion = "2.1.1"
const SchemaID = "builtin:rum.web.app-detection"

func Service(credentials *settings.Credentials) settings.CRUDService[*appdetection.Settings] {
	return settings20.Service(credentials, SchemaID, SchemaVersion, &settings20.ServiceOptions[*appdetection.Settings]{Duplicates: Duplicates})
}

func Duplicates(service settings.RService[*appdetection.Settings], v *appdetection.Settings) (*api.Stub, error) {
	if settings.RejectDuplicate("dynatrace_application_detection_rule_v2") {
		var err error
		var stubs api.Stubs
		if stubs, err = service.List(); err != nil {
			return nil, err
		}
		for _, stub := range stubs {
			config := stub.Value.(*appdetection.Settings)
			if v.Matcher == config.Matcher && v.Pattern == config.Pattern {
				return nil, fmt.Errorf("Application detection rule with filter already exists")
			}
		}
	} else if settings.HijackDuplicate("dynatrace_application_detection_rule_v2") {
		var err error
		var stubs api.Stubs
		if stubs, err = service.List(); err != nil {
			return nil, err
		}
		for _, stub := range stubs {
			config := stub.Value.(*appdetection.Settings)
			if v.Matcher == config.Matcher && v.Pattern == config.Pattern {
				return stub, nil
			}
		}
	}
	return nil, nil
}
