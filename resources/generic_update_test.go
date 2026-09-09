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

package resources_test

import (
	"net/http"
	"testing"

	ctyjson "github.com/hashicorp/go-cty/cty/json"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/export"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/resources"
)

// These tests mimic the SDK during Apply: the applied (new) values live in the
// attributes, while GetRawState() returns the prior state. That prior state is
// what a failed Update restores from.

func TestGeneric_Update_RevertsStateOnFailure(t *testing.T) {
	gen := resources.Generic{
		Descriptor: export.NewResourceDescriptor(MockService(nil, rest.Error{Code: http.StatusBadRequest, Message: "rejected"})),
	}
	res := gen.Resource()

	// prior state = the last good values
	implied := res.CoreConfigSchema().ImpliedType()
	priorRaw, err := ctyjson.Unmarshal([]byte(`{"id":"test","name":"old"}`), implied)
	require.NoError(t, err)

	// attributes hold the planned (new, rejected) value, as the SDK would have
	// applied the diff before calling Update
	st := &terraform.InstanceState{
		ID:         "test",
		Attributes: map[string]string{"id": "test", "name": "new"},
		RawState:   priorRaw,
	}
	d := res.Data(st)
	require.Equal(t, "new", d.Get("name"), "precondition: planned value is applied")

	m := getProviderConfigurationWithFixedEnvironmentURL(t)
	diags := gen.Update(t.Context(), d, m)

	assert.True(t, diags.HasError(), "failed PUT must surface an error")
	assert.Equal(t, "old", d.Get("name"), "state must be reverted to the prior value")
	assert.Equal(t, "old", d.State().Attributes["name"], "persisted state must hold the prior value")
}
