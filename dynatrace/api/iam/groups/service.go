package groups

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
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

func Service(credentials *settings.Credentials) settings.CRUDService[*groups.Group] {
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
func (me *GroupServiceClient) Create(ctx context.Context, group *groups.Group) (*api.Stub, error) {
	var err error
	var responseBytes []byte

	client := iam.NewIAMClient(me)
	if responseBytes, err = client.POST(ctx, fmt.Sprintf("%s/iam/v1/accounts/%s/groups", me.endpointURL, strings.TrimPrefix(me.AccountID(), "urn:dtaccount:")), []*groups.Group{group}, 201, false); err != nil {
		return nil, err
	}

	responseGroups := []*ListGroup{}
	if err = json.Unmarshal(responseBytes, &responseGroups); err != nil {
		return nil, err
	}
	groupID := responseGroups[0].UUID
	groupName := responseGroups[0].Name

	if len(group.Permissions) > 0 {
		if _, err = client.PUT(ctx, fmt.Sprintf("%s/iam/v1/accounts/%s/groups/%s/permissions", me.endpointURL, strings.TrimPrefix(me.AccountID(), "urn:dtaccount:"), groupID), group.Permissions, 200, false); err != nil {
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
	var err error

	client := iam.NewIAMClient(me)
	if _, err = client.PUT(ctx, fmt.Sprintf("%s/iam/v1/accounts/%s/groups/%s", me.endpointURL, strings.TrimPrefix(me.AccountID(), "urn:dtaccount:"), uuid), group, 200, false); err != nil {
		return err
	}

	permissions := []*groups.Permission{}

	if len(group.Permissions) > 0 {
		permissions = group.Permissions
	}
	if _, err = client.PUT(ctx, fmt.Sprintf("%s/iam/v1/accounts/%s/groups/%s/permissions", me.endpointURL, strings.TrimPrefix(me.AccountID(), "urn:dtaccount:"), uuid), permissions, 200, false); err != nil {
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

var cachedGroupStubs []*ListGroup
var groupStubMutex sync.Mutex

func (me *GroupServiceClient) List(ctx context.Context) (api.Stubs, error) {
	groupStubMutex.Lock()
	defer groupStubMutex.Unlock()

	if cachedGroupStubs != nil {
		var stubs api.Stubs
		for _, elem := range cachedGroupStubs {
			stubs = append(stubs, &api.Stub{ID: elem.UUID, Name: elem.Name})
		}
		return stubs, nil
	}

	groupStubs, err := me.listUnguarded(ctx)
	if err != nil {
		return nil, err
	}
	var stubs api.Stubs
	for _, elem := range groupStubs {
		stubs = append(stubs, &api.Stub{ID: elem.UUID, Name: elem.Name})
	}
	return stubs, nil
}

func (me *GroupServiceClient) list(ctx context.Context) ([]*ListGroup, error) {
	groupStubMutex.Lock()
	defer groupStubMutex.Unlock()

	// if cachedGroupStubs != nil {
	// 	return cachedGroupStubs, nil
	// }
	groupStubs, err := me.listUnguarded(ctx)
	if err != nil {
		return nil, err
	}
	cachedGroupStubs = groupStubs
	return cachedGroupStubs, nil
}

func (me *GroupServiceClient) listUnguarded(ctx context.Context) ([]*ListGroup, error) {
	var err error

	client := iam.NewIAMClient(me)
	var response ListGroupsResponse
	accountID := strings.TrimPrefix(me.AccountID(), "urn:dtaccount:")
	if err = iam.GET(client, ctx, fmt.Sprintf("%s/iam/v1/accounts/%s/groups", me.endpointURL, accountID), 200, false, &response); err != nil {
		return nil, err
	}
	return response.Items, nil
}

func (me *GroupServiceClient) Get(ctx context.Context, id string, v *groups.Group) (err error) {
	stubs, err := me.list(ctx)
	if err != nil {
		return err
	}
	for _, listStub := range stubs {
		if listStub.UUID == id {
			accountID := strings.TrimPrefix(me.AccountID(), "urn:dtaccount:")
			client := iam.NewIAMClient(me)
			var groupStub ListGroup
			if err = iam.GET(client, ctx, fmt.Sprintf("%s/iam/v1/accounts/%s/groups/%s/permissions", me.endpointURL, accountID, id), 200, false, &groupStub); err != nil {
				return err
			}

			v.Name = listStub.Name
			v.Description = listStub.Description
			v.FederatedAttributeValues = listStub.FederatedAttributeValues
			// ddd, _ := json.MarshalIndent(groupStub.Permissions, "", "  ")
			// logging.File.Println(string(ddd))
			v.Permissions = groupStub.Permissions
			return nil
		}
	}
	return rest.Error{Code: 404, Message: fmt.Sprintf("no group with id `%s` found", id)}
}

func (me *GroupServiceClient) Delete(ctx context.Context, id string) error {
	_, err := iam.NewIAMClient(me).DELETE(ctx, fmt.Sprintf("%s/iam/v1/accounts/%s/groups/%s", me.endpointURL, strings.TrimPrefix(me.AccountID(), "urn:dtaccount:"), id), 200, false)

	// data sources MAY have cached a list of group IDs
	// Updating the (publicly available) revision signals to them that either a CREATE or DELETE has happened since
	revisionLock.Lock()
	defer revisionLock.Unlock()
	revision = uuid.NewString()

	return err
}
