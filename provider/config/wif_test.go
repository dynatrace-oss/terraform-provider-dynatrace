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
	"testing"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest/wif"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/provider/config"
	"github.com/stretchr/testify/assert"
)

// configuredWith builds a provider configuration around a valid environment URL, so that a test only
// has to state the credentials it is about.
func configuredWith(t *testing.T, credentials map[string]any) *config.ProviderConfiguration {
	t.Helper()

	data := mockResourceData{"dt_env_url": "https://foo.live.dynatrace.com"}
	for key, value := range credentials {
		data[key] = value
	}

	return config.ProviderConfigureGeneric(t.Context(), data)
}

// wifBlock returns a fully configured GitHub WIF block as the list-of-maps structure that
// schema.ResourceData.Get("wif") returns. Pass overrides to change specific fields; use "github" as
// a key with a []interface{} value to replace the github sub-block entirely.
func wifBlock(overrides map[string]any) []interface{} {
	block := map[string]any{
		"vendor":   "github",
		"audience": "dynatrace",
		"github": []interface{}{map[string]any{
			"token_request_url":   "https://token.service.invalid/",
			"token_request_token": "request-token",
		}},
	}
	for key, value := range overrides {
		block[key] = value
	}
	return []interface{}{block}
}

func TestWIFVendorIsParsed(t *testing.T) {
	cfg := configuredWith(t, map[string]any{"wif": wifBlock(nil)})

	assert.Equal(t, wif.VendorGitHub, cfg.Platform.WorkloadIdentityFederationConfig.Vendor)
}

// The counterpart of the test above, differing only in the casing of the value. The export command
// reads the configuration without schema validation, so it has to accept what the schema would have
// normalised away.
func TestWIFVendorIsLowercased(t *testing.T) {
	cfg := configuredWith(t, map[string]any{"wif": wifBlock(map[string]any{"vendor": "GitHub"})})

	assert.Equal(t, wif.VendorGitHub, cfg.Platform.WorkloadIdentityFederationConfig.Vendor)
}

func TestWIFAudienceIsParsed(t *testing.T) {
	cfg := configuredWith(t, map[string]any{"wif": wifBlock(map[string]any{"audience": "https://dynatrace.com"})})

	assert.Equal(t, "https://dynatrace.com", cfg.Platform.WorkloadIdentityFederationConfig.Audience)
}

func TestWIFGitHubTokenRequestURLIsParsed(t *testing.T) {
	cfg := configuredWith(t, map[string]any{"wif": wifBlock(map[string]any{
		"github": []interface{}{map[string]any{
			"token_request_url":   "https://token.actions.githubusercontent.com/",
			"token_request_token": "request-token",
		}},
	})})

	assert.Equal(t, "https://token.actions.githubusercontent.com/", cfg.Platform.WorkloadIdentityFederationConfig.GitHub.TokenRequestURL)
}

func TestWIFGitHubTokenRequestTokenIsParsed(t *testing.T) {
	cfg := configuredWith(t, map[string]any{"wif": wifBlock(map[string]any{
		"github": []interface{}{map[string]any{
			"token_request_url":   "https://token.service.invalid/",
			"token_request_token": "my-request-token",
		}},
	})})

	assert.Equal(t, "my-request-token", cfg.Platform.WorkloadIdentityFederationConfig.GitHub.TokenRequestToken)
}

func TestPlatformValidationAcceptsWIFWithoutOAuth(t *testing.T) {
	cfg := configuredWith(t, map[string]any{"wif": wifBlock(nil)})

	_, err := config.ClientSet(cfg, config.CredValPlatform)

	assert.NoError(t, err)
}

func TestPlatformValidationRejectsWIFWithoutAudience(t *testing.T) {
	cfg := configuredWith(t, map[string]any{"wif": wifBlock(map[string]any{"audience": ""})})

	_, err := config.ClientSet(cfg, config.CredValPlatform)

	assert.EqualError(t, err, " No audience has been specified for Workload Identity Federation. Use either the configuration attribute `wif.audience` or the environment variable `DYNATRACE_WIF_AUDIENCE` for that")
}

func TestPlatformValidationRejectsWIFWithoutGitHubTokenRequestURL(t *testing.T) {
	cfg := configuredWith(t, map[string]any{"wif": wifBlock(map[string]any{
		"github": []interface{}{map[string]any{"token_request_url": "", "token_request_token": "request-token"}},
	})})

	_, err := config.ClientSet(cfg, config.CredValPlatform)

	assert.EqualError(t, err, " No GitHub Actions token request URL has been configured. Use either the configuration attribute `wif.github.token_request_url` or run this job in GitHub Actions with `permissions: { id-token: write }` (which injects `ACTIONS_ID_TOKEN_REQUEST_URL`)")
}

func TestPlatformValidationRejectsWIFWithoutGitHubTokenRequestToken(t *testing.T) {
	cfg := configuredWith(t, map[string]any{"wif": wifBlock(map[string]any{
		"github": []interface{}{map[string]any{"token_request_url": "https://token.service.invalid/", "token_request_token": ""}},
	})})

	_, err := config.ClientSet(cfg, config.CredValPlatform)

	assert.EqualError(t, err, " No GitHub Actions token request token has been configured. Use either the configuration attribute `wif.github.token_request_token` or run this job in GitHub Actions with `permissions: { id-token: write }` (which injects `ACTIONS_ID_TOKEN_REQUEST_TOKEN`)")
}

func TestPlatformValidationRejectsUnsupportedVendor(t *testing.T) {
	cfg := configuredWith(t, map[string]any{"wif": wifBlock(map[string]any{"vendor": "gitlab"})})

	_, err := config.ClientSet(cfg, config.CredValPlatform)

	assert.EqualError(t, err, " `gitlab` is not a supported Workload Identity Federation vendor. The only supported value for `wif.vendor` (`DYNATRACE_WIF_VENDOR`) is `github`")
}

func TestExportValidationAcceptsWIF(t *testing.T) {
	cfg := configuredWith(t, map[string]any{"wif": wifBlock(nil)})

	_, err := config.ClientSet(cfg, config.CredValExport)

	assert.NoError(t, err)
}

// The counterpart of the test above, differing only in the absence of any credential.
func TestExportValidationRejectsMissingCredentials(t *testing.T) {
	cfg := configuredWith(t, map[string]any{})

	_, err := config.ClientSet(cfg, config.CredValExport)

	assert.EqualError(t, err, " No API Token, Platform Token, Workload Identity Federation, or OAuth has been specified for export. More detailed information can be found in the documentation at https://registry.terraform.io/providers/dynatrace-oss/dynatrace/latest/docs#configure-the-dynatrace-provider")
}
