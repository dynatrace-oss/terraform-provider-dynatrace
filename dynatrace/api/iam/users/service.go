package users

import (
	"encoding/json"
	"fmt"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/iam"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/iam/groups"
	users "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/iam/users/settings"
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
func (me *UserServiceClient) Create(user *users.User) (*settings.Stub, error) {
	var err error

	client := iam.NewIAMClient(me)
	if _, err = client.POST(fmt.Sprintf("https://api.dynatrace.com/iam/v1/accounts/%s/users", me.AccountID()), user, 201, false); err != nil {
		return nil, err
	}

	groups := []string{}
	if len(user.Groups) > 0 {
		groups = user.Groups
	}
	if _, err = client.PUT(fmt.Sprintf("https://api.dynatrace.com/iam/v1/accounts/%s/users/%s/groups", me.AccountID(), user.Email), groups, 200, false); err != nil {
		return nil, err
	}

	return &settings.Stub{ID: user.Email, Name: user.Email}, nil
}

type GroupStub struct {
	GroupName string `json:"groupName"`
	UUID      string `json:"uuid"`
}

type GetUserGroupsResponse struct {
	Groups []*GroupStub
}

func (me *UserServiceClient) Get(email string, v *users.User) error {
	var err error
	var responseBytes []byte

	client := iam.NewIAMClient(me)

	if responseBytes, err = client.GET(fmt.Sprintf("https://api.dynatrace.com/iam/v1/accounts/%s/users/%s", me.AccountID(), email), 200, false); err != nil {
		return err
	}

	var response GetUserGroupsResponse
	if err = json.Unmarshal(responseBytes, &response); err != nil {
		return err
	}
	v.Email = email
	v.Groups = []string{}
	groupService := groups.NewGroupService(me.clientID, me.accountID, me.clientID)
	var visibleGroupIDs settings.Stubs
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

	client := iam.NewIAMClient(me)

	if _, err = client.PUT(fmt.Sprintf("https://api.dynatrace.com/iam/v1/accounts/%s/users/%s/groups", me.AccountID(), user.Email), user.Groups, 200, false); err != nil {
		return err
	}
	groups := []string{}
	if len(user.Groups) > 0 {
		groups = user.Groups
	}
	if _, err = client.PUT(fmt.Sprintf("https://api.dynatrace.com/iam/v1/accounts/%s/users/%s/groups", me.AccountID(), user.Email), groups, 200, false); err != nil {
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

func (me *UserServiceClient) List() (settings.Stubs, error) {
	var err error
	var responseBytes []byte

	if responseBytes, err = iam.NewIAMClient(me).GET(fmt.Sprintf("https://api.dynatrace.com/iam/v1/accounts/%s/users", me.AccountID()), 200, false); err != nil {
		return nil, err
	}

	var response ListUsersResponse
	if err = json.Unmarshal(responseBytes, &response); err != nil {
		return nil, err
	}
	var stubs settings.Stubs
	for _, item := range response.Items {
		stubs = append(stubs, &settings.Stub{ID: item.UID, Name: item.Email})
	}
	return stubs, nil
}

func (me *UserServiceClient) Delete(email string) error {
	_, err := iam.NewIAMClient(me).DELETE(fmt.Sprintf("https://api.dynatrace.com/iam/v1/accounts/%s/users/%s", me.AccountID(), email), 200, false)
	return err
}
