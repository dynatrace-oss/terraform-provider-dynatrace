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

// schemaInfo holds only the two root-level fields we need from schema.json.
type schemaInfo struct {
	SchemaID string `json:"schemaId"`
	Version  string `json:"version"`
}

// collectSchemas walks the given root directory and returns one schemaInfo
// for every schema.json file found within it.
func collectSchemas(root string) ([]schemaInfo, error) {
	var infos []schemaInfo

	err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || d.Name() != "schema.json" {
			return nil
		}

		f, err := os.Open(path)
		if err != nil {
			return fmt.Errorf("open %s: %w", path, err)
		}
		defer f.Close()

		var sf schemaInfo
		if err := json.NewDecoder(f).Decode(&sf); err != nil {
			return fmt.Errorf("decode %s: %w", path, err)
		}

		infos = append(infos, sf)
		return nil
	})

	return infos, err
}
