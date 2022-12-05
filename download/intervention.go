package download

import (
	"os"
	"strings"

	"github.com/dtcookie/dynatrace/api/config/dashboards"
	"github.com/dtcookie/dynatrace/api/config/metrics/calculated/service"
	privlocations "github.com/dtcookie/dynatrace/api/config/synthetic/locations"
)

var InterventionInfoMap = map[string]InterventionStruct{
	"dynatrace_dashboard": {
		Move: func(resName string, resourceData ResourceData) {
			for _, resource := range resourceData[resName] {
				dashboard := resource.RESTObject.(*dashboards.Dashboard)
				dbId := dashboard.ID
				dashboard.ID = nil
				environmentURL := os.Getenv("DYNATRACE_ENV_URL")
				environmentURL = strings.TrimSuffix(strings.TrimSuffix(environmentURL, " "), "/")
				apiToken := os.Getenv("DYNATRACE_API_TOKEN")
				client := dashboards.NewService(environmentURL+"/api/config/v1", apiToken)
				errors := client.Validate(dashboard)
				dashboard.ID = dbId
				if len(errors) > 0 && !strings.Contains(errors[0], "Token is missing required scope. Use one of: WriteConfig (Write configuration)") {
					resource.ReqInter.Type = InterventionTypes.Flawed
					errors[0] = "ATTENTION " + strings.ReplaceAll(errors[0], "\n", "")
					resource.ReqInter.Message = errors
					continue
				}
				if dashboard.Metadata.Owner != nil && *dashboard.Metadata.Owner == "Dynatrace" {
					resource.ReqInter.Type = InterventionTypes.ReqAttn
					resource.ReqInter.Message = []string{"ATTENTION " + "Dashboards owned by Dynatrace are automatically excluded to prevent duplicates of OOTB dashboards. Please return to the dashboards folder if this is a custom dashboard."}
				}
			}
		},
	},
	"dynatrace_calculated_service_metric": {
		Move: func(resName string, resourceData ResourceData) {
			reqConditions := []string{"SERVICE_DISPLAY_NAME", "SERVICE_PUBLIC_DOMAIN_NAME", "SERVICE_WEB_APPLICATION_ID", "SERVICE_WEB_CONTEXT_ROOT", "SERVICE_WEB_SERVER_NAME", "SERVICE_WEB_SERVICE_NAME", "SERVICE_WEB_SERVICE_NAMESPACE", "REMOTE_SERVICE_NAME", "REMOTE_ENDPOINT", "AZURE_FUNCTIONS_SITE_NAME", "AZURE_FUNCTIONS_FUNCTION_NAME", "CTG_GATEWAY_URL", "CTG_SERVER_NAME", "ACTOR_SYSTEM", "ESB_APPLICATION_NAME", "SERVICE_TAG", "SERVICE_TYPE", "PROCESS_GROUP_TAG", "PROCESS_GROUP_NAME"}
			for _, resource := range resourceData[resName] {
				dataObj := resource.RESTObject.(*service.CalculatedServiceMetric)
				if len(dataObj.ManagementZones) == 0 && dataObj.Conditions != nil {
					var found bool
					for _, condition := range dataObj.Conditions {
						for _, reqCondition := range reqConditions {
							if string(condition.Attribute) == reqCondition {
								found = true
							}
						}
					}
					if !found {
						resource.ReqInter.Type = InterventionTypes.Flawed
						resource.ReqInter.Message = []string{"ATTENTION " + "The metric needs to either get limited by specifying a Management Zone or by specifying one or more conditions related to SERVICE_DISPLAY_NAME, SERVICE_PUBLIC_DOMAIN_NAME, SERVICE_WEB_APPLICATION_ID, SERVICE_WEB_CONTEXT_ROOT, SERVICE_WEB_SERVER_NAME, SERVICE_WEB_SERVICE_NAME, SERVICE_WEB_SERVICE_NAMESPACE, REMOTE_SERVICE_NAME, REMOTE_ENDPOINT, AZURE_FUNCTIONS_SITE_NAME, AZURE_FUNCTIONS_FUNCTION_NAME, CTG_GATEWAY_URL, CTG_SERVER_NAME, ACTOR_SYSTEM, ESB_APPLICATION_NAME, SERVICE_TAG, SERVICE_TYPE, PROCESS_GROUP_TAG or PROCESS_GROUP_NAME"}
					}
				}

			}
		},
	},
	"dynatrace_credentials": {
		Move: func(resName string, resourceData ResourceData) {
			for _, resource := range resourceData[resName] {
				resource.ReqInter.Type = InterventionTypes.ReqAttn
				resource.ReqInter.Message = []string{"ATTENTION " + "REST API didn't provide credentials"}
			}
		},
	},
	"dynatrace_aws_credentials": {
		Move: func(resName string, resourceData ResourceData) {
			for _, resource := range resourceData[resName] {
				resource.ReqInter.Type = InterventionTypes.ReqAttn
				resource.ReqInter.Message = []string{"ATTENTION " + "REST API didn't provide credentials"}
			}
		},
	},
	"dynatrace_cloudfoundry_credentials": {
		Move: func(resName string, resourceData ResourceData) {
			for _, resource := range resourceData[resName] {
				resource.ReqInter.Type = InterventionTypes.ReqAttn
				resource.ReqInter.Message = []string{"ATTENTION " + "REST API didn't provide credentials"}
			}
		},
	},
	"dynatrace_azure_credentials": {
		Move: func(resName string, resourceData ResourceData) {
			for _, resource := range resourceData[resName] {
				resource.ReqInter.Type = InterventionTypes.ReqAttn
				resource.ReqInter.Message = []string{"ATTENTION " + "REST API didn't provide credentials"}
			}
		},
	},
	"dynatrace_k8s_credentials": {
		Move: func(resName string, resourceData ResourceData) {
			for _, resource := range resourceData[resName] {
				resource.ReqInter.Type = InterventionTypes.ReqAttn
				resource.ReqInter.Message = []string{"ATTENTION " + "REST API didn't provide credentials"}
			}
		},
	},
	"dynatrace_synthetic_location": {
		StripIDs: func(resName string, resourceData ResourceData) {
			for _, resource := range resourceData[resName] {
				dataObj := resource.RESTObject.(*privlocations.PrivateSyntheticLocation)
				if len(dataObj.Nodes) > 0 {
					dataObj.Nodes = []string{}
				}
			}
		},
	},
}

type InterventionStruct struct {
	Move     func(string, ResourceData)
	StripIDs func(string, ResourceData)
}

func (me ResourceData) RequiresIntervention(dlConfig DownloadConfig) error {
	for resName := range me {
		if _, exists := InterventionInfoMap[resName]; exists {
			if InterventionInfoMap[resName].Move != nil {
				InterventionInfoMap[resName].Move(resName, me)
			}
			if dlConfig.Migrate && InterventionInfoMap[resName].StripIDs != nil {
				InterventionInfoMap[resName].StripIDs(resName, me)
			}
		}
	}
	return nil
}
