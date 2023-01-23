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
	"github.com/dynatrace-oss/terraform-provider-dynatrace/datasources/credentials/vault"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/datasources/entities"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/datasources/entity"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/datasources/host"
	ds_iam_groups "github.com/dynatrace-oss/terraform-provider-dynatrace/datasources/iam/groups"
	metricsds "github.com/dynatrace-oss/terraform-provider-dynatrace/datasources/metrics/calculated/service"
	mgmzds "github.com/dynatrace-oss/terraform-provider-dynatrace/datasources/mgmz"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/datasources/process"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/datasources/processgroup"
	reqattrds "github.com/dynatrace-oss/terraform-provider-dynatrace/datasources/requestattributes"
	serviceds "github.com/dynatrace-oss/terraform-provider-dynatrace/datasources/service"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/datasources/synthetic/locations"
	v2alerting "github.com/dynatrace-oss/terraform-provider-dynatrace/datasources/v2alerting"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/export"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/provider/config"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/provider/logging"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/resources"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/resources/customtags"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/resources/environments"
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
			"dynatrace_alerting_profiles":         alerting.DataSource(),
			"dynatrace_alerting_profile":          v2alerting.DataSource(),
			"dynatrace_credentials":               vault.DataSource(),
			"dynatrace_synthetic_locations":       locations.DataSource(),
			"dynatrace_synthetic_location":        locations.UniqueDataSource(),
			"dynatrace_service":                   serviceds.DataSource(),
			"dynatrace_management_zone":           mgmzds.DataSource(),
			"dynatrace_application":               application.DataSource(),
			"dynatrace_host":                      host.DataSource(),
			"dynatrace_process":                   process.DataSource(),
			"dynatrace_process_group":             processgroup.DataSource(),
			"dynatrace_aws_iam_external":          dsaws.DataSource(),
			"dynatrace_request_attribute":         reqattrds.DataSource(),
			"dynatrace_calculated_service_metric": metricsds.DataSource(),
			"dynatrace_iam_group":                 ds_iam_groups.DataSource(),
			"dynatrace_entity":                    entity.DataSource(),
			"dynatrace_entities":                  entities.DataSource(),
		},
		ResourcesMap: map[string]*schema.Resource{
			"dynatrace_custom_service":             resources.NewGeneric(export.ResourceTypes.CustomService).Resource(),
			"dynatrace_dashboard":                  resources.NewGeneric(export.ResourceTypes.Dashboard).Resource(),
			"dynatrace_json_dashboard":             resources.NewGeneric(export.ResourceTypes.JSONDashboard).Resource(),
			"dynatrace_management_zone":            resources.NewGeneric(export.ResourceTypes.ManagementZone).Resource(),
			"dynatrace_management_zone_v2":         resources.NewGeneric(export.ResourceTypes.ManagementZoneV2).Resource(),
			"dynatrace_maintenance_window":         resources.NewGeneric(export.ResourceTypes.MaintenanceWindow).Resource(),
			"dynatrace_maintenance":                resources.NewGeneric(export.ResourceTypes.Maintenance).Resource(),
			"dynatrace_request_attribute":          resources.NewGeneric(export.ResourceTypes.RequestAttribute).Resource(),
			"dynatrace_alerting_profile":           resources.NewGeneric(export.ResourceTypes.AlertingProfile).Resource(),
			"dynatrace_alerting":                   resources.NewGeneric(export.ResourceTypes.Alerting).Resource(),
			"dynatrace_notification":               resources.NewGeneric(export.ResourceTypes.Notification).Resource(),
			"dynatrace_autotag":                    resources.NewGeneric(export.ResourceTypes.AutoTag).Resource(),
			"dynatrace_aws_credentials":            resources.NewGeneric(export.ResourceTypes.AWSCredentials).Resource(),
			"dynatrace_azure_credentials":          resources.NewGeneric(export.ResourceTypes.AzureCredentials).Resource(),
			"dynatrace_k8s_credentials":            resources.NewGeneric(export.ResourceTypes.KubernetesCredentials).Resource(),
			"dynatrace_cloudfoundry_credentials":   resources.NewGeneric(export.ResourceTypes.CloudFoundryCredentials).Resource(),
			"dynatrace_service_anomalies":          resources.NewGeneric(export.ResourceTypes.ServiceAnomalies).Resource(),
			"dynatrace_application_anomalies":      resources.NewGeneric(export.ResourceTypes.ApplicationAnomalies).Resource(),
			"dynatrace_host_anomalies":             resources.NewGeneric(export.ResourceTypes.HostAnomalies).Resource(),
			"dynatrace_database_anomalies":         resources.NewGeneric(export.ResourceTypes.DatabaseAnomalies).Resource(),
			"dynatrace_custom_anomalies":           resources.NewGeneric(export.ResourceTypes.CustomAnomalies).Resource(),
			"dynatrace_metric_events":              resources.NewGeneric(export.ResourceTypes.MetricEvents).Resource(),
			"dynatrace_disk_anomalies":             resources.NewGeneric(export.ResourceTypes.DiskEventAnomalies).Resource(),
			"dynatrace_calculated_service_metric":  resources.NewGeneric(export.ResourceTypes.CalculatedServiceMetric).Resource(),
			"dynatrace_service_naming":             resources.NewGeneric(export.ResourceTypes.ServiceNaming).Resource(),
			"dynatrace_host_naming":                resources.NewGeneric(export.ResourceTypes.HostNaming).Resource(),
			"dynatrace_processgroup_naming":        resources.NewGeneric(export.ResourceTypes.ProcessGroupNaming).Resource(),
			"dynatrace_slo":                        resources.NewGeneric(export.ResourceTypes.SLO).Resource(),
			"dynatrace_span_entry_point":           resources.NewGeneric(export.ResourceTypes.SpanEntryPoint).Resource(),
			"dynatrace_span_capture_rule":          resources.NewGeneric(export.ResourceTypes.SpanCaptureRule).Resource(),
			"dynatrace_span_context_propagation":   resources.NewGeneric(export.ResourceTypes.SpanContextPropagation).Resource(),
			"dynatrace_resource_attributes":        resources.NewGeneric(export.ResourceTypes.ResourceAttributes).Resource(),
			"dynatrace_span_attribute":             resources.NewGeneric(export.ResourceTypes.SpanAttribute).Resource(),
			"dynatrace_dashboard_sharing":          resources.NewGeneric(export.ResourceTypes.DashboardSharing).Resource(),
			"dynatrace_environment":                environments.Resource(),
			"dynatrace_mobile_application":         resources.NewGeneric(export.ResourceTypes.MobileApplication).Resource(),
			"dynatrace_browser_monitor":            resources.NewGeneric(export.ResourceTypes.BrowserMonitor).Resource(),
			"dynatrace_http_monitor":               resources.NewGeneric(export.ResourceTypes.HTTPMonitor).Resource(),
			"dynatrace_web_application":            resources.NewGeneric(export.ResourceTypes.WebApplication).Resource(),
			"dynatrace_application_data_privacy":   resources.NewGeneric(export.ResourceTypes.ApplicationDataPrivacy).Resource(),
			"dynatrace_application_error_rules":    resources.NewGeneric(export.ResourceTypes.ApplicationErrorRules).Resource(),
			"dynatrace_request_naming":             resources.NewGeneric(export.ResourceTypes.RequestNaming).Resource(),
			"dynatrace_request_namings":            resources.NewGeneric(export.ResourceTypes.RequestNamings).Resource(),
			"dynatrace_user_group":                 usergroups.Resource(),
			"dynatrace_user":                       users.Resource(),
			"dynatrace_key_requests":               resources.NewGeneric(export.ResourceTypes.KeyRequests).Resource(),
			"dynatrace_queue_manager":              resources.NewGeneric(export.ResourceTypes.QueueManager).Resource(),
			"dynatrace_ibm_mq_filters":             resources.NewGeneric(export.ResourceTypes.IBMMQFilters).Resource(),
			"dynatrace_queue_sharing_groups":       resources.NewGeneric(export.ResourceTypes.QueueSharingGroups).Resource(),
			"dynatrace_ims_bridges":                resources.NewGeneric(export.ResourceTypes.IMSBridge).Resource(),
			"dynatrace_network_zones":              resources.NewGeneric(export.ResourceTypes.NetworkZones).Resource(),
			"dynatrace_application_detection_rule": resources.NewGeneric(export.ResourceTypes.ApplicationDetection).Resource(),
			"dynatrace_frequent_issues":            resources.NewGeneric(export.ResourceTypes.FrequentIssues).Resource(),
			"dynatrace_ansible_tower_notification": resources.NewGeneric(export.ResourceTypes.AnsibleTowerNotification).Resource(),
			"dynatrace_email_notification":         resources.NewGeneric(export.ResourceTypes.EmailNotification).Resource(),
			"dynatrace_jira_notification":          resources.NewGeneric(export.ResourceTypes.JiraNotification).Resource(),
			"dynatrace_ops_genie_notification":     resources.NewGeneric(export.ResourceTypes.OpsGenieNotification).Resource(),
			"dynatrace_pager_duty_notification":    resources.NewGeneric(export.ResourceTypes.PagerDutyNotification).Resource(),
			"dynatrace_service_now_notification":   resources.NewGeneric(export.ResourceTypes.ServiceNowNotification).Resource(),
			"dynatrace_slack_notification":         resources.NewGeneric(export.ResourceTypes.SlackNotification).Resource(),
			"dynatrace_trello_notification":        resources.NewGeneric(export.ResourceTypes.TrelloNotification).Resource(),
			"dynatrace_victor_ops_notification":    resources.NewGeneric(export.ResourceTypes.VictorOpsNotification).Resource(),
			"dynatrace_webhook_notification":       resources.NewGeneric(export.ResourceTypes.WebHookNotification).Resource(),
			"dynatrace_xmatters_notification":      resources.NewGeneric(export.ResourceTypes.XMattersNotification).Resource(),
			"dynatrace_credentials":                resources.NewGeneric(export.ResourceTypes.Credentials).Resource(),
			"dynatrace_synthetic_location":         resources.NewGeneric(export.ResourceTypes.SyntheticLocation).Resource(),
			"dynatrace_network_zone":               resources.NewGeneric(export.ResourceTypes.NetworkZone).Resource(),
			"dynatrace_iam_user":                   resources.NewGeneric(export.ResourceTypes.IAMUser).Resource(),
			"dynatrace_iam_group":                  resources.NewGeneric(export.ResourceTypes.IAMGroup).Resource(),
			"dynatrace_api_token":                  resources.NewGeneric(export.ResourceTypes.APIToken).Resource(),
			"dynatrace_custom_tags":                customtags.Resource(),
		},
		ConfigureContextFunc: config.ProviderConfigure,
	}
}
