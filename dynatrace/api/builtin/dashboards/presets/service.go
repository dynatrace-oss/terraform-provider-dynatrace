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
	"sync"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api"
	presets "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/dashboards/presets/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings/services/settings20"
)

const SchemaVersion = "0.9.13"
const SchemaID = "builtin:dashboards.presets"

var mu sync.Mutex

func Service(credentials *settings.Credentials) settings.CRUDService[*presets.Settings] {
	return &service{settings20.Service[*presets.Settings](credentials, SchemaID, SchemaVersion)}
}

type service struct {
	service settings.CRUDService[*presets.Settings]
}

func (me *service) List() (api.Stubs, error) {
	mu.Lock()
	defer mu.Unlock()
	return me.service.List()
}

func (me *service) Get(id string, v *presets.Settings) error {
	mu.Lock()
	defer mu.Unlock()
	if stubs, _ := me.service.List(); len(stubs) > 0 {
		return me.service.Get(stubs[0].ID, v)
	}
	return me.service.Get(id, v)
}

func (me *service) SchemaID() string {
	return me.service.SchemaID()
}

func (me *service) Create(v *presets.Settings) (*api.Stub, error) {
	mu.Lock()
	defer mu.Unlock()
	// This schema is flagged with `multiobject=false` - in other words only one
	// object can exist per environment
	// Instead of trying to CREATE the settings we simply update the existing one
	if stubs, _ := me.service.List(); len(stubs) > 0 {
		return stubs[0], me.update(stubs[0].ID, v)
	}
	return me.service.Create(v)
}

func (me *service) update(id string, v *presets.Settings) error {
	return me.service.Update(id, v)
}

func (me *service) Update(id string, v *presets.Settings) error {
	mu.Lock()
	defer mu.Unlock()
	// Just in case we LOCALLY are having a different ID than the ID the
	// environment insists on having we're checking first, whether an object
	// already exists. If so, we're updating using THAT ID.
	if stubs, _ := me.service.List(); len(stubs) > 0 {
		return me.update(stubs[0].ID, v)
	}
	return me.service.Update(id, v)
}

func (me *service) Delete(id string) error {
	mu.Lock()
	defer mu.Unlock()
	return me.service.Delete(id)
}
