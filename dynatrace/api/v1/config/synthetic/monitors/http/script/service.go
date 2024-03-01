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

package script

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/synthetic/monitors/http"
	script "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/synthetic/monitors/http/script/settings"
	httpsettings "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/synthetic/monitors/http/settings"
)

const SchemaID = "v1:synthetic:monitors:http:script"
const BasePath = "/api/v1/synthetic/monitors"

func Service(credentials *settings.Credentials) settings.CRUDService[*script.Settings] {
	return &service{scriptService: settings.NewCRUDService(credentials, SchemaID, settings.DefaultServiceOptions[*script.Settings](http.BasePath)), httpService: http.Service(credentials)}
}

type service struct {
	scriptService settings.CRUDService[*script.Settings]
	httpService   settings.CRUDService[*httpsettings.SyntheticMonitor]
}

func (me *service) Create(v *script.Settings) (*api.Stub, error) {
	return &api.Stub{ID: v.HttpId, Name: v.HttpId}, me.Update(v.HttpId, v)
}

func (me *service) Update(id string, v *script.Settings) error {
	monitorSettings, err := me.getHttp(id)
	if err != nil {
		return err
	}
	monitorSettings.Script = v.Script

	return me.update(id, monitorSettings)
}

func (me *service) getHttp(id string) (*httpsettings.SyntheticMonitor, error) {
	monitorSettings := new(httpsettings.SyntheticMonitor)
	if err := me.httpService.Get(id, monitorSettings); err != nil {
		return nil, err
	}
	return monitorSettings, nil
}

func (me *service) update(id string, v *httpsettings.SyntheticMonitor) error {
	return me.httpService.Update(id, v)
}

func (me *service) Delete(id string) error {
	monitorSettings, err := me.getHttp(id)
	if err != nil {
		return err
	}
	monitorSettings.Script = http.GetTempScript()

	return me.update(id, monitorSettings)
}

func (me *service) Get(id string, v *script.Settings) error {
	monitorSettings, err := me.getHttp(id)
	if err != nil {
		return err
	}
	v.HttpId = id
	v.Script = monitorSettings.Script

	return nil
}

func (me *service) List() (api.Stubs, error) {
	return me.httpService.List()
}

func (me *service) Name() string {
	return me.scriptService.Name()
}

func (me *service) SchemaID() string {
	return me.scriptService.SchemaID()
}
