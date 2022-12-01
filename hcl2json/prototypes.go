package hcl2json

import (
	"github.com/dtcookie/dynatrace/api/config/anomalies/applications"
	"github.com/dtcookie/dynatrace/api/config/anomalies/databaseservices"
	"github.com/dtcookie/dynatrace/api/config/anomalies/diskevents"
	"github.com/dtcookie/dynatrace/api/config/anomalies/hosts"
	"github.com/dtcookie/dynatrace/api/config/anomalies/metricevents"
	services "github.com/dtcookie/dynatrace/api/config/anomalies/services"
	"github.com/dtcookie/dynatrace/api/config/applications/mobile"
	"github.com/dtcookie/dynatrace/api/config/applications/web"
	"github.com/dtcookie/dynatrace/api/config/applications/web/applicationdetectionrules"
	"github.com/dtcookie/dynatrace/api/config/autotags"
	"github.com/dtcookie/dynatrace/api/config/credentials/aws"
	"github.com/dtcookie/dynatrace/api/config/credentials/azure"
	"github.com/dtcookie/dynatrace/api/config/credentials/cloudfoundry"
	"github.com/dtcookie/dynatrace/api/config/credentials/kubernetes"
	"github.com/dtcookie/dynatrace/api/config/credentials/vault"
	"github.com/dtcookie/dynatrace/api/config/customservices"
	"github.com/dtcookie/dynatrace/api/config/dashboards"
	"github.com/dtcookie/dynatrace/api/config/dashboards/sharing"
	"github.com/dtcookie/dynatrace/api/config/managementzones"
	"github.com/dtcookie/dynatrace/api/config/metrics/calculated/service"
	hostNaming "github.com/dtcookie/dynatrace/api/config/naming/hosts"
	processGroupNaming "github.com/dtcookie/dynatrace/api/config/naming/processgroups"
	serviceNaming "github.com/dtcookie/dynatrace/api/config/naming/services"
	"github.com/dtcookie/dynatrace/api/config/requestattributes"
	"github.com/dtcookie/dynatrace/api/config/requestnaming"
	"github.com/dtcookie/dynatrace/api/config/synthetic/locations"
	"github.com/dtcookie/dynatrace/api/config/synthetic/monitors"
	"github.com/dtcookie/dynatrace/api/config/v2/alerting"
	"github.com/dtcookie/dynatrace/api/config/v2/anomalies/frequentissues"
	"github.com/dtcookie/dynatrace/api/config/v2/ibmmq/filters"
	"github.com/dtcookie/dynatrace/api/config/v2/ibmmq/imsbridges"
	"github.com/dtcookie/dynatrace/api/config/v2/ibmmq/queuemanagers"
	"github.com/dtcookie/dynatrace/api/config/v2/ibmmq/queuesharinggroups"
	"github.com/dtcookie/dynatrace/api/config/v2/keyrequests"
	"github.com/dtcookie/dynatrace/api/config/v2/maintenance"
	"github.com/dtcookie/dynatrace/api/config/v2/networkzones"
	"github.com/dtcookie/dynatrace/api/config/v2/notifications"
	"github.com/dtcookie/dynatrace/api/config/v2/slo"
	"github.com/dtcookie/dynatrace/api/config/v2/spans/attributes"
	"github.com/dtcookie/dynatrace/api/config/v2/spans/capture"
	"github.com/dtcookie/dynatrace/api/config/v2/spans/ctxprop"
	"github.com/dtcookie/dynatrace/api/config/v2/spans/entrypoints"
	"github.com/dtcookie/dynatrace/api/config/v2/spans/resattr"
)

var protoTypes = map[string]interface{}{
	"dynatrace_dashboard":                  new(dashboards.Dashboard),
	"dynatrace_alerting":                   new(alerting.Profile),
	"dynatrace_custom_service":             new(customservices.CustomService),
	"dynatrace_management_zone":            new(managementzones.ManagementZone),
	"dynatrace_maintenance":                new(maintenance.MaintenanceWindow),
	"dynatrace_request_attribute":          new(requestattributes.RequestAttribute),
	"dynatrace_autotag":                    new(autotags.AutoTag),
	"dynatrace_aws_credentials":            new(aws.AWSCredentialsConfig),
	"dynatrace_azure_credentials":          new(azure.AzureCredentials),
	"dynatrace_k8s_credentials":            new(kubernetes.KubernetesCredentials),
	"dynatrace_cloudfoundry_credentials":   new(cloudfoundry.CloudFoundryCredentials),
	"dynatrace_service_anomalies":          new(services.AnomalyDetection),
	"dynatrace_application_anomalies":      new(applications.AnomalyDetection),
	"dynatrace_host_anomalies":             new(hosts.AnomalyDetection),
	"dynatrace_database_anomalies":         new(databaseservices.AnomalyDetection),
	"dynatrace_custom_anomalies":           new(metricevents.MetricEvent),
	"dynatrace_disk_anomalies":             new(diskevents.AnomalyDetection),
	"dynatrace_calculated_service_metric":  new(service.CalculatedServiceMetric),
	"dynatrace_service_naming":             new(serviceNaming.NamingRule),
	"dynatrace_host_naming":                new(hostNaming.NamingRule),
	"dynatrace_processgroup_naming":        new(processGroupNaming.NamingRule),
	"dynatrace_slo":                        new(slo.SLO),
	"dynatrace_span_entry_point":           new(entrypoints.SpanEntryPoint),
	"dynatrace_span_capture_rule":          new(capture.SpanCaptureSetting),
	"dynatrace_span_context_propagation":   new(ctxprop.PropagationSetting),
	"dynatrace_resource_attributes":        new(resattr.ResourceAttributes),
	"dynatrace_span_attribute":             new(attributes.SpanAttribute),
	"dynatrace_dashboard_sharing":          new(sharing.DashboardSharing),
	"dynatrace_mobile_application":         new(mobile.NewAppConfig),
	"dynatrace_browser_monitor":            new(monitors.BrowserSyntheticMonitorUpdate),
	"dynatrace_http_monitor":               new(monitors.HTTPSyntheticMonitorUpdate),
	"dynatrace_web_application":            new(web.ApplicationConfig),
	"dynatrace_application_data_privacy":   new(web.ApplicationDataPrivacy),
	"dynatrace_application_error_rules":    new(web.ApplicationErrorRules),
	"dynatrace_request_naming":             new(requestnaming.RequestNaming),
	"dynatrace_request_namings":            new(requestnaming.Order),
	"dynatrace_key_requests":               new(keyrequests.KeyRequest),
	"dynatrace_queue_manager":              new(queuemanagers.QueueManager),
	"dynatrace_ibm_mq_filters":             new(filters.Filters),
	"dynatrace_queue_sharing_groups":       new(queuesharinggroups.QueueSharingGroup),
	"dynatrace_ims_bridges":                new(imsbridges.IMSBridge),
	"dynatrace_network_zones":              new(networkzones.NetworkZones),
	"dynatrace_application_detection_rule": new(applicationdetectionrules.ApplicationDetectionRule),
	"dynatrace_frequent_issues":            new(frequentissues.FrequentIssues),
	"dynatrace_ansible_tower_notification": new(notifications.AnsibleTower),
	"dynatrace_email_notification":         new(notifications.Email),
	"dynatrace_jira_notification":          new(notifications.Jira),
	"dynatrace_ops_genie_notification":     new(notifications.OpsGenie),
	"dynatrace_pager_duty_notification":    new(notifications.PagerDuty),
	"dynatrace_service_now_notification":   new(notifications.ServiceNow),
	"dynatrace_slack_notification":         new(notifications.Slack),
	"dynatrace_trello_notification":        new(notifications.Trello),
	"dynatrace_victor_ops_notification":    new(notifications.VictorOps),
	"dynatrace_webhook_notification":       new(notifications.WebHook),
	"dynatrace_xmatters_notification":      new(notifications.XMatters),
	"dynatrace_credentials":                new(vault.Credentials),
	"dynatrace_synthetic_location":         new(locations.PrivateSyntheticLocation),
}
