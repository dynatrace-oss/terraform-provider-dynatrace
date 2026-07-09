//go:build unit

/**
* @license
* Copyright 2025 Dynatrace LLC
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

package boundaries

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"testing"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/iam"
	testing2 "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/testing"
	"github.com/dynatrace/dynatrace-configuration-as-code-core/api"
	rest2 "github.com/dynatrace/dynatrace-configuration-as-code-core/api/rest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testAccountID = "test-account-id"

func newTestClient(mock iam.IAMClient) *BoundaryServiceClient {
	return &BoundaryServiceClient{
		iamClientGetter: &testing2.MockIAMClientGetter{Client: mock},
		accountID:       testAccountID,
	}
}

func boundaryPageResponse(items []PolicyBoundary) []byte {
	data, _ := json.Marshal(ListPolicyBoundariesResponse{PolicyBoundaries: items})
	return data
}

func boundariesPath() string {
	return fmt.Sprintf("/iam/v1/repo/account/%s/boundaries", testAccountID)
}

func assertPageRequest(t *testing.T, reqURL string, options rest2.RequestOptions, page int) {
	assert.Equal(t, boundariesPath(), reqURL)
	assert.Equal(t, strconv.Itoa(maxPageSize), options.QueryParams.Get("size"))
	assert.Equal(t, strconv.Itoa(page), options.QueryParams.Get("page"))
}

func TestBoundaryServiceClient_List(t *testing.T) {
	t.Run("Returns error on GET failure", func(t *testing.T) {
		mock := &testing2.MockIAMClient{
			GETFunc: func(_ context.Context, _ string, _ rest2.RequestOptions) (api.Response, error) {
				return api.Response{}, assert.AnError
			},
		}
		_, err := newTestClient(mock).List(t.Context())
		assert.ErrorIs(t, err, assert.AnError)
	})

	t.Run("Returns error on invalid JSON response", func(t *testing.T) {
		mock := &testing2.MockIAMClient{
			GETFunc: func(_ context.Context, _ string, _ rest2.RequestOptions) (api.Response, error) {
				return api.Response{Data: []byte("not-json")}, nil
			},
		}
		_, err := newTestClient(mock).List(t.Context())
		expectedErr := &json.SyntaxError{}
		assert.ErrorAs(t, err, &expectedErr)
	})

	t.Run("Returns empty stubs for empty response", func(t *testing.T) {
		mock := &testing2.MockIAMClient{
			GETFunc: func(_ context.Context, _ string, _ rest2.RequestOptions) (api.Response, error) {
				return api.Response{Data: boundaryPageResponse(nil)}, nil
			},
		}
		stubs, err := newTestClient(mock).List(t.Context())
		require.NoError(t, err)
		assert.Empty(t, stubs)
	})

	t.Run("Returns stubs for single page", func(t *testing.T) {
		items := []PolicyBoundary{
			{UUID: "uuid-1", Name: "boundary-1"},
			{UUID: "uuid-2", Name: "boundary-2"},
		}
		callCount := 0
		mock := &testing2.MockIAMClient{
			GETFunc: func(_ context.Context, reqURL string, options rest2.RequestOptions) (api.Response, error) {
				callCount++
				assertPageRequest(t, reqURL, options, 1)
				return api.Response{Data: boundaryPageResponse(items)}, nil
			},
		}

		stubs, err := newTestClient(mock).List(t.Context())
		require.NoError(t, err)
		assert.Equal(t, 1, callCount)
		require.Len(t, stubs, 2)
		assert.Equal(t, "uuid-1", stubs[0].ID)
		assert.Equal(t, "boundary-1", stubs[0].Name)
		assert.Equal(t, "uuid-2", stubs[1].ID)
		assert.Equal(t, "boundary-2", stubs[1].Name)
	})

	t.Run("Paginates until partial page", func(t *testing.T) {
		fullPage := make([]PolicyBoundary, maxPageSize)
		for i := range fullPage {
			fullPage[i] = PolicyBoundary{UUID: fmt.Sprintf("uuid-p1-%d", i), Name: fmt.Sprintf("b-p1-%d", i)}
		}
		page2 := []PolicyBoundary{
			{UUID: "uuid-p2-0", Name: "b-p2-0"},
			{UUID: "uuid-p2-1", Name: "b-p2-1"},
		}

		callCount := 0
		mock := &testing2.MockIAMClient{
			GETFunc: func(_ context.Context, reqURL string, options rest2.RequestOptions) (api.Response, error) {
				callCount++
				switch callCount {
				case 1:
					assertPageRequest(t, reqURL, options, 1)
					return api.Response{Data: boundaryPageResponse(fullPage)}, nil
				case 2:
					assertPageRequest(t, reqURL, options, 2)
					return api.Response{Data: boundaryPageResponse(page2)}, nil
				default:
					return api.Response{}, fmt.Errorf("unexpected page request %d", callCount)
				}
			},
		}

		stubs, err := newTestClient(mock).List(t.Context())
		require.NoError(t, err)
		assert.Equal(t, 2, callCount)
		assert.Len(t, stubs, maxPageSize+2)
	})

	t.Run("Returns error on second page GET failure", func(t *testing.T) {
		fullPage := make([]PolicyBoundary, maxPageSize)
		callCount := 0
		mock := &testing2.MockIAMClient{
			GETFunc: func(_ context.Context, _ string, _ rest2.RequestOptions) (api.Response, error) {
				callCount++
				if callCount == 1 {
					return api.Response{Data: boundaryPageResponse(fullPage)}, nil
				}
				return api.Response{}, assert.AnError
			},
		}

		_, err := newTestClient(mock).List(t.Context())
		assert.ErrorIs(t, err, assert.AnError)
	})
}
