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
	rest2 "github.com/dynatrace/dynatrace-configuration-as-code-core/api/rest"
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
		assert.Equal(t, "API error: badRequest", el.Summary)
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
		assert.Equal(t, "API error: badRequest", el.Summary)
		assert.NotEmpty(t, d.Id())
	})
}

func TestGeneric_Create(t *testing.T) {
	t.Run("returns the error violation message", func(t *testing.T) {
		gen := resources.Generic{
			Descriptor: export.NewResourceDescriptor(MockService(nil, rest.Error{Code: http.StatusBadRequest, ConstraintViolations: []rest.ConstraintViolation{
				{
					Message: "My message",
					Path:    "my/path",
				},
			}})),
		}
		d := schema.TestResourceDataRaw(t, new(mockSchema).Schema(), nil)
		m := &config.ProviderConfiguration{
			EnvironmentURL: "https://dynatrace.com",
		}

		diags := gen.Create(t.Context(), d, m)
		assert.True(t, diags.HasError())
		require.GreaterOrEqual(t, 1, len(diags))
		el := diags[0]
		assert.Equal(t, "API error: my/path: My message", el.Summary)
	})

	t.Run("returns the error message", func(t *testing.T) {
		gen := resources.Generic{
			Descriptor: export.NewResourceDescriptor(MockService(nil, rest.Error{Code: http.StatusBadRequest, Message: "My message"})),
		}
		d := schema.TestResourceDataRaw(t, new(mockSchema).Schema(), nil)
		m := &config.ProviderConfiguration{
			EnvironmentURL: "https://dynatrace.com",
		}

		diags := gen.Create(t.Context(), d, m)
		assert.True(t, diags.HasError())
		require.GreaterOrEqual(t, 1, len(diags))
		el := diags[0]
		assert.Equal(t, "API error: My message", el.Summary)
	})

	t.Run("returns the default error message", func(t *testing.T) {
		apiErr := api.APIError{
			StatusCode: http.StatusBadRequest,
			Body:       []byte("My message"),
			Request: rest2.RequestInfo{
				Method: "POST",
				URL:    "https://dynatrace.com/api/v1/something",
			},
		}
		gen := resources.Generic{
			Descriptor: export.NewResourceDescriptor(MockService(nil, apiErr)),
		}
		d := schema.TestResourceDataRaw(t, new(mockSchema).Schema(), nil)
		m := &config.ProviderConfiguration{
			EnvironmentURL: "https://dynatrace.com",
		}

		diags := gen.Create(t.Context(), d, m)
		assert.True(t, diags.HasError())
		require.GreaterOrEqual(t, 1, len(diags))
		el := diags[0]
		assert.Equal(t, "API error: API request HTTP POST https://dynatrace.com/api/v1/something failed with status code 400: My message", el.Summary)
	})
}
