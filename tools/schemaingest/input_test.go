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

	"github.com/dynatrace-oss/terraform-provider-dynatrace/tools/internal/schema"
)

func TestReadSchemaSupport_ParsesArray(t *testing.T) {
	path := filepath.Join(t.TempDir(), "schemaSupport.json")
	require.NoError(t, os.WriteFile(path, []byte(`[
		{"resourceName":"dynatrace_activegate_token","schemaId":"builtin:activegate-token","version":"1.1"},
		{"resourceName":"dynatrace_activegate_updates","schemaId":"builtin:deployment.activegate.updates","version":"1.0.4"}
	]`), 0o644))

	entries, err := readSchemaSupport(path)
	require.NoError(t, err)
	assert.Equal(t, []schema.Entry{
		{ResourceName: "dynatrace_activegate_token", SchemaID: "builtin:activegate-token", Version: "1.1"},
		{ResourceName: "dynatrace_activegate_updates", SchemaID: "builtin:deployment.activegate.updates", Version: "1.0.4"},
	}, entries)
}

func TestReadSchemaSupport_MissingFile(t *testing.T) {
	_, err := readSchemaSupport(filepath.Join(t.TempDir(), "does-not-exist.json"))
	assert.Error(t, err)
}

func TestReadSchemaSupport_InvalidJSON(t *testing.T) {
	path := filepath.Join(t.TempDir(), "schemaSupport.json")
	require.NoError(t, os.WriteFile(path, []byte(`not json`), 0o644))
	_, err := readSchemaSupport(path)
	assert.Error(t, err)
}

func TestToResults_MapsAndPreservesOrder(t *testing.T) {
	entries := []schema.Entry{
		{ResourceName: "dynatrace_activegate_token", SchemaID: "builtin:activegate-token", Version: "1.1"},
		{ResourceName: "dynatrace_activegate_updates", SchemaID: "builtin:deployment.activegate.updates", Version: "1.0.4"},
	}
	assert.Equal(t, []result{
		{SchemaID: "builtin:activegate-token", Version: "1.1", ResourceName: "dynatrace_activegate_token"},
		{SchemaID: "builtin:deployment.activegate.updates", Version: "1.0.4", ResourceName: "dynatrace_activegate_updates"},
	}, toResults(entries))
}

func TestToResults_Empty(t *testing.T) {
	assert.Empty(t, toResults(nil))
}
