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

func HybridClient(envURL, apiToken string, platform PlatformCredentials) Client {
	return &hybrid_client{envURL: envURL, apiTokenValue: apiToken, platform: platform}
}

type hybrid_client struct {
	envURL        string
	apiTokenValue string
	platform      PlatformCredentials
}

func (me *hybrid_client) environmentURL() string { return me.envURL }

func (me *hybrid_client) apiToken() string { return me.apiTokenValue }

func (me *hybrid_client) platformCredentials() PlatformCredentials { return me.platform }

func (me *hybrid_client) Get(ctx context.Context, url string, expectedStatusCodes ...int) Request {
	req := &hybrid_request{id: uuid.NewString(), ctx: ctx, client: me, url: url, method: http.MethodGet}
	if len(expectedStatusCodes) > 0 {
		req.expect = statuscodes(expectedStatusCodes)
	}
	return req
}

func (me *hybrid_client) Post(ctx context.Context, url string, payload any, expectedStatusCodes ...int) Request {
	req := &hybrid_request{id: uuid.NewString(), ctx: ctx, client: me, url: url, method: http.MethodPost, payload: payload, headers: headers.ContentType.ApplicationJSON}
	if len(expectedStatusCodes) > 0 {
		req.expect = statuscodes(expectedStatusCodes)
	}
	return req
}

func (me *hybrid_client) Put(ctx context.Context, url string, payload any, expectedStatusCodes ...int) Request {
	req := &hybrid_request{id: uuid.NewString(), ctx: ctx, client: me, url: url, method: http.MethodPut, payload: payload, headers: headers.ContentType.ApplicationJSON}
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
	if v := me.ctx.Value("DYNATRACE_HTTP_OAUTH_PREFERENCE"); v != nil {
		if bv, ok := v.(bool); ok {
			isOAuthPreferred = bv
		}
	}

	client := me.client.(environmentClient)

	isAPITokenPossible := containsAPIToken(client.apiToken())
	isOAuthPossible := containsOAuthOrPlatformToken(client.platformCredentials())

	if !isAPITokenPossible && !isOAuthPossible {
		if isOAuthPreferred {
			return NoOAuthCredentialsError
		}
		return NoAPITokenError
	}

	if (isAPITokenPossible && !isOAuthPossible) || (isAPITokenPossible && !isOAuthPreferred) {
		classicRequest := classic_request(*me)
		if client.environmentURL() == TestCaseEnvURL {
			return errors.New("classic")
		}
		return classicRequest.Finish(optionalTarget...)
	}

	if !isOAuthPossible {
		return NoOAuthCredentialsError
	}

	platformRequest := platform_request(*me)
	if client.environmentURL() == TestCaseEnvURL {
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
