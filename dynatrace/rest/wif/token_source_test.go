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
	"slices"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

const concurrentCallers = 10

type concurrentMinter struct {
	calls atomic.Int64
}

func (minter *concurrentMinter) mint(context.Context) (string, error) {
	call := minter.calls.Add(1)

	// Long enough that callers arriving together overlap, without which these tests could not tell
	// single flight minting apart from callers that merely ran one after another.
	time.Sleep(10 * time.Millisecond)

	return jwtWithPayload(fmt.Sprintf(`{"exp":%d,"jti":%d}`, time.Now().Add(time.Hour).Unix(), call)), nil
}

func callConcurrently(call func(caller int)) {
	var arrived, finished sync.WaitGroup
	arrived.Add(concurrentCallers)
	finished.Add(concurrentCallers)

	for caller := range concurrentCallers {
		go func() {
			defer finished.Done()
			arrived.Done()
			arrived.Wait()
			call(caller)
		}()
	}

	finished.Wait()
}

func newTokenSource(t *testing.T, minter minter) *tokenSource {
	return &tokenSource{mintContext: t.Context(), minter: minter}
}

func TestTokenSourceReturnsMintedToken(t *testing.T) {
	rawToken := tokenExpiringIn(time.Hour)

	token, err := newTokenSource(t, &recordingMinter{tokens: []string{rawToken}}).Token()

	require.NoError(t, err)
	assert.Equal(t, rawToken, token.AccessToken)
}

func TestTokenSourceLeavesTokenTypeEmpty(t *testing.T) {
	token, err := newTokenSource(t, &recordingMinter{tokens: []string{tokenExpiringIn(time.Hour)}}).Token()

	require.NoError(t, err)
	assert.Equal(t, "", token.TokenType)
}

func TestTokenSourceReportsExpiryClaimedByToken(t *testing.T) {
	expiry := time.Now().Add(time.Hour)
	rawToken := jwtWithPayload(fmt.Sprintf(`{"exp":%d}`, expiry.Unix()))

	token, err := newTokenSource(t, &recordingMinter{tokens: []string{rawToken}}).Token()

	require.NoError(t, err)
	assert.Equal(t, time.Unix(expiry.Unix(), 0), token.Expiry)
}

func TestTokenSourceReusesTokenUntilRefreshTime(t *testing.T) {
	minter := &recordingMinter{tokens: []string{tokenExpiringIn(time.Hour)}}
	source := newTokenSource(t, minter)

	first, err := source.Token()
	require.NoError(t, err)
	second, err := source.Token()
	require.NoError(t, err)

	assert.Equal(t, first.AccessToken, second.AccessToken)
	assert.Equal(t, 1, minter.calls)
}

// The tests below move refreshAt into the past rather than waiting: the passage of time is the one
// input they cannot supply directly.
func TestTokenSourceMintsReplacementAfterRefreshTime(t *testing.T) {
	minter := &recordingMinter{tokens: []string{tokenExpiringIn(time.Hour), tokenExpiringIn(2 * time.Hour)}}
	source := newTokenSource(t, minter)

	_, err := source.Token()
	require.NoError(t, err)
	source.refreshAt = time.Now().Add(-time.Second)

	second, err := source.Token()

	require.NoError(t, err)
	assert.Equal(t, minter.tokens[1], second.AccessToken)
	assert.Equal(t, 2, minter.calls)
}

func TestTokenSourceMintsOnceForConcurrentCallers(t *testing.T) {
	minter := &concurrentMinter{}
	source := newTokenSource(t, minter)

	callConcurrently(func(int) {
		_, err := source.Token()
		assert.NoError(t, err)
	})

	assert.Equal(t, int64(1), minter.calls.Load())
}

func TestTokenSourceRefreshesAheadOfExpiry(t *testing.T) {
	source := newTokenSource(t, &recordingMinter{tokens: []string{tokenExpiringIn(time.Hour)}})

	token, err := source.Token()

	require.NoError(t, err)
	assert.Equal(t, token.Expiry.Add(-refreshMargin), source.refreshAt)
}

func TestTokenSourceFloorsRefreshTimeForShortLivedToken(t *testing.T) {
	source := newTokenSource(t, &recordingMinter{tokens: []string{tokenExpiringIn(10 * time.Second)}})

	_, err := source.Token()

	require.NoError(t, err)
	assert.WithinDuration(t, time.Now().Add(minimumMintInterval), source.refreshAt, time.Second)
}

func TestTokenSourcePropagatesMintingFailure(t *testing.T) {
	mintErr := errors.New("the token service is unreachable")

	_, err := newTokenSource(t, &recordingMinter{err: mintErr}).Token()

	assert.ErrorIs(t, err, mintErr)
}

func TestTokenSourceServesUnexpiredTokenWhenMintingFails(t *testing.T) {
	source := newTokenSource(t, &recordingMinter{tokens: []string{tokenExpiringIn(time.Hour)}})
	held, err := source.Token()
	require.NoError(t, err)
	source.refreshAt = time.Now().Add(-time.Second)
	source.minter = &recordingMinter{err: errors.New("the token service is unreachable")}

	token, err := source.Token()

	require.NoError(t, err)
	assert.Equal(t, held.AccessToken, token.AccessToken)
}

func TestTokenSourcePropagatesMintingFailureOnceTokenHasExpired(t *testing.T) {
	mintErr := errors.New("the token service is unreachable")
	// Arriving expired is the shortest route to a source holding an expired token; the refresh margin
	// means it never happens this way outside a test.
	source := newTokenSource(t, &recordingMinter{tokens: []string{tokenExpiringIn(-time.Second)}})
	_, err := source.Token()
	require.NoError(t, err)
	source.refreshAt = time.Now().Add(-time.Second)
	source.minter = &recordingMinter{err: mintErr}

	_, err = source.Token()

	assert.ErrorIs(t, err, mintErr)
}

func TestTokenSourcePostponesMintingAfterFailure(t *testing.T) {
	source := newTokenSource(t, &recordingMinter{tokens: []string{tokenExpiringIn(time.Hour)}})
	_, err := source.Token()
	require.NoError(t, err)
	source.refreshAt = time.Now().Add(-time.Second)
	source.minter = &recordingMinter{err: errors.New("the token service is unreachable")}

	_, err = source.Token()

	require.NoError(t, err)
	assert.WithinDuration(t, time.Now().Add(mintRetryInterval), source.refreshAt, time.Second)
}

func TestTokenSourcePostponesMintingNoFurtherThanExpiry(t *testing.T) {
	source := newTokenSource(t, &recordingMinter{tokens: []string{tokenExpiringIn(mintRetryInterval / 2)}})
	held, err := source.Token()
	require.NoError(t, err)
	source.refreshAt = time.Now().Add(-time.Second)
	source.minter = &recordingMinter{err: errors.New("the token service is unreachable")}

	_, err = source.Token()

	require.NoError(t, err)
	assert.Equal(t, held.Expiry, source.refreshAt)
}

func TestTokenSourceRejectsTokenWithoutExpiry(t *testing.T) {
	_, err := newTokenSource(t, &recordingMinter{tokens: []string{jwtWithPayload(`{"aud":"dynatrace"}`)}}).Token()

	assert.EqualError(t, err, "the ID token has no `exp` claim, so the provider cannot tell when to obtain a replacement")
}

func TestTokenSourceReplacesRejectedToken(t *testing.T) {
	minter := &recordingMinter{tokens: []string{tokenExpiringIn(time.Hour), tokenExpiringIn(2 * time.Hour)}}
	source := newTokenSource(t, minter)
	rejected, err := source.Token()
	require.NoError(t, err)

	replacement, err := source.Replace(rejected.AccessToken)

	require.NoError(t, err)
	assert.Equal(t, minter.tokens[1], replacement.AccessToken)
}

// The cancellation is only a marker, identifying the context the minter was handed.
func TestTokenSourceMintsWithTheContextItWasGiven(t *testing.T) {
	minter := &recordingMinter{tokens: []string{tokenExpiringIn(time.Hour)}}
	cancelledContext, cancel := context.WithCancel(t.Context())
	cancel()
	source := &tokenSource{mintContext: cancelledContext, minter: minter}

	_, err := source.Token()

	require.NoError(t, err)
	assert.ErrorIs(t, minter.lastContext.Err(), context.Canceled)
}

func TestTokenSourceKeepsTokenThatWasAlreadyReplaced(t *testing.T) {
	minter := &recordingMinter{tokens: []string{tokenExpiringIn(time.Hour)}}
	source := newTokenSource(t, minter)
	held, err := source.Token()
	require.NoError(t, err)

	replacement, err := source.Replace("the token another request had rejected")

	require.NoError(t, err)
	assert.Equal(t, held.AccessToken, replacement.AccessToken)
	assert.Equal(t, 1, minter.calls)
}

func TestTokenSourceReplacesOnceForConcurrentRejections(t *testing.T) {
	minter := &concurrentMinter{}
	source := newTokenSource(t, minter)
	rejected, err := source.Token()
	require.NoError(t, err)

	replacements := make([]string, concurrentCallers)
	callConcurrently(func(caller int) {
		replacement, err := source.Replace(rejected.AccessToken)
		if assert.NoError(t, err) {
			replacements[caller] = replacement.AccessToken
		}
	})

	assert.Equal(t, int64(2), minter.calls.Load())
	assert.Equal(t, slices.Repeat([]string{source.token.AccessToken}, concurrentCallers), replacements)
}

func TestTokenSourceDoesNotFallBackOnRejectedToken(t *testing.T) {
	mintErr := errors.New("the token service is unreachable")
	source := newTokenSource(t, &recordingMinter{tokens: []string{tokenExpiringIn(time.Hour)}})
	rejected, err := source.Token()
	require.NoError(t, err)
	source.minter = &recordingMinter{err: mintErr}
	_, err = source.Replace(rejected.AccessToken)
	require.ErrorIs(t, err, mintErr)

	_, err = source.Token()

	assert.ErrorIs(t, err, mintErr)
}
