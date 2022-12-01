package download

import (
	api "github.com/dtcookie/dynatrace/api/config"
	"github.com/dtcookie/dynatrace/api/config/applications/web"
	"github.com/dtcookie/dynatrace/api/config/credentials/vault"
	"github.com/dtcookie/dynatrace/api/config/managementzones"
	"github.com/dtcookie/dynatrace/api/config/requestattributes"
	"github.com/dtcookie/dynatrace/api/config/topology/service"
	"github.com/dtcookie/dynatrace/api/config/v2/alerting"
	mgmzapi20 "github.com/dtcookie/dynatrace/api/config/v2/managementzones"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/resources/synthetic/locations"
)

var DataSourceInfoMap = map[string]DataSourceStruct{
	"dynatrace_alerting_profile": {
		RESTClient: func(environmentURL string, apiToken string) DataSourceClient {
			return alerting.NewService(environmentURL+"/api/v2", apiToken)
		},
		MarshallHCL: func(restObject interface{}, dlConfig DownloadConfig) map[string]*DataSourceDetails {
			stubs := restObject.([]*alerting.Profile)
			var restMap = map[string]*DataSourceDetails{}
			for _, stub := range stubs {
				restMap[stub.ID] = &DataSourceDetails{Values: map[string]interface{}{}}
				restMap[stub.ID].Values["name"] = stub.Name
			}
			return restMap
		},
	},
	"dynatrace_application": {
		RESTClient: func(environmentURL string, apiToken string) DataSourceClient {
			return web.NewService(environmentURL+"/api/config/v1", apiToken)
		},
		MarshallHCL: func(restObject interface{}, dlConfig DownloadConfig) map[string]*DataSourceDetails {
			stubs := restObject.(*api.StubList)
			var restMap = map[string]*DataSourceDetails{}
			for _, stub := range stubs.Values {
				restMap[stub.ID] = &DataSourceDetails{Values: map[string]interface{}{}}
				restMap[stub.ID].Values["name"] = stub.Name
			}
			return restMap
		},
	},
	"dynatrace_request_attribute": {
		RESTClient: func(environmentURL string, apiToken string) DataSourceClient {
			return requestattributes.NewService(environmentURL+"/api/config/v1", apiToken)
		},
		MarshallHCL: func(restObject interface{}, dlConfig DownloadConfig) map[string]*DataSourceDetails {
			stubs := restObject.(*api.StubList)
			var restMap = map[string]*DataSourceDetails{}
			for _, stub := range stubs.Values {
				restMap[stub.ID] = &DataSourceDetails{Values: map[string]interface{}{}}
				restMap[stub.ID].Values["name"] = stub.Name
			}
			return restMap
		},
	},
	"dynatrace_management_zone": {
		RESTClient: func(environmentURL string, apiToken string) DataSourceClient {
			return managementzones.NewService(environmentURL+"/api/config/v1", apiToken)
		},
		MarshallHCL: func(restObject interface{}, dlConfig DownloadConfig) map[string]*DataSourceDetails {
			stubs := restObject.([]*api.EntityShortRepresentation)
			var restMap = map[string]*DataSourceDetails{}
			for _, stub := range stubs {
				restMap[stub.ID] = &DataSourceDetails{Values: map[string]interface{}{}}
				restMap[stub.ID].Values["name"] = stub.Name
			}
			v20Client := mgmzapi20.NewService(dlConfig.EnvironmentURL+"/api/v2", dlConfig.APIToken)
			v2Stubs, err := v20Client.List()
			if err != nil {
				return nil
			}
			for _, stub := range v2Stubs {
				for _, restMapEntry := range restMap {
					if stub.Name == restMapEntry.Values["name"] {
						restMapEntry.Values["settings_20_id"] = stub.ID
					}
				}
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
		MarshallHCL: func(restObject interface{}, dlConfig DownloadConfig) map[string]*DataSourceDetails {
			locations := restObject.(*locations.SyntheticLocations)
			var restMap = map[string]*DataSourceDetails{}
			for _, location := range locations.Locations {
				restMap[location.ID] = &DataSourceDetails{Values: map[string]interface{}{}}
				restMap[location.ID].Values["name"] = location.Name
				restMap[location.ID].Values["type"] = location.Type
				if location.CloudPlatform != nil {
					restMap[location.ID].Values["cloud_platform"] = *location.CloudPlatform
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
		MarshallHCL: func(restObject interface{}, dlConfig DownloadConfig) map[string]*DataSourceDetails {
			credentialList := restObject.(*vault.CredentialsList)
			var restMap = map[string]*DataSourceDetails{}
			for _, credential := range credentialList.Credentials {
				restMap[*credential.ID] = &DataSourceDetails{Values: map[string]interface{}{}}
				restMap[*credential.ID].Values["type"] = string(credential.Type)
				restMap[*credential.ID].Values["name"] = credential.Name
				restMap[*credential.ID].Values["scope"] = string(credential.Scope)
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
	MarshallHCL func(interface{}, DownloadConfig) map[string]*DataSourceDetails
	UniqueName  func(map[string]interface{}) string
}

func UniqueDSName(dsName string, values map[string]interface{}) string {
	if DataSourceInfoMap[dsName].UniqueName != nil {
		return DataSourceInfoMap[dsName].UniqueName(values)
	}
	return values["name"].(string)
}
