//go:build integration

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

package api

import (
	"os"
	"path"
	"strings"
	"testing"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/provider"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/require"
)

type TestAccOptions struct {
	ExpectNonEmptyPlan bool
}

func AccEnvsGiven(t *testing.T) bool {
	t.Helper()
	if v := os.Getenv("TF_ACC"); v == "" {
		t.Skip("TF_ACC has not been set for acceptance tests")
		return false
	}
	if v := os.Getenv("DYNATRACE_ENV_URL"); v == "" {
		t.Skip("DYNATRACE_ENV_URL has not been set for acceptance tests")
		return false
	}
	return true
}

// RandomizeResource replaces "#name#" and "${randomize}" with a random string
// Returns:
// 	- The replaced config
//  - The random string that was used
func RandomizeResource(config string) (replacedConfig string, identifier string) {
	name := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	config = strings.ReplaceAll(config, "#name#", name)
	return strings.ReplaceAll(config, "${randomize}", name), name
}

// ReadTfConfig reads a config and replaces "#name#" and "${randomize}" with a random string
// Returns:
// 	- The replaced config
//  - The random string that was used
func ReadTfConfig(t *testing.T, file string) (config string, identifier string) {
	t.Helper()
	content, err := os.ReadFile(file)
	require.NoError(t, err)

	return RandomizeResource(string(content))
}

func readTestData(t *testing.T) []string {
	t.Helper()
	testDataFolders, _ := os.ReadDir("testdata")
	allFiles := make([]string, 0)

	for _, entry := range testDataFolders {
		if !entry.IsDir() {
			continue
		}
		folder := path.Join("testdata", entry.Name())

		entries, err := os.ReadDir(folder)
		require.NoError(t, err)

		for _, entry := range entries {
			if !strings.HasSuffix(entry.Name(), ".tf") || strings.HasSuffix(entry.Name(), "__providers__.tf") {
				continue
			}
			allFiles = append(allFiles, path.Join(folder, entry.Name()))
		}
	}
	return allFiles
}

func TestAcc(t *testing.T, opts ...TestAccOptions) {
	var options TestAccOptions
	if len(opts) > 0 {
		options = opts[0]
	}
	t.Helper()

	if !AccEnvsGiven(t) {
		return
	}

	allFiles := readTestData(t)

	providerFactories := map[string]func() (*schema.Provider, error){
		"dynatrace": func() (*schema.Provider, error) {
			return provider.Provider(), nil
		},
	}

	for _, file := range allFiles {
		subTestName := strings.TrimPrefix(file, "testdata/")

		t.Run(subTestName, func(t *testing.T) {
			t.Helper()

			config, _ := ReadTfConfig(t, file)

			testCase := resource.TestCase{
				ProviderFactories: providerFactories,
				Steps:             []resource.TestStep{{Config: config, ExpectNonEmptyPlan: options.ExpectNonEmptyPlan}},
			}
			resource.Test(t, testCase)
		})
	}
}
