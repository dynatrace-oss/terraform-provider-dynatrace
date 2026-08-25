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
	"fmt"
	"strings"
)

// Validate reports the configuration problems that can be judged without contacting a token service.
// Requesting no Workload Identity Federation at all is not one of them: only a caller that requires
// platform credentials can decide whether its absence is a problem.
func Validate(config Config) error {
	if !config.Configured() {
		return nil
	}

	if len(config.Vendor) > 0 && len(config.StaticToken) > 0 {
		return errors.New("a Workload Identity Federation vendor and a pre-minted OIDC token have both been specified. These options are mutually exclusive. Unset either `wif_vendor` (`DYNATRACE_WIF_VENDOR`) or `wif_oidc_token` (`DYNATRACE_WIF_OIDC_TOKEN`)")
	}

	if len(config.StaticToken) > 0 {
		if strings.Count(config.StaticToken, ".") != jwtSegmentCount-1 {
			return errors.New("the value of `wif_oidc_token` (`DYNATRACE_WIF_OIDC_TOKEN`) is not a JWT: expected three dot-separated segments")
		}
		return nil
	}

	if config.Vendor != VendorGitHub {
		return fmt.Errorf("`%s` is not a supported Workload Identity Federation vendor. The only supported value for `wif_vendor` (`DYNATRACE_WIF_VENDOR`) is `%s`", config.Vendor, VendorGitHub)
	}

	if len(config.Audience) == 0 {
		return errors.New("no audience has been specified for Workload Identity Federation. Use either the configuration attribute `wif_audience` or the environment variable `DYNATRACE_WIF_AUDIENCE` for that")
	}

	return nil
}
