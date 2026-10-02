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

package v2bindings

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/iam/policies"
	bindings "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/iam/v2bindings/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest"
	testing2 "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/testing"
	"github.com/dynatrace/dynatrace-configuration-as-code-core/api"
	rest2 "github.com/dynatrace/dynatrace-configuration-as-code-core/api/rest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	testAccountID  = "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee"
	testGroupID    = "11111111-1111-1111-1111-111111111111"
	readableEnv    = "env12345"
	forbiddenEnv   = "env67890"
	envPolicyID    = "22222222-2222-2222-2222-222222222222"
	absentPolicyID = "33333333-3333-3333-3333-333333333333"
	boundaryID     = "44444444-4444-4444-4444-444444444444"
)

func jsonResponse(t *testing.T, v any) (api.Response, error) {
	data, err := json.Marshal(v)
	require.NoError(t, err)
	return api.Response{StatusCode: http.StatusOK, Data: data}, nil
}

// newMockClient answers the IAM endpoints the binding service uses. Every request for
// forbiddenEnv fails with 403, like an environment the credentials are not permitted to read.
func newMockClient(t *testing.T, boundPolicy string) *testing2.MockIAMClient {
	return &testing2.MockIAMClient{
		AccountIDValue: testAccountID,
		GETFunc: func(_ context.Context, url string, _ rest2.RequestOptions) (api.Response, error) {
			switch {
			case strings.Contains(url, "/environment/"+forbiddenEnv+"/"):
				return api.Response{}, api.APIError{StatusCode: http.StatusForbidden}

			case strings.HasSuffix(url, "/environments"):
				return jsonResponse(t, policies.ListEnvResponse{Data: []policies.DataStub{{ID: forbiddenEnv}, {ID: readableEnv}}})

			case strings.Contains(url, "/bindings/groups/"):
				return jsonResponse(t, GroupPolicyBindings{
					PolicyUuids: []string{boundPolicy},
					BindingsDetails: []BindingDetails{{
						PolicyUUID: boundPolicy,
						GroupUUIDs: []string{testGroupID},
						LevelId:    readableEnv,
						LevelType:  "environment",
						Parameters: map[string]string{"scope": "restricted"},
						Boundaries: []string{boundaryID},
					}},
				})

			case strings.HasSuffix(url, "/policies"):
				return jsonResponse(t, policies.ListPoliciesResponse{})

			case url == fmt.Sprintf("/iam/v1/repo/environment/%s/policies/%s", readableEnv, envPolicyID):
				return jsonResponse(t, policies.PolicyStub{UUID: envPolicyID, Name: "env policy"})
			}
			return api.Response{}, api.APIError{StatusCode: http.StatusNotFound}
		},
	}
}

func bindingID(levelType, levelID string) string {
	return fmt.Sprintf("%s#-#%s#-#%s", testGroupID, levelType, levelID)
}

func TestBindingServiceClient_Get(t *testing.T) {
	t.Run("Resolves a policy although another environment is forbidden", func(t *testing.T) {
		service := &BindingServiceClient{client: newMockClient(t, envPolicyID)}

		binding := &bindings.PolicyBinding{}
		require.NoError(t, service.Get(t.Context(), bindingID("environment", readableEnv), binding))

		assert.Equal(t, testGroupID, binding.GroupID)
		assert.Equal(t, readableEnv, binding.Environment)
		require.Len(t, binding.Policies, 1)
		assert.Equal(t, fmt.Sprintf("%s#-#environment#-#%s", envPolicyID, readableEnv), binding.Policies[0].ID)
		assert.Equal(t, map[string]string{"scope": "restricted"}, binding.Policies[0].Parameters)
		assert.Equal(t, []string{boundaryID}, binding.Policies[0].Boundaries)
	})

	t.Run("Returns the permission error instead of a not found error for an unresolvable policy", func(t *testing.T) {
		service := &BindingServiceClient{client: newMockClient(t, absentPolicyID)}

		err := service.Get(t.Context(), bindingID("environment", readableEnv), &bindings.PolicyBinding{})

		var apiErr api.APIError
		require.ErrorAs(t, err, &apiErr)
		assert.Equal(t, http.StatusForbidden, apiErr.StatusCode)
		assert.False(t, rest.IsNotFoundError(err))
	})
}

func TestBindingServiceClient_List(t *testing.T) {
	t.Run("Keeps bindings of readable environments when one environment is forbidden", func(t *testing.T) {
		mock := &testing2.MockIAMClient{
			AccountIDValue: testAccountID,
			GETFunc: func(_ context.Context, url string, _ rest2.RequestOptions) (api.Response, error) {
				switch {
				case strings.Contains(url, "/environment/"+forbiddenEnv+"/"):
					return api.Response{}, api.APIError{StatusCode: http.StatusForbidden}
				case strings.HasSuffix(url, "/environments"):
					return jsonResponse(t, policies.ListEnvResponse{Data: []policies.DataStub{{ID: forbiddenEnv}, {ID: readableEnv}}})
				case strings.HasSuffix(url, "/bindings"):
					return jsonResponse(t, ListPolicyBindingsResponse{PolicyBindings: []PolicyBindingStub{{Groups: []string{testGroupID}}}})
				}
				return api.Response{}, api.APIError{StatusCode: http.StatusNotFound}
			},
		}

		stubs, err := (&BindingServiceClient{client: mock}).List(t.Context())
		require.NoError(t, err)

		ids := []string{}
		for _, stub := range stubs {
			ids = append(ids, stub.ID)
		}
		assert.ElementsMatch(t, []string{bindingID("account", testAccountID), bindingID("environment", readableEnv)}, ids)
	})
}
