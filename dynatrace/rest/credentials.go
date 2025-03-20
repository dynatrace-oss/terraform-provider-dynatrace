/**
* @license
* Copyright 2020 Dynatrace LLC
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

package rest

const (
	ProdTokenURL   = "https://sso.dynatrace.com/sso/oauth2/token"
	SprintTokenURL = "https://sso-sprint.dynatracelabs.com/sso/oauth2/token"
	DevTokenURL    = "https://sso-dev.dynatracelabs.com/sso/oauth2/token"

	ProdIAMEndpointURL   = "https://api.dynatrace.com"
	SprintIAMEndpointURL = "https://api-hardening.internal.dynatracelabs.com"
	DevIAMEndpointURL    = "https://api-dev.internal.dynatracelabs.com"
)

const TestCaseEnvURL = "go-test"

type OAuthCredentials struct {
	ClientID       string
	ClientSecret   string
	TokenURL       string
	EnvironmentURL string
	PlatformToken  string
}

type Credentials struct {
	URL   string
	Token string
	IAM   struct {
		ClientID     string
		AccountID    string
		ClientSecret string
		TokenURL     string
		EndpointURL  string
	}
	OAuth   OAuthCredentials
	Cluster struct {
		URL   string
		Token string
	}
}

func (c *Credentials) ContainsOAuth() bool {
	return len(c.OAuth.ClientID) > 0 && len(c.OAuth.ClientSecret) > 0
}

func (c *Credentials) ContainsPlatformToken() bool {
	return len(c.OAuth.PlatformToken) > 0
}

func (c *Credentials) ContainsAPIToken() bool {
	return len(c.Token) > 0
}

func (c *Credentials) ContainsClusterURL() bool {
	return len(c.Cluster.URL) > 0
}

func (c *Credentials) ContainsClusterToken() bool {
	return len(c.Cluster.Token) > 0
}

func (c *Credentials) ContainsOAuthOrPlatformToken() bool {
	return c.ContainsOAuth() || c.ContainsPlatformToken()
}
