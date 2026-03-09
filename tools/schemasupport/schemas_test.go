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
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCollectSchemas_SingleFile(t *testing.T) {
	dir := t.TempDir()
	sub := filepath.Join(dir, "alerting-profile")
	require.NoError(t, os.Mkdir(sub, 0o755))
	writeTempFile(t, sub, "schema.json", `{"schemaId":"builtin:alerting.profile","version":"1.2.3"}`)

	infos, err := collectSchemas(dir)
	require.NoError(t, err)
	require.Len(t, infos, 1)
	assert.Equal(t, "builtin:alerting.profile", infos[0].SchemaID)
	assert.Equal(t, "1.2.3", infos[0].Version)
}

func TestCollectSchemas_MultipleNestedFiles(t *testing.T) {
	dir := t.TempDir()
	schemas := []struct {
		subdir  string
		content string
	}{
		{"a", `{"schemaId":"schema-a","version":"1.0.0"}`},
		{"b/c", `{"schemaId":"schema-b","version":"2.0.0"}`},
	}
	for _, s := range schemas {
		sub := filepath.Join(dir, filepath.FromSlash(s.subdir))
		require.NoError(t, os.MkdirAll(sub, 0o755))
		writeTempFile(t, sub, "schema.json", s.content)
	}

	infos, err := collectSchemas(dir)
	require.NoError(t, err)
	require.Len(t, infos, 2)

	ids := map[string]string{}
	for _, info := range infos {
		ids[info.SchemaID] = info.Version
	}
	assert.Equal(t, "1.0.0", ids["schema-a"])
	assert.Equal(t, "2.0.0", ids["schema-b"])
}

func TestCollectSchemas_EmptyDirectory(t *testing.T) {
	infos, err := collectSchemas(t.TempDir())
	require.NoError(t, err)
	assert.Empty(t, infos)
}

func TestCollectSchemas_IgnoresNonSchemaFiles(t *testing.T) {
	dir := t.TempDir()
	writeTempFile(t, dir, "other.json", `{"schemaId":"should-be-ignored","version":"9.9.9"}`)
	writeTempFile(t, dir, "README.md", "# hello")

	infos, err := collectSchemas(dir)
	require.NoError(t, err)
	assert.Empty(t, infos)
}

func TestCollectSchemas_InvalidJSON(t *testing.T) {
	dir := t.TempDir()
	writeTempFile(t, dir, "schema.json", `{ NOT VALID JSON }`)

	_, err := collectSchemas(dir)
	assert.Error(t, err)
}

func TestCollectSchemas_NonexistentRoot(t *testing.T) {
	_, err := collectSchemas(filepath.Join(t.TempDir(), "does-not-exist"))
	assert.Error(t, err)
}
