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

package keyuseractions

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/export/multiuse"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings/services/cache"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings/services/httpcache"

	webservice "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/applications/web"
	keyuseractions "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/applications/web/keyuseractions/settings"
	web "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/applications/web/settings"
)

const SchemaID = "v1:config:applications:web:keyuseractions"

func Service(credentials *settings.Credentials) settings.CRUDService[*keyuseractions.Settings] {
	return &service{
		client:        httpcache.DefaultClient(credentials.URL, credentials.Token, SchemaID),
		webAppService: cache.CRUD(webservice.Service(credentials), true)}
}

type service struct {
	client        rest.Client
	webAppService settings.CRUDService[*web.Application]
}

type KeyUserActionsList struct {
	Values []KeyUserAction `json:"keyUserActionList"`
}

type KeyUserAction struct {
	Name         string `json:"name"`
	ActionType   string `json:"actionType"`
	Domain       string `json:"domain"`
	MEIdentifier string `json:"meIdentifier"`
}

type QueryForWebAppIDResponse struct {
	Entities []struct {
		EntityID string `json:"entityId"`
	} `json:"entities"`
}

func (me *service) Get(ctx context.Context, id string, v *keyuseractions.Settings) error {

	applicationID, realId, err := getIDs(ctx, id, me)
	if err != nil {
		return err
	}

	id = realId

	var kuaList KeyUserActionsList
	if err := me.client.Get(ctx, fmt.Sprintf("/api/config/v1/applications/web/%s/keyUserActions", url.PathEscape(applicationID)), 200).Finish(&kuaList); err != nil {
		return err
	}
	if len(kuaList.Values) > 0 {
		for _, keyUserAction := range kuaList.Values {
			if keyUserAction.MEIdentifier == id {
				v.ApplicationID = applicationID
				if len(keyUserAction.Domain) > 0 {
					v.Domain = &keyUserAction.Domain
				}
				v.Name = keyUserAction.Name
				v.Type = keyuseractions.Type(keyUserAction.ActionType)
				return nil
			}
		}
	}
	return rest.Error{Code: 404, Message: fmt.Sprintf("Key User Action with ID '%s' not found", id)}
}

func getIDs(ctx context.Context, id string, me *service) (string, string, error) {
	found, applicationID, realId := multiuse.ExtractIDParent(id)
	if found {
		return applicationID, realId, nil
	}

	applicationID, err := me.fetchApplicationID(ctx, id)
	if err != nil {
		return "", "", err
	}

	return applicationID, id, nil
}

func (me *service) fetchApplicationID(ctx context.Context, id string) (string, error) {
	var err error
	var response QueryForWebAppIDResponse
	entitySelector := fmt.Sprintf("type(APPLICATION),toRelationships.isApplicationMethodOf(type(APPLICATION_METHOD),entityId(%s))", id)
	if err = me.client.Get(ctx, fmt.Sprintf(`/api/v2/entities?from=now-10y&entitySelector=%s&fields=fromRelationships`, url.QueryEscape(entitySelector)), 200).Finish(&response); err != nil {
		return "", err
	}
	if len(response.Entities) == 0 {
		return "", rest.Error{Code: 404, Message: fmt.Sprintf("Key User Action with ID '%s' not found", id)}
	}
	return response.Entities[0].EntityID, nil
}

func (me *service) Validate(v *keyuseractions.Settings) error {
	return nil
}

func (me *service) Update(ctx context.Context, id string, v *keyuseractions.Settings) error {
	stub, err := me.Create(ctx, v)
	if err != nil {
		return err
	}
	if stub.ID != id {
		return fmt.Errorf("updating key user action '%s' for application '%s' unexpectedly created a new entity", id, v.ApplicationID)
	}
	return nil
}

func (me *service) Delete(ctx context.Context, id string) error {
	applicationID, err := me.fetchApplicationID(ctx, id)
	if err != nil {
		return nil
	}
	me.client.Delete(ctx, fmt.Sprintf("/api/config/v1/applications/web/%s/keyUserActions/%s", url.PathEscape(applicationID), url.PathEscape(id)), 204).Finish()
	var maxTries = 100
	var successes = 0
	var requiredSuccesses = 10

	for i := 0; i < maxTries; i++ {
		var kua keyuseractions.Settings
		if err := me.Get(ctx, id, &kua); err != nil {
			if resterr, ok := err.(rest.Error); ok {
				if resterr.Code == 404 {
					successes++
					if successes >= requiredSuccesses {
						break
					}
					time.Sleep(time.Duration(200+i*100) * time.Millisecond)
					continue
				}
			} else {
				successes = 0
				me.client.Delete(ctx, fmt.Sprintf("/api/config/v1/applications/web/%s/keyUserActions/%s", url.PathEscape(applicationID), url.PathEscape(id)), 204).Finish()
				time.Sleep(time.Duration(200+i*100) * time.Millisecond)
			}
		} else {
			successes = 0
			me.client.Delete(ctx, fmt.Sprintf("/api/config/v1/applications/web/%s/keyUserActions/%s", url.PathEscape(applicationID), url.PathEscape(id)), 204).Finish()
			time.Sleep(time.Duration(200+i*100) * time.Millisecond)
		}
	}
	return nil
}

type KeyUserActionCreateResponse struct {
	ID string `json:"id"`
}

func (me *service) Create(ctx context.Context, v *keyuseractions.Settings) (*api.Stub, error) {
	applicationID := v.ApplicationID
	var createReponse KeyUserActionCreateResponse
	if err := me.client.Post(ctx, fmt.Sprintf("/api/config/v1/applications/web/%s/keyUserActions", url.PathEscape(applicationID)), v, 201).Finish(&createReponse); err != nil {
		return nil, err
	}
	stub := &api.Stub{ID: createReponse.ID}
	var maxTries = 100
	var successes = 0
	var requiredSuccesses = 10

	for i := 0; i < maxTries; i++ {
		var kua keyuseractions.Settings
		if err := me.Get(ctx, stub.ID, &kua); err == nil {
			if kua.Name == v.Name {
				successes++
				if successes >= requiredSuccesses {
					break
				}
				time.Sleep(time.Duration(200+i*100) * time.Millisecond)
				continue
			} else {
				successes = 0
				me.client.Post(ctx, fmt.Sprintf("/api/config/v1/applications/web/%s/keyUserActions", url.PathEscape(applicationID)), v, 201).Finish(&createReponse)
				time.Sleep(time.Duration(200+i*100) * time.Millisecond)
			}
		} else {
			successes = 0
			me.client.Post(ctx, fmt.Sprintf("/api/config/v1/applications/web/%s/keyUserActions", url.PathEscape(applicationID)), v, 201).Finish(&createReponse)
			time.Sleep(time.Duration(200+i*100) * time.Millisecond)
		}
	}
	return stub, nil
}

func (me *service) List(ctx context.Context) (api.Stubs, error) {
	var err error
	var stubs api.Stubs
	var appStubs api.Stubs
	if appStubs, err = me.webAppService.List(ctx); err != nil {
		return nil, err
	}
	for _, appStub := range appStubs {
		var kuaList KeyUserActionsList
		if err := me.client.Get(ctx, fmt.Sprintf("/api/config/v1/applications/web/%s/keyUserActions", url.PathEscape(appStub.ID)), 200).Finish(&kuaList); err != nil {
			return nil, err
		}
		for _, keyUserAction := range kuaList.Values {
			stubs = append(stubs, &api.Stub{
				ID:       keyUserAction.MEIdentifier,
				Name:     fmt.Sprintf("KeyUserAction " + keyUserAction.Name + " for " + appStub.Name),
				ParentID: &appStub.ID,
			})
		}
	}
	return stubs.ToStubs(), nil
}

func (me *service) SchemaID() string {
	return SchemaID
}
