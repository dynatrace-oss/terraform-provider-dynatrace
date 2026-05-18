//go:build unit

/**
* @license
* Copyright 2026 Dynatrace LLC
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

package metricevents_test

import (
	"context"
	"testing"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api"
	metricevents "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/anomalydetection/metricevents/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/export"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/provider/config"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/resources"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// captureService implements settings.CRUDService[*metricevents.Settings] and
// records every payload it receives, so tests can assert on what the resource
// lifecycle would have sent to the Settings 2.0 API.
type captureService struct {
	Stub *api.Stub
	Err  error

	CreateValue *metricevents.Settings
	UpdateID    string
	UpdateValue *metricevents.Settings
	DeleteID    string
	GetValue    *metricevents.Settings
}

func (s *captureService) List(_ context.Context) (api.Stubs, error) { return nil, nil }

func (s *captureService) Get(_ context.Context, _ string, v *metricevents.Settings) error {
	if s.GetValue != nil {
		*v = *s.GetValue
	}
	return s.Err
}

func (s *captureService) SchemaID() string { return "builtin:anomaly-detection.metric-events" }

func (s *captureService) Create(_ context.Context, v *metricevents.Settings) (*api.Stub, error) {
	s.CreateValue = v
	return s.Stub, s.Err
}

func (s *captureService) Update(_ context.Context, id string, v *metricevents.Settings) error {
	s.UpdateID = id
	s.UpdateValue = v
	return s.Err
}

func (s *captureService) Delete(_ context.Context, id string) error {
	s.DeleteID = id
	return s.Err
}

func newTestGeneric(s *captureService) *resources.Generic {
	factory := func(*rest.Credentials) settings.CRUDService[*metricevents.Settings] { return s }
	return &resources.Generic{
		Type:       export.ResourceTypes.MetricEvents,
		Descriptor: export.NewResourceDescriptor(factory),
	}
}

// baseInput returns a minimal valid raw map for dynatrace_metric_events.
// Tests merge in the query_definition (and any other field) needed to exercise
// their scenario.
func baseInput() map[string]any {
	return map[string]any{
		"enabled": true,
		"summary": "test event",
		"event_template": []any{
			map[string]any{
				"title":       "title",
				"description": "description",
				"event_type":  "CUSTOM_ALERT",
			},
		},
		"model_properties": []any{
			map[string]any{
				"alert_condition":    "ABOVE",
				"alert_on_no_data":   false,
				"dealerting_samples": 5,
				"samples":            5,
				"violating_samples":  3,
				"type":               "STATIC_THRESHOLD",
				"threshold":          0.0,
			},
		},
	}
}

func providerConfig() *config.ProviderConfiguration {
	return &config.ProviderConfiguration{EnvironmentURL: "https://example.com"}
}

// TestMetricEvents_Update_DropsStaleEmptyEntityFilter is the regression test
// for the bug where switching `type` from METRIC_KEY to METRIC_SELECTOR shipped
// an empty `entityFilter` to the API. An empty entity_filter can survive in
// state when (a) the resource was first created with type=METRIC_KEY without
// an entity_filter block (UnmarshalHCL auto-creates one), and (b) the
// DiffSuppressFunc on entity_filter masks the 1->0 removal during plan.
// The fix in UnmarshalHCL normalises that empty filter back to nil; this test
// pins the contract by exercising the resource Update with a `d` shaped like
// the one the SDK would hand us after suppression.
func TestMetricEvents_Update_DropsStaleEmptyEntityFilter(t *testing.T) {
	service := &captureService{Stub: &api.Stub{ID: "test-id"}}
	gen := newTestGeneric(service)

	input := baseInput()
	input["query_definition"] = []any{
		map[string]any{
			"type":            "METRIC_SELECTOR",
			"metric_selector": "builtin:host.cpu.usage:avg",
			// Empty entity_filter block — what survives in `d` when the
			// DiffSuppressFunc masks the removal of an auto-created filter.
			"entity_filter": []any{
				map[string]any{},
			},
		},
	}

	d := schema.TestResourceDataRaw(t, gen.Resource().Schema, input)
	d.SetId("test-id")

	diags := gen.Update(t.Context(), d, providerConfig())
	require.False(t, diags.HasError(), "diagnostics: %v", diags)
	require.NotNil(t, service.UpdateValue, "Update was not called on the service")

	qd := service.UpdateValue.QueryDefinition
	require.NotNil(t, qd, "QueryDefinition missing from Update payload")
	assert.Nil(t, qd.EntityFilter, "empty entity_filter must be dropped when type=METRIC_SELECTOR; otherwise the API rejects the payload")
}
