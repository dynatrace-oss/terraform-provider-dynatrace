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

package main

import (
	"os"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/export"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/provider"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
)

// func tf2json(args []string) bool {
// 	if len(args) == 1 {
// 		return false
// 	}
// 	if strings.TrimSpace(args[1]) != "-tf2json" {
// 		return false
// 	}

// 	flag.Bool("tf2json", true, "")
// 	flag.Parse()
// 	tailArgs := flag.Args()

// 	if len(tailArgs) == 0 {
// 		fmt.Println("Usage: terraform-provider-dynatrace -tf2json <folder>")
// 		return true
// 	}

// 	for _, filePath := range tailArgs {
// 		fileInfo, err := os.Stat(filePath)
// 		if err != nil {
// 			fmt.Println(err)
// 			return true
// 		}
// 		if !fileInfo.IsDir() {
// 			if err := tf2jsonFile(filePath, fileInfo, nil); err != nil {
// 				fmt.Println(err)
// 				return true
// 			}
// 			continue
// 		}

// 		filepath.Walk(filePath, tf2jsonFile)
// 	}

// 	return true
// }

// func tf2jsonFile(childPath string, info os.FileInfo, err error) error {
// 	if err != nil {
// 		return err
// 	}
// 	if info.IsDir() {
// 		return nil
// 	}
// 	if !strings.HasSuffix(info.Name(), ".tf") {
// 		return nil
// 	}
// 	if info.Name() == "data_source.tf" {
// 		return nil
// 	}
// 	if info.Name() == "providers.tf" {
// 		return nil
// 	}
// 	if info.Name() == "main.tf" {
// 		return nil
// 	}
// 	if strings.Contains(childPath, ".flawed") {
// 		return nil
// 	}
// 	jsonPath := strings.TrimSuffix(childPath, info.Name()) + info.Name() + ".json"
// 	var module *hcl2json.Module
// 	if module, err = hcl2json.HCL2Config(childPath); err != nil {
// 		fmt.Println(err)
// 		return nil
// 	}
// 	var data []byte
// 	if data, err = json.MarshalIndent(module, "", "  "); err != nil {
// 		return err
// 	}
// 	os.Remove(jsonPath)
// 	if err = os.WriteFile(jsonPath, data, 0644); err != nil {
// 		return err
// 	}
// 	return nil
// }

func main() {
	// if tf2json(os.Args) {
	// 	return
	// }
	defer export.CleanUp.Finish()

	if dynatrace.Export(os.Args) {
		return
	}

	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: func() *schema.Provider {
			return provider.Provider()
		},
	})
}
