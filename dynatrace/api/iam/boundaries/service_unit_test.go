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
	"net/url"
	"strconv"
	"testing"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/iam"
	testing2 "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/testing"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	testAccountID   = "test-account-id"
	testEndpointURL = "https://api.dynatrace.com"
)

func newTestClient(mock iam.IAMClient) *BoundaryServiceClient {
	return &BoundaryServiceClient{
		iamClientGetter: &testing2.MockIAMClientGetter{Client: mock},
		accountID:       testAccountID,
		endpointURL:     testEndpointURL,
	}
}

func boundaryPageResponse(items []PolicyBoundary) []byte {
	data, _ := json.Marshal(ListPolicyBoundariesResponse{PolicyBoundaries: items})
	return data
}

func pageURL(page int) string {
	params := url.Values{}
	params.Set("size", strconv.Itoa(maxPageSize))
	params.Set("page", strconv.Itoa(page))
	return fmt.Sprintf("%s/iam/v1/repo/account/%s/boundaries?%s", testEndpointURL, testAccountID, params.Encode())
}

func TestBoundaryServiceClient_List(t *testing.T) {
	t.Run("Returns error on GET failure", func(t *testing.T) {
		mock := &testing2.MockIAMClient{
			GETFunc: func(_ context.Context, _ string, _ int, _ bool) ([]byte, error) {
				return nil, assert.AnError
			},
		}
		_, err := newTestClient(mock).List(t.Context())
		assert.ErrorIs(t, err, assert.AnError)
	})

	t.Run("Returns error on invalid JSON response", func(t *testing.T) {
		mock := &testing2.MockIAMClient{
			GETFunc: func(_ context.Context, _ string, _ int, _ bool) ([]byte, error) {
				return []byte("not-json"), nil
			},
		}
		_, err := newTestClient(mock).List(t.Context())
		expectedErr := &json.SyntaxError{}
		assert.ErrorAs(t, err, &expectedErr)
	})

	t.Run("Returns empty stubs for empty response", func(t *testing.T) {
		mock := &testing2.MockIAMClient{
			GETFunc: func(_ context.Context, _ string, _ int, _ bool) ([]byte, error) {
				return boundaryPageResponse(nil), nil
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
			GETFunc: func(_ context.Context, reqURL string, _ int, _ bool) ([]byte, error) {
				callCount++
				assert.Equal(t, pageURL(1), reqURL)
				return boundaryPageResponse(items), nil
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
			GETFunc: func(_ context.Context, reqURL string, _ int, _ bool) ([]byte, error) {
				callCount++
				switch callCount {
				case 1:
					assert.Equal(t, pageURL(1), reqURL)
					return boundaryPageResponse(fullPage), nil
				case 2:
					assert.Equal(t, pageURL(2), reqURL)
					return boundaryPageResponse(page2), nil
				default:
					return nil, fmt.Errorf("unexpected page request %d", callCount)
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
			GETFunc: func(_ context.Context, _ string, _ int, _ bool) ([]byte, error) {
				callCount++
				if callCount == 1 {
					return boundaryPageResponse(fullPage), nil
				}
				return nil, assert.AnError
			},
		}

		_, err := newTestClient(mock).List(t.Context())
		assert.ErrorIs(t, err, assert.AnError)
	})
}
