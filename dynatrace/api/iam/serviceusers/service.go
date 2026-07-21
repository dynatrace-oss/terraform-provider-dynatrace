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
	"net/url"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api"
	serviceusers "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/iam/serviceusers/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings"
	rest2 "github.com/dynatrace/dynatrace-configuration-as-code-core/api/rest"
)

type ServiceUserService interface {
	settings.CRUDService[*serviceusers.ServiceUser]
	GetAll(ctx context.Context) ([]ServiceUserStub, error)
}

type serviceUserServiceClient struct {
	client rest.IAMClient
}

func NewService(clientSet rest.ClientSet) (ServiceUserService, error) {
	iamClient, err := clientSet.IAMClient()
	if err != nil {
		return nil, err
	}
	return &serviceUserServiceClient{client: iamClient}, nil
}

func Service(clientSet rest.ClientSet) (settings.CRUDService[*serviceusers.ServiceUser], error) {
	return NewService(clientSet)
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
	response, err := me.client.POST(ctx, fmt.Sprintf("/iam/v1/accounts/%s/service-users", me.client.AccountID()), serviceUser, rest2.RequestOptions{})
	if err != nil {
		return nil, err
	}

	var createResp createResponse
	if err := json.Unmarshal(response.Data, &createResp); err != nil {
		return nil, err
	}

	if err := me.updateGroupAssignments(ctx, createResp.Email, serviceUser.Groups); err != nil {
		deleteErr := me.Delete(ctx, createResp.UID)
		if deleteErr != nil {
			return nil, fmt.Errorf("failed to create service user: %v; additionally failed to clean up service user: %v", err, deleteErr)
		}
		return nil, err
	}

	return &api.Stub{ID: createResp.UID, Name: createResp.Name}, nil
}

// getServiceUserResponse represents the response from getting a service user
type getServiceUserResponse struct {
	UID         string `json:"uid"`
	Email       string `json:"email"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (me *serviceUserServiceClient) Get(ctx context.Context, id string, v *serviceusers.ServiceUser) error {
	response, err := me.client.GET(ctx, fmt.Sprintf("/iam/v1/accounts/%s/service-users/%s", me.client.AccountID(), id), rest2.RequestOptions{})
	if err != nil {
		return err
	}

	var getResp getServiceUserResponse
	if err = json.Unmarshal(response.Data, &getResp); err != nil {
		return err
	}

	groups, err := me.getUserGroups(ctx, getResp.Email)
	if err != nil {
		return err
	}

	v.Email = getResp.Email
	v.Name = getResp.Name
	v.Description = getResp.Description
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
	response, err := me.client.GET(ctx, fmt.Sprintf("/iam/v1/accounts/%s/users/%s", me.client.AccountID(), email), rest2.RequestOptions{})
	if err != nil {
		return nil, err
	}

	var getResp getUserPartialResponse
	if err = json.Unmarshal(response.Data, &getResp); err != nil {
		return nil, err
	}

	groups := make([]string, 0, len(getResp.Groups))
	for _, group := range getResp.Groups {
		groups = append(groups, group.UUID)
	}

	return groups, nil
}

func (me *serviceUserServiceClient) Update(ctx context.Context, id string, serviceUser *serviceusers.ServiceUser) error {
	// Update the service user details
	if _, err := me.client.PUT(ctx, fmt.Sprintf("/iam/v1/accounts/%s/service-users/%s", me.client.AccountID(), id), serviceUser, rest2.RequestOptions{}); err != nil {
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
	_, err := me.client.PUT(ctx, fmt.Sprintf("/iam/v1/accounts/%s/users/%s/groups", me.client.AccountID(), serviceUserEmail), groups, rest2.RequestOptions{})
	return err
}

// ServiceUserStub represents a service user in the list response
type ServiceUserStub struct {
	UID         string `json:"uid"`
	Email       string `json:"email"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// listServiceUsersResponse represents the paginated response from listing service users
type listServiceUsersResponse struct {
	Count    int               `json:"count"`
	Items    []ServiceUserStub `json:"results"`
	NextPage string            `json:"nextPageKey,omitempty"`
}

func (me *serviceUserServiceClient) GetAll(ctx context.Context) ([]ServiceUserStub, error) {
	var stubs []ServiceUserStub
	endpoint := fmt.Sprintf("/iam/v1/accounts/%s/service-users", me.client.AccountID())
	options := rest2.RequestOptions{}

	for {
		response, err := me.client.GET(ctx, endpoint, options)
		if err != nil {
			return nil, err
		}

		var listResp listServiceUsersResponse
		if err := json.Unmarshal(response.Data, &listResp); err != nil {
			return nil, err
		}

		stubs = append(stubs, listResp.Items...)

		// Handle pagination
		if listResp.NextPage == "" {
			break
		}
		options.QueryParams = url.Values{"nextPageKey": {listResp.NextPage}}
	}

	return stubs, nil
}

func (me *serviceUserServiceClient) List(ctx context.Context) (api.Stubs, error) {
	var stubs api.Stubs
	serviceUserStubs, err := me.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	for _, serviceUserStub := range serviceUserStubs {
		stubs = append(stubs, &api.Stub{
			ID:   serviceUserStub.UID,
			Name: serviceUserStub.Name,
		})
	}

	return stubs, nil
}

func (me *serviceUserServiceClient) Delete(ctx context.Context, uid string) error {
	_, err := me.client.DELETE(ctx, fmt.Sprintf("/iam/v1/accounts/%s/service-users/%s", me.client.AccountID(), uid), rest2.RequestOptions{})
	return err
}
