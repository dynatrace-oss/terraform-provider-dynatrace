package settings

type TypePermissions = string

const (
	Read  TypePermissions = "r"
	Write TypePermissions = "w"
)

type TypeAccessor = string

const (
	AllUsers TypeAccessor = "all-users"
	Group    TypeAccessor = "group"
	User     TypeAccessor = "user"
)

type Accessor struct {
	Type TypeAccessor `json:"type"`
	// ID is only needed and valid for TypeAccessor Group and User, but not for AllUsers
	ID string `json:"id,omitempty"`
}

// PermissionObject represents the permissions for a settings object.
// It represents the request payload for POST, and the response payload of a single item for GET requests to the permissions endpoint.
type PermissionObject struct {
	Accessor Accessor `json:"accessor"`
	// Explanation:
	// 	- If Permissions is empty => invalid request
	// 	- If Permissions only contains "w" => invalid request
	// 	- If Permissions contains "r", the Accessor has read access
	// 	- If Permissions contains "w" and "r", the Accessor has write access
	Permissions []TypePermissions `json:"permissions"`
}

// PermissionObjects represents the response payload of a GET request to the permissions endpoint.
type PermissionObjects struct {
	Accessors []PermissionObject `json:"accessors"`
}

// PermissionObjectUpdate represents the updated permissions for a settings object to an accessor.
type PermissionObjectUpdate struct {
	Permissions []TypePermissions `json:"permissions"`
}
