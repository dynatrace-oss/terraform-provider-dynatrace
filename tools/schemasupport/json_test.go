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

package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ---------------------------------------------------------------------------
// expandRows
// ---------------------------------------------------------------------------

func TestExpandRows_MissingSchema(t *testing.T) {
	infos := []schemaInfo{
		{SchemaID: "schema-known", Version: "1.0"},
		{SchemaID: "schema-unknown", Version: "1.1"},
	}
	mapping := map[string][]string{
		"schema-known": {"res_known"},
	}

	rows, missing := expandRows(infos, mapping)
	require.Len(t, missing, 1)
	assert.Equal(t, "schema-unknown", missing[0])
	assert.Len(t, rows, 1)
}

func TestExpandRows_EmptyResourceSlice(t *testing.T) {
	infos := []schemaInfo{
		{SchemaID: "schema-empty", Version: "1.0"},
	}
	mapping := map[string][]string{
		"schema-empty": {},
	}

	_, missing := expandRows(infos, mapping)
	require.Len(t, missing, 1)
	assert.Equal(t, "schema-empty", missing[0], "empty resource slice should be treated as missing")
}

func TestExpandRows_RowContents(t *testing.T) {
	infos := []schemaInfo{
		{SchemaID: "schema-a", Version: "1.0"},
		{SchemaID: "schema-b", Version: "2.0"},
	}
	mapping := map[string][]string{
		"schema-a": {"res_a"},
		"schema-b": {"res_b1", "res_b2"},
	}

	entries, _ := expandRows(infos, mapping)
	require.Len(t, entries, 3)
	assert.ElementsMatch(t, []schemaEntry{
		{
			ResourceName: "res_a",
			SchemaID:     "schema-a",
			Version:      "1.0",
		},
		{
			ResourceName: "res_b1",
			SchemaID:     "schema-b",
			Version:      "2.0",
		},
		{
			ResourceName: "res_b2",
			SchemaID:     "schema-b",
			Version:      "2.0",
		},
	}, entries)
}

func TestExpandRows_EmptyInputs(t *testing.T) {
	rows, missing := expandRows(nil, map[string][]string{})
	assert.Empty(t, rows)
	assert.Empty(t, missing)
}

// ---------------------------------------------------------------------------
// sortRows
// ---------------------------------------------------------------------------

func TestSortRows_ByResourceName(t *testing.T) {
	entries := []schemaEntry{
		{ResourceName: "zzz", SchemaID: "s1", Version: "1.0"},
		{ResourceName: "aaa", SchemaID: "s2", Version: "1.0"},
		{ResourceName: "mmm", SchemaID: "s3", Version: "1.0"},
	}
	sortRows(entries)
	assert.Equal(t, "aaa", entries[0].ResourceName)
	assert.Equal(t, "mmm", entries[1].ResourceName)
	assert.Equal(t, "zzz", entries[2].ResourceName)
}

func TestSortRows_SameResourceNameThenBySchemaId(t *testing.T) {
	entries := []schemaEntry{
		{ResourceName: "res", SchemaID: "z-schema", Version: "1.0"},
		{ResourceName: "res", SchemaID: "a-schema", Version: "1.0"},
		{ResourceName: "res", SchemaID: "m-schema", Version: "1.0"},
	}
	sortRows(entries)
	assert.Equal(t, "a-schema", entries[0].SchemaID)
	assert.Equal(t, "m-schema", entries[1].SchemaID)
	assert.Equal(t, "z-schema", entries[2].SchemaID)
}

func TestSortRows_AlreadySorted(t *testing.T) {
	entries := []schemaEntry{
		{ResourceName: "a", SchemaID: "s1", Version: "1.0"},
		{ResourceName: "b", SchemaID: "s2", Version: "1.0"},
	}
	sortRows(entries)
	assert.Equal(t, "a", entries[0].ResourceName)
	assert.Equal(t, "b", entries[1].ResourceName)
}

func TestSortRows_Empty(t *testing.T) {
	assert.NotPanics(t, func() { sortRows(nil) })
	assert.NotPanics(t, func() { sortRows([]schemaEntry{}) })
}

// ---------------------------------------------------------------------------
// writeJSON
// ---------------------------------------------------------------------------

func TestWriteJSON_HeaderAndRows(t *testing.T) {
	path := filepath.Join(t.TempDir(), "out.json")

	entries := []schemaEntry{
		{ResourceName: "res_a", SchemaID: "schema-a", Version: "1.0"},
		{ResourceName: "res_b", SchemaID: "schema-b", Version: "2.0"},
	}
	require.NoError(t, writeJSON(path, entries))

	content, err := os.ReadFile(path)
	require.NoError(t, err)

	var got []schemaEntry
	require.NoError(t, json.Unmarshal(content, &got))

	require.Len(t, got, 2)
	assert.Equal(t, entries[0], got[0])
	assert.Equal(t, entries[1], got[1])
}

func TestWriteJSON_EmptyRows(t *testing.T) {
	path := filepath.Join(t.TempDir(), "out.json")
	require.NoError(t, writeJSON(path, nil))

	content, err := os.ReadFile(path)
	require.NoError(t, err)

	var got []schemaEntry
	require.NoError(t, json.Unmarshal(content, &got))
	assert.Empty(t, got)
}

func TestWriteJSON_ValidJSON(t *testing.T) {
	path := filepath.Join(t.TempDir(), "out.json")
	require.NoError(t, writeJSON(path, []schemaEntry{{ResourceName: "r", SchemaID: "s", Version: "v"}}))

	content, err := os.ReadFile(path)
	require.NoError(t, err)
	assert.True(t, json.Valid(content), "expected valid JSON in output file")
}

func TestWriteJSON_InvalidPath(t *testing.T) {
	err := writeJSON(filepath.Join(t.TempDir(), "nonexistent", "sub", "out.json"), nil)
	assert.Error(t, err)
}

// ---------------------------------------------------------------------------
// formatMissingError
// ---------------------------------------------------------------------------

func TestFormatMissingError_ContainsAllIDs(t *testing.T) {
	missing := []string{"schema-z", "schema-a", "schema-m"}
	msg := formatMissingError(missing)

	for _, id := range missing {
		assert.Contains(t, msg, id)
	}
}

func TestFormatMissingError_IsSorted(t *testing.T) {
	msg := formatMissingError([]string{"zzz", "aaa", "mmm"})

	idxA := strings.Index(msg, "aaa")
	idxM := strings.Index(msg, "mmm")
	idxZ := strings.Index(msg, "zzz")

	assert.Less(t, idxA, idxM, "aaa should appear before mmm")
	assert.Less(t, idxM, idxZ, "mmm should appear before zzz")
}

func TestFormatMissingError_ContainsBulletPoints(t *testing.T) {
	msg := formatMissingError([]string{"x"})
	assert.Contains(t, msg, "- x")
}

func TestFormatMissingError_Empty(t *testing.T) {
	msg := formatMissingError(nil)
	assert.Contains(t, msg, "schemaIds")
}
