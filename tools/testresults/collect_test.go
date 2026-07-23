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
	"reflect"
	"testing"
)

const sampleJUnitA = `<?xml version="1.0" encoding="UTF-8"?>
<testsuites tests="3" failures="1" errors="0">
  <testsuite name="pkg/one" tests="3" failures="1" errors="0">
    <testcase classname="pkg/one" name="TestPass" time="0.01"></testcase>
    <testcase classname="pkg/one" name="TestFail" time="0.02">
      <failure message="boom">assertion failed</failure>
    </testcase>
    <testcase classname="pkg/one" name="TestSkip" time="0.00">
      <skipped message="skip"></skipped>
    </testcase>
  </testsuite>
</testsuites>`

const sampleJUnitB = `<?xml version="1.0" encoding="UTF-8"?>
<testsuites tests="2" failures="0" errors="1">
  <testsuite name="pkg/two" tests="2" failures="0" errors="1">
    <testcase classname="pkg/two" name="TestErr" time="0.03">
      <error message="panic">runtime error</error>
    </testcase>
    <testcase classname="pkg/two" name="TestOK" time="0.01"></testcase>
  </testsuite>
</testsuites>`

const sampleJUnitAllPass = `<?xml version="1.0" encoding="UTF-8"?>
<testsuites tests="1" failures="0" errors="0">
  <testsuite name="pkg/three" tests="1" failures="0" errors="0">
    <testcase classname="pkg/three" name="TestGood" time="0.01"></testcase>
  </testsuite>
</testsuites>`

func parseString(t *testing.T, dir, name, content string) testSuites {
	t.Helper()
	path := filepath.Join(dir, name)
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("writing %s: %v", name, err)
	}
	suites, err := parseFile(path)
	if err != nil {
		t.Fatalf("parsing %s: %v", name, err)
	}
	return suites
}

func TestParseFile(t *testing.T) {
	dir := t.TempDir()
	suites := parseString(t, dir, "a.xml", sampleJUnitA)
	if len(suites.Suites) != 1 {
		t.Fatalf("expected 1 suite, got %d", len(suites.Suites))
	}
	if got := len(suites.Suites[0].TestCases); got != 3 {
		t.Fatalf("expected 3 testcases, got %d", got)
	}
	if !suites.Suites[0].TestCases[1].failed() {
		t.Errorf("TestFail should be marked failed")
	}
	if suites.Suites[0].TestCases[0].failed() {
		t.Errorf("TestPass should not be marked failed")
	}
}

func TestCollectJUnitFiles(t *testing.T) {
	dir := t.TempDir()
	sub := filepath.Join(dir, "nested")
	if err := os.MkdirAll(sub, 0755); err != nil {
		t.Fatal(err)
	}
	mustWrite(t, filepath.Join(dir, "a.xml"), sampleJUnitA)
	mustWrite(t, filepath.Join(sub, "b.xml"), sampleJUnitB)
	mustWrite(t, filepath.Join(dir, "notes.txt"), "ignore me")

	files, err := collectJUnitFiles(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(files) != 2 {
		t.Fatalf("expected 2 xml files, got %d: %v", len(files), files)
	}
}

func TestFilterFailures(t *testing.T) {
	dir := t.TempDir()
	a := parseString(t, dir, "a.xml", sampleJUnitA)
	b := parseString(t, dir, "b.xml", sampleJUnitB)
	pass := parseString(t, dir, "c.xml", sampleJUnitAllPass)

	result := filterFailures([]testSuites{a, b, pass})

	if result.Tests != 2 {
		t.Errorf("expected 2 failing tests total, got %d", result.Tests)
	}
	if result.Failures != 1 {
		t.Errorf("expected 1 failure, got %d", result.Failures)
	}
	if result.Errors != 1 {
		t.Errorf("expected 1 error, got %d", result.Errors)
	}
	if len(result.Suites) != 2 {
		t.Fatalf("expected 2 suites (all-pass dropped), got %d", len(result.Suites))
	}
	for _, s := range result.Suites {
		for _, tc := range s.TestCases {
			if !tc.failed() {
				t.Errorf("suite %s contains non-failing testcase %s", s.Name, tc.Name)
			}
		}
	}
}

func TestFailedTestNames(t *testing.T) {
	dir := t.TempDir()
	a := parseString(t, dir, "a.xml", sampleJUnitA)
	b := parseString(t, dir, "b.xml", sampleJUnitB)

	result := filterFailures([]testSuites{a, b})
	names := failedTestNames(result)

	want := []string{"pkg/one/TestFail", "pkg/two/TestErr"}
	if !reflect.DeepEqual(names, want) {
		t.Errorf("expected %v, got %v", want, names)
	}
}

func TestFailedTestNamesDeduplicated(t *testing.T) {
	dir := t.TempDir()
	a := parseString(t, dir, "a.xml", sampleJUnitA)
	aDup := parseString(t, dir, "a2.xml", sampleJUnitA)

	result := filterFailures([]testSuites{a, aDup})
	names := failedTestNames(result)

	want := []string{"pkg/one/TestFail"}
	if !reflect.DeepEqual(names, want) {
		t.Errorf("expected deduplicated %v, got %v", want, names)
	}
}

func TestRunProducesOutputs(t *testing.T) {
	dir := t.TempDir()
	mustWrite(t, filepath.Join(dir, "a.xml"), sampleJUnitA)
	mustWrite(t, filepath.Join(dir, "b.xml"), sampleJUnitB)

	outJUnit := filepath.Join(dir, "failed.xml")
	outList := filepath.Join(dir, "failed.txt")

	if err := run(dir, outJUnit, outList); err != nil {
		t.Fatalf("run failed: %v", err)
	}

	// Re-parse the produced JUnit to confirm it only contains failures.
	produced, err := parseFile(outJUnit)
	if err != nil {
		t.Fatalf("parsing produced junit: %v", err)
	}
	if produced.Tests != 2 {
		t.Errorf("expected 2 failing tests in output, got %d", produced.Tests)
	}

	listData, err := os.ReadFile(outList)
	if err != nil {
		t.Fatal(err)
	}
	if string(listData) != "pkg/one/TestFail\npkg/two/TestErr\n" {
		t.Errorf("unexpected list content: %q", string(listData))
	}
}

func mustWrite(t *testing.T, path, content string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("writing %s: %v", path, err)
	}
}
