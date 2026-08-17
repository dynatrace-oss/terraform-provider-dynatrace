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

package wif

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateAcceptsUnconfigured(t *testing.T) {
	assert.NoError(t, Validate(Config{}))
}

func TestValidateAcceptsVendorWithAudience(t *testing.T) {
	assert.NoError(t, Validate(Config{Vendor: VendorGitHub, Audience: "dynatrace"}))
}

func TestValidateRejectsVendorWithoutAudience(t *testing.T) {
	err := Validate(Config{Vendor: VendorGitHub})

	assert.EqualError(t, err, "No audience has been specified for Workload Identity Federation. Use either the configuration attribute `wif_audience` or the environment variable `DYNATRACE_WIF_AUDIENCE` for that")
}

func TestValidateRejectsUnsupportedVendor(t *testing.T) {
	err := Validate(Config{Vendor: "aws", Audience: "dynatrace"})

	assert.EqualError(t, err, "`aws` is not a supported Workload Identity Federation vendor. The only supported value for `wif_vendor` (`DYNATRACE_WIF_VENDOR`) is `github`")
}

// The audience belongs to the vendor that mints a token, so a supplied one is complete without it.
func TestValidateAcceptsStaticTokenWithoutAudience(t *testing.T) {
	assert.NoError(t, Validate(Config{StaticToken: jwtWithPayload(`{"exp":1767225600}`)}))
}

func TestValidateRejectsStaticTokenThatIsNotAJWT(t *testing.T) {
	err := Validate(Config{StaticToken: "not-a-jwt"})

	assert.EqualError(t, err, "The value of `wif_oidc_token` (`DYNATRACE_WIF_OIDC_TOKEN`) is not a JWT: expected three dot-separated segments")
}

func TestValidateRejectsVendorAndStaticTokenTogether(t *testing.T) {
	err := Validate(Config{Vendor: VendorGitHub, Audience: "dynatrace", StaticToken: jwtWithPayload(`{"exp":1767225600}`)})

	assert.EqualError(t, err, "A Workload Identity Federation vendor and a pre-minted OIDC token have both been specified. These options are mutually exclusive. Unset either `wif_vendor` (`DYNATRACE_WIF_VENDOR`) or `wif_oidc_token` (`DYNATRACE_WIF_OIDC_TOKEN`)")
}
