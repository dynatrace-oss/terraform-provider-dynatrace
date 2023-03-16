package mgmz

import (
	"fmt"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest"
)

// ServiceClient TODO: documentation
type ServiceClient struct {
	client rest.Client
}

// NewService creates a new Service Client
// baseURL should look like this: "https://#######.live.dynatrace.com/api/config/v1"
// token is an API Token
func NewService(baseURL string, token string) *ServiceClient {
	return &ServiceClient{client: rest.DefaultClient(baseURL, token)}
}

func (cs *ServiceClient) Create(permission *Permission) (string, error) {
	var err error
	permissionDTO := &PermissionDTO{
		GroupID: permission.GroupID,
		EnvironmentPermissions: []*EnvironmentPerms{
			{
				EnvironmentUUID: permission.EnvironmentID,
				ManagementPermissions: []*ManagementZonePerms{
					{
						ManagementZoneID: permission.ManagementZoneID,
						Permissions:      permission.Permissions,
					},
				},
			},
		},
	}
	if err = cs.client.Put("/groups/managementZones", permissionDTO, 200).Finish(); err != nil {
		return "", err
	}
	return permission.toID(), nil
}

func (cs *ServiceClient) Update(permission *Permission) error {
	var err error
	permissionDTO := &PermissionDTO{
		GroupID: permission.GroupID,
		EnvironmentPermissions: []*EnvironmentPerms{
			{
				EnvironmentUUID: permission.EnvironmentID,
				ManagementPermissions: []*ManagementZonePerms{
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

func (cs *ServiceClient) Delete(id string) error {
	var permission Permission
	if err := permission.readID(id); err != nil {
		return err
	}
	var err error
	permissionDTO := &PermissionDTO{
		GroupID: permission.GroupID,
		EnvironmentPermissions: []*EnvironmentPerms{
			{
				EnvironmentUUID: permission.EnvironmentID,
				ManagementPermissions: []*ManagementZonePerms{
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

func (cs *ServiceClient) Get(id string) (*Permission, error) {
	var err error

	permission := &Permission{}
	if err := permission.readID(id); err != nil {
		return nil, err
	}

	permissionDTO := PermissionDTO{}
	if err = cs.client.Get(fmt.Sprintf("/groups/managementZones/%s", permission.GroupID), 200).Finish(&permissionDTO); err != nil {
		return nil, err
	}
	for _, envperm := range permissionDTO.EnvironmentPermissions {
		for _, mgmzperm := range envperm.ManagementPermissions {
			if mgmzperm.ManagementZoneID == permission.ManagementZoneID {
				permission.Permissions = mgmzperm.Permissions
				return permission, nil
			}
		}
	}
	return nil, rest.Error{Code: 404, Method: "GET", URL: "/groups/managementZones", Message: fmt.Sprintf("No permissions found for group %s, management zone %s, environment %s", permission.GroupID, permission.ManagementZoneID, permission.EnvironmentID), ConstraintViolations: []rest.ConstraintViolation{}}
}
