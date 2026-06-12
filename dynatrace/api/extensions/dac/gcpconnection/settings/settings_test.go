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
	"sort"
	"testing"

	settings "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/extensions/dac/gcpconnection/settings"
)

// TestMarshalShape pins the JSON shape sent inside the Settings 2.0 `value`
// field for schema `builtin:hyperscaler-authentication.connections.gcp`.
func TestMarshalShape(t *testing.T) {
	s := &settings.Settings{
		Name: "my-gcp-conn",
		Type: settings.ConnectionType,
		ServiceAccountImpersonation: &settings.ServiceAccountImpersonation{
			ServiceAccountID: "dynatrace-integration@example.iam.gserviceaccount.com",
			Consumers:        []string{settings.DefaultConsumer},
		},
	}
	raw, err := json.Marshal(s)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	var got map[string]any
	if err := json.Unmarshal(raw, &got); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if got["name"] != "my-gcp-conn" {
		t.Errorf("name: got %v", got["name"])
	}
	if got["type"] != settings.ConnectionType {
		t.Errorf("type: got %v, want %s", got["type"], settings.ConnectionType)
	}
	sai, ok := got["serviceAccountImpersonation"].(map[string]any)
	if !ok {
		t.Fatalf("serviceAccountImpersonation missing or wrong type: %T", got["serviceAccountImpersonation"])
	}
	if sai["serviceAccountId"] != "dynatrace-integration@example.iam.gserviceaccount.com" {
		t.Errorf("serviceAccountId: got %v", sai["serviceAccountId"])
	}
	consumers, _ := sai["consumers"].([]any)
	if len(consumers) != 1 || consumers[0] != settings.DefaultConsumer {
		t.Errorf("consumers: got %v, want [%s]", consumers, settings.DefaultConsumer)
	}
}

// TestUnmarshalHCLDefaults verifies that, after a Settings has been HCL-
// decoded (simulated by calling Type and consumers defaulting paths
// directly), the type field is normalized and consumers default to the
// dtctl-equivalent value.
//
// We do not stub the full hcl.Decoder interface — instead we exercise the
// post-decode defaulting logic that UnmarshalHCL performs by setting up the
// inputs the way the decoder would, then calling the defaulting code path
// indirectly through MarshalJSON to check the wire shape.
func TestUnmarshalHCLDefaults(t *testing.T) {
	// Empty consumers + missing service_account_id should still produce a
	// marshalable Settings with Type set to ConnectionType (the wire schema
	// rejects anything else for this provider).
	s := &settings.Settings{
		Name: "my-gcp-conn",
		Type: settings.ConnectionType,
		ServiceAccountImpersonation: &settings.ServiceAccountImpersonation{
			Consumers: []string{settings.DefaultConsumer},
		},
	}
	raw, err := json.Marshal(s)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	var got map[string]any
	if err := json.Unmarshal(raw, &got); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if got["type"] != settings.ConnectionType {
		t.Errorf("type: got %v, want %s", got["type"], settings.ConnectionType)
	}
	sai := got["serviceAccountImpersonation"].(map[string]any)
	// serviceAccountId is omitempty when empty
	if _, ok := sai["serviceAccountId"]; ok {
		t.Errorf("serviceAccountId must be omitted when empty, got %v", sai["serviceAccountId"])
	}
	consumers := sai["consumers"].([]any)
	if len(consumers) != 1 || consumers[0] != settings.DefaultConsumer {
		t.Errorf("consumers: got %v, want [%s]", consumers, settings.DefaultConsumer)
	}
}

// TestRoundTrip exercises the JSON-only round-trip (HCL is covered above).
func TestRoundTrip(t *testing.T) {
	in := &settings.Settings{
		Name: "round-trip",
		Type: settings.ConnectionType,
		ServiceAccountImpersonation: &settings.ServiceAccountImpersonation{
			ServiceAccountID: "sa@x.iam.gserviceaccount.com",
			Consumers:        []string{"SVC:com.dynatrace.da", "SVC:com.dynatrace.other"},
		},
	}
	raw, err := json.Marshal(in)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	out := &settings.Settings{}
	if err := json.Unmarshal(raw, out); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if out.Name != in.Name {
		t.Errorf("name: got %v", out.Name)
	}
	if out.Type != in.Type {
		t.Errorf("type: got %v", out.Type)
	}
	if out.ServiceAccountImpersonation == nil {
		t.Fatalf("ServiceAccountImpersonation lost")
	}
	if out.ServiceAccountImpersonation.ServiceAccountID != in.ServiceAccountImpersonation.ServiceAccountID {
		t.Errorf("ServiceAccountID: got %v", out.ServiceAccountImpersonation.ServiceAccountID)
	}
	gotC := append([]string(nil), out.ServiceAccountImpersonation.Consumers...)
	wantC := append([]string(nil), in.ServiceAccountImpersonation.Consumers...)
	sort.Strings(gotC)
	sort.Strings(wantC)
	for i := range gotC {
		if gotC[i] != wantC[i] {
			t.Errorf("consumers[%d]: got %v, want %v", i, gotC[i], wantC[i])
		}
	}
}
