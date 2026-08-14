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

// recordingMinter hands out the given tokens in order and counts how often it was asked. Running out
// of tokens is a test error rather than a silent repeat, so that a test claiming a token was reused
// cannot pass by accident.
type recordingMinter struct {
	tokens []string
	err    error
	calls  int
}

func (minter *recordingMinter) mint(ctx context.Context) (string, error) {
	minter.calls++
	if minter.err != nil {
		return "", minter.err
	}
	if minter.calls > len(minter.tokens) {
		return "", fmt.Errorf("minter asked for token %d but only %d were prepared", minter.calls, len(minter.tokens))
	}
	return minter.tokens[minter.calls-1], nil
}

// tokenExpiringIn builds a JWT whose expiry claim lies the given duration ahead.
func tokenExpiringIn(validity time.Duration) string {
	return jwtWithPayload(fmt.Sprintf(`{"exp":%d}`, time.Now().Add(validity).Unix()))
}

// concurrentCallers is how many goroutines the tests below send at a source at once, standing in for
// the resources Terraform applies in parallel.
const concurrentCallers = 10

// concurrentMinter counts what it is asked for and hands out a distinguishable token each time.
type concurrentMinter struct {
	calls atomic.Int64
}

func (minter *concurrentMinter) mint(context.Context) (string, error) {
	call := minter.calls.Add(1)

	// Minting takes long enough that callers arriving together are inside this call at the same time.
	// Without the delay a test could not tell single flight minting apart from callers that merely
	// happened to run one after another.
	time.Sleep(10 * time.Millisecond)

	return jwtWithPayload(fmt.Sprintf(`{"exp":%d,"jti":%d}`, time.Now().Add(time.Hour).Unix(), call)), nil
}

// callConcurrently runs call in concurrentCallers goroutines that all set off together, and returns
// once every one of them has finished.
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

func newTokenSource(minter minter) *tokenSource {
	return &tokenSource{mintContext: context.Background(), minter: minter}
}

func TestTokenSourceReturnsMintedToken(t *testing.T) {
	rawToken := tokenExpiringIn(time.Hour)

	token, err := newTokenSource(&recordingMinter{tokens: []string{rawToken}}).Token()

	require.NoError(t, err)
	assert.Equal(t, rawToken, token.AccessToken)
}

// An empty token type is what makes oauth2 put the token on the wire as a bearer token, which is the
// whole reason these tokens can be used the way a platform token is.
func TestTokenSourceLeavesTokenTypeEmpty(t *testing.T) {
	token, err := newTokenSource(&recordingMinter{tokens: []string{tokenExpiringIn(time.Hour)}}).Token()

	require.NoError(t, err)
	assert.Equal(t, "", token.TokenType)
}

func TestTokenSourceReportsExpiryClaimedByToken(t *testing.T) {
	expiry := time.Now().Add(time.Hour)
	rawToken := jwtWithPayload(fmt.Sprintf(`{"exp":%d}`, expiry.Unix()))

	token, err := newTokenSource(&recordingMinter{tokens: []string{rawToken}}).Token()

	require.NoError(t, err)
	assert.Equal(t, time.Unix(expiry.Unix(), 0), token.Expiry)
}

// Only one token is prepared, so a second trip to the minter fails the test rather than quietly
// returning an equal-looking token.
func TestTokenSourceReusesTokenUntilRefreshTime(t *testing.T) {
	minter := &recordingMinter{tokens: []string{tokenExpiringIn(time.Hour)}}
	source := newTokenSource(minter)

	first, err := source.Token()
	require.NoError(t, err)
	second, err := source.Token()
	require.NoError(t, err)

	assert.Equal(t, first.AccessToken, second.AccessToken)
	assert.Equal(t, 1, minter.calls)
}

// The refresh time is moved into the past rather than waiting for it: the passage of time is the one
// input of this behaviour that a test cannot supply directly.
func TestTokenSourceMintsReplacementAfterRefreshTime(t *testing.T) {
	minter := &recordingMinter{tokens: []string{tokenExpiringIn(time.Hour), tokenExpiringIn(2 * time.Hour)}}
	source := newTokenSource(minter)

	_, err := source.Token()
	require.NoError(t, err)
	source.refreshAt = time.Now().Add(-time.Second)

	second, err := source.Token()

	require.NoError(t, err)
	assert.Equal(t, minter.tokens[1], second.AccessToken)
	assert.Equal(t, 2, minter.calls)
}

// The mutex is held across the network call so that a burst of parallel platform requests produces
// one token rather than one per goroutine. Were the lock dropped, every caller here would find no
// token and go and mint its own.
func TestTokenSourceMintsOnceForConcurrentCallers(t *testing.T) {
	minter := &concurrentMinter{}
	source := newTokenSource(minter)

	callConcurrently(func(int) {
		_, err := source.Token()
		assert.NoError(t, err)
	})

	assert.Equal(t, int64(1), minter.calls.Load())
}

func TestTokenSourceRefreshesAheadOfExpiry(t *testing.T) {
	source := newTokenSource(&recordingMinter{tokens: []string{tokenExpiringIn(time.Hour)}})

	token, err := source.Token()

	require.NoError(t, err)
	assert.Equal(t, token.Expiry.Add(-refreshMargin), source.refreshAt)
}

// A token that already expires inside the refresh margin would otherwise be due for replacement the
// moment it arrives, turning every request into a trip to the token service.
func TestTokenSourceFloorsRefreshTimeForShortLivedToken(t *testing.T) {
	source := newTokenSource(&recordingMinter{tokens: []string{tokenExpiringIn(10 * time.Second)}})

	_, err := source.Token()

	require.NoError(t, err)
	assert.WithinDuration(t, time.Now().Add(minimumMintInterval), source.refreshAt, time.Second)
}

func TestTokenSourcePropagatesMintingFailure(t *testing.T) {
	mintErr := errors.New("the token service is unreachable")

	_, err := newTokenSource(&recordingMinter{err: mintErr}).Token()

	assert.ErrorIs(t, err, mintErr)
}

// The token due for replacement is still good for another two minutes, which is the whole point of
// replacing it that early. Failing a request while holding a token the platform would accept would
// let a moment of trouble at the token service end an apply that is half way through.
func TestTokenSourceServesUnexpiredTokenWhenMintingFails(t *testing.T) {
	source := newTokenSource(&recordingMinter{tokens: []string{tokenExpiringIn(time.Hour)}})
	held, err := source.Token()
	require.NoError(t, err)
	source.refreshAt = time.Now().Add(-time.Second)
	source.minter = &recordingMinter{err: errors.New("the token service is unreachable")}

	token, err := source.Token()

	require.NoError(t, err)
	assert.Equal(t, held.AccessToken, token.AccessToken)
}

// The counterpart of the test above, differing only in whether the token held has any life left.
// Once it has none there is nothing to serve, and the failure has to reach the caller.
func TestTokenSourcePropagatesMintingFailureOnceTokenHasExpired(t *testing.T) {
	mintErr := errors.New("the token service is unreachable")
	// A token that arrives expired is the shortest way to a source holding an expired one; the
	// refresh margin means it never happens this way outside a test.
	source := newTokenSource(&recordingMinter{tokens: []string{tokenExpiringIn(-time.Second)}})
	_, err := source.Token()
	require.NoError(t, err)
	source.refreshAt = time.Now().Add(-time.Second)
	source.minter = &recordingMinter{err: mintErr}

	_, err = source.Token()

	assert.ErrorIs(t, err, mintErr)
}

// Without this, every request made while the token service is down would wait for a minting attempt
// of its own to run into its timeout, one after another.
func TestTokenSourcePostponesMintingAfterFailure(t *testing.T) {
	source := newTokenSource(&recordingMinter{tokens: []string{tokenExpiringIn(time.Hour)}})
	_, err := source.Token()
	require.NoError(t, err)
	source.refreshAt = time.Now().Add(-time.Second)
	source.minter = &recordingMinter{err: errors.New("the token service is unreachable")}

	_, err = source.Token()

	require.NoError(t, err)
	assert.WithinDuration(t, time.Now().Add(mintRetryInterval), source.refreshAt, time.Second)
}

// Postponing past the expiry would hand out an expired token without so much as trying to replace it.
func TestTokenSourcePostponesMintingNoFurtherThanExpiry(t *testing.T) {
	source := newTokenSource(&recordingMinter{tokens: []string{tokenExpiringIn(mintRetryInterval / 2)}})
	held, err := source.Token()
	require.NoError(t, err)
	source.refreshAt = time.Now().Add(-time.Second)
	source.minter = &recordingMinter{err: errors.New("the token service is unreachable")}

	_, err = source.Token()

	require.NoError(t, err)
	assert.Equal(t, held.Expiry, source.refreshAt)
}

func TestTokenSourceRejectsTokenWithoutExpiry(t *testing.T) {
	_, err := newTokenSource(&recordingMinter{tokens: []string{jwtWithPayload(`{"aud":"dynatrace"}`)}}).Token()

	assert.EqualError(t, err, "the ID token has no `exp` claim, so the provider cannot tell when to obtain a replacement")
}

// The platform can turn down a token this source still considers current - a runner whose clock
// differs from the platform's is enough for that. Minting a replacement on the spot is what lets the
// request be sent again instead of failing the resource it belongs to.
func TestTokenSourceReplacesRejectedToken(t *testing.T) {
	minter := &recordingMinter{tokens: []string{tokenExpiringIn(time.Hour), tokenExpiringIn(2 * time.Hour)}}
	source := newTokenSource(minter)
	rejected, err := source.Token()
	require.NoError(t, err)

	replacement, err := source.Replace(rejected.AccessToken)

	require.NoError(t, err)
	assert.Equal(t, minter.tokens[1], replacement.AccessToken)
}

// Only one token is prepared: a token that stops being accepted is rejected on every request in
// flight at that moment, and minting for each of them would mean a token per parallel resource.
func TestTokenSourceKeepsTokenThatWasAlreadyReplaced(t *testing.T) {
	minter := &recordingMinter{tokens: []string{tokenExpiringIn(time.Hour)}}
	source := newTokenSource(minter)
	held, err := source.Token()
	require.NoError(t, err)

	replacement, err := source.Replace("the token another request had rejected")

	require.NoError(t, err)
	assert.Equal(t, held.AccessToken, replacement.AccessToken)
	assert.Equal(t, 1, minter.calls)
}

// The concurrent counterpart of the test above: a token that stops being accepted is rejected on
// every request in flight at that moment, and all of those rejections have to be answered with the
// same replacement. One minting attempt for the setup and one for the replacement, no more.
func TestTokenSourceReplacesOnceForConcurrentRejections(t *testing.T) {
	minter := &concurrentMinter{}
	source := newTokenSource(minter)
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

// A rejected token must not come back through the fallback that serves an unexpired token when
// minting fails: the platform has already refused it, and its expiry says nothing about that.
func TestTokenSourceDoesNotFallBackOnRejectedToken(t *testing.T) {
	mintErr := errors.New("the token service is unreachable")
	source := newTokenSource(&recordingMinter{tokens: []string{tokenExpiringIn(time.Hour)}})
	rejected, err := source.Token()
	require.NoError(t, err)
	source.minter = &recordingMinter{err: mintErr}
	_, err = source.Replace(rejected.AccessToken)
	require.ErrorIs(t, err, mintErr)

	_, err = source.Token()

	assert.ErrorIs(t, err, mintErr)
}
