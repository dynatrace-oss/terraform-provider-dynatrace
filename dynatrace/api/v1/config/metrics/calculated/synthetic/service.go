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

package synthetic

import (
	"time"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api"
	mysettings "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/metrics/calculated/synthetic/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings/services/httpcache"
)

const SchemaID = "v1:config:calculated-metrics-synthetic"
const BasePath = "/api/config/v1/calculatedMetrics/synthetic"

func Service(credentials *settings.Credentials) settings.CRUDService[*mysettings.CalculatedSyntheticMetric] {
	return &service{service: settings.NewCRUDService(
		credentials,
		SchemaID,
		settings.DefaultServiceOptions[*mysettings.CalculatedSyntheticMetric](BasePath),
	), client: httpcache.DefaultClient(credentials.URL, credentials.Token, SchemaID)}
}

type service struct {
	service settings.CRUDService[*mysettings.CalculatedSyntheticMetric]
	client  rest.Client
}

func (me *service) List() (api.Stubs, error) {
	return me.service.List()
}

func (me *service) Get(id string, v *mysettings.CalculatedSyntheticMetric) error {
	return me.service.Get(id, v)
}

func (me *service) SchemaID() string {
	return me.service.SchemaID()
}

func (me *service) Create(v *mysettings.CalculatedSyntheticMetric) (*api.Stub, error) {
	return me.service.Create(v)
}

func (me *service) Update(id string, v *mysettings.CalculatedSyntheticMetric) error {
	return me.service.Update(id, v)
}

func (me *service) Delete(id string) error {
	var err error
	var retry = 10

	for i := 0; i < retry; i++ {
		if err = me.service.Delete(id); err != nil {
			return err
		}
		if err = me.service.Get(id, new(mysettings.CalculatedSyntheticMetric)); err != nil {
			break
		}
		time.Sleep(time.Second * 2)
	}

	return nil
}
