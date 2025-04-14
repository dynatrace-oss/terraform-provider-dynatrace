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
	"errors"
	"net/http"
	"net/url"
	"strings"
	"sync"

	restlogging "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest/logging"
	"github.com/dynatrace/dynatrace-configuration-as-code-core/api/rest"
	"github.com/dynatrace/dynatrace-configuration-as-code-core/clients"
	"github.com/google/uuid"
)

func ClusterV1Client(credentials *Credentials) Client {
	creds := *credentials
	creds.URL = strings.TrimSuffix(credentials.Cluster.URL, "/") + "/api/v1.0/onpremise"
	creds.Token = credentials.Cluster.Token
	return &cluster_v1_client{credentials: &creds}
}

type cluster_v1_client struct {
	credentials *Credentials
}

func (me *cluster_v1_client) Credentials() *Credentials {
	return me.credentials
}

func (me *cluster_v1_client) Get(ctx context.Context, url string, expectedStatusCodes ...int) Request {
	req := &cluster_v1_request{id: uuid.NewString(), ctx: ctx, client: me, url: url, method: http.MethodGet}
	if len(expectedStatusCodes) > 0 {
		req.expect = statuscodes(expectedStatusCodes)
	}
	return req
}

func (me *cluster_v1_client) Post(ctx context.Context, url string, payload any, expectedStatusCodes ...int) Request {
	req := &cluster_v1_request{id: uuid.NewString(), ctx: ctx, client: me, url: url, method: http.MethodPost, payload: payload, headers: Headers.ContentType.ApplicationJSON}
	if len(expectedStatusCodes) > 0 {
		req.expect = statuscodes(expectedStatusCodes)
	}
	return req
}

func (me *cluster_v1_client) Put(ctx context.Context, url string, payload any, expectedStatusCodes ...int) Request {
	req := &cluster_v1_request{id: uuid.NewString(), ctx: ctx, client: me, url: url, method: http.MethodPut, payload: payload, headers: Headers.ContentType.ApplicationJSON}
	if len(expectedStatusCodes) > 0 {
		req.expect = statuscodes(expectedStatusCodes)
	}
	return req
}

func (me *cluster_v1_client) Delete(ctx context.Context, url string, expectedStatusCodes ...int) Request {
	req := &cluster_v1_request{id: uuid.NewString(), ctx: ctx, client: me, url: url, method: http.MethodDelete}
	if len(expectedStatusCodes) > 0 {
		req.expect = statuscodes(expectedStatusCodes)
	}
	return req
}

type cluster_v1_request request

var clusterV1ClientCache = map[string]*rest.Client{}

var clusterV1ClientCacheMutex sync.Mutex

func clusterV1Client(baseURL string, apiToken string) (*rest.Client, error) {
	clusterV1ClientCacheMutex.Lock()
	defer clusterV1ClientCacheMutex.Unlock()

	if client, found := clusterV1ClientCache[baseURL]; found {
		return client, nil
	}

	factory := clients.Factory()
	factory = factory.WithUserAgent("Dynatrace Terraform Provider")
	factory = factory.WithClassicURL(baseURL)
	factory = factory.WithAccessToken(apiToken)
	factory = factory.WithHTTPListener(restlogging.HTTPListener("clustv1 "))

	client, err := factory.CreateClassicClient()

	if client != nil && err == nil {
		clusterV1ClientCache[baseURL] = client
	}

	return client, err
}

func (me *cluster_v1_request) Finish(optionalTarget ...any) error {
	credentials := me.client.Credentials()
	if DYNATRACE_HTTP_LEGACY {
		if !credentials.ContainsClusterURL() {
			return NoClusterURLError
		}
		if !credentials.ContainsClusterToken() {
			return NoClusterTokenError
		}
		legacyRequest := legacy_request(*me)
		if credentials.URL == TestCaseEnvURL {
			return errors.New("legacy")
		}
		return legacyRequest.Finish(optionalTarget...)
	}

	var target any
	if len(optionalTarget) > 0 {
		target = optionalTarget[0]
	}

	PreRequest()

	clusterURL := strings.TrimSuffix(me.client.Credentials().Cluster.URL, "/") + "/api/v1.0/onpremise"

	client, err := clusterV1Client(clusterURL, me.client.Credentials().Cluster.Token)
	if err != nil {
		return err
	}

	fullURL, err := url.Parse(me.url)
	if err != nil {
		return err
	}
	return request(*me).HandleResponse(client, fullURL, target)
}

func (me *cluster_v1_request) Expect(codes ...int) Request {
	me.expect = statuscodes(codes)
	return me
}

func (me *cluster_v1_request) OnResponse(onResponse func(resp *http.Response)) Request {
	me.onResponse = onResponse
	return me
}
