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

func (me ResourceType) AsDataSource() string {
	switch me {
	case ResourceTypes.ManagementZoneV2:
		return "dynatrace_management_zone"
	case ResourceTypes.ManagementZone:
		return "dynatrace_management_zone"
	case ResourceTypes.AlertingProfile:
		return "dynatrace_alerting_profile"
	case ResourceTypes.Alerting:
		return "dynatrace_alerting_profile"
	case ResourceTypes.RequestAttribute:
		return "dynatrace_request_attribute"
	case ResourceTypes.WebApplication:
		return "dynatrace_application"
	case ResourceTypes.MobileApplication:
		return "dynatrace_mobile_application"
	case ResourceTypes.RequestNaming:
		return "dynatrace_request_naming"
	case ResourceTypes.Dashboard:
		return "dynatrace_dashboard"
	case ResourceTypes.JSONDashboard:
		return "dynatrace_dashboard"
	case ResourceTypes.SLO:
		return "dynatrace_slo"
	case ResourceTypes.CalculatedServiceMetric:
		return "dynatrace_calculated_service_metric"
	}
	return ""
}

var ResourceTypes = struct {
	AutoTag                             ResourceType
	AutoTagV2                           ResourceType
	CustomService                       ResourceType
	RequestAttribute                    ResourceType
	ApplicationAnomalies                ResourceType
	DatabaseAnomalies                   ResourceType
	DiskEventAnomalies                  ResourceType
	HostAnomalies                       ResourceType
	ServiceAnomalies                    ResourceType
	CustomAnomalies                     ResourceType
	WebApplication                      ResourceType
	MobileApplication                   ResourceType
	MaintenanceWindow                   ResourceType
	ManagementZone                      ResourceType
	SLO                                 ResourceType
	SLOV2                               ResourceType
	SpanAttribute                       ResourceType
	SpanCaptureRule                     ResourceType
	SpanContextPropagation              ResourceType
	SpanEntryPoint                      ResourceType
	ResourceAttributes                  ResourceType
	JiraNotification                    ResourceType
	WebHookNotification                 ResourceType
	AnsibleTowerNotification            ResourceType
	EmailNotification                   ResourceType
	OpsGenieNotification                ResourceType
	PagerDutyNotification               ResourceType
	ServiceNowNotification              ResourceType
	SlackNotification                   ResourceType
	TrelloNotification                  ResourceType
	VictorOpsNotification               ResourceType
	XMattersNotification                ResourceType
	Alerting                            ResourceType
	FrequentIssues                      ResourceType
	MetricEvents                        ResourceType
	IBMMQFilters                        ResourceType
	IMSBridge                           ResourceType
	QueueManager                        ResourceType
	KeyRequests                         ResourceType
	Maintenance                         ResourceType
	ManagementZoneV2                    ResourceType
	NetworkZones                        ResourceType
	AWSCredentials                      ResourceType
	AzureCredentials                    ResourceType
	CloudFoundryCredentials             ResourceType
	KubernetesCredentials               ResourceType
	Credentials                         ResourceType
	Dashboard                           ResourceType
	JSONDashboard                       ResourceType
	CalculatedServiceMetric             ResourceType
	HostNaming                          ResourceType
	ProcessGroupNaming                  ResourceType
	ServiceNaming                       ResourceType
	NetworkZone                         ResourceType
	RequestNaming                       ResourceType
	BrowserMonitor                      ResourceType
	HTTPMonitor                         ResourceType
	DashboardSharing                    ResourceType
	ApplicationDetection                ResourceType
	ApplicationErrorRules               ResourceType
	ApplicationDataPrivacy              ResourceType
	SyntheticLocation                   ResourceType
	Notification                        ResourceType
	QueueSharingGroups                  ResourceType
	AlertingProfile                     ResourceType
	RequestNamings                      ResourceType
	IAMUser                             ResourceType
	IAMGroup                            ResourceType
	IAMPermission                       ResourceType
	IAMPolicy                           ResourceType
	IAMPolicyBindings                   ResourceType
	ProcessGroupAnomalies               ResourceType
	DDUPool                             ResourceType
	ProcessGroupAlerting                ResourceType
	ServiceAnomaliesV2                  ResourceType
	DatabaseAnomaliesV2                 ResourceType
	ProcessMonitoringRule               ResourceType
	DiskAnomaliesV2                     ResourceType
	DiskSpecificAnomaliesV2             ResourceType
	HostAnomaliesV2                     ResourceType
	CustomAppAnomalies                  ResourceType
	CustomAppCrashRate                  ResourceType
	ProcessMonitoring                   ResourceType
	ProcessAvailability                 ResourceType
	AdvancedProcessGroupDetectionRule   ResourceType
	MobileAppAnomalies                  ResourceType
	MobileAppCrashRate                  ResourceType
	WebAppAnomalies                     ResourceType
	MutedRequests                       ResourceType
	ConnectivityAlerts                  ResourceType
	DeclarativeGrouping                 ResourceType
	HostMonitoring                      ResourceType
	HostProcessGroupMonitoring          ResourceType
	RUMIPLocations                      ResourceType
	CustomAppEnablement                 ResourceType
	MobileAppEnablement                 ResourceType
	WebAppEnablement                    ResourceType
	RUMProcessGroup                     ResourceType
	RUMProviderBreakdown                ResourceType
	UserExperienceScore                 ResourceType
	WebAppResourceCleanup               ResourceType
	UpdateWindows                       ResourceType
	ProcessGroupDetectionFlags          ResourceType
	ProcessGroupMonitoring              ResourceType
	ProcessGroupSimpleDetection         ResourceType
	LogMetrics                          ResourceType
	BrowserMonitorPerformanceThresholds ResourceType
	HttpMonitorPerformanceThresholds    ResourceType
	HttpMonitorCookies                  ResourceType
	SessionReplayWebPrivacy             ResourceType
	SessionReplayResourceCapture        ResourceType
	UsabilityAnalytics                  ResourceType
	SyntheticAvailability               ResourceType
	BrowserMonitorOutageHandling        ResourceType
	HttpMonitorOutageHandling           ResourceType
	CloudAppWorkloadDetection           ResourceType
	MainframeTransactionMonitoring      ResourceType
	MonitoredTechnologiesApache         ResourceType
	MonitoredTechnologiesDotNet         ResourceType
	MonitoredTechnologiesGo             ResourceType
	MonitoredTechnologiesIIS            ResourceType
	MonitoredTechnologiesJava           ResourceType
	MonitoredTechnologiesNGINX          ResourceType
	MonitoredTechnologiesNodeJS         ResourceType
	MonitoredTechnologiesOpenTracing    ResourceType
	MonitoredTechnologiesPHP            ResourceType
	MonitoredTechnologiesVarnish        ResourceType
	MonitoredTechnologiesWSMB           ResourceType
	ProcessVisibility                   ResourceType
	RUMHostHeaders                      ResourceType
	RUMIPDetermination                  ResourceType
	MobileAppRequestErrors              ResourceType
	TransactionStartFilters             ResourceType
	OneAgentFeatures                    ResourceType
	RUMOverloadPrevention               ResourceType
	RUMAdvancedCorrelation              ResourceType
	WebAppBeaconOrigins                 ResourceType
	WebAppResourceTypes                 ResourceType
	GenericTypes                        ResourceType
	GenericRelationships                ResourceType
	SLONormalization                    ResourceType
	DataPrivacy                         ResourceType
	ServiceFailure                      ResourceType
	ServiceHTTPFailure                  ResourceType
	DiskOptions                         ResourceType
	OSServices                          ResourceType
	ExtensionExecutionController        ResourceType
	NetTracerTraffic                    ResourceType
	AIXExtension                        ResourceType
	MetricMetadata                      ResourceType
	MetricQuery                         ResourceType
	ActiveGateToken                     ResourceType
	AuditLog                            ResourceType
	K8sClusterAnomalies                 ResourceType
	K8sNamespaceAnomalies               ResourceType
	K8sNodeAnomalies                    ResourceType
	K8sWorkloadAnomalies                ResourceType
	ContainerBuiltinRule                ResourceType
	ContainerRule                       ResourceType
	ContainerTechnology                 ResourceType
	RemoteEnvironments                  ResourceType
	WebAppCustomErrors                  ResourceType
	WebAppRequestErrors                 ResourceType
	UserSettings                        ResourceType
	DashboardsGeneral                   ResourceType
	DashboardsPresets                   ResourceType
	LogProcessing                       ResourceType
	LogEvents                           ResourceType
	LogTimestamp                        ResourceType
	LogGrail                            ResourceType
	LogCustomAttribute                  ResourceType
	LogSensitiveDataMasking             ResourceType
	LogStorage                          ResourceType
	LogBuckets                          ResourceType
	EULASettings                        ResourceType
	APIDetectionRules                   ResourceType
	ServiceExternalWebRequest           ResourceType
	ServiceExternalWebService           ResourceType
	ServiceFullWebRequest               ResourceType
	ServiceFullWebService               ResourceType
	DashboardsAllowlist                 ResourceType
	FailureDetectionParameters          ResourceType
	FailureDetectionRules               ResourceType
	LogOneAgent                         ResourceType
	IssueTracking                       ResourceType
	GeolocationSettings                 ResourceType
	UserSessionCustomMetrics            ResourceType
	CustomUnits                         ResourceType
	DiskAnalytics                       ResourceType
	NetworkTraffic                      ResourceType
	TokenSettings                       ResourceType
	ExtensionExecutionRemote            ResourceType
	K8sPVCAnomalies                     ResourceType
	UserActionCustomMetrics             ResourceType
	WebAppJavascriptVersion             ResourceType
	WebAppJavascriptUpdates             ResourceType
	OpenTelemetryMetrics                ResourceType
	ActiveGateUpdates                   ResourceType
	OneAgentDefaultVersion              ResourceType
	OneAgentUpdates                     ResourceType
	OwnershipTeams                      ResourceType
	LogCustomSource                     ResourceType
	ApplicationDetectionV2              ResourceType
	Kubernetes                          ResourceType
	CloudFoundry                        ResourceType
	DiskAnomalyDetectionRules           ResourceType
	AWSAnomalies                        ResourceType
	VMwareAnomalies                     ResourceType
	BusinessEventsOneAgent              ResourceType
	BusinessEventsBuckets               ResourceType
	BusinessEventsMetrics               ResourceType
	BusinessEventsProcessing            ResourceType
	WebAppKeyPerformanceCustom          ResourceType
	WebAppKeyPerformanceLoad            ResourceType
	WebAppKeyPerformanceXHR             ResourceType
}{
	"dynatrace_autotag",
	"dynatrace_autotag_v2",
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
	"dynatrace_slo_v2",
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
	"dynatrace_iam_permission",
	"dynatrace_iam_policy",
	"dynatrace_iam_policy_bindings",
	"dynatrace_pg_anomalies",
	"dynatrace_ddu_pool",
	"dynatrace_pg_alerting",
	"dynatrace_service_anomalies_v2",
	"dynatrace_database_anomalies_v2",
	"dynatrace_process_monitoring_rule",
	"dynatrace_disk_anomalies_v2",
	"dynatrace_disk_specific_anomalies_v2",
	"dynatrace_host_anomalies_v2",
	"dynatrace_custom_app_anomalies",
	"dynatrace_custom_app_crash_rate",
	"dynatrace_process_monitoring",
	"dynatrace_process_availability",
	"dynatrace_process_group_detection",
	"dynatrace_mobile_app_anomalies",
	"dynatrace_mobile_app_crash_rate",
	"dynatrace_web_app_anomalies",
	"dynatrace_muted_requests",
	"dynatrace_connectivity_alerts",
	"dynatrace_declarative_grouping",
	"dynatrace_host_monitoring",
	"dynatrace_host_process_group_monitoring",
	"dynatrace_rum_ip_locations",
	"dynatrace_custom_app_enablement",
	"dynatrace_mobile_app_enablement",
	"dynatrace_web_app_enablement",
	"dynatrace_process_group_rum",
	"dynatrace_rum_provider_breakdown",
	"dynatrace_user_experience_score",
	"dynatrace_web_app_resource_cleanup",
	"dynatrace_update_windows",
	"dynatrace_process_group_detection_flags",
	"dynatrace_process_group_monitoring",
	"dynatrace_process_group_simple_detection",
	"dynatrace_log_metrics",
	"dynatrace_browser_monitor_performance",
	"dynatrace_http_monitor_performance",
	"dynatrace_http_monitor_cookies",
	"dynatrace_session_replay_web_privacy",
	"dynatrace_session_replay_resource_capture",
	"dynatrace_usability_analytics",
	"dynatrace_synthetic_availability",
	"dynatrace_browser_monitor_outage",
	"dynatrace_http_monitor_outage",
	"dynatrace_cloudapp_workloaddetection",
	"dynatrace_mainframe_transaction_monitoring",
	"dynatrace_monitored_technologies_apache",
	"dynatrace_monitored_technologies_dotnet",
	"dynatrace_monitored_technologies_go",
	"dynatrace_monitored_technologies_iis",
	"dynatrace_monitored_technologies_java",
	"dynatrace_monitored_technologies_nginx",
	"dynatrace_monitored_technologies_nodejs",
	"dynatrace_monitored_technologies_opentracing",
	"dynatrace_monitored_technologies_php",
	"dynatrace_monitored_technologies_varnish",
	"dynatrace_monitored_technologies_wsmb",
	"dynatrace_process_visibility",
	"dynatrace_rum_host_headers",
	"dynatrace_rum_ip_determination",
	"dynatrace_mobile_app_request_errors",
	"dynatrace_transaction_start_filters",
	"dynatrace_oneagent_features",
	"dynatrace_rum_overload_prevention",
	"dynatrace_rum_advanced_correlation",
	"dynatrace_web_app_beacon_origins",
	"dynatrace_web_app_resource_types",
	"dynatrace_generic_types",
	"dynatrace_generic_relationships",
	"dynatrace_slo_normalization",
	"dynatrace_data_privacy",
	"dynatrace_service_failure",
	"dynatrace_service_http_failure",
	"dynatrace_disk_options",
	"dynatrace_os_services",
	"dynatrace_extension_execution_controller",
	"dynatrace_nettracer",
	"dynatrace_aix_extension",
	"dynatrace_metric_metadata",
	"dynatrace_metric_query",
	"dynatrace_activegate_token",
	"dynatrace_audit_log",
	"dynatrace_k8s_cluster_anomalies",
	"dynatrace_k8s_namespace_anomalies",
	"dynatrace_k8s_node_anomalies",
	"dynatrace_k8s_workload_anomalies",
	"dynatrace_container_builtin_rule",
	"dynatrace_container_rule",
	"dynatrace_container_technology",
	"dynatrace_remote_environments",
	"dynatrace_web_app_custom_errors",
	"dynatrace_web_app_request_errors",
	"dynatrace_user_settings",
	"dynatrace_dashboards_general",
	"dynatrace_dashboards_presets",
	"dynatrace_log_processing",
	"dynatrace_log_events",
	"dynatrace_log_timestamp",
	"dynatrace_log_grail",
	"dynatrace_log_custom_attribute",
	"dynatrace_log_sensitive_data_masking",
	"dynatrace_log_storage",
	"dynatrace_log_buckets",
	"dynatrace_eula_settings",
	"dynatrace_api_detection",
	"dynatrace_service_external_web_request",
	"dynatrace_service_external_web_service",
	"dynatrace_service_full_web_request",
	"dynatrace_service_full_web_service",
	"dynatrace_dashboards_allowlist",
	"dynatrace_failure_detection_parameters",
	"dynatrace_failure_detection_rules",
	"dynatrace_log_oneagent",
	"dynatrace_issue_tracking",
	"dynatrace_geolocation",
	"dynatrace_user_session_metrics",
	"dynatrace_custom_units",
	"dynatrace_disk_analytics",
	"dynatrace_network_traffic",
	"dynatrace_token_settings",
	"dynatrace_extension_execution_remote",
	"dynatrace_k8s_pvc_anomalies",
	"dynatrace_user_action_metrics",
	"dynatrace_web_app_javascript_version",
	"dynatrace_web_app_javascript_updates",
	"dynatrace_opentelemetry_metrics",
	"dynatrace_activegate_updates",
	"dynatrace_oneagent_default_version",
	"dynatrace_oneagent_updates",
	"dynatrace_ownership_teams",
	"dynatrace_log_custom_source",
	"dynatrace_application_detection_rule_v2",
	"dynatrace_kubernetes",
	"dynatrace_cloud_foundry",
	"dynatrace_disk_anomaly_rules",
	"dynatrace_aws_anomalies",
	"dynatrace_vmware_anomalies",
	"dynatrace_business_events_oneagent",
	"dynatrace_business_events_buckets",
	"dynatrace_business_events_metrics",
	"dynatrace_business_events_processing",
	"dynatrace_web_app_key_performance_custom",
	"dynatrace_web_app_key_performance_load",
	"dynatrace_web_app_key_performance_xhr",
}

func (me ResourceType) GetChildren() []ResourceType {
	res := []ResourceType{}
	for k, v := range AllResources {
		if v.Parent != nil && (string(*v.Parent) == string(me)) {
			res = append(res, k)
		}
	}
	return res
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
