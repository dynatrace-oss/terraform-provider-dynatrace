//go:build unit

/*
 * @license
 * Copyright 2026 Dynatrace LLC
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package hcl_test

import (
	"testing"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
)

// assertComputed asserts that the given schema item is in the expected "computed" state.
func assertComputed(t *testing.T, name string, item *schema.Schema) {
	t.Helper()
	assert.True(t, item.Computed, "%s: expected Computed to be true", name)
	assert.False(t, item.Optional, "%s: expected Optional to be false", name)
	assert.False(t, item.Required, "%s: expected Required to be false", name)
	assert.Zero(t, item.MinItems, "%s: expected MinItems to be 0", name)
	assert.Zero(t, item.MaxItems, "%s: expected MaxItems to be 0", name)
}

func TestSetComputedSchema_SetsTopLevelItemsComputed(t *testing.T) {
	in := map[string]*schema.Schema{
		"name": {
			Type:     schema.TypeString,
			Required: true,
		},
		"count": {
			Type:     schema.TypeInt,
			Optional: true,
			MinItems: 1,
			MaxItems: 5,
		},
	}

	out := hcl.SetComputedSchema(in)

	for name, item := range out {
		assertComputed(t, name, item)
	}
}

func TestSetComputedSchema_ReturnsSameMap(t *testing.T) {
	in := map[string]*schema.Schema{
		"name": {Type: schema.TypeString, Required: true},
	}

	out := hcl.SetComputedSchema(in)

	// SetComputedSchema mutates and returns the provided map.
	assert.Len(t, out, len(in))
	assert.Same(t, in["name"], out["name"], "expected the schema items to be mutated in place")
}

func TestSetComputedSchema_RecursesIntoNestedResource(t *testing.T) {
	in := map[string]*schema.Schema{
		"block": {
			Type:     schema.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"inner": {
						Type:     schema.TypeString,
						Required: true,
					},
				},
			},
		},
	}

	out := hcl.SetComputedSchema(in)

	block := out["block"]
	assertComputed(t, "block", block)

	nested, ok := block.Elem.(*schema.Resource)
	assert.True(t, ok, "expected nested Elem to remain a *schema.Resource")
	assertComputed(t, "block.inner", nested.Schema["inner"])
}

func TestSetComputedSchema_RecursesIntoDeeplyNestedResource(t *testing.T) {
	in := map[string]*schema.Schema{
		"outer": {
			Type:     schema.TypeList,
			Optional: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"middle": {
						Type:     schema.TypeList,
						Required: true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"leaf": {
									Type:     schema.TypeBool,
									Optional: true,
								},
							},
						},
					},
				},
			},
		},
	}

	out := hcl.SetComputedSchema(in)

	outer := out["outer"].Elem.(*schema.Resource)
	middle := outer.Schema["middle"]
	assertComputed(t, "outer.middle", middle)

	leaf := middle.Elem.(*schema.Resource).Schema["leaf"]
	assertComputed(t, "outer.middle.leaf", leaf)
}

func TestSetComputedSchema_LeavesNonResourceElemUntouched(t *testing.T) {
	elem := &schema.Schema{
		Type:     schema.TypeString,
		Required: true,
		MaxItems: 3,
	}
	in := map[string]*schema.Schema{
		"tags": {
			Type:     schema.TypeSet,
			Optional: true,
			Elem:     elem,
		},
	}

	out := hcl.SetComputedSchema(in)

	// The top-level item is made computed...
	assertComputed(t, "tags", out["tags"])

	// ...but a plain *schema.Schema Elem is not recursed into.
	assert.Same(t, elem, out["tags"].Elem, "expected the Elem pointer to be preserved")
	assert.False(t, elem.Computed, "expected non-resource Elem to remain unchanged")
	assert.True(t, elem.Required, "expected non-resource Elem Required to remain true")
}

func TestSetComputedSchema_EmptySchema(t *testing.T) {
	in := map[string]*schema.Schema{}

	out := hcl.SetComputedSchema(in)

	assert.Empty(t, out, "expected an empty schema to remain empty")
}
