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

	"github.com/dynatrace-oss/terraform-provider-dynatrace/tools/internal/schema"
)

// readSchemaSupport reads the schemaSupport.json array from path.
func readSchemaSupport(path string) ([]schema.Entry, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read %s: %w", path, err)
	}
	var entries []schema.Entry
	if err := json.Unmarshal(data, &entries); err != nil {
		return nil, fmt.Errorf("parse %s: %w", path, err)
	}
	return entries, nil
}

// toResults converts schema entries to event results, preserving the input order
// (schemaSupport.json is already sorted by resource name then schemaId).
func toResults(entries []schema.Entry) []result {
	results := make([]result, 0, len(entries))
	for _, e := range entries {
		results = append(results, result{
			SchemaID:     e.SchemaID,
			Version:      e.Version,
			ResourceName: e.ResourceName,
		})
	}
	return results
}
