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
		"ACTIONS_ID_TOKEN_REQUEST_URL",
		"ACTIONS_ID_TOKEN_REQUEST_TOKEN",
	} {
		t.Setenv(variable, "")
	}

	return provider.Provider().Validate(terraform.NewResourceConfigRaw(config))
}

// InternalValidate is the only thing that catches a ConflictsWith or RequiredWith naming an
// attribute that does not exist. The other tests in this package skip, so without this the provider
// schema is never checked at all.
func TestProviderSchemaIsInternallyValid(t *testing.T) {
	require.NoError(t, provider.Provider().InternalValidate())
}

func TestWIFRequiresVendor(t *testing.T) {
	diagnostics := validateProviderConfig(t, map[string]any{
		"wif": []interface{}{map[string]any{"audience": "dynatrace"}},
	})

	require.Len(t, diagnostics, 1)
	assert.Equal(t, `The argument "wif.0.vendor" is required, but no definition was found.`, diagnostics[0].Detail)
}

func TestWIFRequiresAudience(t *testing.T) {
	diagnostics := validateProviderConfig(t, map[string]any{
		"wif": []interface{}{map[string]any{"vendor": "github"}},
	})

	require.Len(t, diagnostics, 1)
	assert.Equal(t, `The argument "wif.0.audience" is required, but no definition was found.`, diagnostics[0].Detail)
}

func TestWIFVendorRejectsUnsupportedVendor(t *testing.T) {
	diagnostics := validateProviderConfig(t, map[string]any{
		"wif": []interface{}{map[string]any{
			"vendor":   "gitlab",
			"audience": "dynatrace",
		}},
	})

	require.Len(t, diagnostics, 1)
	assert.Equal(t, `expected wif.0.vendor to be one of ["github"], got gitlab`, diagnostics[0].Summary)
}

func TestWIFVendorAcceptsSupportedVendor(t *testing.T) {
	diagnostics := validateProviderConfig(t, map[string]any{
		"wif": []interface{}{map[string]any{
			"vendor":   "github",
			"audience": "dynatrace",
		}},
	})

	assert.Empty(t, diagnostics)
}
