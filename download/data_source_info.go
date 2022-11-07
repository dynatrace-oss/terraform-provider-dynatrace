package download

import (
	api "github.com/dtcookie/dynatrace/api/config"
	"github.com/dtcookie/dynatrace/api/config/applications/web"
	"github.com/dtcookie/dynatrace/api/config/managementzones"
	"github.com/dtcookie/dynatrace/api/config/topology/service"
)

var DataSourceInfoMap = map[string]DataSourceStruct{
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
}

type DataSourceStruct struct {
	RESTClient  func(string, string) DataSourceClient
	MarshallHCL func(interface{}) map[string]map[string]interface{}
}
