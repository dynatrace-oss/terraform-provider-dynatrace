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

package groups

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/iam"
	groups "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/iam/groups/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings"
	"github.com/google/uuid"
)

// data sources MAY have cached a list of group IDs
// Updating the (publicly available) revision signals to them that either a CREATE or DELETE has happened since
var revision = uuid.NewString()
var revisionLock = sync.Mutex{}

func GetRevision() string {
	revisionLock.Lock()
	defer revisionLock.Unlock()
	return revision
}

type GroupServiceClient struct {
	clientID     string
	accountID    string
	clientSecret string
	tokenURL     string
	endpointURL  string
}

func (me *GroupServiceClient) ClientID() string {
	return me.clientID
}

func (me *GroupServiceClient) AccountID() string {
	return me.accountID
}

func (me *GroupServiceClient) ClientSecret() string {
	return me.clientSecret
}

func (me *GroupServiceClient) TokenURL() string {
	return me.tokenURL
}

func (me *GroupServiceClient) EndpointURL() string {
	return me.endpointURL
}

func NewGroupService(clientID string, accountID string, clientSecret string, tokenURL string, endpointURL string) settings.CRUDService[*groups.Group] {
	return &GroupServiceClient{clientID: clientID, accountID: accountID, clientSecret: clientSecret, tokenURL: tokenURL, endpointURL: endpointURL}
}

func Service(credentials *rest.Credentials) settings.CRUDService[*groups.Group] {
	return &GroupServiceClient{clientID: credentials.IAM.ClientID, accountID: credentials.IAM.AccountID, clientSecret: credentials.IAM.ClientSecret, tokenURL: credentials.IAM.TokenURL, endpointURL: credentials.IAM.EndpointURL}
}

func (me *GroupServiceClient) SchemaID() string {
	return "accounts:iam:groups"
}

func (me *GroupServiceClient) Name() string {
	return me.SchemaID()
}

// TODO ... keep group cache up to date
// UUID                     string             `json:"uuid"`
// Name                     string             `json:"name"`
// Description              string             `json:"description"`
// FederatedAttributeValues []string           `json:"federatedAttributeValues"`
// Permissions              groups.Permissions `json:"permissions"`
func (me *GroupServiceClient) Create(ctx context.Context, group *groups.Group, m any) (*api.Stub, error) {
	var err error
	var responseBytes []byte

	client := iam.NewIAMClient(ctx, me)
	if responseBytes, err = client.POST(ctx, fmt.Sprintf("%s/iam/v1/accounts/%s/groups", me.endpointURL, me.AccountID()), []*groups.Group{group}, 201, false); err != nil {
		return nil, err
	}

	responseGroups := []*ListGroup{}
	if err = json.Unmarshal(responseBytes, &responseGroups); err != nil {
		return nil, err
	}
	groupID := responseGroups[0].UUID
	groupName := responseGroups[0].Name

	if len(group.Permissions) > 0 {
		if _, err = client.PUT(ctx, fmt.Sprintf("%s/iam/v1/accounts/%s/groups/%s/permissions", me.endpointURL, me.AccountID(), groupID), group.Permissions, 200, false); err != nil {
			return nil, err
		}
	}

	// data sources MAY have cached a list of group IDs
	// Updating the (publicly available) revision signals to them that either a CREATE or DELETE has happened since
	revisionLock.Lock()
	defer revisionLock.Unlock()
	revision = uuid.NewString()

	return &api.Stub{ID: groupID, Name: groupName}, nil
}

func (me *GroupServiceClient) Update(ctx context.Context, uuid string, group *groups.Group, m any) error {
	var err error

	client := iam.NewIAMClient(ctx, me)
	if _, err = client.PUT(ctx, fmt.Sprintf("%s/iam/v1/accounts/%s/groups/%s", me.endpointURL, me.AccountID(), uuid), group, 200, false); err != nil {
		return err
	}

	permissions := []*groups.Permission{}

	if len(group.Permissions) > 0 {
		permissions = group.Permissions
	}
	if _, err = client.PUT(ctx, fmt.Sprintf("%s/iam/v1/accounts/%s/groups/%s/permissions", me.endpointURL, me.AccountID(), uuid), permissions, 200, false); err != nil {
		return err
	}

	return nil
}

type ListGroup struct {
	UUID                     string             `json:"uuid"`
	Name                     string             `json:"name"`
	Description              string             `json:"description"`
	FederatedAttributeValues []string           `json:"federatedAttributeValues"`
	Permissions              groups.Permissions `json:"permissions"`
}

type ListGroupsResponse struct {
	Count int          `json:"count:"`
	Items []*ListGroup `json:"items"`
}

func (me *GroupServiceClient) List(ctx context.Context, m any) (api.Stubs, error) {
	client := iam.NewIAMClient(ctx, me)
	var groupStubs ListGroupsResponse
	accountID := me.AccountID()
	if err := iam.GET(client, ctx, fmt.Sprintf("%s/iam/v1/accounts/%s/groups", me.endpointURL, accountID), 200, false, &groupStubs); err != nil {
		return nil, err
	}

	var stubs api.Stubs
	for _, elem := range groupStubs.Items {
		stubs = append(stubs, &api.Stub{ID: elem.UUID, Name: elem.Name})
	}
	return stubs, nil
}

func (me *GroupServiceClient) Get(ctx context.Context, id string, v *groups.Group, m any) (err error) {
	var groupStub ListGroup
	accountID := me.AccountID()
	client := iam.NewIAMClient(ctx, me)
	if err = iam.GET(client, ctx, fmt.Sprintf("%s/iam/v1/accounts/%s/groups/%s/permissions", me.endpointURL, accountID, id), 200, false, &groupStub); err != nil {
		return err
	}

	v.Name = groupStub.Name
	v.Description = groupStub.Description
	v.FederatedAttributeValues = groupStub.FederatedAttributeValues
	v.Permissions = groupStub.Permissions
	return nil
}

func (me *GroupServiceClient) Delete(ctx context.Context, id string, m any) error {
	_, err := iam.NewIAMClient(ctx, me).DELETE(ctx, fmt.Sprintf("%s/iam/v1/accounts/%s/groups/%s", me.endpointURL, me.AccountID(), id), 200, false)

	// data sources MAY have cached a list of group IDs
	// Updating the (publicly available) revision signals to them that either a CREATE or DELETE has happened since
	revisionLock.Lock()
	defer revisionLock.Unlock()
	revision = uuid.NewString()

	return err
}
