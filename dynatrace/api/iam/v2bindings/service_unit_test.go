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

package v2bindings

import (
	"context"
	"net/http"
	"strings"
	"testing"

	bindings "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/iam/v2bindings/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings"
	coreapi "github.com/dynatrace/dynatrace-configuration-as-code-core/api"
	rest2 "github.com/dynatrace/dynatrace-configuration-as-code-core/api/rest"
)

const (
	testAccount     = "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee"
	testGroup       = "11111111-1111-1111-1111-111111111111"
	readableEnv     = "env12345"
	forbiddenEnv    = "env67890"
	policyWithLevel = "22222222-2222-2222-2222-222222222222"
	policyNoLevel   = "33333333-3333-3333-3333-333333333333"
	policyNoDetails = "55555555-5555-5555-5555-555555555555"
)

// fakeIAMClient answers the handful of endpoints the binding service calls.
// Any environment other than readableEnv responds 403, which is what happens
// when the configured credentials do not cover every environment of an account.
type fakeIAMClient struct {
	calls []string
}

func (me *fakeIAMClient) AccountID() string { return testAccount }

func (me *fakeIAMClient) GET(ctx context.Context, url string, options rest2.RequestOptions) (coreapi.Response, error) {
	me.calls = append(me.calls, url)

	switch {
	case strings.Contains(url, "/environments"):
		return body(`{"data": [{"id": "` + readableEnv + `"}, {"id": "` + forbiddenEnv + `"}]}`)

	case strings.Contains(url, "bindings/groups/"):
		// Three policies are bound: one resolvable and carrying boundaries, one
		// whose level cannot be read, and one the API lists without any details.
		return body(`{
			"policyUuids": ["` + policyWithLevel + `", "` + policyNoLevel + `", "` + policyNoDetails + `"],
			"bindingsDetails": [{
				"policyUuid": "` + policyWithLevel + `",
				"groups": ["` + testGroup + `"],
				"levelId": "` + readableEnv + `",
				"levelType": "environment",
				"boundaries": ["44444444-4444-4444-4444-444444444444"]
			}, {
				"policyUuid": "` + policyNoLevel + `",
				"groups": ["` + testGroup + `"],
				"levelId": "` + readableEnv + `",
				"levelType": "environment",
				"parameters": {"scope": "restricted"}
			}]
		}`)

	case strings.Contains(url, "/"+forbiddenEnv+"/"):
		return coreapi.Response{}, coreapi.APIError{StatusCode: http.StatusForbidden}

	case strings.HasSuffix(url, "/policies"):
		return body(`{"policies": []}`)

	case strings.Contains(url, "/policies/"+policyWithLevel):
		if strings.Contains(url, "/environment/"+readableEnv+"/") {
			return body(`{"uuid": "` + policyWithLevel + `", "name": "readable policy"}`)
		}
		return coreapi.Response{}, coreapi.APIError{StatusCode: http.StatusNotFound}

	case strings.Contains(url, "/policies/"+policyNoLevel):
		// Not visible at any level the credentials can read.
		return coreapi.Response{}, coreapi.APIError{StatusCode: http.StatusNotFound}
	}

	return coreapi.Response{}, coreapi.APIError{StatusCode: http.StatusNotFound}
}

func (me *fakeIAMClient) POST(ctx context.Context, url string, payload any, options rest2.RequestOptions) (coreapi.Response, error) {
	return coreapi.Response{}, nil
}

func (me *fakeIAMClient) PUT(ctx context.Context, url string, payload any, options rest2.RequestOptions) (coreapi.Response, error) {
	return coreapi.Response{}, nil
}

func (me *fakeIAMClient) DELETE(ctx context.Context, url string, options rest2.RequestOptions) (coreapi.Response, error) {
	return coreapi.Response{}, nil
}

func body(payload string) (coreapi.Response, error) {
	return coreapi.Response{StatusCode: http.StatusOK, Data: []byte(payload)}, nil
}

// TestGetKeepsBindingsWhenAnEnvironmentIsForbidden covers two cases that used to
// lose data: a policy the credentials cannot resolve to a level, and a bound
// policy the API reports without binding details.
func TestGetKeepsBindingsWhenAnEnvironmentIsForbidden(t *testing.T) {
	client := &fakeIAMClient{}
	service := &BindingServiceClient{client: client}

	binding := &bindings.PolicyBinding{}
	id := joinTestID(testGroup, "environment", readableEnv)
	if err := service.Get(context.Background(), id, binding); err != nil {
		t.Fatalf("Get returned an error although only one environment is forbidden: %v", err)
	}

	if len(binding.Policies) != 3 {
		t.Fatalf("expected all three bound policies, got %d", len(binding.Policies))
	}

	resolved, unresolved, noDetails := binding.Policies[0], binding.Policies[1], binding.Policies[2]

	if !strings.HasPrefix(resolved.ID, policyWithLevel) {
		t.Errorf("first policy should be the resolvable one, got %q", resolved.ID)
	}
	if len(resolved.Boundaries) != 1 {
		t.Errorf("boundaries of the first policy were lost: %v", resolved.Boundaries)
	}
	if unresolved.ID != policyNoLevel {
		t.Errorf("a policy without a resolvable level should fall back to its uuid, got %q", unresolved.ID)
	}
	if unresolved.Parameters["scope"] != "restricted" {
		t.Errorf("parameters of the unresolvable policy were lost: %v", unresolved.Parameters)
	}
	if noDetails.ID != policyNoDetails {
		t.Errorf("a bound policy without binding details was dropped, got %q", noDetails.ID)
	}
}

// TestGetKeepsTheConfiguredSpellingOfAnUnresolvablePolicy makes sure the fallback
// does not rewrite an ID the configuration already uses, which would show up as a
// permanent diff.
func TestGetKeepsTheConfiguredSpellingOfAnUnresolvablePolicy(t *testing.T) {
	service := &BindingServiceClient{client: &fakeIAMClient{}}
	configured := policyNoLevel + "#-#environment#-#" + forbiddenEnv

	ctx := context.WithValue(context.Background(), settings.ContextKeyStateConfig,
		&bindings.PolicyBinding{Policies: []*bindings.Policy{{ID: configured}}})

	binding := &bindings.PolicyBinding{}
	if err := service.Get(ctx, joinTestID(testGroup, "environment", readableEnv), binding); err != nil {
		t.Fatalf("Get returned an error: %v", err)
	}

	for _, policy := range binding.Policies {
		if strings.HasPrefix(policy.ID, policyNoLevel) && policy.ID != configured {
			t.Errorf("expected the configured ID %q to be kept, got %q", configured, policy.ID)
		}
	}
}

func joinTestID(group string, levelType string, levelID string) string {
	return group + "#-#" + levelType + "#-#" + levelID
}
