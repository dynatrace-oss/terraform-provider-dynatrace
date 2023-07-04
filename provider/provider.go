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

package provider

import (
	"context"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/datasources/alerting"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/datasources/application"
	dsaws "github.com/dynatrace-oss/terraform-provider-dynatrace/datasources/aws"
	aws_credentials "github.com/dynatrace-oss/terraform-provider-dynatrace/datasources/credentials/aws"
	aws_supported_services "github.com/dynatrace-oss/terraform-provider-dynatrace/datasources/credentials/aws/supported_services"
	azure_credentials "github.com/dynatrace-oss/terraform-provider-dynatrace/datasources/credentials/azure"
	azure_supported_services "github.com/dynatrace-oss/terraform-provider-dynatrace/datasources/credentials/azure/supported_services"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/datasources/credentials/vault"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/datasources/dashboard"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/datasources/entities"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/datasources/entity"
	failure_detection_parameters "github.com/dynatrace-oss/terraform-provider-dynatrace/datasources/failuredetection/parameters"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/datasources/host"
	ds_iam_groups "github.com/dynatrace-oss/terraform-provider-dynatrace/datasources/iam/groups"
	ds_iam_users "github.com/dynatrace-oss/terraform-provider-dynatrace/datasources/iam/users"
	metricsds "github.com/dynatrace-oss/terraform-provider-dynatrace/datasources/metrics/calculated/service"
	mgmzds "github.com/dynatrace-oss/terraform-provider-dynatrace/datasources/mgmz"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/datasources/mobileapplication"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/datasources/process"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/datasources/processgroup"
	reqattrds "github.com/dynatrace-oss/terraform-provider-dynatrace/datasources/requestattributes"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/datasources/requestnaming"
	serviceds "github.com/dynatrace-oss/terraform-provider-dynatrace/datasources/service"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/datasources/slo"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/datasources/synthetic/locations"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/datasources/synthetic/nodes"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/datasources/tenant"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/datasources/updatewindows"
	v2alerting "github.com/dynatrace-oss/terraform-provider-dynatrace/datasources/v2alerting"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/export"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/provider/config"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/provider/logging"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/resources"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/resources/apitokens"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/resources/bindings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/resources/customtags"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/resources/environments"
	mgmzperm "github.com/dynatrace-oss/terraform-provider-dynatrace/resources/permissions/mgmz"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/resources/policies"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/resources/usergroups"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/resources/users"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ResourceSpecification has no documentation
type ResourceSpecification interface {
	Resource() *schema.Resource
	Create(context.Context, *schema.ResourceData, any) diag.Diagnostics
	Update(context.Context, *schema.ResourceData, any) diag.Diagnostics
	Read(context.Context, *schema.ResourceData, any) diag.Diagnostics
	Delete(context.Context, *schema.ResourceData, any) diag.Diagnostics
}

// Provider function for Dynatrace API
func Provider() *schema.Provider {
	logging.SetOutput()
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"dt_env_url": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{"DYNATRACE_ENV_URL", "DT_ENV_URL"}, nil),
			},
			"dt_api_token": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{"DYNATRACE_API_TOKEN", "DT_API_TOKEN"}, nil),
			},
			"dt_cluster_api_token": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{"DYNATRACE_CLUSTER_API_TOKEN", "DT_CLUSTER_API_TOKEN"}, nil),
			},
			"dt_cluster_url": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{"DYNATRACE_CLUSTER_URL", "DT_CLUSTER_URL"}, nil),
			},
			"iam_client_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{"IAM_CLIENT_ID", "DT_CLIENT_ID"}, nil),
			},
			"iam_account_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{"IAM_ACCOUNT_ID", "DT_ACCOUNT_ID"}, nil),
			},
			"iam_client_secret": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{"IAM_CLIENT_SECRET", "DT_CLIENT_SECRET"}, nil),
			},
		},
		DataSourcesMap: map[string]*schema.Resource{
			"dynatrace_alerting_profiles":            alerting.DataSource(),
			"dynatrace_alerting_profile":             v2alerting.DataSource(),
			"dynatrace_credentials":                  vault.DataSource(),
			"dynatrace_synthetic_locations":          locations.DataSource(),
			"dynatrace_synthetic_location":           locations.UniqueDataSource(),
			"dynatrace_service":                      serviceds.DataSource(),
			"dynatrace_management_zone":              mgmzds.DataSource(),
			"dynatrace_management_zones":             mgmzds.DataSourceMultiple(),
			"dynatrace_application":                  application.DataSource(),
			"dynatrace_mobile_application":           mobileapplication.DataSource(),
			"dynatrace_host":                         host.DataSource(),
			"dynatrace_process":                      process.DataSource(),
			"dynatrace_process_group":                processgroup.DataSource(),
			"dynatrace_aws_iam_external":             dsaws.DataSource(),
			"dynatrace_request_attribute":            reqattrds.DataSource(),
			"dynatrace_calculated_service_metric":    metricsds.DataSource(),
			"dynatrace_iam_group":                    ds_iam_groups.DataSource(),
			"dynatrace_entity":                       entity.DataSource(),
			"dynatrace_entities":                     entities.DataSource(),
			"dynatrace_iam_user":                     ds_iam_users.DataSource(),
			"dynatrace_request_naming":               requestnaming.DataSource(),
			"dynatrace_dashboard":                    dashboard.DataSource(),
			"dynatrace_slo":                          slo.DataSource(),
			"dynatrace_azure_supported_services":     azure_supported_services.DataSource(),
			"dynatrace_aws_supported_services":       aws_supported_services.DataSource(),
			"dynatrace_failure_detection_parameters": failure_detection_parameters.DataSource(),
			"dynatrace_update_windows":               updatewindows.DataSource(),
			"dynatrace_aws_credentials":              aws_credentials.DataSource(),
			"dynatrace_azure_credentials":            azure_credentials.DataSource(),
			"dynatrace_synthetic_nodes":              nodes.DataSource(),
			"dynatrace_tenant":                       tenant.DataSource(),
		},
		ResourcesMap: map[string]*schema.Resource{
			"dynatrace_custom_service":                     resources.NewGeneric(export.ResourceTypes.CustomService).Resource(),
			"dynatrace_dashboard":                          resources.NewGeneric(export.ResourceTypes.Dashboard).Resource(),
			"dynatrace_json_dashboard":                     resources.NewGeneric(export.ResourceTypes.JSONDashboard).Resource(),
			"dynatrace_management_zone":                    resources.NewGeneric(export.ResourceTypes.ManagementZone).Resource(),
			"dynatrace_management_zone_v2":                 resources.NewGeneric(export.ResourceTypes.ManagementZoneV2).Resource(),
			"dynatrace_maintenance_window":                 resources.NewGeneric(export.ResourceTypes.MaintenanceWindow).Resource(),
			"dynatrace_maintenance":                        resources.NewGeneric(export.ResourceTypes.Maintenance).Resource(),
			"dynatrace_request_attribute":                  resources.NewGeneric(export.ResourceTypes.RequestAttribute).Resource(),
			"dynatrace_alerting_profile":                   resources.NewGeneric(export.ResourceTypes.AlertingProfile).Resource(),
			"dynatrace_alerting":                           resources.NewGeneric(export.ResourceTypes.Alerting).Resource(),
			"dynatrace_notification":                       resources.NewGeneric(export.ResourceTypes.Notification).Resource(),
			"dynatrace_autotag":                            resources.NewGeneric(export.ResourceTypes.AutoTag).Resource(),
			"dynatrace_aws_credentials":                    resources.NewGeneric(export.ResourceTypes.AWSCredentials).Resource(),
			"dynatrace_aws_service":                        resources.NewGeneric(export.ResourceTypes.AWSService).Resource(),
			"dynatrace_azure_credentials":                  resources.NewGeneric(export.ResourceTypes.AzureCredentials).Resource(),
			"dynatrace_azure_service":                      resources.NewGeneric(export.ResourceTypes.AzureService).Resource(),
			"dynatrace_k8s_credentials":                    resources.NewGeneric(export.ResourceTypes.KubernetesCredentials).Resource(),
			"dynatrace_cloudfoundry_credentials":           resources.NewGeneric(export.ResourceTypes.CloudFoundryCredentials).Resource(),
			"dynatrace_service_anomalies":                  resources.NewGeneric(export.ResourceTypes.ServiceAnomalies).Resource(),
			"dynatrace_service_anomalies_v2":               resources.NewGeneric(export.ResourceTypes.ServiceAnomaliesV2).Resource(),
			"dynatrace_application_anomalies":              resources.NewGeneric(export.ResourceTypes.ApplicationAnomalies).Resource(),
			"dynatrace_host_anomalies":                     resources.NewGeneric(export.ResourceTypes.HostAnomalies).Resource(),
			"dynatrace_host_anomalies_v2":                  resources.NewGeneric(export.ResourceTypes.HostAnomaliesV2).Resource(),
			"dynatrace_database_anomalies":                 resources.NewGeneric(export.ResourceTypes.DatabaseAnomalies).Resource(),
			"dynatrace_database_anomalies_v2":              resources.NewGeneric(export.ResourceTypes.DatabaseAnomaliesV2).Resource(),
			"dynatrace_custom_anomalies":                   resources.NewGeneric(export.ResourceTypes.CustomAnomalies).Resource(),
			"dynatrace_metric_events":                      resources.NewGeneric(export.ResourceTypes.MetricEvents).Resource(),
			"dynatrace_disk_anomalies":                     resources.NewGeneric(export.ResourceTypes.DiskEventAnomalies).Resource(),
			"dynatrace_disk_anomalies_v2":                  resources.NewGeneric(export.ResourceTypes.DiskAnomaliesV2).Resource(),
			"dynatrace_disk_specific_anomalies_v2":         resources.NewGeneric(export.ResourceTypes.DiskSpecificAnomaliesV2).Resource(),
			"dynatrace_calculated_service_metric":          resources.NewGeneric(export.ResourceTypes.CalculatedServiceMetric).Resource(),
			"dynatrace_service_naming":                     resources.NewGeneric(export.ResourceTypes.ServiceNaming).Resource(),
			"dynatrace_host_naming":                        resources.NewGeneric(export.ResourceTypes.HostNaming).Resource(),
			"dynatrace_processgroup_naming":                resources.NewGeneric(export.ResourceTypes.ProcessGroupNaming).Resource(),
			"dynatrace_slo":                                resources.NewGeneric(export.ResourceTypes.SLO).Resource(),
			"dynatrace_span_entry_point":                   resources.NewGeneric(export.ResourceTypes.SpanEntryPoint).Resource(),
			"dynatrace_span_capture_rule":                  resources.NewGeneric(export.ResourceTypes.SpanCaptureRule).Resource(),
			"dynatrace_span_context_propagation":           resources.NewGeneric(export.ResourceTypes.SpanContextPropagation).Resource(),
			"dynatrace_resource_attributes":                resources.NewGeneric(export.ResourceTypes.ResourceAttributes).Resource(),
			"dynatrace_span_attribute":                     resources.NewGeneric(export.ResourceTypes.SpanAttribute).Resource(),
			"dynatrace_dashboard_sharing":                  resources.NewGeneric(export.ResourceTypes.DashboardSharing).Resource(),
			"dynatrace_environment":                        environments.Resource(),
			"dynatrace_mobile_application":                 resources.NewGeneric(export.ResourceTypes.MobileApplication).Resource(),
			"dynatrace_browser_monitor":                    resources.NewGeneric(export.ResourceTypes.BrowserMonitor).Resource(),
			"dynatrace_http_monitor":                       resources.NewGeneric(export.ResourceTypes.HTTPMonitor).Resource(),
			"dynatrace_web_application":                    resources.NewGeneric(export.ResourceTypes.WebApplication).Resource(),
			"dynatrace_application_data_privacy":           resources.NewGeneric(export.ResourceTypes.ApplicationDataPrivacy).Resource(),
			"dynatrace_application_error_rules":            resources.NewGeneric(export.ResourceTypes.ApplicationErrorRules).Resource(),
			"dynatrace_request_naming":                     resources.NewGeneric(export.ResourceTypes.RequestNaming).Resource(),
			"dynatrace_request_namings":                    resources.NewGeneric(export.ResourceTypes.RequestNamings).Resource(),
			"dynatrace_user_group":                         usergroups.Resource(),
			"dynatrace_user":                               users.Resource(),
			"dynatrace_policy":                             policies.Resource(),
			"dynatrace_policy_bindings":                    bindings.Resource(),
			"dynatrace_mgmz_permission":                    mgmzperm.Resource(),
			"dynatrace_key_requests":                       resources.NewGeneric(export.ResourceTypes.KeyRequests).Resource(),
			"dynatrace_queue_manager":                      resources.NewGeneric(export.ResourceTypes.QueueManager).Resource(),
			"dynatrace_ibm_mq_filters":                     resources.NewGeneric(export.ResourceTypes.IBMMQFilters).Resource(),
			"dynatrace_queue_sharing_groups":               resources.NewGeneric(export.ResourceTypes.QueueSharingGroups).Resource(),
			"dynatrace_ims_bridges":                        resources.NewGeneric(export.ResourceTypes.IMSBridge).Resource(),
			"dynatrace_network_zones":                      resources.NewGeneric(export.ResourceTypes.NetworkZones).Resource(),
			"dynatrace_application_detection_rule":         resources.NewGeneric(export.ResourceTypes.ApplicationDetection).Resource(),
			"dynatrace_frequent_issues":                    resources.NewGeneric(export.ResourceTypes.FrequentIssues).Resource(),
			"dynatrace_ansible_tower_notification":         resources.NewGeneric(export.ResourceTypes.AnsibleTowerNotification).Resource(),
			"dynatrace_email_notification":                 resources.NewGeneric(export.ResourceTypes.EmailNotification).Resource(),
			"dynatrace_jira_notification":                  resources.NewGeneric(export.ResourceTypes.JiraNotification).Resource(),
			"dynatrace_ops_genie_notification":             resources.NewGeneric(export.ResourceTypes.OpsGenieNotification).Resource(),
			"dynatrace_pager_duty_notification":            resources.NewGeneric(export.ResourceTypes.PagerDutyNotification).Resource(),
			"dynatrace_service_now_notification":           resources.NewGeneric(export.ResourceTypes.ServiceNowNotification).Resource(),
			"dynatrace_slack_notification":                 resources.NewGeneric(export.ResourceTypes.SlackNotification).Resource(),
			"dynatrace_trello_notification":                resources.NewGeneric(export.ResourceTypes.TrelloNotification).Resource(),
			"dynatrace_victor_ops_notification":            resources.NewGeneric(export.ResourceTypes.VictorOpsNotification).Resource(),
			"dynatrace_webhook_notification":               resources.NewGeneric(export.ResourceTypes.WebHookNotification).Resource(),
			"dynatrace_xmatters_notification":              resources.NewGeneric(export.ResourceTypes.XMattersNotification).Resource(),
			"dynatrace_credentials":                        resources.NewGeneric(export.ResourceTypes.Credentials).Resource(),
			"dynatrace_synthetic_location":                 resources.NewGeneric(export.ResourceTypes.SyntheticLocation).Resource(),
			"dynatrace_network_zone":                       resources.NewGeneric(export.ResourceTypes.NetworkZone).Resource(),
			"dynatrace_iam_user":                           resources.NewGeneric(export.ResourceTypes.IAMUser).Resource(),
			"dynatrace_iam_group":                          resources.NewGeneric(export.ResourceTypes.IAMGroup).Resource(),
			"dynatrace_iam_permission":                     resources.NewGeneric(export.ResourceTypes.IAMPermission).Resource(),
			"dynatrace_iam_policy":                         resources.NewGeneric(export.ResourceTypes.IAMPolicy).Resource(),
			"dynatrace_iam_policy_bindings":                resources.NewGeneric(export.ResourceTypes.IAMPolicyBindings).Resource(),
			"dynatrace_api_token":                          apitokens.Resource(),
			"dynatrace_custom_tags":                        customtags.Resource(),
			"dynatrace_pg_anomalies":                       resources.NewGeneric(export.ResourceTypes.ProcessGroupAnomalies).Resource(),
			"dynatrace_ddu_pool":                           resources.NewGeneric(export.ResourceTypes.DDUPool).Resource(),
			"dynatrace_pg_alerting":                        resources.NewGeneric(export.ResourceTypes.ProcessGroupAlerting).Resource(),
			"dynatrace_process_monitoring_rule":            resources.NewGeneric(export.ResourceTypes.ProcessMonitoringRule).Resource(),
			"dynatrace_custom_app_anomalies":               resources.NewGeneric(export.ResourceTypes.CustomAppAnomalies).Resource(),
			"dynatrace_custom_app_crash_rate":              resources.NewGeneric(export.ResourceTypes.CustomAppCrashRate).Resource(),
			"dynatrace_process_monitoring":                 resources.NewGeneric(export.ResourceTypes.ProcessMonitoring).Resource(),
			"dynatrace_process_availability":               resources.NewGeneric(export.ResourceTypes.ProcessAvailability).Resource(),
			"dynatrace_process_group_detection":            resources.NewGeneric(export.ResourceTypes.AdvancedProcessGroupDetectionRule).Resource(),
			"dynatrace_mobile_app_anomalies":               resources.NewGeneric(export.ResourceTypes.MobileAppAnomalies).Resource(),
			"dynatrace_mobile_app_crash_rate":              resources.NewGeneric(export.ResourceTypes.MobileAppCrashRate).Resource(),
			"dynatrace_web_app_anomalies":                  resources.NewGeneric(export.ResourceTypes.WebAppAnomalies).Resource(),
			"dynatrace_muted_requests":                     resources.NewGeneric(export.ResourceTypes.MutedRequests).Resource(),
			"dynatrace_connectivity_alerts":                resources.NewGeneric(export.ResourceTypes.ConnectivityAlerts).Resource(),
			"dynatrace_declarative_grouping":               resources.NewGeneric(export.ResourceTypes.DeclarativeGrouping).Resource(),
			"dynatrace_host_monitoring":                    resources.NewGeneric(export.ResourceTypes.HostMonitoring).Resource(),
			"dynatrace_host_process_group_monitoring":      resources.NewGeneric(export.ResourceTypes.HostProcessGroupMonitoring).Resource(),
			"dynatrace_rum_ip_locations":                   resources.NewGeneric(export.ResourceTypes.RUMIPLocations).Resource(),
			"dynatrace_custom_app_enablement":              resources.NewGeneric(export.ResourceTypes.CustomAppEnablement).Resource(),
			"dynatrace_mobile_app_enablement":              resources.NewGeneric(export.ResourceTypes.MobileAppEnablement).Resource(),
			"dynatrace_web_app_enablement":                 resources.NewGeneric(export.ResourceTypes.WebAppEnablement).Resource(),
			"dynatrace_process_group_rum":                  resources.NewGeneric(export.ResourceTypes.RUMProcessGroup).Resource(),
			"dynatrace_rum_provider_breakdown":             resources.NewGeneric(export.ResourceTypes.RUMProviderBreakdown).Resource(),
			"dynatrace_user_experience_score":              resources.NewGeneric(export.ResourceTypes.UserExperienceScore).Resource(),
			"dynatrace_web_app_resource_cleanup":           resources.NewGeneric(export.ResourceTypes.WebAppResourceCleanup).Resource(),
			"dynatrace_update_windows":                     resources.NewGeneric(export.ResourceTypes.UpdateWindows).Resource(),
			"dynatrace_process_group_detection_flags":      resources.NewGeneric(export.ResourceTypes.ProcessGroupDetectionFlags).Resource(),
			"dynatrace_process_group_monitoring":           resources.NewGeneric(export.ResourceTypes.ProcessGroupMonitoring).Resource(),
			"dynatrace_process_group_simple_detection":     resources.NewGeneric(export.ResourceTypes.ProcessGroupSimpleDetection).Resource(),
			"dynatrace_log_metrics":                        resources.NewGeneric(export.ResourceTypes.LogMetrics).Resource(),
			"dynatrace_browser_monitor_performance":        resources.NewGeneric(export.ResourceTypes.BrowserMonitorPerformanceThresholds).Resource(),
			"dynatrace_http_monitor_performance":           resources.NewGeneric(export.ResourceTypes.HttpMonitorPerformanceThresholds).Resource(),
			"dynatrace_http_monitor_cookies":               resources.NewGeneric(export.ResourceTypes.HttpMonitorCookies).Resource(),
			"dynatrace_session_replay_web_privacy":         resources.NewGeneric(export.ResourceTypes.SessionReplayWebPrivacy).Resource(),
			"dynatrace_session_replay_resource_capture":    resources.NewGeneric(export.ResourceTypes.SessionReplayResourceCapture).Resource(),
			"dynatrace_usability_analytics":                resources.NewGeneric(export.ResourceTypes.UsabilityAnalytics).Resource(),
			"dynatrace_synthetic_availability":             resources.NewGeneric(export.ResourceTypes.SyntheticAvailability).Resource(),
			"dynatrace_browser_monitor_outage":             resources.NewGeneric(export.ResourceTypes.BrowserMonitorOutageHandling).Resource(),
			"dynatrace_http_monitor_outage":                resources.NewGeneric(export.ResourceTypes.HttpMonitorOutageHandling).Resource(),
			"dynatrace_cloudapp_workloaddetection":         resources.NewGeneric(export.ResourceTypes.CloudAppWorkloadDetection).Resource(),
			"dynatrace_mainframe_transaction_monitoring":   resources.NewGeneric(export.ResourceTypes.MainframeTransactionMonitoring).Resource(),
			"dynatrace_monitored_technologies_apache":      resources.NewGeneric(export.ResourceTypes.MonitoredTechnologiesApache).Resource(),
			"dynatrace_monitored_technologies_dotnet":      resources.NewGeneric(export.ResourceTypes.MonitoredTechnologiesDotNet).Resource(),
			"dynatrace_monitored_technologies_go":          resources.NewGeneric(export.ResourceTypes.MonitoredTechnologiesGo).Resource(),
			"dynatrace_monitored_technologies_iis":         resources.NewGeneric(export.ResourceTypes.MonitoredTechnologiesIIS).Resource(),
			"dynatrace_monitored_technologies_java":        resources.NewGeneric(export.ResourceTypes.MonitoredTechnologiesJava).Resource(),
			"dynatrace_monitored_technologies_nginx":       resources.NewGeneric(export.ResourceTypes.MonitoredTechnologiesNGINX).Resource(),
			"dynatrace_monitored_technologies_nodejs":      resources.NewGeneric(export.ResourceTypes.MonitoredTechnologiesNodeJS).Resource(),
			"dynatrace_monitored_technologies_opentracing": resources.NewGeneric(export.ResourceTypes.MonitoredTechnologiesOpenTracing).Resource(),
			"dynatrace_monitored_technologies_php":         resources.NewGeneric(export.ResourceTypes.MonitoredTechnologiesPHP).Resource(),
			"dynatrace_monitored_technologies_varnish":     resources.NewGeneric(export.ResourceTypes.MonitoredTechnologiesVarnish).Resource(),
			"dynatrace_monitored_technologies_wsmb":        resources.NewGeneric(export.ResourceTypes.MonitoredTechnologiesWSMB).Resource(),
			"dynatrace_process_visibility":                 resources.NewGeneric(export.ResourceTypes.ProcessVisibility).Resource(),
			"dynatrace_rum_host_headers":                   resources.NewGeneric(export.ResourceTypes.RUMHostHeaders).Resource(),
			"dynatrace_rum_ip_determination":               resources.NewGeneric(export.ResourceTypes.RUMIPDetermination).Resource(),
			"dynatrace_mobile_app_request_errors":          resources.NewGeneric(export.ResourceTypes.MobileAppRequestErrors).Resource(),
			"dynatrace_transaction_start_filters":          resources.NewGeneric(export.ResourceTypes.TransactionStartFilters).Resource(),
			"dynatrace_oneagent_features":                  resources.NewGeneric(export.ResourceTypes.OneAgentFeatures).Resource(),
			"dynatrace_rum_overload_prevention":            resources.NewGeneric(export.ResourceTypes.RUMOverloadPrevention).Resource(),
			"dynatrace_rum_advanced_correlation":           resources.NewGeneric(export.ResourceTypes.RUMAdvancedCorrelation).Resource(),
			"dynatrace_web_app_beacon_origins":             resources.NewGeneric(export.ResourceTypes.WebAppBeaconOrigins).Resource(),
			"dynatrace_web_app_resource_types":             resources.NewGeneric(export.ResourceTypes.WebAppResourceTypes).Resource(),
			"dynatrace_generic_types":                      resources.NewGeneric(export.ResourceTypes.GenericTypes).Resource(),
			"dynatrace_generic_relationships":              resources.NewGeneric(export.ResourceTypes.GenericRelationships).Resource(),
			"dynatrace_slo_normalization":                  resources.NewGeneric(export.ResourceTypes.SLONormalization).Resource(),
			"dynatrace_data_privacy":                       resources.NewGeneric(export.ResourceTypes.DataPrivacy).Resource(),
			"dynatrace_service_failure":                    resources.NewGeneric(export.ResourceTypes.ServiceFailure).Resource(),
			"dynatrace_service_http_failure":               resources.NewGeneric(export.ResourceTypes.ServiceHTTPFailure).Resource(),
			"dynatrace_disk_options":                       resources.NewGeneric(export.ResourceTypes.DiskOptions).Resource(),
			"dynatrace_os_services":                        resources.NewGeneric(export.ResourceTypes.OSServices).Resource(),
			"dynatrace_extension_execution_controller":     resources.NewGeneric(export.ResourceTypes.ExtensionExecutionController).Resource(),
			"dynatrace_nettracer":                          resources.NewGeneric(export.ResourceTypes.NetTracerTraffic).Resource(),
			"dynatrace_aix_extension":                      resources.NewGeneric(export.ResourceTypes.AIXExtension).Resource(),
			"dynatrace_metric_metadata":                    resources.NewGeneric(export.ResourceTypes.MetricMetadata).Resource(),
			"dynatrace_metric_query":                       resources.NewGeneric(export.ResourceTypes.MetricQuery).Resource(),
			"dynatrace_activegate_token":                   resources.NewGeneric(export.ResourceTypes.ActiveGateToken).Resource(),
			"dynatrace_audit_log":                          resources.NewGeneric(export.ResourceTypes.AuditLog).Resource(),
			"dynatrace_k8s_cluster_anomalies":              resources.NewGeneric(export.ResourceTypes.K8sClusterAnomalies).Resource(),
			"dynatrace_k8s_namespace_anomalies":            resources.NewGeneric(export.ResourceTypes.K8sNamespaceAnomalies).Resource(),
			"dynatrace_k8s_node_anomalies":                 resources.NewGeneric(export.ResourceTypes.K8sNodeAnomalies).Resource(),
			"dynatrace_k8s_workload_anomalies":             resources.NewGeneric(export.ResourceTypes.K8sWorkloadAnomalies).Resource(),
			"dynatrace_container_builtin_rule":             resources.NewGeneric(export.ResourceTypes.ContainerBuiltinRule).Resource(),
			"dynatrace_container_rule":                     resources.NewGeneric(export.ResourceTypes.ContainerRule).Resource(),
			"dynatrace_container_technology":               resources.NewGeneric(export.ResourceTypes.ContainerTechnology).Resource(),
			"dynatrace_remote_environments":                resources.NewGeneric(export.ResourceTypes.RemoteEnvironments).Resource(),
			"dynatrace_web_app_custom_errors":              resources.NewGeneric(export.ResourceTypes.WebAppCustomErrors).Resource(),
			"dynatrace_web_app_request_errors":             resources.NewGeneric(export.ResourceTypes.WebAppRequestErrors).Resource(),
			"dynatrace_user_settings":                      resources.NewGeneric(export.ResourceTypes.UserSettings).Resource(),
			"dynatrace_dashboards_general":                 resources.NewGeneric(export.ResourceTypes.DashboardsGeneral).Resource(),
			"dynatrace_dashboards_presets":                 resources.NewGeneric(export.ResourceTypes.DashboardsPresets).Resource(),
			"dynatrace_log_processing":                     resources.NewGeneric(export.ResourceTypes.LogProcessing).Resource(),
			"dynatrace_log_events":                         resources.NewGeneric(export.ResourceTypes.LogEvents).Resource(),
			"dynatrace_log_timestamp":                      resources.NewGeneric(export.ResourceTypes.LogTimestamp).Resource(),
			"dynatrace_log_grail":                          resources.NewGeneric(export.ResourceTypes.LogGrail).Resource(),
			"dynatrace_log_custom_attribute":               resources.NewGeneric(export.ResourceTypes.LogCustomAttribute).Resource(),
			"dynatrace_log_sensitive_data_masking":         resources.NewGeneric(export.ResourceTypes.LogSensitiveDataMasking).Resource(),
			"dynatrace_log_storage":                        resources.NewGeneric(export.ResourceTypes.LogStorage).Resource(),
			"dynatrace_log_buckets":                        resources.NewGeneric(export.ResourceTypes.LogBuckets).Resource(),
			"dynatrace_eula_settings":                      resources.NewGeneric(export.ResourceTypes.EULASettings).Resource(),
			"dynatrace_api_detection":                      resources.NewGeneric(export.ResourceTypes.APIDetectionRules).Resource(),
			"dynatrace_service_external_web_request":       resources.NewGeneric(export.ResourceTypes.ServiceExternalWebRequest).Resource(),
			"dynatrace_service_external_web_service":       resources.NewGeneric(export.ResourceTypes.ServiceExternalWebService).Resource(),
			"dynatrace_service_full_web_request":           resources.NewGeneric(export.ResourceTypes.ServiceFullWebRequest).Resource(),
			"dynatrace_service_full_web_service":           resources.NewGeneric(export.ResourceTypes.ServiceFullWebService).Resource(),
			"dynatrace_dashboards_allowlist":               resources.NewGeneric(export.ResourceTypes.DashboardsAllowlist).Resource(),
			"dynatrace_failure_detection_parameters":       resources.NewGeneric(export.ResourceTypes.FailureDetectionParameters).Resource(),
			"dynatrace_failure_detection_rules":            resources.NewGeneric(export.ResourceTypes.FailureDetectionRules).Resource(),
			"dynatrace_log_oneagent":                       resources.NewGeneric(export.ResourceTypes.LogOneAgent).Resource(),
			"dynatrace_issue_tracking":                     resources.NewGeneric(export.ResourceTypes.IssueTracking).Resource(),
			"dynatrace_geolocation":                        resources.NewGeneric(export.ResourceTypes.GeolocationSettings).Resource(),
			"dynatrace_user_session_metrics":               resources.NewGeneric(export.ResourceTypes.UserSessionCustomMetrics).Resource(),
			"dynatrace_custom_units":                       resources.NewGeneric(export.ResourceTypes.CustomUnits).Resource(),
			"dynatrace_disk_analytics":                     resources.NewGeneric(export.ResourceTypes.DiskAnalytics).Resource(),
			"dynatrace_network_traffic":                    resources.NewGeneric(export.ResourceTypes.NetworkTraffic).Resource(),
			"dynatrace_token_settings":                     resources.NewGeneric(export.ResourceTypes.TokenSettings).Resource(),
			"dynatrace_extension_execution_remote":         resources.NewGeneric(export.ResourceTypes.ExtensionExecutionRemote).Resource(),
			"dynatrace_k8s_pvc_anomalies":                  resources.NewGeneric(export.ResourceTypes.K8sPVCAnomalies).Resource(),
			"dynatrace_user_action_metrics":                resources.NewGeneric(export.ResourceTypes.UserActionCustomMetrics).Resource(),
			"dynatrace_web_app_javascript_version":         resources.NewGeneric(export.ResourceTypes.WebAppJavascriptVersion).Resource(),
			"dynatrace_web_app_javascript_updates":         resources.NewGeneric(export.ResourceTypes.WebAppJavascriptUpdates).Resource(),
			"dynatrace_opentelemetry_metrics":              resources.NewGeneric(export.ResourceTypes.OpenTelemetryMetrics).Resource(),
			"dynatrace_activegate_updates":                 resources.NewGeneric(export.ResourceTypes.ActiveGateUpdates).Resource(),
			"dynatrace_oneagent_default_version":           resources.NewGeneric(export.ResourceTypes.OneAgentDefaultVersion).Resource(),
			"dynatrace_oneagent_updates":                   resources.NewGeneric(export.ResourceTypes.OneAgentUpdates).Resource(),
			"dynatrace_ownership_teams":                    resources.NewGeneric(export.ResourceTypes.OwnershipTeams).Resource(),
			"dynatrace_ownership_config":                   resources.NewGeneric(export.ResourceTypes.OwnershipConfig).Resource(),
			"dynatrace_log_custom_source":                  resources.NewGeneric(export.ResourceTypes.LogCustomSource).Resource(),
			"dynatrace_application_detection_rule_v2":      resources.NewGeneric(export.ResourceTypes.ApplicationDetectionV2).Resource(),
			"dynatrace_kubernetes":                         resources.NewGeneric(export.ResourceTypes.Kubernetes).Resource(),
			"dynatrace_cloud_foundry":                      resources.NewGeneric(export.ResourceTypes.CloudFoundry).Resource(),
			"dynatrace_disk_anomaly_rules":                 resources.NewGeneric(export.ResourceTypes.DiskAnomalyDetectionRules).Resource(),
			"dynatrace_aws_anomalies":                      resources.NewGeneric(export.ResourceTypes.AWSAnomalies).Resource(),
			"dynatrace_vmware_anomalies":                   resources.NewGeneric(export.ResourceTypes.VMwareAnomalies).Resource(),
			"dynatrace_slo_v2":                             resources.NewGeneric(export.ResourceTypes.SLOV2).Resource(),
			"dynatrace_autotag_v2":                         resources.NewGeneric(export.ResourceTypes.AutoTagV2).Resource(),
			"dynatrace_business_events_oneagent":           resources.NewGeneric(export.ResourceTypes.BusinessEventsOneAgent).Resource(),
			"dynatrace_business_events_buckets":            resources.NewGeneric(export.ResourceTypes.BusinessEventsBuckets).Resource(),
			"dynatrace_business_events_metrics":            resources.NewGeneric(export.ResourceTypes.BusinessEventsMetrics).Resource(),
			"dynatrace_business_events_processing":         resources.NewGeneric(export.ResourceTypes.BusinessEventsProcessing).Resource(),
			"dynatrace_builtin_process_monitoring":         resources.NewGeneric(export.ResourceTypes.BuiltinProcessMonitoring).Resource(),
			"dynatrace_limit_outbound_connections":         resources.NewGeneric(export.ResourceTypes.LimitOutboundConnections).Resource(),
			"dynatrace_span_events":                        resources.NewGeneric(export.ResourceTypes.SpanEvents).Resource(),
			"dynatrace_vmware":                             resources.NewGeneric(export.ResourceTypes.VMware).Resource(),
			"dynatrace_web_app_key_performance_custom":     resources.NewGeneric(export.ResourceTypes.WebAppKeyPerformanceCustom).Resource(),
			"dynatrace_web_app_key_performance_load":       resources.NewGeneric(export.ResourceTypes.WebAppKeyPerformanceLoad).Resource(),
			"dynatrace_web_app_key_performance_xhr":        resources.NewGeneric(export.ResourceTypes.WebAppKeyPerformanceXHR).Resource(),
			"dynatrace_custom_device":                      resources.NewGeneric(export.ResourceTypes.CustomDevice).Resource(),
		},
		ConfigureContextFunc: config.ProviderConfigure,
	}
}
