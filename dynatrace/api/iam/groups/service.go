package groups

import (
	"encoding/json"
	"fmt"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/iam"
	groups "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/iam/groups/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings"
)

type GroupServiceClient struct {
	clientID     string
	accountID    string
	clientSecret string
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

func NewGroupService(clientID string, accountID string, clientSecret string) settings.CRUDService[*groups.Group] {
	return &GroupServiceClient{clientID: clientID, accountID: accountID, clientSecret: clientSecret}
}

func Service(credentials *settings.Credentials) settings.CRUDService[*groups.Group] {
	return &GroupServiceClient{clientID: credentials.IAM.ClientID, accountID: credentials.IAM.AccountID, clientSecret: credentials.IAM.ClientSecret}
}

func (me *GroupServiceClient) SchemaID() string {
	return "accounts:iam:groups"
}

func (me *GroupServiceClient) Create(group *groups.Group) (*settings.Stub, error) {
	var err error
	var responseBytes []byte

	client := iam.NewIAMClient(me)
	if responseBytes, err = client.POST(fmt.Sprintf("https://api.dynatrace.com/iam/v1/accounts/%s/groups", me.AccountID()), []*groups.Group{group}, 201, false); err != nil {
		return nil, err
	}

	responseGroups := []*ListGroup{}
	if err = json.Unmarshal(responseBytes, &responseGroups); err != nil {
		return nil, err
	}
	groupID := responseGroups[0].UUID
	groupName := responseGroups[0].Name

	if len(group.Permissions) > 0 {
		if _, err = client.PUT(fmt.Sprintf("https://api.dynatrace.com/iam/v1/accounts/%s/groups/%s/permissions", me.AccountID(), groupID), group.Permissions, 200, false); err != nil {
			return nil, err
		}
	}

	return &settings.Stub{ID: groupID, Name: groupName}, nil
}

func (me *GroupServiceClient) Update(uuid string, group *groups.Group) error {
	var err error

	client := iam.NewIAMClient(me)
	if _, err = client.POST(fmt.Sprintf("https://api.dynatrace.com/iam/v1/accounts/%s/groups/%s", me.AccountID(), uuid), []*groups.Group{group}, 201, false); err != nil {
		return err
	}

	permissions := []*groups.Permission{}

	if len(group.Permissions) > 0 {
		permissions = group.Permissions
	}
	if _, err = client.PUT(fmt.Sprintf("https://api.dynatrace.com/iam/v1/accounts/%s/groups/%s/permissions", me.AccountID(), uuid), permissions, 200, false); err != nil {
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

func (me *GroupServiceClient) List() (settings.Stubs, error) {
	var err error
	var responseBytes []byte

	if responseBytes, err = iam.NewIAMClient(me).GET(fmt.Sprintf("https://api.dynatrace.com/iam/v1/accounts/%s/groups", me.AccountID()), 200, false); err != nil {
		return nil, err
	}

	var response ListGroupsResponse
	if err = json.Unmarshal(responseBytes, &response); err != nil {
		return nil, err
	}
	var stubs settings.Stubs
	for _, elem := range response.Items {
		stubs = append(stubs, &settings.Stub{ID: elem.UUID, Name: elem.Name})
	}
	return stubs, nil
}

func (me *GroupServiceClient) Get(id string, v *groups.Group) (err error) {
	var stubs settings.Stubs

	if stubs, err = me.List(); err != nil {
		return err
	}

	for _, stub := range stubs {
		if stub.ID == id {
			var responseBytes []byte

			if responseBytes, err = iam.NewIAMClient(me).GET(fmt.Sprintf("https://api.dynatrace.com/iam/v1/accounts/%s/groups/%s/permissions", me.AccountID(), id), 200, false); err != nil {
				return err
			}
			var groupStub ListGroup
			if err = json.Unmarshal(responseBytes, &groupStub); err != nil {
				return err
			}
			v.Name = groupStub.Name
			v.Description = groupStub.Description
			v.FederatedAttributeValues = groupStub.FederatedAttributeValues
			v.Permissions = groupStub.Permissions
			return nil
		}
	}

	return rest.Error{Code: 404, Message: fmt.Sprintf("a group with id %s doesn't exist", id)}
}

func (me *GroupServiceClient) Delete(id string) error {
	_, err := iam.NewIAMClient(me).DELETE(fmt.Sprintf("https://api.dynatrace.com/iam/v1/accounts/%s/groups/%s", me.AccountID(), id), 200, false)
	return err
}
