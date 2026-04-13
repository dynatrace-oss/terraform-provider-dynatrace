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
	"os"
	"strconv"
)

// BoolEnvVar represents a boolean environment variable with a default value.
type BoolEnvVar struct {
	Key          string
	DefaultValue bool
}

// Get returns the environment variable's boolean value or its default value if unset or invalid.
func (b BoolEnvVar) Get() bool {
	v, found := os.LookupEnv(b.Key)
	if !found {
		return b.DefaultValue
	}

	value, err := strconv.ParseBool(v)
	if err != nil {
		return b.DefaultValue
	}
	return value
}

// StringEnvVar represents a string environment variable with a default value.
type StringEnvVar struct {
	Key          string
	DefaultValue string
}

// Get returns the environment variable's string value or its default value if unset.
func (s StringEnvVar) Get() string {
	v, found := os.LookupEnv(s.Key)
	if !found {
		return s.DefaultValue
	}
	return v
}
