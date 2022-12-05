package download

import (
	"reflect"

	"github.com/dtcookie/dynatrace/api/config/anomalies/applications"
	"github.com/dtcookie/dynatrace/api/config/anomalies/databaseservices"
	"github.com/dtcookie/dynatrace/api/config/anomalies/diskevents"
	"github.com/dtcookie/dynatrace/api/config/anomalies/hosts"
	"github.com/dtcookie/dynatrace/api/config/anomalies/metricevents"
	"github.com/dtcookie/dynatrace/api/config/anomalies/services"
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
	"github.com/dtcookie/dynatrace/api/config/managementzones"
	"github.com/dtcookie/dynatrace/api/config/metrics/calculated/service"
	hostnaming "github.com/dtcookie/dynatrace/api/config/naming/hosts"
	processgroupnaming "github.com/dtcookie/dynatrace/api/config/naming/processgroups"
	servicenaming "github.com/dtcookie/dynatrace/api/config/naming/services"
	"github.com/dtcookie/dynatrace/api/config/requestattributes"
	"github.com/dtcookie/dynatrace/api/config/requestnaming"
	privlocations "github.com/dtcookie/dynatrace/api/config/synthetic/locations"
	"github.com/dtcookie/dynatrace/api/config/synthetic/monitors"
	servicetopology "github.com/dtcookie/dynatrace/api/config/topology/service"
	"github.com/dtcookie/dynatrace/api/config/v2/alerting"
	"github.com/dtcookie/dynatrace/api/config/v2/anomalies/frequentissues"
	"github.com/dtcookie/dynatrace/api/config/v2/ibmmq/filters"
	"github.com/dtcookie/dynatrace/api/config/v2/ibmmq/imsbridges"
	"github.com/dtcookie/dynatrace/api/config/v2/ibmmq/queuemanagers"
	"github.com/dtcookie/dynatrace/api/config/v2/ibmmq/queuesharinggroups"
	v2maintenance "github.com/dtcookie/dynatrace/api/config/v2/maintenance"
	"github.com/dtcookie/dynatrace/api/config/v2/networkzones"
	"github.com/dtcookie/dynatrace/api/config/v2/notifications"
	"github.com/dtcookie/dynatrace/api/config/v2/slo"
	"github.com/dtcookie/dynatrace/api/config/v2/spans/attributes"
	"github.com/dtcookie/dynatrace/api/config/v2/spans/capture"
	"github.com/dtcookie/dynatrace/api/config/v2/spans/ctxprop"
	"github.com/dtcookie/dynatrace/api/config/v2/spans/entrypoints"
	"github.com/dtcookie/dynatrace/api/config/v2/spans/resattr"
	"github.com/dtcookie/opt"
)

var ResourceInfoMap = map[string]ResourceStruct{
	"dynatrace_alerting": {
		RESTClient: func(environmentURL, apiToken string) []StandardClient {
			clients := []StandardClient{alerting.NewService(environmentURL+"/api/v2", apiToken)}
			return clients
		},
		HardcodedIds: []string{"dynatrace_management_zone"},
		DsReplaceIds: func(resources Resources, dsData DataSourceData) map[string][]*ReplacedID {
			return Replace(resources, "dynatrace_management_zone", dsData, ReplacedID{RefDS: "dynatrace_management_zone", RefRes: "dynatrace_management_zone"})
		},
	},
	"dynatrace_ansible_tower_notification": {
		RESTClient: func(environmentURL, apiToken string) []StandardClient {
			clients := []StandardClient{notifications.NewAnsibleTowerService(environmentURL+"/api/v2", apiToken)}
			return clients
		},
		CustomName: func(dlConfig DownloadConfig, resourceName string, v interface{}, counter NameCounter) string {
			return counter.Numbering(v.(*notifications.Notification).Name)
		},
		HardcodedIds: []string{"dynatrace_alerting_profile"},
		DsReplaceIds: func(resources Resources, dsData DataSourceData) map[string][]*ReplacedID {
			return replaceIdsNotif(resources, dsData)
		},
	},
	"dynatrace_application_anomalies": {
		NoListClient: func(environmentURL, apiToken string) NoListClient {
			return applications.NewService(environmentURL+"/api/config/v1", apiToken)
		},
		CustomName: func(dlConfig DownloadConfig, resourceName string, v interface{}, counter NameCounter) string {
			return resourceName
		},
	},
	"dynatrace_application_data_privacy": {
		RESTClient: func(environmentURL, apiToken string) []StandardClient {
			clients := []StandardClient{web.NewAppDataPrivacyService(environmentURL+"/api/config/v1", apiToken)}
			return clients
		},
		CustomName: func(dlConfig DownloadConfig, resourceName string, v interface{}, counter NameCounter) string {
			client := web.NewService(dlConfig.EnvironmentURL+"/api/config/v1", dlConfig.APIToken)
			stubList, err := client.List()
			if err != nil {
				return ""
			}
			for _, stub := range stubList.Values {
				if stub.ID == *v.(*web.ApplicationDataPrivacy).WebApplicationID {
					return stub.Name
				}
			}
			return ""
		},
		HardcodedIds: []string{"dynatrace_application"},
		DsReplaceIds: func(resources Resources, dsData DataSourceData) map[string][]*ReplacedID {
			return Replace(resources, "dynatrace_application", dsData, ReplacedID{RefDS: "dynatrace_application", RefRes: "dynatrace_web_application"})
		},
	},
	"dynatrace_application_detection_rule": {
		RESTClient: func(environmentURL, apiToken string) []StandardClient {
			clients := []StandardClient{applicationdetectionrules.NewService(environmentURL+"/api/config/v1", apiToken)}
			return clients
		},
		CustomName: func(dlConfig DownloadConfig, resourceName string, v interface{}, counter NameCounter) string {
			client := web.NewService(dlConfig.EnvironmentURL+"/api/config/v1", dlConfig.APIToken)
			stubList, err := client.List()
			if err != nil {
				return ""
			}
			for _, stub := range stubList.Values {
				if stub.ID == v.(*applicationdetectionrules.ApplicationDetectionRule).ApplicationIdentifier {
					return counter.Numbering(stub.Name)
				}
			}
			return ""
		},
		HardcodedIds: []string{"dynatrace_application"},
		DsReplaceIds: func(resources Resources, dsData DataSourceData) map[string][]*ReplacedID {
			return Replace(resources, "dynatrace_application", dsData, ReplacedID{RefDS: "dynatrace_application", RefRes: "dynatrace_web_application"})
		},
	},
	"dynatrace_application_error_rules": {
		RESTClient: func(environmentURL, apiToken string) []StandardClient {
			clients := []StandardClient{web.NewErrorRulesService(environmentURL+"/api/config/v1", apiToken)}
			return clients
		},
		CustomName: func(dlConfig DownloadConfig, resourceName string, v interface{}, counter NameCounter) string {
			client := web.NewService(dlConfig.EnvironmentURL+"/api/config/v1", dlConfig.APIToken)
			stubList, err := client.List()
			if err != nil {
				return ""
			}
			for _, stub := range stubList.Values {
				if stub.ID == v.(*web.ApplicationErrorRules).WebApplicationID {
					return stub.Name
				}
			}
			return ""
		},
		HardcodedIds: []string{"dynatrace_application"},
		DsReplaceIds: func(resources Resources, dsData DataSourceData) map[string][]*ReplacedID {
			return Replace(resources, "dynatrace_application", dsData, ReplacedID{RefDS: "dynatrace_application", RefRes: "dynatrace_web_application"})
		},
	},
	"dynatrace_autotag": {
		RESTClient: func(environmentURL, apiToken string) []StandardClient {
			clients := []StandardClient{autotags.NewService(environmentURL+"/api/config/v1", apiToken)}
			return clients
		},
	},
	"dynatrace_credentials": {
		RESTClient: func(environmentURL, apiToken string) []StandardClient {
			clients := []StandardClient{vault.NewService(environmentURL+"/api/config/v1", apiToken)}
			return clients
		},
		CustomName: func(dlConfig DownloadConfig, resourceName string, v interface{}, counter NameCounter) string {
			return counter.Numbering(v.(*vault.Credentials).Name)
		},
		HardcodedIds: []string{"dynatrace_credentials"},
		DsReplaceIds: func(resources Resources, dsData DataSourceData) map[string][]*ReplacedID {
			ids := map[string][]*ReplacedID{}
			for _, resource := range resources {
				dataObj := resource.RESTObject.(*vault.Credentials)
				switch dataObj.Type {
				case vault.CredentialsTypes.Certificate, vault.CredentialsTypes.PublicCertificate:
					dataObj.Certificate = opt.NewString("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")
					dataObj.CertificateFormat = vault.CertificateFormats.Pkcs12.Ref()
					dataObj.Password = opt.NewString("redacted")
				case vault.CredentialsTypes.Token:
					dataObj.Token = opt.NewString("redacted")
				case vault.CredentialsTypes.UsernamePassword:
					dataObj.Username = opt.NewString("redacted")
					dataObj.Password = opt.NewString("redacted")
				}
				ids = Replace(resources, "dynatrace_credentials", dsData, ReplacedID{RefDS: "dynatrace_credentials", RefRes: "dynatrace_credentials"})
			}
			return ids
		},
	},
	"dynatrace_synthetic_location": {
		RESTClient: func(environmentURL, apiToken string) []StandardClient {
			clients := []StandardClient{privlocations.NewService(environmentURL+"/api/v1", apiToken)}
			return clients
		},
		CustomName: func(dlConfig DownloadConfig, resourceName string, v interface{}, counter NameCounter) string {
			return counter.Numbering(v.(*privlocations.PrivateSyntheticLocation).Name)
		},
	},
	"dynatrace_aws_credentials": {
		RESTClient: func(environmentURL, apiToken string) []StandardClient {
			clients := []StandardClient{aws.NewService(environmentURL+"/api/config/v1", apiToken)}
			return clients
		},
		CustomName: func(dlConfig DownloadConfig, resourceName string, v interface{}, counter NameCounter) string {
			return counter.Numbering(v.(*aws.AWSCredentialsConfig).Label)
		},
		HardcodedIds: []string{"dynatrace_aws_credentials"},
		DsReplaceIds: func(resources Resources, dsData DataSourceData) map[string][]*ReplacedID {
			ids := map[string][]*ReplacedID{}
			for _, resource := range resources {
				dataObj := resource.RESTObject.(*aws.AWSCredentialsConfig)
				if dataObj.AuthenticationData != nil && dataObj.AuthenticationData.KeyBasedAuthentication != nil {
					dataObj.AuthenticationData.KeyBasedAuthentication.SecretKey = opt.NewString("redacted")
				}
			}
			return ids
		},
	},
	"dynatrace_azure_credentials": {
		RESTClient: func(environmentURL, apiToken string) []StandardClient {
			clients := []StandardClient{azure.NewService(environmentURL+"/api/config/v1", apiToken)}
			return clients
		},
		CustomName: func(dlConfig DownloadConfig, resourceName string, v interface{}, counter NameCounter) string {
			return counter.Numbering(v.(*azure.AzureCredentials).Label)
		},
	},
	"dynatrace_browser_monitor": {
		RESTClient: func(environmentURL, apiToken string) []StandardClient {
			clients := []StandardClient{monitors.NewBrowserService(environmentURL+"/api/v1", apiToken)}
			return clients
		},
		CustomName: func(dlConfig DownloadConfig, resourceName string, v interface{}, counter NameCounter) string {
			return counter.Numbering(v.(*monitors.BrowserSyntheticMonitorUpdate).Name)
		},
		HardcodedIds: []string{"dynatrace_application", "dynatrace_synthetic_location", "dynatrace_credentials"},
		DsReplaceIds: func(resources Resources, dsData DataSourceData) map[string][]*ReplacedID {
			ids := map[string][]*ReplacedID{}
			for k, v := range Replace(resources, "dynatrace_application", dsData, ReplacedID{RefDS: "dynatrace_application", RefRes: "dynatrace_web_application"}) {
				ids[k] = v
			}
			for k, v := range Replace(resources, "dynatrace_synthetic_location", dsData, ReplacedID{RefDS: "dynatrace_synthetic_location", RefRes: "dynatrace_synthetic_location", Processed: true}) {
				ids[k] = v
			}
			for k, v := range Replace(resources, "dynatrace_credentials", dsData, ReplacedID{RefDS: "dynatrace_credentials", RefRes: "dynatrace_credentials"}) {
				ids[k] = v
			}
			return ids
		},
	},
	"dynatrace_calculated_service_metric": {
		RESTClient: func(environmentURL, apiToken string) []StandardClient {
			clients := []StandardClient{service.NewService(environmentURL+"/api/config/v1", apiToken)}
			return clients
		},
		HardcodedIds: []string{"dynatrace_request_attribute", "dynatrace_management_zone"},
		DsReplaceIds: func(resources Resources, dsData DataSourceData) map[string][]*ReplacedID {
			ids := map[string][]*ReplacedID{}
			for k, v := range Replace(resources, "dynatrace_request_attribute", dsData, ReplacedID{RefDS: "dynatrace_request_attribute", RefRes: "dynatrace_request_attribute"}) {
				ids[k] = v
			}
			for k, v := range Replace(resources, "dynatrace_management_zone", dsData, ReplacedID{RefDS: "dynatrace_management_zone", RefRes: "dynatrace_management_zone"}, NameReplace) {
				ids[k] = v
			}
			return ids
		},
	},
	"dynatrace_cloudfoundry_credentials": {
		RESTClient: func(environmentURL, apiToken string) []StandardClient {
			clients := []StandardClient{cloudfoundry.NewService(environmentURL+"/api/config/v1", apiToken)}
			return clients
		},
	},
	"dynatrace_custom_anomalies": {
		RESTClient: func(environmentURL, apiToken string) []StandardClient {
			clients := []StandardClient{metricevents.NewService(environmentURL+"/api/config/v1", apiToken)}
			return clients
		},
	},
	"dynatrace_custom_service": {
		RESTClient: func(environmentURL, apiToken string) []StandardClient {
			clients := []StandardClient{
				customservices.NewDotNetService(environmentURL+"/api/config/v1", apiToken),
				customservices.NewGoService(environmentURL+"/api/config/v1", apiToken),
				customservices.NewJavaService(environmentURL+"/api/config/v1", apiToken),
				customservices.NewNodeJSService(environmentURL+"/api/config/v1", apiToken),
				customservices.NewPHPService(environmentURL+"/api/config/v1", apiToken)}
			return clients
		},
	},
	"dynatrace_dashboard": {
		RESTClient: func(environmentURL, apiToken string) []StandardClient {
			clients := []StandardClient{dashboards.NewService(environmentURL+"/api/config/v1", apiToken)}
			return clients
		},
		CustomName: func(dlConfig DownloadConfig, resourceName string, v interface{}, counter NameCounter) string {
			return counter.Numbering(v.(*dashboards.Dashboard).Metadata.Name)
		},
		HardcodedIds: []string{"dynatrace_management_zone"},
		DsReplaceIds: func(resources Resources, dsData DataSourceData) map[string][]*ReplacedID {
			return Replace(resources, "dynatrace_management_zone", dsData, ReplacedID{RefDS: "dynatrace_management_zone", RefRes: "dynatrace_management_zone"})
		},
	},
	"dynatrace_database_anomalies": {
		NoListClient: func(environmentURL, apiToken string) NoListClient {
			return databaseservices.NewService(environmentURL+"/api/config/v1", apiToken)
		},
		CustomName: func(dlConfig DownloadConfig, resourceName string, v interface{}, counter NameCounter) string {
			return resourceName
		},
	},
	"dynatrace_disk_anomalies": {
		RESTClient: func(environmentURL, apiToken string) []StandardClient {
			clients := []StandardClient{diskevents.NewService(environmentURL+"/api/config/v1", apiToken)}
			return clients
		},
	},
	"dynatrace_email_notification": {
		RESTClient: func(environmentURL, apiToken string) []StandardClient {
			clients := []StandardClient{notifications.NewEmailService(environmentURL+"/api/v2", apiToken)}
			return clients
		},
		CustomName: func(dlConfig DownloadConfig, resourceName string, v interface{}, counter NameCounter) string {
			return counter.Numbering(v.(*notifications.Notification).Name)
		},
		HardcodedIds: []string{"dynatrace_alerting_profile"},
		DsReplaceIds: func(resources Resources, dsData DataSourceData) map[string][]*ReplacedID {
			return replaceIdsNotif(resources, dsData)
		},
	},
	"dynatrace_frequent_issues": {
		RESTClient: func(environmentURL, apiToken string) []StandardClient {
			clients := []StandardClient{frequentissues.NewService(environmentURL+"/api/v2", apiToken)}
			return clients
		},
		CustomName: func(dlConfig DownloadConfig, resourceName string, v interface{}, counter NameCounter) string {
			return resourceName
		},
	},
	"dynatrace_host_anomalies": {
		NoListClient: func(environmentURL, apiToken string) NoListClient {
			return hosts.NewService(environmentURL+"/api/config/v1", apiToken)
		},
		CustomName: func(dlConfig DownloadConfig, resourceName string, v interface{}, counter NameCounter) string {
			return resourceName
		},
	},
	"dynatrace_http_monitor": {
		RESTClient: func(environmentURL, apiToken string) []StandardClient {
			clients := []StandardClient{monitors.NewHTTPService(environmentURL+"/api/v1", apiToken)}
			return clients
		},
		CustomName: func(dlConfig DownloadConfig, resourceName string, v interface{}, counter NameCounter) string {
			return counter.Numbering(v.(*monitors.HTTPSyntheticMonitorUpdate).Name)
		},
		HardcodedIds: []string{"dynatrace_application", "dynatrace_synthetic_location", "dynatrace_credentials"},
		DsReplaceIds: func(resources Resources, dsData DataSourceData) map[string][]*ReplacedID {
			ids := map[string][]*ReplacedID{}
			for k, v := range Replace(resources, "dynatrace_credentials", dsData, ReplacedID{RefDS: "dynatrace_credentials"}) {
				ids[k] = v
			}
			for k, v := range Replace(resources, "dynatrace_application", dsData, ReplacedID{RefDS: "dynatrace_application", RefRes: "dynatrace_web_application"}) {
				ids[k] = v
			}
			for k, v := range Replace(resources, "dynatrace_synthetic_location", dsData, ReplacedID{RefDS: "dynatrace_synthetic_location", Processed: true}) {
				ids[k] = v
			}
			return ids
		},
	},
	"dynatrace_host_naming": {
		RESTClient: func(environmentURL, apiToken string) []StandardClient {
			clients := []StandardClient{hostnaming.NewService(environmentURL+"/api/config/v1", apiToken)}
			return clients
		},
	},
	"dynatrace_ibm_mq_filters": {
		RESTClient: func(environmentURL, apiToken string) []StandardClient {
			clients := []StandardClient{filters.NewService(environmentURL+"/api/v2", apiToken)}
			return clients
		},
		CustomName: func(dlConfig DownloadConfig, resourceName string, v interface{}, counter NameCounter) string {
			return resourceName
		},
	},
	"dynatrace_ims_bridges": {
		RESTClient: func(environmentURL, apiToken string) []StandardClient {
			clients := []StandardClient{imsbridges.NewService(environmentURL+"/api/v2", apiToken)}
			return clients
		},
	},
	"dynatrace_jira_notification": {
		RESTClient: func(environmentURL, apiToken string) []StandardClient {
			clients := []StandardClient{notifications.NewJiraService(environmentURL+"/api/v2", apiToken)}
			return clients
		},
		CustomName: func(dlConfig DownloadConfig, resourceName string, v interface{}, counter NameCounter) string {
			return counter.Numbering(v.(*notifications.Notification).Name)
		},
		HardcodedIds: []string{"dynatrace_alerting_profile"},
		DsReplaceIds: func(resources Resources, dsData DataSourceData) map[string][]*ReplacedID {
			return replaceIdsNotif(resources, dsData)
		},
	},
	"dynatrace_key_requests": {
		CustomName: func(dlConfig DownloadConfig, resourceName string, v interface{}, counter NameCounter) string {
			return "for " + v.(servicetopology.Service).DisplayName
		},
	},
	"dynatrace_k8s_credentials": {
		RESTClient: func(environmentURL, apiToken string) []StandardClient {
			clients := []StandardClient{kubernetes.NewService(environmentURL+"/api/config/v1", apiToken)}
			return clients
		},
		CustomName: func(dlConfig DownloadConfig, resourceName string, v interface{}, counter NameCounter) string {
			return counter.Numbering(v.(*kubernetes.KubernetesCredentials).Label)
		},
	},
	"dynatrace_maintenance": {
		RESTClient: func(environmentURL, apiToken string) []StandardClient {
			clients := []StandardClient{v2maintenance.NewService(environmentURL+"/api/v2", apiToken)}
			return clients
		},
		CustomName: func(dlConfig DownloadConfig, resourceName string, v interface{}, counter NameCounter) string {
			return counter.Numbering(v.(*v2maintenance.MaintenanceWindow).GeneralProperties.Name)
		},
		HardcodedIds: []string{"dynatrace_management_zone"},
		DsReplaceIds: func(resources Resources, dsData DataSourceData) map[string][]*ReplacedID {
			return Replace(resources, "dynatrace_management_zone", dsData, ReplacedID{RefDS: "dynatrace_management_zone", RefRes: "dynatrace_management_zone"}, Settings20MatchReplace)
		},
	},
	"dynatrace_management_zone": {
		RESTClient: func(environmentURL, apiToken string) []StandardClient {
			clients := []StandardClient{managementzones.NewService(environmentURL+"/api/config/v1", apiToken)}
			return clients
		},
	},
	"dynatrace_mobile_application": {
		RESTClient: func(environmentURL, apiToken string) []StandardClient {
			clients := []StandardClient{mobile.NewService(environmentURL+"/api/config/v1", apiToken)}
			return clients
		},
		HardcodedIds: []string{"dynatrace_request_attribute"},
		DsReplaceIds: func(resources Resources, dsData DataSourceData) map[string][]*ReplacedID {
			return Replace(resources, "dynatrace_request_attribute", dsData, ReplacedID{RefDS: "dynatrace_request_attribute"})
		},
	},
	"dynatrace_network_zones": {
		RESTClient: func(environmentURL, apiToken string) []StandardClient {
			clients := []StandardClient{networkzones.NewService(environmentURL+"/api/v2", apiToken)}
			return clients
		},
		CustomName: func(dlConfig DownloadConfig, resourceName string, v interface{}, counter NameCounter) string {
			return resourceName
		},
	},
	"dynatrace_ops_genie_notification": {
		RESTClient: func(environmentURL, apiToken string) []StandardClient {
			clients := []StandardClient{notifications.NewOpsgenieService(environmentURL+"/api/v2", apiToken)}
			return clients
		},
		CustomName: func(dlConfig DownloadConfig, resourceName string, v interface{}, counter NameCounter) string {
			return counter.Numbering(v.(*notifications.Notification).Name)
		},
		HardcodedIds: []string{"dynatrace_alerting_profile"},
		DsReplaceIds: func(resources Resources, dsData DataSourceData) map[string][]*ReplacedID {
			return replaceIdsNotif(resources, dsData)
		},
	},
	"dynatrace_pager_duty_notification": {
		RESTClient: func(environmentURL, apiToken string) []StandardClient {
			clients := []StandardClient{notifications.NewPagerDutyService(environmentURL+"/api/v2", apiToken)}
			return clients
		},
		CustomName: func(dlConfig DownloadConfig, resourceName string, v interface{}, counter NameCounter) string {
			return counter.Numbering(v.(*notifications.Notification).Name)
		},
		HardcodedIds: []string{"dynatrace_alerting_profile"},
		DsReplaceIds: func(resources Resources, dsData DataSourceData) map[string][]*ReplacedID {
			return replaceIdsNotif(resources, dsData)
		},
	},
	"dynatrace_processgroup_naming": {
		RESTClient: func(environmentURL, apiToken string) []StandardClient {
			clients := []StandardClient{processgroupnaming.NewService(environmentURL+"/api/config/v1", apiToken)}
			return clients
		},
	},
	"dynatrace_queue_manager": {
		RESTClient: func(environmentURL, apiToken string) []StandardClient {
			clients := []StandardClient{queuemanagers.NewService(environmentURL+"/api/v2", apiToken)}
			return clients
		},
	},
	"dynatrace_queue_sharing_groups": {
		RESTClient: func(environmentURL, apiToken string) []StandardClient {
			clients := []StandardClient{queuesharinggroups.NewService(environmentURL+"/api/v2", apiToken)}
			return clients
		},
	},
	"dynatrace_request_attribute": {
		RESTClient: func(environmentURL, apiToken string) []StandardClient {
			clients := []StandardClient{requestattributes.NewService(environmentURL+"/api/config/v1", apiToken)}
			return clients
		},
	},
	"dynatrace_request_naming": {
		RESTClient: func(environmentURL, apiToken string) []StandardClient {
			clients := []StandardClient{requestnaming.NewService(environmentURL+"/api/config/v1", apiToken)}
			return clients
		},
		CustomName: func(dlConfig DownloadConfig, resourceName string, v interface{}, counter NameCounter) string {
			return counter.Numbering(v.(*requestnaming.RequestNaming).NamingPattern)
		},
	},
	"dynatrace_resource_attributes": {
		RESTClient: func(environmentURL, apiToken string) []StandardClient {
			clients := []StandardClient{resattr.NewService(environmentURL+"/api/v2", apiToken)}
			return clients
		},
		CustomName: func(dlConfig DownloadConfig, resourceName string, v interface{}, counter NameCounter) string {
			return resourceName
		},
	},
	"dynatrace_service_anomalies": {
		NoListClient: func(environmentURL, apiToken string) NoListClient {
			return services.NewService(environmentURL+"/api/config/v1", apiToken)
		},
		CustomName: func(dlConfig DownloadConfig, resourceName string, v interface{}, counter NameCounter) string {
			return resourceName
		},
	},
	"dynatrace_service_naming": {
		RESTClient: func(environmentURL, apiToken string) []StandardClient {
			clients := []StandardClient{servicenaming.NewService(environmentURL+"/api/config/v1", apiToken)}
			return clients
		},
	},
	"dynatrace_service_now_notification": {
		RESTClient: func(environmentURL, apiToken string) []StandardClient {
			clients := []StandardClient{notifications.NewServiceNowService(environmentURL+"/api/v2", apiToken)}
			return clients
		},
		CustomName: func(dlConfig DownloadConfig, resourceName string, v interface{}, counter NameCounter) string {
			return counter.Numbering(v.(*notifications.Notification).Name)
		},
		HardcodedIds: []string{"dynatrace_alerting_profile"},
		DsReplaceIds: func(resources Resources, dsData DataSourceData) map[string][]*ReplacedID {
			return replaceIdsNotif(resources, dsData)
		},
	},
	"dynatrace_slack_notification": {
		RESTClient: func(environmentURL, apiToken string) []StandardClient {
			clients := []StandardClient{notifications.NewSlackService(environmentURL+"/api/v2", apiToken)}
			return clients
		},
		CustomName: func(dlConfig DownloadConfig, resourceName string, v interface{}, counter NameCounter) string {
			return counter.Numbering(v.(*notifications.Notification).Name)
		},
		HardcodedIds: []string{"dynatrace_alerting_profile"},
		DsReplaceIds: func(resources Resources, dsData DataSourceData) map[string][]*ReplacedID {
			return replaceIdsNotif(resources, dsData)
		},
	},
	"dynatrace_slo": {
		RESTClient: func(environmentURL, apiToken string) []StandardClient {
			clients := []StandardClient{slo.NewService(environmentURL+"/api/v2", apiToken)}
			return clients
		},
	},
	"dynatrace_span_attribute": {
		RESTClient: func(environmentURL, apiToken string) []StandardClient {
			clients := []StandardClient{attributes.NewService(environmentURL+"/api/v2", apiToken)}
			return clients
		},
		CustomName: func(dlConfig DownloadConfig, resourceName string, v interface{}, counter NameCounter) string {
			return counter.Numbering(v.(*attributes.SpanAttribute).Key)
		},
	},
	"dynatrace_span_capture_rule": {
		RESTClient: func(environmentURL, apiToken string) []StandardClient {
			clients := []StandardClient{capture.NewService(environmentURL+"/api/v2", apiToken)}
			return clients
		},
		CustomName: func(dlConfig DownloadConfig, resourceName string, v interface{}, counter NameCounter) string {
			return counter.Numbering(v.(*capture.SpanCaptureSetting).SpanCaptureRule.Name)
		},
	},
	"dynatrace_span_context_propagation": {
		RESTClient: func(environmentURL, apiToken string) []StandardClient {
			clients := []StandardClient{ctxprop.NewService(environmentURL+"/api/v2", apiToken)}
			return clients
		},
		CustomName: func(dlConfig DownloadConfig, resourceName string, v interface{}, counter NameCounter) string {
			return counter.Numbering(v.(*ctxprop.PropagationSetting).PropagationRule.Name)
		},
	},
	"dynatrace_span_entry_point": {
		RESTClient: func(environmentURL, apiToken string) []StandardClient {
			clients := []StandardClient{entrypoints.NewService(environmentURL+"/api/v2", apiToken)}
			return clients
		},
		CustomName: func(dlConfig DownloadConfig, resourceName string, v interface{}, counter NameCounter) string {
			return counter.Numbering(v.(*entrypoints.SpanEntryPoint).EntryPointRule.Name)
		},
	},
	"dynatrace_trello_notification": {
		RESTClient: func(environmentURL, apiToken string) []StandardClient {
			clients := []StandardClient{notifications.NewTrelloService(environmentURL+"/api/v2", apiToken)}
			return clients
		},
		CustomName: func(dlConfig DownloadConfig, resourceName string, v interface{}, counter NameCounter) string {
			return counter.Numbering(v.(*notifications.Notification).Name)
		},
		HardcodedIds: []string{"dynatrace_alerting_profile"},
		DsReplaceIds: func(resources Resources, dsData DataSourceData) map[string][]*ReplacedID {
			return replaceIdsNotif(resources, dsData)
		},
	},
	"dynatrace_victor_ops_notification": {
		RESTClient: func(environmentURL, apiToken string) []StandardClient {
			clients := []StandardClient{notifications.NewVictorOpsService(environmentURL+"/api/v2", apiToken)}
			return clients
		},
		CustomName: func(dlConfig DownloadConfig, resourceName string, v interface{}, counter NameCounter) string {
			return counter.Numbering(v.(*notifications.Notification).Name)
		},
		HardcodedIds: []string{"dynatrace_alerting_profile"},
		DsReplaceIds: func(resources Resources, dsData DataSourceData) map[string][]*ReplacedID {
			return replaceIdsNotif(resources, dsData)
		},
	},
	"dynatrace_web_application": {
		RESTClient: func(environmentURL, apiToken string) []StandardClient {
			clients := []StandardClient{web.NewService(environmentURL+"/api/config/v1", apiToken)}
			return clients
		},
		HardcodedIds: []string{"dynatrace_request_attribute"},
		DsReplaceIds: func(resources Resources, dsData DataSourceData) map[string][]*ReplacedID {
			return Replace(resources, "dynatrace_request_attribute", dsData, ReplacedID{RefDS: "dynatrace_request_attribute", RefRes: "dynatrace_request_attribute"})
		}},
	"dynatrace_webhook_notification": {
		RESTClient: func(environmentURL, apiToken string) []StandardClient {
			clients := []StandardClient{notifications.NewWebHookService(environmentURL+"/api/v2", apiToken)}
			return clients
		},
		CustomName: func(dlConfig DownloadConfig, resourceName string, v interface{}, counter NameCounter) string {
			return counter.Numbering(v.(*notifications.Notification).Name)
		},
		HardcodedIds: []string{"dynatrace_alerting_profile"},
		DsReplaceIds: func(resources Resources, dsData DataSourceData) map[string][]*ReplacedID {
			return replaceIdsNotif(resources, dsData)
		},
	},
	"dynatrace_xmatters_notification": {
		RESTClient: func(environmentURL, apiToken string) []StandardClient {
			clients := []StandardClient{notifications.NewXMattersService(environmentURL+"/api/v2", apiToken)}
			return clients
		},
		CustomName: func(dlConfig DownloadConfig, resourceName string, v interface{}, counter NameCounter) string {
			return counter.Numbering(v.(*notifications.Notification).Name)
		},
		HardcodedIds: []string{"dynatrace_alerting_profile"},
		DsReplaceIds: func(resources Resources, dsData DataSourceData) map[string][]*ReplacedID {
			return replaceIdsNotif(resources, dsData)
		},
	},
}

type ResourceStruct struct {
	RESTClient   func(environmentURL, apiToken string) []StandardClient
	NoListClient func(environmentURL, apiToken string) NoListClient
	CustomName   NameFunc
	HardcodedIds []string
	DsReplaceIds DataSourceReplaceFunc
}

type NameFunc func(DownloadConfig, string, interface{}, NameCounter) string

type DataSourceReplaceFunc func(Resources, DataSourceData) map[string][]*ReplacedID

func (me ResourceStruct) Name(dlConfig DownloadConfig, resourceName string, v interface{}, counter NameCounter) string {
	if v == nil {
		return ""
	}
	if me.CustomName != nil {
		return me.CustomName(dlConfig, resourceName, v, counter)
	}
	rv := reflect.ValueOf(v)
	switch rv.Type().Kind() {
	case reflect.Struct:
		field := rv.FieldByName("Name")
		if field.Type().Kind() == reflect.String {
			return counter.Numbering(field.Interface().(string))
		}
	case reflect.Ptr:
		return me.Name(dlConfig, resourceName, rv.Elem().Interface(), counter)
	}
	return ""
}

func replaceIdsNotif(resources Resources, dsData DataSourceData) map[string][]*ReplacedID {
	return Replace(resources, "dynatrace_alerting_profile", dsData, ReplacedID{RefDS: "dynatrace_alerting_profile", RefRes: "dynatrace_alerting"})
}
