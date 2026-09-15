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
	Vendor   Vendor
	Audience string
	GitHub   GitHubConfig
}

// GitHubConfig holds the GitHub Actions OIDC token service credentials. GitHub injects these into
// a job with id-token: write permission; they may also be supplied via the provider configuration.
type GitHubConfig struct {
	TokenRequestURL   string
	TokenRequestToken string
}

// Configured reports whether Workload Identity Federation was requested.
func (config Config) Configured() bool {
	return len(config.Vendor) > 0
}
