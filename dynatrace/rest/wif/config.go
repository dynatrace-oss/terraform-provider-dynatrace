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
	"errors"
	"net/http"
)

// vendorConfig is what one workload identity provider needs in order to issue an OIDC token. Every
// value in it except the audience is discovered from the environment the provider runs in.
type vendorConfig interface {
	createMinter(httpClient *http.Client) (minter, error)
}

var errNoVendorDetected = errors.New("no supported workload identity provider was detected: the only supported one is GitHub Actions, which injects the required `ACTIONS_ID_TOKEN_REQUEST_*` variables only into a job with `permissions: { id-token: write }`")

// inferVendorConfig decides which workload identity provider mints the OIDC token by looking at the
// credentials its token service leaves in the environment. The provider configuration names none of
// them; it only states the audience the token has to be valid for.
func inferVendorConfig(audience string) (vendorConfig, error) {
	if config, detected := inferGitHubConfig(audience); detected {
		return config, nil
	}

	return nil, errNoVendorDetected
}
