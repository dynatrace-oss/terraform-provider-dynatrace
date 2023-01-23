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
	"os/exec"
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
	var environment *export.Environment
	if environment, err = export.Initialize(); err != nil {
		return err
	}
	if err = environment.Export(); err != nil {
		return err
	}

	exePath, _ := exec.LookPath("terraform.exe")
	cmd := exec.Command(exePath, "init")
	cmd.Dir = environment.OutputFolder
	cmd.Start()
	if err = cmd.Wait(); err != nil {
		return err
	}

	if environment.Flags.ImportState {
		fmt.Println("Importing Resources into Terraform State ...")
		if err = environment.ExecuteImport(); err != nil {
			return err
		}
	}

	return nil
}
