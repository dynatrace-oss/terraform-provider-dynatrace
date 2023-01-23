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

package browser

import (
	"strings"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/synthetic/monitors"
	browser "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/synthetic/monitors/browser/settings"
)

const SchemaID = "v1:synthetic:monitors:browser"

func Service(credentials *settings.Credentials) settings.CRUDService[*browser.SyntheticMonitor] {
	return settings.NewCRUDService(credentials, SchemaID, &settings.ServiceOptions[*browser.SyntheticMonitor]{
		Get:            settings.Path("/api/v1/synthetic/monitors/%s"),
		List:           settings.Path("/api/v1/synthetic/monitors?type=BROWSER"),
		CreateURL:      func(v *browser.SyntheticMonitor) string { return "/api/v1/synthetic/monitors" },
		CreateRetry:    CreateRetry,
		Stubs:          &monitors.Monitors{},
		HasNoValidator: true,
		CreateConfirm:  20,
	})
}

func contains(elems []string, search string) bool {
	if len(elems) == 0 {
		return false
	}
	for _, elem := range elems {
		if elem == search {
			return true
		}
	}
	return false
}

func CreateRetry(v *browser.SyntheticMonitor, err error) *browser.SyntheticMonitor {
	message := err.Error()
	if strings.HasPrefix(message, "Unknown application(s):") {
		if len(v.ManuallyAssignedApps) == 0 {
			return nil
		}
		applicationIDs := strings.Split(strings.TrimSpace(strings.TrimPrefix(message, "Unknown application(s):")), ",")
		for idx, applicationID := range applicationIDs {
			applicationIDs[idx] = strings.TrimSpace(applicationID)
		}
		manuallyAssignedApps := []string{}
		for _, manuallyAssignedApp := range v.ManuallyAssignedApps {
			if !contains(applicationIDs, manuallyAssignedApp) {
				manuallyAssignedApps = append(manuallyAssignedApps, manuallyAssignedApp)
			}
		}
		v.ManuallyAssignedApps = manuallyAssignedApps
		return v
	}
	return nil
}
