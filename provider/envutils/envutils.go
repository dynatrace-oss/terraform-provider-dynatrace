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

// ClampedIntEnvVar reads an integer environment variable and clamps it
// to the [Min, Max] range. If the env var is not set or cannot be parsed,
// the DefaultValue is returned.
type ClampedIntEnvVar struct {
	Key          string
	DefaultValue int
	Min          int
	Max          int
}

func (c ClampedIntEnvVar) Get() int {
	v, found := os.LookupEnv(c.Key)
	if !found || len(v) == 0 {
		return c.DefaultValue
	}

	value, err := strconv.Atoi(v)
	if err != nil {
		return c.DefaultValue
	}
	if value > c.Max {
		value = c.Max
	}
	if value < c.Min {
		value = c.Min
	}
	return value
}

// BoundedIntEnvVar reads an integer environment variable and returns the
// DefaultValue if the value is outside the [Min, Max] range (unlike
// ClampedIntEnvVar which clamps to the boundary).
type BoundedIntEnvVar struct {
	Key          string
	DefaultValue int
	Min          int
	Max          int
}

func (b BoundedIntEnvVar) Get() int {
	v, found := os.LookupEnv(b.Key)
	if !found || len(v) == 0 {
		return b.DefaultValue
	}

	value, err := strconv.Atoi(v)
	if err != nil {
		return b.DefaultValue
	}
	if value < b.Min || value > b.Max {
		return b.DefaultValue
	}
	return value
}
