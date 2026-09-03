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
	"net/http"
)

// Vendors differ in how they authenticate to their own token service, so each implementation
// discovers and validates its own credentials in its constructor.
type minter interface {
	mint(ctx context.Context) (string, error)
}

// Reaches no network, so missing credentials surface while the provider is configured.
func newMinter(config Config, httpClient *http.Client) (minter, error) {
	switch config.Vendor {
	case VendorGitHub:
		return newGitHubMinter(config.Audience, httpClient)
	default:
		return nil, unsupportedVendorError(config.Vendor)
	}
}
