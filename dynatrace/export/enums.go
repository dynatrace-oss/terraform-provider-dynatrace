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

package export

import "strings"

type ResourceType string

func (me ResourceType) Trim() string {
	return strings.TrimPrefix(string(me), "dynatrace_")
}

var ResourceTypes = struct {
	AutoTag                  ResourceType
	CustomService            ResourceType
	RequestAttribute         ResourceType
	ApplicationAnomalies     ResourceType
	DatabaseAnomalies        ResourceType
	DiskEventAnomalies       ResourceType
	HostAnomalies            ResourceType
	ServiceAnomalies         ResourceType
	CustomAnomalies          ResourceType
	WebApplication           ResourceType
	MobileApplication        ResourceType
	MaintenanceWindow        ResourceType
	ManagementZone           ResourceType
	SLO                      ResourceType
	SpanAttribute            ResourceType
	SpanCaptureRule          ResourceType
	SpanContextPropagation   ResourceType
	SpanEntryPoint           ResourceType
	ResourceAttributes       ResourceType
	JiraNotification         ResourceType
	WebHookNotification      ResourceType
	AnsibleTowerNotification ResourceType
	EmailNotification        ResourceType
	OpsGenieNotification     ResourceType
	PagerDutyNotification    ResourceType
	ServiceNowNotification   ResourceType
	SlackNotification        ResourceType
	TrelloNotification       ResourceType
	VictorOpsNotification    ResourceType
	XMattersNotification     ResourceType
	Alerting                 ResourceType
	FrequentIssues           ResourceType
	MetricEvents             ResourceType
	APIToken                 ResourceType
	IBMMQFilters             ResourceType
	IMSBridge                ResourceType
	QueueManager             ResourceType
	KeyRequests              ResourceType
	Maintenance              ResourceType
	ManagementZoneV2         ResourceType
	NetworkZones             ResourceType
	AWSCredentials           ResourceType
	AzureCredentials         ResourceType
	CloudFoundryCredentials  ResourceType
	KubernetesCredentials    ResourceType
	Credentials              ResourceType
	Dashboard                ResourceType
	JSONDashboard            ResourceType
	CalculatedServiceMetric  ResourceType
	HostNaming               ResourceType
	ProcessGroupNaming       ResourceType
	ServiceNaming            ResourceType
	NetworkZone              ResourceType
	RequestNaming            ResourceType
	BrowserMonitor           ResourceType
	HTTPMonitor              ResourceType
	DashboardSharing         ResourceType
	ApplicationDetection     ResourceType
	ApplicationErrorRules    ResourceType
	ApplicationDataPrivacy   ResourceType
	SyntheticLocation        ResourceType
	Notification             ResourceType
	QueueSharingGroups       ResourceType
	AlertingProfile          ResourceType
	RequestNamings           ResourceType
	IAMUser                  ResourceType
	IAMGroup                 ResourceType
}{
	"dynatrace_autotag",
	"dynatrace_custom_service",
	"dynatrace_request_attribute",
	"dynatrace_application_anomalies",
	"dynatrace_database_anomalies",
	"dynatrace_disk_anomalies",
	"dynatrace_host_anomalies",
	"dynatrace_service_anomalies",
	"dynatrace_custom_anomalies",
	"dynatrace_web_application",
	"dynatrace_mobile_application",
	"dynatrace_maintenance_window",
	"dynatrace_management_zone",
	"dynatrace_slo",
	"dynatrace_span_attribute",
	"dynatrace_span_capture_rule",
	"dynatrace_span_context_propagation",
	"dynatrace_span_entry_point",
	"dynatrace_resource_attributes",
	"dynatrace_jira_notification",
	"dynatrace_webhook_notification",
	"dynatrace_ansible_tower_notification",
	"dynatrace_email_notification",
	"dynatrace_ops_genie_notification",
	"dynatrace_pager_duty_notification",
	"dynatrace_service_now_notification",
	"dynatrace_slack_notification",
	"dynatrace_trello_notification",
	"dynatrace_victor_ops_notification",
	"dynatrace_xmatters_notification",
	"dynatrace_alerting",
	"dynatrace_frequent_issues",
	"dynatrace_metric_events",
	"dynatrace_api_token",
	"dynatrace_ibm_mq_filters",
	"dynatrace_ims_bridges",
	"dynatrace_queue_manager",
	"dynatrace_key_requests",
	"dynatrace_maintenance",
	"dynatrace_management_zone_v2",
	"dynatrace_network_zones",
	"dynatrace_aws_credentials",
	"dynatrace_azure_credentials",
	"dynatrace_cloudfoundry_credentials",
	"dynatrace_k8s_credentials",
	"dynatrace_credentials",
	"dynatrace_dashboard",
	"dynatrace_json_dashboard",
	"dynatrace_calculated_service_metric",
	"dynatrace_host_naming",
	"dynatrace_processgroup_naming",
	"dynatrace_service_naming",
	"dynatrace_network_zone",
	"dynatrace_request_naming",
	"dynatrace_browser_monitor",
	"dynatrace_http_monitor",
	"dynatrace_dashboard_sharing",
	"dynatrace_application_detection_rule",
	"dynatrace_application_error_rules",
	"dynatrace_application_data_privacy",
	"dynatrace_synthetic_location",
	"dynatrace_notification",
	"dynatrace_queue_sharing_groups",
	"dynatrace_alerting_profile",
	"dynatrace_request_namings",
	"dynatrace_iam_user",
	"dynatrace_iam_group",
}

type ResourceStatus string

func (me ResourceStatus) IsOneOf(stati ...ResourceStatus) bool {
	if len(stati) == 0 {
		return false
	}
	for _, status := range stati {
		if me == status {
			return true
		}
	}
	return false
}

var ResourceStati = struct {
	Downloaded    ResourceStatus
	Erronous      ResourceStatus
	Excluded      ResourceStatus
	Discovered    ResourceStatus
	PostProcessed ResourceStatus
}{
	"Downloaded",
	"Erronous",
	"Excluded",
	"Discovered",
	"PostProcessed",
}

type DataSourceType string

func (me DataSourceType) Trim() string {
	return strings.TrimPrefix(string(me), "dynatrace_")
}

var DataSourceTypes = struct {
	Service          DataSourceType
	AWSIAMExternalID DataSourceType
}{
	"dynatrace_service",
	"dynatrace_aws_iam_external",
}

type ModuleStatus string

func (me ModuleStatus) IsOneOf(stati ...ModuleStatus) bool {
	if len(stati) == 0 {
		return false
	}
	for _, status := range stati {
		if me == status {
			return true
		}
	}
	return false
}

var ModuleStati = struct {
	Untouched  ModuleStatus
	Discovered ModuleStatus
	Downloaded ModuleStatus
	Erronous   ModuleStatus
	Imported   ModuleStatus
}{
	"Untouched",
	"Discovered",
	"Downloaded",
	"Erronous",
	"Imported",
}
