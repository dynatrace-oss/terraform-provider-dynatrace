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
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const auditlogPkg = "dynatrace/api/builtin/auditlog"

// ----- collectResults -----

func TestCollectResults_SinglePassingSuite(t *testing.T) {
	source := t.TempDir()
	writeSchemaPkg(t, source, auditlogPkg, "builtin:audit-log", "1")

	results, errs := collectResults([]testSuite{
		{Name: suiteName(auditlogPkg), Tests: 3, Failures: 0},
	}, source)

	assert.Empty(t, errs)
	assert.Equal(t, []result{{SchemaID: "builtin:audit-log", Version: "1", Status: "PASS"}}, results)
}

func TestCollectResults_FailingSuite(t *testing.T) {
	source := t.TempDir()
	writeSchemaPkg(t, source, auditlogPkg, "builtin:audit-log", "1")

	results, errs := collectResults([]testSuite{
		{Name: suiteName(auditlogPkg), Tests: 3, Failures: 1},
	}, source)

	assert.Empty(t, errs)
	assert.Equal(t, []result{{SchemaID: "builtin:audit-log", Version: "1", Status: "FAIL"}}, results)
}

func TestCollectResults_SameSchemaIsDedupedAndDowngraded(t *testing.T) {
	source := t.TempDir()
	writeSchemaPkg(t, source, auditlogPkg, "builtin:audit-log", "1")

	results, errs := collectResults([]testSuite{
		{Name: suiteName(auditlogPkg), Tests: 3, Failures: 0},
		{Name: suiteName(auditlogPkg), Tests: 2, Failures: 1},
	}, source)

	assert.Empty(t, errs)
	assert.Equal(t, []result{{SchemaID: "builtin:audit-log", Version: "1", Status: "FAIL"}}, results)
}

func TestCollectResults_FailureIsNotResetByLaterPass(t *testing.T) {
	source := t.TempDir()
	writeSchemaPkg(t, source, auditlogPkg, "builtin:audit-log", "1")

	results, errs := collectResults([]testSuite{
		{Name: suiteName(auditlogPkg), Tests: 2, Failures: 1},
		{Name: suiteName(auditlogPkg), Tests: 3, Failures: 0},
	}, source)

	assert.Empty(t, errs)
	assert.Equal(t, []result{{SchemaID: "builtin:audit-log", Version: "1", Status: "FAIL"}}, results)
}

func TestCollectResults_MultipleSchemas(t *testing.T) {
	source := t.TempDir()
	writeSchemaPkg(t, source, auditlogPkg, "builtin:audit-log", "1")
	writeSchemaPkg(t, source, "dynatrace/api/builtin/alerting/profile", "builtin:alerting.profile", "2")

	results, errs := collectResults([]testSuite{
		{Name: suiteName(auditlogPkg), Tests: 3, Failures: 0},
		{Name: suiteName("dynatrace/api/builtin/alerting/profile"), Tests: 1, Failures: 1},
	}, source)

	assert.Empty(t, errs)
	assert.ElementsMatch(t, []result{
		{SchemaID: "builtin:audit-log", Version: "1", Status: "PASS"},
		{SchemaID: "builtin:alerting.profile", Version: "2", Status: "FAIL"},
	}, results)
}

func TestCollectResults_ZeroTestSuiteIsSkippedEvenWithSchema(t *testing.T) {
	source := t.TempDir()
	writeSchemaPkg(t, source, auditlogPkg, "builtin:audit-log", "1")

	results, errs := collectResults([]testSuite{
		{Name: suiteName(auditlogPkg), Tests: 0, Failures: 0},
	}, source)

	assert.Empty(t, errs)
	assert.Empty(t, results)
}

func TestCollectResults_AllSkippedSuiteIsReportedAsPassed(t *testing.T) {
	source := t.TempDir()
	writeSchemaPkg(t, source, auditlogPkg, "builtin:audit-log", "1")

	input := t.TempDir()
	writeTempFile(t, input, "test-results.xml", fmt.Sprintf(
		`<testsuites><testsuite name="%s" tests="3" failures="0" skipped="3"></testsuite></testsuites>`,
		suiteName(auditlogPkg)))
	suites, errs := findAllTestSuites(input)
	require.Empty(t, errs)

	results, errs := collectResults(suites, source)

	assert.Empty(t, errs)
	assert.Equal(t, []result{{SchemaID: "builtin:audit-log", Version: "1", Status: "PASS"}}, results)
}

func TestCollectResults_SuiteNameWithoutModulePrefix(t *testing.T) {
	results, errs := collectResults([]testSuite{
		{Name: "example.com/other/pkg", Tests: 3, Failures: 0},
	}, t.TempDir())

	assert.Empty(t, results)
	require.Len(t, errs, 1)
	assert.Contains(t, errs[0].Error(), "failed to convert suite name to path")
}

func TestCollectResults_PackageWithoutSchema(t *testing.T) {
	source := t.TempDir()
	writeNestedFile(t, source, "dynatrace/rest/client.go", "package rest")

	results, errs := collectResults([]testSuite{
		{Name: suiteName("dynatrace/rest"), Tests: 3, Failures: 0},
	}, source)

	assert.Empty(t, results)
	require.Len(t, errs, 1)
	assert.Contains(t, errs[0].Error(), "failed to find settings schema")
}

func TestCollectResults_BadSuiteDoesNotDiscardGoodResults(t *testing.T) {
	source := t.TempDir()
	writeSchemaPkg(t, source, auditlogPkg, "builtin:audit-log", "1")

	results, errs := collectResults([]testSuite{
		{Name: "example.com/other/pkg", Tests: 1, Failures: 0},
		{Name: suiteName(auditlogPkg), Tests: 3, Failures: 0},
	}, source)

	require.Len(t, errs, 1)
	assert.Equal(t, []result{{SchemaID: "builtin:audit-log", Version: "1", Status: "PASS"}}, results)
}

func TestCollectResults_EmptyInput(t *testing.T) {
	results, errs := collectResults(nil, t.TempDir())

	assert.Empty(t, errs)
	assert.Empty(t, results)
}

// ----- writeResultsAsEvent -----

func TestWriteResultsAsEvent_WritesFullPayload(t *testing.T) {
	t.Setenv("GITHUB_REPOSITORY", "dynatrace-oss/terraform-provider-dynatrace")
	t.Setenv("GITHUB_RUN_ID", "123456")
	output := filepath.Join(t.TempDir(), "test-report.json")

	err := writeResultsAsEvent([]result{
		{SchemaID: "builtin:audit-log", Version: "1", Status: "PASS"},
	}, output)

	require.NoError(t, err)
	data, err := os.ReadFile(output)
	require.NoError(t, err)
	var payload map[string]any
	require.NoError(t, json.Unmarshal(data, &payload))

	assert.Equal(t, "dynatrace-oss/terraform-provider-dynatrace", payload["github.repository"])
	assert.Equal(t, "123456", payload["github.run_id"])
	assert.Equal(t, "casc.e2e.test", payload["event.type"])
	assert.JSONEq(t,
		`[{"dt.settings.schema_id":"builtin:audit-log","dt.settings.schema_version":"1","status":"PASS"}]`,
		payload["results"].(string))
}

func TestWriteResultsAsEvent_UnwritableOutput(t *testing.T) {
	err := writeResultsAsEvent(nil, filepath.Join(t.TempDir(), "missing", "test-report.json"))

	require.Error(t, err)
	assert.Contains(t, err.Error(), "error writing event to")
}
