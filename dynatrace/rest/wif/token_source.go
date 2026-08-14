/**
* @license
* Copyright 2026 Dynatrace LLC
*
* Licensed under the Apache License, Version 2.0 (the "License");
* you may not use this file except in compliance with the License.
* You may obtain a copy of the License at
*
*     http://www.apache.org/licenses/LICENSE-2.0
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
	"sync"
	"time"

	"golang.org/x/oauth2"
)

const (
	// refreshMargin is how long before its expiry a token is replaced. These tokens cannot be
	// refreshed, only re-minted, so the margin has to be wide enough that a request never picks up a
	// token which expires while it is still in flight.
	refreshMargin = 2 * time.Minute

	// minimumMintInterval bounds re-minting when a freshly minted token already expires within the
	// refresh margin, so that a misbehaving token service cannot be asked for a token per request.
	minimumMintInterval = 30 * time.Second

	// mintRetryInterval is how long minting is left alone after it failed while a token that still
	// works is being handed out. Without it, an outage of the token service would put a minting
	// attempt - and its timeout - in front of every single request for as long as it lasts.
	mintRetryInterval = 5 * time.Second
)

// TokenSource hands out the tokens platform requests are authenticated with.
//
// It extends [oauth2.TokenSource] with the opposite direction of information: the platform can
// reject a token that this source still considers current, and Replace is how that answer gets back
// here. A Terraform run easily outlives a federated token, so the ability to obtain a replacement
// mid-run is what keeps a long apply going.
type TokenSource interface {
	oauth2.TokenSource

	// Replace hands out a token to be used instead of one the platform rejected. A source that
	// cannot obtain a replacement returns the rejected token unchanged, which tells the caller that
	// sending the request again is pointless.
	Replace(rejected string) (*oauth2.Token, error)
}

// tokenSource hands out OIDC tokens, obtaining a replacement shortly before the current one expires.
type tokenSource struct {
	// mintContext is captured at construction because oauth2.TokenSource.Token takes no context,
	// while the context the provider is configured with does not outlive the configure call. The
	// minting HTTP client carries its own timeout instead.
	mintContext context.Context
	minter      minter

	mutex     sync.Mutex
	token     *oauth2.Token
	refreshAt time.Time
}

func (source *tokenSource) Token() (*oauth2.Token, error) {
	// The lock is deliberately held across the network call. That makes minting single flight, so a
	// burst of parallel platform requests produces one token rather than one per goroutine.
	source.mutex.Lock()
	defer source.mutex.Unlock()

	if source.token != nil && time.Now().Before(source.refreshAt) {
		return source.token, nil
	}

	token, err := source.mintAndHold()
	if err == nil {
		return token, nil
	}

	// The token due for replacement has not expired yet - having something left to fall back on here
	// is what the refresh margin is for. Handing it out is better than failing a request over a token
	// service that may well answer again a second later, and if the platform rejects it after all,
	// Replace gets another chance at a replacement.
	if fallback := source.unexpiredToken(); fallback != nil {
		source.postponeMinting()
		return fallback, nil
	}

	return nil, err
}

// Replace hands out a token to be used instead of the one the platform rejected, minting one unless
// the token held has been replaced already.
func (source *tokenSource) Replace(rejected string) (*oauth2.Token, error) {
	source.mutex.Lock()
	defer source.mutex.Unlock()

	// The token held may already be a different one: with resources being applied in parallel, a
	// token that stops being accepted is rejected on many requests at once, and one replacement
	// serves all of them.
	if source.token != nil && source.token.AccessToken != rejected {
		return source.token, nil
	}

	// The rejected token is dropped before minting rather than after: it is of no use even if
	// minting fails, and leaving it in place would let the fallback above hand it out again.
	source.token = nil
	source.refreshAt = time.Time{}

	return source.mintAndHold()
}

// mintAndHold obtains a token and takes it into use. The mutex must be held by the caller.
func (source *tokenSource) mintAndHold() (*oauth2.Token, error) {
	rawToken, err := source.minter.mint(source.mintContext)
	if err != nil {
		return nil, err
	}

	expiry, err := expiryOf(rawToken)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	refreshAt := expiry.Add(-refreshMargin)
	if !refreshAt.After(now) {
		refreshAt = now.Add(minimumMintInterval)
	}

	// Expiry stays the expiry the token actually claims; when to replace it is tracked separately,
	// so that the token reports what it is rather than when this source intends to act on it.
	//
	// Leaving TokenType empty is what makes oauth2 send the token as a bearer token, which is the
	// wire format the Dynatrace platform APIs expect - the same one a platform token uses.
	source.token = &oauth2.Token{AccessToken: rawToken, Expiry: expiry}
	source.refreshAt = refreshAt

	return source.token, nil
}

// unexpiredToken returns the token currently held for as long as it can still be used, and nil once
// it has expired or when none is held. The mutex must be held by the caller.
func (source *tokenSource) unexpiredToken() *oauth2.Token {
	if source.token == nil || !time.Now().Before(source.token.Expiry) {
		return nil
	}
	return source.token
}

// postponeMinting schedules the next minting attempt a short while ahead. It never schedules it past
// the expiry of the token held, so that a token is not handed out beyond its expiry without at least
// having tried to replace it. The mutex must be held by the caller.
func (source *tokenSource) postponeMinting() {
	retryAt := time.Now().Add(mintRetryInterval)
	if retryAt.After(source.token.Expiry) {
		retryAt = source.token.Expiry
	}
	source.refreshAt = retryAt
}

// staticTokenSource hands out a token that was supplied ready-made.
type staticTokenSource struct {
	token *oauth2.Token
}

func (source *staticTokenSource) Token() (*oauth2.Token, error) {
	return source.token, nil
}

// Replace returns the token held, which is the one that was rejected: the provider cannot obtain a
// replacement for a token it did not mint itself, so there is nothing to send the request again with.
func (source *staticTokenSource) Replace(string) (*oauth2.Token, error) {
	return source.token, nil
}
