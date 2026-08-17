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

// A platform request builds its own client, and with it a token source, inside every single call, so
// without this cache every API call would mint a token of its own. That is the only reason it exists
// - the day those clients outlive the request they were built for, it can go.
var (
	tokenSourceCacheMutex sync.Mutex
	tokenSourceCache      = map[string]TokenSource{}
)

// TokenSourceFor returns the token source for the given configuration, creating it on first use.
//
// Callers must not wrap the result in an [oauth2.ReuseTokenSource]: that caches until ten seconds
// before the token's stated expiry and would override the refresh margin the source keeps.
func TokenSourceFor(config Config) (TokenSource, error) {
	if err := Validate(config); err != nil {
		return nil, err
	}

	if len(config.Vendor) == 0 {
		// No expiry on purpose: a token the provider did not mint is sent until the platform rejects
		// it, and reporting an expiry would make oauth2 discard it instead.
		return &staticTokenSource{token: &oauth2.Token{AccessToken: config.StaticToken}}, nil
	}

	tokenSourceCacheMutex.Lock()
	defer tokenSourceCacheMutex.Unlock()

	if source, found := tokenSourceCache[config.CacheKey()]; found {
		return source, nil
	}

	minter, err := newMinter(config, mintingHTTPClient())
	if err != nil {
		// Not cached: a cached one would outlive a fix to the environment.
		return nil, err
	}

	source := &tokenSource{mintContext: context.Background(), minter: minter}
	tokenSourceCache[config.CacheKey()] = source

	return source, nil
}
