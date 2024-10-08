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

package items

import (
	"context"
	"fmt"
	"net/url"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api"
	items "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v2/hub/items/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/opt"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings"
)

const SchemaID = "v2:environment:hubitems"

type Options struct {
	Type string
}

func Service(credentials *settings.Credentials, opts Options) settings.RService[*items.HubItemList] {
	return &service{client: rest.DefaultClient(credentials.URL, credentials.Token), opts: opts}
}

type service struct {
	client rest.Client
	opts   Options
}

func (me *service) Get(ctx context.Context, id string, v *items.HubItemList) error {
	queries := map[string]string{}
	if len(me.opts.Type) > 0 {
		queries["itemType"] = me.opts.Type
	}
	nextPageKey := opt.NewString("first")
	for nextPageKey != nil && *nextPageKey != "" {
		hubItemList := items.HubItemList{}
		queryString := ""
		if len(queries) > 0 {
			for key, value := range queries {
				if len(queryString) == 0 {
					queryString = "?"
				} else {
					queryString = queryString + "&"
				}
				queryString = queryString + url.QueryEscape(key) + "=" + url.QueryEscape(value)
			}
		}
		if err := me.client.Get(ctx, fmt.Sprintf(`/api/v2/hub/items%s`, queryString), 200).Finish(&hubItemList); err != nil {
			return err
		}
		v.Items = append(v.Items, hubItemList.Items...)
		nextPageKey = hubItemList.NextPageKey
		queries = map[string]string{}
		if nextPageKey != nil {
			queries["nextPageKey"] = *nextPageKey
		}
	}
	return nil
}

func (me *service) SchemaID() string {
	return SchemaID
}

func (me *service) List(ctx context.Context) (api.Stubs, error) {
	return api.Stubs{&api.Stub{ID: me.SchemaID(), Name: me.SchemaID()}}, nil
}
