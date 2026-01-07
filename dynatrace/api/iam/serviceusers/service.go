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

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/iam"
	serviceusers "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/iam/serviceusers/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings"
)

type iamClientGetter interface {
	New() iam.IAMClient
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

func (me *iamClientGetterImp) New() iam.IAMClient {
	return iam.NewIAMClient(me)
}

type serviceUserServiceClient struct {
	iamClientGetter iamClientGetter
	accountID       string
	endpointURL     string
}

func Service(credentials *rest.Credentials) settings.CRUDService[*serviceusers.ServiceUser] {
	return &serviceUserServiceClient{
		iamClientGetter: &iamClientGetterImp{
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

func (me *serviceUserServiceClient) SchemaID() string {
	return "accounts:iam:serviceusers"
}

// createResponse represents the response from creating a service user
type createResponse struct {
	UID   string `json:"uid"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

func (me *serviceUserServiceClient) Create(ctx context.Context, serviceUser *serviceusers.ServiceUser) (*api.Stub, error) {
	responseBytes, err := me.iamClientGetter.New().POST(ctx, fmt.Sprintf("%s/iam/v1/accounts/%s/service-users", me.endpointURL, me.accountID), serviceUser, 201, false)
	if err != nil {
		return nil, err
	}

	var response createResponse
	if err := json.Unmarshal(responseBytes, &response); err != nil {
		return nil, err
	}

	if err := me.updateGroupAssignments(ctx, response.Email, serviceUser.Groups); err != nil {
		deleteErr := me.Delete(ctx, response.UID)
		if deleteErr != nil {
			return nil, fmt.Errorf("failed to create service user: %v; additionally failed to clean up service user: %v", err, deleteErr)
		}
		return nil, err
	}

	return &api.Stub{ID: response.UID, Name: response.Name}, nil
}

// getServiceUserResponse represents the response from getting a service user
type getServiceUserResponse struct {
	UID         string `json:"uid"`
	Email       string `json:"email"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (me *serviceUserServiceClient) Get(ctx context.Context, id string, v *serviceusers.ServiceUser) error {
	responseBytes, err := me.iamClientGetter.New().GET(ctx, fmt.Sprintf("%s/iam/v1/accounts/%s/service-users/%s", me.endpointURL, me.accountID, id), 200, false)
	if err != nil {
		return err
	}

	var response getServiceUserResponse
	if err = json.Unmarshal(responseBytes, &response); err != nil {
		return err
	}

	groups, err := me.getUserGroups(ctx, response.Email)
	if err != nil {
		return err
	}

	v.Email = response.Email
	v.Name = response.Name
	v.Description = response.Description
	v.Groups = groups

	return nil
}

// groupStub represents a group membership
type groupStub struct {
	UUID string `json:"uuid"`
}

// getUserPartialResponse represents the partial response from getting user information.
type getUserPartialResponse struct {
	Groups []*groupStub `json:"groups"`
}

func (me *serviceUserServiceClient) getUserGroups(ctx context.Context, email string) ([]string, error) {
	responseBytes, err := me.iamClientGetter.New().GET(ctx, fmt.Sprintf("%s/iam/v1/accounts/%s/users/%s", me.endpointURL, me.accountID, email), 200, false)
	if err != nil {
		return nil, err
	}

	var response getUserPartialResponse
	if err = json.Unmarshal(responseBytes, &response); err != nil {
		return nil, err
	}

	groups := make([]string, 0, len(response.Groups))
	for _, group := range response.Groups {
		groups = append(groups, group.UUID)
	}

	return groups, nil
}

func (me *serviceUserServiceClient) Update(ctx context.Context, id string, serviceUser *serviceusers.ServiceUser) error {
	// Update the service user details
	if _, err := me.iamClientGetter.New().PUT(ctx, fmt.Sprintf("%s/iam/v1/accounts/%s/service-users/%s", me.endpointURL, me.accountID, id), serviceUser, 200, false); err != nil {
		return err
	}

	// Update group assignments
	return me.updateGroupAssignments(ctx, serviceUser.Email, serviceUser.Groups)
}

func (me *serviceUserServiceClient) updateGroupAssignments(ctx context.Context, serviceUserEmail string, groups []string) error {
	// no groups must be represented as an empty list rather than nil
	if groups == nil {
		groups = []string{}
	}
	_, err := me.iamClientGetter.New().PUT(ctx, fmt.Sprintf("%s/iam/v1/accounts/%s/users/%s/groups", me.endpointURL, me.accountID, serviceUserEmail), groups, 200, false)
	return err
}

// serviceUserStub represents a service user in the list response
type serviceUserStub struct {
	UID   string `json:"uid"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

// listServiceUsersResponse represents the paginated response from listing service users
type listServiceUsersResponse struct {
	Count    int                `json:"count"`
	Items    []*serviceUserStub `json:"results"`
	NextPage string             `json:"nextPageKey,omitempty"`
}

func (me *serviceUserServiceClient) List(ctx context.Context) (api.Stubs, error) {
	client := me.iamClientGetter.New()

	var stubs api.Stubs
	url := fmt.Sprintf("%s/iam/v1/accounts/%s/service-users", me.endpointURL, me.accountID)

	for {
		responseBytes, err := client.GET(ctx, url, 200, false)
		if err != nil {
			return nil, err
		}

		var response listServiceUsersResponse
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

func (me *serviceUserServiceClient) Delete(ctx context.Context, uid string) error {
	_, err := me.iamClientGetter.New().DELETE(ctx, fmt.Sprintf("%s/iam/v1/accounts/%s/service-users/%s", me.endpointURL, me.accountID, uid), 200, false)
	return err
}
