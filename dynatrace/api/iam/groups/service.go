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
	groups "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/iam/groups/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings"
	rest2 "github.com/dynatrace/dynatrace-configuration-as-code-core/api/rest"
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
	credentials *rest.Credentials
}

func Service(clientSet rest.ClientSet) (settings.CRUDService[*groups.Group], error) {
	return &GroupServiceClient{credentials: clientSet.Credentials()}, nil
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
func (me *GroupServiceClient) Create(ctx context.Context, group *groups.Group) (*api.Stub, error) {
	client := rest.NewIAMClient(ctx, me.credentials)
	response, err := client.POST(ctx, fmt.Sprintf("/iam/v1/accounts/%s/groups", client.AccountID()), []*groups.Group{group}, rest2.RequestOptions{})
	if err != nil {
		return nil, err
	}

	responseGroups := []*ListGroup{}
	if err = json.Unmarshal(response.Data, &responseGroups); err != nil {
		return nil, err
	}
	groupID := responseGroups[0].UUID
	groupName := responseGroups[0].Name

	if len(group.Permissions) > 0 {
		if _, err = client.PUT(ctx, fmt.Sprintf("/iam/v1/accounts/%s/groups/%s/permissions", client.AccountID(), groupID), group.Permissions, rest2.RequestOptions{}); err != nil {
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

func (me *GroupServiceClient) Update(ctx context.Context, uuid string, group *groups.Group) error {
	client := rest.NewIAMClient(ctx, me.credentials)
	if _, err := client.PUT(ctx, fmt.Sprintf("/iam/v1/accounts/%s/groups/%s", client.AccountID(), uuid), group, rest2.RequestOptions{}); err != nil {
		return err
	}

	permissions := []*groups.Permission{}

	if len(group.Permissions) > 0 {
		permissions = group.Permissions
	}
	if _, err := client.PUT(ctx, fmt.Sprintf("/iam/v1/accounts/%s/groups/%s/permissions", client.AccountID(), uuid), permissions, rest2.RequestOptions{}); err != nil {
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

func (me *GroupServiceClient) List(ctx context.Context) (api.Stubs, error) {
	client := rest.NewIAMClient(ctx, me.credentials)
	var groupStubs ListGroupsResponse
	response, err := client.GET(ctx, fmt.Sprintf("/iam/v1/accounts/%s/groups", client.AccountID()), rest2.RequestOptions{})
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(response.Data, &groupStubs); err != nil {
		return nil, err
	}

	var stubs api.Stubs
	for _, elem := range groupStubs.Items {
		stubs = append(stubs, &api.Stub{ID: elem.UUID, Name: elem.Name})
	}
	return stubs, nil
}

func (me *GroupServiceClient) Get(ctx context.Context, id string, v *groups.Group) (err error) {
	var groupStub ListGroup
	client := rest.NewIAMClient(ctx, me.credentials)
	response, err := client.GET(ctx, fmt.Sprintf("/iam/v1/accounts/%s/groups/%s/permissions", client.AccountID(), id), rest2.RequestOptions{})
	if err != nil {
		return err
	}
	if err = json.Unmarshal(response.Data, &groupStub); err != nil {
		return err
	}

	v.Name = groupStub.Name
	v.Description = groupStub.Description
	v.FederatedAttributeValues = groupStub.FederatedAttributeValues
	v.Permissions = groupStub.Permissions
	return nil
}

func (me *GroupServiceClient) Delete(ctx context.Context, id string) error {
	client := rest.NewIAMClient(ctx, me.credentials)
	_, err := client.DELETE(ctx, fmt.Sprintf("/iam/v1/accounts/%s/groups/%s", client.AccountID(), id), rest2.RequestOptions{})

	// data sources MAY have cached a list of group IDs
	// Updating the (publicly available) revision signals to them that either a CREATE or DELETE has happened since
	revisionLock.Lock()
	defer revisionLock.Unlock()
	revision = uuid.NewString()

	return err
}
