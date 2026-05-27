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

	settings "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/extensions/dac/awsmonitoring/settings"
)

// TestMarshalWireShape pins the on-the-wire JSON shape we send to
// /api/v2/extensions/com.dynatrace.extension.da-aws/monitoringConfigurations.
// The shape was extracted from src-knowledge/dtctl/examples/aws_monitoring_config.yaml.
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
		"description":       "my-aws-monitoring",
		"version":           "1.0.0",
		"activationContext": "DATA_ACQUISITION",
	}
	for k, v := range wantTopLevel {
		if !reflect.DeepEqual(value[k], v) {
			t.Errorf("value.%s: got %v, want %v", k, value[k], v)
		}
	}

	fs, _ := value["featureSets"].([]any)
	if len(fs) != 2 {
		t.Errorf("featureSets length: got %d", len(fs))
	}

	aws, _ := value["aws"].(map[string]any)
	if aws == nil {
		t.Fatalf("aws block missing")
	}
	if aws["deploymentRegion"] != "us-east-1" {
		t.Errorf("deploymentRegion: got %v, want us-east-1 (defaulted from first region)", aws["deploymentRegion"])
	}
	creds, _ := aws["credentials"].([]any)
	if len(creds) != 1 {
		t.Fatalf("credentials length: got %d, want 1", len(creds))
	}
	cred := creds[0].(map[string]any)
	if cred["connectionId"] != s.ConnectionID {
		t.Errorf("connectionId mismatch")
	}
	if cred["accountId"] != "123456789012" {
		t.Errorf("accountId: got %v", cred["accountId"])
	}

	mc, _ := aws["metricsConfiguration"].(map[string]any)
	regs, _ := mc["regions"].([]any)
	if len(regs) != 2 {
		t.Errorf("metricsConfiguration.regions length: got %d", len(regs))
	}
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
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	out := &settings.Settings{}
	if err := json.Unmarshal(raw, out); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	sort.Strings(in.FeatureSets)
	sort.Strings(out.FeatureSets)
	if !reflect.DeepEqual(in.Regions, out.Regions) {
		t.Errorf("regions: got %v, want %v", out.Regions, in.Regions)
	}
	if !reflect.DeepEqual(in.FeatureSets, out.FeatureSets) {
		t.Errorf("featureSets: got %v, want %v", out.FeatureSets, in.FeatureSets)
	}
	if out.ConnectionID != in.ConnectionID {
		t.Errorf("connectionId: got %v", out.ConnectionID)
	}
	if out.AccountID != in.AccountID {
		t.Errorf("accountId: got %v", out.AccountID)
	}
	if out.Name != in.Name {
		t.Errorf("name (description): got %v", out.Name)
	}
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
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	var got map[string]any
	if err := json.Unmarshal(raw, &got); err != nil {
		t.Fatalf("decode: %v", err)
	}
	return got["value"].(map[string]any)["aws"].(map[string]any)
}

func TestEnumDefaults(t *testing.T) {
	s := base()
	aws := awsBlock(t, s)
	for k, want := range map[string]any{
		"deploymentScope":   "SINGLE_ACCOUNT",
		"deploymentMode":    "AUTOMATED",
		"configurationMode": "QUICK_START",
	} {
		if aws[k] != want {
			t.Errorf("aws.%s: got %v, want %v", k, aws[k], want)
		}
	}
	if sm, ok := aws["smartscapeConfiguration"].(map[string]any); !ok || sm["enabled"] != false {
		// SmartscapeEnabled zero value is false on a fresh struct; schema default
		// of true only fires through Terraform. Just assert it's a bool.
		if _, ok := sm["enabled"].(bool); !ok {
			t.Errorf("smartscapeConfiguration.enabled missing/bad type: %v", sm)
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

func TestCloudWatchLogsRoundTrip(t *testing.T) {
	s := base()
	s.CloudWatchLogs = &settings.CloudWatchLogsConfig{
		Enabled: true,
		Regions: []string{"eu-central-1", "us-east-1"},
	}
	raw, _ := json.Marshal(s)
	out := &settings.Settings{}
	if err := json.Unmarshal(raw, out); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if out.CloudWatchLogs == nil {
		t.Fatalf("cloud watch logs lost")
	}
	if !out.CloudWatchLogs.Enabled || len(out.CloudWatchLogs.Regions) != 2 {
		t.Errorf("cwl: %+v", out.CloudWatchLogs)
	}
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
	raw, _ := json.Marshal(s)
	out := &settings.Settings{}
	if err := json.Unmarshal(raw, out); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if len(out.CustomNamespaces) != 2 {
		t.Fatalf("namespaces: got %d, want 2", len(out.CustomNamespaces))
	}
	gs := out.CustomNamespaces[0]
	if gs.Namespace != "AWS/GroundStation" || len(gs.Metrics) != 1 {
		t.Fatalf("ground station: %+v", gs)
	}
	m := gs.Metrics[0]
	if m.Name != "AzimuthAngle" || m.Type != "CUSTOM_AWS" || len(m.Aggregations) != 2 {
		t.Errorf("metric: %+v", m)
	}
	if out.CustomNamespaces[1].Metrics[0].Type != "CUSTOM" {
		t.Errorf("custom namespace type mismatch")
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

	// Wire shape check
	aws := awsBlock(t, s)
	dtl, ok := aws["dtLabelsEnrichment"].(map[string]any)
	if !ok {
		t.Fatalf("dtLabelsEnrichment missing: %v", aws)
	}
	if !reflect.DeepEqual(dtl["dt.security_context"], map[string]any{"literal": "my-app"}) {
		t.Errorf("literal entry wrong: %v", dtl["dt.security_context"])
	}
	if !reflect.DeepEqual(dtl["dt.cost.product"], map[string]any{"tagKey": "product"}) {
		t.Errorf("tagKey entry wrong: %v", dtl["dt.cost.product"])
	}
}
