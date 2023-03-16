package mgmz

type PermissionDTO struct {
	GroupID                string              `json:"groupId"`
	EnvironmentPermissions []*EnvironmentPerms `json:"mzPermissionsPerEnvironment"`
}

type EnvironmentPerms struct {
	EnvironmentUUID       string                 `json:"environmentUuid"`
	ManagementPermissions []*ManagementZonePerms `json:"mzPermissions"`
}

type ManagementZonePerms struct {
	ManagementZoneID string   `json:"mzId"`
	Permissions      []string `json:"permissions"`
}
