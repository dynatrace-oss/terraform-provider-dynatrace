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

	settings "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/extensions/dac/azuremonitoring/settings"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func base() *settings.Settings {
	return &settings.Settings{
		Name:             "my-azure-monitoring",
		Enabled:          true,
		ExtensionVersion: "2.0.0",
		Credentials: settings.Credentials{
			{
				ConnectionID:       "conn-objectid",
				ServicePrincipalID: "00000000-0000-0000-0000-000000000001",
				Type:               "FEDERATED",
				Enabled:            true,
			},
		},
		Regions:     []string{"eastus", "westeurope"},
		FeatureSets: []string{"microsoft_compute.virtualmachines_essential"},
	}
}

func azureBlock(t *testing.T, s *settings.Settings) map[string]any {
	t.Helper()
	raw, err := json.Marshal(s)
	require.NoError(t, err, "marshal")
	var got map[string]any
	require.NoError(t, json.Unmarshal(raw, &got), "decode")
	value, ok := got["value"].(map[string]any)
	require.True(t, ok, "value block missing")
	azure, ok := value["azure"].(map[string]any)
	require.True(t, ok, "azure block missing")
	return azure
}

// TestMarshalWireShape pins the on-the-wire JSON shape we send to
// /platform/extensions/v2/extensions/com.dynatrace.extension.da-azure/monitoringConfigurations.
// Shape pinned against the payload the monitoringConfigurations endpoint
// accepts for com.dynatrace.extension.da-azure.
func TestMarshalWireShape(t *testing.T) {
	s := base()

	raw, err := json.Marshal(s)
	require.NoError(t, err, "marshal")

	var got map[string]any
	require.NoError(t, json.Unmarshal(raw, &got), "re-decode")

	assert.Equal(t, settings.DefaultScope, got["scope"])

	value, ok := got["value"].(map[string]any)
	require.True(t, ok, "value: missing or wrong type: %T", got["value"])

	assert.Equal(t, true, value["enabled"])
	assert.Equal(t, "my-azure-monitoring", value["description"])
	assert.Equal(t, "2.0.0", value["version"])
	assert.Equal(t, "DATA_ACQUISITION", value["activationContext"])

	azure, ok := value["azure"].(map[string]any)
	require.True(t, ok, "azure block missing")

	assert.Equal(t, "INCLUDE", azure["subscriptionFilteringMode"])
	assert.Equal(t, "ADVANCED", azure["configurationMode"])
	assert.Equal(t, "AUTOMATED", azure["deploymentMode"])
	assert.Equal(t, "SUBSCRIPTION", azure["deploymentScope"])

	for _, key := range []string{"credentials", "locationFiltering", "subscriptionFiltering", "tagFiltering", "tagEnrichment"} {
		assert.Contains(t, azure, key, "wire shape requires the key even when empty")
	}

	creds, ok := azure["credentials"].([]any)
	require.True(t, ok)
	require.Len(t, creds, 1)
	cred, ok := creds[0].(map[string]any)
	require.True(t, ok)
	assert.Equal(t, "conn-objectid", cred["connectionId"])
	assert.Equal(t, "00000000-0000-0000-0000-000000000001", cred["servicePrincipalId"])
	assert.Equal(t, "FEDERATED", cred["type"])
	assert.Equal(t, true, cred["enabled"])
	assert.Equal(t, "my-azure-monitoring", cred["description"], "defaults to the top-level name")

	assert.ElementsMatch(t, []any{"eastus", "westeurope"}, azure["locationFiltering"])
	assert.ElementsMatch(t, []any{"microsoft_compute.virtualmachines_essential"}, value["featureSets"])
}

func TestEnumDefaults(t *testing.T) {
	s := &settings.Settings{
		Name: "x",
		Credentials: settings.Credentials{
			{ConnectionID: "c", ServicePrincipalID: "spid", Enabled: true},
		},
	}
	azure := azureBlock(t, s)

	assert.Equal(t, "ADVANCED", azure["configurationMode"])
	assert.Equal(t, "AUTOMATED", azure["deploymentMode"])
	assert.Equal(t, "SUBSCRIPTION", azure["deploymentScope"])
	assert.Equal(t, "INCLUDE", azure["subscriptionFilteringMode"])

	raw, err := json.Marshal(s)
	require.NoError(t, err, "marshal")
	var top map[string]any
	require.NoError(t, json.Unmarshal(raw, &top), "decode")
	assert.Equal(t, settings.DefaultScope, top["scope"])

	// Credential gets defaulted type FEDERATED.
	creds, ok := azure["credentials"].([]any)
	require.True(t, ok)
	require.Len(t, creds, 1)
	cred, ok := creds[0].(map[string]any)
	require.True(t, ok)
	assert.Equal(t, "FEDERATED", cred["type"])
}

func TestRoundTrip(t *testing.T) {
	in := base()
	in.SubscriptionFilter = []string{"00000000-0000-0000-0000-000000000abc"}
	raw, err := json.Marshal(in)
	require.NoError(t, err, "marshal")

	out := &settings.Settings{}
	require.NoError(t, json.Unmarshal(raw, out), "unmarshal")

	assert.ElementsMatch(t, in.Regions, out.Regions)
	assert.ElementsMatch(t, in.FeatureSets, out.FeatureSets)
	assert.ElementsMatch(t, in.SubscriptionFilter, out.SubscriptionFilter)

	require.Len(t, out.Credentials, 1)
	assert.Equal(t, in.Credentials[0].ConnectionID, out.Credentials[0].ConnectionID)
	assert.Equal(t, in.Credentials[0].ServicePrincipalID, out.Credentials[0].ServicePrincipalID)
	assert.Equal(t, "FEDERATED", out.Credentials[0].Type)
	assert.Equal(t, in.Name, out.Name, "name (description)")
}

func TestCredentialTypeDefaultedOnUnmarshal(t *testing.T) {
	// Older configurations omit `type` on credentials. Make sure UnmarshalJSON
	// injects FEDERATED so set comparison stays stable.
	raw := []byte(`{
		"scope":"integration-azure",
		"value":{
			"description":"x","enabled":true,"version":"2.0.0","activationContext":"DATA_ACQUISITION",
			"azure":{
				"credentials":[{"connectionId":"c","servicePrincipalId":"sp","enabled":true}],
				"deploymentScope":"SUBSCRIPTION","configurationMode":"ADVANCED","deploymentMode":"AUTOMATED","subscriptionFilteringMode":"INCLUDE",
				"locationFiltering":[],"subscriptionFiltering":[],"tagFiltering":[],"tagEnrichment":[]
			}
		}
	}`)
	out := &settings.Settings{}
	require.NoError(t, json.Unmarshal(raw, out), "unmarshal")

	require.Len(t, out.Credentials, 1)
	assert.Equal(t, "FEDERATED", out.Credentials[0].Type)
}

func TestAPIEchoArraysIgnored(t *testing.T) {
	// `namespaces` and `eventHubsConfiguration` are API-echo arrays — the
	// server always returns them as []; surfacing them would create eternal
	// plan drift. Settings must not model them.
	raw := []byte(`{
		"scope":"integration-azure",
		"value":{
			"description":"x","enabled":true,"version":"2.0.0","activationContext":"DATA_ACQUISITION",
			"azure":{
				"credentials":[{"connectionId":"c","servicePrincipalId":"sp","type":"FEDERATED","enabled":true}],
				"namespaces":[],
				"eventHubsConfiguration":[],
				"deploymentScope":"SUBSCRIPTION","configurationMode":"ADVANCED","deploymentMode":"AUTOMATED","subscriptionFilteringMode":"INCLUDE"
			}
		}
	}`)
	out := &settings.Settings{}
	require.NoError(t, json.Unmarshal(raw, out), "unmarshal")

	// Round-trip → the rendered payload must not carry those keys back.
	azure := azureBlock(t, out)
	for _, forbidden := range []string{"namespaces", "eventHubsConfiguration"} {
		assert.NotContains(t, azure, forbidden, "leaked back into the wire payload (eternal-drift trap)")
	}
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

	azure := azureBlock(t, s)
	dtl, ok := azure["dtLabelsEnrichment"].(map[string]any)
	require.True(t, ok, "dtLabelsEnrichment missing: %v", azure)
	assert.Equal(t, map[string]any{"literal": "my-app"}, dtl["dt.security_context"])
	assert.Equal(t, map[string]any{"tagKey": "product"}, dtl["dt.cost.product"])
}
