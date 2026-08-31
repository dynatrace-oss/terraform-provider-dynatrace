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

package rest_test

import (
	"testing"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest"
	"github.com/dynatrace/dynatrace-configuration-as-code-core/clients"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateClusterClients(t *testing.T) {
	const clusterURL = "https://cluster.example.com"

	for name, create := range map[string]func(baseURL string, apiToken string) error{
		"v1": func(baseURL string, apiToken string) error {
			_, err := rest.CreateClusterV1Client(t.Context(), baseURL, apiToken)
			return err
		},
		"v2": func(baseURL string, apiToken string) error {
			_, err := rest.CreateClusterV2Client(t.Context(), baseURL, apiToken)
			return err
		},
	} {
		t.Run(name+"/no cluster URL", func(t *testing.T) {
			assert.ErrorIs(t, create("", mockToken), rest.NoClusterURLError)
		})

		t.Run(name+"/no cluster URL takes precedence over missing token", func(t *testing.T) {
			assert.ErrorIs(t, create("", ""), rest.NoClusterURLError)
		})

		t.Run(name+"/no cluster API token", func(t *testing.T) {
			assert.ErrorIs(t, create(clusterURL, ""), clients.ErrAccessTokenMissing)
		})

		t.Run(name+"/fully configured", func(t *testing.T) {
			require.NoError(t, create(clusterURL, mockToken))
		})
	}
}
