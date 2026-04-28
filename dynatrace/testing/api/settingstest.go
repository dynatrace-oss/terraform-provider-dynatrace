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
	"github.com/dynatrace-oss/terraform-provider-dynatrace/provider/envutils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/require"
)

type TestCaseAccOptions struct {
	ExternalProviders map[string]resource.ExternalProvider
}
type TestAccOptions struct {
	ExpectNonEmptyPlan bool
	ExternalProviders  map[string]resource.ExternalProvider
}

func AccEnvsGiven(t *testing.T) bool {
	t.Helper()
	if v := envutils.TFAcc.Get(); v == "" {
		t.Skip("TF_ACC has not been set for acceptance tests")
		return false
	}
	if v := envutils.DynatraceEnvURL.Get(); v == "" {
		t.Skip("DYNATRACE_ENV_URL has not been set for acceptance tests")
		return false
	}
	return true
}

// replaceWithIdentifier replaces "#name#" and "${randomize}" with a given identifier
func replaceWithIdentifier(config string, identifier string) string {
	config = strings.ReplaceAll(config, "#name#", identifier)
	return strings.ReplaceAll(config, "${randomize}", identifier)
}

// ReadTfConfig reads a config and replaces "#name#" and "${randomize}" with a random string
// Returns:
//   - The replaced config
//   - The random string that was used
func ReadTfConfig(t *testing.T, file string) (config string, identifier string) {
	identifier = acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	return ReadTfConfigWithIdentifier(t, file, identifier), identifier
}

// ReadTfConfigWithIdentifier reads a config and replaces "#name#" and "${randomize}" with a given identifier
func ReadTfConfigWithIdentifier(t *testing.T, file string, identifier string) string {
	t.Helper()
	content, err := os.ReadFile(file)
	require.NoError(t, err)

	return replaceWithIdentifier(string(content), identifier)
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

// TestAcc executes all test files in the `testdata` folder ("testdata/*/*.tf").
// Fail conditions could be
// - API errors during create and cleanup (GET, POST, DELETE)
// - Inconsistencies (GET after apply)
// - Non-empty plan after create
func TestAcc(t *testing.T, opts ...TestAccOptions) {
	t.Helper()

	if !AccEnvsGiven(t) {
		return
	}

	allFiles := readTestData(t)

	for _, file := range allFiles {
		subTestName := strings.TrimPrefix(file, "testdata/")

		t.Run(subTestName, func(t *testing.T) {
			t.Helper()

			config, _ := ReadTfConfig(t, file)

			testCase := createTestCaseWithOptions(t, config, opts)
			resource.Test(t, testCase)
		})
	}
}

// TestAccSingle executes a single test/example file.  e.g., "testdata/terraform/example.tf"
// Fail conditions could be
// - API errors during create and cleanup (GET, POST, DELETE)
// - Inconsistencies (GET after apply)
// - Non-empty plan after create
func TestAccSingle(t *testing.T, file string, opts ...TestAccOptions) {
	t.Helper()

	if !AccEnvsGiven(t) {
		return
	}

	config, _ := ReadTfConfig(t, file)

	testCase := createTestCaseWithOptions(t, config, opts)
	resource.Test(t, testCase)
}

type testCaseExecution struct {
	TestCases []resource.TestStep
	Name      string
}

func getTestCase(t *testing.T, folder string, testName string) testCaseExecution {
	configCreate, identifier := ReadTfConfig(t, path.Join(folder, "create.tf"))
	configUpdate := ReadTfConfigWithIdentifier(t, path.Join(folder, "update.tf"), identifier)

	return testCaseExecution{
		TestCases: []resource.TestStep{
			{
				Config: configCreate,
			},
			{
				Config: configUpdate,
			},
		},
		Name: testName,
	}
}

func getTestCases(t *testing.T) []testCaseExecution {
	testDataFolders, _ := os.ReadDir("testcases")
	testCases := make([]testCaseExecution, 0)

	for _, entry := range testDataFolders {
		if !entry.IsDir() {
			continue
		}
		folder := path.Join("testcases", entry.Name())

		testCases = append(testCases, getTestCase(t, folder, entry.Name()))
	}
	return testCases
}

// TestAccTestCases reads the `testcases` folder and executes every create.tf and update.tf file per testcase folder inside
// Fail conditions could be
// - API errors during create, update, and cleanup (GET, POST, PUT, DELETE)
// - Inconsistencies (GET after apply)
// - Non-empty plan after create and update
func TestAccTestCases(t *testing.T, opts ...TestCaseAccOptions) {
	t.Helper()

	if !AccEnvsGiven(t) {
		return
	}
	var options TestCaseAccOptions
	if len(opts) > 0 {
		options = opts[0]
	}

	for _, testCase := range getTestCases(t) {
		t.Run(testCase.Name, func(t *testing.T) {
			t.Helper()

			providerFactories := map[string]func() (*schema.Provider, error){
				"dynatrace": func() (*schema.Provider, error) {
					return provider.Provider(), nil
				},
			}

			resource.Test(t, resource.TestCase{
				ProviderFactories: providerFactories,
				ExternalProviders: options.ExternalProviders,
				Steps:             testCase.TestCases,
			})
		})
	}
}

// TestAccTestCase executes a single test case. e.g., "testcases/update-sets" (containing create.tf and update.tf)
// useful different testcase setups or debugging
// Fail conditions could be
// - API errors during create, update, and cleanup (GET, POST, PUT, DELETE)
// - Inconsistencies (GET after apply)
// - Non-empty plan after create and update
func TestAccTestCase(t *testing.T, folder string, opts ...TestCaseAccOptions) {
	t.Helper()
	if !AccEnvsGiven(t) {
		return
	}
	var options TestCaseAccOptions
	if len(opts) > 0 {
		options = opts[0]
	}

	resource.Test(t, resource.TestCase{
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"dynatrace": func() (*schema.Provider, error) {
				return provider.Provider(), nil
			},
		},
		ExternalProviders: options.ExternalProviders,
		Steps:             getTestCase(t, folder, "").TestCases,
	})
}

// TestAccParallel runs all tests in parallel by concatenating all configs into one and running them together.
// Make sure that the tests do not interfere with each other, e.g. by using unique names for resources.
func TestAccParallel(t *testing.T, opts ...TestAccOptions) {
	t.Helper()

	if !AccEnvsGiven(t) {
		return
	}

	allFiles := readTestData(t)

	// concatenate all configs into one, so that they can be run in parallel
	var configs strings.Builder
	for _, file := range allFiles {
		config, _ := ReadTfConfig(t, file)
		configs.WriteString(config)
		configs.WriteByte('\n')
	}

	testCase := createTestCaseWithOptions(t, configs.String(), opts)
	resource.Test(t, testCase)
}

func createTestCaseWithOptions(t *testing.T, config string, opts []TestAccOptions) resource.TestCase {
	t.Helper()
	providerFactories := map[string]func() (*schema.Provider, error){
		"dynatrace": func() (*schema.Provider, error) {
			return provider.Provider(), nil
		},
	}
	var options TestAccOptions
	if len(opts) > 0 {
		options = opts[0]
	}

	return resource.TestCase{
		ProviderFactories: providerFactories,
		ExternalProviders: options.ExternalProviders,
		Steps:             []resource.TestStep{{Config: config, ExpectNonEmptyPlan: options.ExpectNonEmptyPlan}},
	}
}
