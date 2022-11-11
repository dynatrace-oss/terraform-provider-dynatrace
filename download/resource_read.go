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

package download

import (
	"context"
	"fmt"
	"strings"

	"github.com/dtcookie/dynatrace/api/config/dashboards"
	"github.com/dtcookie/dynatrace/api/config/dashboards/sharing"
	servicetopology "github.com/dtcookie/dynatrace/api/config/topology/service"
	"github.com/dtcookie/dynatrace/api/config/v2/keyrequests"
	"github.com/dtcookie/hcl"
)

func (me ResourceData) ProcessRead(dlConfig DownloadConfig) error {
	for resName, resStruct := range ResourceInfoMap {
		if !dlConfig.MatchResource(resName) && len(dlConfig.ResourceNames) != 0 {
			continue
		}
		fmt.Println("Processing read: ", resName)
		if ResourceInfoMap[resName].NoListClient != nil {
			client := resStruct.NoListClient(
				dlConfig.EnvironmentURL,
				dlConfig.APIToken,
			)
			if err := me.readNoList(resName, client); err != nil {
				return err
			}
		} else if resName == "dynatrace_dashboard" {
			clients := resStruct.RESTClient(
				dlConfig.EnvironmentURL,
				dlConfig.APIToken,
			)
			if err := me.readDashboards(dlConfig, resName, clients[0]); err != nil {
				return err
			}
		} else if resName == "dynatrace_key_requests" {
			if err := me.readKeyRequests(dlConfig, resName); err != nil {
				return err
			}
		} else {
			clients := resStruct.RESTClient(
				dlConfig.EnvironmentURL,
				dlConfig.APIToken,
			)
			for _, client := range clients {
				if err := me.read(dlConfig, resName, client); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (me ResourceData) read(dlConfig DownloadConfig, resourceName string, client StandardClient) error {
	resources := Resources{}
	ids, err := client.LIST()
	if err != nil {
		return err
	}
	nameCounter := NameCounter{}
	for _, id := range ids {
		if !dlConfig.MatchID(resourceName, id) && dlConfig.ResourceNames[resourceName] != nil {
			continue
		}
		config, err := client.GET(id)
		if err != nil {
			return err
		}
		if marshaller, ok := config.(hcl.Marshaler); ok {
			resource := Resource{
				ID:         id,
				Name:       ResourceInfoMap[resourceName].Name(dlConfig, resourceName, config, &nameCounter),
				RESTObject: marshaller,
			}
			resources = append(resources, resource)
		}
	}
	if _, found := me[resourceName]; found {
		me[resourceName] = append(me[resourceName], resources...)
	} else {
		me[resourceName] = resources
	}

	return nil
}

func (me ResourceData) readNoList(resourceName string, client NoListClient) error {
	resources := Resources{}
	config, err := client.GET()
	if err != nil {
		return err
	}
	if marshaller, ok := config.(hcl.Marshaler); ok {
		resource := Resource{
			RESTObject: marshaller,
			Name:       ResourceInfoMap[resourceName].Name(DownloadConfig{}, resourceName, config, nil),
		}
		resources = append(resources, resource)
	}
	me[resourceName] = resources

	return nil
}

func (me ResourceData) readDashboards(dlConfig DownloadConfig, resourceName string, client StandardClient) error {
	resources := Resources{}
	ids, err := client.LIST()

	dashboardSharing := Resources{}
	shareRestClient := sharing.NewService(dlConfig.EnvironmentURL+"/api/config/v1", dlConfig.APIToken)
	if err != nil {
		return err
	}

	nameCounter := NameCounter{}
	for _, id := range ids {
		if !dlConfig.MatchID(resourceName, id) && dlConfig.ResourceNames[resourceName] != nil {
			continue
		}
		config, err := client.GET(id)
		if err != nil {
			return err
		}
		var name string
		if marshaller, ok := config.(hcl.Marshaler); ok {
			name = ResourceInfoMap[resourceName].Name(dlConfig, resourceName, config, &nameCounter)
			resource := Resource{
				RESTObject: marshaller,
				Name:       name,
			}
			resources = append(resources, resource)
		}

		shareSettings, err := shareRestClient.GET(context.Background(), id)
		if err != nil {
			if strings.Contains(err.Error(), "Editing or deleting a non user specific dashboard preset is not allowed") {
				continue
			} else {
				return err
			}
		}
		if marshaller, ok := shareSettings.(hcl.Marshaler); ok {
			resource := Resource{
				RESTObject: marshaller,
				Name:       name,
			}
			dataObj := resource.RESTObject.(*sharing.DashboardSharing)
			dataObj.DashboardID = "HCL-UNQUOTE-dynatrace_dashboard." + Escape(config.(*dashboards.Dashboard).Metadata.Name) + ".id"
			dashboardSharing = append(dashboardSharing, resource)
		}
	}
	me[resourceName] = resources

	if len(dashboardSharing) > 0 {
		me["dynatrace_dashboard_sharing"] = dashboardSharing
	}

	return nil
}

func (me ResourceData) readKeyRequests(dlConfig DownloadConfig, resourceName string) error {
	resources := Resources{}
	topRestClient := DataSourceInfoMap["dynatrace_service"].RESTClient(
		dlConfig.EnvironmentURL,
		dlConfig.APIToken,
	)
	services, err := topRestClient.ListInterface()
	if err != nil {
		return err
	}
	restClient := keyrequests.NewService(dlConfig.EnvironmentURL+"/api/v2", dlConfig.APIToken)

	for _, service := range services.(servicetopology.Services) {
		keyRequestID, keyRequest, err := restClient.LIST(service.EntityId)

		kr := keyRequest.(*keyrequests.KeyRequest)
		if err != nil {
			return err
		}
		if kr == nil {
			continue
		}
		if !dlConfig.MatchID(resourceName, keyRequestID) && dlConfig.ResourceNames[resourceName] != nil {
			continue
		}
		if marshaller, ok := keyRequest.(hcl.Marshaler); ok {
			resource := Resource{
				RESTObject: marshaller,
				Name:       ResourceInfoMap[resourceName].Name(dlConfig, resourceName, service, nil),
			}
			resources = append(resources, resource)
		}
	}
	me[resourceName] = resources

	return nil
}
