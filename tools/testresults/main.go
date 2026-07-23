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
	"flag"
	"fmt"
	"log/slog"
	"os"
	"strings"
)

// run via `go run ./tools/testresults --input <dir> --out-junit failed-tests.xml --out-list failed-tests.txt`
//
// It reads every JUnit XML file produced by the e2e test jobs (gotestsum
// `--junitfile`), keeps only the failing testcases, and writes:
//   - a merged JUnit XML containing only the failures (--out-junit)
//   - a newline-delimited list of failed "package/test" ids (--out-list)
//
// The tool always exits 0 (even when failures exist); it only reports results.

func main() {
	input := flag.String("input", "", "directory containing JUnit XML result files (searched recursively)")
	outJUnit := flag.String("out-junit", "failed-tests.xml", "path to write the merged failures-only JUnit XML")
	outList := flag.String("out-list", "failed-tests.txt", "path to write the newline-delimited failed test list")
	flag.Parse()

	if *input == "" {
		slog.Error("--input is required")
		os.Exit(2)
	}

	if err := run(*input, *outJUnit, *outList); err != nil {
		slog.Error("failed to collect test results", "error", err)
		os.Exit(1)
	}
}

func run(inputDir, outJUnitPath, outListPath string) error {
	files, err := collectJUnitFiles(inputDir)
	if err != nil {
		return fmt.Errorf("collecting JUnit files: %w", err)
	}
	if len(files) == 0 {
		slog.Warn("no JUnit XML files found", "dir", inputDir)
	}

	var parsed []testSuites
	for _, f := range files {
		suites, err := parseFile(f)
		if err != nil {
			return err
		}
		parsed = append(parsed, suites)
	}

	failures := filterFailures(parsed)
	names := failedTestNames(failures)

	if err := writeJUnit(outJUnitPath, failures); err != nil {
		return fmt.Errorf("writing JUnit output: %w", err)
	}
	if err := writeList(outListPath, names); err != nil {
		return fmt.Errorf("writing list output: %w", err)
	}

	slog.Info("collected failed tests",
		"filesProcessed", len(files),
		"failedTests", len(names),
		"junitOutput", outJUnitPath,
		"listOutput", outListPath)
	return nil
}

func writeJUnit(path string, suites testSuites) error {
	data, err := xmlMarshalIndent(suites)
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

func writeList(path string, names []string) error {
	content := strings.Join(names, "\n")
	if len(names) > 0 {
		content += "\n"
	}
	return os.WriteFile(path, []byte(content), 0644)
}
