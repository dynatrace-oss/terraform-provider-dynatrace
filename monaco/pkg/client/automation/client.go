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

package automation

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/monaco/pkg/rest"
)

// Response is a "general" Response type holding the ID and the response payload
type Response struct {
	// ID is the identifier that will be used when creating a new automation object
	ID string `json:"id"`
	// Data is the whole body of an automation object
	Data []byte `json:"-"`
}

// UnmarshalJSON de-serializes JSON payload into [Response] type
func (r *Response) UnmarshalJSON(data []byte) error {
	var rawMap map[string]json.RawMessage
	if err := json.Unmarshal(data, &rawMap); err != nil {
		return err
	}
	if err := json.Unmarshal(rawMap["id"], &r.ID); err != nil {
		return err
	}
	r.Data = data
	return nil
}

type listResponse struct {
	Count   int        `json:"count"`
	Results []Response `json:"results"`
}

// Client can be used to interact with the Automation API
type Client struct {
	url       string
	client    *http.Client
	resources map[ResourceType]Resource
}

// NewClient creates a new client to interact with the Automation API
func NewClient(url string, client *http.Client) *Client {
	return &Client{url: url, client: client, resources: resources}
}

// LIST returns all automation objects
func (a Client) LIST(resourceType ResourceType) (res []Response, err error) {
	var retVal []Response
	var result listResponse
	result.Count = 1

	for len(retVal) < result.Count {

		u, err := NextPageURL(a.url, a.resources[resourceType].Path, len(retVal))
		if err != nil {
			return nil, err
		}

		// try to get the list of resources
		resp, err := rest.Get(a.client, u)
		if err != nil {
			return nil, err
		}

		// handle http error
		if !resp.Success() {
			return nil, ResponseErr{
				StatusCode: resp.StatusCode,
				Message:    "Failed to list automation objects",
				Data:       resp.Body,
			}
		}

		// unmarshal and return result
		err = json.Unmarshal(resp.Body, &result)
		if err != nil {
			return nil, err
		}
		retVal = append(retVal, result.Results...)
	}

	if len(retVal) != result.Count {
		log.Printf("[WARN] Total count of items returned for Automation API %q does not match count of actually received items. Expected: %d Got: %d.", resources[resourceType].Path, result.Count, len(retVal))

	}
	return retVal, nil
}

// GET returns one specific automation object
func (a Client) GET(resourceType ResourceType, id string) (res *Response, err error) {
	var resp rest.Response

	if resp, err = rest.Get(a.client, a.url+a.resources[resourceType].Path+"/"+id); err != nil {
		return nil, fmt.Errorf("unable to get object with ID %s: %w", id, err)
	}

	if !resp.Success() {
		return nil, ResponseErr{StatusCode: resp.StatusCode, Message: "Failed to get automation object", Data: resp.Body}
	}

	var e Response
	err = json.Unmarshal(resp.Body, &e)
	return &e, err
}

// UPDATE updates a given automation object
func (a Client) UPDATE(resourceType ResourceType, id string, data []byte) (err error) {
	if id == "" {
		return fmt.Errorf("id must be non empty")
	}
	if err := rmIDField(&data); err != nil {
		return fmt.Errorf("unable to remove id field from payload in order to update object with ID %s: %w", id, err)
	}
	resp, err := rest.Put(a.client, a.url+a.resources[resourceType].Path+"/"+id, data)
	if err != nil {
		return fmt.Errorf("unable to update object with ID %s: %w", id, err)
	}

	if resp.Success() {
		return nil
	}
	return ResponseErr{StatusCode: resp.StatusCode, Message: "Failed to insert automation object", Data: resp.Body}
}

// UPSERT creates a given automation object
func (a Client) INSERT(resourceType ResourceType, data []byte) (id string, err error) {
	if err := rmIDField(&data); err != nil {
		return "", fmt.Errorf("unable to remove id field from payload in when creating object: %w", err)
	}

	resp, err := rest.Post(a.client, a.url+a.resources[resourceType].Path, data)
	if err != nil {
		return "", err
	}

	// handle response err
	if !resp.Success() {
		return id, ResponseErr{StatusCode: resp.StatusCode, Message: "Failed to upsert automation object", Data: resp.Body}
	}

	// de-serialize response
	var e Response
	err = json.Unmarshal(resp.Body, &e)
	return e.ID, err
}

// DELETE removes a given automation object by ID
func (a Client) DELETE(resourceType ResourceType, id string) (err error) {
	if id == "" {
		return fmt.Errorf("id must be non empty")
	}

	var urlParams map[string]string
	if err = rest.DeleteConfig(a.client, a.url+a.resources[resourceType].Path, id, urlParams); err != nil {
		return fmt.Errorf("unable to delete object with ID %s: %w", id, err)
	}

	return nil
}

func setIDField(id string, data *[]byte) (err error) {
	var m map[string]any
	if err = json.Unmarshal(*data, &m); err != nil {
		return err
	}
	m["id"] = id
	*data, err = json.Marshal(m)
	return err
}

func rmIDField(data *[]byte) (err error) {
	var m map[string]any
	if err = json.Unmarshal(*data, &m); err != nil {
		return err
	}
	delete(m, "id")
	*data, err = json.Marshal(m)
	return err
}

func NextPageURL(baseURL, path string, offset int) (string, error) {
	u, e := url.Parse(baseURL)
	if e != nil {
		return "", e
	}

	u.Path, e = url.JoinPath(u.Path, path)
	if e != nil {
		return "", e
	}

	q := u.Query()
	q.Add("offset", strconv.Itoa(offset))
	u.RawQuery = q.Encode()

	return u.String(), nil
}
