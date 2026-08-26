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
	"net/http"
	"net/url"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest/logging"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/provider/version"
	"github.com/google/uuid"

	"github.com/dynatrace/dynatrace-configuration-as-code-core/api/rest"
	"github.com/dynatrace/dynatrace-configuration-as-code-core/clients"
)

func APITokenClient(clientSet ClientSet) Client {
	return &api_token_client{clientSet: clientSet}
}

type api_token_client struct {
	clientSet ClientSet
}

func (me *api_token_client) ClientSet() ClientSet {
	return me.clientSet
}

func (me *api_token_client) Get(ctx context.Context, url string, expectedStatusCodes ...int) Request {
	req := &api_token_request{id: uuid.NewString(), ctx: ctx, client: me, url: url, method: http.MethodGet}
	if len(expectedStatusCodes) > 0 {
		req.expect = statuscodes(expectedStatusCodes)
	}
	return req
}

func (me *api_token_client) Post(ctx context.Context, url string, payload any, expectedStatusCodes ...int) Request {
	req := &api_token_request{id: uuid.NewString(), ctx: ctx, client: me, url: url, method: http.MethodPost, payload: payload}
	if len(expectedStatusCodes) > 0 {
		req.expect = statuscodes(expectedStatusCodes)
	}
	return req
}

func (me *api_token_client) Put(ctx context.Context, url string, payload any, expectedStatusCodes ...int) Request {
	req := &api_token_request{id: uuid.NewString(), ctx: ctx, client: me, url: url, method: http.MethodPut, payload: payload}
	if len(expectedStatusCodes) > 0 {
		req.expect = statuscodes(expectedStatusCodes)
	}
	return req
}

func (me *api_token_client) Delete(ctx context.Context, url string, expectedStatusCodes ...int) Request {
	req := &api_token_request{id: uuid.NewString(), ctx: ctx, client: me, url: url, method: http.MethodDelete}
	if len(expectedStatusCodes) > 0 {
		req.expect = statuscodes(expectedStatusCodes)
	}
	return req
}

func (me *api_token_request) Finish(optionalTarget ...any) error {
	client, err := me.client.ClientSet().APITokenClient()
	if err != nil {
		return err
	}

	pathURL, err := url.Parse(me.url)
	if err != nil {
		return err
	}

	var target any
	if len(optionalTarget) > 0 {
		target = optionalTarget[0]
	}
	return request(*me).HandleResponse(client, pathURL, target)
}

type api_token_request request

func (me *api_token_request) Expect(codes ...int) Request {
	me.expect = statuscodes(codes)
	return me
}

func (me *api_token_request) OnResponse(onResponse func(resp *http.Response)) Request {
	me.onResponse = onResponse
	return me
}

func CreateAPITokenClient(ctx context.Context, classicURL string, apiToken string) (*rest.Client, error) {
	return clients.Factory().
		WithUserAgent(version.UserAgent()).
		WithClassicURL(classicURL).
		WithAccessToken(apiToken).
		WithHTTPListener(logging.HTTPListener("classic")).
		WithRetryOptions(defaultRetryOptions).
		CreateClassicClientWithContext(ctx)
}
