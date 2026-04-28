/**
* @license
* Copyright 2026 Dynatrace LLC
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

package envutils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBoolEnvVar_Get(t *testing.T) {
	t.Run("returns default when env var is not set", func(t *testing.T) {
		b := BoolEnvVar{Key: "TEST_BOOL_UNSET", DefaultValue: true}
		assert.Equal(t, true, b.Get())
	})

	t.Run("returns true when env var is 'true'", func(t *testing.T) {
		t.Setenv("TEST_BOOL_TRUE", "true")
		b := BoolEnvVar{Key: "TEST_BOOL_TRUE", DefaultValue: false}
		assert.Equal(t, true, b.Get())
	})

	t.Run("returns false when env var is 'false'", func(t *testing.T) {
		t.Setenv("TEST_BOOL_FALSE", "false")
		b := BoolEnvVar{Key: "TEST_BOOL_FALSE", DefaultValue: true}
		assert.Equal(t, false, b.Get())
	})

	t.Run("returns true when env var is '1'", func(t *testing.T) {
		t.Setenv("TEST_BOOL_ONE", "1")
		b := BoolEnvVar{Key: "TEST_BOOL_ONE", DefaultValue: false}
		assert.Equal(t, true, b.Get())
	})

	t.Run("returns false when env var is '0'", func(t *testing.T) {
		t.Setenv("TEST_BOOL_ZERO", "0")
		b := BoolEnvVar{Key: "TEST_BOOL_ZERO", DefaultValue: true}
		assert.Equal(t, false, b.Get())
	})

	t.Run("returns default when env var is invalid", func(t *testing.T) {
		t.Setenv("TEST_BOOL_INVALID", "notabool")
		b := BoolEnvVar{Key: "TEST_BOOL_INVALID", DefaultValue: true}
		assert.Equal(t, true, b.Get())
	})

	t.Run("returns default (false) when env var is invalid", func(t *testing.T) {
		t.Setenv("TEST_BOOL_INVALID2", "notabool")
		b := BoolEnvVar{Key: "TEST_BOOL_INVALID2", DefaultValue: false}
		assert.Equal(t, false, b.Get())
	})
}

func TestStringEnvVar_Get(t *testing.T) {
	t.Run("returns default when env var is not set", func(t *testing.T) {
		s := StringEnvVar{Key: "TEST_STRING_UNSET", DefaultValue: "default"}
		assert.Equal(t, "default", s.Get())
	})

	t.Run("returns env var value when set", func(t *testing.T) {
		t.Setenv("TEST_STRING_SET", "custom")
		s := StringEnvVar{Key: "TEST_STRING_SET", DefaultValue: "default"}
		assert.Equal(t, "custom", s.Get())
	})

	t.Run("returns empty string when env var is set to empty", func(t *testing.T) {
		t.Setenv("TEST_STRING_EMPTY", "")
		s := StringEnvVar{Key: "TEST_STRING_EMPTY", DefaultValue: "default"}
		assert.Equal(t, "", s.Get())
	})

	t.Run("returns empty default when env var is not set", func(t *testing.T) {
		s := StringEnvVar{Key: "TEST_STRING_EMPTY_DEFAULT", DefaultValue: ""}
		assert.Equal(t, "", s.Get())
	})
}

func TestClampedIntEnvVar_Get(t *testing.T) {
	t.Run("returns default when env var is not set", func(t *testing.T) {
		c := ClampedIntEnvVar{Key: "TEST_CLAMPED_UNSET", DefaultValue: 50, Min: 10, Max: 100}
		assert.Equal(t, 50, c.Get())
	})

	t.Run("returns parsed value when within range", func(t *testing.T) {
		t.Setenv("TEST_CLAMPED_VALID", "75")
		c := ClampedIntEnvVar{Key: "TEST_CLAMPED_VALID", DefaultValue: 50, Min: 10, Max: 100}
		assert.Equal(t, 75, c.Get())
	})

	t.Run("clamps to max when value exceeds max", func(t *testing.T) {
		t.Setenv("TEST_CLAMPED_HIGH", "200")
		c := ClampedIntEnvVar{Key: "TEST_CLAMPED_HIGH", DefaultValue: 50, Min: 10, Max: 100}
		assert.Equal(t, 100, c.Get())
	})

	t.Run("clamps to min when value is below min", func(t *testing.T) {
		t.Setenv("TEST_CLAMPED_LOW", "5")
		c := ClampedIntEnvVar{Key: "TEST_CLAMPED_LOW", DefaultValue: 50, Min: 10, Max: 100}
		assert.Equal(t, 10, c.Get())
	})

	t.Run("returns default when env var is invalid", func(t *testing.T) {
		t.Setenv("TEST_CLAMPED_INVALID", "notanint")
		c := ClampedIntEnvVar{Key: "TEST_CLAMPED_INVALID", DefaultValue: 50, Min: 10, Max: 100}
		assert.Equal(t, 50, c.Get())
	})

	t.Run("returns default when env var is empty", func(t *testing.T) {
		t.Setenv("TEST_CLAMPED_EMPTY", "")
		c := ClampedIntEnvVar{Key: "TEST_CLAMPED_EMPTY", DefaultValue: 50, Min: 10, Max: 100}
		assert.Equal(t, 50, c.Get())
	})
}

func TestMultiStringEnvVar_Get(t *testing.T) {
	t.Run("returns default when no env vars are set", func(t *testing.T) {
		m := MultiStringEnvVar{Keys: []string{"TEST_MULTI_A", "TEST_MULTI_B"}, DefaultValue: "fallback"}
		assert.Equal(t, "fallback", m.Get())
	})

	t.Run("returns first set env var", func(t *testing.T) {
		t.Setenv("TEST_MULTI_FIRST_A", "valueA")
		t.Setenv("TEST_MULTI_FIRST_B", "valueB")
		m := MultiStringEnvVar{Keys: []string{"TEST_MULTI_FIRST_A", "TEST_MULTI_FIRST_B"}}
		assert.Equal(t, "valueA", m.Get())
	})

	t.Run("skips empty env vars and returns first non-empty", func(t *testing.T) {
		t.Setenv("TEST_MULTI_SKIP_A", "")
		t.Setenv("TEST_MULTI_SKIP_B", "valueB")
		m := MultiStringEnvVar{Keys: []string{"TEST_MULTI_SKIP_A", "TEST_MULTI_SKIP_B"}}
		assert.Equal(t, "valueB", m.Get())
	})

	t.Run("returns empty default when no keys match", func(t *testing.T) {
		m := MultiStringEnvVar{Keys: []string{"TEST_MULTI_NONE_A", "TEST_MULTI_NONE_B"}}
		assert.Equal(t, "", m.Get())
	})
}

func TestBoundedIntEnvVar_Get(t *testing.T) {
	t.Run("returns default when env var is not set", func(t *testing.T) {
		b := BoundedIntEnvVar{Key: "TEST_BOUNDED_UNSET", DefaultValue: 50, Min: 10, Max: 100}
		assert.Equal(t, 50, b.Get())
	})

	t.Run("returns parsed value when within range", func(t *testing.T) {
		t.Setenv("TEST_BOUNDED_VALID", "75")
		b := BoundedIntEnvVar{Key: "TEST_BOUNDED_VALID", DefaultValue: 50, Min: 10, Max: 100}
		assert.Equal(t, 75, b.Get())
	})

	t.Run("returns default when value exceeds max", func(t *testing.T) {
		t.Setenv("TEST_BOUNDED_HIGH", "200")
		b := BoundedIntEnvVar{Key: "TEST_BOUNDED_HIGH", DefaultValue: 50, Min: 10, Max: 100}
		assert.Equal(t, 50, b.Get())
	})

	t.Run("returns default when value is below min", func(t *testing.T) {
		t.Setenv("TEST_BOUNDED_LOW", "5")
		b := BoundedIntEnvVar{Key: "TEST_BOUNDED_LOW", DefaultValue: 50, Min: 10, Max: 100}
		assert.Equal(t, 50, b.Get())
	})

	t.Run("returns default when env var is invalid", func(t *testing.T) {
		t.Setenv("TEST_BOUNDED_INVALID", "notanint")
		b := BoundedIntEnvVar{Key: "TEST_BOUNDED_INVALID", DefaultValue: 50, Min: 10, Max: 100}
		assert.Equal(t, 50, b.Get())
	})

	t.Run("returns default when env var is empty", func(t *testing.T) {
		t.Setenv("TEST_BOUNDED_EMPTY", "")
		b := BoundedIntEnvVar{Key: "TEST_BOUNDED_EMPTY", DefaultValue: 50, Min: 10, Max: 100}
		assert.Equal(t, 50, b.Get())
	})

	t.Run("returns value at exact min boundary", func(t *testing.T) {
		t.Setenv("TEST_BOUNDED_MIN", "10")
		b := BoundedIntEnvVar{Key: "TEST_BOUNDED_MIN", DefaultValue: 50, Min: 10, Max: 100}
		assert.Equal(t, 10, b.Get())
	})

	t.Run("returns value at exact max boundary", func(t *testing.T) {
		t.Setenv("TEST_BOUNDED_MAX", "100")
		b := BoundedIntEnvVar{Key: "TEST_BOUNDED_MAX", DefaultValue: 50, Min: 10, Max: 100}
		assert.Equal(t, 100, b.Get())
	})
}
