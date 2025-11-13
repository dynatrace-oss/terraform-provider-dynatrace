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

var defaultCreateConfirm = 8
var createConfirm = settings.GetIntEnv("DYNATRACE_CREATE_CONFIRM_SYNTHETIC_MONITORS_V2", defaultCreateConfirm, 1, 50)

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

	// The code below validates that the resource is fully created and available
	// First loop: There could be a small delay until the API returns a successful GET after create
	// Second loop: Tag configuration takes a while to propagate across the cluster
	successes := 0
	for {
		if err = client.Get(ctx, fmt.Sprintf("%s/%s", BasePath, url.PathEscape(resp.EntityId)), 200).Finish(); err == nil {
			successes = successes + 1
			if successes >= createConfirm {
				break
			}
			time.Sleep(time.Millisecond * 100)
		} else {
			successes = 0
			time.Sleep(time.Millisecond * 500)
		}
	}

	if len(v.Tags) > 0 {
		successes = 0
		for {
			validateMonitor := monitors.Settings{}
			if err = client.Get(ctx, fmt.Sprintf("%s/%s", BasePath, url.PathEscape(resp.EntityId)), 200).Finish(&validateMonitor); err == nil {
				if len(v.Tags) == len(validateMonitor.Tags) {
					successes = successes + 1
					if successes >= createConfirm {
						break
					}
				} else {
					successes = 0
				}
				time.Sleep(time.Millisecond * 100)
			}
		}
	}

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
	return rest.APITokenClient(me.credentials).Put(ctx, fmt.Sprintf("%s/%s", BasePath, url.PathEscape(id)), v, 200).Finish()
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
