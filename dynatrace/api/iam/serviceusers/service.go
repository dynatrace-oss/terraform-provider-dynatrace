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

type iamClientGetter interface {
	NewIAMClient() iam.IAMClient
}

type iamClientGetterImp struct {
	clientID     string
	accountID    string
	clientSecret string
	tokenURL     string
	endpointURL  string
}

func (me *iamClientGetterImp) ClientID() string {
	return me.clientID
}

func (me *iamClientGetterImp) AccountID() string {
	return me.accountID
}

func (me *iamClientGetterImp) ClientSecret() string {
	return me.clientSecret
}

func (me *iamClientGetterImp) TokenURL() string {
	return me.tokenURL
}

func (me *iamClientGetterImp) EndpointURL() string {
	return me.endpointURL
}

func (me *iamClientGetterImp) NewIAMClient() iam.IAMClient {
	return iam.NewIAMClient(me)
}

type ServiceUserServiceClient struct {
	iamClient   iamClientGetter
	accountID   string
	endpointURL string
}

func Service(credentials *rest.Credentials) settings.CRUDService[*serviceusers.ServiceUser] {
	return &ServiceUserServiceClient{
		iamClient: &iamClientGetterImp{
			clientID:     credentials.IAM.ClientID,
			accountID:    credentials.IAM.AccountID,
			clientSecret: credentials.IAM.ClientSecret,
			tokenURL:     credentials.IAM.TokenURL,
			endpointURL:  credentials.IAM.EndpointURL,
		},
		accountID:   credentials.IAM.AccountID,
		endpointURL: credentials.IAM.EndpointURL,
	}
}

func (me *ServiceUserServiceClient) SchemaID() string {
	return "accounts:iam:serviceusers"
}

// CreateResponse represents the response from creating a service user
type CreateResponse struct {
	UID   string `json:"uid"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

func (me *ServiceUserServiceClient) Create(ctx context.Context, serviceUser *serviceusers.ServiceUser) (*api.Stub, error) {
	responseBytes, err := me.iamClient.NewIAMClient().POST(ctx, fmt.Sprintf("%s/iam/v1/accounts/%s/service-users", me.endpointURL, me.accountID), serviceUser, 201, false)
	if err != nil {
		return nil, err
	}

	var response CreateResponse
	if err := json.Unmarshal(responseBytes, &response); err != nil {
		return nil, err
	}

	if err := me.updateGroupAssignments(ctx, serviceUser.Email, serviceUser.Groups); err != nil {
		deleteErr := me.Delete(ctx, response.UID)
		if deleteErr != nil {
			return nil, fmt.Errorf("failed to create service user: %v; additionally failed to clean up service user: %v", err, deleteErr)
		}
		return nil, err
	}

	return &api.Stub{ID: response.UID, Name: response.Name}, nil
}

// GroupStub represents a group membership
type GroupStub struct {
	GroupName string `json:"groupName"`
	UUID      string `json:"uuid"`
}

// GetServiceUserResponse represents the response from getting a service user
type GetServiceUserResponse struct {
	UID         string       `json:"uid"`
	Email       string       `json:"email"`
	Name        string       `json:"name"`
	Description string       `json:"description"`
	Groups      []*GroupStub `json:"groups"`
}

func (me *ServiceUserServiceClient) Get(ctx context.Context, uid string, v *serviceusers.ServiceUser) error {
	responseBytes, err := me.iamClient.NewIAMClient().GET(ctx, fmt.Sprintf("%s/iam/v1/accounts/%s/service-users/%s", me.endpointURL, me.accountID, uid), 200, false)
	if err != nil {
		return err
	}

	var response GetServiceUserResponse
	if err = json.Unmarshal(responseBytes, &response); err != nil {
		return err
	}

	v.UID = response.UID
	v.Email = response.Email
	v.Name = response.Name
	v.Description = response.Description
	v.Groups = []string{}
	for _, group := range response.Groups {
		v.Groups = append(v.Groups, group.UUID)
	}

	return nil
}

func (me *ServiceUserServiceClient) Update(ctx context.Context, uid string, serviceUser *serviceusers.ServiceUser) error {
	// Update the service user details
	if _, err := me.iamClient.NewIAMClient().PUT(ctx, fmt.Sprintf("%s/iam/v1/accounts/%s/service-users/%s", me.endpointURL, me.accountID, uid), serviceUser, 200, false); err != nil {
		return err
	}

	// Update group assignments
	return me.updateGroupAssignments(ctx, serviceUser.Email, serviceUser.Groups)
}

func (me *ServiceUserServiceClient) updateGroupAssignments(ctx context.Context, serviceUserEmail string, groups []string) error {
	_, err := me.iamClient.NewIAMClient().PUT(ctx, fmt.Sprintf("%s/iam/v1/accounts/%s/users/%s/groups", me.endpointURL, me.accountID, serviceUserEmail), groups, 200, false)
	return err
}

// ServiceUserStub represents a service user in the list response
type ServiceUserStub struct {
	UID   string `json:"uid"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

// ListServiceUsersResponse represents the paginated response from listing service users
type ListServiceUsersResponse struct {
	Count    int                `json:"count"`
	Items    []*ServiceUserStub `json:"results"`
	NextPage string             `json:"nextPageKey,omitempty"`
}

func (me *ServiceUserServiceClient) List(ctx context.Context) (api.Stubs, error) {
	client := me.iamClient.NewIAMClient()

	var stubs api.Stubs
	url := fmt.Sprintf("%s/iam/v1/accounts/%s/service-users", me.endpointURL, me.accountID)

	for {
		responseBytes, err := client.GET(ctx, url, 200, false)
		if err != nil {
			return nil, err
		}

		var response ListServiceUsersResponse
		if err := json.Unmarshal(responseBytes, &response); err != nil {
			return nil, err
		}

		for _, item := range response.Items {
			stubs = append(stubs, &api.Stub{ID: item.UID, Name: item.Name})
		}

		// Handle pagination
		if response.NextPage == "" {
			break
		}
		url = fmt.Sprintf("%s/iam/v1/accounts/%s/service-users?nextPageKey=%s", me.endpointURL, me.accountID, response.NextPage)
	}

	return stubs, nil
}

func (me *ServiceUserServiceClient) Delete(ctx context.Context, uid string) error {
	_, err := me.iamClient.NewIAMClient().DELETE(ctx, fmt.Sprintf("%s/iam/v1/accounts/%s/service-users/%s", me.endpointURL, me.accountID, uid), 200, false)
	if err != nil {
		// ignore error if service user does not exist
		if strings.Contains(err.Error(), fmt.Sprintf("User %s does not exist", uid)) {
			return nil
		}
		return err
	}
	return nil
}
