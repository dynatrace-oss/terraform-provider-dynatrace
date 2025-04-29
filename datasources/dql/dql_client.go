/*
 * @license
 * Copyright 2023 Dynatrace LLC
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package dql

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/dynatrace/dynatrace-configuration-as-code-core/api/rest"
)

const (
	queryPath = "/platform/storage/query/v1/query:execute"
	pollPath  = "/platform/storage/query/v1/query:poll"
)

type DQLClient struct {
	client *rest.Client
}

func NewDQLClient(client *rest.Client) *DQLClient {
	c := &DQLClient{
		client: client,
	}
	return c
}

func (c *DQLClient) Fetch(ctx context.Context, data []byte) (*http.Response, error) {
	resp, err := c.client.POST(ctx, queryPath, bytes.NewReader(data), rest.RequestOptions{ContentType: "application/json"})
	if err != nil {
		return nil, fmt.Errorf("unable to fetch query result: %w", err)
	}
	return resp, err
}

func (c *DQLClient) Poll(ctx context.Context, requestToken string) (*http.Response, error) {
	query := url.Values{}
	query.Add("request-token", requestToken)

	resp, err := c.client.GET(ctx, pollPath, rest.RequestOptions{QueryParams: query})
	if err != nil {
		return nil, fmt.Errorf("unable to poll for query result: %w", err)
	}
	return resp, err
}
