package settings

import (
	"testing"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
)

func TestGroupAccessor_MarshalHCL(t *testing.T) {
	cases := []struct {
		name     string
		accessor *GroupAccessor
		expected hcl.Properties
	}{
		{
			name: "read-access",
			accessor: &GroupAccessor{
				"my-id",
				HCLAccessorRead,
			},
			expected: map[string]any{"id": "my-id", "access": HCLAccessorRead},
		},
		{
			name: "write-access",
			accessor: &GroupAccessor{
				"my-id",
				HCLAccessorWrite,
			},
			expected: map[string]any{"id": "my-id", "access": HCLAccessorWrite},
		},
		{
			name:     "empty",
			accessor: &GroupAccessor{},
			expected: map[string]any{"id": "", "access": ""},
		},
		{
			name: "empty id",
			accessor: &GroupAccessor{
				"",
				HCLAccessorRead,
			},
			expected: map[string]any{"id": "", "access": HCLAccessorRead},
		},
		{
			name: "empty access",
			accessor: &GroupAccessor{
				"my-id",
				"",
			},
			expected: map[string]any{"id": "my-id", "access": ""},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var actual = hcl.Properties{}
			err := c.accessor.MarshalHCL(actual)
			assert.NoError(t, err)
			assert.Equal(t, c.expected, actual)
		})
	}
}

func TestGroupAccessor_UnmarshalHCL(t *testing.T) {
	s := new(GroupAccessor).Schema()

	cases := []struct {
		name     string
		input    map[string]interface{}
		expected *GroupAccessor
	}{
		{
			name:  "read-access",
			input: map[string]interface{}{"id": "my-id", "access": HCLAccessorRead},
			expected: &GroupAccessor{
				ID:     "my-id",
				Access: HCLAccessorRead,
			},
		},
		{
			name:  "write-access",
			input: map[string]interface{}{"id": "my-id", "access": HCLAccessorWrite},
			expected: &GroupAccessor{
				ID:     "my-id",
				Access: HCLAccessorWrite,
			},
		},
		{
			name:     "empty",
			input:    map[string]interface{}{"id": "", "access": ""},
			expected: &GroupAccessor{},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			d := schema.TestResourceDataRaw(t, s, c.input)
			assert.NotNil(t, d)

			var actual GroupAccessor
			decoder := hcl.DecoderFrom(d)
			err := actual.UnmarshalHCL(decoder)
			assert.Equal(t, c.expected, &actual)
			assert.NoError(t, err)
		})
	}
}

func TestGroups_MarshalHCL(t *testing.T) {
	groups := Groups{
		{ID: "group-1", Access: HCLAccessorRead},
		{ID: "group-2", Access: HCLAccessorWrite},
	}

	properties := hcl.Properties{}
	err := groups.MarshalHCL(properties)
	assert.NoError(t, err)

	groupList, ok := properties["group"].([]interface{})
	assert.True(t, ok)
	assert.Len(t, groupList, 2)

	group1 := groupList[0].(hcl.Properties)
	assert.Equal(t, "group-1", group1["id"])
	assert.Equal(t, HCLAccessorRead, group1["access"])

	group2 := groupList[1].(hcl.Properties)
	assert.Equal(t, "group-2", group2["id"])
	assert.Equal(t, HCLAccessorWrite, group2["access"])
}

func TestGroups_UnmarshalHCL(t *testing.T) {
	s := new(Groups).Schema()

	input := map[string]interface{}{
		"group": []interface{}{
			map[string]interface{}{
				"id":     "group-1",
				"access": HCLAccessorRead,
			},
			map[string]interface{}{
				"id":     "group-2",
				"access": HCLAccessorWrite,
			},
		},
	}

	d := schema.TestResourceDataRaw(t, s, input)
	assert.NotNil(t, d)

	var groups Groups
	decoder := hcl.DecoderFrom(d)
	err := groups.UnmarshalHCL(decoder)
	assert.NoError(t, err)
	assert.Len(t, groups, 2)

	// Define expected groups
	expectedGroup1 := &GroupAccessor{ID: "group-1", Access: HCLAccessorRead}
	expectedGroup2 := &GroupAccessor{ID: "group-2", Access: HCLAccessorWrite}

	// Check that both expected groups are in the slice, regardless of order
	assert.Contains(t, groups, expectedGroup1)
	assert.Contains(t, groups, expectedGroup2)
}
