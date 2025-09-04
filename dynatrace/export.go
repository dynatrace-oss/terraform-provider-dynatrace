/**
* @license
* Copyright 2020 Dynatrace LLC
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

package dynatrace

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/export"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/provider/config"
)

func Export(args []string, cfgGetter config.Getter) bool {
	if len(args) == 1 {
		return false
	}

	if strings.TrimSpace(args[1]) != "-export" {
		return false
	}

	if len(args) > 2 {
		if strings.TrimSpace(args[2]) == "-list-exclusions" {
			for _, group := range export.GetExcludeListedResourceGroups() {
				fmt.Println(group.Reason)
				// Calculate the maximum length of the name field
				maxNameLength := 0
				for _, exclusion := range group.Exclusions {
					if len(exclusion.ResourceType) > maxNameLength {
						maxNameLength = len(exclusion.ResourceType)
					}
				}

				for _, exclusion := range group.Exclusions {
					if len(exclusion.Reason) == 0 {
						fmt.Printf("  %s\n", exclusion.ResourceType)
					} else {
						fmt.Printf("  %-*s  ... %s\n", maxNameLength, exclusion.ResourceType, exclusion.Reason)
					}
				}
				fmt.Println()
			}
			return true
		}
	}

	// defer export.CleanUp.Finish()
	if err := runExport(cfgGetter); err != nil {
		fmt.Println(err.Error())
	}
	return true
}

func runExport(cfgGetter config.Getter) (err error) {
	start := time.Now()
	defer func() {
		fmt.Printf("... finished after %v seconds\n", int64(time.Since(start).Seconds()))
	}()
	os.Remove("terraform-provider-dynatrace.export.log")
	os.Remove("terraform-provider-dynatrace.warnings.log")

	// This ensures that every Settings 2.0 resource
	// that offers ordering using the `insert_after`
	// attribute won't produce hardcoded IDs
	// when exported.
	export.AddInsertAfterWeakIDDependencies()

	var environment *export.Environment
	if environment, err = export.Initialize(cfgGetter); err != nil {
		return err
	}

	err = environment.RunQuickInit()
	if err != nil {
		return err
	}

	if err = environment.Export(); err != nil {
		return err
	}

	err = environment.FinishExport()
	if err != nil {
		return err
	}

	return nil
}
