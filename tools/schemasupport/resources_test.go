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

const sampleMarkdown = `# Supported Resources

| Resource Name | API Endpoint |
|---|---|
| dynatrace_alerting | Settings API (schema: builtin:alerting.profile) |
| dynatrace_autotag | Settings API (schema: builtin:tags.auto-tagging) |
| dynatrace_classic | Classic Config API /api/config/v1/something |
| dynatrace_multi | Settings API (schema: builtin:multi.schema) |
| dynatrace_multi_2 | Settings API (schema: builtin:multi.schema) |
`

func TestParseResourceMapping_HappyPath(t *testing.T) {
	path := writeTempFile(t, t.TempDir(), "supported-resources.md", sampleMarkdown)

	mapping, err := parseResourceMapping(path)
	require.NoError(t, err)

	assert.Equal(t, []string{"dynatrace_alerting"}, mapping["builtin:alerting.profile"])
	assert.Equal(t, []string{"dynatrace_autotag"}, mapping["builtin:tags.auto-tagging"])
}

func TestParseResourceMapping_SkipsNonSettingsRows(t *testing.T) {
	path := writeTempFile(t, t.TempDir(), "supported-resources.md", sampleMarkdown)

	mapping, err := parseResourceMapping(path)
	require.NoError(t, err)

	// Classic Config API row must not appear.
	for schemaID := range mapping {
		assert.NotContains(t, schemaID, "classic", "classic resource should have been skipped")
	}
}

func TestParseResourceMapping_MultipleResourcesSameSchema(t *testing.T) {
	path := writeTempFile(t, t.TempDir(), "supported-resources.md", sampleMarkdown)

	mapping, err := parseResourceMapping(path)
	require.NoError(t, err)

	resources := mapping["builtin:multi.schema"]
	require.Len(t, resources, 2)
	assert.Contains(t, resources, "dynatrace_multi")
	assert.Contains(t, resources, "dynatrace_multi_2")
}

func TestParseResourceMapping_SkipsHeaderAndSeparatorRows(t *testing.T) {
	md := `| Resource Name | API Endpoint |
|---|---|
| dynatrace_x | Settings API (schema: builtin:x) |
`
	path := writeTempFile(t, t.TempDir(), "supported-resources.md", md)

	mapping, err := parseResourceMapping(path)
	require.NoError(t, err)

	assert.Contains(t, mapping, "builtin:x", "expected builtin:x to be found")
	for k := range mapping {
		assert.NotContains(t, k, "Resource Name", "header row was not skipped")
	}
}

func TestParseResourceMapping_NonexistentFile(t *testing.T) {
	_, err := parseResourceMapping(filepath.Join(t.TempDir(), "missing.md"))
	assert.Error(t, err)
}

func TestParseResourceMapping_EmptyFile(t *testing.T) {
	path := writeTempFile(t, t.TempDir(), "supported-resources.md", "")

	mapping, err := parseResourceMapping(path)
	require.NoError(t, err)
	assert.Empty(t, mapping)
}
