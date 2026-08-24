/**
* @license
* Copyright 2025 Dynatrace LLC
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

import (
	"context"
	"errors"
	"net/url"
	"strings"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest/logging"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/provider/version"
	"golang.org/x/oauth2/clientcredentials"

	"github.com/dynatrace/dynatrace-configuration-as-code-core/api/rest"
	"github.com/dynatrace/dynatrace-configuration-as-code-core/clients"
)

var eligiblePlatformRequests = map[string]string{
	"/api/v2/": "/platform/classic/environment-api/v2/",
	"/api/v1/": "/platform/classic/environment-api/v1/",
}

type platform_request request

var NoPlatformCredentialsErr = errors.New("neither oauth credentials nor platform token present")

func CreatePlatformClient(ctx context.Context, platformURL string, credentials *Credentials) (*rest.Client, error) {
	factory := clients.Factory().
		WithPlatformURL(platformURL).
		WithRateLimiter(true).
		WithRetryOptions(defaultRetryOptions).
		WithUserAgent(version.UserAgent())

	if credentials.ContainsPlatformToken() {
		return factory.
			WithHTTPListener(logging.HTTPListener("plat/tok")).
			WithPlatformToken(credentials.Platform.PlatformToken).
			CreatePlatformClient(ctx)
	}

	if credentials.ContainsOAuth() {
		return factory.
			WithHTTPListener(logging.HTTPListener("platform")).
			WithOAuthCredentials(clientcredentials.Config{
				ClientID:     credentials.Platform.ClientID,
				ClientSecret: credentials.Platform.ClientSecret,
				TokenURL:     credentials.Platform.TokenURL,
			}).
			CreatePlatformClient(NewContextWithOAuthRetryClient(ctx))
	}

	return nil, NoPlatformCredentialsErr
}

func (me *platform_request) Finish(optionalTarget ...any) error {
	var target any
	if len(optionalTarget) > 0 {
		target = optionalTarget[0]
	}

	client, err := selectClient(me.client.ClientSet(), me.url)
	if err != nil {
		return err
	}

	meURL := me.url
	for eligiblePlatformURL, replacement := range eligiblePlatformRequests {
		if strings.HasPrefix(meURL, eligiblePlatformURL) {
			meURL = strings.ReplaceAll(meURL, eligiblePlatformURL, replacement)
			// We can exit early here
			// Once we have found a URL eligible to replace with its platform counterpart
			// we shouldn't expect a second one
			break
		}
	}

	fullURL, err := url.Parse(client.BaseURL().String() + meURL)
	if err != nil {
		return err
	}
	return request(*me).HandleResponse(client, fullURL, target)
}

func selectClient(clientSet ClientSet, url string) (*rest.Client, error) {
	if strings.Contains(url, "/api/config/v1/") {
		return clientSet.ClassicPlatformClient()
	}
	return clientSet.PlatformClient()
}
