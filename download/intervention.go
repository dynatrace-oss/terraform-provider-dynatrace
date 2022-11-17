package download

import (
	"github.com/dtcookie/dynatrace/api/config/metrics/calculated/service"
)

var InterventionInfoMap = map[string]InterventionStruct{
	"dynatrace_calculated_service_metric": {
		Move: func(resName string, resourceData ResourceData) {
			reqConditions := []string{"SERVICE_DISPLAY_NAME", "SERVICE_PUBLIC_DOMAIN_NAME", "SERVICE_WEB_APPLICATION_ID", "SERVICE_WEB_CONTEXT_ROOT", "SERVICE_WEB_SERVER_NAME", "SERVICE_WEB_SERVICE_NAME", "SERVICE_WEB_SERVICE_NAMESPACE", "REMOTE_SERVICE_NAME", "REMOTE_ENDPOINT", "AZURE_FUNCTIONS_SITE_NAME", "AZURE_FUNCTIONS_FUNCTION_NAME", "CTG_GATEWAY_URL", "CTG_SERVER_NAME", "ACTOR_SYSTEM", "ESB_APPLICATION_NAME", "SERVICE_TAG", "SERVICE_TYPE", "PROCESS_GROUP_TAG", "PROCESS_GROUP_NAME"}
			for idx, resource := range resourceData[resName] {
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
						resourceData[resName][idx].ReqInter = true
					}
				}

			}
		},
	},
}

type InterventionStruct struct {
	Move func(string, ResourceData)
}

func (me ResourceData) RequiresIntervention() error {
	for resName := range me {
		if _, exists := InterventionInfoMap[resName]; exists {
			InterventionInfoMap[resName].Move(resName, me)
		}
	}
	return nil
}
