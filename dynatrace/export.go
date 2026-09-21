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
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/export"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/provider/config"
	"golang.org/x/exp/slices"
)

// IsExport return true if the arguments start with "-export"
func IsExport(args []string) bool {
	return len(args) > 1 && args[1] == "-export"
}

func Export(args []string, cfgGetter config.Getter) error {
	if slices.ContainsFunc(args, func(arg string) bool { return strings.TrimSpace(arg) == "-list-exclusions" }) {
		if len(args) > 3 {
			fmt.Println("-list-exclusions cannot be combined with other flags\nUsage: terraform-provider-dynatrace -export -list-exclusions")
			return nil
		}

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
		return nil
	}
	// defer export.CleanUp.Finish()
	return runExport(cfgGetter, args)
}

func runExport(cfgGetter config.Getter, args []string) (err error) {
	start := time.Now()
	defer func() {
		fmt.Printf("... finished after %v seconds\n", int64(time.Since(start).Seconds()))
	}()
	os.Remove("terraform-provider-dynatrace.export.log")
	os.Remove("terraform-provider-dynatrace.warnings.log")

	pc := config.ProviderConfigureGeneric(context.Background(), cfgGetter)

	// Ensure every ordered Settings 2.0 resource (with `insert_after` attribute) won't produce hardcoded IDs when exported.
	export.AddInsertAfterWeakIDDependencies(export.AllResources, pc)

	var environment *export.Environment
	if environment, err = export.Initialize(pc, args); err != nil {
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
