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

package serviceusers

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/iam"
	serviceusers "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/iam/serviceusers/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings"
)

type ServiceUserServiceClient struct {
	clientID     string
	accountID    string
	clientSecret string
	tokenURL     string
	endpointURL  string
}

func (me *ServiceUserServiceClient) ClientID() string {
	return me.clientID
}

func (me *ServiceUserServiceClient) AccountID() string {
	return me.accountID
}

func (me *ServiceUserServiceClient) ClientSecret() string {
	return me.clientSecret
}

func (me *ServiceUserServiceClient) TokenURL() string {
	return me.tokenURL
}

func (me *ServiceUserServiceClient) EndpointURL() string {
	return me.endpointURL
}

func NewServiceUserService(clientID string, accountID string, clientSecret string, tokenURL string, endpointURL string) *ServiceUserServiceClient {
	return &ServiceUserServiceClient{clientID: clientID, accountID: accountID, clientSecret: clientSecret, tokenURL: tokenURL, endpointURL: endpointURL}
}

func Service(credentials *rest.Credentials) settings.CRUDService[*serviceusers.ServiceUser] {
	return &ServiceUserServiceClient{clientID: credentials.IAM.ClientID, accountID: credentials.IAM.AccountID, clientSecret: credentials.IAM.ClientSecret, tokenURL: credentials.IAM.TokenURL, endpointURL: credentials.IAM.EndpointURL}
}

func (me *ServiceUserServiceClient) SchemaID() string {
	return "accounts:iam:service-users"
}

type CreateServiceUserRequest struct {
	Name        string  `json:"name"`
	Description *string `json:"description,omitempty"`
}

type CreateServiceUserResponse struct {
	ID string `json:"id"`
}

func (me *ServiceUserServiceClient) Create(ctx context.Context, serviceUser *serviceusers.ServiceUser) (*api.Stub, error) {
	var err error
	var responseBytes []byte

	client := iam.NewIAMClient(me)

	// Create the service user
	request := CreateServiceUserRequest{
		Name:        serviceUser.ServiceUserName,
		Description: serviceUser.Description,
	}

	if responseBytes, err = client.POST(ctx, fmt.Sprintf("%s/iam/v1/accounts/%s/service-users", me.endpointURL, strings.TrimPrefix(me.AccountID(), "urn:dtaccount:")), request, 201, false); err != nil {
		return nil, err
	}

	var response CreateServiceUserResponse
	if err = json.Unmarshal(responseBytes, &response); err != nil {
		return nil, err
	}

	// Update groups if provided
	groups := []string{}
	if len(serviceUser.Groups) > 0 {
		groups = serviceUser.Groups
	}
	if _, err = client.PUT(ctx, fmt.Sprintf("%s/iam/v1/accounts/%s/service-users/%s/groups", me.endpointURL, strings.TrimPrefix(me.AccountID(), "urn:dtaccount:"), response.ID), groups, 200, false); err != nil {
		return nil, err
	}

	return &api.Stub{ID: response.ID, Name: serviceUser.ServiceUserName}, nil
}

type GroupStub struct {
	GroupName string `json:"groupName"`
	UUID      string `json:"uuid"`
}

type GetServiceUserResponse struct {
	ID          string       `json:"id"`
	Name        string       `json:"name"`
	Description *string      `json:"description,omitempty"`
	Groups      []*GroupStub `json:"groups,omitempty"`
}

func (me *ServiceUserServiceClient) Get(ctx context.Context, id string, v *serviceusers.ServiceUser) error {
	var err error
	var responseBytes []byte

	client := iam.NewIAMClient(me)

	if responseBytes, err = client.GET(ctx, fmt.Sprintf("%s/iam/v1/accounts/%s/service-users/%s", me.endpointURL, strings.TrimPrefix(me.AccountID(), "urn:dtaccount:"), id), 200, false); err != nil {
		if err != nil && strings.Contains(err.Error(), fmt.Sprintf("Service user %s not found", id)) {
			return rest.Error{Code: 404, Message: err.Error()}
		}
		return err
	}

	var response GetServiceUserResponse
	if err = json.Unmarshal(responseBytes, &response); err != nil {
		return err
	}

	v.ID = response.ID
	v.ServiceUserName = response.Name
	v.Description = response.Description
	v.Groups = []string{}
	for _, group := range response.Groups {
		v.Groups = append(v.Groups, group.UUID)
	}

	return nil
}

func (me *ServiceUserServiceClient) Update(ctx context.Context, id string, serviceUser *serviceusers.ServiceUser) error {
	var err error

	client := iam.NewIAMClient(me)

	// Update service user details
	request := CreateServiceUserRequest{
		Name:        serviceUser.ServiceUserName,
		Description: serviceUser.Description,
	}

	if _, err = client.PUT(ctx, fmt.Sprintf("%s/iam/v1/accounts/%s/service-users/%s", me.endpointURL, strings.TrimPrefix(me.AccountID(), "urn:dtaccount:"), id), request, 200, false); err != nil {
		return err
	}

	// Update groups
	groups := []string{}
	if len(serviceUser.Groups) > 0 {
		groups = serviceUser.Groups
	}
	if _, err = client.PUT(ctx, fmt.Sprintf("%s/iam/v1/accounts/%s/service-users/%s/groups", me.endpointURL, strings.TrimPrefix(me.AccountID(), "urn:dtaccount:"), id), groups, 200, false); err != nil {
		return err
	}

	return nil
}

type ServiceUserStub struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type ListServiceUsersResponse struct {
	Count int               `json:"count"`
	Items []ServiceUserStub `json:"items"`
}

func (me *ServiceUserServiceClient) List(ctx context.Context) (api.Stubs, error) {
	var err error
	var responseBytes []byte

	if responseBytes, err = iam.NewIAMClient(me).GET(ctx, fmt.Sprintf("%s/iam/v1/accounts/%s/service-users", me.endpointURL, strings.TrimPrefix(me.AccountID(), "urn:dtaccount:")), 200, false); err != nil {
		return nil, err
	}

	var response ListServiceUsersResponse
	if err = json.Unmarshal(responseBytes, &response); err != nil {
		return nil, err
	}
	var stubs api.Stubs
	for _, item := range response.Items {
		stubs = append(stubs, &api.Stub{ID: item.ID, Name: item.Name})
	}
	return stubs, nil
}

func (me *ServiceUserServiceClient) Delete(ctx context.Context, id string) error {
	_, err := iam.NewIAMClient(me).DELETE(ctx, fmt.Sprintf("%s/iam/v1/accounts/%s/service-users/%s", me.endpointURL, strings.TrimPrefix(me.AccountID(), "urn:dtaccount:"), id), 200, false)
	if err != nil && strings.Contains(err.Error(), fmt.Sprintf("Service user %s not found", id)) {
		return nil
	}
	return err
}
