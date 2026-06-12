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
	"reflect"
	"sort"
	"testing"

	settings "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/extensions/dac/gcpmonitoring/settings"
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
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	var got map[string]any
	if err := json.Unmarshal(raw, &got); err != nil {
		t.Fatalf("decode: %v", err)
	}
	return got["value"].(map[string]any)["googleCloud"].(map[string]any)
}

// TestMarshalWireShape pins the on-the-wire JSON shape we send to
// /platform/extensions/v2/extensions/com.dynatrace.extension.da-gcp/monitoringConfigurations.
// Shape derived from dtctl pkg/resources/gcpmonitoringconfig.
func TestMarshalWireShape(t *testing.T) {
	s := base()

	raw, err := json.Marshal(s)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	var got map[string]any
	if err := json.Unmarshal(raw, &got); err != nil {
		t.Fatalf("re-decode: %v", err)
	}

	if got["scope"] != settings.DefaultScope {
		t.Fatalf("scope: got %v, want %s", got["scope"], settings.DefaultScope)
	}

	value, ok := got["value"].(map[string]any)
	if !ok {
		t.Fatalf("value: missing or wrong type: %T", got["value"])
	}
	wantTopLevel := map[string]any{
		"enabled":           true,
		"description":       "my-gcp-monitoring",
		"version":           "2.0.0",
		"activationContext": "DATA_ACQUISITION",
	}
	for k, v := range wantTopLevel {
		if !reflect.DeepEqual(value[k], v) {
			t.Errorf("value.%s: got %v, want %v", k, value[k], v)
		}
	}

	gc, _ := value["googleCloud"].(map[string]any)
	if gc == nil {
		t.Fatalf("googleCloud block missing")
	}
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
		if _, ok := gc[key]; !ok {
			t.Errorf("googleCloud.%s missing — wire shape requires the key", key)
		}
	}

	creds, _ := gc["credentials"].([]any)
	if len(creds) != 1 {
		t.Fatalf("credentials length: got %d, want 1", len(creds))
	}
	cred := creds[0].(map[string]any)
	wantCred := map[string]any{
		"connectionId":   "conn-objectid",
		"serviceAccount": "dynatrace-integration@example.iam.gserviceaccount.com",
		"enabled":        true,
	}
	for k, v := range wantCred {
		if !reflect.DeepEqual(cred[k], v) {
			t.Errorf("credentials[0].%s: got %v, want %v", k, cred[k], v)
		}
	}
	// description defaults to top-level name when not set
	if cred["description"] != "my-gcp-monitoring" {
		t.Errorf("credentials[0].description: got %v, want defaulted to name", cred["description"])
	}
	// GCP credentials must NOT carry a `type` field — there is only one auth mode.
	if _, ok := cred["type"]; ok {
		t.Errorf("credentials[0].type must be absent on GCP (has only one auth mode)")
	}

	loc, _ := gc["locationFiltering"].([]any)
	if len(loc) != 2 {
		t.Errorf("locationFiltering length: got %d, want 2", len(loc))
	}

	fs, _ := value["featureSets"].([]any)
	if len(fs) != 1 || fs[0] != "compute_engine_essential" {
		t.Errorf("featureSets: got %v", fs)
	}

	// smartscapeConfiguration must be the object {enabled: bool}, not a plain bool.
	sc, ok := gc["smartscapeConfiguration"].(map[string]any)
	if !ok {
		t.Fatalf("smartscapeConfiguration must be an object, got %T", gc["smartscapeConfiguration"])
	}
	if sc["enabled"] != true {
		t.Errorf("smartscapeConfiguration.enabled: got %v, want true", sc["enabled"])
	}

	// observabilityScopesEnabled defaults to false — must not be emitted at all
	// (omitempty semantics — saves wire bytes and matches dtctl).
	if _, ok := gc["observabilityScopesEnabled"]; ok {
		t.Errorf("observabilityScopesEnabled must be omitted when false")
	}
}

func TestApplyDefaults(t *testing.T) {
	s := &settings.Settings{
		Name: "x",
		Credentials: settings.Credentials{
			{ConnectionID: "c", ServiceAccount: "sa@x.iam.gserviceaccount.com", Enabled: true},
		},
	}
	raw, _ := json.Marshal(s)
	var top map[string]any
	_ = json.Unmarshal(raw, &top)
	if top["scope"] != settings.DefaultScope {
		t.Errorf("scope default: got %v, want %s", top["scope"], settings.DefaultScope)
	}
	if v := top["value"].(map[string]any)["activationContext"]; v != settings.DefaultActivationContext {
		t.Errorf("activationContext default: got %v, want %s", v, settings.DefaultActivationContext)
	}
	gc := top["value"].(map[string]any)["googleCloud"].(map[string]any)
	creds := gc["credentials"].([]any)
	if creds[0].(map[string]any)["description"] != "x" {
		t.Errorf("credential.description default: got %v, want 'x' (top-level name)", creds[0].(map[string]any)["description"])
	}
}

func TestRoundTrip(t *testing.T) {
	in := base()
	in.ProjectFilter = []string{"my-prod-project", "my-staging-project"}
	in.FolderFilter = []string{"folders/123"}
	in.TagEnrichment = []string{"tagKeys/owner"}
	in.LabelEnrichment = []string{"team"}
	raw, err := json.Marshal(in)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	out := &settings.Settings{}
	if err := json.Unmarshal(raw, out); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	sort.Strings(in.Regions)
	sort.Strings(out.Regions)
	if !reflect.DeepEqual(in.Regions, out.Regions) {
		t.Errorf("regions: got %v, want %v", out.Regions, in.Regions)
	}
	sort.Strings(in.ProjectFilter)
	sort.Strings(out.ProjectFilter)
	if !reflect.DeepEqual(in.ProjectFilter, out.ProjectFilter) {
		t.Errorf("projectFilter: got %v, want %v", out.ProjectFilter, in.ProjectFilter)
	}
	sort.Strings(in.FolderFilter)
	sort.Strings(out.FolderFilter)
	if !reflect.DeepEqual(in.FolderFilter, out.FolderFilter) {
		t.Errorf("folderFilter: got %v, want %v", out.FolderFilter, in.FolderFilter)
	}
	sort.Strings(in.TagEnrichment)
	sort.Strings(out.TagEnrichment)
	if !reflect.DeepEqual(in.TagEnrichment, out.TagEnrichment) {
		t.Errorf("tagEnrichment: got %v, want %v", out.TagEnrichment, in.TagEnrichment)
	}
	sort.Strings(in.LabelEnrichment)
	sort.Strings(out.LabelEnrichment)
	if !reflect.DeepEqual(in.LabelEnrichment, out.LabelEnrichment) {
		t.Errorf("labelEnrichment: got %v, want %v", out.LabelEnrichment, in.LabelEnrichment)
	}
	sort.Strings(in.FeatureSets)
	sort.Strings(out.FeatureSets)
	if !reflect.DeepEqual(in.FeatureSets, out.FeatureSets) {
		t.Errorf("featureSets: got %v, want %v", out.FeatureSets, in.FeatureSets)
	}
	if len(out.Credentials) != 1 || out.Credentials[0].ConnectionID != in.Credentials[0].ConnectionID {
		t.Errorf("credential lost: %+v", out.Credentials)
	}
	if out.Credentials[0].ServiceAccount != in.Credentials[0].ServiceAccount {
		t.Errorf("serviceAccount mismatch: got %v", out.Credentials[0].ServiceAccount)
	}
	if out.Name != in.Name {
		t.Errorf("name (description): got %v", out.Name)
	}
	if !out.SmartscapeEnabled {
		t.Errorf("smartscape: got %v, want true (round-tripped)", out.SmartscapeEnabled)
	}
}

// TestTagsVsLabelsSeparation guards spec §5: GCP tags (`tagKeys/…` resource-
// manager tags) and labels (per-resource key/value pairs) are two distinct
// filtering inputs that must NOT collide on the wire.
func TestTagsVsLabelsSeparation(t *testing.T) {
	s := base()
	s.TagFilters = settings.TagFilters{
		{Key: "tagKeys/env", Value: "tagValues/prod", Condition: "INCLUDE"},
	}
	s.LabelFilters = settings.TagFilters{
		{Key: "team", Value: "infra", Condition: "EXCLUDE"},
	}
	gc := googleCloudBlock(t, s)

	tags, _ := gc["tagFiltering"].([]any)
	labels, _ := gc["labelFiltering"].([]any)
	if len(tags) != 1 {
		t.Fatalf("tagFiltering length: got %d, want 1", len(tags))
	}
	if len(labels) != 1 {
		t.Fatalf("labelFiltering length: got %d, want 1", len(labels))
	}
	tag := tags[0].(map[string]any)
	if tag["key"] != "tagKeys/env" || tag["condition"] != "INCLUDE" {
		t.Errorf("tagFiltering[0] mismatch: %+v", tag)
	}
	label := labels[0].(map[string]any)
	if label["key"] != "team" || label["condition"] != "EXCLUDE" {
		t.Errorf("labelFiltering[0] mismatch: %+v", label)
	}

	// Round-trip
	raw, _ := json.Marshal(s)
	out := &settings.Settings{}
	if err := json.Unmarshal(raw, out); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if len(out.TagFilters) != 1 || out.TagFilters[0].Key != "tagKeys/env" {
		t.Errorf("tag filters lost: %+v", out.TagFilters)
	}
	if len(out.LabelFilters) != 1 || out.LabelFilters[0].Key != "team" {
		t.Errorf("label filters lost: %+v", out.LabelFilters)
	}
}

// TestAPIEchoArraysGuard validates spec §5 gotcha #2: server returns empty
// arrays (`featureSetConfiguration`, `resources` when nothing set, etc.) that
// must NOT surface as state — otherwise plan drift is eternal.
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
	if err := json.Unmarshal(raw, out); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	// All optional slices must come back as nil (not empty), so that
	// Terraform sees "unchanged" rather than "moved from null to []".
	if out.ProjectFilter != nil {
		t.Errorf("ProjectFilter: got %v, want nil", out.ProjectFilter)
	}
	if out.FolderFilter != nil {
		t.Errorf("FolderFilter: got %v, want nil", out.FolderFilter)
	}
	if out.TagFilters != nil {
		t.Errorf("TagFilters: got %v, want nil", out.TagFilters)
	}
	if out.LabelFilters != nil {
		t.Errorf("LabelFilters: got %v, want nil", out.LabelFilters)
	}
	if out.TagEnrichment != nil {
		t.Errorf("TagEnrichment: got %v, want nil", out.TagEnrichment)
	}
	if out.LabelEnrichment != nil {
		t.Errorf("LabelEnrichment: got %v, want nil", out.LabelEnrichment)
	}
	if out.ResourceAutodiscovery != nil {
		t.Errorf("ResourceAutodiscovery: got %v, want nil", out.ResourceAutodiscovery)
	}
	if out.FeatureSets != nil {
		t.Errorf("FeatureSets: got %v, want nil", out.FeatureSets)
	}
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
	if err := json.Unmarshal(raw, out); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if !out.SmartscapeEnabled {
		t.Errorf("smartscape missing in response: got %v, want forced true", out.SmartscapeEnabled)
	}

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
	if err := json.Unmarshal(raw, out); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if !out.SmartscapeEnabled {
		t.Errorf("smartscape explicit false echoed by API: got %v, want forced true (hidden-attribute drift guard)", out.SmartscapeEnabled)
	}

	// 3. Caller stuffs SmartscapeEnabled=false into the struct → wire payload still true.
	s := &settings.Settings{
		Name:              "x",
		ExtensionVersion:  "2.0.0",
		Credentials:       settings.Credentials{{ConnectionID: "c", ServiceAccount: "sa@x.iam.gserviceaccount.com", Enabled: true}},
		SmartscapeEnabled: false,
	}
	gc := googleCloudBlock(t, s)
	sc, ok := gc["smartscapeConfiguration"].(map[string]any)
	if !ok {
		t.Fatalf("smartscapeConfiguration must be an object, got %T", gc["smartscapeConfiguration"])
	}
	if sc["enabled"] != true {
		t.Errorf("smartscapeConfiguration.enabled: got %v, want true (hardcoded)", sc["enabled"])
	}
}

// TestObservabilityScopesEnabled exercises the omitempty boolean: true →
// emitted, false → omitted entirely.
func TestObservabilityScopesEnabled(t *testing.T) {
	s := base()
	s.ObservabilityScopesEnabled = true
	gc := googleCloudBlock(t, s)
	if gc["observabilityScopesEnabled"] != true {
		t.Errorf("observabilityScopesEnabled: got %v, want true", gc["observabilityScopesEnabled"])
	}

	// Round-trip true.
	raw, _ := json.Marshal(s)
	out := &settings.Settings{}
	if err := json.Unmarshal(raw, out); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if !out.ObservabilityScopesEnabled {
		t.Errorf("ObservabilityScopesEnabled: got %v, want true (round-tripped)", out.ObservabilityScopesEnabled)
	}
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
	res, _ := gc["resources"].([]any)
	if len(res) != 2 {
		t.Fatalf("resources length: got %d, want 2", len(res))
	}
	r0 := res[0].(map[string]any)
	if r0["resourceType"] != "compute.googleapis.com/Instance" || r0["autoDiscoveryEnabled"] != true {
		t.Errorf("resources[0] mismatch: %+v", r0)
	}
	exc, _ := r0["autodiscoveryExcludeMetricType"].([]any)
	if len(exc) != 1 || exc[0] != "compute.googleapis.com/instance/disk/read_bytes_count" {
		t.Errorf("resources[0].autodiscoveryExcludeMetricType: got %v", exc)
	}
	r1 := res[1].(map[string]any)
	if r1["autoDiscoveryEnabled"] != false {
		t.Errorf("resources[1].autoDiscoveryEnabled: got %v, want false", r1["autoDiscoveryEnabled"])
	}
	// When ExcludeMetricType is empty the key is omitted.
	if _, ok := r1["autodiscoveryExcludeMetricType"]; ok {
		t.Errorf("resources[1].autodiscoveryExcludeMetricType must be omitted when empty")
	}

	// Round-trip
	raw, _ := json.Marshal(s)
	out := &settings.Settings{}
	if err := json.Unmarshal(raw, out); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if len(out.ResourceAutodiscovery) != 2 {
		t.Fatalf("ResourceAutodiscovery lost: %+v", out.ResourceAutodiscovery)
	}
	if out.ResourceAutodiscovery[0].ResourceType != "compute.googleapis.com/Instance" {
		t.Errorf("ResourceAutodiscovery[0].ResourceType: got %v", out.ResourceAutodiscovery[0].ResourceType)
	}
	if len(out.ResourceAutodiscovery[0].ExcludeMetricType) != 1 {
		t.Errorf("ResourceAutodiscovery[0].ExcludeMetricType lost: %v", out.ResourceAutodiscovery[0].ExcludeMetricType)
	}
}
