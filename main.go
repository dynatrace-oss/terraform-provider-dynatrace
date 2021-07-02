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

package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/provider"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
)

var resArr = []string{
	"dynatrace_custom_service",
	"dynatrace_dashboard",
	"dynatrace_management_zone",
	"dynatrace_maintenance_window",
	"dynatrace_request_attribute",
	"dynatrace_alerting_profile",
	"dynatrace_notification",
	"dynatrace_autotag",
	"dynatrace_aws_credentials",
	"dynatrace_azure_credentials",
	"dynatrace_k8s_credentials",
	"dynatrace_service_anomalies",
	"dynatrace_application_anomalies",
	"dynatrace_host_anomalies",
	"dynatrace_database_anomalies",
	"dynatrace_custom_anomalies",
	"dynatrace_disk_anomalies",
	"dynatrace_calculated_service_metric",
	"dynatrace_service_naming",
	"dynatrace_host_naming",
	"dynatrace_processgroup_naming",
	"dynatrace_slo",
}

func matchRes(keyVal string) (string, string) {
	res1 := ""
	res2 := ""
	for _, entry := range resArr {
		if strings.HasPrefix(keyVal, entry) {
			res1 = entry
			if strings.HasPrefix(keyVal, entry+"=") {
				res2 = keyVal[len(entry)+1:]
			}
		}
	}
	return res1, res2
}

func export(args []string) bool {
	if len(args) == 1 {
		return false
	}
	if strings.TrimSpace(args[1]) != "export" {
		return false
	}
	environmentURL := os.Getenv("DYNATRACE_ENV_URL")
	if environmentURL == "" {
		fmt.Println("The environment variable DYNATRACE_ENV_URL needs to be set")
		os.Exit(0)
	}
	apiToken := os.Getenv("DYNATRACE_API_TOKEN")
	if apiToken == "" {
		fmt.Println("The environment variable DYNATRACE_API_TOKEN needs to be set")
		os.Exit(0)
	}
	targetFolder := os.Getenv("DYNATRACE_TARGET_FOLDER")
	if targetFolder == "" {
		fmt.Println("The environment variable DYNATRACE_TARGET_FOLDER has not been set - using folder 'configuration' as default")
		targetFolder = "configuration"
	}
	if len(args) == 2 {
		return downloadWith(environmentURL, apiToken, targetFolder, nil)
	}

	m := map[string][]string{}
	for idx := 2; idx < len(args); idx++ {
		key, id := matchRes(args[idx])
		if key == "" {
			fmt.Println("Unknown resource `" + args[idx] + "`")
			os.Exit(0)
		}
		stored, ok := m[key]
		if ok {
			if stored != nil {
				if id == "" {
					m[key] = nil
				} else {
					stored = append(stored, id)
					m[key] = stored
				}
			}
		} else {
			if id == "" {
				m[key] = nil
			} else {
				stored = []string{id}
				m[key] = stored
			}
		}
	}
	return downloadWith(environmentURL, apiToken, targetFolder, m)
}

func download(args []string) bool {
	if len(args) == 1 {
		return false
	}
	if strings.TrimSpace(args[1]) != "download" {
		return false
	}

	if len(args) < 3 {
		fmt.Println("Usage: terraform-provider-dynatrace download <environment-url> <api-token> [<target-folder>")
		os.Exit(0)
	}
	targetFolder := "configuration"
	environmentURL := args[2]
	apiToken := args[3]
	if len(args) > 4 {
		targetFolder = args[4]
	}
	return downloadWith(environmentURL, apiToken, targetFolder, nil)
}

func downloadWith(environmentURL string, apiToken string, targetFolder string, m map[string][]string) bool {
	if m == nil {
		m = map[string][]string{
			"dynatrace_custom_service":            nil,
			"dynatrace_dashboard":                 nil,
			"dynatrace_management_zone":           nil,
			"dynatrace_maintenance_window":        nil,
			"dynatrace_request_attribute":         nil,
			"dynatrace_alerting_profile":          nil,
			"dynatrace_notification":              nil,
			"dynatrace_autotag":                   nil,
			"dynatrace_aws_credentials":           nil,
			"dynatrace_azure_credentials":         nil,
			"dynatrace_k8s_credentials":           nil,
			"dynatrace_service_anomalies":         nil,
			"dynatrace_application_anomalies":     nil,
			"dynatrace_host_anomalies":            nil,
			"dynatrace_database_anomalies":        nil,
			"dynatrace_custom_anomalies":          nil,
			"dynatrace_disk_anomalies":            nil,
			"dynatrace_calculated_service_metric": nil,
			"dynatrace_service_naming":            nil,
			"dynatrace_host_naming":               nil,
			"dynatrace_processgroup_naming":       nil,
			"dynatrace_slo":                       nil,
		}
	}
	if err := os.RemoveAll(targetFolder); err != nil {
		fmt.Println(err.Error())
		os.Exit(0)
	}

	if rids, ok := m["dynatrace_autotag"]; ok {
		if err := importAutoTags(targetFolder+"/autotags", environmentURL, apiToken, rids); err != nil {
			fmt.Println(err.Error())
			os.Exit(0)
		}
	}
	if rids, ok := m["dynatrace_dashboard"]; ok {
		if err := importDashboards(targetFolder+"/dashboards", environmentURL, apiToken, rids); err != nil {
			fmt.Println(err.Error())
			os.Exit(0)
		}
	}
	if rids, ok := m["dynatrace_management_zone"]; ok {
		if err := importManagementZones(targetFolder+"/management_zones", environmentURL, apiToken, rids); err != nil {
			fmt.Println(err.Error())
			os.Exit(0)
		}
	}
	if rids, ok := m["dynatrace_custom_service"]; ok {
		if err := importCustomServices(targetFolder+"/custom_services", environmentURL, apiToken, rids); err != nil {
			fmt.Println(err.Error())
			os.Exit(0)
		}
	}
	if rids, ok := m["dynatrace_alerting_profile"]; ok {
		if err := importAlertingProfiles(targetFolder+"/alerting_profiles", environmentURL, apiToken, rids); err != nil {
			fmt.Println(err.Error())
			os.Exit(0)
		}
	}
	if rids, ok := m["dynatrace_request_attribute"]; ok {
		if err := importRequestAttributes(targetFolder+"/request_attributes", environmentURL, apiToken, rids); err != nil {
			fmt.Println(err.Error())
			os.Exit(0)
		}
	}
	if rids, ok := m["dynatrace_notification"]; ok {
		if err := importNotificationConfigs(targetFolder+"/notifications", environmentURL, apiToken, rids); err != nil {
			fmt.Println(err.Error())
			os.Exit(0)
		}
	}
	if rids, ok := m["dynatrace_aws_credentials"]; ok {
		if err := importAWSCredentials(targetFolder+"/credentials", environmentURL, apiToken, rids); err != nil {
			fmt.Println(err.Error())
			os.Exit(0)
		}
	}
	if rids, ok := m["dynatrace_azure_credentials"]; ok {
		if err := importAzureCredentials(targetFolder+"/credentials", environmentURL, apiToken, rids); err != nil {
			fmt.Println(err.Error())
			os.Exit(0)
		}
	}
	if rids, ok := m["dynatrace_k8s_credentials"]; ok {
		if err := importK8sCredentials(targetFolder+"/credentials", environmentURL, apiToken, rids); err != nil {
			fmt.Println(err.Error())
			os.Exit(0)
		}
	}
	if rids, ok := m["dynatrace_maintenance_window"]; ok {
		if err := importMaintenance(targetFolder+"/maintenance", environmentURL, apiToken, rids); err != nil {
			fmt.Println(err.Error())
			os.Exit(0)
		}
	}
	if rids, ok := m["dynatrace_disk_anomalies"]; ok {
		if err := importDiskAnomalies(targetFolder+"/anomalies/disks", environmentURL, apiToken, rids); err != nil {
			fmt.Println(err.Error())
			os.Exit(0)
		}
	}
	if rids, ok := m["dynatrace_custom_anomalies"]; ok {
		if err := importMetricAnomalies(targetFolder+"/anomalies/custom", environmentURL, apiToken, rids); err != nil {
			fmt.Println(err.Error())
			os.Exit(0)
		}
	}
	if _, ok := m["dynatrace_database_anomalies"]; ok {
		if err := importDatabaseAnomalies(targetFolder+"/anomalies", environmentURL, apiToken); err != nil {
			fmt.Println(err.Error())
			os.Exit(0)
		}
	}
	if _, ok := m["dynatrace_host_anomalies"]; ok {
		if err := importHostAnomalies(targetFolder+"/anomalies", environmentURL, apiToken); err != nil {
			fmt.Println(err.Error())
			os.Exit(0)
		}
	}
	if _, ok := m["dynatrace_application_anomalies"]; ok {
		if err := importApplicationAnomalies(targetFolder+"/anomalies", environmentURL, apiToken); err != nil {
			fmt.Println(err.Error())
			os.Exit(0)
		}
	}
	if _, ok := m["dynatrace_service_anomalies"]; ok {
		if err := importServiceAnomalies(targetFolder+"/anomalies", environmentURL, apiToken); err != nil {
			fmt.Println(err.Error())
			os.Exit(0)
		}
	}
	if rids, ok := m["dynatrace_calculated_service_metric"]; ok {
		if err := importCalculatedServiceMetrics(targetFolder+"/metrics/calculated/service", environmentURL, apiToken, rids); err != nil {
			fmt.Println(err.Error())
			os.Exit(0)
		}
	}
	if rids, ok := m["dynatrace_service_naming"]; ok {
		if err := importServiceNamings(targetFolder+"/naming/services", environmentURL, apiToken, rids); err != nil {
			fmt.Println(err.Error())
			os.Exit(0)
		}
	}
	if rids, ok := m["dynatrace_host_naming"]; ok {
		if err := importHostNamings(targetFolder+"/naming/hosts", environmentURL, apiToken, rids); err != nil {
			fmt.Println(err.Error())
			os.Exit(0)
		}
	}
	if rids, ok := m["dynatrace_processgroup_naming"]; ok {
		if err := importProcessGroupNamings(targetFolder+"/naming/process_groups", environmentURL, apiToken, rids); err != nil {
			fmt.Println(err.Error())
			os.Exit(0)
		}
	}
	if rids, ok := m["dynatrace_slo"]; ok {
		if err := importSLOs(targetFolder+"/slo", environmentURL, apiToken, rids); err != nil {
			fmt.Println(err.Error())
			os.Exit(0)
		}
	}

	return true
}

func main() {
	if export(os.Args) {
		return
	}
	if download(os.Args) {
		return
	}
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: func() *schema.Provider {
			return provider.Provider()
		},
	})
}
