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
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

// writeTempFile creates a file with the given content inside dir and returns
// its path. The test is stopped immediately if writing fails.
func writeTempFile(t *testing.T, dir, name, content string) string {
	t.Helper()
	path := filepath.Join(dir, name)
	require.NoError(t, os.WriteFile(path, []byte(content), 0o644))
	return path
}

// writeNestedFile creates a file at the slash-separated relPath below root,
// creating any missing parent directories. Both walkers under test recurse, so
// most fixtures need more than a flat directory.
func writeNestedFile(t *testing.T, root, relPath, content string) string {
	t.Helper()
	path := filepath.Join(root, filepath.FromSlash(relPath))
	require.NoError(t, os.MkdirAll(filepath.Dir(path), 0o755))
	require.NoError(t, os.WriteFile(path, []byte(content), 0o644))
	return path
}

// writeSchemaPkg fakes a settings package by writing a schema.json below root.
// Note that version is a JSON string in the real schemas (e.g. "1"), not a number.
func writeSchemaPkg(t *testing.T, root, relPath, schemaID, version string) string {
	t.Helper()
	content := fmt.Sprintf(`{"schemaId":%q,"version":%q}`, schemaID, version)
	return writeNestedFile(t, root, relPath+"/schema.json", content)
}

// suiteName builds the JUnit suite name gotestsum produces for a package, which
// is its full Go import path.
func suiteName(relPath string) string {
	if relPath == "" {
		return moduleName
	}
	return moduleName + "/" + relPath
}
