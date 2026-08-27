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
	"maps"
	"os"
	"slices"
	"sort"
)

const moduleName = "github.com/dynatrace-oss/terraform-provider-dynatrace"

// run via `go run ./tools/testresults --input <dir> --sourceRoot <dir> --output <file>`
func main() {
	input := flag.String("input", "", "directory containing the JUnit result artifacts to collect")
	sourceRoot := flag.String("sourceRoot", "", "path to the source code, used to convert JUnit test suite names to paths")
	output := flag.String("output", "", "file path the JSON event payload is written to")

	flag.Parse()

	if *input == "" {
		fmt.Fprintln(os.Stderr, "error: --input is required")
		os.Exit(1)
	}

	if *sourceRoot == "" {
		fmt.Fprintln(os.Stderr, "error: --sourceRoot is required")
		os.Exit(1)
	}

	if *output == "" {
		fmt.Fprintln(os.Stderr, "error: --output is required")
		os.Exit(1)
	}

	testSuites, errs := findAllTestSuites(*input)
	printErrorsAsWarnings(errs)

	results, errs := collectResults(testSuites, *sourceRoot)
	printErrorsAsWarnings(errs)

	if err := writeResultsAsEvent(results, *output); err != nil {
		fmt.Fprintf(os.Stderr, "error writing results as event: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("wrote event to %s\n", *output)
}

func printErrorsAsWarnings(errs []error) {
	for _, err := range errs {
		fmt.Fprintln(os.Stderr, "warning:", err)
	}
}

func collectResults(testSuites []testSuite, sourceRoot string) ([]result, []error) {
	resultsBySchema := map[string]result{}
	errs := []error{}
	for _, suite := range testSuites {
		// ignore suites with no tests
		if suite.Tests == 0 {
			continue
		}

		suitePath, err := nameToPath(suite.Name, sourceRoot, moduleName)
		if err != nil {
			errs = append(errs, fmt.Errorf("failed to convert suite name to path: %w", err))
			continue
		}

		schema, err := findSettingsSchemaIDAndVersion(suitePath)
		if err != nil {
			errs = append(errs, fmt.Errorf("failed to find settings schema in %s for %s: %w", suitePath, suite.Name, err))
			continue
		}

		r, found := resultsBySchema[schema.SchemaID]
		if !found {
			r = result{
				SchemaID: schema.SchemaID,
				Version:  schema.Version,
				Status:   "PASS",
			}
		}

		// if any suite failed, mark the schema as failed
		if suite.Failures > 0 {
			r.Status = "FAIL"
		}

		resultsBySchema[schema.SchemaID] = r
	}
	results := slices.Collect(maps.Values(resultsBySchema))
	sort.Slice(results, func(i, j int) bool { return results[i].SchemaID < results[j].SchemaID })
	return results, errs
}

func writeResultsAsEvent(results []result, output string) error {
	ee := getEnvironmentEvent()
	event, err := buildEvent(ee, results)
	if err != nil {
		return fmt.Errorf("error building event: %w", err)
	}

	if err := writeEvent(output, event); err != nil {
		return fmt.Errorf("error writing event to %s: %w", output, err)
	}

	return nil
}
