/**
* @license
* Copyright 2023 Dynatrace LLC
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

package sensitive

import (
	"testing"
)

func TestGetBoolEnv(t *testing.T) {
	tests := []struct {
		name         string
		key          string
		value        string
		defaultValue bool
		expected     bool
	}{
		{"empty value uses default true", "TEST_VAR_SENSITIVE_1", "", true, true},
		{"empty value uses default false", "TEST_VAR_SENSITIVE_2", "", false, false},
		{"true string", "TEST_VAR_SENSITIVE_3", "true", false, true},
		{"false string", "TEST_VAR_SENSITIVE_4", "false", true, false},
		{"1 string", "TEST_VAR_SENSITIVE_5", "1", false, true},
		{"0 string", "TEST_VAR_SENSITIVE_6", "0", true, false},
		{"TRUE uppercase", "TEST_VAR_SENSITIVE_7", "TRUE", false, true},
		{"FALSE uppercase", "TEST_VAR_SENSITIVE_8", "FALSE", true, false},
		{"invalid value uses default", "TEST_VAR_SENSITIVE_9", "invalid", true, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set the value if provided (t.Setenv automatically handles cleanup)
			if tt.value != "" {
				t.Setenv(tt.key, tt.value)
			}

			result := getBoolEnv(tt.key, tt.defaultValue)
			if result != tt.expected {
				t.Errorf("getBoolEnv(%q, %v) = %v; want %v", tt.key, tt.defaultValue, result, tt.expected)
			}
		})
	}
}
