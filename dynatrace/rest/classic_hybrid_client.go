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
	"context"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/provider/envutils"
)

// ClassicHybridClient prefers classic over platform requests and doesn't consider the DYNATRACE_HTTP_OAUTH_PREFERENCE env variable.
func ClassicHybridClient(clientSet ClientSet) Client {
	credentials := clientSet.Credentials()
	return &classic_hybrid_client{client: HybridClient(credentials), credentials: credentials}
}

type classic_hybrid_client struct {
	client      Client
	credentials *Credentials
}

func (me *classic_hybrid_client) Credentials() *Credentials {
	return me.credentials
}

func (me *classic_hybrid_client) Get(ctx context.Context, url string, expectedStatusCodes ...int) Request {
	return me.client.Get(newDisableOAuthPreferenceContext(ctx), url, expectedStatusCodes...)
}

func (me *classic_hybrid_client) Post(ctx context.Context, url string, payload any, expectedStatusCodes ...int) Request {
	return me.client.Post(newDisableOAuthPreferenceContext(ctx), url, payload, expectedStatusCodes...)
}

func (me *classic_hybrid_client) Put(ctx context.Context, url string, payload any, expectedStatusCodes ...int) Request {
	return me.client.Put(newDisableOAuthPreferenceContext(ctx), url, payload, expectedStatusCodes...)
}

func (me *classic_hybrid_client) Delete(ctx context.Context, url string, expectedStatusCodes ...int) Request {
	return me.client.Delete(newDisableOAuthPreferenceContext(ctx), url, expectedStatusCodes...)
}

func newDisableOAuthPreferenceContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, envutils.DynatraceHTTPOAuthPreference.Key, false)
}
