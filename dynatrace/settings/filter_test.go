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

package settings_test

import (
	"testing"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings"
	"github.com/stretchr/testify/assert"
)

func TestEscapeFilterString(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{name: "plain value is unchanged", input: "my-pipeline", expected: "my-pipeline"},
		{name: "empty value", input: "", expected: ""},
		{name: "single quote", input: "it's", expected: `it\'s`},
		{name: "backslash", input: `a\b`, expected: `a\\b`},
		// Without escaping the backslash, the trailing "\" would escape the closing quote of the literal.
		{name: "trailing backslash", input: `a\`, expected: `a\\`},
		// The backslash must be escaped before the quote, otherwise the quote's escape would be doubled.
		{name: "backslash followed by quote", input: `\'`, expected: `\\\'`},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expected, settings.EscapeFilterString(tc.input))
		})
	}
}
