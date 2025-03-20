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
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api"
	slo "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/monitoring/slo/settings"
	slo_env2_service "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v2/slo"
	slo_env2 "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v2/slo/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings/services/settings20"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/shutdown"
	"github.com/google/uuid"
)

const SchemaVersion = "6.0.14"
const SchemaID = "builtin:monitoring.slo"

func Service(credentials *rest.Credentials) settings.CRUDService[*slo.Settings] {
	return &service{credentials: credentials, client: rest.HybridClient(credentials)}
}

type service struct {
	credentials *rest.Credentials
	client      rest.Client
}

func (me *service) SchemaID() string {
	return SchemaID
}

func (me *service) Create(ctx context.Context, v *slo.Settings) (*api.Stub, error) {
	slo := me.convertToEnvV2(v)

	service := slo_env2_service.Service(me.credentials)
	stub, err := service.Create(ctx, &slo)
	if err != nil {
		return nil, err
	}
	stub.LegacyID = &stub.ID

	retries := 12
	for i := 1; i <= retries; i++ {
		if stubs, err := me.List(ctx); err == nil {
			for _, listStub := range stubs {
				if listStub.Name == stub.Name {
					stub.ID = listStub.ID
				}
			}
		}
		if len(stub.ID) > 0 {
			break
		}
		time.Sleep(5 * time.Second)
	}
	if len(stub.ID) == 0 {
		if err := service.Delete(ctx, *stub.LegacyID); err != nil {
			return nil, err
		}
		return nil, errors.New("SLO creation failed, unable to retrieve ID. Please create a GitHub issue.")
	}

	return stub, nil
}

func (me *service) Update(ctx context.Context, id string, v *slo.Settings) error {
	legacyId := settings.LegacyID(id)
	slo := me.convertToEnvV2(v)

	service := slo_env2_service.Service(me.credentials)
	err := service.Update(ctx, legacyId, &slo)
	if err != nil {
		return err
	}

	return nil
}

func (me *service) Validate(ctx context.Context, v *slo.Settings) error {
	soc := settings20.SettingsObjectCreate{
		SchemaID:      SchemaID,
		SchemaVersion: SchemaVersion,
		Scope:         "environment",
		Value:         v,
	}

	if err := me.client.Post(ctx, "/api/v2/settings/objects?validateOnly=true", []settings20.SettingsObjectCreate{soc}).Expect(200).Finish(); err != nil {
		return err
	}

	return nil
}

func (me *service) Delete(ctx context.Context, id string) error {
	legacyId := settings.LegacyID(id)

	service := slo_env2_service.Service(me.credentials)
	err := service.Delete(ctx, legacyId)
	if err != nil {
		return err
	}

	return nil
}

type sloGet struct {
	slo_env2.SLO
	MetricKey string `json:"metricKey"`
}

func (me *service) Get(ctx context.Context, id string, v *slo.Settings) error {
	err := me.get(ctx, id, v)
	if err != nil {
		if err.Error() == "Cannot access a disabled SLO." {
			settingsObject := settings20.SettingsObject{}
			req := me.client.Get(ctx, fmt.Sprintf("/api/v2/settings/objects/%s", id), 200)
			if err = req.Finish(&settingsObject); err != nil {
				return err
			}
			if err = json.Unmarshal(settingsObject.Value, v); err != nil {
				return err
			}
		}
		return err
	}
	return nil
}

func (me *service) get(ctx context.Context, id string, v *slo.Settings) error {
	var legacyId string
	if _, err := uuid.Parse(id); err == nil {
		legacyId = id
	} else {
		legacyId = settings.LegacyID(id)
	}
	if len(legacyId) == 0 {
		legacyId = id
	}

	slo := new(sloGet)

	client := rest.HybridClient(me.credentials)
	req := client.Get(ctx, fmt.Sprintf("/api/v2/slo/%s", url.PathEscape(legacyId)), 200)
	if err := req.Finish(slo); err != nil {
		return err
	}

	*v = me.convertToSettings20(slo)
	v.LegacyID = &legacyId

	return nil
}

func (me *service) List(ctx context.Context) (api.Stubs, error) {
	return me.listSettings20(ctx)
}

func (me *service) listSettings20(ctx context.Context) (api.Stubs, error) {
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
		req := me.client.Get(ctx, urlStr, 200)
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
