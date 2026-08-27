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
	"path/filepath"
)

type schemaIDAndVersion struct {
	SchemaID string `json:"schemaId"`
	Version  string `json:"version"`
}

func findSettingsSchemaIDAndVersion(path string) (schemaIDAndVersion, error) {
	schemaPath := filepath.Join(path, "schema.json")
	if _, err := os.Stat(schemaPath); err != nil {
		return schemaIDAndVersion{}, fmt.Errorf("schema.json not found in %s: %w", path, err)
	}
	return parseSettingsSchemaFile(schemaPath)
}

func parseSettingsSchemaFile(s string) (schemaIDAndVersion, error) {
	data, err := os.ReadFile(s)
	if err != nil {
		return schemaIDAndVersion{}, fmt.Errorf("error reading settings schema file %s: %w", s, err)
	}
	var sch schemaIDAndVersion
	if err := json.Unmarshal(data, &sch); err != nil {
		return schemaIDAndVersion{}, fmt.Errorf("error parsing settings schema file %s: %w", s, err)
	}
	return sch, nil
}
