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

package provider_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/provider"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/provider/envutils"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/stretchr/testify/require"
)

func TestRefreshOrderedSettings(t *testing.T) {
	t.Setenv(envutils.DynatraceDisableOrderingSupport.Key, "false")
	const endpoint = "/api/v2/settings/objects"
	const resource = "dynatrace_log_storage"
	var pages atomic.Int32
	var objects atomic.Int32
	var moved, deleted, failNextPage atomic.Bool
	failNextPage.Store(true)
	api := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch {
		case r.Method == http.MethodDelete && r.URL.Path == endpoint+"/C":
			deleted.Store(true)
			w.WriteHeader(http.StatusNoContent)
		case r.Method == http.MethodGet && strings.HasPrefix(r.URL.Path, endpoint+"/"):
			objects.Add(1)
			fmt.Fprint(w, `{"scope":"environment","value":{"config-item-title":"test"}}`)
		case r.URL.Path == endpoint && r.URL.Query().Get("schemaIds") == "builtin:logmonitoring.log-storage-settings":
			pages.Add(1)
			if deleted.Load() {
				fmt.Fprint(w, `{"items":[{"objectId":"A"},{"objectId":"B"}]}`)
			} else if moved.Load() {
				fmt.Fprint(w, `{"items":[{"objectId":"C"},{"objectId":"A"},{"objectId":"B"}]}`)
			} else {
				fmt.Fprint(w, `{"items":[{"objectId":"A"},{"objectId":"B"}],"nextPageKey":"page2"}`)
			}
		case r.URL.Path == endpoint && r.URL.Query().Get("nextPageKey") == "page2":
			pages.Add(1)
			if failNextPage.Swap(false) {
				http.Error(w, "list failed on page two", http.StatusBadRequest)
				return
			}
			fmt.Fprint(w, `{"items":[{"objectId":"C"}]}`)
		default:
			http.Error(w, "unexpected request", http.StatusBadRequest)
		}
	}))
	defer api.Close()

	server := provider.NewServer()
	schemas, err := server.GetProviderSchema(t.Context(), &tfprotov5.GetProviderSchemaRequest{})
	require.NoError(t, err)
	resourceType := schemas.ResourceSchemas[resource].ValueType()
	providerConfig := rpcValue(t, schemas.Provider.ValueType(), fmt.Sprintf(`{"dt_env_url":%q,"dt_api_token":"fixture"}`, api.URL))
	configure := func() {
		resp, err := server.ConfigureProvider(t.Context(), &tfprotov5.ConfigureProviderRequest{Config: providerConfig})
		require.NoError(t, err)
		require.Empty(t, resp.Diagnostics)
	}
	read := func(t *testing.T, id string) *tfprotov5.ReadResourceResponse {
		t.Helper()
		state := rpcValue(t, resourceType, fmt.Sprintf(`{"id":%q}`, id))
		resp, err := server.ReadResource(t.Context(), &tfprotov5.ReadResourceRequest{TypeName: resource, CurrentState: state})
		require.NoError(t, err)
		return resp
	}
	checkPrevious := func(t *testing.T, id, expected string) {
		t.Helper()
		resp := read(t, id)
		require.Empty(t, resp.Diagnostics)
		state, err := resp.NewState.Unmarshal(resourceType)
		require.NoError(t, err)
		var attrs map[string]tftypes.Value
		require.NoError(t, state.As(&attrs))
		var previous string
		require.NoError(t, attrs["insert_after"].As(&previous))
		require.Equal(t, expected, previous)
	}

	configure()
	require.NotEmpty(t, read(t, "C").Diagnostics)
	t.Run("concurrent reads retry and share every page", func(t *testing.T) {
		for id, expected := range map[string]string{"A": "", "B": "A", "C": "B"} {
			t.Run(id, func(t *testing.T) {
				t.Parallel()
				checkPrevious(t, id, expected)
			})
		}
	})
	require.EqualValues(t, 4, pages.Load()) // Two pages in the failed scan, two in the successful one.
	require.EqualValues(t, 4, objects.Load())

	moved.Store(true)
	configure()
	checkPrevious(t, "A", "C")
	require.EqualValues(t, 5, pages.Load())

	null := rpcValue(t, resourceType, "null")
	resp, err := server.ApplyResourceChange(t.Context(), &tfprotov5.ApplyResourceChangeRequest{
		TypeName: resource, PriorState: rpcValue(t, resourceType, `{"id":"C"}`), PlannedState: null, Config: null,
	})
	require.NoError(t, err)
	require.Empty(t, resp.Diagnostics)
	require.True(t, deleted.Load())
	checkPrevious(t, "A", "")
	require.EqualValues(t, 6, pages.Load())
}

func rpcValue(t *testing.T, typ tftypes.Type, json string) *tfprotov5.DynamicValue {
	t.Helper()
	value, err := tftypes.ValueFromJSON([]byte(json), typ)
	require.NoError(t, err)
	wire, err := tfprotov5.NewDynamicValue(typ, value)
	require.NoError(t, err)
	return &wire
}
