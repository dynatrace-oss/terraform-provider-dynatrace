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
	"github.com/dynatrace-oss/terraform-provider-dynatrace/datasources/apitoken"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/datasources/application"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/datasources/appsec/attackalerting"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/datasources/appsec/vulnerabilityalerting"
	autotag "github.com/dynatrace-oss/terraform-provider-dynatrace/datasources/autotag"
	dsaws "github.com/dynatrace-oss/terraform-provider-dynatrace/datasources/aws"
	aws_credentials "github.com/dynatrace-oss/terraform-provider-dynatrace/datasources/credentials/aws"
	aws_supported_services "github.com/dynatrace-oss/terraform-provider-dynatrace/datasources/credentials/aws/supported_services"
	azure_credentials "github.com/dynatrace-oss/terraform-provider-dynatrace/datasources/credentials/azure"
	azure_supported_services "github.com/dynatrace-oss/terraform-provider-dynatrace/datasources/credentials/azure/supported_services"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/datasources/credentials/vault"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/datasources/dashboard"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/datasources/deployment/lambdaagent"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/datasources/documents/document"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/datasources/entities"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/datasources/entity"
	failure_detection_parameters "github.com/dynatrace-oss/terraform-provider-dynatrace/datasources/failuredetection/parameters"
	genericsettingsds "github.com/dynatrace-oss/terraform-provider-dynatrace/datasources/generic/settings"
	geocities "github.com/dynatrace-oss/terraform-provider-dynatrace/datasources/geographicregions/cities"
	geocountries "github.com/dynatrace-oss/terraform-provider-dynatrace/datasources/geographicregions/countries"
	georegions "github.com/dynatrace-oss/terraform-provider-dynatrace/datasources/geographicregions/regions"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/datasources/host"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/datasources/hub/items"
	ds_iam_groups "github.com/dynatrace-oss/terraform-provider-dynatrace/datasources/iam/groups"
	ds_iam_policies "github.com/dynatrace-oss/terraform-provider-dynatrace/datasources/iam/policies"
	ds_iam_users "github.com/dynatrace-oss/terraform-provider-dynatrace/datasources/iam/users"
	metricsds "github.com/dynatrace-oss/terraform-provider-dynatrace/datasources/metrics/calculated/service"
	mgmzds "github.com/dynatrace-oss/terraform-provider-dynatrace/datasources/mgmz"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/datasources/mobileapplication"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/datasources/process"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/datasources/processgroup"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/datasources/remoteenvironments"
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
	"github.com/dynatrace-oss/terraform-provider-dynatrace/resources/activegatetokens"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/resources/apitokens"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/resources/autotagrules"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/resources/backup"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/resources/bindings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/resources/customtags"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/resources/environments"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/resources/generic"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/resources/goldenstate"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/resources/internetproxy"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/resources/networkzones"
	mgmzperm "github.com/dynatrace-oss/terraform-provider-dynatrace/resources/permissions/mgmz"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/resources/policies"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/resources/preferences"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/resources/publicendpoints"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/resources/remoteaccess"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/resources/smtp"
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
	prv := &schema.Provider{
		Schema: map[string]*schema.Schema{
			"dt_env_url": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{"DYNATRACE_ENV_URL", "DT_ENV_URL", "DYNATRACE_ENVIRONMENT_URL", "DT_ENVIRONMENT_URL"}, nil),
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
			"client_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{"DT_CLIENT_ID", "DYNATRACE_CLIENT_ID"}, nil),
			},
			"account_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{"DT_ACCOUNT_ID", "DYNATRACE_ACCOUNT_ID"}, nil),
			},
			"client_secret": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{"DYNATRACE_CLIENT_SECRET", "DT_CLIENT_SECRET"}, nil),
			},
			"iam_client_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{"IAM_CLIENT_ID", "DYNATRACE_IAM_CLIENT_ID", "DT_IAM_CLIENT_ID", "DT_CLIENT_ID", "DYNATRACE_CLIENT_ID"}, nil),
			},
			"iam_account_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{"IAM_ACCOUNT_ID", "DYNATRACE_IAM_ACCOUNT_ID", "DT_IAM_ACCOUNT_ID", "DT_ACCOUNT_ID", "DYNATRACE_ACCOUNT_ID"}, nil),
			},
			"iam_client_secret": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{"IAM_CLIENT_SECRET", "DYNATRACE_IAM_CLIENT_SECRET", "DT_IAM_CLIENT_SECRET", "DYNATRACE_CLIENT_SECRET", "DT_CLIENT_SECRET"}, nil),
			},
			"iam_endpoint_url": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{"IAM_ENDPOINT_URL", "DYNATRACE_IAM_ENDPOINT_URL", "DT_IAM_ENDPOINT_URL", "DYNATRACE_ENDPOINT_URL", "DT_ENDPOINT_URL"}, nil),
			},
			"iam_token_url": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{"IAM_TOKEN_URL", "DYNATRACE_IAM_TOKEN_URL", "DT_IAM_TOKEN_URL", "DYNATRACE_TOKEN_URL", "DT_TOKEN_URL"}, nil),
			},
			"automation_client_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{"AUTOMATION_CLIENT_ID", "DYNATRACE_AUTOMATION_CLIENT_ID", "DT_AUTOMATION_CLIENT_ID", "DT_CLIENT_ID", "DYNATRACE_CLIENT_ID"}, nil),
			},
			"automation_client_secret": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{"AUTOMATION_CLIENT_SECRET", "DYNATRACE_AUTOMATION_CLIENT_SECRET", "DT_AUTOMATION_CLIENT_SECRET", "DYNATRACE_CLIENT_SECRET", "DT_CLIENT_SECRET"}, nil),
			},
			"automation_token_url": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{"AUTOMATION_TOKEN_URL", "DT_AUTOMATION_TOKEN_URL", "DYNATRACE_AUTOMATION_TOKEN_URL"}, nil),
				Description: "The URL that provides the Bearer tokens when accessing the Automation REST API. This is optional configuration when `dt_env_url` already specifies a SaaS Environment like `https://#####.live.dynatrace.com` or `https://#####.apps.dynatrace.com`",
			},
			"automation_env_url": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The URL of the Dynatrace Environment with Platform capabilities turned on (`https://#####.apps.dynatrace.com)`. This is optional configuration when `dt_env_url` already specifies a SaaS Environment like `https://#####.live.dynatrace.com` or `https://#####.apps.dynatrace.com`",
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{"AUTOMATION_ENVIRONMENT_URL", "DT_AUTOMATION_ENVIRONMENT_URL", "DYNATRACE_AUTOMATION_ENVIRONMENT_URL", "DYNATRACE_AUTOMATION_ENV_URL", "DT_AUTOMATION_ENV_URL"}, nil),
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
			"dynatrace_management_zone_v2":           mgmzds.DataSourceV2(),
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
			"dynatrace_iam_groups":                   ds_iam_groups.DataSourceMulti(),
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
			"dynatrace_vulnerability_alerting":       vulnerabilityalerting.DataSource(),
			"dynatrace_attack_alerting":              attackalerting.DataSource(),
			"dynatrace_remote_environments":          remoteenvironments.DataSource(),
			"dynatrace_hub_items":                    items.DataSource(),
			"dynatrace_documents":                    document.DataSource(),
			"dynatrace_iam_policies":                 ds_iam_policies.DataSource(),
			"dynatrace_iam_policy":                   ds_iam_policies.DataSourceSingle(),
			"dynatrace_lambda_agent_version":         lambdaagent.DataSource(),
			"dynatrace_autotag":                      autotag.DataSource(),
			"dynatrace_generic_settings":             genericsettingsds.DataSourceMultiple(),
			"dynatrace_generic_setting":              genericsettingsds.DataSource(),
			"dynatrace_api_tokens":                   apitoken.DataSourceMultiple(),
			"dynatrace_api_token":                    apitoken.DataSource(),
			"dynatrace_geo_countries":                geocountries.DataSource(),
			"dynatrace_geo_regions":                  georegions.DataSource(),
			"dynatrace_geo_cities":                   geocities.DataSource(),
		},
		ResourcesMap: map[string]*schema.Resource{
			"dynatrace_custom_service":                      resources.NewGeneric(export.ResourceTypes.CustomService).Resource(),
			"dynatrace_dashboard":                           resources.NewGeneric(export.ResourceTypes.Dashboard).Resource(),
			"dynatrace_json_dashboard":                      resources.NewGeneric(export.ResourceTypes.JSONDashboard).Resource(),
			"dynatrace_json_dashboard_base":                 resources.NewGeneric(export.ResourceTypes.JSONDashboardBase).Resource(),
			"dynatrace_management_zone":                     resources.NewGeneric(export.ResourceTypes.ManagementZone).Resource(),
			"dynatrace_management_zone_v2":                  resources.NewGeneric(export.ResourceTypes.ManagementZoneV2).Resource(),
			"dynatrace_maintenance_window":                  resources.NewGeneric(export.ResourceTypes.MaintenanceWindow).Resource(),
			"dynatrace_maintenance":                         resources.NewGeneric(export.ResourceTypes.Maintenance).Resource(),
			"dynatrace_request_attribute":                   resources.NewGeneric(export.ResourceTypes.RequestAttribute).Resource(),
			"dynatrace_alerting_profile":                    resources.NewGeneric(export.ResourceTypes.AlertingProfile).Resource(),
			"dynatrace_alerting":                            resources.NewGeneric(export.ResourceTypes.Alerting).Resource(),
			"dynatrace_notification":                        resources.NewGeneric(export.ResourceTypes.Notification).Resource(),
			"dynatrace_autotag":                             resources.NewGeneric(export.ResourceTypes.AutoTag).Resource(),
			"dynatrace_aws_credentials":                     resources.NewGeneric(export.ResourceTypes.AWSCredentials).Resource(),
			"dynatrace_aws_service":                         resources.NewGeneric(export.ResourceTypes.AWSService).Resource(),
			"dynatrace_azure_credentials":                   resources.NewGeneric(export.ResourceTypes.AzureCredentials).Resource(),
			"dynatrace_azure_service":                       resources.NewGeneric(export.ResourceTypes.AzureService).Resource(),
			"dynatrace_k8s_credentials":                     resources.NewGeneric(export.ResourceTypes.KubernetesCredentials).Resource(),
			"dynatrace_cloudfoundry_credentials":            resources.NewGeneric(export.ResourceTypes.CloudFoundryCredentials).Resource(),
			"dynatrace_service_anomalies":                   resources.NewGeneric(export.ResourceTypes.ServiceAnomalies).Resource(),
			"dynatrace_service_anomalies_v2":                resources.NewGeneric(export.ResourceTypes.ServiceAnomaliesV2).Resource(),
			"dynatrace_application_anomalies":               resources.NewGeneric(export.ResourceTypes.ApplicationAnomalies).Resource(),
			"dynatrace_host_anomalies":                      resources.NewGeneric(export.ResourceTypes.HostAnomalies).Resource(),
			"dynatrace_host_anomalies_v2":                   resources.NewGeneric(export.ResourceTypes.HostAnomaliesV2).Resource(),
			"dynatrace_database_anomalies":                  resources.NewGeneric(export.ResourceTypes.DatabaseAnomalies).Resource(),
			"dynatrace_database_anomalies_v2":               resources.NewGeneric(export.ResourceTypes.DatabaseAnomaliesV2).Resource(),
			"dynatrace_custom_anomalies":                    resources.NewGeneric(export.ResourceTypes.CustomAnomalies).Resource(),
			"dynatrace_metric_events":                       resources.NewGeneric(export.ResourceTypes.MetricEvents).Resource(),
			"dynatrace_disk_anomalies":                      resources.NewGeneric(export.ResourceTypes.DiskEventAnomalies).Resource(),
			"dynatrace_disk_anomalies_v2":                   resources.NewGeneric(export.ResourceTypes.DiskAnomaliesV2).Resource(),
			"dynatrace_disk_specific_anomalies_v2":          resources.NewGeneric(export.ResourceTypes.DiskSpecificAnomaliesV2).Resource(),
			"dynatrace_calculated_service_metric":           resources.NewGeneric(export.ResourceTypes.CalculatedServiceMetric).Resource(),
			"dynatrace_calculated_web_metric":               resources.NewGeneric(export.ResourceTypes.CalculatedWebMetric).Resource(),
			"dynatrace_calculated_mobile_metric":            resources.NewGeneric(export.ResourceTypes.CalculatedMobileMetric).Resource(),
			"dynatrace_calculated_synthetic_metric":         resources.NewGeneric(export.ResourceTypes.CalculatedSyntheticMetric).Resource(),
			"dynatrace_service_naming":                      resources.NewGeneric(export.ResourceTypes.ServiceNaming).Resource(),
			"dynatrace_host_naming":                         resources.NewGeneric(export.ResourceTypes.HostNaming).Resource(),
			"dynatrace_processgroup_naming":                 resources.NewGeneric(export.ResourceTypes.ProcessGroupNaming).Resource(),
			"dynatrace_slo":                                 resources.NewGeneric(export.ResourceTypes.SLO).Resource(),
			"dynatrace_span_entry_point":                    resources.NewGeneric(export.ResourceTypes.SpanEntryPoint).Resource(),
			"dynatrace_span_capture_rule":                   resources.NewGeneric(export.ResourceTypes.SpanCaptureRule).Resource(),
			"dynatrace_span_context_propagation":            resources.NewGeneric(export.ResourceTypes.SpanContextPropagation).Resource(),
			"dynatrace_resource_attributes":                 resources.NewGeneric(export.ResourceTypes.ResourceAttributes).Resource(),
			"dynatrace_span_attribute":                      resources.NewGeneric(export.ResourceTypes.SpanAttribute).Resource(),
			"dynatrace_dashboard_sharing":                   resources.NewGeneric(export.ResourceTypes.DashboardSharing).Resource(),
			"dynatrace_environment":                         environments.Resource(),
			"dynatrace_mobile_application":                  resources.NewGeneric(export.ResourceTypes.MobileApplication).Resource(),
			"dynatrace_browser_monitor":                     resources.NewGeneric(export.ResourceTypes.BrowserMonitor).Resource(),
			"dynatrace_http_monitor":                        resources.NewGeneric(export.ResourceTypes.HTTPMonitor).Resource(),
			"dynatrace_web_application":                     resources.NewGeneric(export.ResourceTypes.WebApplication).Resource(),
			"dynatrace_application_data_privacy":            resources.NewGeneric(export.ResourceTypes.ApplicationDataPrivacy).Resource(),
			"dynatrace_application_error_rules":             resources.NewGeneric(export.ResourceTypes.ApplicationErrorRules).Resource(),
			"dynatrace_request_naming":                      resources.NewGeneric(export.ResourceTypes.RequestNaming).Resource(),
			"dynatrace_request_namings":                     resources.NewGeneric(export.ResourceTypes.RequestNamings).Resource(),
			"dynatrace_user_group":                          usergroups.Resource(),
			"dynatrace_user":                                users.Resource(),
			"dynatrace_policy":                              policies.Resource(),
			"dynatrace_policy_bindings":                     bindings.Resource(),
			"dynatrace_mgmz_permission":                     mgmzperm.Resource(),
			"dynatrace_key_requests":                        resources.NewGeneric(export.ResourceTypes.KeyRequests).Resource(),
			"dynatrace_queue_manager":                       resources.NewGeneric(export.ResourceTypes.QueueManager).Resource(),
			"dynatrace_ibm_mq_filters":                      resources.NewGeneric(export.ResourceTypes.IBMMQFilters).Resource(),
			"dynatrace_queue_sharing_groups":                resources.NewGeneric(export.ResourceTypes.QueueSharingGroups).Resource(),
			"dynatrace_ims_bridges":                         resources.NewGeneric(export.ResourceTypes.IMSBridge).Resource(),
			"dynatrace_network_zones":                       resources.NewGeneric(export.ResourceTypes.NetworkZones).Resource(),
			"dynatrace_application_detection_rule":          resources.NewGeneric(export.ResourceTypes.ApplicationDetection).Resource(),
			"dynatrace_frequent_issues":                     resources.NewGeneric(export.ResourceTypes.FrequentIssues).Resource(),
			"dynatrace_ansible_tower_notification":          resources.NewGeneric(export.ResourceTypes.AnsibleTowerNotification).Resource(),
			"dynatrace_email_notification":                  resources.NewGeneric(export.ResourceTypes.EmailNotification).Resource(),
			"dynatrace_jira_notification":                   resources.NewGeneric(export.ResourceTypes.JiraNotification).Resource(),
			"dynatrace_ops_genie_notification":              resources.NewGeneric(export.ResourceTypes.OpsGenieNotification).Resource(),
			"dynatrace_pager_duty_notification":             resources.NewGeneric(export.ResourceTypes.PagerDutyNotification).Resource(),
			"dynatrace_service_now_notification":            resources.NewGeneric(export.ResourceTypes.ServiceNowNotification).Resource(),
			"dynatrace_slack_notification":                  resources.NewGeneric(export.ResourceTypes.SlackNotification).Resource(),
			"dynatrace_trello_notification":                 resources.NewGeneric(export.ResourceTypes.TrelloNotification).Resource(),
			"dynatrace_victor_ops_notification":             resources.NewGeneric(export.ResourceTypes.VictorOpsNotification).Resource(),
			"dynatrace_webhook_notification":                resources.NewGeneric(export.ResourceTypes.WebHookNotification).Resource(),
			"dynatrace_xmatters_notification":               resources.NewGeneric(export.ResourceTypes.XMattersNotification).Resource(),
			"dynatrace_credentials":                         resources.NewGeneric(export.ResourceTypes.Credentials).Resource(),
			"dynatrace_synthetic_location":                  resources.NewGeneric(export.ResourceTypes.SyntheticLocation).Resource(),
			"dynatrace_network_zone":                        resources.NewGeneric(export.ResourceTypes.NetworkZone).Resource(),
			"dynatrace_iam_user":                            resources.NewGeneric(export.ResourceTypes.IAMUser, resources.CredValIAM).Resource(),
			"dynatrace_iam_group":                           resources.NewGeneric(export.ResourceTypes.IAMGroup, resources.CredValIAM).Resource(),
			"dynatrace_iam_permission":                      resources.NewGeneric(export.ResourceTypes.IAMPermission, resources.CredValIAM).Resource(),
			"dynatrace_iam_policy":                          resources.NewGeneric(export.ResourceTypes.IAMPolicy, resources.CredValIAM).Resource(),
			"dynatrace_iam_policy_bindings":                 resources.NewGeneric(export.ResourceTypes.IAMPolicyBindings, resources.CredValIAM).Resource(),
			"dynatrace_iam_policy_bindings_v2":              resources.NewGeneric(export.ResourceTypes.IAMPolicyBindingsV2, resources.CredValIAM).Resource(),
			"dynatrace_api_token":                           apitokens.Resource(),
			"dynatrace_custom_tags":                         customtags.Resource(),
			"dynatrace_pg_anomalies":                        resources.NewGeneric(export.ResourceTypes.ProcessGroupAnomalies).Resource(),
			"dynatrace_ddu_pool":                            resources.NewGeneric(export.ResourceTypes.DDUPool).Resource(),
			"dynatrace_pg_alerting":                         resources.NewGeneric(export.ResourceTypes.ProcessGroupAlerting).Resource(),
			"dynatrace_process_monitoring_rule":             resources.NewGeneric(export.ResourceTypes.ProcessMonitoringRule).Resource(),
			"dynatrace_custom_app_anomalies":                resources.NewGeneric(export.ResourceTypes.CustomAppAnomalies).Resource(),
			"dynatrace_custom_app_crash_rate":               resources.NewGeneric(export.ResourceTypes.CustomAppCrashRate).Resource(),
			"dynatrace_process_monitoring":                  resources.NewGeneric(export.ResourceTypes.ProcessMonitoring).Resource(),
			"dynatrace_process_availability":                resources.NewGeneric(export.ResourceTypes.ProcessAvailability).Resource(),
			"dynatrace_process_group_detection":             resources.NewGeneric(export.ResourceTypes.AdvancedProcessGroupDetectionRule).Resource(),
			"dynatrace_mobile_app_anomalies":                resources.NewGeneric(export.ResourceTypes.MobileAppAnomalies).Resource(),
			"dynatrace_mobile_app_crash_rate":               resources.NewGeneric(export.ResourceTypes.MobileAppCrashRate).Resource(),
			"dynatrace_web_app_anomalies":                   resources.NewGeneric(export.ResourceTypes.WebAppAnomalies).Resource(),
			"dynatrace_muted_requests":                      resources.NewGeneric(export.ResourceTypes.MutedRequests).Resource(),
			"dynatrace_connectivity_alerts":                 resources.NewGeneric(export.ResourceTypes.ConnectivityAlerts).Resource(),
			"dynatrace_declarative_grouping":                resources.NewGeneric(export.ResourceTypes.DeclarativeGrouping).Resource(),
			"dynatrace_host_monitoring":                     resources.NewGeneric(export.ResourceTypes.HostMonitoring).Resource(),
			"dynatrace_host_process_group_monitoring":       resources.NewGeneric(export.ResourceTypes.HostProcessGroupMonitoring).Resource(),
			"dynatrace_rum_ip_locations":                    resources.NewGeneric(export.ResourceTypes.RUMIPLocations).Resource(),
			"dynatrace_custom_app_enablement":               resources.NewGeneric(export.ResourceTypes.CustomAppEnablement).Resource(),
			"dynatrace_mobile_app_enablement":               resources.NewGeneric(export.ResourceTypes.MobileAppEnablement).Resource(),
			"dynatrace_web_app_enablement":                  resources.NewGeneric(export.ResourceTypes.WebAppEnablement).Resource(),
			"dynatrace_process_group_rum":                   resources.NewGeneric(export.ResourceTypes.RUMProcessGroup).Resource(),
			"dynatrace_rum_provider_breakdown":              resources.NewGeneric(export.ResourceTypes.RUMProviderBreakdown).Resource(),
			"dynatrace_user_experience_score":               resources.NewGeneric(export.ResourceTypes.UserExperienceScore).Resource(),
			"dynatrace_web_app_resource_cleanup":            resources.NewGeneric(export.ResourceTypes.WebAppResourceCleanup).Resource(),
			"dynatrace_update_windows":                      resources.NewGeneric(export.ResourceTypes.UpdateWindows).Resource(),
			"dynatrace_process_group_detection_flags":       resources.NewGeneric(export.ResourceTypes.ProcessGroupDetectionFlags).Resource(),
			"dynatrace_process_group_monitoring":            resources.NewGeneric(export.ResourceTypes.ProcessGroupMonitoring).Resource(),
			"dynatrace_process_group_simple_detection":      resources.NewGeneric(export.ResourceTypes.ProcessGroupSimpleDetection).Resource(),
			"dynatrace_log_metrics":                         resources.NewGeneric(export.ResourceTypes.LogMetrics).Resource(),
			"dynatrace_browser_monitor_performance":         resources.NewGeneric(export.ResourceTypes.BrowserMonitorPerformanceThresholds).Resource(),
			"dynatrace_http_monitor_performance":            resources.NewGeneric(export.ResourceTypes.HttpMonitorPerformanceThresholds).Resource(),
			"dynatrace_http_monitor_cookies":                resources.NewGeneric(export.ResourceTypes.HttpMonitorCookies).Resource(),
			"dynatrace_session_replay_web_privacy":          resources.NewGeneric(export.ResourceTypes.SessionReplayWebPrivacy).Resource(),
			"dynatrace_session_replay_resource_capture":     resources.NewGeneric(export.ResourceTypes.SessionReplayResourceCapture).Resource(),
			"dynatrace_usability_analytics":                 resources.NewGeneric(export.ResourceTypes.UsabilityAnalytics).Resource(),
			"dynatrace_synthetic_availability":              resources.NewGeneric(export.ResourceTypes.SyntheticAvailability).Resource(),
			"dynatrace_browser_monitor_outage":              resources.NewGeneric(export.ResourceTypes.BrowserMonitorOutageHandling).Resource(),
			"dynatrace_http_monitor_outage":                 resources.NewGeneric(export.ResourceTypes.HttpMonitorOutageHandling).Resource(),
			"dynatrace_cloudapp_workloaddetection":          resources.NewGeneric(export.ResourceTypes.CloudAppWorkloadDetection).Resource(),
			"dynatrace_mainframe_transaction_monitoring":    resources.NewGeneric(export.ResourceTypes.MainframeTransactionMonitoring).Resource(),
			"dynatrace_monitored_technologies_apache":       resources.NewGeneric(export.ResourceTypes.MonitoredTechnologiesApache).Resource(),
			"dynatrace_monitored_technologies_dotnet":       resources.NewGeneric(export.ResourceTypes.MonitoredTechnologiesDotNet).Resource(),
			"dynatrace_monitored_technologies_go":           resources.NewGeneric(export.ResourceTypes.MonitoredTechnologiesGo).Resource(),
			"dynatrace_monitored_technologies_iis":          resources.NewGeneric(export.ResourceTypes.MonitoredTechnologiesIIS).Resource(),
			"dynatrace_monitored_technologies_java":         resources.NewGeneric(export.ResourceTypes.MonitoredTechnologiesJava).Resource(),
			"dynatrace_monitored_technologies_nginx":        resources.NewGeneric(export.ResourceTypes.MonitoredTechnologiesNGINX).Resource(),
			"dynatrace_monitored_technologies_nodejs":       resources.NewGeneric(export.ResourceTypes.MonitoredTechnologiesNodeJS).Resource(),
			"dynatrace_monitored_technologies_opentracing":  resources.NewGeneric(export.ResourceTypes.MonitoredTechnologiesOpenTracing).Resource(),
			"dynatrace_monitored_technologies_php":          resources.NewGeneric(export.ResourceTypes.MonitoredTechnologiesPHP).Resource(),
			"dynatrace_monitored_technologies_varnish":      resources.NewGeneric(export.ResourceTypes.MonitoredTechnologiesVarnish).Resource(),
			"dynatrace_monitored_technologies_wsmb":         resources.NewGeneric(export.ResourceTypes.MonitoredTechnologiesWSMB).Resource(),
			"dynatrace_process_visibility":                  resources.NewGeneric(export.ResourceTypes.ProcessVisibility).Resource(),
			"dynatrace_rum_host_headers":                    resources.NewGeneric(export.ResourceTypes.RUMHostHeaders).Resource(),
			"dynatrace_rum_ip_determination":                resources.NewGeneric(export.ResourceTypes.RUMIPDetermination).Resource(),
			"dynatrace_mobile_app_request_errors":           resources.NewGeneric(export.ResourceTypes.MobileAppRequestErrors).Resource(),
			"dynatrace_transaction_start_filters":           resources.NewGeneric(export.ResourceTypes.TransactionStartFilters).Resource(),
			"dynatrace_oneagent_features":                   resources.NewGeneric(export.ResourceTypes.OneAgentFeatures).Resource(),
			"dynatrace_rum_overload_prevention":             resources.NewGeneric(export.ResourceTypes.RUMOverloadPrevention).Resource(),
			"dynatrace_rum_advanced_correlation":            resources.NewGeneric(export.ResourceTypes.RUMAdvancedCorrelation).Resource(),
			"dynatrace_web_app_beacon_origins":              resources.NewGeneric(export.ResourceTypes.WebAppBeaconOrigins).Resource(),
			"dynatrace_web_app_resource_types":              resources.NewGeneric(export.ResourceTypes.WebAppResourceTypes).Resource(),
			"dynatrace_generic_types":                       resources.NewGeneric(export.ResourceTypes.GenericTypes).Resource(),
			"dynatrace_generic_relationships":               resources.NewGeneric(export.ResourceTypes.GenericRelationships).Resource(),
			"dynatrace_slo_normalization":                   resources.NewGeneric(export.ResourceTypes.SLONormalization).Resource(),
			"dynatrace_data_privacy":                        resources.NewGeneric(export.ResourceTypes.DataPrivacy).Resource(),
			"dynatrace_service_failure":                     resources.NewGeneric(export.ResourceTypes.ServiceFailure).Resource(),
			"dynatrace_service_http_failure":                resources.NewGeneric(export.ResourceTypes.ServiceHTTPFailure).Resource(),
			"dynatrace_disk_options":                        resources.NewGeneric(export.ResourceTypes.DiskOptions).Resource(),
			"dynatrace_os_services":                         resources.NewGeneric(export.ResourceTypes.OSServices).Resource(),
			"dynatrace_extension_execution_controller":      resources.NewGeneric(export.ResourceTypes.ExtensionExecutionController).Resource(),
			"dynatrace_nettracer":                           resources.NewGeneric(export.ResourceTypes.NetTracerTraffic).Resource(),
			"dynatrace_aix_extension":                       resources.NewGeneric(export.ResourceTypes.AIXExtension).Resource(),
			"dynatrace_metric_metadata":                     resources.NewGeneric(export.ResourceTypes.MetricMetadata).Resource(),
			"dynatrace_metric_query":                        resources.NewGeneric(export.ResourceTypes.MetricQuery).Resource(),
			"dynatrace_activegate_token":                    resources.NewGeneric(export.ResourceTypes.ActiveGateToken).Resource(),
			"dynatrace_ag_token":                            activegatetokens.Resource(),
			"dynatrace_audit_log":                           resources.NewGeneric(export.ResourceTypes.AuditLog).Resource(),
			"dynatrace_k8s_cluster_anomalies":               resources.NewGeneric(export.ResourceTypes.K8sClusterAnomalies).Resource(),
			"dynatrace_k8s_namespace_anomalies":             resources.NewGeneric(export.ResourceTypes.K8sNamespaceAnomalies).Resource(),
			"dynatrace_k8s_node_anomalies":                  resources.NewGeneric(export.ResourceTypes.K8sNodeAnomalies).Resource(),
			"dynatrace_k8s_workload_anomalies":              resources.NewGeneric(export.ResourceTypes.K8sWorkloadAnomalies).Resource(),
			"dynatrace_container_builtin_rule":              resources.NewGeneric(export.ResourceTypes.ContainerBuiltinRule).Resource(),
			"dynatrace_container_rule":                      resources.NewGeneric(export.ResourceTypes.ContainerRule).Resource(),
			"dynatrace_container_technology":                resources.NewGeneric(export.ResourceTypes.ContainerTechnology).Resource(),
			"dynatrace_container_registry":                  resources.NewGeneric(export.ResourceTypes.ContainerRegistry).Resource(),
			"dynatrace_remote_environments":                 resources.NewGeneric(export.ResourceTypes.RemoteEnvironments).Resource(),
			"dynatrace_web_app_custom_errors":               resources.NewGeneric(export.ResourceTypes.WebAppCustomErrors).Resource(),
			"dynatrace_web_app_request_errors":              resources.NewGeneric(export.ResourceTypes.WebAppRequestErrors).Resource(),
			"dynatrace_user_settings":                       resources.NewGeneric(export.ResourceTypes.UserSettings).Resource(),
			"dynatrace_dashboards_general":                  resources.NewGeneric(export.ResourceTypes.DashboardsGeneral).Resource(),
			"dynatrace_dashboards_presets":                  resources.NewGeneric(export.ResourceTypes.DashboardsPresets).Resource(),
			"dynatrace_log_processing":                      resources.NewGeneric(export.ResourceTypes.LogProcessing).Resource(),
			"dynatrace_log_events":                          resources.NewGeneric(export.ResourceTypes.LogEvents).Resource(),
			"dynatrace_log_timestamp":                       resources.NewGeneric(export.ResourceTypes.LogTimestamp).Resource(),
			"dynatrace_log_grail":                           resources.NewGeneric(export.ResourceTypes.LogGrail).Resource(),
			"dynatrace_log_custom_attribute":                resources.NewGeneric(export.ResourceTypes.LogCustomAttribute).Resource(),
			"dynatrace_log_sensitive_data_masking":          resources.NewGeneric(export.ResourceTypes.LogSensitiveDataMasking).Resource(),
			"dynatrace_log_storage":                         resources.NewGeneric(export.ResourceTypes.LogStorage).Resource(),
			"dynatrace_log_buckets":                         resources.NewGeneric(export.ResourceTypes.LogBuckets).Resource(),
			"dynatrace_log_security_context":                resources.NewGeneric(export.ResourceTypes.LogSecurityContext).Resource(),
			"dynatrace_eula_settings":                       resources.NewGeneric(export.ResourceTypes.EULASettings).Resource(),
			"dynatrace_api_detection":                       resources.NewGeneric(export.ResourceTypes.APIDetectionRules).Resource(),
			"dynatrace_service_external_web_request":        resources.NewGeneric(export.ResourceTypes.ServiceExternalWebRequest).Resource(),
			"dynatrace_service_external_web_service":        resources.NewGeneric(export.ResourceTypes.ServiceExternalWebService).Resource(),
			"dynatrace_service_full_web_request":            resources.NewGeneric(export.ResourceTypes.ServiceFullWebRequest).Resource(),
			"dynatrace_service_full_web_service":            resources.NewGeneric(export.ResourceTypes.ServiceFullWebService).Resource(),
			"dynatrace_dashboards_allowlist":                resources.NewGeneric(export.ResourceTypes.DashboardsAllowlist).Resource(),
			"dynatrace_failure_detection_parameters":        resources.NewGeneric(export.ResourceTypes.FailureDetectionParameters).Resource(),
			"dynatrace_failure_detection_rules":             resources.NewGeneric(export.ResourceTypes.FailureDetectionRules).Resource(),
			"dynatrace_log_oneagent":                        resources.NewGeneric(export.ResourceTypes.LogOneAgent).Resource(),
			"dynatrace_issue_tracking":                      resources.NewGeneric(export.ResourceTypes.IssueTracking).Resource(),
			"dynatrace_geolocation":                         resources.NewGeneric(export.ResourceTypes.GeolocationSettings).Resource(),
			"dynatrace_user_session_metrics":                resources.NewGeneric(export.ResourceTypes.UserSessionCustomMetrics).Resource(),
			"dynatrace_custom_units":                        resources.NewGeneric(export.ResourceTypes.CustomUnits).Resource(),
			"dynatrace_disk_analytics":                      resources.NewGeneric(export.ResourceTypes.DiskAnalytics).Resource(),
			"dynatrace_network_traffic":                     resources.NewGeneric(export.ResourceTypes.NetworkTraffic).Resource(),
			"dynatrace_token_settings":                      resources.NewGeneric(export.ResourceTypes.TokenSettings).Resource(),
			"dynatrace_extension_execution_remote":          resources.NewGeneric(export.ResourceTypes.ExtensionExecutionRemote).Resource(),
			"dynatrace_k8s_pvc_anomalies":                   resources.NewGeneric(export.ResourceTypes.K8sPVCAnomalies).Resource(),
			"dynatrace_user_action_metrics":                 resources.NewGeneric(export.ResourceTypes.UserActionCustomMetrics).Resource(),
			"dynatrace_web_app_javascript_version":          resources.NewGeneric(export.ResourceTypes.WebAppJavascriptVersion).Resource(),
			"dynatrace_web_app_javascript_updates":          resources.NewGeneric(export.ResourceTypes.WebAppJavascriptUpdates).Resource(),
			"dynatrace_opentelemetry_metrics":               resources.NewGeneric(export.ResourceTypes.OpenTelemetryMetrics).Resource(),
			"dynatrace_activegate_updates":                  resources.NewGeneric(export.ResourceTypes.ActiveGateUpdates).Resource(),
			"dynatrace_oneagent_default_version":            resources.NewGeneric(export.ResourceTypes.OneAgentDefaultVersion).Resource(),
			"dynatrace_oneagent_updates":                    resources.NewGeneric(export.ResourceTypes.OneAgentUpdates).Resource(),
			"dynatrace_ownership_teams":                     resources.NewGeneric(export.ResourceTypes.OwnershipTeams).Resource(),
			"dynatrace_ownership_config":                    resources.NewGeneric(export.ResourceTypes.OwnershipConfig).Resource(),
			"dynatrace_log_custom_source":                   resources.NewGeneric(export.ResourceTypes.LogCustomSource).Resource(),
			"dynatrace_application_detection_rule_v2":       resources.NewGeneric(export.ResourceTypes.ApplicationDetectionV2).Resource(),
			"dynatrace_kubernetes":                          resources.NewGeneric(export.ResourceTypes.Kubernetes).Resource(),
			"dynatrace_cloud_foundry":                       resources.NewGeneric(export.ResourceTypes.CloudFoundry).Resource(),
			"dynatrace_disk_anomaly_rules":                  resources.NewGeneric(export.ResourceTypes.DiskAnomalyDetectionRules).Resource(),
			"dynatrace_aws_anomalies":                       resources.NewGeneric(export.ResourceTypes.AWSAnomalies).Resource(),
			"dynatrace_vmware_anomalies":                    resources.NewGeneric(export.ResourceTypes.VMwareAnomalies).Resource(),
			"dynatrace_slo_v2":                              resources.NewGeneric(export.ResourceTypes.SLOV2).Resource(),
			"dynatrace_autotag_v2":                          resources.NewGeneric(export.ResourceTypes.AutoTagV2).Resource(),
			"dynatrace_business_events_oneagent":            resources.NewGeneric(export.ResourceTypes.BusinessEventsOneAgent).Resource(),
			"dynatrace_business_events_buckets":             resources.NewGeneric(export.ResourceTypes.BusinessEventsBuckets).Resource(),
			"dynatrace_business_events_metrics":             resources.NewGeneric(export.ResourceTypes.BusinessEventsMetrics).Resource(),
			"dynatrace_business_events_processing":          resources.NewGeneric(export.ResourceTypes.BusinessEventsProcessing).Resource(),
			"dynatrace_business_events_security_context":    resources.NewGeneric(export.ResourceTypes.BusinessEventsSecurityContext).Resource(),
			"dynatrace_builtin_process_monitoring":          resources.NewGeneric(export.ResourceTypes.BuiltinProcessMonitoring).Resource(),
			"dynatrace_limit_outbound_connections":          resources.NewGeneric(export.ResourceTypes.LimitOutboundConnections).Resource(),
			"dynatrace_span_events":                         resources.NewGeneric(export.ResourceTypes.SpanEvents).Resource(),
			"dynatrace_vmware":                              resources.NewGeneric(export.ResourceTypes.VMware).Resource(),
			"dynatrace_web_app_key_performance_custom":      resources.NewGeneric(export.ResourceTypes.WebAppKeyPerformanceCustom).Resource(),
			"dynatrace_web_app_key_performance_load":        resources.NewGeneric(export.ResourceTypes.WebAppKeyPerformanceLoad).Resource(),
			"dynatrace_web_app_key_performance_xhr":         resources.NewGeneric(export.ResourceTypes.WebAppKeyPerformanceXHR).Resource(),
			"dynatrace_mobile_app_key_performance":          resources.NewGeneric(export.ResourceTypes.MobileAppKeyPerformance).Resource(),
			"dynatrace_custom_device":                       resources.NewGeneric(export.ResourceTypes.CustomDevice).Resource(),
			"dynatrace_k8s_monitoring":                      resources.NewGeneric(export.ResourceTypes.K8sMonitoring).Resource(),
			"dynatrace_host_monitoring_mode":                resources.NewGeneric(export.ResourceTypes.HostMonitoringMode).Resource(),
			"dynatrace_ip_address_masking":                  resources.NewGeneric(export.ResourceTypes.IPAddressMasking).Resource(),
			"dynatrace_automation_workflow":                 resources.NewGeneric(export.ResourceTypes.AutomationWorkflow).Resource(),
			"dynatrace_automation_business_calendar":        resources.NewGeneric(export.ResourceTypes.AutomationBusinessCalendar).Resource(),
			"dynatrace_automation_scheduling_rule":          resources.NewGeneric(export.ResourceTypes.AutomationSchedulingRule).Resource(),
			"dynatrace_vulnerability_settings":              resources.NewGeneric(export.ResourceTypes.AppSecVulnerabilitySettings).Resource(),
			"dynatrace_vulnerability_third_party":           resources.NewGeneric(export.ResourceTypes.AppSecVulnerabilityThirdParty).Resource(),
			"dynatrace_vulnerability_code":                  resources.NewGeneric(export.ResourceTypes.AppSecVulnerabilityCode).Resource(),
			"dynatrace_appsec_notification":                 resources.NewGeneric(export.ResourceTypes.AppSecNotification).Resource(),
			"dynatrace_vulnerability_alerting":              resources.NewGeneric(export.ResourceTypes.AppSecVulnerabilityAlerting).Resource(),
			"dynatrace_attack_alerting":                     resources.NewGeneric(export.ResourceTypes.AppSecAttackAlerting).Resource(),
			"dynatrace_attack_settings":                     resources.NewGeneric(export.ResourceTypes.AppSecAttackSettings).Resource(),
			"dynatrace_attack_rules":                        resources.NewGeneric(export.ResourceTypes.AppSecAttackRules).Resource(),
			"dynatrace_attack_allowlist":                    resources.NewGeneric(export.ResourceTypes.AppSecAttackAllowlist).Resource(),
			"dynatrace_unified_services_metrics":            resources.NewGeneric(export.ResourceTypes.UnifiedServicesMetrics).Resource(),
			"dynatrace_unified_services_opentel":            resources.NewGeneric(export.ResourceTypes.UnifiedServicesOpenTel).Resource(),
			"dynatrace_autotag_rules":                       autotagrules.Resource(),
			"dynatrace_generic_setting":                     generic.Resource(),
			"dynatrace_managed_smtp":                        smtp.Resource(),
			"dynatrace_managed_internet_proxy":              internetproxy.Resource(),
			"dynatrace_managed_preferences":                 preferences.Resource(),
			"dynatrace_platform_bucket":                     resources.NewGeneric(export.ResourceTypes.PlatformBucket).Resource(),
			"dynatrace_managed_public_endpoints":            publicendpoints.Resource(),
			"dynatrace_managed_backup":                      backup.Resource(),
			"dynatrace_managed_remote_access":               remoteaccess.Resource(),
			"dynatrace_key_user_action":                     resources.NewGeneric(export.ResourceTypes.KeyUserAction).Resource(),
			"dynatrace_url_based_sampling":                  resources.NewGeneric(export.ResourceTypes.UrlBasedSampling).Resource(),
			"dynatrace_host_monitoring_advanced":            resources.NewGeneric(export.ResourceTypes.HostMonitoringAdvanced).Resource(),
			"dynatrace_attribute_allow_list":                resources.NewGeneric(export.ResourceTypes.AttributeAllowList).Resource(),
			"dynatrace_attribute_block_list":                resources.NewGeneric(export.ResourceTypes.AttributeBlockList).Resource(),
			"dynatrace_attribute_masking":                   resources.NewGeneric(export.ResourceTypes.AttributeMasking).Resource(),
			"dynatrace_attributes_preferences":              resources.NewGeneric(export.ResourceTypes.AttributesPreferences).Resource(),
			"dynatrace_oneagent_side_masking":               resources.NewGeneric(export.ResourceTypes.OneAgentSideMasking).Resource(),
			"dynatrace_hub_subscriptions":                   resources.NewGeneric(export.ResourceTypes.HubSubscriptions).Resource(),
			"dynatrace_mobile_notifications":                resources.NewGeneric(export.ResourceTypes.MobileNotifications).Resource(),
			"dynatrace_crashdump_analytics":                 resources.NewGeneric(export.ResourceTypes.CrashdumpAnalytics).Resource(),
			"dynatrace_app_monitoring":                      resources.NewGeneric(export.ResourceTypes.AppMonitoring).Resource(),
			"dynatrace_grail_security_context":              resources.NewGeneric(export.ResourceTypes.GrailSecurityContext).Resource(),
			"dynatrace_site_reliability_guardian":           resources.NewGeneric(export.ResourceTypes.SiteReliabilityGuardian).Resource(),
			"dynatrace_automation_workflow_jira":            resources.NewGeneric(export.ResourceTypes.JiraForWorkflows).Resource(),
			"dynatrace_automation_workflow_slack":           resources.NewGeneric(export.ResourceTypes.SlackForWorkflows).Resource(),
			"dynatrace_kubernetes_app":                      resources.NewGeneric(export.ResourceTypes.KubernetesApp).Resource(),
			"dynatrace_grail_metrics_allowall":              resources.NewGeneric(export.ResourceTypes.GrailMetricsAllowall).Resource(),
			"dynatrace_grail_metrics_allowlist":             resources.NewGeneric(export.ResourceTypes.GrailMetricsAllowlist).Resource(),
			"dynatrace_web_app_beacon_endpoint":             resources.NewGeneric(export.ResourceTypes.WebAppBeaconEndpoint).Resource(),
			"dynatrace_web_app_custom_config_properties":    resources.NewGeneric(export.ResourceTypes.WebAppCustomConfigProperties).Resource(),
			"dynatrace_web_app_injection_cookie":            resources.NewGeneric(export.ResourceTypes.WebAppInjectionCookie).Resource(),
			"dynatrace_http_monitor_script":                 resources.NewGeneric(export.ResourceTypes.HTTPMonitorScript).Resource(),
			"dynatrace_managed_network_zones":               networkzones.Resource(),
			"dynatrace_hub_extension_config":                resources.NewGeneric(export.ResourceTypes.HubExtensionConfig).Resource(),
			"dynatrace_hub_extension_active_version":        resources.NewGeneric(export.ResourceTypes.HubActiveExtensionVersion).Resource(),
			"dynatrace_document":                            resources.NewGeneric(export.ResourceTypes.Documents).Resource(),
			"dynatrace_direct_shares":                       resources.NewGeneric(export.ResourceTypes.DirectShares).Resource(),
			"dynatrace_db_app_feature_flags":                resources.NewGeneric(export.ResourceTypes.DatabaseAppFeatureFlags).Resource(),
			"dynatrace_infraops_app_feature_flags":          resources.NewGeneric(export.ResourceTypes.InfraOpsAppFeatureFlags).Resource(),
			"dynatrace_ebpf_service_discovery":              resources.NewGeneric(export.ResourceTypes.EBPFServiceDiscovery).Resource(),
			"dynatrace_davis_anomaly_detectors":             resources.NewGeneric(export.ResourceTypes.DavisAnomalyDetectors).Resource(),
			"dynatrace_log_debug_settings":                  resources.NewGeneric(export.ResourceTypes.LogDebugSettings).Resource(),
			"dynatrace_infraops_app_settings":               resources.NewGeneric(export.ResourceTypes.InfraOpsAppSettings).Resource(),
			"dynatrace_disk_edge_anomaly_detectors":         resources.NewGeneric(export.ResourceTypes.DiskEdgeAnomalyDetectors).Resource(),
			"dynatrace_report":                              resources.NewGeneric(export.ResourceTypes.Reports).Resource(),
			"dynatrace_network_monitor":                     resources.NewGeneric(export.ResourceTypes.NetworkMonitor).Resource(),
			"dynatrace_network_monitor_outage":              resources.NewGeneric(export.ResourceTypes.NetworkMonitorOutageHandling).Resource(),
			"dynatrace_hub_permissions":                     resources.NewGeneric(export.ResourceTypes.HubPermissions).Resource(),
			"dynatrace_automation_workflow_k8s_connections": resources.NewGeneric(export.ResourceTypes.K8sAutomationConnections).Resource(),
			"dynatrace_business_events_oneagent_outgoing":   resources.NewGeneric(export.ResourceTypes.BusinessEventsOneAgentOutgoing).Resource(),
			"dynatrace_oneagent_default_mode":               resources.NewGeneric(export.ResourceTypes.OneAgentDefaultMode).Resource(),
			"dynatrace_web_app_custom_injection":            resources.NewGeneric(export.ResourceTypes.WebAppCustomInjectionRules).Resource(),
			"dynatrace_discovery_default_rules":             resources.NewGeneric(export.ResourceTypes.DiscoveryDefaultRules).Resource(),
			"dynatrace_discovery_feature_flags":             resources.NewGeneric(export.ResourceTypes.DiscoveryFeatureFlags).Resource(),
			"dynatrace_histogram_metrics":                   resources.NewGeneric(export.ResourceTypes.HistogramMetrics).Resource(),
			"dynatrace_kubernetes_enrichment":               resources.NewGeneric(export.ResourceTypes.KubernetesEnrichment).Resource(),
			"dynatrace_devobs_git_onprem":                   resources.NewGeneric(export.ResourceTypes.DevObsGitOnPrem).Resource(),
			"dynatrace_automation_workflow_aws_connections": resources.NewGeneric(export.ResourceTypes.AWSAutomationConnections).Resource(),
			"dynatrace_devobs_agent_optin":                  resources.NewGeneric(export.ResourceTypes.DevObsAgentOptin).Resource(),
			"dynatrace_devobs_data_masking":                 resources.NewGeneric(export.ResourceTypes.DevObsDataMasking).Resource(),
			"dynatrace_davis_copilot":                       resources.NewGeneric(export.ResourceTypes.DavisCoPilot).Resource(),
			"dynatrace_golden_state":                        goldenstate.Resource(),
			"dynatrace_openpipeline_logs":                   resources.NewGeneric(export.ResourceTypes.OpenPipelineLogs).Resource(),
			"dynatrace_openpipeline_events":                 resources.NewGeneric(export.ResourceTypes.OpenPipelineEvents).Resource(),
			"dynatrace_openpipeline_security_events":        resources.NewGeneric(export.ResourceTypes.OpenPipelineSecurityEvents).Resource(),
			"dynatrace_openpipeline_business_events":        resources.NewGeneric(export.ResourceTypes.OpenPipelineBusinessEvents).Resource(),
			"dynatrace_openpipeline_sdlc_events":            resources.NewGeneric(export.ResourceTypes.OpenPipelineSDLCEvents).Resource(),
			"dynatrace_cloud_development_environments":      resources.NewGeneric(export.ResourceTypes.CloudDevelopmentEnvironments).Resource(),
			"dynatrace_kubernetes_spm":                      resources.NewGeneric(export.ResourceTypes.KubernetesSPM).Resource(),
			"dynatrace_log_agent_feature_flags":             resources.NewGeneric(export.ResourceTypes.LogAgentFeatureFlags).Resource(),
			"dynatrace_problem_record_propagation_rules":    resources.NewGeneric(export.ResourceTypes.ProblemRecordPropagationRules).Resource(),
			"dynatrace_problem_fields":                      resources.NewGeneric(export.ResourceTypes.ProblemFields).Resource(),
			"dynatrace_automation_controller_connections":   resources.NewGeneric(export.ResourceTypes.AutomationControllerConnections).Resource(),
			"dynatrace_event_driven_ansible_connections":    resources.NewGeneric(export.ResourceTypes.EventDrivenAnsibleConnections).Resource(),
		},
		ConfigureContextFunc: config.ProviderConfigure,
	}

	incubator(prv)
	return prv
}
