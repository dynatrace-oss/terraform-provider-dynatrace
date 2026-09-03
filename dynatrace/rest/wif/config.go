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

var (
	errNotConfigured        = errors.New("no Workload Identity Federation has been configured")
	errVendorAndStaticToken = errors.New("a vendor and a pre-minted OIDC token are mutually exclusive, specify only one of them")
	errStaticTokenNotAJWT   = errors.New("the pre-minted OIDC token is not a JWT: expected three dot-separated segments")
	errNoAudience           = errors.New("no audience has been specified")
)

// Vendor identifies the workload identity provider that issues the OIDC token.
type Vendor = string

// VendorGitHub obtains tokens from the GitHub Actions OIDC token service.
const VendorGitHub Vendor = "github"

type Config struct {
	Vendor      Vendor
	Audience    string
	StaticToken string
}

// Configured reports whether any form of Workload Identity Federation was requested.
func (config Config) Configured() bool {
	return len(config.Vendor) > 0 || len(config.StaticToken) > 0
}

// Validate reports whether the configuration describes a usable setup, as far as that can be judged
// without contacting a token service. Callers for which Workload Identity Federation is optional have
// to guard with [Config.Configured] - only they can decide whether its absence is a problem.
func (config Config) Validate() error {
	if !config.Configured() {
		return errNotConfigured
	}

	if len(config.Vendor) > 0 && len(config.StaticToken) > 0 {
		return errVendorAndStaticToken
	}

	if len(config.StaticToken) > 0 {
		if strings.Count(config.StaticToken, ".") != jwtSegmentCount-1 {
			return errStaticTokenNotAJWT
		}
		return nil
	}

	if config.Vendor != VendorGitHub {
		return unsupportedVendorError(config.Vendor)
	}

	if len(config.Audience) == 0 {
		return errNoAudience
	}

	return nil
}

func unsupportedVendorError(vendor Vendor) error {
	return fmt.Errorf("`%s` is not a supported vendor, the only supported one is `%s`", vendor, VendorGitHub)
}
