package mgmz

import (
	"context"
	"fmt"
	"strings"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api"
	mgmz "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/cluster/v1/permissions/mgmz/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings"
)

const SchemaID = "accounts:groups-mgmz"

func Service(credentials *settings.Credentials) settings.CRUDService[*mgmz.Permission] {
	return &service{
		serviceClient: NewService(fmt.Sprintf("%s%s", credentials.Cluster.URL, "/api/v1.0/onpremise"), credentials.Cluster.Token),
	}
}

// ServiceClient TODO: documentation
type ServiceClient struct {
	client rest.Client
}

func (me *service) Create(ctx context.Context, v *mgmz.Permission) (*api.Stub, error) {
	return me.serviceClient.Create(ctx, v)
}

func (me *service) Update(ctx context.Context, id string, v *mgmz.Permission) error {
	return me.serviceClient.Update(ctx, v)
}

func (me *service) Delete(ctx context.Context, id string) error {
	return me.serviceClient.Delete(ctx, id)
}

func (me *service) List(ctx context.Context) (api.Stubs, error) {
	return me.serviceClient.List(ctx)
}

func (me *service) Get(ctx context.Context, id string, v *mgmz.Permission) error {
	return me.serviceClient.Get(ctx, id, v)
}

func (me *service) SchemaID() string {
	return SchemaID
}

func (cs *ServiceClient) SchemaID() string {
	return SchemaID
}

// NewService creates a new Service Client
// baseURL should look like this: "https://#######.live.dynatrace.com/api/config/v1"
// token is an API Token
func NewService(baseURL string, token string) *ServiceClient {
	return &ServiceClient{client: rest.DefaultClient(baseURL, token)}
}

type service struct {
	serviceClient *ServiceClient
}

func (cs *ServiceClient) Create(ctx context.Context, permission *mgmz.Permission) (*api.Stub, error) {
	var err error
	permissionDTO := &mgmz.PermissionDTO{
		GroupID: permission.GroupID,
		EnvironmentPermissions: []*mgmz.EnvironmentPerms{
			{
				EnvironmentUUID: permission.EnvironmentID,
				ManagementPermissions: []*mgmz.ManagementZonePerms{
					{
						ManagementZoneID: permission.ManagementZoneID,
						Permissions:      permission.Permissions,
					},
				},
			},
		},
	}
	if err = cs.client.Put("/groups/managementZones", permissionDTO, 200).Finish(); err != nil {
		return nil, err
	}
	return &api.Stub{ID: toID(permission), Name: toID(permission)}, nil
}

func (cs *ServiceClient) Update(ctx context.Context, permission *mgmz.Permission) error {
	var err error
	permissionDTO := &mgmz.PermissionDTO{
		GroupID: permission.GroupID,
		EnvironmentPermissions: []*mgmz.EnvironmentPerms{
			{
				EnvironmentUUID: permission.EnvironmentID,
				ManagementPermissions: []*mgmz.ManagementZonePerms{
					{
						ManagementZoneID: permission.ManagementZoneID,
						Permissions:      permission.Permissions,
					},
				},
			},
		},
	}
	if err = cs.client.Put("/groups/managementZones", permissionDTO, 200).Finish(); err != nil {
		return err
	}
	return nil
}

func (cs *ServiceClient) Delete(ctx context.Context, id string) error {
	var permission mgmz.Permission
	if err := readID(id, &permission); err != nil {
		return err
	}
	var err error
	permissionDTO := &mgmz.PermissionDTO{
		GroupID: permission.GroupID,
		EnvironmentPermissions: []*mgmz.EnvironmentPerms{
			{
				EnvironmentUUID: permission.EnvironmentID,
				ManagementPermissions: []*mgmz.ManagementZonePerms{
					{
						ManagementZoneID: permission.ManagementZoneID,
						Permissions:      []string{},
					},
				},
			},
		},
	}
	if err = cs.client.Put("/groups/managementZones", permissionDTO, 200).Finish(); err != nil {
		// if either environment or management zone don't exist anymore, it safe to say the permission doesn't exist anymore
		if err.Error() == fmt.Sprintf("Management zone with ID %s doesn't exist within environment %s", permission.ManagementZoneID, permission.EnvironmentID) {
			return nil
		}
		if err.Error() == fmt.Sprintf("Environment with UUID %s doesn't exist", permission.EnvironmentID) {
			return nil
		}
		return err
	}
	return nil
}

func (cs *ServiceClient) Get(ctx context.Context, id string, permission *mgmz.Permission) error {
	var err error

	if err := readID(id, permission); err != nil {
		return err
	}

	permissionDTO := mgmz.PermissionDTO{}
	if err = cs.client.Get(fmt.Sprintf("/groups/managementZones/%s", permission.GroupID), 200).Finish(&permissionDTO); err != nil {
		return err
	}
	for _, envperm := range permissionDTO.EnvironmentPermissions {
		for _, mgmzperm := range envperm.ManagementPermissions {
			if mgmzperm.ManagementZoneID == permission.ManagementZoneID {
				permission.Permissions = mgmzperm.Permissions
				return nil
			}
		}
	}
	return rest.Error{Code: 404, Method: "GET", URL: "/groups/managementZones", Message: fmt.Sprintf("No permissions found for group %s, management zone %s, environment %s", permission.GroupID, permission.ManagementZoneID, permission.EnvironmentID), ConstraintViolations: []rest.ConstraintViolation{}}
}

type ListResponse []struct {
	GroupId                     string `json:"groupId"`
	MzPermissionsPerEnvironment []struct {
		EnvironmentUuid string `json:"environmentUuid"`
		MzPermissions   []struct {
			MzId        string   `json:"mzId"`
			Permissions []string `json:"permissions"`
		} `json:"mzPermissions"`
	} `json:"mzPermissionsPerEnvironment"`
}

func (cs *ServiceClient) List(ctx context.Context) (api.Stubs, error) {
	var err error
	var stubs api.Stubs

	response := ListResponse{}
	if err = cs.client.Get("/groups/managementZones", 200).Finish(&response); err != nil {
		return nil, err
	}

	for _, group := range response {
		for _, mzPermissionsPerEnv := range group.MzPermissionsPerEnvironment {
			for _, mzPermissions := range mzPermissionsPerEnv.MzPermissions {
				if len(mzPermissions.Permissions) > 0 {
					joinedId := fmt.Sprintf("%s#-#%s#-#%s", mzPermissionsPerEnv.EnvironmentUuid, group.GroupId, mzPermissions.MzId)
					stubs = append(stubs, &api.Stub{ID: joinedId, Name: joinedId})
				}
			}
		}
	}

	return stubs, nil
}

func toID(permission *mgmz.Permission) string {
	return fmt.Sprintf("%s#-#%s#-#%s", permission.EnvironmentID, permission.GroupID, permission.ManagementZoneID)
}

func readID(id string, permission *mgmz.Permission) error {
	parts := strings.Split(id, "#-#")
	if len(parts) != 3 {
		return fmt.Errorf("%s is not a valid permission ID", id)
	}
	permission.EnvironmentID = parts[0]
	permission.GroupID = parts[1]
	permission.ManagementZoneID = parts[2]
	return nil
}
