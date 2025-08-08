package permissions

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v2/settings/objects/permissions/convert"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v2/settings/objects/permissions/reconcile"
	permissions "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v2/settings/objects/permissions/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings"

	permissions2 "github.com/dynatrace/dynatrace-configuration-as-code-core/clients/settings/permissions"
)

func Service(credentials *rest.Credentials) settings.CRUDService[*permissions.SettingPermissions] {
	return &ServiceImpl{Credentials: credentials}
}

type ServiceImpl struct {
	Credentials    *rest.Credentials
	Client         permissions.PermissionClient
	SettingsClient SettingsClient
}

func (me *ServiceImpl) getClient(ctx context.Context) (permissions.PermissionClient, error) {
	if me.Client != nil {
		return me.Client, nil
	}
	restClient, err := rest.CreatePlatformClient(ctx, me.Credentials.OAuth.EnvironmentURL, me.Credentials)
	if err != nil {
		return nil, err
	}

	me.Client = permissions2.NewClient(restClient)
	return me.Client, nil
}

func (me *ServiceImpl) getSettingsClient(ctx context.Context) (SettingsClient, error) {
	if me.SettingsClient != nil {
		return me.SettingsClient, nil
	}
	restClient, err := rest.CreatePlatformClient(ctx, me.Credentials.OAuth.EnvironmentURL, me.Credentials)
	if err != nil {
		return nil, err
	}

	me.SettingsClient = NewSettingsClient(restClient)
	return me.SettingsClient, nil
}

func (me *ServiceImpl) Get(ctx context.Context, objectID string, v *permissions.SettingPermissions) error {
	currentPermissions, err := me.get(ctx, objectID)
	if err != nil {
		return err
	}

	v.SettingsObjectID = objectID
	v.Users = currentPermissions.Users
	v.Groups = currentPermissions.Groups
	v.AllUsers = currentPermissions.AllUsers
	return nil
}

func (me *ServiceImpl) get(ctx context.Context, objectID string) (*permissions.SettingPermissions, error) {
	client, err := me.getClient(ctx)
	if err != nil {
		return nil, err
	}
	req, err := client.GetAllAccessors(ctx, objectID, false)
	if err != nil {
		return nil, err
	}
	var permissionObjects permissions.PermissionObjects
	err = json.Unmarshal(req.Data, &permissionObjects)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal permissions response: %w", err)
	}
	return convert.DTOToHCL(permissionObjects, objectID), nil
}

func (me *ServiceImpl) SchemaID() string {
	return "settings:permissions"
}

func (me *ServiceImpl) List(ctx context.Context) (api.Stubs, error) {
	client, err := me.getSettingsClient(ctx)
	if err != nil {
		return nil, err
	}
	schemaIds, err := client.GetSchemaIDsWithOwnerBasedAccessControl(ctx)
	if err != nil {
		return nil, err
	}

	stubs := api.Stubs{}
	for _, schemaID := range schemaIds {
		objectIDs, err := client.ListObjectsIDsOfSchema(ctx, schemaID)
		if err != nil {
			return nil, err
		}

		for _, id := range objectIDs {
			stubs = append(stubs, &api.Stub{
				ID:   id,
				Name: id,
			})
		}
	}
	return stubs, nil
}

func (me *ServiceImpl) Upsert(ctx context.Context, v *permissions.SettingPermissions) (*api.Stub, error) {
	client, err := me.getClient(ctx)
	if err != nil {
		return nil, err
	}
	data, err := me.get(ctx, v.SettingsObjectID)
	if err != nil {
		return nil, err
	}

	err = reconcile.CompareAndUpdate(ctx, client, data, v)
	if err != nil {
		return nil, err
	}

	return &api.Stub{ID: v.SettingsObjectID}, nil
}

func (me *ServiceImpl) Create(ctx context.Context, v *permissions.SettingPermissions) (*api.Stub, error) {
	return me.Upsert(ctx, v)
}

func (me *ServiceImpl) Update(ctx context.Context, _ string, v *permissions.SettingPermissions) error {
	_, err := me.Upsert(ctx, v)
	return err
}

func (me *ServiceImpl) Delete(ctx context.Context, id string) error {
	// via the Upsert we can delete the permissions by passing an empty permissions object
	emptyPermissions := &permissions.SettingPermissions{
		SettingsObjectID: id,
		Users:            nil,
		Groups:           nil,
		AllUsers:         permissions.HCLAccessorNone,
	}
	_, err := me.Upsert(ctx, emptyPermissions)
	return err
}
