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

package mobile

import (
	"context"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings"

	mobile "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/applications/mobile/settings"
)

const SchemaID = "v1:config:applications:mobile"

func Service(credentials *rest.Credentials) settings.CRUDService[*mobile.Application] {
	return settings.NewAPITokenService(
		credentials,
		SchemaID,
		&settings.ServiceOptions[*mobile.Application]{
			Get:           settings.Path("/api/config/v1/applications/mobile/%s"),
			List:          settings.Path("/api/config/v1/applications/mobile"),
			CreateConfirm: 20,
			OnChanged:     StoreKeyUserActionsAndSessionProperties,
			CompleteGet:   LoadKeyUserActionsAndSessionProperties,
			OnBeforeUpdate: func(id string, v *mobile.Application) error {
				v.ApplicationType = nil // Application Type cannot be changed on existing Mobile Apps
				v.ApplicationID = nil   // Application ID cannot be changed on existing Mobile Apps
				return nil
			},
			Duplicates: Duplicates,
		},
	)
}

func Duplicates(ctx context.Context, service settings.RService[*mobile.Application], v *mobile.Application) (*api.Stub, error) {
	if settings.RejectDuplicate("dynatrace_mobile_application") {
		var err error
		var stubs api.Stubs
		if stubs, err = service.List(ctx); err != nil {
			return nil, err
		}
		for _, stub := range stubs {
			if v.Name == stub.Name {
				return nil, fmt.Errorf("Mobile Application named '%s' already exists", v.Name)
			}
		}
	} else if settings.HijackDuplicate("dynatrace_mobile_application") {
		var err error
		var stubs api.Stubs
		if stubs, err = service.List(ctx); err != nil {
			return nil, err
		}
		for _, stub := range stubs {
			if v.Name == stub.Name {
				return stub, nil
			}
		}
	}
	return nil, nil
}

func LoadKeyUserActionsAndSessionProperties(ctx context.Context, client rest.Client, id string, v *mobile.Application) error {
	var err error
	resp := struct {
		KeyUserActions []struct {
			Name string `json:"name"`
		} `json:"keyUserActions"`
	}{}
	if err = client.Get(ctx, fmt.Sprintf("/api/config/v1/applications/mobile/%s/keyUserActions", id), 200).Finish(&resp); err != nil {
		return err
	}
	names := []string{}
	for _, item := range resp.KeyUserActions {
		names = append(names, item.Name)
	}
	if len(names) > 0 {
		v.KeyUserActions = names
	}

	presp := struct {
		SessionProperties []struct {
			Key         string `json:"key"`
			DisplayName string `json:"displayName"`
		} `json:"sessionProperties"`
		UserActionProperties []struct {
			Key         string `json:"key"`
			DisplayName string `json:"displayName"`
		} `json:"userActionProperties"`
	}{}
	if err = client.Get(ctx, fmt.Sprintf("/api/config/v1/applications/mobile/%s/userActionAndSessionProperties", id), 200).Finish(&presp); err != nil {
		return err
	}
	remoteProperties := map[string]*mobile.UserActionAndSessionProperty{}
	propKeys := map[string]string{}
	for _, v := range presp.SessionProperties {
		propKeys[v.Key] = v.Key
	}
	for _, v := range presp.UserActionProperties {
		propKeys[v.Key] = v.Key
	}
	for propKey := range propKeys {
		var property mobile.UserActionAndSessionProperty
		if err = client.Get(ctx, fmt.Sprintf("/api/config/v1/applications/mobile/%s/userActionAndSessionProperties/%s", id, url.PathEscape(propKey)), 200).Finish(&property); err != nil {
			return err
		}
		remoteProperties[propKey] = &property
	}
	if len(remoteProperties) > 0 {
		v.Properties = mobile.UserActionAndSessionProperties{}
		for _, property := range remoteProperties {
			v.Properties = append(v.Properties, property)
		}
	}
	return nil
}

func StoreKeyUserActionsAndSessionProperties(ctx context.Context, client rest.Client, id string, v *mobile.Application) error {
	var err error

	req := client.Get(ctx, fmt.Sprintf("/api/config/v1/applications/mobile/%s/keyUserActions", id), 200)
	remoteKeyUserActions := map[string]string{}
	resp := struct {
		KeyUserActions []struct {
			Name string `json:"name"`
		} `json:"keyUserActions"`
	}{}
	if err = req.Finish(&resp); err != nil {
		return err
	}
	for _, item := range resp.KeyUserActions {
		remoteKeyUserActions[item.Name] = item.Name
	}
	keyUserActionsToDelete := map[string]string{}
	for keyUserAction := range remoteKeyUserActions {
		keyUserActionsToDelete[keyUserAction] = keyUserAction
	}
	keyUserActionsToAdd := []string{}
	applicationConfig := v
	if len(applicationConfig.KeyUserActions) > 0 {
		for _, keyUserAction := range applicationConfig.KeyUserActions {
			delete(keyUserActionsToDelete, keyUserAction)
			if _, found := remoteKeyUserActions[keyUserAction]; !found {
				keyUserActionsToAdd = append(keyUserActionsToAdd, keyUserAction)
			}
		}
	}
	for keyUserAction := range keyUserActionsToDelete {
		req := client.Delete(ctx, fmt.Sprintf("/api/config/v1/applications/mobile/%s/keyUserActions/%s", url.PathEscape(id), url.PathEscape(keyUserAction)), 204)
		if err = req.Finish(); err != nil {
			return err
		}
	}
	for _, keyUserAction := range keyUserActionsToAdd {
		req := client.Post(ctx, fmt.Sprintf("/api/config/v1/applications/mobile/%s/keyUserActions/%s", url.PathEscape(id), url.PathEscape(keyUserAction)), map[string]any{}, 200)
		if err = req.Finish(); err != nil {
			return err
		}
	}

	if len(keyUserActionsToAdd) > 0 {
		var maxTries = 40
		var successes = 0
		var requiredSuccesses = 5

		var response = struct {
			Entities []struct {
				DisplayName string `json:"displayName"`
			} `json:"entities"`
			TotalCount int `json:"totalCount"`
		}{}

		for i := 0; i < maxTries; i++ {
			if err = client.Get(ctx, fmt.Sprintf(`/api/v2/entities?pageSize=4000&from=now-3y&&entitySelector=type("DEVICE_APPLICATION_METHOD"),fromRelationships.isDeviceApplicationMethodOf(entityId("%s"))&fields=fromRelationships`, id), 200).Finish(&response); err != nil {
				return err
			}

			success := true
			for _, kua := range applicationConfig.KeyUserActions {
				found := false
				for _, respEntity := range response.Entities {
					if kua == respEntity.DisplayName {
						found = true
					}
				}
				if !found {
					success = false
					break
				}
			}

			if success {
				successes++
				if successes >= requiredSuccesses {
					break
				}
				time.Sleep(200 * time.Millisecond)
				continue
			} else {
				successes = 0
				time.Sleep(10 * time.Second)
			}
		}
	}

	getPropertiesResponse := struct {
		SessionProperties []struct {
			Key         string `json:"key"`
			DisplayName string `json:"displayName"`
		} `json:"sessionProperties"`
		UserActionProperties []struct {
			Key         string `json:"key"`
			DisplayName string `json:"displayName"`
		} `json:"userActionProperties"`
	}{}
	if err = client.Get(ctx, fmt.Sprintf("/api/config/v1/applications/mobile/%s/userActionAndSessionProperties", id), 200).Finish(&getPropertiesResponse); err != nil {
		return err
	}
	propKeys := map[string]string{}
	for _, v := range getPropertiesResponse.SessionProperties {
		propKeys[v.Key] = v.Key
	}
	for _, v := range getPropertiesResponse.UserActionProperties {
		propKeys[v.Key] = v.Key
	}
	remoteProperties := map[string]*mobile.UserActionAndSessionProperty{}
	for propKey := range propKeys {
		var property mobile.UserActionAndSessionProperty
		if err = client.Get(ctx, fmt.Sprintf("/api/config/v1/applications/mobile/%s/userActionAndSessionProperties/%s", url.PathEscape(id), url.PathEscape(propKey)), 200).Finish(&property); err != nil {
			return err
		}
		remoteProperties[propKey] = &property
	}
	propsToDelete := map[string]string{}
	for propKey := range remoteProperties {
		propsToDelete[propKey] = propKey
	}
	propsToUpdate := map[string]*mobile.UserActionAndSessionProperty{}
	propsToCreate := map[string]*mobile.UserActionAndSessionProperty{}
	if len(applicationConfig.Properties) > 0 {
		for _, property := range applicationConfig.Properties {
			propKey := property.Key
			delete(propsToDelete, propKey)
			if _, found := remoteProperties[propKey]; found {
				propsToUpdate[propKey] = property
			} else {
				propsToCreate[propKey] = property
			}
		}
	}
	for propKey := range propsToDelete {
		if err = client.Delete(ctx, fmt.Sprintf("/api/config/v1/applications/mobile/%s/userActionAndSessionProperties/%s", url.PathEscape(id), url.PathEscape(propKey)), 204).Finish(); err != nil {
			return err
		}
	}
	for _, property := range propsToCreate {
		if err = client.Post(ctx, fmt.Sprintf("/api/config/v1/applications/mobile/%s/userActionAndSessionProperties", url.PathEscape(id)), property, 201, 204).Finish(); err != nil {
			return err
		}
	}
	for propKey, property := range propsToUpdate {
		if err = client.Put(ctx, fmt.Sprintf("/api/config/v1/applications/mobile/%s/userActionAndSessionProperties/%s", id, url.PathEscape(propKey)), property, 201, 204).Finish(); err != nil {
			if !strings.Contains(err.Error(), "No Content (PUT)") {
				return err
			}
		}
	}
	return nil
}
