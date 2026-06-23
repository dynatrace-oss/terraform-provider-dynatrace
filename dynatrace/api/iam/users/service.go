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

package users

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/iam"
	users "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/iam/users/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings"
	rest2 "github.com/dynatrace/dynatrace-configuration-as-code-core/api/rest"
)

type UserServiceClient struct {
	clientID     string
	accountID    string
	clientSecret string
	tokenURL     string
	endpointURL  string
}

func (me *UserServiceClient) ClientID() string {
	return me.clientID
}

func (me *UserServiceClient) AccountID() string {
	return me.accountID
}

func (me *UserServiceClient) ClientSecret() string {
	return me.clientSecret
}

func (me *UserServiceClient) TokenURL() string {
	return me.tokenURL
}

func (me *UserServiceClient) EndpointURL() string {
	return me.endpointURL
}

func NewUserService(clientID string, accountID string, clientSecret string, tokenURL string, endpointURL string) *UserServiceClient {
	return &UserServiceClient{clientID: clientID, accountID: accountID, clientSecret: clientSecret, tokenURL: tokenURL, endpointURL: endpointURL}
}

func Service(credentials *rest.Credentials) settings.CRUDService[*users.User] {

	return &UserServiceClient{clientID: credentials.IAM.ClientID, accountID: credentials.IAM.AccountID, clientSecret: credentials.IAM.ClientSecret, tokenURL: credentials.IAM.TokenURL, endpointURL: credentials.IAM.EndpointURL}
}

func (me *UserServiceClient) SchemaID() string {
	return "accounts:iam:users"
}

func (me *UserServiceClient) Create(ctx context.Context, user *users.User) (*api.Stub, error) {
	client := iam.NewIAMClient(ctx, me)
	if _, err := client.POST(ctx, fmt.Sprintf("/iam/v1/accounts/%s/users", me.AccountID()), user, rest2.RequestOptions{}); err != nil {
		return nil, err
	}

	groups := []string{}
	if len(user.Groups) > 0 {
		groups = user.Groups
	}
	if _, err := client.PUT(ctx, fmt.Sprintf("/iam/v1/accounts/%s/users/%s/groups", me.AccountID(), user.Email), groups, rest2.RequestOptions{}); err != nil {
		return nil, err
	}

	return &api.Stub{ID: user.Email, Name: user.Email}, nil
}

type GroupStub struct {
	GroupName string `json:"groupName"`
	UUID      string `json:"uuid"`
}

type GetUserGroupsResponse struct {
	Groups []*GroupStub
	UID    string `json:"uid"`
}

func (me *UserServiceClient) Get(ctx context.Context, email string, v *users.User) error {
	client := iam.NewIAMClient(ctx, me)

	response, err := client.GET(ctx, fmt.Sprintf("/iam/v1/accounts/%s/users/%s", me.AccountID(), email), rest2.RequestOptions{})
	if err != nil {
		if strings.Contains(err.Error(), fmt.Sprintf("User %s not found", email)) {
			return rest.Error{Code: 404, Message: err.Error()}
		}
		return err
	}

	var userGroups GetUserGroupsResponse
	if err = json.Unmarshal(response.Data, &userGroups); err != nil {
		return err
	}
	v.Email = email
	v.Groups = []string{}
	v.UID = userGroups.UID
	for _, group := range userGroups.Groups {
		v.Groups = append(v.Groups, group.UUID)
	}

	return nil
}

func (me *UserServiceClient) Update(ctx context.Context, email string, user *users.User) error {
	groups := []string{}
	if len(user.Groups) > 0 {
		groups = user.Groups
	}
	if _, err := iam.NewIAMClient(ctx, me).PUT(ctx, fmt.Sprintf("/iam/v1/accounts/%s/users/%s/groups", me.AccountID(), user.Email), groups, rest2.RequestOptions{}); err != nil {
		return err
	}

	return nil
}

type UserStub struct {
	UID   string `json:"uid"`
	Email string `json:"email"`
}

type ListUsersResponse struct {
	Count int        `json:"count:"`
	Items []UserStub `json:"items"`
}

func (me *UserServiceClient) List(ctx context.Context) (api.Stubs, error) {
	response, err := iam.NewIAMClient(ctx, me).GET(ctx, fmt.Sprintf("/iam/v1/accounts/%s/users", me.AccountID()), rest2.RequestOptions{})
	if err != nil {
		return nil, err
	}

	var listResp ListUsersResponse
	if err = json.Unmarshal(response.Data, &listResp); err != nil {
		return nil, err
	}
	var stubs api.Stubs
	for _, item := range listResp.Items {
		stubs = append(stubs, &api.Stub{ID: item.Email, Name: item.Email})
	}
	return stubs, nil
}

func (me *UserServiceClient) Delete(ctx context.Context, email string) error {
	_, err := iam.NewIAMClient(ctx, me).DELETE(ctx, fmt.Sprintf("/iam/v1/accounts/%s/users/%s", me.AccountID(), email), rest2.RequestOptions{})
	if err != nil && strings.Contains(err.Error(), fmt.Sprintf("User %s not found", email)) {
		return nil
	}
	return err
}
