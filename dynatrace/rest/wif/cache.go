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

	"golang.org/x/oauth2"
)

// tokenSourceCache holds one token source per configuration for the lifetime of the process.
//
// It is what keeps token minting away from the request path: the provider builds a platform client
// per request, so without this cache each client would come with a fresh token source and every
// single API call would mint a new token.
var (
	tokenSourceCacheMutex sync.Mutex
	tokenSourceCache      = map[string]TokenSource{}
)

// TokenSourceFor returns the token source for the given configuration, creating it on first use.
//
// The returned source is consulted on every request, and it decides on its own when to replace the
// token it holds. Callers must not wrap it in an [oauth2.ReuseTokenSource]: that caches until ten
// seconds before the token's stated expiry and would override the refresh margin the source keeps.
func TokenSourceFor(config Config) (TokenSource, error) {
	if err := Validate(config); err != nil {
		return nil, err
	}

	if len(config.Vendor) == 0 {
		// A pre-minted token carries no expiry here on purpose: the provider cannot obtain a
		// replacement for a token it did not mint, so it sends this one until the platform rejects
		// it. There is also nothing to cache, and a credential is not worth putting into a map key.
		return &staticTokenSource{token: &oauth2.Token{AccessToken: config.StaticToken}}, nil
	}

	tokenSourceCacheMutex.Lock()
	defer tokenSourceCacheMutex.Unlock()

	if source, found := tokenSourceCache[config.CacheKey()]; found {
		return source, nil
	}

	minter, err := newMinter(config, mintingHTTPClient())
	if err != nil {
		// Configuration errors are deliberately not cached. They are deterministic and cheap to
		// reproduce, and a cached one would outlive a fix to the environment.
		return nil, err
	}

	source := &tokenSource{mintContext: context.Background(), minter: minter}
	tokenSourceCache[config.CacheKey()] = source

	return source, nil
}
