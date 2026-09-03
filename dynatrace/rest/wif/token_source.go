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
	"time"

	"golang.org/x/oauth2"
)

// refreshMargin is how long before its stated expiry a token is replaced: wide enough that a token handed
// out just before it is due cannot expire in flight, and wide enough to absorb a clock differing from
// the platform's.
const refreshMargin = 2 * time.Minute

// TokenSourceFor returns a token source for the given configuration. Every mint is bounded by ctx, so
// ctx has to outlive the requests the source ends up authenticating.
func TokenSourceFor(ctx context.Context, config Config) (oauth2.TokenSource, error) {
	if err := config.Validate(); err != nil {
		return nil, err
	}

	if len(config.StaticToken) > 0 {
		// No expiry: reporting one would make oauth2 discard a token the provider cannot re-mint.
		return oauth2.StaticTokenSource(&oauth2.Token{AccessToken: config.StaticToken}), nil
	}

	minter, err := newMinter(config, mintingHTTPClient())
	if err != nil {
		return nil, err
	}

	return reusingTokenSource(ctx, minter), nil
}

func reusingTokenSource(ctx context.Context, minter minter) oauth2.TokenSource {
	return oauth2.ReuseTokenSourceWithExpiry(nil, &mintingTokenSource{mintContext: ctx, minter: minter}, refreshMargin)
}

// mintingTokenSource mints on every call; when that is necessary is decided by the source wrapping it.
type mintingTokenSource struct {
	// Held rather than passed: oauth2.TokenSource.Token takes no context, and this one must outlive
	// every request that uses the source.
	mintContext context.Context
	minter      minter
}

func (source *mintingTokenSource) Token() (*oauth2.Token, error) {
	rawToken, err := source.minter.mint(source.mintContext)
	if err != nil {
		return nil, err
	}

	expiry, err := expiryOf(rawToken)
	if err != nil {
		return nil, err
	}

	// An empty TokenType is what makes oauth2 send this as a bearer token.
	return &oauth2.Token{AccessToken: rawToken, Expiry: expiry}, nil
}
