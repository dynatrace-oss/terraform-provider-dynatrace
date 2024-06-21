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

package document

import (
	"bytes"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/monaco/pkg/rest"
)

type Response struct {
	ID        string `json:"id"`
	Actor     string `json:"actor"`
	Owner     string `json:"owner"`
	Name      string `json:"name"`
	Type      string `json:"type"`
	IsPrivate bool   `json:"isPrivate"`
	Content   string `json:"content"`
	Version   int    `json:"version"`
}

type listResponse struct {
	TotalCount  int        `json:"totalCount"`
	Documents   []Response `json:"documents"`
	NextPageKey string     `json:"nextPageKey"`
}

type Client struct {
	url       string
	client    *http.Client
	resources map[ResourceType]Resource
}

// NewClient creates a new client to interact with the Automation API
func NewClient(url string, client *http.Client) *Client {
	return &Client{url: url, client: client, resources: resources}
}

// LIST returns all document objects
func (a Client) LIST(resourceType ResourceType) (res []Response, err error) {
	var retVal []Response
	var result listResponse
	result.NextPageKey = "initial"

	for result.NextPageKey != "" {
		if result.NextPageKey == "initial" {
			result.NextPageKey = ""
		}

		u, err := NextPageURL(a.url, a.resources[resourceType].Path, result.NextPageKey)
		if err != nil {
			return nil, err
		}

		resp, err := rest.Get(a.client, u)
		if err != nil {
			return nil, err
		}

		if !resp.Success() {
			return nil, ResponseErr{
				StatusCode: resp.StatusCode,
				Message:    "Failed to list automation objects",
				Data:       resp.Body,
			}
		}

		err = json.Unmarshal(resp.Body, &result)
		if err != nil {
			return nil, err
		}

		retVal = append(retVal, result.Documents...)
	}

	return retVal, nil
}

// GET returns one specific automation object
func (a Client) GET(resourceType ResourceType, id string) (res *Response, err error) {
	var resp rest.Response
	var e Response

	if resp, err = rest.Get(a.client, a.url+a.resources[resourceType].Path+"/"+id); err != nil {
		return nil, fmt.Errorf("unable to get object with ID %s: %w", id, err)
	}

	if !resp.Success() {
		return nil, ResponseErr{StatusCode: resp.StatusCode, Message: "Failed to get automation object", Data: resp.Body}
	}

	contentType := resp.Headers["Content-Type"][0]
	boundaryIndex := strings.Index(contentType, "boundary=")
	if boundaryIndex == -1 {
		return nil, fmt.Errorf("no boundary parameter found in Content-Type header")
	}
	boundary := contentType[boundaryIndex+len("boundary="):]

	r := multipart.NewReader(bytes.NewReader(resp.Body), boundary)

	form, err := r.ReadForm(0)
	if err != nil {
		return nil, fmt.Errorf("unable to read multipart form: %w", err)
	}

	if len(form.Value["metadata"]) == 0 {
		return nil, fmt.Errorf("metadata field not found in response")
	}

	err = json.Unmarshal([]byte(form.Value["metadata"][0]), &e)
	if err != nil {
		return &e, fmt.Errorf("unable to unmarshal metadata: %w", err)
	}

	file, err := form.File["content"][0].Open()
	if err != nil {
		return &e, fmt.Errorf("unable to open file: %w", err)
	}
	defer file.Close()

	fileContent := new(bytes.Buffer)
	_, err = fileContent.ReadFrom(file)
	if err != nil {
		return &e, fmt.Errorf("unable to read file: %w", err)
	}

	e.Content = fileContent.String()

	return &e, err
}

// UPDATE updates a given document object
func (a Client) UPDATE(resourceType ResourceType, id string, data *bytes.Buffer, contentType string) (err error) {
	if id == "" {
		return fmt.Errorf("id must be non empty")
	}

	doc, err := a.GET(resourceType, id)
	if err != nil {
		return fmt.Errorf("unable to get object with ID %s: %w", id, err)
	}

	var urlParams = make(map[string]string)
	urlParams["optimistic-locking-version"] = fmt.Sprint(doc.Version)

	resp, err := rest.PatchMultiPartFile(a.client, a.url+a.resources[resourceType].Path+"/"+id, data, contentType, urlParams)
	if err != nil {
		return fmt.Errorf("unable to update object with ID %s: %w", id, err)
	}

	if resp.Success() {
		return nil
	}
	return ResponseErr{StatusCode: resp.StatusCode, Message: "Failed to update document object", Data: resp.Body}
}

// INSERT creates a given document object
func (a Client) INSERT(resourceType ResourceType, data *bytes.Buffer, contentType string) (id string, err error) {
	resp, err := rest.PostMultiPartFile(a.client, a.url+a.resources[resourceType].Path, data, contentType)
	if err != nil {
		return "", err
	}

	if !resp.Success() {
		return id, ResponseErr{StatusCode: resp.StatusCode, Message: "Failed to insert document object", Data: resp.Body}
	}

	var e Response
	err = json.Unmarshal(resp.Body, &e)

	return e.ID, err
}

// DELETE removes a given document object by ID
func (a Client) DELETE(resourceType ResourceType, id string) (err error) {
	if id == "" {
		return fmt.Errorf("id must be non empty")
	}

	doc, err := a.GET(resourceType, id)
	if err != nil {
		return fmt.Errorf("unable to get object with ID %s: %w", id, err)
	}

	var urlParams = make(map[string]string)
	urlParams["optimistic-locking-version"] = fmt.Sprint(doc.Version)

	if err = rest.DeleteConfig(a.client, a.url+a.resources[resourceType].Path, id, urlParams); err != nil {
		return fmt.Errorf("unable to delete object with ID %s: %w", id, err)
	}

	return nil
}

func NextPageURL(baseURL, path string, nextPageKey string) (string, error) {
	u, e := url.Parse(baseURL)
	if e != nil {
		return "", e
	}

	u.Path, e = url.JoinPath(u.Path, path)
	if e != nil {
		return "", e
	}

	q := u.Query()
	q.Add("page-key", nextPageKey)
	u.RawQuery = q.Encode()

	return u.String(), nil
}

func (r *listResponse) UnmarshalJSON(data []byte) error {
	type Alias listResponse

	var aux struct {
		*Alias
		NextPageKey *string `json:"nextPageKey"`
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	if aux.NextPageKey == nil {
		aux.NextPageKey = new(string)
	}

	r.TotalCount = aux.TotalCount
	r.Documents = aux.Documents
	r.NextPageKey = *aux.NextPageKey

	return nil
}
