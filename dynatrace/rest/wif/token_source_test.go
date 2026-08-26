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
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/provider/envutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/oauth2"
)

// Running out of prepared tokens is a test error rather than a silent repeat, so that a test
// claiming a token was reused cannot pass by accident.
type recordingMinter struct {
	tokens      []string
	err         error
	calls       int
	lastContext context.Context
}

func (minter *recordingMinter) mint(ctx context.Context) (string, error) {
	minter.calls++
	minter.lastContext = ctx
	if minter.err != nil {
		return "", minter.err
	}
	if minter.calls > len(minter.tokens) {
		return "", fmt.Errorf("minter asked for token %d but only %d were prepared", minter.calls, len(minter.tokens))
	}
	return minter.tokens[minter.calls-1], nil
}

func tokenExpiringIn(validity time.Duration) string {
	return jwtWithPayload(fmt.Sprintf(`{"exp":%d}`, time.Now().Add(validity).Unix()))
}

// A vendor minter reads these while it is being built.
func setMintingCredentials(t *testing.T) {
	t.Helper()

	t.Setenv(envutils.ActionsIDTokenRequestURL.Key, "https://token.service.invalid/")
	t.Setenv(envutils.ActionsIDTokenRequestToken.Key, "request-token")
}

func newMintingTokenSource(t *testing.T, minter minter) *mintingTokenSource {
	return &mintingTokenSource{mintContext: t.Context(), minter: minter}
}

func TestMintingTokenSourceReturnsMintedToken(t *testing.T) {
	rawToken := tokenExpiringIn(time.Hour)

	token, err := newMintingTokenSource(t, &recordingMinter{tokens: []string{rawToken}}).Token()

	require.NoError(t, err)
	assert.Equal(t, rawToken, token.AccessToken)
}

func TestMintingTokenSourceLeavesTokenTypeEmpty(t *testing.T) {
	token, err := newMintingTokenSource(t, &recordingMinter{tokens: []string{tokenExpiringIn(time.Hour)}}).Token()

	require.NoError(t, err)
	assert.Equal(t, "", token.TokenType)
}

func TestMintingTokenSourceReportsExpiryClaimedByToken(t *testing.T) {
	expiry := time.Now().Add(time.Hour)
	rawToken := jwtWithPayload(fmt.Sprintf(`{"exp":%d}`, expiry.Unix()))

	token, err := newMintingTokenSource(t, &recordingMinter{tokens: []string{rawToken}}).Token()

	require.NoError(t, err)
	assert.Equal(t, time.Unix(expiry.Unix(), 0), token.Expiry)
}

// A token is minted on every call. Reuse is the job of the source wrapping this one.
func TestMintingTokenSourceMintsOnEveryCall(t *testing.T) {
	minter := &recordingMinter{tokens: []string{tokenExpiringIn(time.Hour), tokenExpiringIn(2 * time.Hour)}}
	source := newMintingTokenSource(t, minter)

	_, err := source.Token()
	require.NoError(t, err)
	second, err := source.Token()

	require.NoError(t, err)
	assert.Equal(t, minter.tokens[1], second.AccessToken)
	assert.Equal(t, 2, minter.calls)
}

func TestMintingTokenSourcePropagatesMintingFailure(t *testing.T) {
	mintErr := errors.New("the token service is unreachable")

	_, err := newMintingTokenSource(t, &recordingMinter{err: mintErr}).Token()

	assert.ErrorIs(t, err, mintErr)
}

func TestMintingTokenSourceRejectsTokenWithoutExpiry(t *testing.T) {
	_, err := newMintingTokenSource(t, &recordingMinter{tokens: []string{jwtWithPayload(`{"aud":"dynatrace"}`)}}).Token()

	assert.ErrorIs(t, err, errNoExpiryClaim)
}

// The cancellation is only a marker, identifying the context the minter was handed.
func TestMintingTokenSourceMintsWithTheContextItWasGiven(t *testing.T) {
	minter := &recordingMinter{tokens: []string{tokenExpiringIn(time.Hour)}}
	cancelledContext, cancel := context.WithCancel(t.Context())
	cancel()
	source := &mintingTokenSource{mintContext: cancelledContext, minter: minter}

	_, err := source.Token()

	require.NoError(t, err)
	assert.ErrorIs(t, minter.lastContext.Err(), context.Canceled)
}

// Only one token is prepared, so a second trip to the minter fails the test.
func TestReusingTokenSourceKeepsTokenUntilItIsDue(t *testing.T) {
	minter := &recordingMinter{tokens: []string{tokenExpiringIn(time.Hour)}}
	source := reusingTokenSource(t.Context(), minter)

	first, err := source.Token()
	require.NoError(t, err)
	second, err := source.Token()

	require.NoError(t, err)
	assert.Equal(t, first.AccessToken, second.AccessToken)
	assert.Equal(t, 1, minter.calls)
}

// A token whose remaining life is shorter than the margin is due the moment it arrives. Were the
// margin not passed on, oauth2 would apply its ten second default and reuse this one instead.
func TestReusingTokenSourceReplacesTokenInsideTheRefreshMargin(t *testing.T) {
	minter := &recordingMinter{tokens: []string{tokenExpiringIn(refreshMargin / 2), tokenExpiringIn(time.Hour)}}
	source := reusingTokenSource(t.Context(), minter)

	_, err := source.Token()
	require.NoError(t, err)
	second, err := source.Token()

	require.NoError(t, err)
	assert.Equal(t, minter.tokens[1], second.AccessToken)
	assert.Equal(t, 2, minter.calls)
}

func TestReusingTokenSourcePropagatesMintingFailure(t *testing.T) {
	mintErr := errors.New("the token service is unreachable")

	_, err := reusingTokenSource(t.Context(), &recordingMinter{err: mintErr}).Token()

	assert.ErrorIs(t, err, mintErr)
}

// Handing oauth2 a token without an expiry is what keeps it from discarding a supplied one.
func TestStaticTokenSourceHandsOutTheSuppliedTokenWithoutExpiry(t *testing.T) {
	token, err := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: "supplied.id.token"}).Token()

	require.NoError(t, err)
	assert.Equal(t, "supplied.id.token", token.AccessToken)
	assert.True(t, token.Expiry.IsZero())
}

func TestTokenSourceForBuildsSourceForVendor(t *testing.T) {
	setMintingCredentials(t)

	source, err := TokenSourceFor(t.Context(), Config{Vendor: VendorGitHub, Audience: "dynatrace"})

	require.NoError(t, err)
	assert.NotNil(t, source)
}

func TestTokenSourceForReportsMissingVendorCredentials(t *testing.T) {
	t.Setenv(envutils.ActionsIDTokenRequestURL.Key, "")
	t.Setenv(envutils.ActionsIDTokenRequestToken.Key, "")

	_, err := TokenSourceFor(t.Context(), Config{Vendor: VendorGitHub, Audience: "dynatrace"})

	assert.ErrorIs(t, err, errMissingCredential)
}

func TestTokenSourceForReturnsSuppliedStaticToken(t *testing.T) {
	rawToken := jwtWithPayload(`{"exp":1767225600}`)

	source, err := TokenSourceFor(t.Context(), Config{StaticToken: rawToken})
	require.NoError(t, err)

	token, err := source.Token()
	require.NoError(t, err)
	assert.Equal(t, rawToken, token.AccessToken)
}

// Reporting an expiry would make oauth2 discard the token and leave the request unauthenticated.
func TestTokenSourceForLeavesSuppliedStaticTokenWithoutExpiry(t *testing.T) {
	source, err := TokenSourceFor(t.Context(), Config{StaticToken: jwtWithPayload(`{"exp":1767225600}`)})
	require.NoError(t, err)

	token, err := source.Token()
	require.NoError(t, err)
	assert.True(t, token.Expiry.IsZero())
}

// Without this, an unconfigured setup used to yield a token source handing out an empty bearer token.
func TestTokenSourceForRejectsUnconfiguredFederation(t *testing.T) {
	_, err := TokenSourceFor(t.Context(), Config{})

	assert.ErrorIs(t, err, errNotConfigured)
}

func TestTokenSourceForRejectsInvalidConfiguration(t *testing.T) {
	_, err := TokenSourceFor(t.Context(), Config{Vendor: VendorGitHub})

	assert.ErrorIs(t, err, errNoAudience)
}
