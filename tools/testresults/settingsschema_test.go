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
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ----- parseSettingsSchemaFile -----

func TestParseSettingsSchemaFile_Valid(t *testing.T) {
	path := writeTempFile(t, t.TempDir(), "schema.json", `{"schemaId":"builtin:audit-log","version":"1"}`)

	sch, err := parseSettingsSchemaFile(path)

	require.NoError(t, err)
	assert.Equal(t, schemaIDAndVersion{SchemaID: "builtin:audit-log", Version: "1"}, sch)
}

func TestParseSettingsSchemaFile_IgnoresUnknownFields(t *testing.T) {
	path := writeTempFile(t, t.TempDir(), "schema.json",
		`{"schemaId":"builtin:audit-log","version":"1","displayName":"Audit log","properties":{}}`)

	sch, err := parseSettingsSchemaFile(path)

	require.NoError(t, err)
	assert.Equal(t, schemaIDAndVersion{SchemaID: "builtin:audit-log", Version: "1"}, sch)
}

func TestParseSettingsSchemaFile_MalformedJSON(t *testing.T) {
	path := writeTempFile(t, t.TempDir(), "schema.json", `{ NOT VALID JSON }`)

	_, err := parseSettingsSchemaFile(path)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "error parsing settings schema file")
}

func TestParseSettingsSchemaFile_MissingFile(t *testing.T) {
	_, err := parseSettingsSchemaFile(filepath.Join(t.TempDir(), "missing.json"))

	require.Error(t, err)
	assert.Contains(t, err.Error(), "error reading settings schema file")
}

// ----- findSettingsSchemaIDAndVersion -----

func TestFindSettingsSchemaIDAndVersion_Found(t *testing.T) {
	dir := t.TempDir()
	writeNestedFile(t, dir, "schema.json", `{"schemaId":"builtin:audit-log","version":"1"}`)

	sch, err := findSettingsSchemaIDAndVersion(dir)

	require.NoError(t, err)
	assert.Equal(t, schemaIDAndVersion{SchemaID: "builtin:audit-log", Version: "1"}, sch)
}

func TestFindSettingsSchemaIDAndVersion_NotFound(t *testing.T) {
	_, err := findSettingsSchemaIDAndVersion(t.TempDir())

	require.Error(t, err)
	assert.Contains(t, err.Error(), "schema.json not found in")
}

func TestFindSettingsSchemaIDAndVersion_NonexistentPath(t *testing.T) {
	_, err := findSettingsSchemaIDAndVersion(filepath.Join(t.TempDir(), "does-not-exist"))

	require.Error(t, err)
	assert.Contains(t, err.Error(), "schema.json not found in")
}
