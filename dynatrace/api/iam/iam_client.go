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

package iam

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest/logging"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/provider/version"
	"golang.org/x/oauth2/clientcredentials"

	"github.com/dynatrace/dynatrace-configuration-as-code-core/api"
	"github.com/dynatrace/dynatrace-configuration-as-code-core/api/auth"
	rest2 "github.com/dynatrace/dynatrace-configuration-as-code-core/api/rest"
)

// IAMClient performs HTTP requests against the IAM REST API. Each method returns an api.Response
// for any successful (2xx) response. Any non-2xx response is returned as an api.APIError, which
// call sites can inspect (e.g. via errors.AsType) when a specific non-2xx code is expected.
type IAMClient interface {
	POST(ctx context.Context, url string, payload any, options rest2.RequestOptions) (api.Response, error)
	PUT(ctx context.Context, url string, payload any, options rest2.RequestOptions) (api.Response, error)
	GET(ctx context.Context, url string, options rest2.RequestOptions) (api.Response, error)
	DELETE(ctx context.Context, url string, options rest2.RequestOptions) (api.Response, error)
	// AccountID returns the account the client is scoped to. The provider only supports a single
	// account, so it is held here rather than being repeated by every service.
	AccountID() string
}

type iamClient struct {
	accountID string
	client    *rest2.Client
}

// NewIAMClientForCredentials builds an IAM client directly from the IAM section of the provider credentials.
func NewIAMClientForCredentials(ctx context.Context, credentials *rest.Credentials) IAMClient {
	return newIAMClient(ctx, credentials.IAM.ClientID, credentials.IAM.AccountID, credentials.IAM.ClientSecret, credentials.IAM.TokenURL, credentials.IAM.EndpointURL)
}

func newIAMClient(ctx context.Context, clientID string, accountID string, clientSecret string, tokenURL string, endpointURL string) IAMClient {
	oauthConfig := clientcredentials.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		TokenURL:     tokenURL,
	}
	httpClient := auth.NewOAuthClient(rest.NewContextWithOAuthRetryClient(ctx), &oauthConfig)

	opts := []rest2.Option{
		rest2.WithHTTPListener(logging.HTTPListener("iam")),
		rest2.WithRateLimiter(),
		rest2.WithRetryOptions(&rest2.RetryOptions{
			DelayAfterRetry: time.Second * 60,
			MaxRetries:      10,
			ShouldRetryFunc: func(resp *http.Response) bool {
				return rest2.RetryIfTooManyRequestsOrServiceUnavailable(resp) || resp.StatusCode == http.StatusBadGateway || resp.StatusCode == http.StatusGatewayTimeout
			},
		}),
	}
	eUrl, _ := url.Parse(endpointURL)
	client := rest2.NewClient(eUrl, httpClient, opts...)
	client.SetHeader("User-Agent", version.UserAgent())
	return &iamClient{accountID: accountID, client: client}
}

func (me *iamClient) AccountID() string {
	return me.accountID
}

func (me *iamClient) POST(ctx context.Context, url string, payload any, options rest2.RequestOptions) (api.Response, error) {
	body, err := marshalPayload(payload)
	if err != nil {
		return api.Response{}, err
	}
	httpResponse, err := me.client.POST(ctx, url, body, options)
	if err != nil {
		return api.Response{}, err
	}
	return api.NewResponseFromHTTPResponse(httpResponse)

}

func (me *iamClient) PUT(ctx context.Context, url string, payload any, options rest2.RequestOptions) (api.Response, error) {
	body, err := marshalPayload(payload)
	if err != nil {
		return api.Response{}, err
	}
	httpResponse, err := me.client.PUT(ctx, url, body, options)
	if err != nil {
		return api.Response{}, err
	}
	return api.NewResponseFromHTTPResponse(httpResponse)
}

func (me *iamClient) GET(ctx context.Context, url string, options rest2.RequestOptions) (api.Response, error) {
	httpResponse, err := me.client.GET(ctx, url, options)
	if err != nil {
		return api.Response{}, err
	}
	return api.NewResponseFromHTTPResponse(httpResponse)
}

func (me *iamClient) DELETE(ctx context.Context, url string, options rest2.RequestOptions) (api.Response, error) {
	httpResponse, err := me.client.DELETE(ctx, url, options)
	if err != nil {
		return api.Response{}, err
	}
	return api.NewResponseFromHTTPResponse(httpResponse)
}

// marshalPayload serializes a request payload to a JSON body. A nil payload yields a nil body.
func marshalPayload(payload any) (io.Reader, error) {
	if payload == nil {
		return nil, nil
	}
	requestBody, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(requestBody), nil
}
