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

	"github.com/dynatrace-oss/terraform-provider-dynatrace/tools/internal/githubevent"
)

// run via `go run ./tools/schemaingest/ --file schemaSupport.json --out event.json --release`
func main() {
	file := flag.String("file", "", "path to the schemaSupport.json file to read")
	out := flag.String("out", "", "file path the JSON event payload is written to")
	release := flag.Bool("release", false, "mark the event as originating from a release")

	flag.Parse()

	if *file == "" {
		_, _ = fmt.Fprintln(os.Stderr, "error: --file is required")
		os.Exit(1)
	}

	if *out == "" {
		_, _ = fmt.Fprintln(os.Stderr, "error: --out is required")
		os.Exit(1)
	}

	entries, err := readSchemaSupport(*file)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "error reading schema support file: %v\n", err)
		os.Exit(1)
	}

	e, err := buildEvent(githubevent.Get(), *release, toResults(entries))
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "error building event: %v\n", err)
		os.Exit(1)
	}

	if err := writeEvent(*out, e); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "error writing event to %s: %v\n", *out, err)
		os.Exit(1)
	}

	fmt.Printf("wrote event to %s\n", *out)
}
