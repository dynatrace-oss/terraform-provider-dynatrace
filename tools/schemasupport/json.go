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
	"sort"
	"strings"
)

type schemaEntry struct {
	ResourceName string `json:"resourceName"`
	SchemaID     string `json:"schemaId"`
	Version      string `json:"version"`
}

// expandRows joins the schema infos with the resource mapping and returns one
// entry per resource name per schema.
// It also returns the schemaIds that had no matching resource name.
func expandRows(infos []schemaInfo, schemaToResources map[string][]string) (entries []schemaEntry, missing []string) {
	for _, info := range infos {
		resources, found := schemaToResources[info.SchemaID]
		if !found || len(resources) == 0 {
			missing = append(missing, info.SchemaID)
			continue
		}
		for _, res := range resources {
			entries = append(entries, schemaEntry{ResourceName: res, SchemaID: info.SchemaID, Version: info.Version})
		}
	}
	return entries, missing
}

// sortRows sorts entries first by resource name, then by schemaId.
func sortRows(entries []schemaEntry) {
	sort.Slice(entries, func(i, j int) bool {
		if entries[i].ResourceName != entries[j].ResourceName {
			return entries[i].ResourceName < entries[j].ResourceName
		}
		return entries[i].SchemaID < entries[j].SchemaID
	})
}

// writeJSON writes entries as a JSON array of {resourceName, schemaId, version} objects to path.
func writeJSON(path string, entries []schemaEntry) error {
	data, err := json.MarshalIndent(entries, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal JSON: %w", err)
	}
	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("write %s: %w", path, err)
	}
	return nil
}

// formatMissingError returns an error message listing all schemaIds that have
// no matching resource name.
func formatMissingError(missing []string) string {
	sort.Strings(missing)
	var sb strings.Builder
	sb.WriteString("the following schemaIds have no matching resource name in supported-resources.md:")
	for _, id := range missing {
		sb.WriteString("\n  - ")
		sb.WriteString(id)
	}
	return sb.String()
}
