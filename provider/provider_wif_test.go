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

package provider_test

import (
	"sort"
	"testing"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/provider"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// validateProviderConfig runs the provider block through schema validation.
//
// The Workload Identity Federation variables are cleared first because the SDK resolves DefaultFunc
// before validating, so anything left in the environment of the machine running the tests would take
// part in the very rules under test here. MultiEnvDefaultFunc treats an empty value as unset.
func validateProviderConfig(t *testing.T, config map[string]any) diag.Diagnostics {
	t.Helper()

	for _, variable := range []string{
		"DYNATRACE_WIF_VENDOR", "DT_WIF_VENDOR",
		"DYNATRACE_WIF_AUDIENCE", "DT_WIF_AUDIENCE",
		"DYNATRACE_WIF_OIDC_TOKEN", "DT_WIF_OIDC_TOKEN",
	} {
		t.Setenv(variable, "")
	}

	return provider.Provider().Validate(terraform.NewResourceConfigRaw(config))
}

func detailsOf(diagnostics diag.Diagnostics) []string {
	details := make([]string, 0, len(diagnostics))
	for _, diagnostic := range diagnostics {
		details = append(details, diagnostic.Detail)
	}
	sort.Strings(details)

	return details
}

// InternalValidate is the only thing that catches a ConflictsWith or RequiredWith naming an
// attribute that does not exist. The other tests in this package skip, so without this the provider
// schema is never checked at all.
func TestProviderSchemaIsInternallyValid(t *testing.T) {
	require.NoError(t, provider.Provider().InternalValidate())
}

func TestWIFVendorConflictsWithWIFOIDCToken(t *testing.T) {
	diagnostics := validateProviderConfig(t, map[string]any{
		"wif_vendor":     "github",
		"wif_audience":   "dynatrace",
		"wif_oidc_token": "header.payload.signature",
	})

	// Both attributes declare the conflict, so each one is reported. The order they arrive in follows
	// map iteration, hence the sort.
	require.Len(t, diagnostics, 2)
	assert.Equal(t, []string{
		`"wif_oidc_token": conflicts with wif_vendor`,
		`"wif_vendor": conflicts with wif_oidc_token`,
	}, detailsOf(diagnostics))
}

// A missing audience is not a schema rule, because the audience may just as well arrive through an
// environment variable, and because the export command never runs schema validation. The credential
// validation rejects it instead - see TestPlatformValidationRejectsWIFWithoutAudience.
func TestWIFVendorWithoutAudienceIsNotASchemaError(t *testing.T) {
	diagnostics := validateProviderConfig(t, map[string]any{"wif_vendor": "github"})

	assert.Empty(t, diagnostics)
}

func TestWIFVendorRejectsUnsupportedVendor(t *testing.T) {
	diagnostics := validateProviderConfig(t, map[string]any{
		"wif_vendor":   "gitlab",
		"wif_audience": "dynatrace",
	})

	require.Len(t, diagnostics, 1)
	assert.Equal(t, `expected wif_vendor to be one of ["github"], got gitlab`, diagnostics[0].Summary)
}

func TestWIFVendorAcceptsSupportedVendor(t *testing.T) {
	diagnostics := validateProviderConfig(t, map[string]any{
		"wif_vendor":   "github",
		"wif_audience": "dynatrace",
	})

	assert.Empty(t, diagnostics)
}

// A supplied token stands on its own: it needs no vendor and no audience.
func TestWIFOIDCTokenNeedsNoOtherAttribute(t *testing.T) {
	diagnostics := validateProviderConfig(t, map[string]any{"wif_oidc_token": "header.payload.signature"})

	assert.Empty(t, diagnostics)
}
