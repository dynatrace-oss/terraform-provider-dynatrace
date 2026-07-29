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

	settings "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/extensions/dac/azuremonitoring/settings"
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
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	var got map[string]any
	if err := json.Unmarshal(raw, &got); err != nil {
		t.Fatalf("decode: %v", err)
	}
	return got["value"].(map[string]any)["azure"].(map[string]any)
}

// TestMarshalWireShape pins the on-the-wire JSON shape we send to
// /platform/extensions/v2/extensions/com.dynatrace.extension.da-azure/monitoringConfigurations.
// Shape derived from dtctl pkg/resources/azuremonitoringconfig.
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
		"description":       "my-azure-monitoring",
		"version":           "2.0.0",
		"activationContext": "DATA_ACQUISITION",
	}
	for k, v := range wantTopLevel {
		if !reflect.DeepEqual(value[k], v) {
			t.Errorf("value.%s: got %v, want %v", k, value[k], v)
		}
	}

	azure, _ := value["azure"].(map[string]any)
	if azure == nil {
		t.Fatalf("azure block missing")
	}
	wantAzure := map[string]any{
		"subscriptionFilteringMode": "INCLUDE",
		"configurationMode":         "ADVANCED",
		"deploymentMode":            "AUTOMATED",
		"deploymentScope":           "SUBSCRIPTION",
	}
	for k, v := range wantAzure {
		if azure[k] != v {
			t.Errorf("azure.%s: got %v, want %v", k, azure[k], v)
		}
	}

	for _, key := range []string{"credentials", "locationFiltering", "subscriptionFiltering", "tagFiltering", "tagEnrichment"} {
		if _, ok := azure[key]; !ok {
			t.Errorf("azure.%s missing — wire shape requires the key even when empty", key)
		}
	}

	creds, _ := azure["credentials"].([]any)
	if len(creds) != 1 {
		t.Fatalf("credentials length: got %d, want 1", len(creds))
	}
	cred := creds[0].(map[string]any)
	wantCred := map[string]any{
		"connectionId":       "conn-objectid",
		"servicePrincipalId": "00000000-0000-0000-0000-000000000001",
		"type":               "FEDERATED",
		"enabled":            true,
	}
	for k, v := range wantCred {
		if !reflect.DeepEqual(cred[k], v) {
			t.Errorf("credentials[0].%s: got %v, want %v", k, cred[k], v)
		}
	}
	// description defaults to top-level name when not set
	if cred["description"] != "my-azure-monitoring" {
		t.Errorf("credentials[0].description: got %v, want defaulted to name", cred["description"])
	}

	loc, _ := azure["locationFiltering"].([]any)
	if len(loc) != 2 {
		t.Errorf("locationFiltering length: got %d, want 2", len(loc))
	}

	fs, _ := value["featureSets"].([]any)
	if len(fs) != 1 || fs[0] != "microsoft_compute.virtualmachines_essential" {
		t.Errorf("featureSets: got %v", fs)
	}
}

func TestEnumDefaults(t *testing.T) {
	s := &settings.Settings{
		Name: "x",
		Credentials: settings.Credentials{
			{ConnectionID: "c", ServicePrincipalID: "spid", Enabled: true},
		},
	}
	azure := azureBlock(t, s)
	for k, want := range map[string]any{
		"configurationMode":         "ADVANCED",
		"deploymentMode":            "AUTOMATED",
		"deploymentScope":           "SUBSCRIPTION",
		"subscriptionFilteringMode": "INCLUDE",
	} {
		if azure[k] != want {
			t.Errorf("azure.%s: got %v, want %v", k, azure[k], want)
		}
	}

	raw, _ := json.Marshal(s)
	var top map[string]any
	_ = json.Unmarshal(raw, &top)
	if top["scope"] != settings.DefaultScope {
		t.Errorf("scope: got %v, want %s", top["scope"], settings.DefaultScope)
	}

	// Credential gets defaulted type FEDERATED.
	creds := azure["credentials"].([]any)
	if creds[0].(map[string]any)["type"] != "FEDERATED" {
		t.Errorf("credential.type default: got %v, want FEDERATED", creds[0].(map[string]any)["type"])
	}
}

func TestRoundTrip(t *testing.T) {
	in := base()
	in.SubscriptionFilter = []string{"00000000-0000-0000-0000-000000000abc"}
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
	sort.Strings(in.FeatureSets)
	sort.Strings(out.FeatureSets)
	if !reflect.DeepEqual(in.FeatureSets, out.FeatureSets) {
		t.Errorf("featureSets: got %v, want %v", out.FeatureSets, in.FeatureSets)
	}
	if len(out.Credentials) != 1 || out.Credentials[0].ConnectionID != in.Credentials[0].ConnectionID {
		t.Errorf("credential lost: %+v", out.Credentials)
	}
	if out.Credentials[0].ServicePrincipalID != in.Credentials[0].ServicePrincipalID {
		t.Errorf("servicePrincipalId mismatch")
	}
	if out.Credentials[0].Type != "FEDERATED" {
		t.Errorf("credential.type: got %v", out.Credentials[0].Type)
	}
	if len(out.SubscriptionFilter) != 1 || out.SubscriptionFilter[0] != in.SubscriptionFilter[0] {
		t.Errorf("subscriptionFilter lost: %+v", out.SubscriptionFilter)
	}
	if out.Name != in.Name {
		t.Errorf("name (description): got %v", out.Name)
	}
}

func TestCredentialTypeDefaultedOnUnmarshal(t *testing.T) {
	// Older configs in the live dump omit `type` on credentials. Make sure
	// UnmarshalJSON injects FEDERATED so set comparison stays stable.
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
	if err := json.Unmarshal(raw, out); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if len(out.Credentials) != 1 {
		t.Fatalf("credentials lost: %+v", out.Credentials)
	}
	if out.Credentials[0].Type != "FEDERATED" {
		t.Errorf("credential.type: got %q, want FEDERATED", out.Credentials[0].Type)
	}
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
	if err := json.Unmarshal(raw, out); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	// Round-trip → the rendered payload must not carry those keys back.
	rendered, err := json.Marshal(out)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	var got map[string]any
	_ = json.Unmarshal(rendered, &got)
	azure := got["value"].(map[string]any)["azure"].(map[string]any)
	for _, forbidden := range []string{"namespaces", "eventHubsConfiguration"} {
		if _, ok := azure[forbidden]; ok {
			t.Errorf("azure.%s leaked back into the wire payload (eternal-drift trap)", forbidden)
		}
	}
}

func TestTagFilterRoundTrip(t *testing.T) {
	s := base()
	s.TagFilters = settings.TagFilters{
		{Key: "env", Value: "prod", Condition: "INCLUDE"},
		{Key: "team", Value: "infra", Condition: "EXCLUDE"},
	}
	raw, _ := json.Marshal(s)
	out := &settings.Settings{}
	if err := json.Unmarshal(raw, out); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if len(out.TagFilters) != 2 {
		t.Fatalf("tag filters: got %d, want 2", len(out.TagFilters))
	}
	if out.TagFilters[0].Key != "env" || out.TagFilters[0].Condition != "INCLUDE" {
		t.Errorf("tag filter[0] mismatch: %+v", out.TagFilters[0])
	}
	if out.TagFilters[1].Condition != "EXCLUDE" {
		t.Errorf("tag filter[1] condition mismatch: %+v", out.TagFilters[1])
	}
}

func TestTagEnrichmentRoundTrip(t *testing.T) {
	s := base()
	s.TagEnrichment = []string{"owner", "cost-center"}
	raw, _ := json.Marshal(s)
	out := &settings.Settings{}
	if err := json.Unmarshal(raw, out); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	sort.Strings(out.TagEnrichment)
	want := []string{"cost-center", "owner"}
	if !reflect.DeepEqual(out.TagEnrichment, want) {
		t.Errorf("tag enrichment: got %v, want %v", out.TagEnrichment, want)
	}
}

func TestDtLabelEnrichmentRoundTrip(t *testing.T) {
	s := base()
	s.DTLabelEnrichments = settings.DTLabelEnrichments{
		{Label: "dt.security_context", Literal: "my-app"},
		{Label: "dt.cost.product", TagKey: "product"},
	}
	raw, _ := json.Marshal(s)
	out := &settings.Settings{}
	if err := json.Unmarshal(raw, out); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if len(out.DTLabelEnrichments) != 2 {
		t.Fatalf("labels: %d", len(out.DTLabelEnrichments))
	}
	// UnmarshalJSON sorts label keys alphabetically.
	if out.DTLabelEnrichments[0].Label != "dt.cost.product" || out.DTLabelEnrichments[0].TagKey != "product" {
		t.Errorf("label[0]: %+v", out.DTLabelEnrichments[0])
	}
	if out.DTLabelEnrichments[1].Label != "dt.security_context" || out.DTLabelEnrichments[1].Literal != "my-app" {
		t.Errorf("label[1]: %+v", out.DTLabelEnrichments[1])
	}

	azure := azureBlock(t, s)
	dtl, ok := azure["dtLabelsEnrichment"].(map[string]any)
	if !ok {
		t.Fatalf("dtLabelsEnrichment missing: %v", azure)
	}
	if !reflect.DeepEqual(dtl["dt.security_context"], map[string]any{"literal": "my-app"}) {
		t.Errorf("literal entry wrong: %v", dtl["dt.security_context"])
	}
	if !reflect.DeepEqual(dtl["dt.cost.product"], map[string]any{"tagKey": "product"}) {
		t.Errorf("tagKey entry wrong: %v", dtl["dt.cost.product"])
	}
}
