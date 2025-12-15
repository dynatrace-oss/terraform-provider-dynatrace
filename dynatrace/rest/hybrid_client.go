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

	"github.com/google/uuid"
)

func HybridClient(credentials *Credentials) Client {
	return &hybrid_client{credentials: credentials}
}

type hybrid_client struct {
	credentials *Credentials
}

func (me *hybrid_client) Credentials() *Credentials {
	return me.credentials
}

func (me *hybrid_client) Get(ctx context.Context, url string, expectedStatusCodes ...int) Request {
	req := &hybrid_request{id: uuid.NewString(), ctx: ctx, client: me, url: url, method: http.MethodGet}
	if len(expectedStatusCodes) > 0 {
		req.expect = statuscodes(expectedStatusCodes)
	}
	return req
}

func (me *hybrid_client) Post(ctx context.Context, url string, payload any, expectedStatusCodes ...int) Request {
	req := &hybrid_request{id: uuid.NewString(), ctx: ctx, client: me, url: url, method: http.MethodPost, payload: payload, headers: Headers.ContentType.ApplicationJSON}
	if len(expectedStatusCodes) > 0 {
		req.expect = statuscodes(expectedStatusCodes)
	}
	return req
}

func (me *hybrid_client) Put(ctx context.Context, url string, payload any, expectedStatusCodes ...int) Request {
	req := &hybrid_request{id: uuid.NewString(), ctx: ctx, client: me, url: url, method: http.MethodPut, payload: payload, headers: Headers.ContentType.ApplicationJSON}
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
	isOAuthPreferred := DYNATRACE_HTTP_OAUTH_PREFERENCE
	if v := me.ctx.Value("DYNATRACE_HTTP_OAUTH_PREFERENCE"); v != nil {
		if bv, ok := v.(bool); ok {
			isOAuthPreferred = bv
		}
	}

	credentials := me.client.Credentials()

	if DYNATRACE_HTTP_LEGACY {
		if !credentials.ContainsAPIToken() {
			return NoAPITokenError
		}
		legacyRequest := legacy_request(*me)
		if credentials.URL == TestCaseEnvURL {
			return errors.New("legacy")
		}
		return legacyRequest.Finish(optionalTarget...)
	}

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
		classicRequest := classic_request(*me)
		if credentials.URL == TestCaseEnvURL {
			return errors.New("classic")
		}
		return classicRequest.Finish(optionalTarget...)
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

func (me *hybrid_request) SetHeader(name string, value string) {
	if me.headers == nil {
		me.headers = map[string]string{}
	}
	me.headers[name] = value
}

func NewPreferOAuthContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, "DYNATRACE_HTTP_OAUTH_PREFERENCE", true)
}
