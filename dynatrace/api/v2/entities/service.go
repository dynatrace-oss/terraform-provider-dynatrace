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
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/shutdown"

	entities "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v2/entities/settings"
)

const SchemaID = "v2:environment:entities"

func Service(entityType string, entityName string, entitySelector string, from string, to string, credentials *settings.Credentials) settings.RService[*entities.Settings] {
	return &service{entityType: entityType, entityName: entityName, entitySelector: entitySelector, from: from, to: to, client: rest.DefaultClient(credentials.URL, credentials.Token)}
}

type service struct {
	client         rest.Client
	entityType     string
	entityName     string
	entitySelector string
	from           string
	to             string
}

func (me *service) Get(ctx context.Context, id string, v *entities.Settings) (err error) {
	from := me.from
	to := me.to
	if len(from) == 0 {
		from = "now-3y"
	}
	if len(to) > 0 {
		to = fmt.Sprintf("&to=%s", url.QueryEscape(to))
	}

	var dataObj entities.Settings
	if len(me.entitySelector) > 0 {
		if err = me.client.Get(ctx, fmt.Sprintf(`/api/v2/entities?pageSize=4000%s&from=%s&entitySelector=%s&fields=tags,properties,lastSeenTms`, to, url.QueryEscape(from), url.QueryEscape(me.entitySelector)), 200).Finish(&dataObj); err != nil {
			return err
		}
	} else {
		entitySelector := ""
		if len(me.entityType) > 0 {
			entitySelector = entitySelector + fmt.Sprintf("type(\"%s\")", me.entityType)
		}
		if len(me.entityName) > 0 {
			entitySelector = entitySelector + fmt.Sprintf(",entityName.equals(\"%s\")", me.entityName)
		}
		// shouldn't happen - just sanity
		// in case there was no type but a name
		entitySelector = strings.TrimPrefix(entitySelector, ",")
		if err = me.client.Get(ctx, fmt.Sprintf(`/api/v2/entities?pageSize=4000%s&from=%s&entitySelector=%s&fields=tags,properties,lastSeenTms`, to, url.QueryEscape(from), url.QueryEscape(entitySelector)), 200).Finish(&dataObj); err != nil {
			return err
		}
	}
	if shutdown.System.Stopped() {
		return nil
	}
	if dataObj.NextPageKey != nil {
		key := *dataObj.NextPageKey
		for {
			var tempObj entities.Settings
			if err = me.client.Get(ctx, fmt.Sprintf("/api/v2/entities?nextPageKey=%s", url.PathEscape(key)), 200).Finish(&tempObj); err != nil {
				return err
			}
			dataObj.Entities = append(dataObj.Entities, tempObj.Entities...)
			if tempObj.NextPageKey == nil {
				break
			}
			key = *tempObj.NextPageKey
		}
	}
	*v = dataObj
	return nil
}

func (me *service) SchemaID() string {
	return fmt.Sprintf("%s-%s", SchemaID, me.entityType)
}

func (me *service) List(ctx context.Context) (api.Stubs, error) {
	return api.Stubs{&api.Stub{ID: me.SchemaID(), Name: me.SchemaID()}}, nil
}
