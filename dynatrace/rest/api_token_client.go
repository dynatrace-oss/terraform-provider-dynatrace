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
	"net/url"
	"sync"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest/logging"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/provider/version"
	"github.com/google/uuid"

	"github.com/dynatrace/dynatrace-configuration-as-code-core/api/rest"
	"github.com/dynatrace/dynatrace-configuration-as-code-core/clients"
)

var NoClassicURLDefinedErr = errors.New("no Environment URL has been specified. Use either the environment variable `DYNATRACE_ENV_URL` or the configuration attribute `dt_env_url` of the provider for that")

func APITokenClient(credentials *Credentials) Client {
	return &api_token_client{credentials: credentials}
}

type api_token_client struct {
	credentials *Credentials
}

func (me *api_token_client) Credentials() *Credentials {
	return me.credentials
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
	credentials := me.client.Credentials()
	if !credentials.ContainsAPIToken() {
		return NoAPITokenError
	}

	var target any
	if len(optionalTarget) > 0 {
		target = optionalTarget[0]
	}

	client, err := createAPITokenClient(me.client.Credentials().ClassicEnvironmentURL, me.client.Credentials().Token)
	if err != nil {
		return err
	}

	pathURL, err := url.Parse(me.url)
	if err != nil {
		return err
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

var apiTokenClientCache = map[string]*rest.Client{}

var apiTokenClientCacheMutex sync.Mutex

func createAPITokenClient(classicURL string, apiToken string) (*rest.Client, error) {
	if classicURL == "" {
		// sanity check - this should not be empty anymore at this point
		return nil, NoClassicURLDefinedErr
	}

	apiTokenClientCacheMutex.Lock()
	defer apiTokenClientCacheMutex.Unlock()

	if client, found := apiTokenClientCache[classicURL]; found {
		return client, nil
	}

	client, err := clients.Factory().
		WithUserAgent(version.UserAgent()).
		WithClassicURL(classicURL).
		WithAccessToken(apiToken).
		WithHTTPListener(logging.HTTPListener("classic")).
		WithRetryOptions(defaultRetryOptions).
		CreateClassicClient()

	if err != nil {
		return nil, err
	}

	apiTokenClientCache[classicURL] = client

	return client, nil
}
