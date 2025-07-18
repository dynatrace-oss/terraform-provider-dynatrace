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

package customprocessmonitoring

import (
	"context"
	"fmt"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api"
	customprocessmonitoring "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/process/monitoring/custom/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings/services/settings20"
)

const SchemaVersion = "0.4.9"
const SchemaID = "builtin:process.custom-process-monitoring-rule"

func Service(credentials *rest.Credentials) settings.CRUDService[*customprocessmonitoring.Settings] {
	return settings20.Service(credentials, SchemaID, SchemaVersion, &settings20.ServiceOptions[*customprocessmonitoring.Settings]{Duplicates: Duplicates})
}

func Duplicates(ctx context.Context, service settings.RService[*customprocessmonitoring.Settings], v *customprocessmonitoring.Settings) (*api.Stub, error) {
	if settings.RejectDuplicate("dynatrace_process_monitoring_rule") {
		var err error
		var stubs api.Stubs
		if stubs, err = service.List(ctx); err != nil {
			return nil, err
		}
		for _, stub := range stubs {
			config := stub.Value.(*customprocessmonitoring.Settings)
			if v.HostGroupID == config.HostGroupID && v.Condition.Item == config.Condition.Item && v.Condition.Operator == config.Condition.Operator {
				if (v.Condition.Value == nil && config.Condition.Value == nil) || (v.Condition.Value != nil && config.Condition.Value != nil && *v.Condition.Value == *config.Condition.Value) {
					return nil, fmt.Errorf("Process monitoring rule with condition already exists")
				}
			}
		}
	} else if settings.HijackDuplicate("dynatrace_process_monitoring_rule") {
		var err error
		var stubs api.Stubs
		if stubs, err = service.List(ctx); err != nil {
			return nil, err
		}
		for _, stub := range stubs {
			config := stub.Value.(*customprocessmonitoring.Settings)
			if v.HostGroupID == config.HostGroupID && v.Condition.Item == config.Condition.Item && v.Condition.Operator == config.Condition.Operator {
				if (v.Condition.Value == nil && config.Condition.Value == nil) || (v.Condition.Value != nil && config.Condition.Value != nil && *v.Condition.Value == *config.Condition.Value) {
					return stub, nil
				}
			}
		}
	}
	return nil, nil
}
