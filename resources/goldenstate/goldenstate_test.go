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

package goldenstate_test

import (
	"context"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/export"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/provider"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/provider/config"
	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func validateHCL(hclString string) ([]string, error) {
	parser := hclparse.NewParser()
	hclFile, diag := parser.ParseHCL([]byte(hclString), "test.hcl")
	if diag.HasErrors() {
		return nil, fmt.Errorf("invalid HCL: %s", diag.Error())
	}
	resourceBlocks := []string{}
	hclBody := hclFile.Body.(*hclsyntax.Body)
	for _, block := range hclBody.Blocks {
		if block.Type == "resource" && len(block.Labels) == 2 {
			resourceBlocks = append(resourceBlocks, fmt.Sprintf("%s.%s", block.Labels[0], block.Labels[1]))
		}
	}

	return resourceBlocks, nil
}

type GoldenResource struct {
	ResourceType                export.ResourceType
	AbsPathGoldenConfigFile     string
	AbsPathGarbageConfigFile    string
	AbsPathAdditionalConfigFile string
	GarbageResourceBlocks       []string
	GoldenResourceBlocks        []string
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func findGoldenResources() ([]GoldenResource, error) {
	goldenResources := []GoldenResource{}

	entries, err := os.ReadDir("testdata")
	if err != nil {
		return nil, err
	}
	for _, entry := range entries {
		folderName := entry.Name()
		if !strings.HasPrefix(folderName, "dynatrace_") {
			return nil, err
		}
		if !entry.IsDir() {
			return nil, err
		}
		expectedResourceType := strings.TrimSpace(folderName)
		_ = expectedResourceType
		if fileExists(path.Join("testdata", folderName, "skip")) {
			continue
		}
		files, err := os.ReadDir(path.Join("testdata", folderName))
		if err != nil {
			return nil, err
		}
		absPathGoldenConfigFile := ""
		absPathGarbageConfigFile := ""
		absPathAdditionalConfigFile := ""

		garbageResourceBlocks := []string{}
		goldenResourceBlocks := []string{}
		for _, file := range files {
			if file.IsDir() {
				continue
			}
			fileName := file.Name()
			if fileName == "garbage.config.tf" {
				fileData, err := os.ReadFile(path.Join("testdata", folderName, fileName))
				if err != nil {
					return nil, err
				}
				resourceBlocks, err := validateHCL(string(fileData))
				if err != nil {
					return nil, fmt.Errorf("%s contains invalid HCL code", path.Join("testdata", folderName, fileName))
				}
				if len(resourceBlocks) == 0 {
					return nil, fmt.Errorf("'%s' doesn't contain any resource blocks", path.Join("testdata", folderName, fileName))
				}
				for _, resourceBlock := range resourceBlocks {
					expectedPrefix := fmt.Sprintf("%s.", expectedResourceType)
					if !strings.HasPrefix(resourceBlock, expectedPrefix) {
						return nil, fmt.Errorf("in file '%s' resource block '%s' doesn't begin with '%s' as expected", path.Join("testdata", folderName, fileName), resourceBlock, expectedPrefix)
					}
				}
				garbageResourceBlocks = append(garbageResourceBlocks, resourceBlocks...)
				absPathGarbageConfigFile, _ = filepath.Abs(path.Join("testdata", folderName, fileName))
			} else if fileName == "golden.config.tf" {
				fileData, err := os.ReadFile(path.Join("testdata", folderName, fileName))
				if err != nil {
					return nil, err
				}
				resourceBlocks, err := validateHCL(string(fileData))
				if err != nil {
					return nil, fmt.Errorf("%s contains invalid HCL code", path.Join("testdata", folderName, fileName))
				}
				if len(resourceBlocks) == 0 {
					return nil, fmt.Errorf("'%s' doesn't contain any resource blocks", path.Join("testdata", folderName, fileName))
				}
				for _, resourceBlock := range resourceBlocks {
					expectedPrefix := fmt.Sprintf("%s.", expectedResourceType)
					if !strings.HasPrefix(resourceBlock, expectedPrefix) {
						return nil, fmt.Errorf("in file '%s' resource block '%s' doesn't begin with '%s' as expected", path.Join("testdata", folderName, fileName), resourceBlock, expectedPrefix)
					}
				}
				goldenResourceBlocks = append(goldenResourceBlocks, resourceBlocks...)
				absPathGoldenConfigFile, _ = filepath.Abs(path.Join("testdata", folderName, fileName))
			} else if fileName == "supporting.config.tf" {
				fileData, err := os.ReadFile(path.Join("testdata", folderName, fileName))
				if err != nil {
					return nil, err
				}
				if _, err := validateHCL(string(fileData)); err != nil {
					return nil, fmt.Errorf("%s contains invalid HCL code", path.Join("testdata", folderName, fileName))
				}
				absPathAdditionalConfigFile, _ = filepath.Abs(path.Join("testdata", folderName, fileName))
			}
		}
		if len(absPathGoldenConfigFile) == 0 {
			return nil, fmt.Errorf("'%s' didn't contain a file `%s`", folderName, "golden.config.tf")
		}
		if len(goldenResourceBlocks) == 0 {
			return nil, fmt.Errorf("'%s' didn't any resource blocks", absPathGoldenConfigFile)
		}
		if len(absPathGarbageConfigFile) == 0 {
			return nil, fmt.Errorf("'%s' didn't contain a file `%s`", folderName, "garbage.config.tf")
		}
		if len(garbageResourceBlocks) == 0 {
			return nil, fmt.Errorf("'%s' didn't any resource blocks", absPathGarbageConfigFile)
		}
		goldenResource := GoldenResource{
			ResourceType:                export.ResourceType(expectedResourceType),
			AbsPathGoldenConfigFile:     absPathGoldenConfigFile,
			AbsPathGarbageConfigFile:    absPathGarbageConfigFile,
			AbsPathAdditionalConfigFile: absPathAdditionalConfigFile,
			GarbageResourceBlocks:       garbageResourceBlocks,
			GoldenResourceBlocks:        goldenResourceBlocks,
		}
		goldenResources = append(goldenResources, goldenResource)
	}
	return goldenResources, nil
}

func TestGolden(t *testing.T) {
	if os.Getenv("DYNATRACE_GOLDEN_STATE_ENABLED") != "true" {
		t.Skip("The Golden State feature isn't enabled. Skipping this test")
	}

	provider := provider.Provider()

	credentials := createCredentials(&ConfigGetter{provider})
	if credentials == nil {
		return
	}

	goldenResources, err := findGoldenResources()
	if err != nil {
		t.Error(err)
		return
	}
	_ = goldenResources

	// t.Skip("skipping early")

	for _, goldenResource := range goldenResources {
		resourceType := goldenResource.ResourceType
		t.Run(string(resourceType), func(t *testing.T) {
			t.Parallel()
			serviceCache := NewServiceCache(t, credentials)
			service := serviceCache.GetService(resourceType)

			stubs, err := service.List(context.Background())
			if err != nil {
				t.Errorf("Listing currently existing settings for resource '%s' failed: %s", resourceType, err.Error())
				return
			}
			if len(stubs) > 0 {
				t.Errorf("Precheck discovered existing settings for resource '%s'. Either this test is getting executed against an environment that is not completely empty (used for other purposes) or a freshly provisioned environment contains default settings for that resource.", resourceType)
				return
			}

			resourceIDs := map[string]string{}
			testCase := &resource.TestCase{
				IsUnitTest: true,
				ProviderFactories: map[string]func() (*schema.Provider, error){
					"dynatrace": func() (*schema.Provider, error) {
						return provider, nil
					},
				},
				Steps: []resource.TestStep{
					{
						Config: genTestCaseConfigFull(t, goldenResource),
						Check: func(s *terraform.State) error {
							for name, resourceState := range s.RootModule().Resources {
								resourceIDs[name] = resourceState.Primary.ID
							}
							return nil
						},
					},
					{Config: genTestCaseConfigReduced(t, goldenResource)},
					{Config: genTestCaseConfigReducedGoldenState(t, goldenResource)},
				},
			}
			resource.Test(t, *testCase)

			// CLEANUP
			deletedIDs, deletionErrors := cleanup(serviceCache, resourceIDs)
			if len(deletedIDs) > 0 {
				for _, deletedID := range deletedIDs {
					t.Logf("For Resource '%s' this ID had to get cleaned up unexpectedly: %s", resourceType, deletedID)
				}
				t.Fail()
			}
			if len(deletionErrors) > 0 {
				for _, deletionError := range deletionErrors {
					t.Logf("During cleanup for Resource '%s' the following unexpected error was received: %s", resourceType, deletionError.Error())
				}
				t.Fail()
			}
		})
	}
}

func cleanup(serviceCache *ServiceCache, resourceIDs map[string]string) ([]string, []error) {
	deletedIDs := []string{}
	errors := []error{}
	for resourceName, resourceID := range resourceIDs {
		resourceType := strings.Split(resourceName, ".")[0]
		deleted, err := deleteResource(export.ResourceType(resourceType), resourceID, serviceCache)
		if deleted {
			deletedIDs = append(deletedIDs, resourceID)
		}
		if err != nil {
			errors = append(errors, err)
		}
	}
	return deletedIDs, errors
}

func existsResource(resourceType export.ResourceType, id string, serviceCache *ServiceCache) bool {
	service := serviceCache.GetService(resourceType)
	dummy := serviceCache.NewSettings(resourceType)
	err := service.Get(context.Background(), id, dummy)
	return err == nil
}

func deleteResource(resourceType export.ResourceType, id string, serviceCache *ServiceCache) (bool, error) {
	if !existsResource(resourceType, id, serviceCache) {
		return false, nil
	}
	service := serviceCache.GetService(resourceType)
	err := service.Delete(context.Background(), id)
	if err != nil {
		if restErr, ok := err.(rest.Error); ok {
			if restErr.Code != 404 {
				return false, fmt.Errorf("ERROR CODE: %d", restErr.Code)
			}
			return false, nil
		}
		if restErr, ok := err.(*rest.Error); ok {
			if restErr.Code != 404 {
				return false, fmt.Errorf("ERROR CODE: %d", restErr.Code)
			}
			return false, nil
		}
		return false, fmt.Errorf("Not A REST ERROR: %T", err)
	}
	return true, nil
}

func genTestCaseConfigFull(t *testing.T, goldenResource GoldenResource) string {
	goldenConfigData, _ := os.ReadFile(goldenResource.AbsPathGoldenConfigFile)
	goldenConfigStr := string(goldenConfigData)
	garbageConfigData, _ := os.ReadFile(goldenResource.AbsPathGarbageConfigFile)
	garbageConfigStr := string(garbageConfigData)
	additionalConfigStr := ""
	if len(goldenResource.AbsPathAdditionalConfigFile) > 0 {
		additionalConfigData, _ := os.ReadFile(goldenResource.AbsPathAdditionalConfigFile)
		additionalConfigStr = string(additionalConfigData)
	}
	result := fmt.Sprintf("%s\n%s\n%s", goldenConfigStr, garbageConfigStr, additionalConfigStr)
	if os.Getenv("DYNATRACE_VERBOSE_TEST_LOGGING") == "true" {
		t.Log("----------------------------------------------------------")
		t.Log(result)
	}
	return result
}

func genTestCaseConfigReduced(t *testing.T, goldenResource GoldenResource) string {
	goldenConfigData, _ := os.ReadFile(goldenResource.AbsPathGoldenConfigFile)
	goldenConfigStr := string(goldenConfigData)
	garbageConfigStr := ""
	for _, garbageResourceBlock := range goldenResource.GarbageResourceBlocks {
		garbageRemovalStr := fmt.Sprintf(`removed {
  from = %s
   lifecycle {
     destroy = false
   }
}`, garbageResourceBlock)
		garbageConfigStr = fmt.Sprintf("%s\n%s", garbageConfigStr, garbageRemovalStr)
	}
	additionalConfigStr := ""
	if len(goldenResource.AbsPathAdditionalConfigFile) > 0 {
		additionalConfigData, _ := os.ReadFile(goldenResource.AbsPathAdditionalConfigFile)
		additionalConfigStr = string(additionalConfigData)
	}
	result := fmt.Sprintf("%s\n%s\n%s", goldenConfigStr, garbageConfigStr, additionalConfigStr)
	if os.Getenv("DYNATRACE_VERBOSE_TEST_LOGGING") == "true" {
		t.Log("----------------------------------------------------------")
		t.Log(result)
	}
	return result
}

func genTestCaseConfigReducedGoldenState(t *testing.T, goldenResource GoldenResource) string {
	goldenConfigData, _ := os.ReadFile(goldenResource.AbsPathGoldenConfigFile)
	goldenConfigStr := string(goldenConfigData)

	goldenResourceBlockIDs := []string{}
	for _, goldenResourceBlock := range goldenResource.GoldenResourceBlocks {
		goldenResourceBlockIDs = append(goldenResourceBlockIDs, fmt.Sprintf("%s.id", goldenResourceBlock))
	}

	goldenStateConfigStr := fmt.Sprintf(`
resource "dynatrace_golden_state" "golden_state" {
    mode = "DELETE"
	%s = [
	  %s
	]
}`, goldenResource.ResourceType, strings.Join(goldenResourceBlockIDs, ", "))

	additionalConfigStr := ""
	if len(goldenResource.AbsPathAdditionalConfigFile) > 0 {
		additionalConfigData, _ := os.ReadFile(goldenResource.AbsPathAdditionalConfigFile)
		additionalConfigStr = string(additionalConfigData)
	}
	result := fmt.Sprintf("%s\n%s\n%s", goldenConfigStr, goldenStateConfigStr, additionalConfigStr)
	if os.Getenv("DYNATRACE_VERBOSE_TEST_LOGGING") == "true" {
		t.Log("----------------------------------------------------------")
		t.Log(result)
	}
	return result
}

// ----------- CREDENTIALS -----------

type ConfigGetter struct {
	provider *schema.Provider
}

func (me *ConfigGetter) Get(key string) any {
	schema, found := me.provider.Schema[key]
	if !found {
		return ""
	}
	if schema.DefaultFunc == nil {
		return ""
	}
	result, _ := schema.DefaultFunc()
	return result
}

func createCredentials(getter config.Getter) *settings.Credentials {
	configResult, _ := config.ProviderConfigureGeneric(context.Background(), getter)
	creds, _ := config.Credentials(configResult, config.CredValNone)
	return creds
}

// ----------- SERVICE CACHE -----------

type ServiceCache struct {
	mu          sync.Mutex
	t           *testing.T
	credentials *settings.Credentials
	services    map[export.ResourceType]settings.CRUDService[settings.Settings]
}

func NewServiceCache(t *testing.T, credentials *settings.Credentials) *ServiceCache {
	return &ServiceCache{t: t, credentials: credentials, services: map[export.ResourceType]settings.CRUDService[settings.Settings]{}}
}

func (me *ServiceCache) GetService(resourceType export.ResourceType) settings.CRUDService[settings.Settings] {
	me.mu.Lock()
	defer me.mu.Unlock()

	if service, found := me.services[resourceType]; found {
		return service
	}
	resourceDescriptor, found := export.AllResources[resourceType]
	if !found {
		return nil
	}
	if service := resourceDescriptor.Service(me.credentials); service != nil {
		me.services[resourceType] = service
		return service
	}

	me.t.Logf("[WARNING] No Service could get created for resource type '%s'", resourceType)
	return nil
}

func (me *ServiceCache) NewSettings(resourceType export.ResourceType) settings.Settings {
	me.mu.Lock()
	defer me.mu.Unlock()

	resourceDescriptor, found := export.AllResources[resourceType]
	if !found {
		return nil
	}
	return resourceDescriptor.NewSettings()
}
