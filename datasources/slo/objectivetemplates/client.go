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

package objectivetemplates

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/dynatrace/dynatrace-configuration-as-code-core/api"
	"github.com/dynatrace/dynatrace-configuration-as-code-core/api/rest"
)

type client interface {
	List(ctx context.Context, queryParams map[string][]string) (*http.Response, error)
}

// Client is the HTTP client to be used for interacting with the SLO objective templates API
type Client struct {
	client client
}

// NewClient creates a new SLO objective templates client
func NewClient(client *rest.Client) *Client {
	c := &Client{client: NewSLOObjectiveTemplatesClient(client)}
	return c
}

// Response contains the API response
type Response struct {
	api.Response
	// Metadata
}

func (c Client) List(ctx context.Context, queryParams map[string][]string) (api.Response, error) {
	resp, err := processHttpResponse(c.client.List(ctx, queryParams))
	if err != nil {
		return resp, err
	}
	return resp, nil
}

func processHttpResponse(resp *http.Response, err error) (api.Response, error) {
	if resp != nil {
		defer resp.Body.Close()
	}

	if err != nil {
		return api.Response{}, err
	}

	if !rest.IsSuccess(resp) {
		return api.Response{}, api.NewAPIErrorFromResponse(resp)
	}

	var body []byte
	body, err = io.ReadAll(resp.Body)
	if err != nil {
		return api.Response{}, fmt.Errorf("unable to read API response body: %w", err)
	}

	return api.NewResponseFromHTTPResponseAndBody(resp, body), nil
}
