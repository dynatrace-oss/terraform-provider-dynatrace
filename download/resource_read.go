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

	"github.com/dtcookie/dynatrace/api/config/dashboards/sharing"
	servicetopology "github.com/dtcookie/dynatrace/api/config/topology/service"
	"github.com/dtcookie/dynatrace/api/config/v2/keyrequests"
	"github.com/dtcookie/hcl"
	"github.com/google/uuid"
)

func (me ResourceData) ProcessRead(dlConfig DownloadConfig) error {
	for resName, resStruct := range ResourceInfoMap {
		if len(dlConfig.ResourceNames) != 0 {
			if (!dlConfig.Exclude && !dlConfig.MatchResource(resName)) || (dlConfig.Exclude && dlConfig.MatchResource(resName)) {
				continue
			}
		}
		if resName == "dynatrace_dashboard" && !dlConfig.MatchResource(resName) {
			continue
		}
		if resName == "dynatrace_json_dashboard" && !dlConfig.MatchResource(resName) {
			continue
		}
		if resName == "dynatrace_iam_user" && !dlConfig.MatchResource(resName) {
			continue
		}
		if resName == "dynatrace_iam_group" && !dlConfig.MatchResource(resName) {
			continue
		}
		if resName == "dynatrace_management_zone" && !dlConfig.MatchResource(resName) {
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
			if err := me.readDashboards(dlConfig, resName, clients[0], nil); err != nil {
				return err
			}
		} else if resName == "dynatrace_key_requests" {
			if err := me.readKeyRequests(dlConfig, resName, nil); err != nil {
				return err
			}
		} else {
			var clients []StandardClient
			if resStruct.IAMClient != nil {
				if len(dlConfig.IAMClientID) == 0 || len(dlConfig.IAMAccountID) == 0 || len(dlConfig.IAMClientSecret) == 0 {
					continue
				}
				clients = resStruct.IAMClient(
					dlConfig.IAMClientID,
					dlConfig.IAMAccountID,
					dlConfig.IAMClientSecret,
				)
			} else {
				clients = resStruct.RESTClient(
					dlConfig.EnvironmentURL,
					dlConfig.APIToken,
				)
			}
			for _, client := range clients {
				if err := me.read(dlConfig, resName, client, nil); err != nil {
					return err
				}
			}

		}
	}

	return nil
}

func (me ResourceData) ProcessRepIdRead(dlConfig DownloadConfig, replacedIds ReplacedIDs) error {
	for _, replacedId := range replacedIds {
		for resName, repId := range replacedId {
			if !containsProcessedRepId(repId) {
				if !dlConfig.Exclude {
					if _, exists := dlConfig.ResourceNames[resName]; exists {
						for _, repIdStruct := range repId {
							repIdStruct.Processed = true
						}
						continue
					}
				} else {
					if _, exists := dlConfig.ResourceNames[resName]; !exists {
						for _, repIdStruct := range repId {
							repIdStruct.Processed = true
						}
						continue
					}
				}
				fmt.Println("Processing read: ", resName)
				if ResourceInfoMap[resName].NoListClient != nil {
					client := ResourceInfoMap[resName].NoListClient(
						dlConfig.EnvironmentURL,
						dlConfig.APIToken,
					)
					if err := me.readNoList(resName, client); err != nil {
						return err
					}
				} else if resName == "dynatrace_dashboard" {
					clients := ResourceInfoMap[resName].RESTClient(
						dlConfig.EnvironmentURL,
						dlConfig.APIToken,
					)
					if err := me.readDashboards(dlConfig, resName, clients[0], repId); err != nil {
						return err
					}
				} else if resName == "dynatrace_key_requests" {
					if err := me.readKeyRequests(dlConfig, resName, repId); err != nil {
						return err
					}
				} else {
					var clients []StandardClient
					if ResourceInfoMap[resName].IAMClient != nil {
						if len(dlConfig.IAMClientID) == 0 || len(dlConfig.IAMAccountID) == 0 || len(dlConfig.IAMClientSecret) == 0 {
							continue
						}
						clients = ResourceInfoMap[resName].IAMClient(
							dlConfig.IAMClientID,
							dlConfig.IAMAccountID,
							dlConfig.IAMClientSecret,
						)
					} else {
						clients = ResourceInfoMap[resName].RESTClient(
							dlConfig.EnvironmentURL,
							dlConfig.APIToken,
						)
					}
					for _, client := range clients {
						if err := me.read(dlConfig, resName, client, repId); err != nil {
							return err
						}
					}
				}
			}
		}
	}

	return nil
}

func (me ResourceData) read(dlConfig DownloadConfig, resourceName string, client StandardClient, replacedIds []*ReplacedID) error {
	resources := Resources{}
	ids, err := client.LIST()
	if err != nil {
		return err
	}
	nameCounter := NewNameCounter()
	resNameCounter := NewNameCounter().Replace(ResourceName)
	for idx, id := range ids {
		if replacedIds != nil && !containsRepId(replacedIds, id) {
			continue
		}
		if dlConfig.ResourceNames[resourceName] != nil {
			if (!dlConfig.Exclude && !dlConfig.MatchID(resourceName, id)) || (dlConfig.Exclude && dlConfig.MatchID(resourceName, id)) {
				continue
			}
		}
		if dlConfig.Verbose {
			fmt.Println("  ", id, "[", idx+1, "of", len(ids), "]")
		}

		config, err := client.GET(id)
		if err != nil {
			return err
		}
		if marshaller, ok := config.(hcl.Marshaler); ok {
			name := ResourceInfoMap[resourceName].Name(dlConfig, resourceName, config, nameCounter)
			resource := Resource{
				ID:         id,
				Name:       name,
				RESTObject: marshaller,
				UniqueName: resNameCounter.Numbering(escape(name)),
			}
			resources[resource.UniqueName] = &resource
		}
	}
	if _, found := me[resourceName]; found {
		for k, v := range resources {
			me[resourceName][k] = v
		}
	} else {
		me[resourceName] = resources
	}

	return nil
}

func (me ResourceData) readNoList(resourceName string, client NoListClient) error {
	resNameCounter := NewNameCounter().Replace(ResourceName)
	resources := Resources{}
	config, err := client.GET()
	if err != nil {
		return err
	}
	if marshaller, ok := config.(hcl.Marshaler); ok {
		name := ResourceInfoMap[resourceName].Name(DownloadConfig{}, resourceName, config, nil)
		resource := Resource{
			ID:         uuid.New().String(),
			RESTObject: marshaller,
			Name:       name,
			UniqueName: resNameCounter.Numbering(escape(name)),
		}
		resources[resource.UniqueName] = &resource
	}
	me[resourceName] = resources

	return nil
}

func (me ResourceData) readDashboards(dlConfig DownloadConfig, resourceName string, client StandardClient, replacedIds []*ReplacedID) error {
	resources := Resources{}
	ids, err := client.LIST()

	dashboardSharing := Resources{}
	shareRestClient := sharing.NewService(dlConfig.EnvironmentURL+"/api/config/v1", dlConfig.APIToken)
	if err != nil {
		return err
	}

	nameCounter := NewNameCounter()
	resNameCounter := NewNameCounter().Replace(ResourceName)
	for idx, id := range ids {
		if replacedIds != nil && !containsRepId(replacedIds, id) {
			continue
		}
		if dlConfig.ResourceNames[resourceName] != nil {
			if (!dlConfig.Exclude && !dlConfig.MatchID(resourceName, id)) || (dlConfig.Exclude && dlConfig.MatchID(resourceName, id)) {
				continue
			}
		}

		if dlConfig.Verbose {
			fmt.Println("  ", id, "[", idx+1, "of", len(ids), "]")
		}

		config, err := client.GET(id)
		if err != nil {
			return err
		}
		var name string
		var uniqueName string
		if marshaller, ok := config.(hcl.Marshaler); ok {
			name = ResourceInfoMap[resourceName].Name(dlConfig, resourceName, config, nameCounter)
			uniqueName = resNameCounter.Numbering(escape(name))
			resource := Resource{
				ID:         id,
				RESTObject: marshaller,
				Name:       name,
				UniqueName: uniqueName,
			}
			resources[resource.UniqueName] = &resource
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
				ID:         id,
				RESTObject: marshaller,
				Name:       name,
				UniqueName: uniqueName,
			}
			dataObj := resource.RESTObject.(*sharing.DashboardSharing)
			if dataObj.Preset {
				continue
			}
			dataObj.DashboardID = "HCL-UNQUOTE-dynatrace_dashboard." + uniqueName + ".id"
			dashboardSharing[resource.UniqueName] = &resource
		}
	}
	me[resourceName] = resources

	if len(dashboardSharing) > 0 {
		me["dynatrace_dashboard_sharing"] = dashboardSharing
	}

	return nil
}

func (me ResourceData) readKeyRequests(dlConfig DownloadConfig, resourceName string, replacedIds []*ReplacedID) error {
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
	resNameCounter := NewNameCounter().Replace(ResourceName)
	for _, service := range services.(servicetopology.Services) {
		keyRequestID, keyRequest, err := restClient.LISTSVC(service.EntityId)

		kr := keyRequest.(*keyrequests.KeyRequest)
		if err != nil {
			return err
		}
		if kr == nil {
			continue
		}
		if replacedIds != nil && !containsRepId(replacedIds, keyRequestID) {
			continue
		}
		if dlConfig.ResourceNames[resourceName] != nil {
			if (!dlConfig.Exclude && !dlConfig.MatchID(resourceName, keyRequestID)) || (dlConfig.Exclude && dlConfig.MatchID(resourceName, keyRequestID)) {
				continue
			}
		}
		if marshaller, ok := keyRequest.(hcl.Marshaler); ok {
			name := ResourceInfoMap[resourceName].Name(dlConfig, resourceName, service, nil)
			resource := Resource{
				ID:         keyRequestID,
				RESTObject: marshaller,
				Name:       name,
				UniqueName: resNameCounter.Numbering(escape(name)),
			}
			resources[resource.UniqueName] = &resource
		}
	}
	me[resourceName] = resources

	return nil
}

func containsRepId(replacedId []*ReplacedID, id string) bool {
	for _, repId := range replacedId {
		if strings.HasPrefix(repId.ID, "GEOLOCATION-") {
			repId.Processed = true
		}
	}

	for _, repId := range replacedId {
		if id == repId.ID {
			repId.Processed = true
			return true
		}
	}
	return false
}

func containsProcessedRepId(replacedId []*ReplacedID) bool {
	for _, repId := range replacedId {
		if !repId.Processed {
			return false
		}
	}
	return true
}
