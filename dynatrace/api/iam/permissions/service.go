package permissions

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/iam"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/iam/groups"
	permissions "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/iam/permissions/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings"
)

type PermissionServiceClient struct {
	clientID     string
	accountID    string
	clientSecret string
	tokenURL     string
	endpointURL  string
}

func (me *PermissionServiceClient) ClientID() string {
	return me.clientID
}

func (me *PermissionServiceClient) AccountID() string {
	return me.accountID
}

func (me *PermissionServiceClient) ClientSecret() string {
	return me.clientSecret
}

func (me *PermissionServiceClient) TokenURL() string {
	return me.tokenURL
}

func (me *PermissionServiceClient) EndpointURL() string {
	return me.endpointURL
}

func Service(credentials *settings.Credentials) settings.CRUDService[*permissions.Permission] {
	return &PermissionServiceClient{clientID: credentials.IAM.ClientID, accountID: credentials.IAM.AccountID, clientSecret: credentials.IAM.ClientSecret, tokenURL: credentials.IAM.TokenURL, endpointURL: credentials.IAM.EndpointURL}
}

func (me *PermissionServiceClient) SchemaID() string {
	return "accounts:iam:permissions"
}

func (me *PermissionServiceClient) Name() string {
	return me.SchemaID()
}

func (me *PermissionServiceClient) Create(ctx context.Context, permission *permissions.Permission) (*api.Stub, error) {
	var err error

	client := iam.NewIAMClient(me)
	scope := ""
	scopeType := ""
	if len(permission.Account) > 0 {
		scope = permission.Account
		scopeType = "account"
	} else if len(permission.ManagementZone) > 0 {
		scope = permission.Environment + ":" + permission.ManagementZone
		scopeType = "management-zone"
	} else if len(permission.Environment) > 0 {
		scope = permission.Environment
		scopeType = "tenant"
	}
	payload := []permissions.PermissionDTO{{
		GroupID:   permission.GroupID,
		Scope:     scope,
		ScopeType: scopeType,
		Name:      permission.Name,
	}}
	if _, err = client.POST(ctx, fmt.Sprintf("%s/iam/v1/accounts/%s/groups/%s/permissions", me.endpointURL, strings.TrimPrefix(me.AccountID(), "urn:dtaccount:"), permission.GroupID), payload, 201, false); err != nil {
		return nil, err
	}

	return &api.Stub{ID: payload[0].ToID(permission.GroupID), Name: permission.Name}, nil
}

type GetGroupPermissionsResponse struct {
	Permissions []*permissions.PermissionDTO
}

func (me *PermissionServiceClient) Get(ctx context.Context, id string, v *permissions.Permission) error {
	var err error
	var responseBytes []byte

	client := iam.NewIAMClient(me)

	parts := strings.Split(id, "#-#")
	if len(parts) < 4 {
		return fmt.Errorf("'%s' is not a valid permission ID", id)
	}
	groupID := parts[0]
	name := parts[1]
	scope := parts[2]
	scopeType := parts[3]

	if responseBytes, err = client.GET(ctx, fmt.Sprintf("%s/iam/v1/accounts/%s/groups/%s/permissions", me.endpointURL, strings.TrimPrefix(me.AccountID(), "urn:dtaccount:"), groupID), 200, false); err != nil {
		return err
	}

	var response GetGroupPermissionsResponse
	if err = json.Unmarshal(responseBytes, &response); err != nil {
		return err
	}
	if len(response.Permissions) > 0 {
		for _, permission := range response.Permissions {
			permissionID := permission.ToID(groupID)
			if permissionID == id {
				v.GroupID = groupID
				v.Name = permission.Name
				if permission.ScopeType == "management-zone" {
					v.Environment = strings.Split(permission.Scope, ":")[0]
					v.ManagementZone = strings.Split(permission.Scope, ":")[1]
				} else if permission.ScopeType == "tenant" {
					v.Environment = permission.Scope
				} else if permission.ScopeType == "account" {
					v.Account = permission.Scope
				}
				return nil
			}
		}
	}

	return rest.Error{Code: 404, Message: fmt.Sprintf("there exists no permission for group %s with name %s, scope %s and scope type %s", groupID, name, scope, scopeType)}
}

func (me *PermissionServiceClient) Update(ctx context.Context, email string, permission *permissions.Permission) error {
	return errors.New("permissions are not expected to get updated - only destroy and create are possible")
}

func (me *PermissionServiceClient) List(ctx context.Context) (api.Stubs, error) {
	groupsService := groups.NewGroupService(me.clientID, me.accountID, me.clientSecret, me.tokenURL, me.endpointURL)
	groupStubs, err := groupsService.List(ctx)
	if err != nil {
		return nil, err
	}

	var stubs api.Stubs

	client := iam.NewIAMClient(me)
	for _, groupStub := range groupStubs {
		groupID := groupStub.ID

		accountID := strings.TrimPrefix(me.AccountID(), "urn:dtaccount:")

		var response GetGroupPermissionsResponse
		if err = iam.GET(client, ctx, fmt.Sprintf("%s/iam/v1/accounts/%s/groups/%s/permissions", me.endpointURL, accountID, groupID), 200, false, &response); err != nil {
			return nil, err
		}

		if len(response.Permissions) > 0 {
			for _, permission := range response.Permissions {
				permissionID := strings.Join([]string{groupID, permission.Name, permission.Scope, permission.ScopeType}, "#-#")
				stubs = append(stubs, &api.Stub{ID: permissionID, Name: permissionID})
			}
		}
	}

	return stubs, nil
}

func (me *PermissionServiceClient) Delete(ctx context.Context, id string) error {
	parts := strings.Split(id, "#-#")
	if len(parts) < 4 {
		return fmt.Errorf("'%s' is not a valid permission ID", id)
	}
	groupID := parts[0]
	name := parts[1]
	scope := parts[2]
	scopeType := parts[3]

	_, err := iam.NewIAMClient(me).DELETE(ctx, fmt.Sprintf("%s/iam/v1/accounts/%s/groups/%s/permissions?scope=%s&permission-name=%s&scope-type=%s", me.endpointURL, strings.TrimPrefix(me.AccountID(), "urn:dtaccount:"), groupID, url.QueryEscape(scope), url.QueryEscape(name), url.QueryEscape(scopeType)), 200, false)
	if err != nil && strings.Contains(err.Error(), fmt.Sprintf("Permission %s not found", id)) {
		return nil
	}
	return err
}
