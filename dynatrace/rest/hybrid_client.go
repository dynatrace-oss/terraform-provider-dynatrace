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
	"net/http"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/provider/envutils"

	"github.com/google/uuid"
)

func HybridClient(clientSet ClientSet) Client {
	return &hybrid_client{clientSet: clientSet}
}

type hybrid_client struct {
	clientSet ClientSet
}

func (me *hybrid_client) ClientSet() ClientSet {
	return me.clientSet
}

func (me *hybrid_client) Get(ctx context.Context, url string, expectedStatusCodes ...int) Request {
	req := &hybrid_request{id: uuid.NewString(), ctx: ctx, client: me, url: url, method: http.MethodGet}
	if len(expectedStatusCodes) > 0 {
		req.expect = statuscodes(expectedStatusCodes)
	}
	return req
}

func (me *hybrid_client) Post(ctx context.Context, url string, payload any, expectedStatusCodes ...int) Request {
	req := &hybrid_request{id: uuid.NewString(), ctx: ctx, client: me, url: url, method: http.MethodPost, payload: payload}
	if len(expectedStatusCodes) > 0 {
		req.expect = statuscodes(expectedStatusCodes)
	}
	return req
}

func (me *hybrid_client) Put(ctx context.Context, url string, payload any, expectedStatusCodes ...int) Request {
	req := &hybrid_request{id: uuid.NewString(), ctx: ctx, client: me, url: url, method: http.MethodPut, payload: payload}
	if len(expectedStatusCodes) > 0 {
		req.expect = statuscodes(expectedStatusCodes)
	}
	return req
}

func (me *hybrid_client) Delete(ctx context.Context, url string, expectedStatusCodes ...int) Request {
	req := &hybrid_request{id: uuid.NewString(), ctx: ctx, client: me, url: url, method: http.MethodDelete}
	if len(expectedStatusCodes) > 0 {
		req.expect = statuscodes(expectedStatusCodes)
	}
	return req
}

func (me *hybrid_request) Finish(optionalTarget ...any) error {
	isOAuthPreferred := envutils.DynatraceHTTPOAuthPreference.Get()
	if v := me.ctx.Value(envutils.DynatraceHTTPOAuthPreference.Key); v != nil {
		if bv, ok := v.(bool); ok {
			isOAuthPreferred = bv
		}
	}

	credentials := me.client.ClientSet().Credentials()

	if !credentials.ContainsAPIToken() && !credentials.ContainsOAuthOrPlatformToken() {
		if isOAuthPreferred {
			return NoOAuthCredentialsError
		}
		return NoAPITokenError
	}

	isAPITokenPossible := credentials.ContainsAPIToken()
	isOAuthPossible := credentials.ContainsOAuthOrPlatformToken()

	if (isAPITokenPossible && !isOAuthPossible) || (isAPITokenPossible && !isOAuthPreferred) {
		if !credentials.ContainsAPIToken() {
			return NoAPITokenError
		}
		apiTokenRequest := api_token_request(*me)
		return apiTokenRequest.Finish(optionalTarget...)
	}

	if !credentials.ContainsOAuthOrPlatformToken() {
		return NoOAuthCredentialsError
	}

	platformRequest := platform_request(*me)
	if credentials.URL == TestCaseEnvURL {
		return errors.New("platform")
	}
	return platformRequest.Finish(optionalTarget...)
}

type hybrid_request request

func (me *hybrid_request) Expect(codes ...int) Request {
	me.expect = statuscodes(codes)
	return me
}

func (me *hybrid_request) OnResponse(onResponse func(resp *http.Response)) Request {
	me.onResponse = onResponse
	return me
}

func NewPreferOAuthContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, envutils.DynatraceHTTPOAuthPreference.Key, true)
}
