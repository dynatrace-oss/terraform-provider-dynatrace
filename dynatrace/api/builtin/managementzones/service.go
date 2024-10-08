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

package managementzones

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings/services/settings20"

	managementzones "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/managementzones/settings"
	slo "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/monitoring/slo"
	slosettings "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/monitoring/slo/settings"
)

const SchemaID = "builtin:management-zones"
const SchemaVersion = "1.0.12"

func Service(credentials *settings.Credentials) settings.CRUDService[*managementzones.Settings] {
	return &service{
		service:     settings20.Service(credentials, SchemaID, SchemaVersion, &settings20.ServiceOptions[*managementzones.Settings]{LegacyID: settings.LegacyLongDecode}),
		client:      rest.DefaultClient(credentials.URL, credentials.Token),
		credentials: credentials,
	}
}

type service struct {
	service     settings.ListIDCRUDService[*managementzones.Settings]
	client      rest.Client
	credentials *settings.Credentials
}

const DefaultNumRequiredSuccesses = 5
const MinNumRequiredSuccesses = 5
const MaxNumRequiredSuccesses = 25

const DefaultMaxConfirmationRetries = 50
const MaxMaxConfirmationRetries = 150
const MinMaxConfirmationRetries = 50

func getEnv(key string, def int, min int, max int) int {
	value := os.Getenv(key)
	if len(value) == 0 {
		return def
	}
	iValue, err := strconv.Atoi(strings.TrimSpace(value))
	if err != nil {
		return def
	}
	if iValue > max {
		iValue = max
	}
	if iValue < min {
		iValue = min
	}
	return iValue
}

func (me *service) Create(ctx context.Context, v *managementzones.Settings) (*api.Stub, error) {
	stub, err := me.service.Create(ctx, v)
	if err != nil {
		return nil, err
	}

	// The Dynatrace API is unprepared for the management zone to be utilized right after create despite a subsequent successful GET.
	// Utilizing the SLO endpoint, the following code does a validate only call with a dependency of the created management zone.
	sloValue := slosettings.Settings{
		Enabled:          false,
		Name:             "TerraformManagementZoneTest",
		MetricName:       "terraform_management_zone_test",
		MetricExpression: "builtin:apps.web.action.count.load.browser",
		EvaluationType:   "AGGREGATE",
		Filter:           fmt.Sprintf("mzName(%s)", v.Name),
		EvaluationWindow: "-1w",
		TargetSuccess:    99.98,
		TargetWarning:    99.99,
		ErrorBudgetBurnRate: &slosettings.ErrorBudgetBurnRate{
			BurnRateVisualizationEnabled: false,
		},
	}

	validator := slo.Service(me.credentials).(settings.Validator[*slosettings.Settings])

	maxConfirmationRetries := getEnv("DT_MGMZ_RETRIES", DefaultMaxConfirmationRetries, MinMaxConfirmationRetries, MaxMaxConfirmationRetries)
	numRequiredSuccesses := getEnv("DT_MGMZ_SUCCESSES", DefaultNumRequiredSuccesses, MinNumRequiredSuccesses, MaxNumRequiredSuccesses)
	success := 0
	for i := 0; i < maxConfirmationRetries; i++ {
		if err := validator.Validate(ctx, &sloValue); err == nil {
			success++
			if success >= numRequiredSuccesses {
				break
			}
		} else {
			success = 0
			time.Sleep(500)
		}
	}

	return stub, nil
}

func (me *service) Update(ctx context.Context, id string, v *managementzones.Settings) error {
	return me.service.Update(ctx, id, v)
}

func (me *service) Delete(ctx context.Context, id string) error {
	return me.service.Delete(ctx, id)
}

func (me *service) Get(ctx context.Context, id string, v *managementzones.Settings) error {
	return me.service.Get(ctx, id, v)
}

func (me *service) List(ctx context.Context) (api.Stubs, error) {
	return me.service.List(ctx)
}

func (me *service) ListIDs(ctx context.Context) (api.Stubs, error) {
	return me.service.ListIDs(ctx)
}

func (me *service) SchemaID() string {
	return me.service.SchemaID()
}
