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

func Service(credentials *settings.Credentials) settings.CRUDService[*users.User] {

	return &UserServiceClient{clientID: credentials.IAM.ClientID, accountID: credentials.IAM.AccountID, clientSecret: credentials.IAM.ClientSecret, tokenURL: credentials.IAM.TokenURL, endpointURL: credentials.IAM.EndpointURL}
}

func (me *UserServiceClient) SchemaID() string {
	return "accounts:iam:users"
}

func (me *UserServiceClient) Create(ctx context.Context, user *users.User) (*api.Stub, error) {
	var err error

	client := iam.NewIAMClient(me)
	if _, err = client.POST(ctx, fmt.Sprintf("%s/iam/v1/accounts/%s/users", me.endpointURL, strings.TrimPrefix(me.AccountID(), "urn:dtaccount:")), user, 201, false); err != nil {
		if err.Error() == "User already exists" {
			if err = me.Update(ctx, user.Email, user); err != nil {
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
	if _, err = client.PUT(ctx, fmt.Sprintf("%s/iam/v1/accounts/%s/users/%s/groups", me.endpointURL, strings.TrimPrefix(me.AccountID(), "urn:dtaccount:"), user.Email), groups, 200, false); err != nil {
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
	var err error
	var responseBytes []byte

	client := iam.NewIAMClient(me)

	if responseBytes, err = client.GET(ctx, fmt.Sprintf("%s/iam/v1/accounts/%s/users/%s", me.endpointURL, strings.TrimPrefix(me.AccountID(), "urn:dtaccount:"), email), 200, false); err != nil {
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
	for _, group := range response.Groups {
		v.Groups = append(v.Groups, group.UUID)
	}

	return nil
}

func (me *UserServiceClient) Update(ctx context.Context, email string, user *users.User) error {
	var err error

	groups := []string{}
	if len(user.Groups) > 0 {
		groups = user.Groups
	}
	if _, err = iam.NewIAMClient(me).PUT(ctx, fmt.Sprintf("%s/iam/v1/accounts/%s/users/%s/groups", me.endpointURL, strings.TrimPrefix(me.AccountID(), "urn:dtaccount:"), user.Email), groups, 200, false); err != nil {
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
	var err error
	var responseBytes []byte

	if responseBytes, err = iam.NewIAMClient(me).GET(ctx, fmt.Sprintf("%s/iam/v1/accounts/%s/users", me.endpointURL, strings.TrimPrefix(me.AccountID(), "urn:dtaccount:")), 200, false); err != nil {
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

func (me *UserServiceClient) Delete(ctx context.Context, email string) error {
	_, err := iam.NewIAMClient(me).DELETE(ctx, fmt.Sprintf("%s/iam/v1/accounts/%s/users/%s", me.endpointURL, strings.TrimPrefix(me.AccountID(), "urn:dtaccount:"), email), 200, false)
	if err != nil && strings.Contains(err.Error(), fmt.Sprintf("User %s not found", email)) {
		return nil
	}
	return err
}
