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
	"context"
	"time"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api"
	mysettings "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/metrics/calculated/synthetic/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings"
)

const SchemaID = "v1:config:calculated-metrics-synthetic"
const BasePath = "/api/config/v1/calculatedMetrics/synthetic"

func Service(credentials *rest.Credentials) settings.CRUDService[*mysettings.CalculatedSyntheticMetric] {
	return &service{service: settings.NewAPITokenService(
		credentials,
		SchemaID,
		settings.DefaultServiceOptions[*mysettings.CalculatedSyntheticMetric](BasePath),
	), client: rest.APITokenClient(credentials)}
}

type service struct {
	service settings.CRUDService[*mysettings.CalculatedSyntheticMetric]
	client  rest.Client
}

func (me *service) List(ctx context.Context) (api.Stubs, error) {
	return me.service.List(ctx)
}

func (me *service) Get(ctx context.Context, id string, v *mysettings.CalculatedSyntheticMetric) error {
	return me.service.Get(ctx, id, v)
}

func (me *service) SchemaID() string {
	return me.service.SchemaID()
}

func (me *service) Create(ctx context.Context, v *mysettings.CalculatedSyntheticMetric) (*api.Stub, error) {
	return me.service.Create(ctx, v)
}

func (me *service) Update(ctx context.Context, id string, v *mysettings.CalculatedSyntheticMetric) error {
	return me.service.Update(ctx, id, v)
}

func (me *service) Delete(ctx context.Context, id string) error {
	var err error
	var retry = 10

	for i := 0; i < retry; i++ {
		if err = me.service.Delete(ctx, id); err != nil {
			return err
		}
		if err = me.service.Get(ctx, id, new(mysettings.CalculatedSyntheticMetric)); err != nil {
			break
		}
		time.Sleep(time.Second * 2)
	}

	return nil
}
