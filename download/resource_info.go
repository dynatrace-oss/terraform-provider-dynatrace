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
	"github.com/dtcookie/dynatrace/api/config/customservices"
	"github.com/dtcookie/dynatrace/api/config/dashboards"
	"github.com/dtcookie/dynatrace/api/config/managementzones"
	"github.com/dtcookie/dynatrace/api/config/metrics/calculated/service"
	hostnaming "github.com/dtcookie/dynatrace/api/config/naming/hosts"
	processgroupnaming "github.com/dtcookie/dynatrace/api/config/naming/processgroups"
	servicenaming "github.com/dtcookie/dynatrace/api/config/naming/services"
	"github.com/dtcookie/dynatrace/api/config/requestattributes"
	"github.com/dtcookie/dynatrace/api/config/requestnaming"
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
		DsReplaceIds: func(resources Resources, dsData DataSourceData) []ReplacedID {
			var ids = []ReplacedID{}
			for _, resource := range resources {
				dataObj := resource.RESTObject.(*alerting.Profile)
				dsName := "dynatrace_management_zone"
				for id, appInfo := range dsData[dsName].RESTMap {
					if dataObj.ManagementZone != nil && *dataObj.ManagementZone == id {
						dataObj.ManagementZone = opt.NewString("HCL-UNQUOTE-data.dynatrace_management_zone." + escape(appInfo["name"].(string)) + ".id")
						ids = append(ids, ReplacedID{id, dsName})
					}
				}
			}
			return ids
		},
		// ResReplaceIds: func(resName string, resourceData ResourceData) {
		// 	for _, resource := range resourceData[resName] {
		// 		dataObj := resource.RESTObject.(*alerting.Profile)
		// 		for _, resource := range resourceData["dynatrace_management_zone"] {
		// 			resourceObj := resource.RESTObject.(*managementzones.ManagementZone)
		// 			if dataObj.ManagementZone != nil && *dataObj.ManagementZone == *resourceObj.ID {
		// 				dataObj.ManagementZone = opt.NewString("HCL-UNQUOTE-dynatrace_management_zone." + escape(resourceObj.Name) + ".id")
		// 			}
		// 		}
		// 	}
		// },
	},
	"dynatrace_ansible_tower_notification": {
		RESTClient: func(environmentURL, apiToken string) []StandardClient {
			clients := []StandardClient{notifications.NewAnsibleTowerService(environmentURL+"/api/v2", apiToken)}
			return clients
		},
		CustomName: func(dlConfig DownloadConfig, resourceName string, v interface{}, counter NameCounter) string {
			return counter.Numbering(v.(*notifications.Notification).Name)
		},
		HardcodedIds: []string{"dynatrace_alerting"},
		DsReplaceIds: func(resources Resources, dsData DataSourceData) []ReplacedID {
			return replaceIdsNotif(resources, dsData)
		},
		// ResReplaceIds: func(resName string, resourceData ResourceData) {
		// 	for _, resource := range resourceData[resName] {
		// 		dataObj := resource.RESTObject.(*notifications.Notification)
		// 		for _, resource := range resourceData["dynatrace_alerting"] {
		// 			resourceObj := resource.RESTObject.(*alerting.Profile)
		// 			if dataObj.ProfileID != "" && dataObj.ProfileID == resourceObj.ID {
		// 				dataObj.ProfileID = "HCL-UNQUOTE-dynatrace_alerting." + escape(resourceObj.Name) + ".id"
		// 			}
		// 		}
		// 	}
		// },
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
		DsReplaceIds: func(resources Resources, dsData DataSourceData) []ReplacedID {
			var ids = []ReplacedID{}
			for _, resource := range resources {
				dataObj := resource.RESTObject.(*web.ApplicationDataPrivacy)
				dsName := "dynatrace_application"
				for id, appInfo := range dsData[dsName].RESTMap {
					if dataObj.WebApplicationID != nil && *dataObj.WebApplicationID == id {
						dataObj.WebApplicationID = opt.NewString("HCL-UNQUOTE-data.dynatrace_application." + escape(appInfo["name"].(string)) + ".id")
						ids = append(ids, ReplacedID{id, dsName})
					}
				}
			}
			return ids
		},
		// ResReplaceIds: func(resName string, resourceData ResourceData) {
		// 	for _, resource := range resourceData[resName] {
		// 		dataObj := resource.RESTObject.(*web.ApplicationDataPrivacy)
		// 		for _, resource := range resourceData["dynatrace_web_application"] {
		// 			resourceObj := resource.RESTObject.(*web.ApplicationConfig)
		// 			if dataObj.WebApplicationID != nil && *dataObj.WebApplicationID == *resourceObj.ID {
		// 				dataObj.WebApplicationID = opt.NewString("HCL-UNQUOTE-dynatrace_web_application." + escape(resourceObj.Name) + ".id")
		// 			}
		// 		}
		// 	}
		// },
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
		DsReplaceIds: func(resources Resources, dsData DataSourceData) []ReplacedID {
			var ids = []ReplacedID{}
			for _, resource := range resources {
				dataObj := resource.RESTObject.(*applicationdetectionrules.ApplicationDetectionRule)
				dsName := "dynatrace_application"
				for id, appInfo := range dsData[dsName].RESTMap {
					if dataObj.ApplicationIdentifier != "" && dataObj.ApplicationIdentifier == id {
						dataObj.ApplicationIdentifier = "HCL-UNQUOTE-data.dynatrace_application." + escape(appInfo["name"].(string)) + ".id"
						ids = append(ids, ReplacedID{id, dsName})
					}
				}
			}
			return ids
		},
		// ResReplaceIds: func(resName string, resourceData ResourceData) {
		// 	for _, resource := range resourceData[resName] {
		// 		dataObj := resource.RESTObject.(*applicationdetectionrules.ApplicationDetectionRule)
		// 		for _, resource := range resourceData["dynatrace_web_application"] {
		// 			resourceObj := resource.RESTObject.(*web.ApplicationConfig)
		// 			if dataObj.ApplicationIdentifier != "" && dataObj.ApplicationIdentifier == *resourceObj.ID {
		// 				dataObj.ApplicationIdentifier = "HCL-UNQUOTE-dynatrace_web_application." + escape(resourceObj.Name) + ".id"
		// 			}
		// 		}
		// 	}
		// },
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
		DsReplaceIds: func(resources Resources, dsData DataSourceData) []ReplacedID {
			var ids = []ReplacedID{}
			for _, resource := range resources {
				dataObj := resource.RESTObject.(*web.ApplicationErrorRules)
				dsName := "dynatrace_application"
				for id, appInfo := range dsData[dsName].RESTMap {
					if dataObj.WebApplicationID != "" && dataObj.WebApplicationID == id {
						dataObj.WebApplicationID = "HCL-UNQUOTE-data.dynatrace_application." + escape(appInfo["name"].(string)) + ".id"
						ids = append(ids, ReplacedID{id, dsName})
					}
				}
			}
			return ids
		},
		// ResReplaceIds: func(resName string, resourceData ResourceData) {
		// 	for _, resource := range resourceData[resName] {
		// 		dataObj := resource.RESTObject.(*web.ApplicationErrorRules)
		// 		for _, resource := range resourceData["dynatrace_web_application"] {
		// 			resourceObj := resource.RESTObject.(*web.ApplicationConfig)
		// 			if dataObj.WebApplicationID != "" && dataObj.WebApplicationID == *resourceObj.ID {
		// 				dataObj.WebApplicationID = "HCL-UNQUOTE-dynatrace_web_application." + escape(resourceObj.Name) + ".id"
		// 			}
		// 		}
		// 	}
		// },
	},
	"dynatrace_autotag": {
		RESTClient: func(environmentURL, apiToken string) []StandardClient {
			clients := []StandardClient{autotags.NewService(environmentURL+"/api/config/v1", apiToken)}
			return clients
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
		HardcodedIds: []string{"dynatrace_application"},
		DsReplaceIds: func(resources Resources, dsData DataSourceData) []ReplacedID {
			var ids = []ReplacedID{}
			for _, resource := range resources {
				dataObj := resource.RESTObject.(*monitors.BrowserSyntheticMonitorUpdate)
				for idx, assignedApp := range dataObj.ManuallyAssignedApps {
					dsName := "dynatrace_application"
					for id, appInfo := range dsData[dsName].RESTMap {
						if assignedApp == id {
							dataObj.ManuallyAssignedApps[idx] = "HCL-UNQUOTE-" + "data.dynatrace_application." + escape(appInfo["name"].(string)) + ".id"
							ids = append(ids, ReplacedID{id, dsName})
							break
						}
					}
				}
			}
			return ids
		},
		// ResReplaceIds: func(resName string, resourceData ResourceData) {
		// 	for _, resource := range resourceData[resName] {
		// 		dataObj := resource.RESTObject.(*monitors.BrowserSyntheticMonitorUpdate)
		// 		for idx, assignedApp := range dataObj.ManuallyAssignedApps {
		// 			for _, resource := range resourceData["dynatrace_web_application"] {
		// 				if assignedApp == resource.ID {
		// 					dataObj.ManuallyAssignedApps[idx] = "HCL-UNQUOTE-" + "dynatrace_web_application." + escape(resource.Name) + ".id"
		// 				}
		// 			}
		// 		}
		// 	}
		// },
	},
	"dynatrace_calculated_service_metric": {
		RESTClient: func(environmentURL, apiToken string) []StandardClient {
			clients := []StandardClient{service.NewService(environmentURL+"/api/config/v1", apiToken)}
			return clients
		},
		// HardcodedIds: []string{"dynatrace_request_attribute"},
		// DsReplaceIds: func(resources Resources, dsData DataSourceData) []string {
		// 	var ids = []string{}
		// 	for _, resource := range resources {
		// 		dataObj := resource.RESTObject.(*notifications.Notification)
		// 		for id, dsObj := range dsData["dynatrace_request_attribute"].RESTMap {
		// 			if dataObj.ProfileID != "" && dataObj.ProfileID == id {
		// 				dataObj.ProfileID = "HCL-UNQUOTE-data.dynatrace_alerting." + escape(dsObj["name"].(string)) + ".id"
		// 				ids = append(ids, id)
		// 			}
		// 		}
		// 	}
		// 	return ids
		// },
		// ResReplaceIds: func(resName string, resourceData ResourceData) {
		// 	for _, resource := range resourceData[resName] {
		// 		dataObj := resource.RESTObject.(*service.CalculatedServiceMetric)
		// 		if dataObj.MetricDefinition != nil && dataObj.MetricDefinition.RequestAttribute != nil {
		// 			for _, resource := range resourceData["dynatrace_request_attribute"] {
		// 				resourceObj := resource.RESTObject.(*requestattributes.RequestAttribute)
		// 				if *dataObj.MetricDefinition.RequestAttribute == resourceObj.Name {
		// 					dataObj.MetricDefinition.RequestAttribute = opt.NewString("HCL-UNQUOTE-dynatrace_request_attribute." + escape(resourceObj.Name) + ".name")
		// 				}
		// 			}
		// 		}
		// 	}
		// },
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
		DsReplaceIds: func(resources Resources, dsData DataSourceData) []ReplacedID {
			var ids = []ReplacedID{}
			for _, resource := range resources {
				dataObj := resource.RESTObject.(*dashboards.Dashboard)
				if dataObj.Metadata != nil && dataObj.Metadata.Filter != nil && dataObj.Metadata.Filter.ManagementZone != nil {
					dsName := "dynatrace_management_zone"
					for id, appInfo := range dsData[dsName].RESTMap {
						if dataObj.Metadata.Filter.ManagementZone.ID == id {
							dataObj.Metadata.Filter.ManagementZone.ID = "HCL-UNQUOTE-data.dynatrace_management_zone." + escape(appInfo["name"].(string)) + ".id"
							dataObj.Metadata.Filter.ManagementZone.Name = nil
							dataObj.Metadata.Filter.ManagementZone.Description = nil
							ids = append(ids, ReplacedID{id, dsName})
						}
					}
				}
				if dataObj.Tiles != nil {
					for _, tile := range dataObj.Tiles {
						if tile.Filter != nil && tile.Filter.ManagementZone != nil {
							dsName := "dynatrace_management_zone"
							for id, appInfo := range dsData[dsName].RESTMap {
								if tile.Filter.ManagementZone.ID == id {
									tile.Filter.ManagementZone.ID = "HCL-UNQUOTE-data.dynatrace_management_zone." + escape(appInfo["name"].(string)) + ".id"
									tile.Filter.ManagementZone.Name = nil
									tile.Filter.ManagementZone.Description = nil
									ids = append(ids, ReplacedID{id, dsName})
								}
							}
						}
					}
				}
			}
			return ids
		},
		// ResReplaceIds: func(resName string, resourceData ResourceData) {
		// 	for _, resource := range resourceData[resName] {
		// 		dataObj := resource.RESTObject.(*dashboards.Dashboard)
		// 		if dataObj.Metadata != nil && dataObj.Metadata.Filter != nil && dataObj.Metadata.Filter.ManagementZone != nil {
		// 			for _, resource := range resourceData["dynatrace_management_zone"] {
		// 				resourceObj := resource.RESTObject.(*managementzones.ManagementZone)
		// 				if dataObj.Metadata.Filter.ManagementZone.ID == *resourceObj.ID {
		// 					dataObj.Metadata.Filter.ManagementZone.ID = "HCL-UNQUOTE-dynatrace_management_zone." + escape(resourceObj.Name) + ".id"
		// 					dataObj.Metadata.Filter.ManagementZone.Name = nil
		// 					dataObj.Metadata.Filter.ManagementZone.Description = nil
		// 				}
		// 			}
		// 		}
		// 		if dataObj.Tiles != nil {
		// 			for _, tile := range dataObj.Tiles {
		// 				if tile.Filter != nil && tile.Filter.ManagementZone != nil {
		// 					for _, resource := range resourceData["dynatrace_management_zone"] {
		// 						resourceObj := resource.RESTObject.(*managementzones.ManagementZone)
		// 						if tile.Filter.ManagementZone.ID == *resourceObj.ID {
		// 							tile.Filter.ManagementZone.ID = "HCL-UNQUOTE-dynatrace_management_zone." + escape(resourceObj.Name) + ".id"
		// 							tile.Filter.ManagementZone.Name = nil
		// 							tile.Filter.ManagementZone.Description = nil
		// 						}
		// 					}
		// 				}
		// 			}
		// 		}
		// 	}
		// },
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
		HardcodedIds: []string{"dynatrace_alerting"},
		DsReplaceIds: func(resources Resources, dsData DataSourceData) []ReplacedID {
			return replaceIdsNotif(resources, dsData)
		},
		// ResReplaceIds: func(resName string, resourceData ResourceData) {
		// 	for _, resource := range resourceData[resName] {
		// 		dataObj := resource.RESTObject.(*notifications.Notification)
		// 		for _, resource := range resourceData["dynatrace_alerting"] {
		// 			resourceObj := resource.RESTObject.(*alerting.Profile)
		// 			if dataObj.ProfileID != "" && dataObj.ProfileID == resourceObj.ID {
		// 				dataObj.ProfileID = "HCL-UNQUOTE-dynatrace_alerting." + escape(resourceObj.Name) + ".id"
		// 			}
		// 		}
		// 	}
		// },
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
		HardcodedIds: []string{"dynatrace_alerting"},
		DsReplaceIds: func(resources Resources, dsData DataSourceData) []ReplacedID {
			return replaceIdsNotif(resources, dsData)
		},
		// ResReplaceIds: func(resName string, resourceData ResourceData) {
		// 	for _, resource := range resourceData[resName] {
		// 		dataObj := resource.RESTObject.(*notifications.Notification)
		// 		for _, resource := range resourceData["dynatrace_alerting"] {
		// 			resourceObj := resource.RESTObject.(*alerting.Profile)
		// 			if dataObj.ProfileID != "" && dataObj.ProfileID == resourceObj.ID {
		// 				dataObj.ProfileID = "HCL-UNQUOTE-dynatrace_alerting." + escape(resourceObj.Name) + ".id"
		// 			}
		// 		}
		// 	}
		// },
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
		HardcodedIds: []string{"dynatrace_alerting"},
		DsReplaceIds: func(resources Resources, dsData DataSourceData) []ReplacedID {
			return replaceIdsNotif(resources, dsData)
		},
		// ResReplaceIds: func(resName string, resourceData ResourceData) {
		// 	for _, resource := range resourceData[resName] {
		// 		dataObj := resource.RESTObject.(*notifications.Notification)
		// 		for _, resource := range resourceData["dynatrace_alerting"] {
		// 			resourceObj := resource.RESTObject.(*alerting.Profile)
		// 			if dataObj.ProfileID != "" && dataObj.ProfileID == resourceObj.ID {
		// 				dataObj.ProfileID = "HCL-UNQUOTE-dynatrace_alerting." + escape(resourceObj.Name) + ".id"
		// 			}
		// 		}
		// 	}
		// },
	},
	"dynatrace_pager_duty_notification": {
		RESTClient: func(environmentURL, apiToken string) []StandardClient {
			clients := []StandardClient{notifications.NewPagerDutyService(environmentURL+"/api/v2", apiToken)}
			return clients
		},
		CustomName: func(dlConfig DownloadConfig, resourceName string, v interface{}, counter NameCounter) string {
			return counter.Numbering(v.(*notifications.Notification).Name)
		},
		HardcodedIds: []string{"dynatrace_alerting"},
		DsReplaceIds: func(resources Resources, dsData DataSourceData) []ReplacedID {
			return replaceIdsNotif(resources, dsData)
		},
		// ResReplaceIds: func(resName string, resourceData ResourceData) {
		// 	for _, resource := range resourceData[resName] {
		// 		dataObj := resource.RESTObject.(*notifications.Notification)
		// 		for _, resource := range resourceData["dynatrace_alerting"] {
		// 			resourceObj := resource.RESTObject.(*alerting.Profile)
		// 			if dataObj.ProfileID != "" && dataObj.ProfileID == resourceObj.ID {
		// 				dataObj.ProfileID = "HCL-UNQUOTE-dynatrace_alerting." + escape(resourceObj.Name) + ".id"
		// 			}
		// 		}
		// 	}
		// },
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
		HardcodedIds: []string{"dynatrace_alerting"},
		DsReplaceIds: func(resources Resources, dsData DataSourceData) []ReplacedID {
			return replaceIdsNotif(resources, dsData)
		},
		// ResReplaceIds: func(resName string, resourceData ResourceData) {
		// 	for _, resource := range resourceData[resName] {
		// 		dataObj := resource.RESTObject.(*notifications.Notification)
		// 		for _, resource := range resourceData["dynatrace_alerting"] {
		// 			resourceObj := resource.RESTObject.(*alerting.Profile)
		// 			if dataObj.ProfileID != "" && dataObj.ProfileID == resourceObj.ID {
		// 				dataObj.ProfileID = "HCL-UNQUOTE-dynatrace_alerting." + escape(resourceObj.Name) + ".id"
		// 			}
		// 		}
		// 	}
		// },
	},
	"dynatrace_slack_notification": {
		RESTClient: func(environmentURL, apiToken string) []StandardClient {
			clients := []StandardClient{notifications.NewSlackService(environmentURL+"/api/v2", apiToken)}
			return clients
		},
		CustomName: func(dlConfig DownloadConfig, resourceName string, v interface{}, counter NameCounter) string {
			return counter.Numbering(v.(*notifications.Notification).Name)
		},
		HardcodedIds: []string{"dynatrace_alerting"},
		DsReplaceIds: func(resources Resources, dsData DataSourceData) []ReplacedID {
			return replaceIdsNotif(resources, dsData)
		},
		// ResReplaceIds: func(resName string, resourceData ResourceData) {
		// 	for _, resource := range resourceData[resName] {
		// 		dataObj := resource.RESTObject.(*notifications.Notification)
		// 		for _, resource := range resourceData["dynatrace_alerting"] {
		// 			resourceObj := resource.RESTObject.(*alerting.Profile)
		// 			if dataObj.ProfileID != "" && dataObj.ProfileID == resourceObj.ID {
		// 				dataObj.ProfileID = "HCL-UNQUOTE-dynatrace_alerting." + escape(resourceObj.Name) + ".id"
		// 			}
		// 		}
		// 		dataObj.Slack.URL = "######"
		// 	}
		// },
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
		HardcodedIds: []string{"dynatrace_alerting"},
		DsReplaceIds: func(resources Resources, dsData DataSourceData) []ReplacedID {
			return replaceIdsNotif(resources, dsData)
		},
		// ResReplaceIds: func(resName string, resourceData ResourceData) {
		// 	for _, resource := range resourceData[resName] {
		// 		dataObj := resource.RESTObject.(*notifications.Notification)
		// 		for _, resource := range resourceData["dynatrace_alerting"] {
		// 			resourceObj := resource.RESTObject.(*alerting.Profile)
		// 			if dataObj.ProfileID != "" && dataObj.ProfileID == resourceObj.ID {
		// 				dataObj.ProfileID = "HCL-UNQUOTE-dynatrace_alerting." + escape(resourceObj.Name) + ".id"
		// 			}
		// 		}
		// 	}
		// },
	},
	"dynatrace_victor_ops_notification": {
		RESTClient: func(environmentURL, apiToken string) []StandardClient {
			clients := []StandardClient{notifications.NewVictorOpsService(environmentURL+"/api/v2", apiToken)}
			return clients
		},
		CustomName: func(dlConfig DownloadConfig, resourceName string, v interface{}, counter NameCounter) string {
			return counter.Numbering(v.(*notifications.Notification).Name)
		},
		HardcodedIds: []string{"dynatrace_alerting"},
		DsReplaceIds: func(resources Resources, dsData DataSourceData) []ReplacedID {
			return replaceIdsNotif(resources, dsData)
		},
		// ResReplaceIds: func(resName string, resourceData ResourceData) {
		// 	for _, resource := range resourceData[resName] {
		// 		dataObj := resource.RESTObject.(*notifications.Notification)
		// 		for _, resource := range resourceData["dynatrace_alerting"] {
		// 			resourceObj := resource.RESTObject.(*alerting.Profile)
		// 			if dataObj.ProfileID != "" && dataObj.ProfileID == resourceObj.ID {
		// 				dataObj.ProfileID = "HCL-UNQUOTE-dynatrace_alerting." + escape(resourceObj.Name) + ".id"
		// 			}
		// 		}
		// 	}
		// },
	},
	"dynatrace_web_application": {
		RESTClient: func(environmentURL, apiToken string) []StandardClient {
			clients := []StandardClient{web.NewService(environmentURL+"/api/config/v1", apiToken)}
			return clients
		},
	},
	"dynatrace_webhook_notification": {
		RESTClient: func(environmentURL, apiToken string) []StandardClient {
			clients := []StandardClient{notifications.NewWebHookService(environmentURL+"/api/v2", apiToken)}
			return clients
		},
		CustomName: func(dlConfig DownloadConfig, resourceName string, v interface{}, counter NameCounter) string {
			return counter.Numbering(v.(*notifications.Notification).Name)
		},
		HardcodedIds: []string{"dynatrace_alerting"},
		DsReplaceIds: func(resources Resources, dsData DataSourceData) []ReplacedID {
			return replaceIdsNotif(resources, dsData)
		},
		// ResReplaceIds: func(resName string, resourceData ResourceData) {
		// 	for _, resource := range resourceData[resName] {
		// 		dataObj := resource.RESTObject.(*notifications.Notification)
		// 		for _, resource := range resourceData["dynatrace_alerting"] {
		// 			resourceObj := resource.RESTObject.(*alerting.Profile)
		// 			if dataObj.ProfileID != "" && dataObj.ProfileID == resourceObj.ID {
		// 				dataObj.ProfileID = "HCL-UNQUOTE-dynatrace_alerting." + escape(resourceObj.Name) + ".id"
		// 			}
		// 		}
		// 	}
		// },
	},
	"dynatrace_xmatters_notification": {
		RESTClient: func(environmentURL, apiToken string) []StandardClient {
			clients := []StandardClient{notifications.NewXMattersService(environmentURL+"/api/v2", apiToken)}
			return clients
		},
		CustomName: func(dlConfig DownloadConfig, resourceName string, v interface{}, counter NameCounter) string {
			return counter.Numbering(v.(*notifications.Notification).Name)
		},
		HardcodedIds: []string{"dynatrace_alerting"},
		DsReplaceIds: func(resources Resources, dsData DataSourceData) []ReplacedID {
			return replaceIdsNotif(resources, dsData)
		},
		// ResReplaceIds: func(resName string, resourceData ResourceData) {
		// 	for _, resource := range resourceData[resName] {
		// 		dataObj := resource.RESTObject.(*notifications.Notification)
		// 		for _, resource := range resourceData["dynatrace_alerting"] {
		// 			resourceObj := resource.RESTObject.(*alerting.Profile)
		// 			if dataObj.ProfileID != "" && dataObj.ProfileID == resourceObj.ID {
		// 				dataObj.ProfileID = "HCL-UNQUOTE-dynatrace_alerting." + escape(resourceObj.Name) + ".id"
		// 			}
		// 		}
		// 	}
		// },
	},
}

type ResourceStruct struct {
	RESTClient   func(environmentURL, apiToken string) []StandardClient
	NoListClient func(environmentURL, apiToken string) NoListClient
	CustomName   NameFunc
	HardcodedIds []string
	DsReplaceIds DataSourceReplaceFunc
	// ResReplaceIds ResourceReplaceFunc
}

type NameFunc func(DownloadConfig, string, interface{}, NameCounter) string

type DataSourceReplaceFunc func(Resources, DataSourceData) []ReplacedID

// type ResourceReplaceFunc func(string, ResourceData)

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

func replaceIdsNotif(resources Resources, dsData DataSourceData) []ReplacedID {
	var ids = []ReplacedID{}
	for _, resource := range resources {
		dataObj := resource.RESTObject.(*notifications.Notification)
		dsName := "dynatrace_alerting"
		for id, dsObj := range dsData[dsName].RESTMap {
			if dataObj.ProfileID != "" && dataObj.ProfileID == id {
				dataObj.ProfileID = "HCL-UNQUOTE-data.dynatrace_alerting." + escape(dsObj["name"].(string)) + ".id"
				ids = append(ids, ReplacedID{id, dsName})
			}
		}
	}
	return ids
}
