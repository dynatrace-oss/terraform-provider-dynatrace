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

// simpleEntry is a plain struct with no Preconditioner implementation.
type simpleEntry struct {
	Name  string
	Value int
}

// preconditionEntry implements Preconditioner and sets a default field on HandlePreconditions.
type preconditionEntry struct {
	Name  string
	Value int
}

func (p *preconditionEntry) HandlePreconditions() error {
	if p.Name == "" {
		p.Name = "default-name"
	}
	return nil
}

func TestFilterEmpty_RemovesDefaultValue(t *testing.T) {
	defaultVal := simpleEntry{Name: "phantom", Value: 0}
	entries := []*simpleEntry{
		{Name: "phantom", Value: 0},
		{Name: "real", Value: 42},
	}

	result := hcl.FilterEmpty(entries, defaultVal)

	assert.Len(t, result, 1)
	assert.Equal(t, "real", result[0].Name)
}

func TestFilterEmpty_RemovesZeroValue(t *testing.T) {
	defaultVal := simpleEntry{Name: "phantom"}
	entries := []*simpleEntry{
		{}, // zero value — should be removed
		{Name: "real", Value: 1},
	}

	result := hcl.FilterEmpty(entries, defaultVal)

	assert.Len(t, result, 1)
	assert.Equal(t, "real", result[0].Name)
}

func TestFilterEmpty_KeepsAllWhenNoneMatchDefault(t *testing.T) {
	defaultVal := simpleEntry{Name: "phantom"}
	entries := []*simpleEntry{
		{Name: "a", Value: 1},
		{Name: "b", Value: 2},
	}

	result := hcl.FilterEmpty(entries, defaultVal)

	assert.Len(t, result, 2)
}

func TestFilterEmpty_RemovesAllWhenAllMatchDefault(t *testing.T) {
	defaultVal := simpleEntry{Name: "phantom"}
	entries := []*simpleEntry{
		{Name: "phantom"},
		{},
	}

	result := hcl.FilterEmpty(entries, defaultVal)

	assert.Empty(t, result)
}

func TestFilterEmpty_WithPreconditioner_RemovesPostPreconditionDefault(t *testing.T) {
	// defaultValue has empty Name; after HandlePreconditions Name becomes "default-name"
	defaultVal := preconditionEntry{Name: "", Value: 0}

	entries := []*preconditionEntry{
		{Name: "default-name", Value: 0}, // matches post-precondition default → removed
		{Name: "custom", Value: 5},       // kept
	}

	result := hcl.FilterEmpty(entries, defaultVal)

	assert.Len(t, result, 1)
	assert.Equal(t, "custom", result[0].Name)
}

func TestFilterEmpty_WithPreconditioner_RemovesZeroValue(t *testing.T) {
	defaultVal := preconditionEntry{Name: "", Value: 0}

	entries := []*preconditionEntry{
		{},                        // zero value → removed
		{Name: "real", Value: 99}, // kept
	}

	result := hcl.FilterEmpty(entries, defaultVal)

	assert.Len(t, result, 1)
	assert.Equal(t, "real", result[0].Name)
}

func TestFilterEmpty_WithPreconditioner_KeepsNonDefaultEntries(t *testing.T) {
	defaultVal := preconditionEntry{Name: "", Value: 0}

	entries := []*preconditionEntry{
		{Name: "alpha", Value: 1},
		{Name: "beta", Value: 2},
	}

	result := hcl.FilterEmpty(entries, defaultVal)

	assert.Len(t, result, 2)
}
