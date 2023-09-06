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

package slo

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api"
	slo "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/monitoring/slo/settings"
	slo_env2_service "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v2/slo"
	slo_env2 "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v2/slo/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings/services/httpcache"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings/services/settings20"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/shutdown"
)

const SchemaVersion = "6.0.13"
const SchemaID = "builtin:monitoring.slo"

func Service(credentials *settings.Credentials) settings.CRUDService[*slo.Settings] {
	return &service{
		credentials: credentials,
		client:      rest.DefaultClient(credentials.URL, credentials.Token),
	}
}

type service struct {
	credentials *settings.Credentials
	client      rest.Client
}

func (me *service) Name() string {
	return me.SchemaID()
}

func (me *service) SchemaID() string {
	return SchemaID
}

func (me *service) Create(v *slo.Settings) (*api.Stub, error) {
	slo := me.convertToEnvV2(v)

	service := slo_env2_service.Service(me.credentials)
	stub, err := service.Create(&slo)
	if err != nil {
		return nil, err
	}
	stub.LegacyID = &stub.ID

	if stubs, err := me.List(); err == nil {
		for _, listStub := range stubs {
			if listStub.Name == stub.Name {
				stub.ID = listStub.ID
			}
		}
	} else {
		return nil, err
	}

	return stub, nil
}

func (me *service) Update(id string, v *slo.Settings) error {
	legacyId := settings.LegacyID(id)
	slo := me.convertToEnvV2(v)

	service := slo_env2_service.Service(me.credentials)
	err := service.Update(legacyId, &slo)
	if err != nil {
		return err
	}

	return nil
}

func (me *service) Validate(v *slo.Settings) error {
	return nil
}

func (me *service) Delete(id string) error {
	legacyId := settings.LegacyID(id)

	service := slo_env2_service.Service(me.credentials)
	err := service.Delete(legacyId)
	if err != nil {
		return err
	}

	return nil
}

type sloGet struct {
	slo_env2.SLO
	MetricKey string `json:"metricKey"`
}

func (me *service) Get(id string, v *slo.Settings) error {
	legacyId := settings.LegacyID(id)
	slo := new(sloGet)

	service := slo_env2_service.Service(me.credentials)
	client := httpcache.DefaultClient(me.credentials.URL, me.credentials.Token, service.SchemaID())
	req := client.Get(fmt.Sprintf("/api/v2/slo/%s", url.PathEscape(legacyId)), 200)
	if err := req.Finish(slo); err != nil {
		return err
	}

	*v = me.convertToSettings20(slo)
	v.LegacyID = &legacyId

	return nil
}

func (me *service) List() (api.Stubs, error) {
	return me.listSettings20()
}

func (me *service) listSettings20() (api.Stubs, error) {
	var err error

	stubs := api.Stubs{}
	nextPage := true

	var nextPageKey *string
	for nextPage {
		var sol settings20.SettingsObjectList
		var urlStr string
		if nextPageKey != nil {
			urlStr = fmt.Sprintf("/api/v2/settings/objects?nextPageKey=%s", url.QueryEscape(*nextPageKey))
		} else {
			urlStr = fmt.Sprintf("/api/v2/settings/objects?schemaIds=%s&fields=%s&pageSize=100", url.QueryEscape(me.SchemaID()), url.QueryEscape("objectId,value,scope,schemaVersion"))
		}
		req := me.client.Get(urlStr, 200)
		if err = req.Finish(&sol); err != nil {
			return nil, err
		}
		if shutdown.System.Stopped() {
			return stubs, nil
		}

		if len(sol.Items) > 0 {
			for _, item := range sol.Items {
				newItem := slo.Settings{}
				if err = json.Unmarshal(item.Value, &newItem); err != nil {
					return nil, err
				}
				settings.SetScope(&newItem, item.Scope)
				itemName := newItem.Name
				stub := &api.Stub{ID: item.ObjectID, Name: itemName, Value: newItem}
				if len(itemName) > 0 {
					stubs = append(stubs, stub)
				}
			}
		}
		nextPageKey = sol.NextPageKey
		nextPage = (nextPageKey != nil)
	}

	return stubs, nil
}

func (me *service) convertToEnvV2(v *slo.Settings) slo_env2.SLO {
	slo := slo_env2.SLO{
		Name:             v.Name,
		Enabled:          v.Enabled,
		Description:      v.CustomDescription,
		MetricExpression: &v.MetricExpression,
		EvaluationType:   string(v.EvaluationType),
		Filter:           &v.Filter,
		Target:           v.TargetSuccess,
		Warning:          v.TargetWarning,
		Timeframe:        v.EvaluationWindow,
		ErrorBudgetBurnRate: &slo_env2.ErrorBudgetBurnRate{
			BurnRateVisualizationEnabled: &v.ErrorBudgetBurnRate.BurnRateVisualizationEnabled,
			FastBurnThreshold:            v.ErrorBudgetBurnRate.FastBurnThreshold,
		},
	}
	if v.MetricName != "" {
		slo.MetricName = &v.MetricName
	}
	return slo
}

func (me *service) convertToSettings20(v1 *sloGet) slo.Settings {
	slo := &slo.Settings{
		Name:                v1.Name,
		MetricName:          strings.TrimPrefix(v1.MetricKey, "func:slo."),
		Enabled:             v1.Enabled,
		CustomDescription:   v1.Description,
		EvaluationType:      slo.SloEvaluationType(v1.EvaluationType),
		TargetSuccess:       v1.Target,
		TargetWarning:       v1.Warning,
		EvaluationWindow:    v1.Timeframe,
		ErrorBudgetBurnRate: new(slo.ErrorBudgetBurnRate),
	}
	if v1.MetricExpression != nil {
		slo.MetricExpression = *v1.MetricExpression
	}
	if v1.Filter != nil {
		slo.Filter = *v1.Filter
	}
	if v1.ErrorBudgetBurnRate.BurnRateVisualizationEnabled != nil {
		slo.ErrorBudgetBurnRate.BurnRateVisualizationEnabled = *v1.ErrorBudgetBurnRate.BurnRateVisualizationEnabled
	}
	if v1.ErrorBudgetBurnRate.FastBurnThreshold != nil {
		slo.ErrorBudgetBurnRate.FastBurnThreshold = v1.ErrorBudgetBurnRate.FastBurnThreshold
	}
	return *slo
}
