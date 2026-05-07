/**
* @license
* Copyright 2026 Dynatrace LLC
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
	"regexp"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/provider/envutils"
	"github.com/spf13/afero"
)

var ANSI_ESCAPE_PATTERN = regexp.MustCompile(`\x1B(?:[@-Z\\-_]|\[[0-?]*[ -/]*[@-~])`)

type IgnoreResourcesMap map[string]map[string][]string

var IGNORE_RESOURCES_MAP *IgnoreResourcesMap = nil

func LoadIgnoreResourcesMap() error {
	if envutils.DynatraceExportIgnoreResources.Get() == "" {
		return nil
	}

	if IGNORE_RESOURCES_MAP == nil {
		// pass
	} else {
		return nil
	}

	fs := afero.NewOsFs()

	IGNORE_RESOURCES_MAP = &IgnoreResourcesMap{}

	data, err := afero.ReadFile(fs, envutils.DynatraceExportIgnoreResources.Get())
	if err != nil {
		fmt.Printf("\nIgnore Resource file not found: %s", envutils.DynatraceExportIgnoreResources.Get())
		return nil
	}

	err = json.Unmarshal(data, IGNORE_RESOURCES_MAP)
	if err != nil {
		fmt.Printf("\nThe Ignore Resource file wasn't formatted properly: %s, error: %s", envutils.DynatraceExportIgnoreResources.Get(), err)
		return nil
	}

	fmt.Printf("\nThe Ignore Resource file was read properly: %s", envutils.DynatraceExportIgnoreResources.Get())
	for module_name, resourceMap := range *IGNORE_RESOURCES_MAP {

		fmt.Printf("\n- %s", module_name)

		for resource_name, error_lines := range resourceMap {
			fmt.Printf("\n  - %s", resource_name)

			for _, error_line := range error_lines {

				fmt.Printf("\n    %s", removeANSIColors(error_line))
			}
		}
	}
	fmt.Println()
	fmt.Println()

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
