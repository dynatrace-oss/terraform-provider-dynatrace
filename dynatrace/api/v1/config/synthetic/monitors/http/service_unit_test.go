//go:build unit

/*
 * @license
 * Copyright 2026 Dynatrace LLC
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package http_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/synthetic/monitors"
	http2 "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/synthetic/monitors/http"
	httpSettings "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/synthetic/monitors/http/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest"
	mocktesting "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockService = mocktesting.MockCRUDService[*httpSettings.SyntheticMonitor]

func tempScriptMonitor() *httpSettings.SyntheticMonitor {
	return &httpSettings.SyntheticMonitor{Script: http2.GetTempScript()}
}

func TestService_Get(t *testing.T) {
	t.Run("Returns error when underlying Get fails", func(t *testing.T) {
		mock := &mockService{
			GetFunc: func(ctx context.Context, id string, v *httpSettings.SyntheticMonitor) error {
				return assert.AnError
			},
		}
		svc := http2.WithPredefinedService(mock)
		err := svc.Get(t.Context(), "id-1", &httpSettings.SyntheticMonitor{})
		assert.ErrorIs(t, err, assert.AnError)
	})

	t.Run("Returns settings as-is when script is not the auto-generated one", func(t *testing.T) {
		userDesc := "user-defined"
		realScript := &httpSettings.Script{Requests: httpSettings.Requests{&httpSettings.Request{
			Description: &userDesc,
			URL:         "https://example.com",
			Method:      "GET",
		}}}
		mock := &mockService{
			GetFunc: func(ctx context.Context, id string, v *httpSettings.SyntheticMonitor) error {
				v.Script = realScript
				return nil
			},
		}
		svc := http2.WithPredefinedService(mock)
		v := &httpSettings.SyntheticMonitor{}
		require.NoError(t, svc.Get(t.Context(), "id-1", v))
		assert.Same(t, realScript, v.Script)
		assert.Nil(t, v.NoScript)
	})

	t.Run("Substitutes auto-generated temp script with NoScript=true", func(t *testing.T) {
		mock := &mockService{
			GetFunc: func(ctx context.Context, id string, v *httpSettings.SyntheticMonitor) error {
				*v = *tempScriptMonitor()
				return nil
			},
		}
		svc := http2.WithPredefinedService(mock)
		v := &httpSettings.SyntheticMonitor{}
		require.NoError(t, svc.Get(t.Context(), "id-1", v))
		assert.Nil(t, v.Script)
		require.NotNil(t, v.NoScript)
		assert.True(t, *v.NoScript)
	})
}

func TestService_Create(t *testing.T) {
	t.Run("Returns error when underlying Create fails", func(t *testing.T) {
		mock := &mockService{
			CreateFunc: func(ctx context.Context, v *httpSettings.SyntheticMonitor) (*api.Stub, error) {
				return nil, assert.AnError
			},
		}
		svc := http2.WithPredefinedService(mock)
		_, err := svc.Create(t.Context(), &httpSettings.SyntheticMonitor{})
		assert.ErrorIs(t, err, assert.AnError)
	})

	t.Run("Substitutes temp script when NoScript=true and Script is nil", func(t *testing.T) {
		var captured *httpSettings.SyntheticMonitor
		mock := &mockService{
			CreateFunc: func(ctx context.Context, v *httpSettings.SyntheticMonitor) (*api.Stub, error) {
				captured = v
				return nil, assert.AnError
			},
		}
		svc := http2.WithPredefinedService(mock)
		_, err := svc.Create(t.Context(), &httpSettings.SyntheticMonitor{NoScript: new(true)})

		require.ErrorIs(t, err, assert.AnError)
		require.NotNil(t, captured)
		require.NotNil(t, captured.Script)
		require.Len(t, captured.Script.Requests, 1)
		assert.Equal(t, *http2.GetTempScript().Requests[0].Description, *captured.Script.Requests[0].Description)
	})

	t.Run("Does not overwrite existing Script when NoScript=true", func(t *testing.T) {
		existing := &httpSettings.Script{Requests: httpSettings.Requests{&httpSettings.Request{
			URL:    "https://example.com",
			Method: "GET",
		}}}
		var captured *httpSettings.SyntheticMonitor
		mock := &mockService{
			CreateFunc: func(ctx context.Context, v *httpSettings.SyntheticMonitor) (*api.Stub, error) {
				captured = v
				return nil, assert.AnError
			},
		}
		svc := http2.WithPredefinedService(mock)
		noScript := true
		_, err := svc.Create(t.Context(), &httpSettings.SyntheticMonitor{NoScript: &noScript, Script: existing})

		require.ErrorIs(t, err, assert.AnError)
		require.NotNil(t, captured)
		assert.Same(t, existing, captured.Script)
	})

	t.Run("Propagates non-retryable Get error from validateReady", func(t *testing.T) {
		mock := &mockService{
			CreateFunc: func(ctx context.Context, v *httpSettings.SyntheticMonitor) (*api.Stub, error) {
				return &api.Stub{ID: "new-id"}, nil
			},
			GetFunc: func(ctx context.Context, id string, v *httpSettings.SyntheticMonitor) error {
				return assert.AnError
			},
		}
		svc := http2.WithPredefinedService(mock)
		_, err := svc.Create(t.Context(), &httpSettings.SyntheticMonitor{})
		assert.ErrorIs(t, err, assert.AnError)
	})
}

func TestService_Update(t *testing.T) {
	t.Run("Returns error when underlying Update fails", func(t *testing.T) {
		mock := &mockService{
			UpdateFunc: func(ctx context.Context, id string, v *httpSettings.SyntheticMonitor) error {
				return assert.AnError
			},
		}
		svc := http2.WithPredefinedService(mock)
		err := svc.Update(t.Context(), "id-1", &httpSettings.SyntheticMonitor{})
		assert.ErrorIs(t, err, assert.AnError)
	})

	t.Run("Pre-fetches existing script when NoScript=true and Script is nil", func(t *testing.T) {
		existing := &httpSettings.Script{Requests: httpSettings.Requests{&httpSettings.Request{
			URL:    "https://example.com",
			Method: "GET",
		}}}
		var captured *httpSettings.SyntheticMonitor
		mock := &mockService{
			GetFunc: func(ctx context.Context, id string, v *httpSettings.SyntheticMonitor) error {
				v.Script = existing
				return nil
			},
			UpdateFunc: func(ctx context.Context, id string, v *httpSettings.SyntheticMonitor) error {
				captured = v
				return assert.AnError
			},
		}
		svc := http2.WithPredefinedService(mock)
		noScript := true
		err := svc.Update(t.Context(), "id-1", &httpSettings.SyntheticMonitor{NoScript: &noScript})

		require.ErrorIs(t, err, assert.AnError)
		require.NotNil(t, captured)
		assert.Same(t, existing, captured.Script)
	})

	t.Run("Skips pre-fetch when Script is already provided", func(t *testing.T) {
		getCalled := false
		mock := &mockService{
			GetFunc: func(ctx context.Context, id string, v *httpSettings.SyntheticMonitor) error {
				getCalled = true
				return nil
			},
			UpdateFunc: func(ctx context.Context, id string, v *httpSettings.SyntheticMonitor) error {
				return assert.AnError
			},
		}
		svc := http2.WithPredefinedService(mock)
		noScript := true
		err := svc.Update(t.Context(), "id-1", &httpSettings.SyntheticMonitor{
			NoScript: &noScript,
			Script:   &httpSettings.Script{},
		})
		require.ErrorIs(t, err, assert.AnError)
		assert.False(t, getCalled, "pre-fetch Get should not be called when Script is provided")
	})

	t.Run("Returns error when pre-fetch Get fails", func(t *testing.T) {
		mock := &mockService{
			GetFunc: func(ctx context.Context, id string, v *httpSettings.SyntheticMonitor) error {
				return assert.AnError
			},
		}
		svc := http2.WithPredefinedService(mock)
		noScript := true
		err := svc.Update(t.Context(), "id-1", &httpSettings.SyntheticMonitor{NoScript: &noScript})
		assert.ErrorIs(t, err, assert.AnError)
	})

	t.Run("Propagates non-retryable Get error from validateReady", func(t *testing.T) {
		mock := &mockService{
			UpdateFunc: func(ctx context.Context, id string, v *httpSettings.SyntheticMonitor) error {
				return nil
			},
			GetFunc: func(ctx context.Context, id string, v *httpSettings.SyntheticMonitor) error {
				return assert.AnError
			},
		}
		svc := http2.WithPredefinedService(mock)
		err := svc.Update(t.Context(), "id-1", &httpSettings.SyntheticMonitor{})
		assert.ErrorIs(t, err, assert.AnError)
	})
}

func TestService_Delete(t *testing.T) {
	t.Run("Returns error when underlying Delete fails", func(t *testing.T) {
		mock := &mockService{
			DeleteFunc: func(ctx context.Context, id string) error {
				return assert.AnError
			},
		}
		svc := http2.WithPredefinedService(mock)
		err := svc.Delete(t.Context(), "id-1")
		assert.ErrorIs(t, err, assert.AnError)
	})

	t.Run("Returns nil and forwards id on success", func(t *testing.T) {
		var capturedID string
		mock := &mockService{
			DeleteFunc: func(ctx context.Context, id string) error {
				capturedID = id
				return nil
			},
		}
		svc := http2.WithPredefinedService(mock)
		require.NoError(t, svc.Delete(t.Context(), "id-1"))
		assert.Equal(t, "id-1", capturedID)
	})
}

func TestService_ReadyCheck(t *testing.T) {
	t.Run("Returns retryable error on 404", func(t *testing.T) {
		mock := &mockService{
			GetFunc: func(ctx context.Context, id string, v *httpSettings.SyntheticMonitor) error {
				return rest.Error{Code: http.StatusNotFound}
			},
		}
		svc := http2.WithPredefinedService(mock)

		result := svc.ReadyCheck(t.Context(), "id-1", false)()
		require.NotNil(t, result)
		assert.True(t, result.Retryable)
		assert.True(t, rest.IsNotFoundError(result.Err), "wrapped err should be a 404")
	})

	t.Run("Returns non-retryable error on other Get failure", func(t *testing.T) {
		mock := &mockService{
			GetFunc: func(ctx context.Context, id string, v *httpSettings.SyntheticMonitor) error {
				return assert.AnError
			},
		}
		svc := http2.WithPredefinedService(mock)

		result := svc.ReadyCheck(t.Context(), "id-1", false)()
		require.NotNil(t, result)
		assert.False(t, result.Retryable)
		assert.ErrorIs(t, result.Err, assert.AnError)
	})

	t.Run("Returns retryable consistency error when tags expected but missing", func(t *testing.T) {
		mock := &mockService{
			GetFunc: func(ctx context.Context, id string, v *httpSettings.SyntheticMonitor) error {
				v.Tags = nil
				return nil
			},
		}
		svc := http2.WithPredefinedService(mock)

		result := svc.ReadyCheck(t.Context(), "id-1", true)()
		require.NotNil(t, result)
		assert.True(t, result.Retryable)
		assert.ErrorIs(t, result.Err, http2.ErrConsistencyRetry)
	})

	t.Run("Counts as success when tags expected and present", func(t *testing.T) {
		mock := &mockService{
			GetFunc: func(ctx context.Context, id string, v *httpSettings.SyntheticMonitor) error {
				v.Tags = monitors.TagsWithSourceInfo{{Key: "env"}}
				return nil
			},
		}
		svc := http2.WithPredefinedService(mock)

		check := svc.ReadyCheck(t.Context(), "id-1", true)
		for i := 1; i < http2.DesiredReadySuccesses; i++ {
			result := check()
			require.NotNil(t, result, "iteration %d", i)
			assert.True(t, result.Retryable, "iteration %d should still be retrying", i)
			assert.ErrorIs(t, result.Err, http2.ErrConsistencyRetry)
		}
		assert.Nil(t, check(), "final iteration should succeed")
	})

	t.Run("Returns nil only after desired number of successes", func(t *testing.T) {
		callCount := 0
		mock := &mockService{
			GetFunc: func(ctx context.Context, id string, v *httpSettings.SyntheticMonitor) error {
				callCount++
				return nil
			},
		}
		svc := http2.WithPredefinedService(mock)

		check := svc.ReadyCheck(t.Context(), "id-1", false)
		for i := 1; i < http2.DesiredReadySuccesses; i++ {
			result := check()
			require.NotNil(t, result, "iteration %d", i)
			assert.True(t, result.Retryable)
			assert.ErrorIs(t, result.Err, http2.ErrConsistencyRetry)
		}
		final := check()
		assert.Nil(t, final)
		assert.Equal(t, http2.DesiredReadySuccesses, callCount)
	})
}
