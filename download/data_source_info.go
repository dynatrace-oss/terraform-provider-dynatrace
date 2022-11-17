package download

import (
	api "github.com/dtcookie/dynatrace/api/config"
	"github.com/dtcookie/dynatrace/api/config/applications/web"
	"github.com/dtcookie/dynatrace/api/config/credentials/vault"
	"github.com/dtcookie/dynatrace/api/config/managementzones"
	"github.com/dtcookie/dynatrace/api/config/requestattributes"
	"github.com/dtcookie/dynatrace/api/config/topology/service"
	"github.com/dtcookie/dynatrace/api/config/v2/alerting"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/resources/synthetic/locations"
)

var DataSourceInfoMap = map[string]DataSourceStruct{
	"dynatrace_alerting_profile": {
		RESTClient: func(environmentURL string, apiToken string) DataSourceClient {
			return alerting.NewService(environmentURL+"/api/v2", apiToken)
		},
		MarshallHCL: func(restObject interface{}) map[string]map[string]interface{} {
			stubs := restObject.([]*alerting.Profile)
			var restMap = map[string]map[string]interface{}{}
			for _, stub := range stubs {
				restMap[stub.ID] = map[string]interface{}{}
				restMap[stub.ID]["name"] = stub.Name
			}
			return restMap
		},
	},
	"dynatrace_application": {
		RESTClient: func(environmentURL string, apiToken string) DataSourceClient {
			return web.NewService(environmentURL+"/api/config/v1", apiToken)
		},
		MarshallHCL: func(restObject interface{}) map[string]map[string]interface{} {
			stubs := restObject.(*api.StubList)
			var restMap = map[string]map[string]interface{}{}
			for _, stub := range stubs.Values {
				restMap[stub.ID] = map[string]interface{}{}
				restMap[stub.ID]["name"] = stub.Name
			}
			return restMap
		},
	},
	"dynatrace_request_attribute": {
		RESTClient: func(environmentURL string, apiToken string) DataSourceClient {
			return requestattributes.NewService(environmentURL+"/api/config/v1", apiToken)
		},
		MarshallHCL: func(restObject interface{}) map[string]map[string]interface{} {
			stubs := restObject.(*api.StubList)
			var restMap = map[string]map[string]interface{}{}
			for _, stub := range stubs.Values {
				restMap[stub.ID] = map[string]interface{}{}
				restMap[stub.ID]["name"] = stub.Name
			}
			return restMap
		},
	},
	"dynatrace_management_zone": {
		RESTClient: func(environmentURL string, apiToken string) DataSourceClient {
			return managementzones.NewService(environmentURL+"/api/config/v1", apiToken)
		},
		MarshallHCL: func(restObject interface{}) map[string]map[string]interface{} {
			stubs := restObject.([]*api.EntityShortRepresentation)
			var restMap = map[string]map[string]interface{}{}
			for _, stub := range stubs {
				restMap[stub.ID] = map[string]interface{}{}
				restMap[stub.ID]["name"] = stub.Name
			}
			return restMap
		},
	},
	"dynatrace_service": {
		RESTClient: func(environmentURL, apiToken string) DataSourceClient {
			return service.NewService(environmentURL+"/api/v1", apiToken)
		},
	},
	"dynatrace_synthetic_location": {
		RESTClient: func(environmentURL string, apiToken string) DataSourceClient {
			return locations.NewService(environmentURL+"/api/v1", apiToken)
		},
		MarshallHCL: func(restObject interface{}) map[string]map[string]interface{} {
			locations := restObject.(*locations.SyntheticLocations)
			var restMap = map[string]map[string]interface{}{}
			for _, location := range locations.Locations {
				restMap[location.ID] = map[string]interface{}{}
				restMap[location.ID]["name"] = location.Name
				restMap[location.ID]["type"] = location.Type
				if location.CloudPlatform != nil {
					restMap[location.ID]["cloud_platform"] = *location.CloudPlatform
				}
			}
			return restMap
		},
		UniqueName: func(values map[string]interface{}) string {
			if values["cloud_platform"] != nil {
				return values["name"].(string) + "_" + string(values["cloud_platform"].(locations.CloudPlatform))
			}
			return values["name"].(string)
		},
	},
	"dynatrace_credentials": {
		RESTClient: func(environmentURL string, apiToken string) DataSourceClient {
			return vault.NewService(environmentURL+"/api/v2", apiToken)
		},
		MarshallHCL: func(restObject interface{}) map[string]map[string]interface{} {
			credentialList := restObject.(*vault.CredentialsList)
			var restMap = map[string]map[string]interface{}{}
			for _, credential := range credentialList.Credentials {
				restMap[*credential.ID] = map[string]interface{}{}
				restMap[*credential.ID]["type"] = string(credential.Type)
				restMap[*credential.ID]["name"] = credential.Name
				restMap[*credential.ID]["scope"] = string(credential.Scope)
			}
			return restMap
		},
		UniqueName: func(values map[string]interface{}) string {
			return values["name"].(string)
		},
	},
}

type DataSourceStruct struct {
	RESTClient  func(string, string) DataSourceClient
	MarshallHCL func(interface{}) map[string]map[string]interface{}
	UniqueName  func(map[string]interface{}) string
}

func UniqueDSName(dsName string, values map[string]interface{}) string {
	if DataSourceInfoMap[dsName].UniqueName != nil {
		return DataSourceInfoMap[dsName].UniqueName(values)
	}
	return values["name"].(string)
}
