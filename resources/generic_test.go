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

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/export"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/provider/config"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/resources"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/dynatrace/dynatrace-configuration-as-code-core/api"
)

func TestGeneric_Read(t *testing.T) {
	t.Run("resource not found error should set the ID to empty", func(t *testing.T) {
		gen := resources.Generic{
			Descriptor: export.NewResourceDescriptor(MockService(nil, rest.Error{Code: http.StatusNotFound})),
		}
		d := schema.TestResourceDataRaw(t, new(mockSchema).Schema(), nil)
		d.SetId("test")
		m := &config.ProviderConfiguration{
			EnvironmentURL: "https://dynatrace.com",
		}

		diags := gen.Read(t.Context(), d, m)
		assert.False(t, diags.HasError())
		assert.Empty(t, d.Id())
	})

	t.Run("resource not found core error should set the ID to empty", func(t *testing.T) {
		gen := resources.Generic{
			Descriptor: export.NewResourceDescriptor(MockService(nil, api.APIError{StatusCode: http.StatusNotFound})),
		}
		d := schema.TestResourceDataRaw(t, new(mockSchema).Schema(), nil)
		d.SetId("test")
		m := &config.ProviderConfiguration{
			EnvironmentURL: "https://dynatrace.com",
		}

		diags := gen.Read(t.Context(), d, m)
		assert.False(t, diags.HasError())
		assert.Empty(t, d.Id())
	})

	t.Run("Non not found errors should error", func(t *testing.T) {
		gen := resources.Generic{
			Descriptor: export.NewResourceDescriptor(MockService(nil, rest.Error{Code: http.StatusBadRequest, Message: "badRequest"})),
		}
		d := schema.TestResourceDataRaw(t, new(mockSchema).Schema(), nil)
		d.SetId("test")
		m := &config.ProviderConfiguration{
			EnvironmentURL: "https://dynatrace.com",
		}

		diags := gen.Read(t.Context(), d, m)
		assert.True(t, diags.HasError())
		require.GreaterOrEqual(t, 1, len(diags))
		el := diags[0]
		assert.Equal(t, "badRequest", el.Summary)
		assert.NotEmpty(t, d.Id())
	})
}

func TestGeneric_Delete(t *testing.T) {
	t.Run("resource not found error should set the ID to empty", func(t *testing.T) {
		gen := resources.Generic{
			Descriptor: export.NewResourceDescriptor(MockService(nil, rest.Error{Code: http.StatusNotFound})),
		}
		d := schema.TestResourceDataRaw(t, new(mockSchema).Schema(), nil)
		d.SetId("test")
		m := &config.ProviderConfiguration{
			EnvironmentURL: "https://dynatrace.com",
		}

		diags := gen.Delete(t.Context(), d, m)
		assert.False(t, diags.HasError())
		assert.Empty(t, d.Id())
	})

	t.Run("resource not found core error should set the ID to empty", func(t *testing.T) {
		gen := resources.Generic{
			Descriptor: export.NewResourceDescriptor(MockService(nil, api.APIError{StatusCode: http.StatusNotFound})),
		}
		d := schema.TestResourceDataRaw(t, new(mockSchema).Schema(), nil)
		d.SetId("test")
		m := &config.ProviderConfiguration{
			EnvironmentURL: "https://dynatrace.com",
		}

		diags := gen.Delete(t.Context(), d, m)
		assert.False(t, diags.HasError())
		assert.Empty(t, d.Id())
	})

	t.Run("Non not found errors should error", func(t *testing.T) {
		gen := resources.Generic{
			Descriptor: export.NewResourceDescriptor(MockService(nil, rest.Error{Code: http.StatusBadRequest, Message: "badRequest"})),
		}
		d := schema.TestResourceDataRaw(t, new(mockSchema).Schema(), nil)
		d.SetId("test")
		m := &config.ProviderConfiguration{
			EnvironmentURL: "https://dynatrace.com",
		}

		diags := gen.Delete(t.Context(), d, m)
		assert.True(t, diags.HasError())
		require.GreaterOrEqual(t, 1, len(diags))
		el := diags[0]
		assert.Equal(t, "badRequest", el.Summary)
		assert.NotEmpty(t, d.Id())
	})
}
