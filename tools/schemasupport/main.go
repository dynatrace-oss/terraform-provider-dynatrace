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
	"fmt"
	"os"
	"path/filepath"
)

const (
	schemasDir     = "dynatrace"
	supportedResMd = "documentation/supported-resources.md"
	outputJSON     = "schemaSupport.json"
)

// run via `go run ./tools/schemasupport/`

func main() {
	infos, err := collectSchemas(schemasDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error collecting schemas: %v\n", err)
		os.Exit(1)
	}

	schemaToResources, err := parseResourceMapping(filepath.FromSlash(supportedResMd))
	if err != nil {
		fmt.Fprintf(os.Stderr, "error parsing supported-resources.md: %v\n", err)
		os.Exit(1)
	}

	rows, missing := expandRows(infos, schemaToResources)
	if len(missing) > 0 {
		fmt.Fprintf(os.Stderr, "error: %s\n", formatMissingError(missing))
		os.Exit(1)
	}

	sortRows(rows)

	if err := writeJSON(outputJSON, rows); err != nil {
		fmt.Fprintf(os.Stderr, "error writing JSON: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("wrote %s (%d rows)\n", outputJSON, len(rows))
}
