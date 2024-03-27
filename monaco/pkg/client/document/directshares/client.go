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

package directshare

import (
	"encoding/json"
	"fmt"
	"net/http"

	directshares "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/documents/directshares/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/monaco/pkg/rest"
)

type Response struct {
	ID         string                  `json:"id"`
	DocumentId string                  `json:"documentId"`
	Access     []string                `json:"access"`
	UserCount  int                     `json:"userCount"`
	GroupCount int                     `json:"groupCount"`
	Recipients directshares.Recipients `json:"recipients"`
}

type Client struct {
	url       string
	client    *http.Client
	resources map[ResourceType]Resource
}

type tempRecipients struct {
	Recipients directshares.Recipients `json:"recipients"`
}

type tempRecipient struct {
	Recipient directshares.Recipient `json:"recipient"`
}

type IDs struct {
	Ids []string `json:"ids"`
}

func (r *tempRecipients) Contains(recipient tempRecipient) bool {
	for _, rec := range r.Recipients {
		if rec.Id == recipient.Recipient.Id && rec.Type == recipient.Recipient.Type {
			return true
		}
	}
	return false
}

func NewClient(url string, client *http.Client) *Client {
	return &Client{url: url, client: client, resources: resources}
}

// LIST returns all document objects
func (a Client) LIST(resourceType ResourceType) (res []Response, err error) {
	return nil, nil // not implemented
}

// GET returns one specific automation object
func (a Client) GET(resourceType ResourceType, id string) (res *Response, err error) {
	var resp rest.Response
	var resp_recipients rest.Response
	var e Response

	if resp, err = rest.Get(a.client, a.url+a.resources[resourceType].Path+"/"+id); err != nil {
		return nil, fmt.Errorf("unable to get object with ID %s: %w", id, err)
	}

	if !resp.Success() {
		return nil, ResponseErr{StatusCode: resp.StatusCode, Message: "Failed to get automation object", Data: resp.Body}
	}

	err = json.Unmarshal(resp.Body, &e)
	if err != nil {
		return &e, fmt.Errorf("unable to unmarshal response: %w", err)
	}

	if resp_recipients, err = rest.Get(a.client, a.url+a.resources[resourceType].Path+"/"+id+"/recipients"); err != nil {
		return &e, fmt.Errorf("unable to get recipients for object with ID %s: %w", id, err)
	}

	if !resp_recipients.Success() {
		return &e, ResponseErr{StatusCode: resp_recipients.StatusCode, Message: "Failed to get recipients for object", Data: resp_recipients.Body}
	}

	var tempRecipients tempRecipients
	err = json.Unmarshal(resp_recipients.Body, &tempRecipients)
	if err != nil {
		return &e, fmt.Errorf("unable to unmarshal recipients response: %w", err)
	}

	e.Recipients = tempRecipients.Recipients

	return &e, err
}

// UPDATE updates a given automation object
func (a Client) UPDATE(resourceType ResourceType, id string, data []byte) (err error) {
	if id == "" {
		return fmt.Errorf("id must be non empty")
	}

	ds, err := a.GET(resourceType, id)
	if err != nil {
		return fmt.Errorf("unable to get object with ID %s: %w", id, err)
	}

	dsData := tempRecipients{}
	if json.Unmarshal(data, &dsData) != nil {
		return fmt.Errorf("unable to unmarshal data")
	}

	dsRemote := tempRecipients{}
	dsRemote.Recipients = ds.Recipients

	var toAdd tempRecipients
	for _, rec := range dsData.Recipients {
		if !dsRemote.Contains(tempRecipient{Recipient: *rec}) {
			toAdd.Recipients = append(toAdd.Recipients, rec)
		}
	}

	toAddData, err := json.Marshal(toAdd)
	if err != nil {
		return fmt.Errorf("unable to marshal toAdd: %w", err)
	}

	var toRemove IDs
	for _, rec := range dsRemote.Recipients {
		if !dsData.Contains(tempRecipient{Recipient: *rec}) {
			toRemove.Ids = append(toRemove.Ids, rec.Id)
		}
	}

	toRemoveData, err := json.Marshal(toRemove)
	if err != nil {
		return fmt.Errorf("unable to marshal toRemove: %w", err)
	}

	resp_add, err := rest.Post(a.client, a.url+a.resources[resourceType].Path+"/"+id+"/recipients/add", toAddData)
	if err != nil {
		return fmt.Errorf("unable to update object, add recipientes, with ID %s: %w", id, err)
	}

	resp_remove, err := rest.Post(a.client, a.url+a.resources[resourceType].Path+"/"+id+"/recipients/remove", toRemoveData)
	if err != nil {
		return fmt.Errorf("unable to update object, remove recipients, with ID %s: %w", id, err)
	}

	if (resp_add.Success() || len(toAdd.Recipients) == 0) && (resp_remove.Success() || len(toRemove.Ids) == 0) {
		return nil
	}

	return fmt.Errorf("unable to update object with ID %s: %w", id, err)
}

// INSERT creates a given document object
func (a Client) INSERT(resourceType ResourceType, data []byte) (id string, err error) {
	resp, err := rest.Post(a.client, a.url+a.resources[resourceType].Path, data)
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

// DELETE removes a given automation object by ID
func (a Client) DELETE(resourceType ResourceType, id string) (err error) {
	if id == "" {
		return fmt.Errorf("id must be non empty")
	}

	if err != nil {
		return fmt.Errorf("unable to get object with ID %s: %w", id, err)
	}

	var urlParams = make(map[string]string)

	if err = rest.DeleteConfig(a.client, a.url+a.resources[resourceType].Path, id, urlParams); err != nil {
		return fmt.Errorf("unable to delete object with ID %s: %w", id, err)
	}

	return nil
}
