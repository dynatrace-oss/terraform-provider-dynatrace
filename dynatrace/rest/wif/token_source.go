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
	// Wide enough that a token handed out just before it is due can still not expire in flight.
	refreshMargin = 2 * time.Minute

	// Floor for a token that already expires within the refresh margin when it arrives, so that a
	// misbehaving token service cannot be asked for a token per request.
	minimumMintInterval = 30 * time.Second

	// How long minting is left alone after it failed while a usable token is still being served.
	mintRetryInterval = 5 * time.Second
)

// TokenSource hands out the tokens platform requests are authenticated with.
type TokenSource interface {
	oauth2.TokenSource

	// Replace hands out a token to be used instead of one the platform rejected. A source that
	// cannot obtain a replacement returns the rejected token unchanged, which tells the caller that
	// sending the request again is pointless.
	Replace(rejected string) (*oauth2.Token, error)
}

type tokenSource struct {
	// Held rather than passed, because oauth2.TokenSource.Token takes no context. Must outlive every
	// request that uses this source, so it bounds the provider's lifetime rather than one request.
	mintContext context.Context
	minter      minter

	mutex     sync.Mutex
	token     *oauth2.Token
	refreshAt time.Time
}

func (source *tokenSource) Token() (*oauth2.Token, error) {
	// Held across the network call, which makes minting single flight: a burst of parallel requests
	// produces one token rather than one per goroutine.
	source.mutex.Lock()
	defer source.mutex.Unlock()

	if source.token != nil && time.Now().Before(source.refreshAt) {
		return source.token, nil
	}

	token, err := source.mintAndHold()
	if err == nil {
		return token, nil
	}

	// Serving a token the platform still accepts beats failing the request over a token service that
	// may well answer again a second later. Having one left to serve here is what the margin is for.
	if fallback := source.unexpiredToken(); fallback != nil {
		source.postponeMinting()
		return fallback, nil
	}

	return nil, err
}

func (source *tokenSource) Replace(rejected string) (*oauth2.Token, error) {
	source.mutex.Lock()
	defer source.mutex.Unlock()

	// One replacement serves every rejection of the same token, and with resources applied in
	// parallel they arrive together.
	if source.token != nil && source.token.AccessToken != rejected {
		return source.token, nil
	}

	// Dropped before minting rather than after, or a failed mint would leave it for Token to serve.
	source.token = nil
	source.refreshAt = time.Time{}

	return source.mintAndHold()
}

// The mutex must be held by the caller.
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

	// An empty TokenType is what makes oauth2 send this as a bearer token.
	source.token = &oauth2.Token{AccessToken: rawToken, Expiry: expiry}
	source.refreshAt = refreshAt

	return source.token, nil
}

// The mutex must be held by the caller.
func (source *tokenSource) unexpiredToken() *oauth2.Token {
	if source.token == nil || !time.Now().Before(source.token.Expiry) {
		return nil
	}
	return source.token
}

// The mutex must be held by the caller.
func (source *tokenSource) postponeMinting() {
	retryAt := time.Now().Add(mintRetryInterval)
	// Never past the expiry, or a token would be served beyond it without an attempt to replace it.
	if retryAt.After(source.token.Expiry) {
		retryAt = source.token.Expiry
	}
	source.refreshAt = retryAt
}

type staticTokenSource struct {
	token *oauth2.Token
}

func (source *staticTokenSource) Token() (*oauth2.Token, error) {
	return source.token, nil
}

// Replace returns the rejected token itself: the provider cannot replace a token it did not mint.
func (source *staticTokenSource) Replace(string) (*oauth2.Token, error) {
	return source.token, nil
}
