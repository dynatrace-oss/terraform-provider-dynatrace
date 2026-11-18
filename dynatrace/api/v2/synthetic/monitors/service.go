/**
* @license
* Copyright 2026 Dynatrace LLC
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

package monitors

import (
	"context"
	"fmt"
	"net/url"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api"
	monitors "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v2/synthetic/monitors/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings"
)

const SchemaID = "v2:synthetic:monitors:network"
const BasePath = "/api/v2/synthetic/monitors"

var _ settings.ReadAfterWriteSkipper = new(service)

func Service(clientSet rest.ClientSet) (settings.CRUDService[*monitors.Settings], error) {
	return &service{clientSet: clientSet}, nil
}

type service struct {
	clientSet rest.ClientSet
}

func (me *service) SkipReadAfterWrite() bool {
	// There is a huge eventual consistency time. Therefore, we ignore and skip the GET after create/update
	return true
}

func (me *service) Create(ctx context.Context, v *monitors.Settings) (*api.Stub, error) {
	var err error

	resp := struct {
		EntityId string `json:"entityId"`
	}{}
	client := rest.APITokenClient(me.clientSet)
	if err = client.Post(ctx, BasePath, v).Finish(&resp); err != nil {
		return nil, err
	}

	return &api.Stub{ID: resp.EntityId, Name: v.Name}, nil
}

func (me *service) Get(ctx context.Context, id string, v *monitors.Settings) error {
	if err := rest.APITokenClient(me.clientSet).Get(ctx, fmt.Sprintf("%s/%s", BasePath, url.PathEscape(id))).Finish(v); err != nil {
		return err
	}
	me.filterRuleBasedTags(v)

	return nil
}

// filterRuleBasedTags filters out tags coming from auto-tagging-rules and other sources
// Only USER can be modified and should be returned
func (me *service) filterRuleBasedTags(v *monitors.Settings) {
	if len(v.Tags) == 0 {
		return
	}

	tags := make(monitors.TagsWithSourceInfo, 0, len(v.Tags))
	for _, t := range v.Tags {
		if t.Source == nil || *t.Source == monitors.TagSources.User {
			tags = append(tags, t)
		}
	}
	v.Tags = tags
}

type monitorList struct {
	Monitors []struct {
		EntityId string `json:"entityId"`
		Name     string `json:"name"`
		Type     string `json:"Type"`
		Enabled  bool   `json:"Enabled"`
	} `json:"monitors"`
}

func (me *service) List(ctx context.Context) (api.Stubs, error) {
	var err error
	var monitors monitorList

	if err = rest.APITokenClient(me.clientSet).Get(ctx, BasePath).Finish(&monitors); err != nil {
		return nil, err
	}
	stubs := api.Stubs{}
	for _, monitor := range monitors.Monitors {
		stubs = append(stubs, &api.Stub{ID: monitor.EntityId, Name: monitor.Name})
	}
	return stubs, nil
}

func (me *service) Validate(v *monitors.Settings) error {
	return nil // no endpoint for that
}

func (me *service) Update(ctx context.Context, id string, v *monitors.Settings) error {
	return rest.APITokenClient(me.clientSet).Put(ctx, fmt.Sprintf("%s/%s", BasePath, url.PathEscape(id)), v).Finish()
}

func (me *service) Delete(ctx context.Context, id string) error {
	return rest.APITokenClient(me.clientSet).Delete(ctx, fmt.Sprintf("%s/%s", BasePath, url.PathEscape(id))).Finish()
}

func (me *service) New() *monitors.Settings {
	return new(monitors.Settings)
}

func (me *service) SchemaID() string {
	return SchemaID
}
