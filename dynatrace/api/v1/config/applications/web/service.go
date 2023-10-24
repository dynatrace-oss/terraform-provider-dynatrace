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
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings/services/httpcache"

	web "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/applications/web/settings"
)

const SchemaID = "v1:config:applications:web"

var DefaultCreateConfirmTimeout = 60
var createConfirmTimeout = settings.GetIntEnv("DYNATRACE_CREATE_CONFIRM_WEB_APPLICATION", DefaultCreateConfirmTimeout, 20, 300)

func Service(credentials *settings.Credentials) settings.CRUDService[*web.Application] {
	return &service{
		service: settings.NewCRUDService(
			credentials,
			SchemaID,
			&settings.ServiceOptions[*web.Application]{
				Get:           settings.Path("/api/config/v1/applications/web/%s"),
				List:          settings.Path("/api/config/v1/applications/web"),
				CreateConfirm: createConfirmTimeout,
				Duplicates:    Duplicates,
			},
		),
		client: httpcache.DefaultClient(credentials.URL, credentials.Token, SchemaID),
	}
}

type service struct {
	service settings.CRUDService[*web.Application]
	client  rest.Client
}

func (me *service) List() (api.Stubs, error) {
	return me.service.List()
}

func (me *service) GetWithContext(ctx context.Context, id string, v *web.Application) error {
	var stateKeyUserActions web.KeyUserActions
	cfg := ctx.Value(settings.ContextKeyStateConfig)
	if appConfig, ok := cfg.(*web.Application); ok {
		stateKeyUserActions = appConfig.KeyUserActions
	}
	if err := me.service.Get(id, v); err != nil {
		return err
	}
	var err error
	var kual web.KeyUserActionList
	req := me.client.Get(fmt.Sprintf("/api/config/v1/applications/web/%s/keyUserActions", url.PathEscape(id)), 200)
	if err = req.Finish(&kual); err != nil {
		return err
	}
	actions := web.KeyUserActions{}
	if len(stateKeyUserActions) > 0 {
		for _, stateKeyUserAction := range stateKeyUserActions {
			for _, onlineKeyUserAction := range kual.KeyUserActions {
				if stateKeyUserAction.Equals(onlineKeyUserAction) {
					actions = append(actions, stateKeyUserAction)
					break
				}
			}
		}
	}
	if len(actions) > 0 {
		v.KeyUserActions = actions
	}
	return nil
}

func (me *service) Get(id string, v *web.Application) error {
	return me.service.Get(id, v)
}

func (me *service) SchemaID() string {
	return me.service.SchemaID()
}

func (me *service) CreateWithContext(ctx context.Context, v *web.Application) (*api.Stub, error) {
	stub, err := me.service.Create(v)
	if err != nil {
		return stub, err
	}
	if len(v.KeyUserActions) > 0 {
		for _, keyUserAction := range v.KeyUserActions {
			req := me.client.Post(fmt.Sprintf("/api/config/v1/applications/web/%s/keyUserActions", url.PathEscape(stub.ID)), keyUserAction, 201)
			if err = req.Finish(); err != nil {
				return nil, err
			}
		}
	}
	return stub, me.pollUntilKeyUserActionsCreated(stub.ID, v.KeyUserActions)
}

func (me *service) Create(v *web.Application) (*api.Stub, error) {
	return me.service.Create(v)
}

func (me *service) UpdateWithContext(ctx context.Context, id string, v *web.Application) error {
	var stateKeyUserActions web.KeyUserActions
	cfg := ctx.Value(settings.ContextKeyStateConfig)
	if appConfig, ok := cfg.(*web.Application); ok {
		stateKeyUserActions = appConfig.KeyUserActions
	}
	if err := me.service.Update(id, v); err != nil {
		return err
	}
	var err error
	var remoteKeyUserActions map[string]*web.KeyUserAction
	if remoteKeyUserActions, err = me.fetchKeyUserActions(id); err != nil {
		return err
	}
	keyUserActionsToCreate := web.KeyUserActions{}
	if len(v.KeyUserActions) > 0 {
		for _, configuredKeyUserAction := range v.KeyUserActions {
			found := false
			for _, remoteKeyUserAction := range remoteKeyUserActions {
				if remoteKeyUserAction.Equals(configuredKeyUserAction) {
					found = true
					break
				}
			}
			if !found {
				keyUserActionsToCreate = append(keyUserActionsToCreate, configuredKeyUserAction)
			}
		}
	}
	for _, keyUserAction := range keyUserActionsToCreate {
		req := me.client.Post(fmt.Sprintf("/api/config/v1/applications/web/%s/keyUserActions", url.PathEscape(id)), keyUserAction, 201)
		if err = req.Finish(); err != nil {
			return err
		}
	}

	// from now on we only take key user actions into consideration that were present in the state
	if len(stateKeyUserActions) > 0 {
		// we don't want to delete key user actions not known within the state
		for remoteKeyUserActionID, remoteKeyUserAction := range remoteKeyUserActions {
			found := false
			for _, stateKeyUserAction := range stateKeyUserActions {
				if remoteKeyUserAction.Equals(stateKeyUserAction) {
					found = true
					break
				}
			}
			if !found {
				delete(remoteKeyUserActions, remoteKeyUserActionID)
			}
		}
		// key user actions present remotely but absent within the configuration need to get deleted
		keyUserActionsToDelete := []string{}
		if len(v.KeyUserActions) == 0 {
			for keyUserActionID := range remoteKeyUserActions {
				keyUserActionsToDelete = append(keyUserActionsToDelete, keyUserActionID)
			}
		} else {
			for keyUserActionID, remoteKeyUserAction := range remoteKeyUserActions {
				found := false
				for _, configuredKeyUserAction := range v.KeyUserActions {
					if remoteKeyUserAction.Equals(configuredKeyUserAction) {
						found = true
						break
					}
				}
				if !found {
					keyUserActionsToDelete = append(keyUserActionsToDelete, keyUserActionID)
				}
			}
		}
		// execute deletions
		for _, keyUserActionID := range keyUserActionsToDelete {
			var err error
			req := me.client.Delete(fmt.Sprintf("/api/config/v1/applications/web/%s/keyUserActions/%s", url.PathEscape(id), url.PathEscape(keyUserActionID)), 204)
			if err = req.Finish(); err != nil {
				return err
			}
		}
	}
	return me.pollUntilKeyUserActionsCreated(id, keyUserActionsToCreate)
}

// to stay consisent we need to poll
// newly created key user actions may not be "online" right away
func (me *service) pollUntilKeyUserActionsCreated(id string, keyUserActionsToCreate web.KeyUserActions) error {

	var err error
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
			if err = me.client.Get(fmt.Sprintf(`/api/v2/entities?pageSize=4000&from=now-3y&&entitySelector=type("APPLICATION_METHOD"),fromRelationships.isApplicationMethodOf(entityId("%s"))&fields=fromRelationships`, id), 200).Finish(&response); err != nil {
				return err
			}

			success := true
			for _, kua := range keyUserActionsToCreate {
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

func (me *service) fetchKeyUserActions(id string) (map[string]*web.KeyUserAction, error) {
	actions := map[string]*web.KeyUserAction{}
	var err error
	req := me.client.Get(fmt.Sprintf("/api/config/v1/applications/web/%s/keyUserActions", url.PathEscape(id))).Expect(200)
	kual := struct {
		KeyUserActions []struct {
			ID     string                `json:"meIdentifier"`
			Name   string                `json:"name"`
			Type   web.KeyUserActionType `json:"actionType"`
			Domain *string               `json:"domain,omitempty"`
		} `json:"keyUserActionList"`
	}{}
	if err = req.Finish(&kual); err != nil {
		return nil, err
	}
	for _, kua := range kual.KeyUserActions {
		actions[kua.ID] = &web.KeyUserAction{Name: kua.Name, Type: kua.Type, Domain: kua.Domain}
	}
	return actions, nil
}

func (me *service) Update(id string, v *web.Application) error {
	return me.service.Update(id, v)
}

func (me *service) Delete(id string) error {
	return me.service.Delete(id)
}

func (me *service) Name() string {
	return me.service.Name()
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
				return nil, fmt.Errorf("a Web Application named '%s' already exists", v.Name)
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
