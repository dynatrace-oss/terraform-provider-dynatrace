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

package config_test

import (
	"testing"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/provider/config"
	"github.com/stretchr/testify/assert"
)

type mockResourceData map[string]any

func (mrd mockResourceData) Get(k string) any {
	return mrd[k]
}

func TestProviderConfigureGeneric(t *testing.T) {
	t.Run("Direct fields are passed through and URLs cleaned", func(t *testing.T) {
		cfg := config.ProviderConfigureGeneric(t.Context(), mockResourceData{
			"dt_env_url":           "https://foo.live.dynatrace.com",
			"dt_api_token":         "api-token",
			"dt_cluster_api_token": "cluster-token",
			"dt_cluster_url":       "https://cluster.example.com/ ",
			"platform_token":       "platform-token",
		})
		assert.Equal(t, "api-token", cfg.APIToken)
		assert.Equal(t, "cluster-token", cfg.ClusterAPIToken)
		assert.Equal(t, "https://cluster.example.com", cfg.ClusterAPIV2URL)
		assert.Equal(t, "platform-token", cfg.Platform.PlatformToken)
	})

	t.Run("Environment URL normalization", func(t *testing.T) {
		for _, tc := range []struct {
			name             string
			envURL           string
			expectedClassic  string
			expectedPlatform string
		}{
			{
				name:             "Saas live",
				envURL:           "https://foo.live.dynatrace.com",
				expectedClassic:  "https://foo.live.dynatrace.com",
				expectedPlatform: "https://foo.apps.dynatrace.com",
			},
			{
				name:             "Saas apps normalized to live",
				envURL:           "https://foo.apps.dynatrace.com",
				expectedClassic:  "https://foo.live.dynatrace.com",
				expectedPlatform: "https://foo.apps.dynatrace.com",
			},
			{
				name:             "Trailing slash and space trimmed",
				envURL:           "https://foo.live.dynatrace.com/ ",
				expectedClassic:  "https://foo.live.dynatrace.com",
				expectedPlatform: "https://foo.apps.dynatrace.com",
			},
			{
				name:             "Sprint",
				envURL:           "https://foo.sprint.dynatracelabs.com",
				expectedClassic:  "https://foo.sprint.dynatracelabs.com",
				expectedPlatform: "https://foo.sprint.apps.dynatracelabs.com",
			},
			{
				name:             "Sprint apps normalized",
				envURL:           "https://foo.sprint.apps.dynatracelabs.com",
				expectedClassic:  "https://foo.sprint.dynatracelabs.com",
				expectedPlatform: "https://foo.sprint.apps.dynatracelabs.com",
			},
			{
				name:             "Dev",
				envURL:           "https://foo.dev.dynatracelabs.com",
				expectedClassic:  "https://foo.dev.dynatracelabs.com",
				expectedPlatform: "https://foo.dev.apps.dynatracelabs.com",
			},
			{
				name:             "Dev apps normalized",
				envURL:           "https://foo.dev.apps.dynatracelabs.com",
				expectedClassic:  "https://foo.dev.dynatracelabs.com",
				expectedPlatform: "https://foo.dev.apps.dynatracelabs.com",
			},
			{
				name:             "Managed passthrough",
				envURL:           "https://managed.example.com/e/abc",
				expectedClassic:  "https://managed.example.com/e/abc",
				expectedPlatform: "https://managed.example.com/e/abc",
			},
			{
				name:             "Empty",
				envURL:           "",
				expectedClassic:  "",
				expectedPlatform: "",
			},
		} {
			t.Run(tc.name, func(t *testing.T) {
				cfg := config.ProviderConfigureGeneric(t.Context(), mockResourceData{"dt_env_url": tc.envURL})
				assert.Equal(t, tc.expectedClassic, cfg.EnvironmentURL)
				assert.Equal(t, tc.expectedPlatform, cfg.Platform.EnvironmentURL)
			})
		}
	})

	t.Run("Explicit platform environment URL overrides derivation", func(t *testing.T) {
		cfg := config.ProviderConfigureGeneric(t.Context(), mockResourceData{
			"dt_env_url":         "https://foo.live.dynatrace.com",
			"automation_env_url": "https://custom.apps.example.com/",
		})
		assert.Equal(t, "https://foo.live.dynatrace.com", cfg.EnvironmentURL)
		assert.Equal(t, "https://custom.apps.example.com", cfg.Platform.EnvironmentURL)
	})

	t.Run("Client ID precedence", func(t *testing.T) {
		for _, tc := range []struct {
			name             string
			data             mockResourceData
			expectedIAM      string
			expectedPlatform string
		}{
			{
				name: "Specific keys win per audience",
				data: mockResourceData{
					"iam_client_id":        "iam",
					"client_id":            "generic",
					"automation_client_id": "automation",
				},
				expectedIAM:      "iam",
				expectedPlatform: "automation",
			},
			{
				name: "Generic client_id shared by both",
				data: mockResourceData{
					"client_id": "generic",
				},
				expectedIAM:      "generic",
				expectedPlatform: "generic",
			},
			{
				name: "Only iam_client_id falls through to platform",
				data: mockResourceData{
					"iam_client_id": "iam",
				},
				expectedIAM:      "iam",
				expectedPlatform: "iam",
			},
			{
				name: "Only automation_client_id falls through to IAM",
				data: mockResourceData{
					"automation_client_id": "automation",
				},
				expectedIAM:      "automation",
				expectedPlatform: "automation",
			},
		} {
			t.Run(tc.name, func(t *testing.T) {
				cfg := config.ProviderConfigureGeneric(t.Context(), tc.data)
				assert.Equal(t, tc.expectedIAM, cfg.IAM.ClientID)
				assert.Equal(t, tc.expectedPlatform, cfg.Platform.ClientID)
			})
		}
	})

	t.Run("Client secret precedence", func(t *testing.T) {
		for _, tc := range []struct {
			name             string
			data             mockResourceData
			expectedIAM      string
			expectedPlatform string
		}{
			{
				name: "Specific keys win per audience",
				data: mockResourceData{
					"iam_client_secret":        "iam",
					"client_secret":            "generic",
					"automation_client_secret": "automation",
				},
				expectedIAM:      "iam",
				expectedPlatform: "automation",
			},
			{
				name: "Generic client_secret shared by both",
				data: mockResourceData{
					"client_secret": "generic",
				},
				expectedIAM:      "generic",
				expectedPlatform: "generic",
			},
			{
				name: "Only iam_client_secret falls through to platform",
				data: mockResourceData{
					"iam_client_secret": "iam",
				},
				expectedIAM:      "iam",
				expectedPlatform: "iam",
			},
			{
				name: "Only automation_client_secret falls through to IAM",
				data: mockResourceData{
					"automation_client_secret": "automation",
				},
				expectedIAM:      "automation",
				expectedPlatform: "automation",
			},
		} {
			t.Run(tc.name, func(t *testing.T) {
				cfg := config.ProviderConfigureGeneric(t.Context(), tc.data)
				assert.Equal(t, tc.expectedIAM, cfg.IAM.ClientSecret)
				assert.Equal(t, tc.expectedPlatform, cfg.Platform.ClientSecret)
			})
		}
	})

	t.Run("Account ID precedence and urn prefix trimming", func(t *testing.T) {
		for _, tc := range []struct {
			name     string
			data     mockResourceData
			expected string
		}{
			{
				name: "iam_account_id wins",
				data: mockResourceData{
					"iam_account_id": "iam",
					"account_id":     "generic",
				},
				expected: "iam",
			},
			{
				name: "account_id fallback",
				data: mockResourceData{
					"account_id": "generic",
				},
				expected: "generic",
			},
			{
				name: "URN prefix trimmed",
				data: mockResourceData{
					"iam_account_id": "urn:dtaccount:1234",
				},
				expected: "1234",
			},
			{
				name:     "Empty",
				data:     mockResourceData{},
				expected: "",
			},
		} {
			t.Run(tc.name, func(t *testing.T) {
				assert.Equal(t, tc.expected, config.ProviderConfigureGeneric(t.Context(), tc.data).IAM.AccountID)
			})
		}
	})

	t.Run("Token URL precedence", func(t *testing.T) {
		for _, tc := range []struct {
			name             string
			data             mockResourceData
			expectedIAM      string
			expectedPlatform string
		}{
			{
				name: "Specific keys win per audience",
				data: mockResourceData{
					"iam_token_url":        "https://iam.example.com",
					"token_url":            "https://generic.example.com",
					"automation_token_url": "https://automation.example.com",
				},
				expectedIAM:      "https://iam.example.com",
				expectedPlatform: "https://automation.example.com",
			},
			{
				name: "Generic token_url shared by both",
				data: mockResourceData{
					"token_url": "https://generic.example.com",
				},
				expectedIAM:      "https://generic.example.com",
				expectedPlatform: "https://generic.example.com",
			},
			{
				name: "Only automation_token_url shared by both",
				data: mockResourceData{
					"automation_token_url": "https://automation.example.com",
				},
				expectedIAM:      "https://automation.example.com",
				expectedPlatform: "https://automation.example.com",
			},
			{
				name: "Only iam_token_url shared by both",
				data: mockResourceData{
					"iam_token_url": "https://iam.example.com",
				},
				expectedIAM:      "https://iam.example.com",
				expectedPlatform: "https://iam.example.com",
			},
			{
				name: "Trailing slash trimmed",
				data: mockResourceData{
					"token_url": "https://generic.example.com/",
				},
				expectedIAM:      "https://generic.example.com",
				expectedPlatform: "https://generic.example.com",
			},
		} {
			t.Run(tc.name, func(t *testing.T) {
				cfg := config.ProviderConfigureGeneric(t.Context(), tc.data)
				assert.Equal(t, tc.expectedIAM, cfg.IAM.TokenURL)
				assert.Equal(t, tc.expectedPlatform, cfg.Platform.TokenURL)
			})
		}
	})

	t.Run("Token and IAM endpoint URL defaults derived from environment", func(t *testing.T) {
		for _, tc := range []struct {
			name             string
			envURL           string
			expectedTokenURL string
			expectedEndpoint string
		}{
			{
				name:             "Saas",
				envURL:           "https://foo.live.dynatrace.com",
				expectedTokenURL: config.ProdTokenURL,
				expectedEndpoint: config.ProdIAMEndpointURL,
			},
			{
				name:             "Sprint",
				envURL:           "https://foo.sprint.dynatracelabs.com",
				expectedTokenURL: config.SprintTokenURL,
				expectedEndpoint: config.SprintIAMEndpointURL,
			},
			{
				name:             "Dev",
				envURL:           "https://foo.dev.dynatracelabs.com",
				expectedTokenURL: config.DevTokenURL,
				expectedEndpoint: config.DevIAMEndpointURL,
			},
			{
				name:             "Managed defaults to prod",
				envURL:           "https://managed.example.com/e/abc",
				expectedTokenURL: config.ProdTokenURL,
				expectedEndpoint: config.ProdIAMEndpointURL,
			},
			{
				name:             "Empty defaults to prod",
				envURL:           "",
				expectedTokenURL: config.ProdTokenURL,
				expectedEndpoint: config.ProdIAMEndpointURL,
			},
		} {
			t.Run(tc.name, func(t *testing.T) {
				cfg := config.ProviderConfigureGeneric(t.Context(), mockResourceData{"dt_env_url": tc.envURL})
				assert.Equal(t, tc.expectedTokenURL, cfg.IAM.TokenURL)
				assert.Equal(t, tc.expectedTokenURL, cfg.Platform.TokenURL)
				assert.Equal(t, tc.expectedEndpoint, cfg.IAM.EndpointURL)
			})
		}
	})

	t.Run("Explicit IAM endpoint URL overrides derivation", func(t *testing.T) {
		cfg := config.ProviderConfigureGeneric(t.Context(), mockResourceData{
			"dt_env_url":       "https://foo.sprint.dynatracelabs.com",
			"iam_endpoint_url": "https://custom-endpoint.example.com/",
		})
		assert.Equal(t, "https://custom-endpoint.example.com", cfg.IAM.EndpointURL)
	})
}
