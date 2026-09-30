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

package settings_test

import (
	"encoding/json"
	"testing"

	settings "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/extensions/dac/awsmonitoring/settings"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestMarshalWireShape pins the on-the-wire JSON shape we send to
// /platform/extensions/v2/extensions/com.dynatrace.extension.da-aws/monitoringConfigurations.
// The shape is pinned against the payload the monitoringConfigurations endpoint
// accepts for com.dynatrace.extension.da-aws.
func TestMarshalWireShape(t *testing.T) {
	s := &settings.Settings{
		Name:             "my-aws-monitoring",
		Enabled:          true,
		ExtensionVersion: "1.0.0",
		ConnectionID:     "vu9U3hXa3q0AAAABACdidWlsdGluOmh5cGVyc2NhbGVyLWF1dGhlbnRpY2F0aW9uOmF3cw",
		AccountID:        "123456789012",
		Regions:          []string{"us-east-1", "eu-central-1"},
		FeatureSets:      []string{"EC2_essential", "RDS_essential"},
	}

	raw, err := json.Marshal(s)
	require.NoError(t, err, "marshal")

	var got map[string]any
	require.NoError(t, json.Unmarshal(raw, &got), "re-decode")

	assert.Equal(t, settings.DefaultScope, got["scope"])

	value, ok := got["value"].(map[string]any)
	require.True(t, ok, "value: missing or wrong type: %T", got["value"])

	assert.Equal(t, true, value["enabled"])
	assert.Equal(t, "my-aws-monitoring", value["description"])
	assert.Equal(t, "1.0.0", value["version"])
	assert.Equal(t, "DATA_ACQUISITION", value["activationContext"])
	assert.ElementsMatch(t, []any{"EC2_essential", "RDS_essential"}, value["featureSets"])

	aws, ok := value["aws"].(map[string]any)
	require.True(t, ok, "aws block missing")
	assert.Equal(t, "us-east-1", aws["deploymentRegion"], "defaulted from first region")

	creds, ok := aws["credentials"].([]any)
	require.True(t, ok)
	require.Len(t, creds, 1)
	cred, ok := creds[0].(map[string]any)
	require.True(t, ok)
	assert.Equal(t, s.ConnectionID, cred["connectionId"])
	assert.Equal(t, "123456789012", cred["accountId"])

	mc, ok := aws["metricsConfiguration"].(map[string]any)
	require.True(t, ok, "metricsConfiguration missing")
	assert.ElementsMatch(t, []any{"us-east-1", "eu-central-1"}, mc["regions"])
}

func TestRoundTrip(t *testing.T) {
	in := &settings.Settings{
		Name:             "x",
		Enabled:          true,
		ExtensionVersion: "1.2.3",
		ConnectionID:     "conn-abc",
		AccountID:        "111122223333",
		Regions:          []string{"us-east-1"},
		FeatureSets:      []string{"S3_essential"},
		DeploymentRegion: "us-east-1",
	}
	raw, err := json.Marshal(in)
	require.NoError(t, err, "marshal")

	out := &settings.Settings{}
	require.NoError(t, json.Unmarshal(raw, out), "unmarshal")

	assert.ElementsMatch(t, in.Regions, out.Regions)
	assert.ElementsMatch(t, in.FeatureSets, out.FeatureSets)
	assert.Equal(t, in.ConnectionID, out.ConnectionID)
	assert.Equal(t, in.AccountID, out.AccountID)
	assert.Equal(t, in.Name, out.Name, "name (description)")
}

func base() *settings.Settings {
	return &settings.Settings{
		Name:             "x",
		Enabled:          true,
		ExtensionVersion: "1.0.0",
		ConnectionID:     "c",
		AccountID:        "111122223333",
		Regions:          []string{"eu-central-1"},
	}
}

func awsBlock(t *testing.T, s *settings.Settings) map[string]any {
	t.Helper()
	raw, err := json.Marshal(s)
	require.NoError(t, err, "marshal")
	var got map[string]any
	require.NoError(t, json.Unmarshal(raw, &got), "decode")
	value, ok := got["value"].(map[string]any)
	require.True(t, ok, "value block missing")
	aws, ok := value["aws"].(map[string]any)
	require.True(t, ok, "aws block missing")
	return aws
}

func TestEnumDefaults(t *testing.T) {
	aws := awsBlock(t, base())

	assert.Equal(t, "SINGLE_ACCOUNT", aws["deploymentScope"])
	assert.Equal(t, "AUTOMATED", aws["deploymentMode"])
	assert.Equal(t, "QUICK_START", aws["configurationMode"])

	// smartscape_enabled is hidden from the user-facing schema and
	// force-set to true in applyDefaults; the wire payload must reflect that
	// regardless of struct zero-value.
	sm, ok := aws["smartscapeConfiguration"].(map[string]any)
	require.True(t, ok, "smartscapeConfiguration missing/wrong type: %v", aws["smartscapeConfiguration"])
	assert.Equal(t, true, sm["enabled"], "hardcoded")
}

func TestTagFilterRoundTrip(t *testing.T) {
	s := base()
	s.TagFilters = settings.TagFilters{
		{Key: "env", Value: "prod", Condition: "INCLUDE"},
		{Key: "team", Value: "infra", Condition: "EXCLUDE"},
	}
	raw, err := json.Marshal(s)
	require.NoError(t, err, "marshal")

	out := &settings.Settings{}
	require.NoError(t, json.Unmarshal(raw, out), "unmarshal")

	assert.ElementsMatch(t, s.TagFilters, out.TagFilters)
}

func TestTagEnrichmentRoundTrip(t *testing.T) {
	s := base()
	s.TagEnrichment = []string{"owner", "cost-center"}
	raw, err := json.Marshal(s)
	require.NoError(t, err, "marshal")

	out := &settings.Settings{}
	require.NoError(t, json.Unmarshal(raw, out), "unmarshal")

	assert.ElementsMatch(t, []string{"owner", "cost-center"}, out.TagEnrichment)
}

func TestCloudWatchLogsRoundTrip(t *testing.T) {
	s := base()
	s.CloudWatchLogs = &settings.CloudWatchLogsConfig{
		Enabled: true,
		Regions: []string{"eu-central-1", "us-east-1"},
	}
	raw, err := json.Marshal(s)
	require.NoError(t, err, "marshal")

	out := &settings.Settings{}
	require.NoError(t, json.Unmarshal(raw, out), "unmarshal")

	require.NotNil(t, out.CloudWatchLogs, "cloud watch logs lost")
	assert.True(t, out.CloudWatchLogs.Enabled)
	assert.ElementsMatch(t, s.CloudWatchLogs.Regions, out.CloudWatchLogs.Regions)
}

func TestCustomNamespaceWithMetricRoundTrip(t *testing.T) {
	s := base()
	s.CustomNamespaces = settings.CustomNamespaces{
		{
			Namespace:            "AWS/GroundStation",
			AutoDiscoveryEnabled: false,
			Metrics: []*settings.CustomMetric{
				{
					Name:         "AzimuthAngle",
					Unit:         "Count",
					Dimensions:   []string{"SatelliteId"},
					Aggregations: []string{"Sum", "SampleCount"},
					Type:         "CUSTOM_AWS",
				},
			},
		},
		{
			Namespace:            "MyApp/Metrics",
			AutoDiscoveryEnabled: false,
			Metrics: []*settings.CustomMetric{
				{
					Name:         "queue.depth",
					Unit:         "Count",
					Aggregations: []string{"Average"},
					Type:         "CUSTOM",
				},
			},
		},
	}
	raw, err := json.Marshal(s)
	require.NoError(t, err, "marshal")

	out := &settings.Settings{}
	require.NoError(t, json.Unmarshal(raw, out), "unmarshal")

	require.Len(t, out.CustomNamespaces, 2)

	gs := out.CustomNamespaces[0]
	assert.Equal(t, "AWS/GroundStation", gs.Namespace)
	require.Len(t, gs.Metrics, 1)

	m := gs.Metrics[0]
	assert.Equal(t, "AzimuthAngle", m.Name)
	assert.Equal(t, "CUSTOM_AWS", m.Type)
	assert.ElementsMatch(t, []string{"Sum", "SampleCount"}, m.Aggregations)
	assert.ElementsMatch(t, []string{"SatelliteId"}, m.Dimensions)

	require.Len(t, out.CustomNamespaces[1].Metrics, 1)
	assert.Equal(t, "CUSTOM", out.CustomNamespaces[1].Metrics[0].Type)
}

func TestDtLabelEnrichmentRoundTrip(t *testing.T) {
	s := base()
	s.DTLabelEnrichments = settings.DTLabelEnrichments{
		{Label: "dt.security_context", Literal: "my-app"},
		{Label: "dt.cost.product", TagKey: "product"},
	}
	raw, err := json.Marshal(s)
	require.NoError(t, err, "marshal")

	out := &settings.Settings{}
	require.NoError(t, json.Unmarshal(raw, out), "unmarshal")

	// The wire representation is a JSON object, so the decoded order carries no
	// meaning — compare as a set.
	assert.ElementsMatch(t, s.DTLabelEnrichments, out.DTLabelEnrichments)

	aws := awsBlock(t, s)
	dtl, ok := aws["dtLabelsEnrichment"].(map[string]any)
	require.True(t, ok, "dtLabelsEnrichment missing: %v", aws)
	assert.Equal(t, map[string]any{"literal": "my-app"}, dtl["dt.security_context"])
	assert.Equal(t, map[string]any{"tagKey": "product"}, dtl["dt.cost.product"])
}
