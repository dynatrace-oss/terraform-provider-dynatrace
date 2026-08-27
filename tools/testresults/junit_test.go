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
	"encoding/xml"
	"fmt"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ----- nameToPath -----

func TestNameToPath_NestedPackage(t *testing.T) {
	path, err := nameToPath(suiteName("dynatrace/api/builtin/auditlog"), "base", moduleName)

	require.NoError(t, err)
	assert.Equal(t, filepath.Join("base", "dynatrace", "api", "builtin", "auditlog"), path)
}

func TestNameToPath_ModuleItselfResolvesToBasePath(t *testing.T) {
	path, err := nameToPath(moduleName, "base", moduleName)

	require.NoError(t, err)
	assert.Equal(t, "base", path)
}

func TestNameToPath_ConvertsSlashesToSeparator(t *testing.T) {
	path, err := nameToPath(suiteName("a/b/c"), "base", moduleName)

	require.NoError(t, err)
	sep := string(filepath.Separator)
	assert.Equal(t, "base"+sep+"a"+sep+"b"+sep+"c", path)
}

func TestNameToPath_MissingPrefix(t *testing.T) {
	_, err := nameToPath("example.com/other/pkg", "base", moduleName)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "example.com/other/pkg")
	assert.Contains(t, err.Error(), moduleName)
}

func TestNameToPath_AcceptsPrefixWithoutSeparatorBoundary(t *testing.T) {
	path, err := nameToPath(moduleName+"-other/pkg", "base", moduleName)

	require.NoError(t, err)
	assert.Equal(t, filepath.Join("base", "-other", "pkg"), path)
}

// ----- parseFile -----

func TestParseFile_WellFormed(t *testing.T) {
	content := fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<testsuites tests="5" failures="1" errors="0" time="1.500000">
  <testsuite name="%s" tests="3" failures="0" time="0.500000"></testsuite>
  <testsuite name="%s" tests="2" failures="1" time="1.000000"></testsuite>
</testsuites>`, suiteName("a"), suiteName("b"))
	path := writeTempFile(t, t.TempDir(), "test-results.xml", content)

	suites, err := parseFile(path)

	require.NoError(t, err)
	require.Len(t, suites.Suites, 2)
	assert.Equal(t, suiteName("a"), suites.Suites[0].Name)
	assert.Equal(t, 3, suites.Suites[0].Tests)
	assert.Equal(t, 0, suites.Suites[0].Failures)
	assert.Equal(t, suiteName("b"), suites.Suites[1].Name)
	assert.Equal(t, 2, suites.Suites[1].Tests)
	assert.Equal(t, 1, suites.Suites[1].Failures)
}

func TestParseFile_EmptyTestSuites(t *testing.T) {
	path := writeTempFile(t, t.TempDir(), "test-results.xml", `<testsuites></testsuites>`)

	suites, err := parseFile(path)

	require.NoError(t, err)
	assert.Empty(t, suites.Suites)
}

func TestParseFile_MalformedXML(t *testing.T) {
	path := writeTempFile(t, t.TempDir(), "test-results.xml", `<testsuites><testsuite`)

	_, err := parseFile(path)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "error parsing junit file")
}

func TestParseFile_MissingFile(t *testing.T) {
	_, err := parseFile(filepath.Join(t.TempDir(), "missing.xml"))

	require.Error(t, err)
	assert.Contains(t, err.Error(), "error reading junit file")
}

func TestParseFile_TestMainFailureIsInvisible(t *testing.T) {
	content := fmt.Sprintf(`<testsuites tests="0" failures="0" errors="1">
  <testsuite name="%s" tests="0" failures="0" time="0.420000">
    <testcase classname="%s" name="TestMain" time="0.000000">
      <failure message="Failed">panic: boom</failure>
    </testcase>
  </testsuite>
</testsuites>`, suiteName("a"), suiteName("a"))
	path := writeTempFile(t, t.TempDir(), "test-results.xml", content)

	suites, err := parseFile(path)

	require.NoError(t, err)
	require.Len(t, suites.Suites, 1)
	assert.Equal(t, 0, suites.Suites[0].Tests)
	assert.Equal(t, 0, suites.Suites[0].Failures)
}

func TestParseFile_SkippedCountIsInvisible(t *testing.T) {
	content := fmt.Sprintf(
		`<testsuites><testsuite name="%s" tests="3" failures="0" skipped="3"></testsuite></testsuites>`,
		suiteName("a"))
	path := writeTempFile(t, t.TempDir(), "test-results.xml", content)

	suites, err := parseFile(path)

	require.NoError(t, err)
	require.Len(t, suites.Suites, 1)
	assert.Equal(t, 3, suites.Suites[0].Tests)
	assert.Equal(t, 0, suites.Suites[0].Failures)

	remarshalled, err := xml.Marshal(suites)
	require.NoError(t, err)
	assert.NotContains(t, string(remarshalled), "skipped")
}

func TestParseFile_RootErrorsAttributeIsDropped(t *testing.T) {
	content := fmt.Sprintf(
		`<testsuites tests="1" failures="0" errors="2"><testsuite name="%s" tests="1" failures="0"></testsuite></testsuites>`,
		suiteName("a"))
	path := writeTempFile(t, t.TempDir(), "test-results.xml", content)

	suites, err := parseFile(path)

	require.NoError(t, err)

	remarshalled, err := xml.Marshal(suites)
	require.NoError(t, err)
	assert.NotContains(t, string(remarshalled), "errors")
}

// ----- findXMLFiles -----

func TestFindXMLFiles_RecursesAndFiltersByExtension(t *testing.T) {
	dir := t.TempDir()
	writeNestedFile(t, dir, "unit/test-results.xml", `<testsuites></testsuites>`)
	writeNestedFile(t, dir, "v2/nested/test-results.xml", `<testsuites></testsuites>`)
	writeNestedFile(t, dir, "notes.txt", "ignore me")
	writeNestedFile(t, dir, "report.json", "{}")

	files, err := findXMLFiles(dir)

	require.NoError(t, err)
	assert.ElementsMatch(t, []string{
		filepath.Join(dir, "unit", "test-results.xml"),
		filepath.Join(dir, "v2", "nested", "test-results.xml"),
	}, files)
}

func TestFindXMLFiles_EmptyDirectory(t *testing.T) {
	files, err := findXMLFiles(t.TempDir())

	require.NoError(t, err)
	assert.Empty(t, files)
}

func TestFindXMLFiles_NonexistentRoot(t *testing.T) {
	_, err := findXMLFiles(filepath.Join(t.TempDir(), "does-not-exist"))

	assert.Error(t, err)
}

// ----- findAllTestSuites -----

func TestFindAllTestSuites_FlattensAcrossFiles(t *testing.T) {
	dir := t.TempDir()
	writeNestedFile(t, dir, "unit/test-results.xml", fmt.Sprintf(
		`<testsuites><testsuite name="%s" tests="3" failures="0"></testsuite></testsuites>`, suiteName("a")))
	writeNestedFile(t, dir, "v2/test-results.xml", fmt.Sprintf(
		`<testsuites><testsuite name="%s" tests="2" failures="1"></testsuite><testsuite name="%s" tests="1" failures="0"></testsuite></testsuites>`,
		suiteName("b"), suiteName("c")))

	suites, errs := findAllTestSuites(dir)

	assert.Empty(t, errs)
	names := make([]string, 0, len(suites))
	for _, s := range suites {
		names = append(names, s.Name)
	}
	assert.ElementsMatch(t, []string{suiteName("a"), suiteName("b"), suiteName("c")}, names)
}

func TestFindAllTestSuites_MalformedFileIsReportedButOthersSurvive(t *testing.T) {
	dir := t.TempDir()
	writeNestedFile(t, dir, "good/test-results.xml", fmt.Sprintf(
		`<testsuites><testsuite name="%s" tests="3" failures="0"></testsuite></testsuites>`, suiteName("a")))
	writeNestedFile(t, dir, "bad/test-results.xml", `<testsuites><testsuite`)

	suites, errs := findAllTestSuites(dir)

	require.Len(t, errs, 1)
	assert.Contains(t, errs[0].Error(), "error parsing junit file")
	require.Len(t, suites, 1)
	assert.Equal(t, suiteName("a"), suites[0].Name)
}

func TestFindAllTestSuites_NonexistentRoot(t *testing.T) {
	suites, errs := findAllTestSuites(filepath.Join(t.TempDir(), "does-not-exist"))

	require.Len(t, errs, 1)
	assert.Empty(t, suites)
}
