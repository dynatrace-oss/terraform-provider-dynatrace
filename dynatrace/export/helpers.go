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

package export

import (
	"os"
	"os/exec"
	"strings"
	"unicode"
)

func format(name string, force bool) {
	if force || os.Getenv("DYNATRACE_FORMAT_HCL_FILES") == "true" {
		exePath, _ := exec.LookPath("terraform.exe")
		// log.Println(exePath, "fmt", name)
		cmd := exec.Command(exePath, "fmt", name)
		cmd.Start()
		cmd.Wait()
	}
}

func fileSystemName(filename string) string {
	filename = strings.ReplaceAll(filename, "/", "_")
	filename = strings.ReplaceAll(filename, ":", "_")
	filename = strings.ReplaceAll(filename, "|", "_")
	filename = strings.ReplaceAll(filename, "<", "_")
	filename = strings.ReplaceAll(filename, ">", "_")
	filename = strings.ReplaceAll(filename, "\"", "_")
	filename = strings.ReplaceAll(filename, "?", "_")
	filename = strings.ReplaceAll(filename, "*", "_")
	return filename
}

func toTerraformName(s string) string {
	result := ""

	// s = strings.ReplaceAll(s, "@", "_at_")
	s = strings.ReplaceAll(s, "@", "_")
	// s = strings.ReplaceAll(s, "/", "_slash_")
	s = strings.ReplaceAll(s, "/", "_")
	// s = strings.ReplaceAll(s, "(", "_lrb_")
	s = strings.ReplaceAll(s, "(", "_")
	// s = strings.ReplaceAll(s, ")", "_rrb_")
	s = strings.ReplaceAll(s, ")", "_")
	// s = strings.ReplaceAll(s, "{", "_lcb_")
	s = strings.ReplaceAll(s, "{", "_")
	// s = strings.ReplaceAll(s, "}", "_rcb_")
	s = strings.ReplaceAll(s, "}", "_")
	// s = strings.ReplaceAll(s, "[", "_ldb_")
	s = strings.ReplaceAll(s, "[", "_")
	// s = strings.ReplaceAll(s, "]", "_rdb_")
	s = strings.ReplaceAll(s, "]", "_")
	// s = strings.ReplaceAll(s, "|", "_P_")
	s = strings.ReplaceAll(s, "|", "_")
	// s = strings.ReplaceAll(s, ":", "_colon_")
	s = strings.ReplaceAll(s, ":", "_")
	// s = strings.ReplaceAll(s, ".", "_dot_")
	s = strings.ReplaceAll(s, ".", "_")
	for _, ch := range s {
		if unicode.IsDigit(ch) {
			result = result + string(ch)
		} else if unicode.IsLetter(ch) {
			result = result + string(ch)
		} else if ch == '_' {
			result = result + string(ch)
		} else if ch == '-' {
			result = result + string(ch)
		} else {
			result = result + "_"
		}
	}
	first := []rune(result)[0]
	if !unicode.IsLetter(first) && first != '_' {
		result = "_" + result
	}
	for strings.Contains(result, "__") {
		result = strings.ReplaceAll(result, "__", "_")
	}
	return strings.TrimSuffix(result, "_")
}
