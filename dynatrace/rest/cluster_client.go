/**
* @license
* Copyright 2020 Dynatrace LLC
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

	"github.com/google/uuid"
)

func ClusterClient(baseURL string, credentials *Credentials) Client {
	return &cluster_client{baseURL: baseURL, credentials: credentials}
}

type cluster_client struct {
	baseURL     string
	credentials *Credentials
}

func (me *cluster_client) Credentials() *Credentials {
	return me.credentials
}

func (me *cluster_client) Get(ctx context.Context, url string, expectedStatusCodes ...int) Request {
	req := &cluster_request{id: uuid.NewString(), ctx: ctx, client: me, url: url, method: "GET"}
	if len(expectedStatusCodes) > 0 {
		req.expect = statuscodes(expectedStatusCodes)
	}
	return req
}

func (me *cluster_client) Post(ctx context.Context, url string, payload any, expectedStatusCodes ...int) Request {
	req := &cluster_request{id: uuid.NewString(), ctx: ctx, client: me, url: url, method: "POST", payload: payload, headers: map[string]string{"Content-Type": "application/json"}}
	if len(expectedStatusCodes) > 0 {
		req.expect = statuscodes(expectedStatusCodes)
	}
	return req
}

func (me *cluster_client) Put(ctx context.Context, url string, payload any, expectedStatusCodes ...int) Request {
	req := &cluster_request{id: uuid.NewString(), ctx: ctx, client: me, url: url, method: "PUT", payload: payload, headers: map[string]string{"Content-Type": "application/json"}}
	if len(expectedStatusCodes) > 0 {
		req.expect = statuscodes(expectedStatusCodes)
	}
	return req
}

func (me *cluster_client) Delete(ctx context.Context, url string, expectedStatusCodes ...int) Request {
	req := &cluster_request{id: uuid.NewString(), ctx: ctx, client: me, url: url, method: "DELETE"}
	if len(expectedStatusCodes) > 0 {
		req.expect = statuscodes(expectedStatusCodes)
	}
	return req
}
