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
	"os"
)

// run via `go run ./tools/testresults --input <dir> --output <file>`
//
// This is a mock: it ignores --input and writes a fixed set of CASC_E2E_TEST
// events (populated with the ambient GitHub run context from the environment)
// to --output. Parsing real JUnit results from --input is a follow-up story.

var mockTestNames = []string{
	"github.com/dynatrace-oss/terraform-provider-dynatrace/mock/TestMockOne",
	"github.com/dynatrace-oss/terraform-provider-dynatrace/mock/TestMockTwo",
}

func main() {
	input := flag.String("input", "", "directory containing the JUnit result artifacts to collect")
	output := flag.String("output", "", "file path the JSON event payload is written to")
	flag.Parse()

	if *output == "" {
		fmt.Fprintln(os.Stderr, "error: --output is required")
		os.Exit(1)
	}

	events := buildEvents(mockTestNames, metadataFromEnvironment())

	if err := writeEvents(*output, events); err != nil {
		fmt.Fprintf(os.Stderr, "error writing events to %s: %v\n", *output, err)
		os.Exit(1)
	}

	fmt.Printf("wrote %d event(s) to %s (input %q ignored by mock)\n", len(events), *output, *input)
}
