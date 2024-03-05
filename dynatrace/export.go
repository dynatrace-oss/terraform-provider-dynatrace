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

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/export"
)

func Export(args []string) bool {
	if len(args) == 1 {
		return false
	}
	if strings.TrimSpace(args[1]) != "-export" {
		return false
	}

	// defer export.CleanUp.Finish()
	if err := runExport(); err != nil {
		fmt.Println(err.Error())
	}
	return true
}

func runExport() (err error) {
	os.Remove("terraform-provider-dynatrace.export.log")
	os.Remove("terraform-provider-dynatrace.warnings.log")
	var environment *export.Environment
	if environment, err = export.Initialize(); err != nil {
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
