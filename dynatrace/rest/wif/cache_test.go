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

	"github.com/dynatrace-oss/terraform-provider-dynatrace/provider/envutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// withEmptyCache gives a test the cache to itself, since it outlives any single one of them.
func withEmptyCache(t *testing.T) {
	t.Helper()

	clear(tokenSourceCache)
	t.Cleanup(func() { clear(tokenSourceCache) })

	// A vendor minter reads its credentials from the environment while it is being built.
	t.Setenv(envutils.ActionsIDTokenRequestURL.Key, "https://token.service.invalid/")
	t.Setenv(envutils.ActionsIDTokenRequestToken.Key, "request-token")
}

func TestTokenSourceForReusesSourceForSameConfiguration(t *testing.T) {
	withEmptyCache(t)
	config := Config{Vendor: VendorGitHub, Audience: "dynatrace"}

	first, err := TokenSourceFor(t.Context(), config)
	require.NoError(t, err)
	second, err := TokenSourceFor(t.Context(), config)
	require.NoError(t, err)

	assert.Same(t, first, second)
}

func TestTokenSourceForSeparatesSourcesByAudience(t *testing.T) {
	withEmptyCache(t)

	first, err := TokenSourceFor(t.Context(), Config{Vendor: VendorGitHub, Audience: "dynatrace"})
	require.NoError(t, err)
	second, err := TokenSourceFor(t.Context(), Config{Vendor: VendorGitHub, Audience: "another-audience"})
	require.NoError(t, err)

	assert.NotSame(t, first, second)
}

func TestTokenSourceForReturnsSuppliedStaticToken(t *testing.T) {
	withEmptyCache(t)
	rawToken := jwtWithPayload(`{"exp":1767225600}`)

	source, err := TokenSourceFor(t.Context(), Config{StaticToken: rawToken})
	require.NoError(t, err)

	token, err := source.Token()
	require.NoError(t, err)
	assert.Equal(t, rawToken, token.AccessToken)
}

// Reporting an expiry would make oauth2 discard the token and leave the request unauthenticated.
func TestTokenSourceForLeavesSuppliedStaticTokenWithoutExpiry(t *testing.T) {
	withEmptyCache(t)

	source, err := TokenSourceFor(t.Context(), Config{StaticToken: jwtWithPayload(`{"exp":1767225600}`)})
	require.NoError(t, err)

	token, err := source.Token()
	require.NoError(t, err)
	assert.True(t, token.Expiry.IsZero())
}

func TestTokenSourceForCannotReplaceSuppliedStaticToken(t *testing.T) {
	withEmptyCache(t)
	rawToken := jwtWithPayload(`{"exp":1767225600}`)

	source, err := TokenSourceFor(t.Context(), Config{StaticToken: rawToken})
	require.NoError(t, err)

	token, err := source.Replace(rawToken)
	require.NoError(t, err)
	assert.Equal(t, rawToken, token.AccessToken)
}

// Caching it would make a credential a key in a map that lives as long as the process.
func TestTokenSourceForDoesNotCacheSuppliedStaticToken(t *testing.T) {
	withEmptyCache(t)

	_, err := TokenSourceFor(t.Context(), Config{StaticToken: jwtWithPayload(`{"exp":1767225600}`)})

	require.NoError(t, err)
	assert.Empty(t, tokenSourceCache)
}

func TestTokenSourceForRejectsInvalidConfiguration(t *testing.T) {
	withEmptyCache(t)

	_, err := TokenSourceFor(t.Context(), Config{Vendor: VendorGitHub})

	assert.EqualError(t, err, "no audience has been specified for Workload Identity Federation. Use either the configuration attribute `wif_audience` or the environment variable `DYNATRACE_WIF_AUDIENCE` for that")
}

// A cached failure would outlive a fix to the environment.
func TestTokenSourceForDoesNotCacheCredentialFailure(t *testing.T) {
	withEmptyCache(t)
	t.Setenv(envutils.ActionsIDTokenRequestURL.Key, "")

	_, err := TokenSourceFor(t.Context(), Config{Vendor: VendorGitHub, Audience: "dynatrace"})

	require.Error(t, err)
	assert.Empty(t, tokenSourceCache)
}
