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

package http

import (
	"context"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/opt"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/synthetic/monitors"
	http "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/synthetic/monitors/http/settings"
)

const SchemaID = "v1:synthetic:monitors:http"
const BasePath = "/api/v1/synthetic/monitors"

func Service(credentials *rest.Credentials) settings.CRUDService[*http.SyntheticMonitor] {
	return &service{service: settings.NewAPITokenService(credentials, SchemaID, &settings.ServiceOptions[*http.SyntheticMonitor]{
		Get:            settings.Path("/api/v1/synthetic/monitors/%s"),
		List:           settings.Path("/api/v1/synthetic/monitors?type=HTTP"),
		CreateURL:      func(v *http.SyntheticMonitor) string { return "/api/v1/synthetic/monitors" },
		Stubs:          &monitors.Monitors{},
		HasNoValidator: true,
		CreateConfirm:  30,
	})}
}

type service struct {
	service settings.CRUDService[*http.SyntheticMonitor]
}

func GetTempScript() *http.Script {
	return &http.Script{Requests: http.Requests{&http.Request{
		Description: opt.NewString("--terraform-auto-generated--"),
		URL:         "http://localhost",
		Method:      "OPTIONS",
	}}}
}

func (me *service) Create(ctx context.Context, v *http.SyntheticMonitor) (*api.Stub, error) {
	if v.NoScript != nil && *v.NoScript && v.Script == nil {
		v.Script = GetTempScript()
	}
	return me.service.Create(ctx, v)
}

func (me *service) Update(ctx context.Context, id string, v *http.SyntheticMonitor) error {
	if v.NoScript != nil && *v.NoScript && v.Script == nil {
		monitorSettings := new(http.SyntheticMonitor)
		if err := me.service.Get(ctx, id, monitorSettings); err != nil {
			return err
		}
		v.Script = monitorSettings.Script
	}
	return me.service.Update(ctx, id, v)
}

func (me *service) Delete(ctx context.Context, id string) error {
	return me.service.Delete(ctx, id)
}

func (me *service) Get(ctx context.Context, id string, v *http.SyntheticMonitor) error {
	if err := me.service.Get(ctx, id, v); err != nil {
		return err
	}
	if v.Script != nil && len(v.Script.Requests) == 1 && *v.Script.Requests[0].Description == *GetTempScript().Requests[0].Description {
		v.NoScript = opt.NewBool(true)
		v.Script = nil
	}
	return nil
}

func (me *service) List(ctx context.Context) (api.Stubs, error) {
	return me.service.List(ctx)
}

func (me *service) SchemaID() string {
	return me.service.SchemaID()
}
