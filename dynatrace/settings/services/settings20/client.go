// @license
// Copyright 2023 Dynatrace LLC
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package settings20

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/provider/logging"
	"github.com/dynatrace/dynatrace-configuration-as-code-core/api"

	"github.com/dynatrace/dynatrace-configuration-as-code-core/api/rest"
	"github.com/go-logr/logr"
)

const endpointPath = "api/v2/settings/objects"

type Response struct {
	api.Response
	ID    string
	Items []struct {
		ID    string          `json:"objectId"`
		Scope string          `json:"scope"`
		Value json.RawMessage `json:"value"`
	}
}

type listResponse struct {
	api.Response
	Items       []json.RawMessage `json:"items"`
	NextPageKey string            `json:"nextPageKey"`
	PageSize    int               `json:"pageSize"`
	TotalCount  int               `json:"totalCount"`
}

type retrySettings struct {
	maxRetries   int
	waitDuration time.Duration
}

type Client struct {
	tokenClient *client
	oAuthClient *client
}

type client struct {
	client        *rest.Client
	schemaID      string
	retrySettings retrySettings
}

// Option represents a functional Option for the Client.
type Option func(*Client)

// WithUpdateRetrySettings sets the maximum number of retries as well as duration between retries for Update HTTP requests
func WithUpdateRetrySettings(maxRetries int, waitDuration time.Duration) Option {
	return func(c *Client) {
		if c.tokenClient != nil {
			c.tokenClient.retrySettings = retrySettings{maxRetries: maxRetries, waitDuration: waitDuration}
		}
		if c.oAuthClient != nil {
			c.oAuthClient.retrySettings = retrySettings{maxRetries: maxRetries, waitDuration: waitDuration}
		}
	}
}

func NewClient(tokenClient *rest.Client, oAuthClient *rest.Client, schemaID string, option ...Option) *Client {
	var tc *client
	var oac *client
	if tokenClient != nil {
		tc = &client{client: tokenClient, schemaID: schemaID, retrySettings: retrySettings{maxRetries: 5, waitDuration: time.Second}}
	}
	if oAuthClient != nil {
		oac = &client{client: oAuthClient, schemaID: schemaID, retrySettings: retrySettings{maxRetries: 5, waitDuration: time.Second}}
	}
	c := &Client{tokenClient: tc, oAuthClient: oac}
	for _, o := range option {
		o(c)
	}
	return c
}

func (c client) Get(ctx context.Context, id string) (Response, error) {
	resp, err := c.get(ctx, id)
	if err != nil {
		return Response{}, err
	}
	return Response{api.Response{StatusCode: resp.StatusCode, Data: resp.Payload, Request: resp.RequestInfo}, id, nil}, nil
}

func (c client) List(ctx context.Context) (Response, error) {
	items := []struct {
		ID    string          `json:"objectId"`
		Scope string          `json:"scope"`
		Value json.RawMessage `json:"value"`
	}{}
	nextPageKey := ""

	var resp listResponse
	var err error
	for {
		resp, err = c.list(ctx, nextPageKey)
		if err != nil {
			return Response{}, err
		}

		for _, v := range resp.Items {
			var item struct {
				ID    string          `json:"objectId"`
				Scope string          `json:"scope"`
				Value json.RawMessage `json:"value"`
			}
			if err := json.Unmarshal(v, &item); err != nil {
				return Response{}, err
			}
			items = append(items, item)
		}
		if len(resp.NextPageKey) == 0 {
			break
		}
		nextPageKey = resp.NextPageKey
	}
	return Response{
		Response: api.Response{
			StatusCode: resp.StatusCode,
			Request:    resp.Request,
			Data:       resp.Response.Data,
		},
		Items: items,
	}, nil
}

func (c client) Create(ctx context.Context, scope string, data []byte) (Response, error) {
	soc := SettingsObjectCreate{
		SchemaID: c.schemaID,
		Scope:    "environment",
		Value:    json.RawMessage(data),
	}
	if len(scope) > 0 {
		soc.Scope = scope
	}
	dsoc, e := json.Marshal([]SettingsObjectCreate{soc})
	if e != nil {
		return Response{api.Response{StatusCode: 0, Data: nil, Request: rest.RequestInfo{}}, "", nil}, e
	}
	resp, err := c.create(ctx, dsoc)
	if err != nil {
		return Response{}, err
	}
	tmp := []struct {
		ID string `json:"objectId"`
	}{}
	if err := json.Unmarshal(resp.Payload, &tmp); err != nil {
		return Response{Response: api.Response{StatusCode: resp.StatusCode, Data: resp.Payload, Request: resp.RequestInfo}, ID: ""}, err
	}
	return Response{Response: api.Response{StatusCode: resp.StatusCode, Data: resp.Payload, Request: resp.RequestInfo}, ID: tmp[0].ID}, nil
}

func (c client) Update(ctx context.Context, id string, data []byte) (Response, error) {
	sou := SettingsObjectUpdate{
		Value: json.RawMessage(data),
	}
	dsou, e := json.Marshal(sou)
	if e != nil {
		return Response{api.Response{StatusCode: 0, Data: nil, Request: rest.RequestInfo{}}, "", nil}, e
	}

	logger := logr.FromContextOrDiscard(ctx)

	var resp rest.Response
	var err error

	for i := 0; i < c.retrySettings.maxRetries; i++ {
		logger.V(1).Info(fmt.Sprintf("Trying to update setting with id %q (%d/%d retries)", id, i+1, c.retrySettings.maxRetries))

		resp, err = c.update(ctx, id, dsou)
		if err != nil {
			return Response{}, err
		}

		if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden || resp.StatusCode == http.StatusBadRequest {
			return Response{api.Response{StatusCode: resp.StatusCode, Data: resp.Payload, Request: resp.RequestInfo}, id, nil}, nil
		}

		if resp.IsSuccess() {
			logger.Info(fmt.Sprintf("Updated setting with id %q", id))
			return Response{api.Response{StatusCode: resp.StatusCode, Data: resp.Payload, Request: resp.RequestInfo}, id, nil}, nil
		}
		time.Sleep(c.retrySettings.waitDuration)
	}
	return Response{Response: api.Response{StatusCode: resp.StatusCode, Data: resp.Payload, Request: resp.RequestInfo}}, err
}

func (c client) Delete(ctx context.Context, id string) (Response, error) {
	if id == "" {
		return Response{}, fmt.Errorf("id must be non-empty")
	}
	path, err := url.JoinPath(endpointPath, id)
	if err != nil {
		return Response{}, fmt.Errorf("failed to create URL: %w", err)
	}
	resp, err := c.client.DELETE(ctx, path, rest.RequestOptions{})
	if err != nil {
		return Response{}, fmt.Errorf("unable to delete object with id %q: %w", id, err)
	}
	return Response{api.Response{StatusCode: resp.StatusCode, Data: resp.Payload, Request: resp.RequestInfo}, id, nil}, err
}

func (c client) create(ctx context.Context, data []byte) (rest.Response, error) {
	r, err := c.client.POST(ctx, endpointPath, bytes.NewReader(data), rest.RequestOptions{})
	if err != nil {
		return rest.Response{}, fmt.Errorf("failed to create object: %w", err)
	}
	return r, nil
}

func (c client) get(ctx context.Context, id string) (rest.Response, error) {
	path, err := url.JoinPath(endpointPath, id)
	if err != nil {
		return rest.Response{}, fmt.Errorf("failed to create URL: %w", err)
	}
	return c.client.GET(ctx, path, rest.RequestOptions{})
}

func (c client) list(ctx context.Context, nextPageKey string) (listResponse, error) {
	var resp rest.Response
	var err error
	var requestOptions rest.RequestOptions
	if len(nextPageKey) == 0 {
		requestOptions.QueryParams = url.Values{
			"schemaIds": []string{c.schemaID},
			"fields":    []string{"objectId,scope,value"},
			"pageSize":  []string{"500"},
		}
	} else {
		requestOptions.QueryParams = url.Values{
			"nextPageKey": []string{nextPageKey},
		}
	}
	resp, err = c.client.GET(ctx, endpointPath, requestOptions)
	if err != nil {
		return listResponse{}, fmt.Errorf("failed to list objects:%w", err)
	}
	l, err := unmarshalJSONList(&resp)
	if err != nil {
		return listResponse{}, fmt.Errorf("failed to parse list response:%w", err)
	}
	return l, nil
}

func (c client) update(ctx context.Context, id string, data []byte) (rest.Response, error) {
	var err error

	// construct path for PUT request
	path, err := url.JoinPath(endpointPath, id)
	if err != nil {
		return rest.Response{}, fmt.Errorf("failed to join URL: %w", err)
	}

	// make PUT request
	return c.client.PUT(ctx, path, bytes.NewReader(data), rest.RequestOptions{})
}

// unmarshalJSONList unmarshals JSON data into a listResponse struct.
func unmarshalJSONList(raw *rest.Response) (listResponse, error) {
	var r listResponse
	err := json.Unmarshal(raw.Payload, &r)
	if err != nil {
		return listResponse{}, fmt.Errorf("failed to unmarshal response: %w", err)
	}
	r.Data = raw.Payload
	r.StatusCode = raw.StatusCode
	r.Request = raw.RequestInfo
	return r, nil
}

func (c Client) Get(ctx context.Context, id string) (Response, error) {
	if c.tokenClient != nil {
		return c.tokenClient.Get(ctx, id)
	}
	if c.oAuthClient != nil {
		return c.oAuthClient.Get(ctx, id)
	}
	return Response{}, errors.New("no client configured")
}

func (c Client) List(ctx context.Context) (Response, error) {
	if c.tokenClient != nil {
		return c.tokenClient.List(ctx)
	}
	if c.oAuthClient != nil {
		return c.oAuthClient.List(ctx)
	}
	return Response{}, errors.New("no client configured")
}

type SettingsObjectResponse struct {
	Code  int `json:"code"`
	Error struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

func (sor SettingsObjectResponse) RequiresOAuth() bool {
	if sor.Code != http.StatusBadRequest {
		return false
	}
	if sor.Error.Code != http.StatusBadRequest {
		return false
	}
	return sor.Error.Message == "No OAuth token provided"
}

func CreateRequiresOAuth(response Response) bool {
	if response.StatusCode != http.StatusBadRequest {
		return false
	}
	if len(response.Data) == 0 {
		return false
	}
	sor := []SettingsObjectResponse{{}}
	if err := json.Unmarshal(response.Data, &sor); err != nil {
		return false
	}
	if len(sor) == 0 {
		return false
	}
	for _, entry := range sor {
		if entry.RequiresOAuth() {
			return true
		}
	}
	return false
}

func UpdateRequiresOAuth(response Response) bool {
	if response.StatusCode != http.StatusBadRequest {
		return false
	}
	if len(response.Data) == 0 {
		return false
	}
	sor := SettingsObjectResponse{}
	if err := json.Unmarshal(response.Data, &sor); err != nil {
		logging.File.Println(err)
		return false
	}
	return sor.RequiresOAuth()
}

func (c Client) Create(ctx context.Context, scope string, data []byte) (Response, error) {
	if c.tokenClient != nil {
		response, err := c.tokenClient.Create(ctx, scope, data)
		if err != nil {
			return response, err
		}
		if CreateRequiresOAuth(response) && c.oAuthClient != nil {
			return c.oAuthClient.Create(ctx, scope, data)
		}
		return response, err
	}
	if c.oAuthClient != nil {
		return c.oAuthClient.Create(ctx, scope, data)
	}
	return Response{}, errors.New("no client configured")
}

func (c Client) Update(ctx context.Context, id string, data []byte) (Response, error) {
	if c.tokenClient != nil {
		response, err := c.tokenClient.Update(ctx, id, data)
		if err != nil {
			return response, err
		}
		if UpdateRequiresOAuth(response) && c.oAuthClient != nil {
			return c.oAuthClient.Update(ctx, id, data)
		}
		return response, err
	}
	if c.oAuthClient != nil {
		return c.oAuthClient.Update(ctx, id, data)
	}
	return Response{}, errors.New("no client configured")
}

func (c Client) Delete(ctx context.Context, id string) (Response, error) {
	if c.tokenClient != nil {
		return c.tokenClient.Delete(ctx, id)
	}
	if c.oAuthClient != nil {
		return c.oAuthClient.Delete(ctx, id)
	}
	return Response{}, errors.New("no client configured")
}
