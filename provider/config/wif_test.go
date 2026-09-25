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

package config_test

import (
	"maps"
	"testing"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/provider/config"
	"github.com/stretchr/testify/assert"
)

// configuredWith builds a provider configuration around a valid environment URL, so that a test only
// has to state the credentials it is about.
func configuredWith(t *testing.T, credentials map[string]any) *config.ProviderConfiguration {
	t.Helper()

	data := mockResourceData{"dt_env_url": "https://foo.live.dynatrace.com"}
	maps.Copy(data, credentials)

	return config.ProviderConfigureGeneric(t.Context(), data)
}

func TestWIFAudienceIsParsed(t *testing.T) {
	cfg := configuredWith(t, map[string]any{"wif_audience": "https://dynatrace.com"})

	assert.Equal(t, "https://dynatrace.com", cfg.Platform.WorkloadIdentityFederationAudience)
}

func TestWIFAudienceIsTrimmed(t *testing.T) {
	cfg := configuredWith(t, map[string]any{"wif_audience": "  https://dynatrace.com  "})

	assert.Equal(t, "https://dynatrace.com", cfg.Platform.WorkloadIdentityFederationAudience)
}

func TestPlatformValidationAcceptsWIFWithoutOAuth(t *testing.T) {
	cfg := configuredWith(t, map[string]any{"wif_audience": "dynatrace"})

	_, err := config.ClientSet(cfg, config.CredValPlatform)

	assert.NoError(t, err)
}

// The counterpart of the test above, differing only in the absence of the audience: without it there
// is nothing left to authenticate a platform request with.
func TestPlatformValidationRejectsMissingCredentials(t *testing.T) {
	cfg := configuredWith(t, map[string]any{})

	_, err := config.ClientSet(cfg, config.CredValPlatform)

	assert.EqualError(t, err, " No OAuth Client ID for the Automation API has been specified. Use either the environment variable `DT_AUTOMATION_CLIENT_ID` or the configuration attribute `automation_client_id` of the provider for that")
}

func TestExportValidationAcceptsWIF(t *testing.T) {
	cfg := configuredWith(t, map[string]any{"wif_audience": "dynatrace"})

	_, err := config.ClientSet(cfg, config.CredValExport)

	assert.NoError(t, err)
}

// The counterpart of the test above, differing only in the absence of any credential.
func TestExportValidationRejectsMissingCredentials(t *testing.T) {
	cfg := configuredWith(t, map[string]any{})

	_, err := config.ClientSet(cfg, config.CredValExport)

	assert.EqualError(t, err, " No API Token, Platform Token, Workload Identity Federation, or OAuth has been specified for export. More detailed information can be found in the documentation at https://registry.terraform.io/providers/dynatrace-oss/dynatrace/latest/docs#configure-the-dynatrace-provider")
}
