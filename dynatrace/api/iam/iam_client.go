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
	POST(ctx context.Context, url string, payload any, options rest2.RequestOptions, expectedResponseCode int) ([]byte, error)
	PUT(ctx context.Context, url string, payload any, options rest2.RequestOptions, expectedResponseCode int) ([]byte, error)
	PUT_MULTI_RESPONSE(ctx context.Context, url string, payload any, options rest2.RequestOptions, expectedResponseCodes []int) ([]byte, error)
	GET(ctx context.Context, url string, options rest2.RequestOptions, expectedResponseCode int) ([]byte, error)
	DELETE(ctx context.Context, url string, options rest2.RequestOptions, expectedResponseCode int) ([]byte, error)
	DELETE_MULTI_RESPONSE(ctx context.Context, url string, options rest2.RequestOptions, expectedResponseCodes []int) ([]byte, error)
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
				return rest2.RetryIfTooManyRequestsOrServiceUnavailable(resp) || resp.StatusCode == http.StatusBadGateway || resp.StatusCode == http.StatusGatewayTimeout
			},
		}),
	}
	eUrl, _ := url.Parse(a.EndpointURL())
	client := rest2.NewClient(eUrl, httpClient, opts...)
	client.SetHeader("User-Agent", version.UserAgent())
	return &iamClient{a, client}
}

func (me *iamClient) POST(ctx context.Context, url string, payload any, options rest2.RequestOptions, expectedResponseCode int) ([]byte, error) {
	body, err := marshalPayload(payload)
	if err != nil {
		return nil, err
	}
	return handleResponse(me.client.POST(ctx, url, body, options))([]int{expectedResponseCode})
}

func (me *iamClient) PUT(ctx context.Context, url string, payload any, options rest2.RequestOptions, expectedResponseCode int) ([]byte, error) {
	body, err := marshalPayload(payload)
	if err != nil {
		return nil, err
	}
	return handleResponse(me.client.PUT(ctx, url, body, options))([]int{expectedResponseCode})
}

func (me *iamClient) PUT_MULTI_RESPONSE(ctx context.Context, url string, payload any, options rest2.RequestOptions, expectedResponseCodes []int) ([]byte, error) {
	body, err := marshalPayload(payload)
	if err != nil {
		return nil, err
	}
	return handleResponse(me.client.PUT(ctx, url, body, options))(expectedResponseCodes)
}

func (me *iamClient) GET(ctx context.Context, url string, options rest2.RequestOptions, expectedResponseCode int) ([]byte, error) {
	return handleResponse(me.client.GET(ctx, url, options))([]int{expectedResponseCode})
}

func (me *iamClient) DELETE(ctx context.Context, url string, options rest2.RequestOptions, expectedResponseCode int) ([]byte, error) {
	return handleResponse(me.client.DELETE(ctx, url, options))([]int{expectedResponseCode})
}

func (me *iamClient) DELETE_MULTI_RESPONSE(ctx context.Context, url string, options rest2.RequestOptions, expectedResponseCodes []int) ([]byte, error) {
	return handleResponse(me.client.DELETE(ctx, url, options))(expectedResponseCodes)
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

// handleResponse reads the response body and validates the status code against the expected
// codes. It is curried so it can wrap a rest client call directly, e.g.
// handleResponse(client.GET(...))(expectedResponseCodes).
func handleResponse(response *http.Response, err error) func(expectedResponseCodes []int) ([]byte, error) {
	return func(expectedResponseCodes []int) ([]byte, error) {
		if err != nil {
			return nil, err
		}
		defer response.Body.Close()

		responseBytes, err := io.ReadAll(response.Body)
		if err != nil {
			return nil, err
		}

		if slices.Contains(expectedResponseCodes, response.StatusCode) {
			return responseBytes, nil
		}

		return nil, rest.Error{Code: response.StatusCode, Message: fmt.Sprintf("response code %d (expected: %d)", response.StatusCode, expectedResponseCodes)}
	}
}
