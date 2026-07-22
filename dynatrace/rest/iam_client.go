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

package rest

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest/logging"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/provider/version"
	"golang.org/x/oauth2/clientcredentials"

	"github.com/dynatrace/dynatrace-configuration-as-code-core/api"
	rest2 "github.com/dynatrace/dynatrace-configuration-as-code-core/api/rest"
	"github.com/dynatrace/dynatrace-configuration-as-code-core/clients"
)

// IAMClient performs HTTP requests against the IAM REST API. Each method returns an api.Response
// for any successful (2xx) response. Any non-2xx response is returned as an api.APIError, which
// call sites can inspect (e.g. via errors.AsType) when a specific non-2xx code is expected.
type IAMClient interface {
	POST(ctx context.Context, url string, payload any, options rest2.RequestOptions) (api.Response, error)
	PUT(ctx context.Context, url string, payload any, options rest2.RequestOptions) (api.Response, error)
	GET(ctx context.Context, url string, options rest2.RequestOptions) (api.Response, error)
	DELETE(ctx context.Context, url string, options rest2.RequestOptions) (api.Response, error)
	// AccountID returns the account ID associated with the client. It is a constant fixed when the
	// client is created.
	AccountID() string
}

type iamClient struct {
	client    *rest2.Client
	accountID string
}

func NewIAMClient(ctx context.Context, credentials *Credentials) (IAMClient, error) {
	client, err := clients.Factory().
		WithHTTPListener(logging.HTTPListener("iam")).
		WithOAuthCredentials(clientcredentials.Config{
			ClientID:     credentials.IAM.ClientID,
			ClientSecret: credentials.IAM.ClientSecret,
			TokenURL:     credentials.IAM.TokenURL,
		}).
		WithUserAgent(version.UserAgent()).
		WithRateLimiter(true).
		WithRetryOptions(&rest2.RetryOptions{
			DelayAfterRetry: time.Second * 60,
			MaxRetries:      10,
			ShouldRetryFunc: func(resp *http.Response) bool {
				return rest2.RetryIfTooManyRequestsOrServiceUnavailable(resp) || resp.StatusCode == http.StatusBadGateway || resp.StatusCode == http.StatusGatewayTimeout
			},
		}).
		WithAccountURL(credentials.IAM.EndpointURL).
		AccountRestClient(NewContextWithOAuthRetryClient(ctx))
	if err != nil {
		return nil, err
	}

	return &iamClient{client: client, accountID: credentials.IAM.AccountID}, nil
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
