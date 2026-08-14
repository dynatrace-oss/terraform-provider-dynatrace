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
	"fmt"
	"net/http"
)

// minter obtains a freshly issued OIDC token from a workload identity provider.
//
// This is the extension point for further vendors: a new vendor is one case in newMinter plus one
// file implementing this interface. The interface is deliberately narrow, because vendors differ in
// how they authenticate to their own token service - GitHub passes credentials to the job through
// environment variables, others use a metadata endpoint, a mounted token file or a signed request.
// Each implementation therefore discovers and validates its own credentials in its constructor,
// rather than being handed a shape that only fits one vendor.
type minter interface {
	mint(ctx context.Context) (string, error)
}

// newMinter builds the minter for the configured vendor. It fails without any network access when
// the vendor's credentials are not available, so that the problem is reported while the provider is
// being configured rather than on the first platform request.
func newMinter(config Config, httpClient *http.Client) (minter, error) {
	switch config.Vendor {
	case VendorGitHub:
		return newGitHubMinter(config.Audience, httpClient)
	default:
		// Validate rejects unsupported vendors before this is reached; this branch keeps the switch
		// total should a caller ever build a minter without validating first.
		return nil, ConfigError{fmt.Sprintf("`%s` is not a supported Workload Identity Federation vendor. The only supported value for `wif_vendor` (`DYNATRACE_WIF_VENDOR`) is `%s`", config.Vendor, VendorGitHub)}
	}
}
