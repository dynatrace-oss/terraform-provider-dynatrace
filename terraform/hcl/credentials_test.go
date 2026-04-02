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
	"github.com/stretchr/testify/assert"
)

func TestReplaceCredPlaceholders(t *testing.T) {
	t.Run("replaces placeholder string with tfState value", func(t *testing.T) {
		// ***NNN*** is the credential placeholder pattern
		assert.Equal(t, "my-secret", hcl.ReplaceCredPlaceholders("my-secret", "***123***"))
	})

	t.Run("all valid placeholder digit counts (3 digits)", func(t *testing.T) {
		for _, placeholder := range []string{"***000***", "***001***", "***999***"} {
			assert.Equal(t, "secret", hcl.ReplaceCredPlaceholders("secret", placeholder), "placeholder %q", placeholder)
		}
	})

	t.Run("does not replace non-placeholder string", func(t *testing.T) {
		assert.Equal(t, "api-value", hcl.ReplaceCredPlaceholders("tf-value", "api-value"))
	})

	t.Run("does not replace string with too few digits", func(t *testing.T) {
		assert.Equal(t, "***12***", hcl.ReplaceCredPlaceholders("tf-value", "***12***"))
	})

	t.Run("does not replace string with too many digits", func(t *testing.T) {
		assert.Equal(t, "***1234***", hcl.ReplaceCredPlaceholders("tf-value", "***1234***"))
	})

	t.Run("does not replace string with wrong asterisk count", func(t *testing.T) {
		assert.Equal(t, "**123***", hcl.ReplaceCredPlaceholders("tf-value", "**123***"))
	})

	t.Run("does not replace string with non-digit characters in placeholder", func(t *testing.T) {
		assert.Equal(t, "***abc***", hcl.ReplaceCredPlaceholders("tf-value", "***abc***"))
	})

	t.Run("apiState not a string when tfState is string returns apiState unchanged", func(t *testing.T) {
		assert.Equal(t, 42, hcl.ReplaceCredPlaceholders("tf-value", 42))
	})

	t.Run("map: replaces placeholder values recursively", func(t *testing.T) {
		tfState := map[string]any{
			"password": "real-password",
			"username": "real-user",
		}
		apiState := map[string]any{
			"password": "***456***",
			"username": "real-user",
		}
		hcl.ReplaceCredPlaceholders(tfState, apiState)
		assert.Equal(t, "real-password", apiState["password"])
		assert.Equal(t, "real-user", apiState["username"])
	})

	t.Run("map: key present in apiState but missing in tfState is kept unchanged", func(t *testing.T) {
		tfState := map[string]any{
			"username": "real-user",
		}
		apiState := map[string]any{
			"password": "***789***",
			"username": "real-user",
		}
		hcl.ReplaceCredPlaceholders(tfState, apiState)
		// "password" exists in apiState but not in tfState → kept as placeholder
		assert.Equal(t, "***789***", apiState["password"])
	})

	t.Run("map: apiState not a map returns apiState unchanged", func(t *testing.T) {
		tfState := map[string]any{"key": "value"}
		assert.Equal(t, "not-a-map", hcl.ReplaceCredPlaceholders(tfState, "not-a-map"))
	})

	t.Run("slice: replaces placeholder values element-wise", func(t *testing.T) {
		tfState := []any{"real-secret", "normal-value"}
		apiState := []any{"***321***", "normal-value"}
		result := hcl.ReplaceCredPlaceholders(tfState, apiState).([]any)
		assert.Equal(t, "real-secret", result[0])
		assert.Equal(t, "normal-value", result[1])
	})

	t.Run("slice: length mismatch returns apiState unchanged", func(t *testing.T) {
		tfState := []any{"a", "b"}
		apiState := []any{"***111***"}
		result := hcl.ReplaceCredPlaceholders(tfState, apiState).([]any)
		assert.Equal(t, "***111***", result[0])
	})

	t.Run("slice: apiState not a slice returns apiState unchanged", func(t *testing.T) {
		tfState := []any{"a"}
		assert.Equal(t, "not-a-slice", hcl.ReplaceCredPlaceholders(tfState, "not-a-slice"))
	})

	t.Run("deeply nested map with slice containing placeholder", func(t *testing.T) {
		tfState := map[string]any{
			"config": map[string]any{
				"tokens": []any{"token-a", "token-b"},
			},
		}
		apiState := map[string]any{
			"config": map[string]any{
				"tokens": []any{"***100***", "***200***"},
			},
		}
		result := hcl.ReplaceCredPlaceholders(tfState, apiState).(map[string]any)
		tokens := result["config"].(map[string]any)["tokens"].([]any)
		assert.Equal(t, "token-a", tokens[0])
		assert.Equal(t, "token-b", tokens[1])
	})

	t.Run("nil tfState and nil apiState returns nil", func(t *testing.T) {
		assert.Nil(t, hcl.ReplaceCredPlaceholders(nil, nil))
	})

	t.Run("nil tfState with non-nil apiState returns apiState", func(t *testing.T) {
		assert.Equal(t, "***123***", hcl.ReplaceCredPlaceholders(nil, "***123***"))
	})
}
