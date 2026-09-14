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

// Vendor identifies the workload identity provider that issues the OIDC token.
type Vendor = string

// VendorGitHub obtains tokens from the GitHub Actions OIDC token service.
const VendorGitHub Vendor = "github"

type Config struct {
	Vendor      Vendor
	Audience    string
	StaticToken string
	// GitHubTokenRequestURL and GitHubTokenRequestToken are the credentials GitHub injects into a
	// job with id-token: write permission. They may also be supplied via the provider configuration.
	// Consulted only when Vendor is VendorGitHub.
	GitHubTokenRequestURL   string
	GitHubTokenRequestToken string
}

// Configured reports whether any form of Workload Identity Federation was requested.
func (config Config) Configured() bool {
	return len(config.Vendor) > 0 || len(config.StaticToken) > 0
}
