//go:build unit

package settings

import (
	"testing"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
)

func TestSettingPermissions_MarshalHCL(t *testing.T) {
	cases := []struct {
		name        string
		permissions *SettingPermissions
		expected    hcl.Properties
	}{
		{
			name: "full-permissions",
			permissions: &SettingPermissions{
				SettingsObjectID: "obj-123",
				AllUsers:         HCLAccessorNone,
				Users: Users{
					{UID: "user1", Access: HCLAccessorWrite},
				},
				Groups: Groups{
					{ID: "group1", Access: HCLAccessorRead},
					{ID: "group2", Access: HCLAccessorWrite},
				},
			},
			expected: hcl.Properties{
				"settings_object_id": "obj-123",
				"all_users":          HCLAccessorNone,
				"users": []interface{}{
					hcl.Properties{
						"user": []interface{}{
							hcl.Properties{
								"access": HCLAccessorWrite,
								"uid":    "user1",
							},
						},
					},
				},
				"groups": []interface{}{
					hcl.Properties{
						"group": []interface{}{
							hcl.Properties{
								"access": HCLAccessorRead,
								"id":     "group1",
							},
							hcl.Properties{
								"access": HCLAccessorWrite,
								"id":     "group2",
							},
						},
					},
				},
			},
		},
		{
			name: "minimal-permissions",
			permissions: &SettingPermissions{
				SettingsObjectID: "obj-123",
				AllUsers:         HCLAccessorNone,
			},
			expected: hcl.Properties{
				"settings_object_id": "obj-123",
				"all_users":          HCLAccessorNone,
				"users":              nil,
				"groups":             nil,
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var actual = hcl.Properties{}
			err := c.permissions.MarshalHCL(actual)

			assert.NoError(t, err)
			assert.Equal(t, c.expected, actual)
		})
	}
}

func TestSettingPermissions_UnmarshalHCL(t *testing.T) {
	s := new(SettingPermissions).Schema()

	cases := []struct {
		name     string
		input    map[string]interface{}
		expected SettingPermissions
	}{
		{
			name: "full-permissions",
			input: map[string]interface{}{
				"settings_object_id": "obj-123",
				"all_users":          HCLAccessorNone,
				"users": []interface{}{
					map[string]interface{}{
						"user": []interface{}{
							map[string]interface{}{
								"uid":    "user1",
								"access": HCLAccessorWrite,
							},
							map[string]interface{}{
								"uid":    "user2",
								"access": HCLAccessorWrite,
							},
						},
					},
				},
				"groups": []interface{}{
					map[string]interface{}{
						"group": []interface{}{
							map[string]interface{}{
								"id":     "group1",
								"access": HCLAccessorRead,
							},
							map[string]interface{}{
								"id":     "group2",
								"access": HCLAccessorWrite,
							},
						},
					},
				},
			},
			expected: SettingPermissions{
				SettingsObjectID: "obj-123",
				AllUsers:         HCLAccessorNone,
				Users: Users{
					{UID: "user1", Access: HCLAccessorWrite},
					{UID: "user2", Access: HCLAccessorWrite},
				},
				Groups: Groups{
					{ID: "group1", Access: HCLAccessorRead},
					{ID: "group2", Access: HCLAccessorWrite},
				},
			},
		},
		{
			name: "minimal-permissions",
			input: map[string]interface{}{
				"settings_object_id": "obj-123",
				"all_users":          HCLAccessorNone,
			},
			expected: SettingPermissions{
				SettingsObjectID: "obj-123",
				AllUsers:         HCLAccessorNone,
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			d := schema.TestResourceDataRaw(t, s, c.input)
			assert.NotNil(t, d)

			var actual SettingPermissions
			decoder := hcl.DecoderFrom(d)
			err := actual.UnmarshalHCL(decoder)

			assert.NoError(t, err)
			assertSettingPermissionsAreEqual(t, c.expected, actual)
		})
	}
}

func assertSettingPermissionsAreEqual(t *testing.T, this SettingPermissions, that SettingPermissions) {
	assert.Equal(t, this.SettingsObjectID, that.SettingsObjectID)
	assert.Equal(t, this.AllUsers, that.AllUsers)
	assert.ElementsMatch(t, this.Users, that.Users)
	assert.ElementsMatch(t, this.Groups, that.Groups)
}
