package users

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/iam"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/iam/groups"
	users "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/iam/users/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings"
)

type UserServiceClient struct {
	clientID     string
	accountID    string
	clientSecret string
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

func NewUserService(clientID string, accountID string, clientSecret string) *UserServiceClient {
	return &UserServiceClient{clientID: clientID, accountID: accountID, clientSecret: clientSecret}
}

func Service(credentials *settings.Credentials) settings.CRUDService[*users.User] {
	return &UserServiceClient{clientID: credentials.IAM.ClientID, accountID: credentials.IAM.AccountID, clientSecret: credentials.IAM.ClientSecret}
}

func (me *UserServiceClient) SchemaID() string {
	return "accounts:iam:users"
}

func (me *UserServiceClient) Name() string {
	return me.SchemaID()
}

func (me *UserServiceClient) Create(user *users.User) (*api.Stub, error) {
	var err error

	client := iam.NewIAMClient(me)
	if _, err = client.POST(fmt.Sprintf("https://api.dynatrace.com/iam/v1/accounts/%s/users", strings.TrimPrefix(me.AccountID(), "urn:dtaccount:")), user, 201, false); err != nil {
		if err.Error() == "User already exists" {
			if err = me.Update(user.Email, user); err != nil {
				return nil, err
			}
			return &api.Stub{ID: user.Email, Name: user.Email}, nil
		}
		return nil, err
	}

	groups := []string{}
	if len(user.Groups) > 0 {
		groups = user.Groups
	}
	if _, err = client.PUT(fmt.Sprintf("https://api.dynatrace.com/iam/v1/accounts/%s/users/%s/groups", strings.TrimPrefix(me.AccountID(), "urn:dtaccount:"), user.Email), groups, 200, false); err != nil {
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

func (me *UserServiceClient) Get(email string, v *users.User) error {
	var err error
	var responseBytes []byte

	client := iam.NewIAMClient(me)

	if responseBytes, err = client.GET(fmt.Sprintf("https://api.dynatrace.com/iam/v1/accounts/%s/users/%s", strings.TrimPrefix(me.AccountID(), "urn:dtaccount:"), email), 200, false); err != nil {
		if err != nil && strings.Contains(err.Error(), fmt.Sprintf("User %s not found", email)) {
			return rest.Error{Code: 404, Message: err.Error()}
		}
		return err
	}

	var response GetUserGroupsResponse
	if err = json.Unmarshal(responseBytes, &response); err != nil {
		return err
	}
	v.Email = email
	v.Groups = []string{}
	v.UID = response.UID
	groupService := groups.NewGroupService(me.clientID, me.accountID, me.clientID)
	var visibleGroupIDs api.Stubs
	if visibleGroupIDs, err = groupService.List(); err != nil {
		return err
	}
	for _, group := range response.Groups {
		for _, stub := range visibleGroupIDs {
			if stub.ID == group.UUID {
				v.Groups = append(v.Groups, group.UUID)
				break
			}
		}
	}

	return nil
}

func (me *UserServiceClient) Update(email string, user *users.User) error {
	var err error

	groups := []string{}
	if len(user.Groups) > 0 {
		groups = user.Groups
	}
	if _, err = iam.NewIAMClient(me).PUT(fmt.Sprintf("https://api.dynatrace.com/iam/v1/accounts/%s/users/%s/groups", strings.TrimPrefix(me.AccountID(), "urn:dtaccount:"), user.Email), groups, 200, false); err != nil {
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

func (me *UserServiceClient) List() (api.Stubs, error) {
	var err error
	var responseBytes []byte

	if responseBytes, err = iam.NewIAMClient(me).GET(fmt.Sprintf("https://api.dynatrace.com/iam/v1/accounts/%s/users", strings.TrimPrefix(me.AccountID(), "urn:dtaccount:")), 200, false); err != nil {
		return nil, err
	}

	var response ListUsersResponse
	if err = json.Unmarshal(responseBytes, &response); err != nil {
		return nil, err
	}
	var stubs api.Stubs
	for _, item := range response.Items {
		stubs = append(stubs, &api.Stub{ID: item.Email, Name: item.Email})
	}
	return stubs, nil
}

func (me *UserServiceClient) Delete(email string) error {
	_, err := iam.NewIAMClient(me).DELETE(fmt.Sprintf("https://api.dynatrace.com/iam/v1/accounts/%s/users/%s", strings.TrimPrefix(me.AccountID(), "urn:dtaccount:"), email), 200, false)
	if err != nil && strings.Contains(err.Error(), fmt.Sprintf("User %s not found", email)) {
		return nil
	}
	return err
}
