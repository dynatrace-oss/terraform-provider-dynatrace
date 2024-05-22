/**
* @license
* Copyright 2023 Dynatrace LLC
*
* Licensed under the Apache License, Version 2.0 (the "License");
* you may not use this file except in compliance with the License.
* You may obtain a copy of the License at
*
*     http://www.apache.org/licenses/LICENSE-2.0
*
* Unless required by applicable law or agreed to in writing, software
* distributed under the License is distributed on an "AS IS" BASIS,
* WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
* See the License for the specific language governing permissions and
* limitations under the License.
 */
package export

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"

	"github.com/spf13/afero"
)

var EXPORT_IGNORE_RESOURCES_PATH = os.Getenv("DYNATRACE_EXPORT_IGNORE_RESOURCES")
var ANSI_ESCAPE_PATTERN = regexp.MustCompile(`\x1B(?:[@-Z\\-_]|\[[0-?]*[ -/]*[@-~])`)

type IgnoreResourcesMap map[string]map[string][]string

var IGNORE_RESOURCES_MAP *IgnoreResourcesMap = nil

func LoadIgnoreResourcesMap() error {
	if EXPORT_IGNORE_RESOURCES_PATH == "" {
		return nil
	}

	if IGNORE_RESOURCES_MAP == nil {
		// pass
	} else {
		return nil
	}

	fs := afero.NewOsFs()

	IGNORE_RESOURCES_MAP = &IgnoreResourcesMap{}

	data, err := afero.ReadFile(fs, EXPORT_IGNORE_RESOURCES_PATH)
	if err != nil {
		fmt.Printf("\nIgnore Resource file not found: %s", EXPORT_IGNORE_RESOURCES_PATH)
		return nil
	}

	err = json.Unmarshal(data, IGNORE_RESOURCES_MAP)
	if err != nil {
		fmt.Printf("\nThe Ignore Resource file wasn't formatted properly: %s, error: %s", EXPORT_IGNORE_RESOURCES_PATH, err)
		return nil
	}

	fmt.Printf("\nThe Ignore Resource file was read properly: %s", EXPORT_IGNORE_RESOURCES_PATH)
	for module_name, resourceMap := range *IGNORE_RESOURCES_MAP {

		fmt.Printf("\n- %s", module_name)

		for resource_name, error_lines := range resourceMap {
			fmt.Printf("\n  - %s", resource_name)

			for _, error_line := range error_lines {

				fmt.Printf("\n    %s", removeANSIColors(error_line))
			}
		}
	}
	fmt.Println("\n")

	return nil
}

func removeANSIColors(logContent string) string {
	cleanedContent := ANSI_ESCAPE_PATTERN.ReplaceAllString(logContent, "")

	return cleanedContent
}

func IsIgnoredResource(resourceType ResourceType, id string) bool {
	if IGNORE_RESOURCES_MAP == nil {
		return false
	}

	resourceMap, exists := (*IGNORE_RESOURCES_MAP)[resourceType.Trim()]

	if exists {
		// pass
	} else {
		return false
	}

	_, exists = resourceMap[id]

	return exists
}
