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
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const suppliedToken = "header.payload.signature"

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

func diagnosticsFor(t *testing.T, credentials map[string]any) diag.Diagnostics {
	t.Helper()

	data := mockResourceData{"dt_env_url": "https://foo.live.dynatrace.com"}
	for key, value := range credentials {
		data[key] = value
	}
	_, diagnostics := config.ProviderConfigureGenericWithDiagnostics(t.Context(), data)

	return diagnostics
}

func TestWIFVendorIsParsed(t *testing.T) {
	cfg := configuredWith(t, map[string]any{"wif_vendor": "github", "wif_audience": "dynatrace"})

	assert.Equal(t, wif.VendorGitHub, cfg.Platform.WIF.Vendor)
}

// The counterpart of the test above, differing only in the casing of the value. The export command
// reads the configuration without schema validation, so it has to accept what the schema would have
// normalised away.
func TestWIFVendorIsLowercased(t *testing.T) {
	cfg := configuredWith(t, map[string]any{"wif_vendor": "GitHub", "wif_audience": "dynatrace"})

	assert.Equal(t, wif.VendorGitHub, cfg.Platform.WIF.Vendor)
}

func TestWIFAudienceIsParsed(t *testing.T) {
	cfg := configuredWith(t, map[string]any{"wif_vendor": "github", "wif_audience": "https://dynatrace.com"})

	assert.Equal(t, "https://dynatrace.com", cfg.Platform.WIF.Audience)
}

func TestWIFOIDCTokenIsParsed(t *testing.T) {
	cfg := configuredWith(t, map[string]any{"wif_oidc_token": suppliedToken})

	assert.Equal(t, suppliedToken, cfg.Platform.WIF.StaticToken)
}

// The counterpart of the test above, differing only in the trailing newline that a token picks up on
// its way through a CI secret. "Bearer <token>\n" is not a valid header value.
func TestWIFOIDCTokenIsTrimmed(t *testing.T) {
	cfg := configuredWith(t, map[string]any{"wif_oidc_token": suppliedToken + "\n"})

	assert.Equal(t, suppliedToken, cfg.Platform.WIF.StaticToken)
}

func TestNoDiagnosticsWhenOnlyWIFIsConfigured(t *testing.T) {
	diagnostics := diagnosticsFor(t, map[string]any{"wif_vendor": "github", "wif_audience": "dynatrace"})

	assert.Empty(t, diagnostics)
}

// The counterpart of the test above, differing only in the leftover platform token.
func TestWarnsWhenWIFAndPlatformTokenAreBothConfigured(t *testing.T) {
	diagnostics := diagnosticsFor(t, map[string]any{
		"wif_vendor":     "github",
		"wif_audience":   "dynatrace",
		"platform_token": "platform-token",
	})

	require.Len(t, diagnostics, 1)
	assert.Equal(t, diag.Warning, diagnostics[0].Severity)
	assert.Equal(t, "Workload Identity Federation takes precedence over other platform credentials", diagnostics[0].Summary)
	assert.Equal(t, "Workload Identity Federation is configured and is used for all Dynatrace platform API requests. `platform_token` will be ignored for platform requests. OAuth credentials are still used for IAM (Account Management) requests.", diagnostics[0].Detail)
}

func TestWarnsWhenWIFAndOAuthCredentialsAreBothConfigured(t *testing.T) {
	diagnostics := diagnosticsFor(t, map[string]any{
		"wif_vendor":    "github",
		"wif_audience":  "dynatrace",
		"client_id":     "dt0s02.CLIENT",
		"client_secret": "dt0s02.SECRET",
	})

	require.Len(t, diagnostics, 1)
	assert.Equal(t, diag.Warning, diagnostics[0].Severity)
	assert.Equal(t, "Workload Identity Federation is configured and is used for all Dynatrace platform API requests. `client_id` and `client_secret` will be ignored for platform requests. OAuth credentials are still used for IAM (Account Management) requests.", diagnostics[0].Detail)
}

// Reporting that an unusable setup takes precedence over a credential that works would send the user
// looking in the wrong place, so the error replaces the warning rather than joining it.
func TestReportsErrorInsteadOfWarningWhenWIFIsUnusable(t *testing.T) {
	diagnostics := diagnosticsFor(t, map[string]any{
		"wif_vendor":     "github",
		"platform_token": "platform-token",
	})

	require.Len(t, diagnostics, 1)
	assert.Equal(t, diag.Error, diagnostics[0].Severity)
	assert.Equal(t, "Invalid Workload Identity Federation configuration", diagnostics[0].Summary)
	assert.Equal(t, "No audience has been specified for Workload Identity Federation. Use either the configuration attribute `wif_audience` or the environment variable `DYNATRACE_WIF_AUDIENCE` for that", diagnostics[0].Detail)
}

func TestPlatformValidationAcceptsWIFWithoutOAuth(t *testing.T) {
	cfg := configuredWith(t, map[string]any{"wif_vendor": "github", "wif_audience": "dynatrace"})

	_, err := config.ClientSet(cfg, config.CredValPlatform)

	assert.NoError(t, err)
}

func TestPlatformValidationRejectsWIFWithoutAudience(t *testing.T) {
	cfg := configuredWith(t, map[string]any{"wif_vendor": "github"})

	_, err := config.ClientSet(cfg, config.CredValPlatform)

	assert.EqualError(t, err, " No audience has been specified for Workload Identity Federation. Use either the configuration attribute `wif_audience` or the environment variable `DYNATRACE_WIF_AUDIENCE` for that")
}

func TestPlatformValidationRejectsUnsupportedVendor(t *testing.T) {
	cfg := configuredWith(t, map[string]any{"wif_vendor": "gitlab", "wif_audience": "dynatrace"})

	_, err := config.ClientSet(cfg, config.CredValPlatform)

	assert.EqualError(t, err, " `gitlab` is not a supported Workload Identity Federation vendor. The only supported value for `wif_vendor` (`DYNATRACE_WIF_VENDOR`) is `github`")
}

// The schema rejects this combination too, but the export command never runs schema validation.
func TestPlatformValidationRejectsVendorAndOIDCTokenTogether(t *testing.T) {
	cfg := configuredWith(t, map[string]any{
		"wif_vendor":     "github",
		"wif_audience":   "dynatrace",
		"wif_oidc_token": suppliedToken,
	})

	_, err := config.ClientSet(cfg, config.CredValPlatform)

	assert.EqualError(t, err, " A Workload Identity Federation vendor and a pre-minted OIDC token have both been specified. These options are mutually exclusive. Unset either `wif_vendor` (`DYNATRACE_WIF_VENDOR`) or `wif_oidc_token` (`DYNATRACE_WIF_OIDC_TOKEN`)")
}

func TestExportValidationAcceptsWIF(t *testing.T) {
	cfg := configuredWith(t, map[string]any{"wif_vendor": "github", "wif_audience": "dynatrace"})

	_, err := config.ClientSet(cfg, config.CredValExport)

	assert.NoError(t, err)
}

// The counterpart of the test above, differing only in the absence of any credential.
func TestExportValidationRejectsMissingCredentials(t *testing.T) {
	cfg := configuredWith(t, map[string]any{})

	_, err := config.ClientSet(cfg, config.CredValExport)

	assert.EqualError(t, err, " No API Token, Platform Token, Workload Identity Federation, or OAuth has been specified for export. More detailed information can be found in the documentation at https://registry.terraform.io/providers/dynatrace-oss/dynatrace/latest/docs#configure-the-dynatrace-provider")
}
