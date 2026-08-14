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

import "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest/wif"

const TestCaseEnvURL = "go-test"

type IAMCredentials struct {
	ClientID     string
	AccountID    string
	ClientSecret string
	TokenURL     string
	EndpointURL  string
}

type PlatformCredentials struct {
	ClientID       string
	ClientSecret   string
	TokenURL       string
	EnvironmentURL string
	PlatformToken  string
	WIF            wif.Config
}

type ClusterCredentials struct {
	URL   string
	Token string
}

type Credentials struct {
	URL      string
	Token    string
	IAM      IAMCredentials
	Platform PlatformCredentials
	Cluster  ClusterCredentials
}

func (c *Credentials) ContainsOAuth() bool {
	return len(c.Platform.ClientID) > 0 && len(c.Platform.ClientSecret) > 0
}

func (c *Credentials) ContainsPlatformToken() bool {
	return len(c.Platform.PlatformToken) > 0
}

func (c *Credentials) ContainsAPIToken() bool {
	return len(c.Token) > 0
}

// ContainsWIF reports whether Workload Identity Federation was requested. It says nothing about
// whether obtaining a token will succeed: a configuration that asks for it but cannot reach its
// token service has to fail with that error, rather than silently falling back to another credential
// and authenticating as somebody else.
func (c *Credentials) ContainsWIF() bool {
	return c.Platform.WIF.Configured()
}

// ContainsPlatformCredentials reports whether any credential usable against the platform APIs is
// present. It answers whether a request may take the platform path, not which credential that
// request will end up using - that order is decided in CreatePlatformClient.
func (c *Credentials) ContainsPlatformCredentials() bool {
	return c.ContainsWIF() || c.ContainsOAuth() || c.ContainsPlatformToken()
}
