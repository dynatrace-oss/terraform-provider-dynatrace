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
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

var reSchema = regexp.MustCompile(`\(schema:\s*([^)]+)\)`)

// parseResourceMapping reads a supported-resources markdown file and returns
// a map of schemaId -> []resourceName for every table row that references a
// Settings API schema.
func parseResourceMapping(path string) (map[string][]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open %s: %w", path, err)
	}
	defer f.Close()

	result := map[string][]string{}

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()

		// Only process pipe-delimited data rows (skip headers, separators, blank lines).
		if !strings.HasPrefix(line, "|") {
			continue
		}
		// Separator rows look like |---|---|
		if strings.Contains(line, "---") {
			continue
		}

		cols := strings.Split(line, "|")
		// cols[0] is empty (before first |), cols[1] = resource name, cols[2] = API endpoint
		if len(cols) < 3 {
			continue
		}

		resourceName := strings.TrimSpace(cols[1])
		endpoint := strings.TrimSpace(cols[2])

		if resourceName == "" || strings.EqualFold(resourceName, "resource name") {
			continue
		}

		matches := reSchema.FindStringSubmatch(endpoint)
		if matches == nil {
			// No schema in this row — resource uses a different API; skip.
			continue
		}
		schemaID := strings.TrimSpace(matches[1])
		result[schemaID] = append(result[schemaID], resourceName)
	}

	return result, scanner.Err()
}
