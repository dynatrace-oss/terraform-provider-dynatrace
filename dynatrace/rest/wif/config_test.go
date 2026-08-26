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

// Callers for which federation is optional guard with Configured, so absence is a problem by the time
// Validate is asked.
func TestValidateRejectsUnconfigured(t *testing.T) {
	assert.ErrorIs(t, Config{}.Validate(), errNotConfigured)
}

func TestValidateAcceptsVendorWithAudience(t *testing.T) {
	assert.NoError(t, Config{Vendor: VendorGitHub, Audience: "dynatrace"}.Validate())
}

func TestValidateRejectsVendorWithoutAudience(t *testing.T) {
	err := Config{Vendor: VendorGitHub}.Validate()

	assert.ErrorIs(t, err, errNoAudience)
}

func TestValidateRejectsUnsupportedVendor(t *testing.T) {
	err := Config{Vendor: "aws", Audience: "dynatrace"}.Validate()

	assert.EqualError(t, err, "`aws` is not a supported vendor, the only supported one is `github`")
}

// The audience belongs to the vendor that mints a token, so a supplied one is complete without it.
func TestValidateAcceptsStaticTokenWithoutAudience(t *testing.T) {
	assert.NoError(t, Config{StaticToken: jwtWithPayload(`{"exp":1767225600}`)}.Validate())
}

func TestValidateRejectsStaticTokenThatIsNotAJWT(t *testing.T) {
	err := Config{StaticToken: "not-a-jwt"}.Validate()

	assert.ErrorIs(t, err, errStaticTokenNotAJWT)
}

func TestValidateRejectsVendorAndStaticTokenTogether(t *testing.T) {
	err := Config{Vendor: VendorGitHub, Audience: "dynatrace", StaticToken: jwtWithPayload(`{"exp":1767225600}`)}.Validate()

	assert.ErrorIs(t, err, errVendorAndStaticToken)
}
