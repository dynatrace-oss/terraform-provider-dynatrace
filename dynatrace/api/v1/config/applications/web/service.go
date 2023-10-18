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

package web

import (
	"fmt"
	"net/url"
	"time"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings"

	web "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/applications/web/settings"
)

const SchemaID = "v1:config:applications:web"

var DefaultCreateConfirmTimeout = 60
var createConfirmTimeout = settings.GetIntEnv("DYNATRACE_CREATE_CONFIRM_WEB_APPLICATION", DefaultCreateConfirmTimeout, 20, 300)

func Service(credentials *settings.Credentials) settings.CRUDService[*web.Application] {
	return settings.NewCRUDService(
		credentials,
		SchemaID,
		&settings.ServiceOptions[*web.Application]{
			Get:           settings.Path("/api/config/v1/applications/web/%s"),
			List:          settings.Path("/api/config/v1/applications/web"),
			CreateConfirm: createConfirmTimeout,
			CompleteGet:   LoadKeyUserActions,
			OnChanged:     SaveKeyUserActions,
			Duplicates:    Duplicates,
		},
	)
}

func Duplicates(service settings.RService[*web.Application], v *web.Application) (*api.Stub, error) {
	if settings.RejectDuplicate("dynatrace_web_application") {
		var err error
		var stubs api.Stubs
		if stubs, err = service.List(); err != nil {
			return nil, err
		}
		for _, stub := range stubs {
			if v.Name == stub.Name {
				return nil, fmt.Errorf("Web Application named '%s' already exists", v.Name)
			}
		}
	} else if settings.HijackDuplicate("dynatrace_web_application") {
		var err error
		var stubs api.Stubs
		if stubs, err = service.List(); err != nil {
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

func SaveKeyUserActions(client rest.Client, id string, v *web.Application) error {
	var err error
	req := client.Get(fmt.Sprintf("/api/config/v1/applications/web/%s/keyUserActions", url.PathEscape(id))).Expect(200)
	kual := struct {
		KeyUserActions []struct {
			ID     string                `json:"meIdentifier"`
			Name   string                `json:"name"`
			Type   web.KeyUserActionType `json:"actionType"`
			Domain *string               `json:"domain,omitempty"`
		} `json:"keyUserActionList"`
	}{}
	if err = req.Finish(&kual); err != nil {
		return err
	}
	remoteKeyUserActions := []*web.KeyUserAction{}
	for _, kua := range kual.KeyUserActions {
		remoteKeyUserActions = append(remoteKeyUserActions, &web.KeyUserAction{Name: kua.Name, Type: kua.Type, Domain: kua.Domain})
	}
	keyUserActionsToDelete := []string{}
	keyUserActionsToCreate := []*web.KeyUserAction{}

	for _, kua := range v.KeyUserActions {
		found := false
		for _, rkua := range remoteKeyUserActions {
			if kua.Equals(rkua) {
				found = true
				break
			}
		}
		if !found {
			keyUserActionsToCreate = append(keyUserActionsToCreate, kua)
		}
	}
	for _, rxkua := range kual.KeyUserActions {
		rkua := &web.KeyUserAction{Name: rxkua.Name, Type: rxkua.Type, Domain: rxkua.Domain}
		found := false
		for _, kua := range v.KeyUserActions {
			if rkua.Equals(kua) {
				found = true
				break
			}
		}
		if !found {
			keyUserActionsToDelete = append(keyUserActionsToDelete, rxkua.ID)
		}
	}
	for _, kuaID := range keyUserActionsToDelete {
		req := client.Delete(fmt.Sprintf("/api/config/v1/applications/web/%s/keyUserActions/%s", url.PathEscape(id), url.PathEscape(kuaID)), 204)
		if err = req.Finish(); err != nil {
			return err
		}
	}
	for _, keyUserAction := range keyUserActionsToCreate {
		req := client.Post(fmt.Sprintf("/api/config/v1/applications/web/%s/keyUserActions", url.PathEscape(id)), keyUserAction, 201)
		if err = req.Finish(); err != nil {
			return err
		}
	}

	if len(keyUserActionsToCreate) > 0 {
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
			if err = client.Get(fmt.Sprintf(`/api/v2/entities?pageSize=4000&from=now-3y&&entitySelector=type("APPLICATION_METHOD"),fromRelationships.isApplicationMethodOf(entityId("%s"))&fields=fromRelationships`, id), 200).Finish(&response); err != nil {
				return err
			}

			success := true
			for _, kua := range v.KeyUserActions {
				found := false
				for _, respEntity := range response.Entities {
					if kua.Name == respEntity.DisplayName {
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

	return nil
}

func LoadKeyUserActions(client rest.Client, id string, v *web.Application) error {
	var err error
	var kual web.KeyUserActionList
	req := client.Get(fmt.Sprintf("/api/config/v1/applications/web/%s/keyUserActions", url.PathEscape(id)), 200)
	if err = req.Finish(&kual); err != nil {
		return err
	}
	actions := []*web.KeyUserAction{}
	for _, action := range kual.KeyUserActions {
		actions = append(actions, action)
	}
	if len(actions) > 0 {
		v.KeyUserActions = actions
	}
	return nil
}
