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

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings"

	web "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/applications/web/settings"
)

const SchemaID = "v1:config:applications:web"

func Service(credentials *settings.Credentials) settings.CRUDService[*web.Application] {
	return settings.NewCRUDService(
		credentials,
		SchemaID,
		&settings.ServiceOptions[*web.Application]{
			Get:           settings.Path("/api/config/v1/applications/web/%s"),
			List:          settings.Path("/api/config/v1/applications/web"),
			CreateConfirm: 20,
			CompleteGet:   LoadKeyUserActions,
			OnChanged:     SaveKeyUserActions,
		},
	)
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

	successes := 0
	numRequiredSuccesses := 20
	for {
		testKual := web.KeyUserActionList{}
		req = client.Get(fmt.Sprintf("/api/config/v1/applications/web/%s/keyUserActions", url.PathEscape(id)), 200)
		if err = req.Finish(&testKual); err != nil {
			return err
		}

		if len(testKual.KeyUserActions) == len(v.KeyUserActions) {
			successes = successes + 1
			if successes >= numRequiredSuccesses {
				break
			}
			time.Sleep(time.Millisecond * 200)
		} else {
			successes = 0
			time.Sleep(time.Second * 10)
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
