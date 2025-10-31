//go:build unit

/**
* @license
* Copyright 2020 Dynatrace LLC
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

package segments_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/grail/segments"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestList(t *testing.T) {
	t.Run("Returns an error if the client creation fails", func(t *testing.T) {
		_, err := segments.Service(&rest.Credentials{}).List(t.Context())
		assert.ErrorIs(t, err, rest.NoPlatformCredentialsErr)
	})

	t.Run("Returns the list of segments ignoring ready-made ones", func(t *testing.T) {
		apiResponse := `{"filterSegments": [
			{"uid": "1", "name": "ready-made-false", "isPublic": false, "isReadyMade": false},
			{"uid": "2", "name": "ready-made-true", "isPublic": false, "isReadyMade": true},
			{"uid": "3", "name": "ready-made-missing","isPublic": false}]}`

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, http.MethodGet, r.Method)
			require.Equal(t, "/platform/storage/filter-segments/v1/filter-segments:lean", r.URL.Path)

			w.WriteHeader(http.StatusOK)
			w.Write([]byte(apiResponse))
		}))
		defer server.Close()

		result, err := segments.Service(&rest.Credentials{
			OAuth: rest.OAuthCredentials{
				EnvironmentURL: server.URL,
				PlatformToken:  "token",
			},
		}).List(t.Context())

		require.NoError(t, err)
		assert.Len(t, result, 2)
		for _, segment := range result {
			assert.True(t, segment.ID != "2", "ready-made segment must be ignored")
		}
	})
}
