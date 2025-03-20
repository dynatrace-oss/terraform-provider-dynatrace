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

package presets

import (
	"context"
	"sync"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api"
	presets "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/dashboards/presets/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings/services/settings20"
)

const SchemaVersion = "0.9.13"
const SchemaID = "builtin:dashboards.presets"

var mu sync.Mutex

func Service(credentials *rest.Credentials) settings.CRUDService[*presets.Settings] {
	return &service{settings20.Service[*presets.Settings](credentials, SchemaID, SchemaVersion)}
}

type service struct {
	service settings.CRUDService[*presets.Settings]
}

func (me *service) List(ctx context.Context) (api.Stubs, error) {
	mu.Lock()
	defer mu.Unlock()
	return me.service.List(ctx)
}

func (me *service) Get(ctx context.Context, id string, v *presets.Settings) error {
	mu.Lock()
	defer mu.Unlock()
	if stubs, _ := me.service.List(ctx); len(stubs) > 0 {
		return me.service.Get(ctx, stubs[0].ID, v)
	}
	return me.service.Get(ctx, id, v)
}

func (me *service) SchemaID() string {
	return me.service.SchemaID()
}

func (me *service) Create(ctx context.Context, v *presets.Settings) (*api.Stub, error) {
	mu.Lock()
	defer mu.Unlock()
	// This schema is flagged with `multiobject=false` - in other words only one
	// object can exist per environment
	// Instead of trying to CREATE the settings we simply update the existing one
	if stubs, _ := me.service.List(ctx); len(stubs) > 0 {
		return stubs[0], me.update(ctx, stubs[0].ID, v)
	}
	return me.service.Create(ctx, v)
}

func (me *service) update(ctx context.Context, id string, v *presets.Settings) error {
	return me.service.Update(ctx, id, v)
}

func (me *service) Update(ctx context.Context, id string, v *presets.Settings) error {
	mu.Lock()
	defer mu.Unlock()
	// Just in case we LOCALLY are having a different ID than the ID the
	// environment insists on having we're checking first, whether an object
	// already exists. If so, we're updating using THAT ID.
	if stubs, _ := me.service.List(ctx); len(stubs) > 0 {
		return me.update(ctx, stubs[0].ID, v)
	}
	return me.service.Update(ctx, id, v)
}

func (me *service) Delete(ctx context.Context, id string) error {
	mu.Lock()
	defer mu.Unlock()
	return me.service.Delete(ctx, id)
}
