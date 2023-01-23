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

package entities

import (
	"fmt"
	"net/url"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings"

	entities "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v2/entities/settings"
)

const SchemaID = "v2:environment:entities"

func Service(entityType string, credentials *settings.Credentials) settings.RService[*entities.Settings] {
	return &service{entityType: entityType, client: rest.DefaultClient(credentials.URL, credentials.Token)}
}

type service struct {
	client     rest.Client
	entityType string
}

func (me *service) Get(id string, v *entities.Settings) (err error) {
	var dataObj entities.Settings
	if err = me.client.Get(fmt.Sprintf(`/api/v2/entities?pageSize=100&entitySelector=type("%s")&fields=tags`, url.QueryEscape(me.entityType)), 200).Finish(&dataObj); err != nil {
		return err
	}
	if dataObj.NextPageKey != nil {
		key := *dataObj.NextPageKey
		for {
			var tempObj entities.Settings
			if err = me.client.Get(fmt.Sprintf("/api/v2/entities?nextPageKey=%s", url.PathEscape(key)), 200).Finish(&tempObj); err != nil {
				return err
			}
			dataObj.Entities = append(dataObj.Entities, tempObj.Entities...)
			if tempObj.NextPageKey == nil {
				break
			}
			key = *tempObj.NextPageKey
		}
	}
	return nil
}

func (me *service) SchemaID() string {
	return fmt.Sprintf("%s-%s", SchemaID, me.entityType)
}

func (me *service) List() (settings.Stubs, error) {
	return settings.Stubs{&settings.Stub{ID: me.SchemaID(), Name: me.SchemaID()}}, nil
}
