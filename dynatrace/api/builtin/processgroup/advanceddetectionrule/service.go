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

package advanceddetectionrule

import (
	"context"
	"fmt"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api"
	advanceddetectionrule "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/processgroup/advanceddetectionrule/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings/services/settings20"
)

const SchemaVersion = "0.0.5"
const SchemaID = "builtin:process-group.advanced-detection-rule"

func Service(credentials *settings.Credentials) settings.CRUDService[*advanceddetectionrule.Settings] {
	return settings20.Service(credentials, SchemaID, SchemaVersion, &settings20.ServiceOptions[*advanceddetectionrule.Settings]{Duplicates: Duplicates})
}

func Duplicates(ctx context.Context, service settings.RService[*advanceddetectionrule.Settings], v *advanceddetectionrule.Settings) (*api.Stub, error) {
	if settings.RejectDuplicate("dynatrace_process_group_detection") {
		var err error
		var stubs api.Stubs
		if stubs, err = service.List(ctx); err != nil {
			return nil, err
		}
		for _, stub := range stubs {
			config := stub.Value.(*advanceddetectionrule.Settings)
			if v.ProcessDetection.Property == config.ProcessDetection.Property && v.ProcessDetection.ContainedString == config.ProcessDetection.ContainedString {
				if (v.ProcessDetection.RestrictToProcessType == nil && config.ProcessDetection.RestrictToProcessType == nil) || (v.ProcessDetection.RestrictToProcessType != nil && config.ProcessDetection.RestrictToProcessType != nil && *v.ProcessDetection.RestrictToProcessType == *config.ProcessDetection.RestrictToProcessType) {
					return nil, fmt.Errorf("Advanced detection rule with condition already exists")
				}
			}
		}
	} else if settings.HijackDuplicate("dynatrace_process_group_detection") {
		var err error
		var stubs api.Stubs
		if stubs, err = service.List(ctx); err != nil {
			return nil, err
		}
		for _, stub := range stubs {
			config := stub.Value.(*advanceddetectionrule.Settings)
			if v.ProcessDetection.Property == config.ProcessDetection.Property && v.ProcessDetection.ContainedString == config.ProcessDetection.ContainedString {
				if (v.ProcessDetection.RestrictToProcessType == nil && config.ProcessDetection.RestrictToProcessType == nil) || (v.ProcessDetection.RestrictToProcessType != nil && config.ProcessDetection.RestrictToProcessType != nil && *v.ProcessDetection.RestrictToProcessType == *config.ProcessDetection.RestrictToProcessType) {
					return stub, nil
				}
			}
		}
	}
	return nil, nil
}
