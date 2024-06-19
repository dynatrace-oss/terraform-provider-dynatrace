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

package mobile

import (
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api"
	mysettings "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/metrics/calculated/mobile/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings/services/httpcache"
)

const SchemaID = "v1:config:calculated-metrics-mobile"

func Service(credentials *settings.Credentials) settings.CRUDService[*mysettings.CalculatedMobileMetric] {
	return &service{client: httpcache.DefaultClient(credentials.URL, credentials.Token, SchemaID)}
}

type service struct {
	client rest.Client
}

func (me *service) Get(id string, v *mysettings.CalculatedMobileMetric) error {
	return me.client.Get(fmt.Sprintf("/api/config/v1/calculatedMetrics/mobile/%s", url.PathEscape(id)), 200).Finish(v)
}

func (me *service) List() (api.Stubs, error) {
	var err error

	req := me.client.Get("/api/config/v1/calculatedMetrics/mobile", 200)
	var stubList api.StubList
	if err = req.Finish(&stubList); err != nil {
		return nil, err
	}
	return stubList.Values, nil
}

func (me *service) Validate(v *mysettings.CalculatedMobileMetric) error {
	var err error
	client := me.client

	req := client.Post("/api/config/v1/calculatedMetrics/mobile/validator", v, 204)
	if err = req.Finish(); err != nil {
		return err
	}

	return nil
}

func (me *service) Create(v *mysettings.CalculatedMobileMetric) (*api.Stub, error) {
	var err error
	client := me.client
	var stub api.Stub

	req := client.Post("/api/config/v1/calculatedMetrics/mobile", v, 201)
	if err = req.Finish(&stub); err != nil {
		return nil, err
	}

	return &stub, nil
}

func (me *service) Update(id string, v *mysettings.CalculatedMobileMetric) error {
	if err := me.client.Put(fmt.Sprintf("/api/config/v1/calculatedMetrics/mobile/%s", url.PathEscape(id)), v, 204).Finish(); err != nil {
		return err
	}
	return nil
}

func (me *service) Delete(id string) error {
	var err error
	attempts := 30

	for i := 0; i < attempts; i++ {
		if err = me.client.Delete(fmt.Sprintf("/api/config/v1/calculatedMetrics/mobile/%s", url.PathEscape(id)), 204, 200).Finish(); err != nil {
			if strings.Contains(err.Error(), fmt.Sprintf("Metric with key \"%s\" does not exist", id)) {
				return nil
			}
			if !strings.Contains(err.Error(), fmt.Sprintf("Unable to delete Mrum metric with key: \"%s\" from DemMetricsConfigPersistence", id)) {
				return err
			}
		} else {
			break
		}
		time.Sleep(2 * time.Second)
	}

	if err != nil {
		return err
	}
	return nil
}

func (me *service) SchemaID() string {
	return SchemaID
}
