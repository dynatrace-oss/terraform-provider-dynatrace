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
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/address"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/export"
)

const ENV_VAR_CUSTOM_PROVIDER_LOCATION = "DYNATRACE_CUSTOM_PROVIDER_LOCATION"

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

func readStuff(scanner *bufio.Scanner) {
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}
}

func runExport() (err error) {
	var environment *export.Environment
	if environment, err = export.Initialize(); err != nil {
		return err
	}
	if err = environment.Export(); err != nil {
		return err
	}

	exePath, _ := exec.LookPath("terraform")
	cmdOptions := []string{"init", "-no-color"}

	customProviderLocation := os.Getenv(ENV_VAR_CUSTOM_PROVIDER_LOCATION)
	if len(customProviderLocation) != 0 && customProviderLocation != "" {
		cmdOptions = append(cmdOptions, fmt.Sprint("-plugin-dir=", customProviderLocation))
	}

	cmd := exec.Command(exePath, cmdOptions...)
	cmd.Dir = environment.OutputFolder
	outs, err := cmd.StderrPipe()
	if err != nil {
		return err
	}
	err = cmd.Start()
	if err != nil {
		fmt.Println("Terraform CLI not installed - skipping import")
		return nil
	} else {
		fmt.Println("Executing 'terraform init'")
		defer func() {
			cmd.Wait()
			if environment.Flags.ImportState || environment.Flags.ImportStateV2 {
				fmt.Println("Importing Resources into Terraform State ...")
				if err = environment.ExecuteImport(); err != nil {
					return
				}
			}
		}()

		go readStuff(bufio.NewScanner(outs))
	}

	address.SaveOriginalMap(environment.OutputFolder)
	address.SaveCompletedMap(environment.OutputFolder)

	return nil
}
