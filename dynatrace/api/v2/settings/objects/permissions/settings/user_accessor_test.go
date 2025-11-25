/**
* @license
* Copyright 2025 Dynatrace LLC
*
* Licensed under the Apache License, Version 2.0 (the "License");
* you may not use this file except in compliance with the License.
* You may obtain a copy of the License at
*
*     http://www.apache.org/licenses/LICENSE-2.0
*
* Unless required by applicable law or agreed to in writing, software
* distributed under the License is distributed on an "AS IS" BASIS,
* WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
* See the License for the specific language governing permissions and
* limitations under the License.
 */

//go:build unit

package settings

import (
	"testing"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
)

func TestUserAccessor_MarshalHCL(t *testing.T) {
	cases := []struct {
		name     string
		accessor *UserAccessor
		expected hcl.Properties
	}{
		{
			name: "HCLAccessorRead-access",
			accessor: &UserAccessor{
				"my-id",
				HCLAccessorRead,
			},
			expected: map[string]any{"uid": "my-id", "access": HCLAccessorRead},
		},
		{
			name: "write-access",
			accessor: &UserAccessor{
				"my-id",
				HCLAccessorWrite,
			},
			expected: map[string]any{"uid": "my-id", "access": HCLAccessorWrite},
		},
		{
			name:     "empty",
			accessor: &UserAccessor{},
			expected: map[string]any{"uid": "", "access": ""},
		},
		{
			name: "empty UID",
			accessor: &UserAccessor{
				"",
				HCLAccessorRead,
			},
			expected: map[string]any{"uid": "", "access": HCLAccessorRead},
		},
		{
			name: "empty access",
			accessor: &UserAccessor{
				"my-id",
				"",
			},
			expected: map[string]any{"uid": "my-id", "access": ""},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {

			var actual = hcl.Properties{}
			_ = c.accessor.MarshalHCL(actual)
			assert.Equal(t, c.expected, actual)
		})
	}
}

func TestUserAccessor_UnmarshalHCL(t *testing.T) {
	s := new(UserAccessor).Schema()

	cases := []struct {
		name     string
		input    map[string]interface{}
		expected *UserAccessor
	}{
		{
			name:  "HCLAccessorRead-access",
			input: map[string]interface{}{"uid": "my-id", "access": HCLAccessorRead},
			expected: &UserAccessor{
				UID:    "my-id",
				Access: HCLAccessorRead,
			},
		},
		{
			name:  "write-access",
			input: map[string]interface{}{"uid": "my-id", "access": HCLAccessorWrite},
			expected: &UserAccessor{
				UID:    "my-id",
				Access: HCLAccessorWrite,
			},
		},
		{
			name:     "empty",
			input:    map[string]interface{}{"uid": "", "access": ""},
			expected: &UserAccessor{},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			d := schema.TestResourceDataRaw(t, s, c.input)
			assert.NotNil(t, d)

			var actual UserAccessor
			decoder := hcl.DecoderFrom(d)
			err := actual.UnmarshalHCL(decoder)
			assert.Equal(t, c.expected, &actual)
			assert.NoError(t, err)
		})
	}
}

func TestUsers_MarshalHCL(t *testing.T) {
	users := Users{
		{UID: "user1", Access: HCLAccessorRead},
		{UID: "user2", Access: HCLAccessorWrite},
	}

	properties := hcl.Properties{}
	err := users.MarshalHCL(properties)
	assert.NoError(t, err)

	userList, ok := properties["user"].([]interface{})
	assert.True(t, ok)
	assert.Len(t, userList, 2)

	user1 := userList[0].(hcl.Properties)
	assert.Equal(t, "user1", user1["uid"])
	assert.Equal(t, HCLAccessorRead, user1["access"])

	user2 := userList[1].(hcl.Properties)
	assert.Equal(t, "user2", user2["uid"])
	assert.Equal(t, HCLAccessorWrite, user2["access"])
}

func TestUsers_UnmarshalHCL(t *testing.T) {
	s := new(Users).Schema()

	input := map[string]interface{}{
		"user": []interface{}{
			map[string]interface{}{
				"uid":    "user1",
				"access": HCLAccessorRead,
			},
			map[string]interface{}{
				"uid":    "user2",
				"access": HCLAccessorWrite,
			},
		},
	}

	d := schema.TestResourceDataRaw(t, s, input)
	assert.NotNil(t, d)

	var users Users
	decoder := hcl.DecoderFrom(d)
	err := users.UnmarshalHCL(decoder)
	assert.NoError(t, err)
	assert.Len(t, users, 2)

	// Define expected users
	expectedUser1 := &UserAccessor{UID: "user1", Access: HCLAccessorRead}
	expectedUser2 := &UserAccessor{UID: "user2", Access: HCLAccessorWrite}

	// Check that both expected users are in the slice, regardless of order
	assert.Contains(t, users, expectedUser1)
	assert.Contains(t, users, expectedUser2)
}
