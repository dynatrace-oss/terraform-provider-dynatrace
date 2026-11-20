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
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"slices"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/provider"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/provider/envutils"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	logStorageResource = "dynatrace_log_storage"
	logStorageSchema   = "builtin:logmonitoring.log-storage-settings"
	settingsEndpoint   = "/api/v2/settings/objects"
)

func TestRefreshOrder(t *testing.T) {
	t.Setenv(envutils.DynatraceDisableOrderingSupport.Key, "false")
	types := loadSchemaTypes(t)

	t.Run("reads in one refresh share a single scan", func(t *testing.T) {
		h := newHarness(t, types, "A", "B", "C")

		// Parallel reads share one scan; scheduling decides whether any of them actually waits on it.
		t.Run("parallel reads", func(t *testing.T) {
			for id, previous := range map[string]string{"A": "", "B": "A", "C": "B"} {
				t.Run(id, func(t *testing.T) {
					t.Parallel()
					assert.Equal(t, previous, h.insertAfter(t, id))
				})
			}
		})

		assert.EqualValues(t, 1, h.api.scans.Load(), "the order is listed once, not once per read")
	})

	t.Run("a failed scan is not cached", func(t *testing.T) {
		h := newHarness(t, types, "A", "B", "C")
		// The scan fails on page two, after page one already returned A and B.
		h.api.failNextPage.Store(true)

		require.NotEmpty(t, h.read(t, "C").Diagnostics)
		require.False(t, h.api.failNextPage.Load(), "the read failed because page two failed")
		require.EqualValues(t, 1, h.api.scans.Load())

		assert.Equal(t, "B", h.insertAfter(t, "C"), "the next read scans again instead of reusing the error or the partial page")
		assert.EqualValues(t, 2, h.api.scans.Load())
	})

	t.Run("reconfigure starts a new cache", func(t *testing.T) {
		h := newHarness(t, types, "A", "B", "C")
		require.Equal(t, "", h.insertAfter(t, "A"))

		h.api.setOrder("C", "A", "B")
		assert.Equal(t, "", h.insertAfter(t, "A"), "the same refresh keeps using the cached order")

		h.configure(t)
		assert.Equal(t, "C", h.insertAfter(t, "A"), "a new provider configuration lists again")
	})

	t.Run("apply disables the cache", func(t *testing.T) {
		h := newHarness(t, types, "C", "A", "B")
		require.Equal(t, "C", h.insertAfter(t, "A"))

		h.delete(t, "C")

		assert.Equal(t, "", h.insertAfter(t, "A"), "after a write, reads see the order without the deleted object")
		assert.EqualValues(t, 2, h.api.scans.Load())

		assert.Equal(t, "A", h.insertAfter(t, "B"))
		assert.EqualValues(t, 3, h.api.scans.Load(), "without a cache, every read lists again")
	})
}

// fakeSettingsAPI serves the Settings 2.0 endpoints for one ordered schema, paging over objectIDs.
// It ignores the client's pageSize and uses its own to force multi-page scans.
type fakeSettingsAPI struct {
	mu           sync.Mutex
	objectIDs    []string
	pageSize     int
	failNextPage atomic.Bool  // fails the next follow-up page request (nextPageKey set)
	scans        atomic.Int32 // number of first-page list requests
}

type fakeListItem struct {
	ObjectID string `json:"objectId"`
}

type fakeListResponse struct {
	Items       []fakeListItem `json:"items"`
	NextPageKey *string        `json:"nextPageKey,omitempty"`
}

func (f *fakeSettingsAPI) setOrder(ids ...string) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.objectIDs = slices.Clone(ids)
}

func (f *fakeSettingsAPI) remove(id string) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.objectIDs = slices.DeleteFunc(f.objectIDs, func(o string) bool { return o == id })
}

func (f *fakeSettingsAPI) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	objectID, isObject := strings.CutPrefix(r.URL.Path, settingsEndpoint+"/")
	switch {
	case r.URL.Path == settingsEndpoint && r.Method == http.MethodGet:
		f.list(w, r)
	case isObject && r.Method == http.MethodGet:
		fmt.Fprint(w, `{"scope":"environment","value":{"config-item-title":"test"}}`)
	case isObject && r.Method == http.MethodDelete:
		f.remove(objectID)
		w.WriteHeader(http.StatusNoContent)
	default:
		http.Error(w, "unexpected request", http.StatusBadRequest)
	}
}

func (f *fakeSettingsAPI) list(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	offset := 0
	if key := query.Get("nextPageKey"); key != "" {
		if f.failNextPage.Swap(false) {
			http.Error(w, "list failed on a follow-up page", http.StatusBadRequest)
			return
		}
		var err error
		if offset, err = strconv.Atoi(key); err != nil {
			http.Error(w, "invalid nextPageKey", http.StatusBadRequest)
			return
		}
	} else {
		if query.Get("schemaIds") != logStorageSchema {
			http.Error(w, "unexpected schema", http.StatusBadRequest)
			return
		}
		f.scans.Add(1)
	}

	f.mu.Lock()
	offset = min(offset, len(f.objectIDs))
	end := min(offset+f.pageSize, len(f.objectIDs))
	resp := fakeListResponse{Items: []fakeListItem{}}
	for _, id := range f.objectIDs[offset:end] {
		resp.Items = append(resp.Items, fakeListItem{ObjectID: id})
	}
	if end < len(f.objectIDs) {
		next := strconv.Itoa(end)
		resp.NextPageKey = &next
	}
	f.mu.Unlock()

	_ = json.NewEncoder(w).Encode(resp)
}

type schemaTypes struct {
	provider, resource tftypes.Type
}

// loadSchemaTypes builds the provider once to get the wire types; each harness still gets its own server.
func loadSchemaTypes(t *testing.T) schemaTypes {
	t.Helper()
	resp, err := provider.NewServer().GetProviderSchema(t.Context(), &tfprotov5.GetProviderSchemaRequest{})
	require.NoError(t, err)
	require.Contains(t, resp.ResourceSchemas, logStorageResource)
	return schemaTypes{
		provider: resp.Provider.ValueType(),
		resource: resp.ResourceSchemas[logStorageResource].ValueType(),
	}
}

// harness wraps one configured provider server talking to a fakeSettingsAPI over gRPC request types.
type harness struct {
	server tfprotov5.ProviderServer
	api    *fakeSettingsAPI
	types  schemaTypes
	config *tfprotov5.DynamicValue
}

func newHarness(t *testing.T, types schemaTypes, objectIDs ...string) *harness {
	t.Helper()
	api := &fakeSettingsAPI{pageSize: 2}
	api.setOrder(objectIDs...)
	srv := httptest.NewServer(api)
	t.Cleanup(srv.Close)

	h := &harness{
		server: provider.NewServer(),
		api:    api,
		types:  types,
		config: dynamicValue(t, types.provider, fmt.Sprintf(`{"dt_env_url":%q,"dt_api_token":"fixture"}`, srv.URL)),
	}
	h.configure(t)
	return h
}

func (h *harness) configure(t *testing.T) {
	t.Helper()
	resp, err := h.server.ConfigureProvider(t.Context(), &tfprotov5.ConfigureProviderRequest{Config: h.config})
	require.NoError(t, err)
	require.Empty(t, resp.Diagnostics)
}

func (h *harness) read(t *testing.T, id string) *tfprotov5.ReadResourceResponse {
	t.Helper()
	resp, err := h.server.ReadResource(t.Context(), &tfprotov5.ReadResourceRequest{
		TypeName:     logStorageResource,
		CurrentState: dynamicValue(t, h.types.resource, fmt.Sprintf(`{"id":%q}`, id)),
	})
	require.NoError(t, err)
	return resp
}

// insertAfter refreshes the object and returns its computed insert_after.
func (h *harness) insertAfter(t *testing.T, id string) string {
	t.Helper()
	resp := h.read(t, id)
	require.Empty(t, resp.Diagnostics)

	state, err := resp.NewState.Unmarshal(h.types.resource)
	require.NoError(t, err)
	var attrs map[string]tftypes.Value
	require.NoError(t, state.As(&attrs))
	var previous string
	require.NoError(t, attrs["insert_after"].As(&previous))
	return previous
}

// delete destroys the object through ApplyResourceChange, like `terraform destroy` would.
func (h *harness) delete(t *testing.T, id string) {
	t.Helper()
	null := dynamicValue(t, h.types.resource, "null")
	resp, err := h.server.ApplyResourceChange(t.Context(), &tfprotov5.ApplyResourceChangeRequest{
		TypeName:     logStorageResource,
		PriorState:   dynamicValue(t, h.types.resource, fmt.Sprintf(`{"id":%q}`, id)),
		PlannedState: null,
		Config:       null,
	})
	require.NoError(t, err)
	require.Empty(t, resp.Diagnostics)
}

func dynamicValue(t *testing.T, typ tftypes.Type, raw string) *tfprotov5.DynamicValue {
	t.Helper()
	value, err := tftypes.ValueFromJSON([]byte(raw), typ)
	require.NoError(t, err)
	wire, err := tfprotov5.NewDynamicValue(typ, value)
	require.NoError(t, err)
	return &wire
}
