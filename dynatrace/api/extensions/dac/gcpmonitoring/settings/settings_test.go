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

	settings "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/extensions/dac/gcpmonitoring/settings"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func base() *settings.Settings {
	return &settings.Settings{
		Name:             "my-gcp-monitoring",
		Enabled:          true,
		ExtensionVersion: "2.0.0",
		Credentials: settings.Credentials{
			{
				ConnectionID:   "conn-objectid",
				ServiceAccount: "dynatrace-integration@example.iam.gserviceaccount.com",
				Enabled:        true,
			},
		},
		Regions:           []string{"us-central1", "europe-west1"},
		FeatureSets:       []string{"compute_engine_essential"},
		SmartscapeEnabled: true,
	}
}

func googleCloudBlock(t *testing.T, s *settings.Settings) map[string]any {
	t.Helper()
	raw, err := json.Marshal(s)
	require.NoError(t, err, "marshal")
	var got map[string]any
	require.NoError(t, json.Unmarshal(raw, &got), "decode")
	value, ok := got["value"].(map[string]any)
	require.True(t, ok, "value block missing")
	gc, ok := value["googleCloud"].(map[string]any)
	require.True(t, ok, "googleCloud block missing")
	return gc
}

// TestMarshalWireShape pins the on-the-wire JSON shape we send to
// /platform/extensions/v2/extensions/com.dynatrace.extension.da-gcp/monitoringConfigurations.
// Shape pinned against the payload the monitoringConfigurations endpoint
// accepts for com.dynatrace.extension.da-gcp.
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
	assert.Equal(t, "my-gcp-monitoring", value["description"])
	assert.Equal(t, "2.0.0", value["version"])
	assert.Equal(t, "DATA_ACQUISITION", value["activationContext"])

	gc, ok := value["googleCloud"].(map[string]any)
	require.True(t, ok, "googleCloud block missing")

	// Wire shape requires these keys even when empty so the API does not
	// rewrite them with server-side defaults.
	for _, key := range []string{
		"credentials",
		"locationFiltering",
		"projectFiltering",
		"folderFiltering",
		"tagFiltering",
		"labelFiltering",
		"tagEnrichment",
		"labelEnrichment",
		"smartscapeConfiguration",
		"resources",
	} {
		assert.Contains(t, gc, key, "wire shape requires the key")
	}

	creds, ok := gc["credentials"].([]any)
	require.True(t, ok)
	require.Len(t, creds, 1)
	cred, ok := creds[0].(map[string]any)
	require.True(t, ok)
	assert.Equal(t, "conn-objectid", cred["connectionId"])
	assert.Equal(t, "dynatrace-integration@example.iam.gserviceaccount.com", cred["serviceAccount"])
	assert.Equal(t, true, cred["enabled"])
	assert.Equal(t, "my-gcp-monitoring", cred["description"], "defaults to the top-level name")
	assert.NotContains(t, cred, "type", "GCP has only one auth mode")

	assert.ElementsMatch(t, []any{"us-central1", "europe-west1"}, gc["locationFiltering"])
	assert.ElementsMatch(t, []any{"compute_engine_essential"}, value["featureSets"])

	// smartscapeConfiguration must be the object {enabled: bool}, not a plain bool.
	sc, ok := gc["smartscapeConfiguration"].(map[string]any)
	require.True(t, ok, "smartscapeConfiguration must be an object, got %T", gc["smartscapeConfiguration"])
	assert.Equal(t, true, sc["enabled"])

	// observabilityScopesEnabled defaults to false — must not be emitted at all
	// (omitempty semantics — keeps the payload minimal).
	assert.NotContains(t, gc, "observabilityScopesEnabled", "must be omitted when false")
}

func TestApplyDefaults(t *testing.T) {
	s := &settings.Settings{
		Name: "x",
		Credentials: settings.Credentials{
			{ConnectionID: "c", ServiceAccount: "sa@x.iam.gserviceaccount.com", Enabled: true},
		},
	}
	raw, err := json.Marshal(s)
	require.NoError(t, err, "marshal")

	var top map[string]any
	require.NoError(t, json.Unmarshal(raw, &top), "decode")
	assert.Equal(t, settings.DefaultScope, top["scope"])

	value, ok := top["value"].(map[string]any)
	require.True(t, ok, "value block missing")
	assert.Equal(t, settings.DefaultActivationContext, value["activationContext"])

	gc, ok := value["googleCloud"].(map[string]any)
	require.True(t, ok, "googleCloud block missing")
	creds, ok := gc["credentials"].([]any)
	require.True(t, ok)
	require.Len(t, creds, 1)
	cred, ok := creds[0].(map[string]any)
	require.True(t, ok)
	assert.Equal(t, "x", cred["description"], "defaults to the top-level name")
}

func TestRoundTrip(t *testing.T) {
	in := base()
	in.ProjectFilter = []string{"my-prod-project", "my-staging-project"}
	in.FolderFilter = []string{"folders/123"}
	in.TagEnrichment = []string{"tagKeys/owner"}
	in.LabelEnrichment = []string{"team"}
	raw, err := json.Marshal(in)
	require.NoError(t, err, "marshal")

	out := &settings.Settings{}
	require.NoError(t, json.Unmarshal(raw, out), "unmarshal")

	assert.ElementsMatch(t, in.Regions, out.Regions)
	assert.ElementsMatch(t, in.ProjectFilter, out.ProjectFilter)
	assert.ElementsMatch(t, in.FolderFilter, out.FolderFilter)
	assert.ElementsMatch(t, in.TagEnrichment, out.TagEnrichment)
	assert.ElementsMatch(t, in.LabelEnrichment, out.LabelEnrichment)
	assert.ElementsMatch(t, in.FeatureSets, out.FeatureSets)

	require.Len(t, out.Credentials, 1)
	assert.Equal(t, in.Credentials[0].ConnectionID, out.Credentials[0].ConnectionID)
	assert.Equal(t, in.Credentials[0].ServiceAccount, out.Credentials[0].ServiceAccount)
	assert.Equal(t, in.Name, out.Name, "name (description)")
	assert.True(t, out.SmartscapeEnabled, "round-tripped")
}

// TestTagsVsLabelsSeparation guards that GCP tags (`tagKeys/…` resource-manager
// tags) and labels (per-resource key/value pairs) are two distinct filtering
// inputs that must NOT collide on the wire.
func TestTagsVsLabelsSeparation(t *testing.T) {
	s := base()
	s.TagFilters = settings.TagFilters{
		{Key: "tagKeys/env", Value: "tagValues/prod", Condition: "INCLUDE"},
	}
	s.LabelFilters = settings.TagFilters{
		{Key: "team", Value: "infra", Condition: "EXCLUDE"},
	}
	gc := googleCloudBlock(t, s)

	tags, ok := gc["tagFiltering"].([]any)
	require.True(t, ok)
	require.Len(t, tags, 1)
	labels, ok := gc["labelFiltering"].([]any)
	require.True(t, ok)
	require.Len(t, labels, 1)

	tag, ok := tags[0].(map[string]any)
	require.True(t, ok)
	assert.Equal(t, "tagKeys/env", tag["key"])
	assert.Equal(t, "INCLUDE", tag["condition"])

	label, ok := labels[0].(map[string]any)
	require.True(t, ok)
	assert.Equal(t, "team", label["key"])
	assert.Equal(t, "EXCLUDE", label["condition"])

	raw, err := json.Marshal(s)
	require.NoError(t, err, "marshal")
	out := &settings.Settings{}
	require.NoError(t, json.Unmarshal(raw, out), "unmarshal")

	assert.ElementsMatch(t, s.TagFilters, out.TagFilters)
	assert.ElementsMatch(t, s.LabelFilters, out.LabelFilters)
}

// TestAPIEchoArraysGuard validates that the empty arrays the server echoes
// back (`featureSetConfiguration`, `resources` when nothing is set, and the
// filtering lists) do NOT surface as state — otherwise plan drift is eternal.
func TestAPIEchoArraysGuard(t *testing.T) {
	raw := []byte(`{
		"scope":"integration-gcp",
		"value":{
			"description":"x","enabled":true,"version":"2.0.0","activationContext":"DATA_ACQUISITION",
			"googleCloud":{
				"credentials":[{"connectionId":"c","serviceAccount":"sa@x.iam.gserviceaccount.com","enabled":true}],
				"locationFiltering":[],
				"projectFiltering":[],
				"folderFiltering":[],
				"tagFiltering":[],
				"labelFiltering":[],
				"tagEnrichment":[],
				"labelEnrichment":[],
				"resources":[],
				"smartscapeConfiguration":{"enabled":true}
			},
			"featureSetConfiguration":[]
		}
	}`)
	out := &settings.Settings{}
	require.NoError(t, json.Unmarshal(raw, out), "unmarshal")

	// All optional slices must come back as nil (not empty), so that
	// Terraform sees "unchanged" rather than "moved from null to []".
	assert.Nil(t, out.ProjectFilter)
	assert.Nil(t, out.FolderFilter)
	assert.Nil(t, out.TagFilters)
	assert.Nil(t, out.LabelFilters)
	assert.Nil(t, out.TagEnrichment)
	assert.Nil(t, out.LabelEnrichment)
	assert.Nil(t, out.ResourceAutodiscovery)
	assert.Nil(t, out.FeatureSets)
}

// TestSmartscapeAlwaysTrue guards the hidden-attribute contract: smartscape
// is intentionally not user-configurable. Whatever the API echoes, and
// whatever a caller stuffs into the struct, the wire payload must always
// re-send the canonical {enabled:true} so plans stay stable.
func TestSmartscapeAlwaysTrue(t *testing.T) {
	// 1. Response with smartscapeConfiguration missing → forced to true.
	raw := []byte(`{
		"scope":"integration-gcp",
		"value":{
			"description":"x","enabled":true,"version":"2.0.0","activationContext":"DATA_ACQUISITION",
			"googleCloud":{
				"credentials":[{"connectionId":"c","serviceAccount":"sa@x.iam.gserviceaccount.com","enabled":true}]
			}
		}
	}`)
	out := &settings.Settings{}
	require.NoError(t, json.Unmarshal(raw, out), "unmarshal")
	assert.True(t, out.SmartscapeEnabled, "smartscape missing in response must be forced true")

	// 2. Response with smartscapeConfiguration.enabled=false → still forced true.
	raw = []byte(`{
		"scope":"integration-gcp",
		"value":{
			"description":"x","enabled":true,"version":"2.0.0","activationContext":"DATA_ACQUISITION",
			"googleCloud":{
				"credentials":[{"connectionId":"c","serviceAccount":"sa@x.iam.gserviceaccount.com","enabled":true}],
				"smartscapeConfiguration":{"enabled":false}
			}
		}
	}`)
	out = &settings.Settings{}
	require.NoError(t, json.Unmarshal(raw, out), "unmarshal")
	assert.True(t, out.SmartscapeEnabled, "explicit false echoed by the API must be forced true")

	// 3. Caller stuffs SmartscapeEnabled=false into the struct → wire payload still true.
	s := &settings.Settings{
		Name:              "x",
		ExtensionVersion:  "2.0.0",
		Credentials:       settings.Credentials{{ConnectionID: "c", ServiceAccount: "sa@x.iam.gserviceaccount.com", Enabled: true}},
		SmartscapeEnabled: false,
	}
	gc := googleCloudBlock(t, s)
	sc, ok := gc["smartscapeConfiguration"].(map[string]any)
	require.True(t, ok, "smartscapeConfiguration must be an object, got %T", gc["smartscapeConfiguration"])
	assert.Equal(t, true, sc["enabled"], "hardcoded")
}

// TestObservabilityScopesEnabled exercises the omitempty boolean: true →
// emitted, false → omitted entirely.
func TestObservabilityScopesEnabled(t *testing.T) {
	s := base()
	s.ObservabilityScopesEnabled = true
	gc := googleCloudBlock(t, s)
	assert.Equal(t, true, gc["observabilityScopesEnabled"])

	raw, err := json.Marshal(s)
	require.NoError(t, err, "marshal")
	out := &settings.Settings{}
	require.NoError(t, json.Unmarshal(raw, out), "unmarshal")
	assert.True(t, out.ObservabilityScopesEnabled, "round-tripped")
}

// TestResourceAutodiscoveryRoundTrip exercises per-resource-type overrides
// (`resources[]` on the wire, `resource_autodiscovery` blocks in HCL).
func TestResourceAutodiscoveryRoundTrip(t *testing.T) {
	s := base()
	s.ResourceAutodiscovery = settings.ResourceAutodiscoveries{
		{
			ResourceType:         "compute.googleapis.com/Instance",
			AutoDiscoveryEnabled: true,
			ExcludeMetricType:    []string{"compute.googleapis.com/instance/disk/read_bytes_count"},
		},
		{
			ResourceType:         "storage.googleapis.com/Bucket",
			AutoDiscoveryEnabled: false,
		},
	}
	gc := googleCloudBlock(t, s)
	res, ok := gc["resources"].([]any)
	require.True(t, ok)
	require.Len(t, res, 2)

	r0, ok := res[0].(map[string]any)
	require.True(t, ok)
	assert.Equal(t, "compute.googleapis.com/Instance", r0["resourceType"])
	assert.Equal(t, true, r0["autoDiscoveryEnabled"])
	assert.ElementsMatch(t, []any{"compute.googleapis.com/instance/disk/read_bytes_count"}, r0["autodiscoveryExcludeMetricType"])

	r1, ok := res[1].(map[string]any)
	require.True(t, ok)
	assert.Equal(t, false, r1["autoDiscoveryEnabled"])
	assert.NotContains(t, r1, "autodiscoveryExcludeMetricType", "must be omitted when empty")

	raw, err := json.Marshal(s)
	require.NoError(t, err, "marshal")
	out := &settings.Settings{}
	require.NoError(t, json.Unmarshal(raw, out), "unmarshal")

	assert.ElementsMatch(t, s.ResourceAutodiscovery, out.ResourceAutodiscovery)
}
