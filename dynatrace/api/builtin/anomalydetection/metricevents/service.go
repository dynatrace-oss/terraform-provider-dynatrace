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

package metricevents

import (
	"context"
	"fmt"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api"
	metricevents "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/anomalydetection/metricevents/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings/services/settings20"
)

const SchemaID = "builtin:anomaly-detection.metric-events"
const SchemaVersion = "1.0.16"

func Service(credentials *settings.Credentials) settings.CRUDService[*metricevents.Settings] {
	return settings20.Service(credentials, SchemaID, SchemaVersion, &settings20.ServiceOptions[*metricevents.Settings]{LegacyID: settings.LegacyObjIDDecode, Duplicates: Duplicates})
}

func Duplicates(ctx context.Context, service settings.RService[*metricevents.Settings], v *metricevents.Settings) (*api.Stub, error) {
	if settings.RejectDuplicate("dynatrace_metric_events") {
		var err error
		var stubs api.Stubs
		if stubs, err = service.List(ctx); err != nil {
			return nil, err
		}
		for _, stub := range stubs {
			if v.Summary == stub.Name {
				return nil, fmt.Errorf("Metric event named '%s' already exists", v.Summary)
			}
		}
	} else if settings.HijackDuplicate("dynatrace_metric_events") {
		var err error
		var stubs api.Stubs
		if stubs, err = service.List(ctx); err != nil {
			return nil, err
		}
		for _, stub := range stubs {
			if v.Summary == stub.Name {
				return stub, nil
			}
		}
	}
	return nil, nil
}
