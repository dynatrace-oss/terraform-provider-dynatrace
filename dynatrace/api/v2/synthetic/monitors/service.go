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

package monitors

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api"
	monitors "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v2/synthetic/monitors/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings"
)

const SchemaID = "v2:synthetic:monitors:network"
const BasePath = "/api/v2/synthetic/monitors"

// TimeForUpdateConsistency Time to wait after create/update operations to allow for consistency in subsequent GET operations
var TimeForUpdateConsistency = 5 * time.Second

func Service(credentials *rest.Credentials) settings.CRUDService[*monitors.Settings] {
	return &service{credentials: credentials}
}

type service struct {
	credentials *rest.Credentials
}

func (me *service) Create(ctx context.Context, v *monitors.Settings) (*api.Stub, error) {
	var err error

	resp := struct {
		EntityId string `json:"entityId"`
	}{}
	client := rest.APITokenClient(me.credentials)
	if err = client.Post(ctx, BasePath, v, 201).Finish(&resp); err != nil {
		return nil, err
	}

	time.Sleep(TimeForUpdateConsistency)
	return &api.Stub{ID: resp.EntityId, Name: v.Name}, nil
}

func (me *service) Get(ctx context.Context, id string, v *monitors.Settings) error {
	if err := rest.APITokenClient(me.credentials).Get(ctx, fmt.Sprintf("%s/%s", BasePath, url.PathEscape(id)), 200).Finish(v); err != nil {
		return err
	}

	return nil
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

	if err = rest.APITokenClient(me.credentials).Get(ctx, BasePath, 200).Finish(&monitors); err != nil {
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
	err := rest.APITokenClient(me.credentials).Put(ctx, fmt.Sprintf("%s/%s", BasePath, url.PathEscape(id)), v, 200).Finish()
	if err != nil {
		return err
	}
	time.Sleep(TimeForUpdateConsistency)
	return nil
}

func (me *service) Delete(ctx context.Context, id string) error {
	return rest.APITokenClient(me.credentials).Delete(ctx, fmt.Sprintf("%s/%s", BasePath, url.PathEscape(id)), 204).Finish()
}

func (me *service) New() *monitors.Settings {
	return new(monitors.Settings)
}

func (me *service) SchemaID() string {
	return SchemaID
}
