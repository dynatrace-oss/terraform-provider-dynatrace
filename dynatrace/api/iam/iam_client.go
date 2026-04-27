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
	"fmt"
	"io"
	"net/http"
	"net/url"
	"slices"
	"time"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest/logging"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/provider/version"
	"golang.org/x/oauth2/clientcredentials"

	"github.com/dynatrace/dynatrace-configuration-as-code-core/api/auth"
	rest2 "github.com/dynatrace/dynatrace-configuration-as-code-core/api/rest"
)

type IAMClient interface {
	POST(ctx context.Context, url string, payload any, expectedResponseCode int, forceNewBearer bool) ([]byte, error)
	PUT(ctx context.Context, url string, payload any, expectedResponseCode int, forceNewBearer bool) ([]byte, error)
	PUT_MULTI_RESPONSE(ctx context.Context, url string, payload any, expectedResponseCodes []int, forceNewBearer bool) ([]byte, error)
	GET(ctx context.Context, url string, expectedResponseCode int, forceNewBearer bool) ([]byte, error)
	DELETE(ctx context.Context, url string, expectedResponseCode int, forceNewBearer bool) ([]byte, error)
	DELETE_MULTI_RESPONSE(ctx context.Context, url string, expectedResponseCodes []int, forceNewBearer bool) ([]byte, error)
}

type iamClient struct {
	auth   Authenticator
	client *rest2.Client
}

func NewIAMClient(ctx context.Context, a Authenticator) IAMClient {
	oauthConfig := clientcredentials.Config{
		ClientID:     a.ClientID(),
		ClientSecret: a.ClientSecret(),
		TokenURL:     a.TokenURL(),
	}
	httpClient := auth.NewOAuthClient(rest.NewContextWithOAuthRetryClient(ctx), &oauthConfig)

	opts := []rest2.Option{
		rest2.WithHTTPListener(logging.HTTPListener("iam")),
		rest2.WithRateLimiter(),
		rest2.WithRetryOptions(&rest2.RetryOptions{
			DelayAfterRetry: time.Second * 60,
			MaxRetries:      10,
			ShouldRetryFunc: func(resp *http.Response) bool {
				return resp.StatusCode == 429 || resp.StatusCode == 504
			},
		}),
	}
	eUrl, _ := url.Parse(a.EndpointURL())
	client := rest2.NewClient(eUrl, httpClient, opts...)
	return &iamClient{a, client}
}

func (me *iamClient) POST(ctx context.Context, url string, payload any, expectedResponseCode int, forceNewBearer bool) ([]byte, error) {
	return me.request(ctx, url, http.MethodPost, []int{expectedResponseCode}, forceNewBearer, 0, payload, map[string]string{"Content-Type": "application/json"})
}

func (me *iamClient) PUT(ctx context.Context, url string, payload any, expectedResponseCode int, forceNewBearer bool) ([]byte, error) {
	return me.request(ctx, url, http.MethodPut, []int{expectedResponseCode}, forceNewBearer, 0, payload, map[string]string{"Content-Type": "application/json"})
}

func (me *iamClient) PUT_MULTI_RESPONSE(ctx context.Context, url string, payload any, expectedResponseCodes []int, forceNewBearer bool) ([]byte, error) {
	return me.request(ctx, url, http.MethodPut, expectedResponseCodes, forceNewBearer, 0, payload, map[string]string{"Content-Type": "application/json"})
}

func (me *iamClient) GET(ctx context.Context, url string, expectedResponseCode int, forceNewBearer bool) ([]byte, error) {
	return me.request(ctx, url, http.MethodGet, []int{expectedResponseCode}, forceNewBearer, 0, nil, nil)
}

func (me *iamClient) DELETE(ctx context.Context, url string, expectedResponseCode int, forceNewBearer bool) ([]byte, error) {
	return me.request(ctx, url, http.MethodDelete, []int{expectedResponseCode}, forceNewBearer, 0, nil, nil)
}

func (me *iamClient) DELETE_MULTI_RESPONSE(ctx context.Context, url string, expectedResponseCodes []int, forceNewBearer bool) ([]byte, error) {
	return me.request(ctx, url, http.MethodDelete, expectedResponseCodes, forceNewBearer, 0, nil, nil)
}

func (me *iamClient) request(ctx context.Context, u string, method string, expectedResponseCodes []int, forceNewBearer bool, forceNewBearerRetryCount int, payload any, headers map[string]string) ([]byte, error) {
	var err error
	var httpRequest *http.Request
	var httpResponse *http.Response
	var responseBytes []byte
	var requestBody []byte

	if requestBody, err = json.Marshal(payload); err != nil {
		return nil, err
	}

	var body io.Reader

	if payload != nil {
		body = bytes.NewReader(requestBody)
	}

	if httpRequest, err = http.NewRequest(method, u, body); err != nil {
		return nil, err
	}

	for k, v := range headers {
		httpRequest.Header.Add(k, v)
	}

	httpRequest.Header.Set("User-Agent", version.UserAgent())

	if httpResponse, err = me.client.Do(httpRequest); err != nil {
		return nil, err
	}

	if responseBytes, err = io.ReadAll(httpResponse.Body); err != nil {
		return nil, err
	}

	if slices.Contains(expectedResponseCodes, httpResponse.StatusCode) {
		return responseBytes, nil
	}

	return nil, rest.Error{Code: httpResponse.StatusCode, Message: fmt.Sprintf("response code %d (expected: %d)", httpResponse.StatusCode, expectedResponseCodes)}
}
