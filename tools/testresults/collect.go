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
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// collectJUnitFiles returns all *.xml files found (recursively) under dir.
func collectJUnitFiles(dir string) ([]string, error) {
	var files []string
	err := filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		if strings.EqualFold(filepath.Ext(path), ".xml") {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	sort.Strings(files)
	return files, nil
}

// parseFile reads and unmarshals a single JUnit XML file into a testSuites.
func parseFile(path string) (testSuites, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return testSuites{}, err
	}
	var suites testSuites
	if err := xml.Unmarshal(data, &suites); err != nil {
		return testSuites{}, fmt.Errorf("parsing %s: %w", path, err)
	}
	return suites, nil
}

// filterFailures returns a testSuites containing only the failing testcases of
// the input suites, with recomputed counts. Suites without failures are dropped.
func filterFailures(inputs []testSuites) testSuites {
	result := testSuites{Name: "failed-e2e-tests"}
	for _, in := range inputs {
		for _, suite := range in.Suites {
			var failed []testCase
			for _, tc := range suite.TestCases {
				if tc.failed() {
					failed = append(failed, tc)
				}
			}
			if len(failed) == 0 {
				continue
			}
			outSuite := testSuite{
				Name:      suite.Name,
				Time:      suite.Time,
				Timestamp: suite.Timestamp,
				TestCases: failed,
			}
			outSuite.Tests = len(failed)
			for _, tc := range failed {
				if tc.Error != nil {
					outSuite.Errors++
				} else {
					outSuite.Failures++
				}
			}
			result.Suites = append(result.Suites, outSuite)
			result.Tests += outSuite.Tests
			result.Failures += outSuite.Failures
			result.Errors += outSuite.Errors
		}
	}
	return result
}

// failedTestNames returns a sorted, de-duplicated list of "package/test"
// identifiers for every failing testcase in the given suites.
func failedTestNames(suites testSuites) []string {
	seen := map[string]struct{}{}
	var names []string
	for _, suite := range suites.Suites {
		for _, tc := range suite.TestCases {
			pkg := tc.ClassName
			if pkg == "" {
				pkg = suite.Name
			}
			id := tc.Name
			if pkg != "" {
				id = pkg + "/" + tc.Name
			}
			if _, ok := seen[id]; ok {
				continue
			}
			seen[id] = struct{}{}
			names = append(names, id)
		}
	}
	sort.Strings(names)
	return names
}
