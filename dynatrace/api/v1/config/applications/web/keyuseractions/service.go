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
	"fmt"
	"net/url"
	"time"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api"
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

// {
// 	"totalCount": 1,
// 	"pageSize": 50,
// 	"entities": [
// 	  {
// 		"entityId": "APPLICATION_METHOD-889EFBCD62C0A9EC",
// 		"type": "APPLICATION_METHOD",
// 		"displayName": "Loading of page /example",
// 		"fromRelationships": {
// 		  "isApplicationMethodOf": [
// 			{
// 			  "id": "APPLICATION-0DA2132F908F4A55",
// 			  "type": "APPLICATION"
// 			}
// 		  ],
// 		  "isApplicationMethodOfGroup": [
// 			{
// 			  "id": "APPLICATION_METHOD_GROUP-28C0BF69FAA99BF4",
// 			  "type": "APPLICATION_METHOD_GROUP"
// 			}
// 		  ]
// 		}
// 	  }
// 	]
//   }

type QueryForWebAppIDResponse struct {
	Entities []struct {
		FromRelationships struct {
			IsApplicationMethodOf []struct {
				ID string `json:"id"`
			} `json:"isApplicationMethodOf"`
		} `json:"fromRelationships"`
	} `json:"entities"`
}

func (me *service) Get(id string, v *keyuseractions.Settings) error {
	applicationID, err := me.fetchApplicationID(id)
	if err != nil {
		return err
	}
	var kuaList KeyUserActionsList
	if err := me.client.Get(fmt.Sprintf("/api/config/v1/applications/web/%s/keyUserActions", url.PathEscape(applicationID)), 200).Finish(&kuaList); err != nil {
		return err
	}
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
	return rest.Error{Code: 404, Message: fmt.Sprintf("Key User Action with ID '%s' not found", id)}
}

func (me *service) fetchApplicationID(id string) (string, error) {
	var err error
	var response QueryForWebAppIDResponse
	entitySelector := fmt.Sprintf("type(APPLICATION_METHOD),entityId(%s)", id)
	if err = me.client.Get(fmt.Sprintf(`/api/v2/entities?from=now-3y&entitySelector=%s&fields=fromRelationships`, entitySelector), 200).Finish(&response); err != nil {
		return "", err
	}
	if len(response.Entities) == 0 {
		return "", rest.Error{Code: 404, Message: fmt.Sprintf("Key User Action with ID '%s' not found", id)}
	}
	applicationMethodOfs := response.Entities[0].FromRelationships.IsApplicationMethodOf
	if len(applicationMethodOfs) == 0 {
		return "", rest.Error{Code: 404, Message: fmt.Sprintf("Key User Action with ID '%s' not found (reason: no relationship to web app found)", id)}
	}
	applicationID := applicationMethodOfs[0].ID
	if len(applicationID) == 0 {
		return "", rest.Error{Code: 404, Message: fmt.Sprintf("Key User Action with ID '%s' not found (reason: no relationship to web app found)", id)}
	}

	return applicationID, nil

}

func (me *service) Validate(v *keyuseractions.Settings) error {
	return nil
}

func (me *service) Update(id string, v *keyuseractions.Settings) error {
	stub, err := me.Create(v)
	if err != nil {
		return err
	}
	if stub.ID != id {
		return fmt.Errorf("updating key user action '%s' for application '%s' unexpectedly created a new entity", id, v.ApplicationID)
	}
	return nil
}

func (me *service) Delete(id string) error {
	applicationID, err := me.fetchApplicationID(id)
	if err != nil {
		return nil
	}
	return me.client.Delete(fmt.Sprintf("/api/config/v1/applications/web/%s/keyUserActions/%s", url.PathEscape(applicationID), url.PathEscape(id)), 204).Finish()
}

type KeyUserActionCreateResponse struct {
	ID string `json:"id"`
}

func (me *service) Create(v *keyuseractions.Settings) (*api.Stub, error) {
	applicationID := v.ApplicationID
	var createReponse KeyUserActionCreateResponse
	if err := me.client.Post(fmt.Sprintf("/api/config/v1/applications/web/%s/keyUserActions", url.PathEscape(applicationID)), v, 201).Finish(&createReponse); err != nil {
		return nil, err
	}
	stub := &api.Stub{ID: createReponse.ID}
	var maxTries = 100
	var successes = 0
	var requiredSuccesses = 5

	for i := 0; i < maxTries; i++ {
		var kua keyuseractions.Settings
		if err := me.Get(stub.ID, &kua); err == nil {
			if kua.Name == v.Name {
				successes++
				if successes >= requiredSuccesses {
					break
				}
				time.Sleep(time.Duration(200+i*20) * time.Millisecond)
				continue
			} else {
				successes = 0
				time.Sleep(time.Duration(200+i*20) * time.Millisecond)
			}
		} else {
			successes = 0
			time.Sleep(time.Duration(200+i*20) * time.Millisecond)
		}
	}
	return stub, nil
}

func (me *service) List() (api.Stubs, error) {
	var err error
	var stubs api.Stubs
	var appStubs api.Stubs
	if appStubs, err = me.webAppService.List(); err != nil {
		return nil, err
	}
	for _, appStub := range appStubs {
		var kuaList KeyUserActionsList
		if err := me.client.Get(fmt.Sprintf("/api/config/v1/applications/web/%s/keyUserActions", url.PathEscape(appStub.ID)), 200).Finish(&kuaList); err != nil {
			return nil, err
		}
		for _, keyUserAction := range kuaList.Values {
			stubs = append(stubs, &api.Stub{ID: keyUserAction.MEIdentifier, Name: fmt.Sprintf("KeyUserAction " + keyUserAction.Name + " for " + appStub.Name)})
		}
	}
	return stubs.ToStubs(), nil
}

func (me *service) SchemaID() string {
	return SchemaID
}

func (me *service) Name() string {
	return SchemaID
}
