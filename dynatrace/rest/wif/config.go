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

// Package wif obtains OIDC tokens for Workload Identity Federation against the Dynatrace platform
// APIs.
//
// The token is sent to Dynatrace verbatim as a bearer token, exactly like a platform token - there
// is no token exchange. What differs from a platform token is the lifetime: these tokens live
// minutes, carry a fixed expiry and cannot be refreshed, only re-minted. That is why this package
// hands out an [golang.org/x/oauth2.TokenSource] rather than a token.
//
// This package deliberately does not import the surrounding rest package. Beyond avoiding an import
// cycle, that keeps the calls to a token service off the provider's HTTP logging, which writes whole
// response bodies when DYNATRACE_HTTP_RESPONSE is set - and the response of a token service is the
// token itself.
//
// The error messages produced here are user-facing provider diagnostics and therefore name the
// provider's configuration attributes and environment variables.
package wif

// Vendor identifies the workload identity provider that issues the OIDC token.
type Vendor string

// VendorGitHub obtains tokens from the GitHub Actions OIDC token service.
const VendorGitHub Vendor = "github"

// Config describes how a token is obtained. Vendor and StaticToken are mutually exclusive: either a
// vendor mints tokens on demand, or a pre-minted token is supplied.
//
// These fields are limited to what is universal across vendors. Whatever a single vendor needs in
// order to reach its own token service - environment variables, a metadata endpoint, a mounted
// token file - is discovered by that vendor's minter and is not carried here.
type Config struct {
	Vendor      Vendor
	Audience    string
	StaticToken string
}

// Configured reports whether any form of Workload Identity Federation was requested.
func (config Config) Configured() bool {
	return len(config.Vendor) > 0 || len(config.StaticToken) > 0
}

// CacheKey identifies the token source belonging to this configuration.
//
// The audience is part of the key because it is the identity of the minted token: two Dynatrace
// environments addressed with the same audience legitimately share one token.
//
// StaticToken is deliberately not part of the key. A pre-minted token is never cached, so its value
// never becomes a key in a map that lives as long as the process.
func (config Config) CacheKey() string {
	return string(config.Vendor) + "\x00" + config.Audience
}
