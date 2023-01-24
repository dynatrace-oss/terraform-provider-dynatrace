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

package dashboards

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings/services/cache"

	dashboards "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/dashboards/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/jsondashboards"
)

const SchemaID = "v1:config:dashboards"

func Service(credentials *settings.Credentials) settings.CRUDService[*dashboards.Dashboard] {
	return &service{service: cache.CRUD(jsondashboards.Service(credentials), true)}
}

type service struct {
	service settings.CRUDService[*dashboards.JSONDashboard]
}

func (me *service) NoCache() bool {
	return true
}

func (me *service) List() (settings.Stubs, error) {
	return me.service.List()
}

func (me *service) Get(id string, v *dashboards.Dashboard) error {
	var err error
	var data []byte
	jsondb := settings.NewSettings(me.service.(settings.RService[*dashboards.JSONDashboard]))
	if err = me.service.Get(id, jsondb); err != nil {
		return err
	}
	if data, err = settings.ToJSON(jsondb); err != nil {
		return err
	}
	return settings.FromJSON(data, v)
}

func (me *service) Validate(v *dashboards.Dashboard) error {
	var err error
	var data []byte
	jsondb := settings.NewSettings(me.service.(settings.RService[*dashboards.JSONDashboard]))
	if data, err = settings.ToJSON(v); err != nil {
		return err
	}
	if err = settings.FromJSON(data, jsondb); err != nil {
		return err
	}
	if validator, ok := me.service.(settings.Validator[*dashboards.JSONDashboard]); ok {
		return validator.Validate(jsondb)
	}
	return nil
}

func (me *service) Create(v *dashboards.Dashboard) (*settings.Stub, error) {
	var err error
	var data []byte
	jsondb := settings.NewSettings(me.service.(settings.RService[*dashboards.JSONDashboard]))
	if data, err = settings.ToJSON(v); err != nil {
		return nil, err
	}
	if err = settings.FromJSON(data, jsondb); err != nil {
		return nil, err
	}
	return me.service.Create(jsondb)
}

func (me *service) Update(id string, v *dashboards.Dashboard) error {
	var err error
	var data []byte
	jsondb := settings.NewSettings(me.service.(settings.RService[*dashboards.JSONDashboard]))
	if data, err = settings.ToJSON(v); err != nil {
		return err
	}
	if err = settings.FromJSON(data, jsondb); err != nil {
		return err
	}
	return me.service.Update(id, jsondb)
}

func (me *service) Delete(id string) error {
	return me.service.Delete(id)
}

func (me *service) SchemaID() string {
	return SchemaID
}
