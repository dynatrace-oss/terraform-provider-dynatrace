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
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/dtcookie/dynatrace/api/config/v2/notifications"
	"github.com/dtcookie/dynatrace/rest"
	downloadv2 "github.com/dynatrace-oss/terraform-provider-dynatrace/download"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/provider"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
)

var clusterResArr = []string{
	"dynatrace_environment",
}

var resArr = []string{
	"dynatrace_custom_service",
	"dynatrace_dashboard",
	"dynatrace_management_zone",
	// "dynatrace_maintenance_window",
	"dynatrace_request_attribute",
	"dynatrace_alerting",
	// "dynatrace_notification", // removed in favor of Settings 2.0 Notifications
	"dynatrace_autotag",
	"dynatrace_aws_credentials",
	"dynatrace_azure_credentials",
	"dynatrace_k8s_credentials",
	"dynatrace_cloudfoundry_credentials",
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
	"dynatrace_span_entry_point",
	"dynatrace_span_capture_rule",
	"dynatrace_span_context_propagation",
	"dynatrace_resource_attributes",
	"dynatrace_span_attribute",
	"dynatrace_mobile_application",
	// "dynatrace_credentials",
	"dynatrace_browser_monitor",
	"dynatrace_http_monitor",
	"dynatrace_web_application",
	"dynatrace_application_data_privacy",
	"dynatrace_application_error_rules",
	"dynatrace_request_naming",
	"dynatrace_key_requests",
	"dynatrace_queue_manager",
	"dynatrace_ibm_mq_filters",
	"dynatrace_queue_sharing_groups",
	"dynatrace_ims_bridges",
	"dynatrace_network_zones",
	"dynatrace_application_detection_rule",
	"dynatrace_maintenance",
	"dynatrace_frequent_issues",
	"dynatrace_ansible_tower_notification",
	"dynatrace_email_notification",
	"dynatrace_jira_notification",
	"dynatrace_ops_genie_notification",
	"dynatrace_pager_duty_notification",
	"dynatrace_service_now_notification",
	"dynatrace_slack_notification",
	"dynatrace_trello_notification",
	"dynatrace_victor_ops_notification",
	"dynatrace_webhook_notification",
	"dynatrace_xmatters_notification",
}

func matchClusterRes(keyVal string) (string, string) {
	res1 := ""
	res2 := ""
	for _, entry := range clusterResArr {
		if strings.HasPrefix(keyVal, entry) {
			res1 = entry
			if strings.HasPrefix(keyVal, entry+"=") {
				res2 = keyVal[len(entry)+1:]
			}
		}
	}
	return res1, res2
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

func clusterExport(args []string) bool {
	if len(args) == 1 {
		return false
	}
	if strings.TrimSpace(args[1]) != "cluster_export" {
		return false
	}
	clusterURL := os.Getenv("DYNATRACE_CLUSTER_URL")
	if clusterURL == "" {
		fmt.Println("The environment variable DYNATRACE_CLUSTER_URL needs to be set")
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
	debug := os.Getenv("DYNATRACE_DEBUG")
	if debug == "true" {
		rest.Verbose = true
	}
	if len(args) == 2 {
		return clusterDownloadWith(clusterURL, apiToken, targetFolder, nil)
	}
	m := map[string][]string{}
	for idx := 2; idx < len(args); idx++ {
		key, id := matchClusterRes(args[idx])
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
	return clusterDownloadWith(clusterURL, apiToken, targetFolder, m)
}

func exportv2(args []string) bool {
	if len(args) == 1 {
		return false
	}
	if strings.TrimSpace(args[1]) != "-exportv2" {
		return false
	}
	environmentURL := os.Getenv("DYNATRACE_ENV_URL")
	if environmentURL == "" {
		fmt.Println("The environment variable DYNATRACE_ENV_URL needs to be set")
		os.Exit(0)
	}
	environmentURL = strings.TrimSuffix(strings.TrimSuffix(environmentURL, " "), "/")
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

	flag.Bool("exportv2", true, "")
	// fileArg := flag.Bool("single", false, "single file resource output `true` or `false` (default \"false\")")
	refArg := flag.Bool("ref", false, "enable data sources and dependencies")
	comIdArg := flag.Bool("id", false, "enable commented ids")
	// repIdArg := flag.String("ref", "hardcoded", "references option `hardcoded`, `resource` or `datasource`")
	flag.Parse()
	tailArgs := flag.Args()

	// if *repIdArg != "hardcoded" && *repIdArg != "resource" && *repIdArg != "datasource" {
	// 	fmt.Println("Unknown replace id option `" + *repIdArg + "`")
	// 	flag.CommandLine.Usage()
	// 	os.Exit(0)
	// }

	resArgs := map[string][]string{}
	for _, idx := range tailArgs {
		key, id := downloadv2.ValidateResource(idx)
		if key == "" {
			fmt.Println("Unknown resource `" + idx + "`")
			os.Exit(0)
		}
		stored, ok := resArgs[key]
		if ok {
			if stored != nil {
				if id == "" {
					resArgs[key] = nil
				} else {
					stored = append(stored, id)
					resArgs[key] = stored
				}
			}
		} else {
			if id == "" {
				resArgs[key] = nil
			} else {
				stored = []string{id}
				resArgs[key] = stored
			}
		}
	}

	// return downloadv2.Download(environmentURL, apiToken, targetFolder, *fileArg, *comIdArg, *repIdArg, resArgs)
	return downloadv2.Download(environmentURL, apiToken, targetFolder, *refArg, *comIdArg, resArgs)
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
	environmentURL = strings.TrimSuffix(strings.TrimSuffix(environmentURL, " "), "/")
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

func clusterDownloadWith(clusterURL string, apiToken string, targetFolder string, m map[string][]string) bool {
	if m == nil {
		m = map[string][]string{}
		for _, resn := range clusterResArr {
			m[resn] = nil
		}
	}
	if err := os.RemoveAll(targetFolder); err != nil {
		fmt.Println(err.Error())
		os.Exit(0)
	}
	if rids, ok := m["dynatrace_environment"]; ok {
		if err := importEnvironments(targetFolder+"/environments", clusterURL, apiToken, rids); err != nil {
			fmt.Println(err.Error())
			os.Exit(0)
		}
	}
	return true
}

func downloadWith(environmentURL string, apiToken string, targetFolder string, m map[string][]string) bool {
	if m == nil {
		m = map[string][]string{}
		for _, resn := range resArr {
			m[resn] = nil
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
	if rids, ok := m["dynatrace_alerting"]; ok {
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
	if rids, ok := m["dynatrace_cloudfoundry_credentials"]; ok {
		if err := importCloudFoundryCredentials(targetFolder+"/credentials", environmentURL, apiToken, rids); err != nil {
			fmt.Println(err.Error())
			os.Exit(0)
		}
	}
	/* if rids, ok := m["dynatrace_maintenance_window"]; ok {
		if err := importMaintenance(targetFolder+"/maintenance", environmentURL, apiToken, rids); err != nil {
			fmt.Println(err.Error())
			os.Exit(0)
		}
	} */
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
	if rids, ok := m["dynatrace_span_entry_point"]; ok {
		if err := importSpanEntryPoints(targetFolder+"/span/entry_points", environmentURL, apiToken, rids); err != nil {
			fmt.Println(err.Error())
			os.Exit(0)
		}
	}
	if rids, ok := m["dynatrace_span_capture_rule"]; ok {
		if err := importSpanCaptureRules(targetFolder+"/span/capture", environmentURL, apiToken, rids); err != nil {
			fmt.Println(err.Error())
			os.Exit(0)
		}
	}
	if rids, ok := m["dynatrace_span_context_propagation"]; ok {
		if err := importSpanContextPropagation(targetFolder+"/span/ctxprop", environmentURL, apiToken, rids); err != nil {
			fmt.Println(err.Error())
			os.Exit(0)
		}
	}
	if rids, ok := m["dynatrace_resource_attributes"]; ok {
		if err := importResourceAttributes(targetFolder+"/span/resource_attributes", environmentURL, apiToken, rids); err != nil {
			fmt.Println(err.Error())
			os.Exit(0)
		}
	}
	if rids, ok := m["dynatrace_span_attribute"]; ok {
		if err := importSpanAttributes(targetFolder+"/span/span_attributes", environmentURL, apiToken, rids); err != nil {
			fmt.Println(err.Error())
			os.Exit(0)
		}
	}
	if rids, ok := m["dynatrace_mobile_application"]; ok {
		if err := importMobileApps(targetFolder+"/applications/mobile", environmentURL, apiToken, rids); err != nil {
			fmt.Println(err.Error())
			os.Exit(0)
		}
	}
	// if rids, ok := m["dynatrace_credentials"]; ok {
	// 	if err := importVaultCredentials(targetFolder+"/credentials", environmentURL, apiToken, rids); err != nil {
	// 		fmt.Println(err.Error())
	// 		os.Exit(0)
	// 	}
	// }
	if rids, ok := m["dynatrace_browser_monitor"]; ok {
		if err := importBrowserMonitors(targetFolder+"/synthetic/browser", environmentURL, apiToken, rids); err != nil {
			fmt.Println(err.Error())
			os.Exit(0)
		}
	}
	if rids, ok := m["dynatrace_http_monitor"]; ok {
		if err := importHTTPMonitors(targetFolder+"/synthetic/http", environmentURL, apiToken, rids); err != nil {
			fmt.Println(err.Error())
			os.Exit(0)
		}
	}
	if rids, ok := m["dynatrace_web_application"]; ok {
		if err := importWebApps(targetFolder+"/applications/web", environmentURL, apiToken, rids); err != nil {
			fmt.Println(err.Error())
			os.Exit(0)
		}
	}
	if rids, ok := m["dynatrace_application_data_privacy"]; ok {
		if err := importApplicationDataPrivacy(targetFolder+"/applications/web", environmentURL, apiToken, rids); err != nil {
			fmt.Println(err.Error())
			os.Exit(0)
		}
	}
	if rids, ok := m["dynatrace_application_error_rules"]; ok {
		if err := importApplicationErrorRules(targetFolder+"/applications/web", environmentURL, apiToken, rids); err != nil {
			fmt.Println(err.Error())
			os.Exit(0)
		}
	}
	if rids, ok := m["dynatrace_request_naming"]; ok {
		if err := importRequestNamings(targetFolder+"/requestnaming", environmentURL, apiToken, rids); err != nil {
			fmt.Println(err.Error())
			os.Exit(0)
		}
	}
	if rids, ok := m["dynatrace_key_requests"]; ok {
		if err := importKeyRequests(targetFolder+"/keyrequests", environmentURL, apiToken, rids); err != nil {
			fmt.Println(err.Error())
			os.Exit(0)
		}
	}
	if rids, ok := m["dynatrace_queue_manager"]; ok {
		if err := importMQQueueManagers(targetFolder+"/ibmmq/queuemanagers", environmentURL, apiToken, rids); err != nil {
			fmt.Println(err.Error())
			os.Exit(0)
		}
	}
	if rids, ok := m["dynatrace_ibm_mq_filters"]; ok {
		if err := importMQFilters(targetFolder+"/ibmmq/filters", environmentURL, apiToken, rids); err != nil {
			fmt.Println(err.Error())
			os.Exit(0)
		}
	}
	if rids, ok := m["dynatrace_queue_sharing_groups"]; ok {
		if err := importMQQueueSharingGroups(targetFolder+"/ibmmq/queuesharinggroups", environmentURL, apiToken, rids); err != nil {
			fmt.Println(err.Error())
			os.Exit(0)
		}
	}
	if rids, ok := m["dynatrace_ims_bridges"]; ok {
		if err := importMQIMSBridges(targetFolder+"/ibmmq/imsbridges", environmentURL, apiToken, rids); err != nil {
			fmt.Println(err.Error())
			os.Exit(0)
		}
	}
	if rids, ok := m["dynatrace_network_zones"]; ok {
		if err := importNetworkZones(targetFolder+"/networkzones", environmentURL, apiToken, rids); err != nil {
			fmt.Println(err.Error())
			os.Exit(0)
		}
	}
	if rids, ok := m["dynatrace_application_detection_rule"]; ok {
		if err := importApplicationDetectionRules(targetFolder+"/applications/web/applicationdetectionrules", environmentURL, apiToken, rids); err != nil {
			fmt.Println(err.Error())
			os.Exit(0)
		}
	}
	if rids, ok := m["dynatrace_maintenance"]; ok {
		if err := importMaintenanceV2(targetFolder+"/maintenance", environmentURL, apiToken, rids); err != nil {
			fmt.Println(err.Error())
			os.Exit(0)
		}
	}
	if rids, ok := m["dynatrace_frequent_issues"]; ok {
		if err := importFrequentIssues(targetFolder+"/anomalies", environmentURL, apiToken, rids); err != nil {
			fmt.Println(err.Error())
			os.Exit(0)
		}
	}
	if rids, ok := m["dynatrace_ansible_tower_notification"]; ok {
		if err := importNotifications(targetFolder+"/notifications", environmentURL, apiToken, rids, notifications.Types.AnsibleTower, "ansible", "dynatrace_ansible_tower_notification"); err != nil {
			fmt.Println(err.Error())
			os.Exit(0)
		}
	}
	if rids, ok := m["dynatrace_email_notification"]; ok {
		if err := importNotifications(targetFolder+"/notifications", environmentURL, apiToken, rids, notifications.Types.Email, "email", "dynatrace_email_notification"); err != nil {
			fmt.Println(err.Error())
			os.Exit(0)
		}
	}
	if rids, ok := m["dynatrace_jira_notification"]; ok {
		if err := importNotifications(targetFolder+"/notifications", environmentURL, apiToken, rids, notifications.Types.Jira, "jira", "dynatrace_jira_notification"); err != nil {
			fmt.Println(err.Error())
			os.Exit(0)
		}
	}
	if rids, ok := m["dynatrace_ops_genie_notification"]; ok {
		if err := importNotifications(targetFolder+"/notifications", environmentURL, apiToken, rids, notifications.Types.OpsGenie, "opsgenie", "dynatrace_ops_genie_notification"); err != nil {
			fmt.Println(err.Error())
			os.Exit(0)
		}
	}
	if rids, ok := m["dynatrace_pager_duty_notification"]; ok {
		if err := importNotifications(targetFolder+"/notifications", environmentURL, apiToken, rids, notifications.Types.PagerDuty, "pagerduty", "dynatrace_pager_duty_notification"); err != nil {
			fmt.Println(err.Error())
			os.Exit(0)
		}
	}
	if rids, ok := m["dynatrace_service_now_notification"]; ok {
		if err := importNotifications(targetFolder+"/notifications", environmentURL, apiToken, rids, notifications.Types.ServiceNow, "servicenow", "dynatrace_service_now_notification"); err != nil {
			fmt.Println(err.Error())
			os.Exit(0)
		}
	}
	if rids, ok := m["dynatrace_slack_notification"]; ok {
		if err := importNotifications(targetFolder+"/notifications", environmentURL, apiToken, rids, notifications.Types.Slack, "slack", "dynatrace_slack_notification"); err != nil {
			fmt.Println(err.Error())
			os.Exit(0)
		}
	}
	if rids, ok := m["dynatrace_trello_notification"]; ok {
		if err := importNotifications(targetFolder+"/notifications", environmentURL, apiToken, rids, notifications.Types.Slack, "trello", "dynatrace_trello_notification"); err != nil {
			fmt.Println(err.Error())
			os.Exit(0)
		}
	}
	if rids, ok := m["dynatrace_victor_ops_notification"]; ok {
		if err := importNotifications(targetFolder+"/notifications", environmentURL, apiToken, rids, notifications.Types.VictorOps, "victorops", "dynatrace_victor_ops_notification"); err != nil {
			fmt.Println(err.Error())
			os.Exit(0)
		}
	}
	if rids, ok := m["dynatrace_webhook_notification"]; ok {
		if err := importNotifications(targetFolder+"/notifications", environmentURL, apiToken, rids, notifications.Types.WebHook, "webhook", "dynatrace_webhook_notification"); err != nil {
			fmt.Println(err.Error())
			os.Exit(0)
		}
	}
	if rids, ok := m["dynatrace_xmatters_notification"]; ok {
		if err := importNotifications(targetFolder+"/notifications", environmentURL, apiToken, rids, notifications.Types.XMatters, "xmatters", "dynatrace_xmatters_notification"); err != nil {
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
	if exportv2(os.Args) {
		return
	}
	if clusterExport(os.Args) {
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
