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

package main_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"reflect"
	"strings"
	"testing"

	"github.com/dtcookie/dynatrace/api/config/managementzones"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/config"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/testbase"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

type ResourceTest interface {
	ResourceKey() string
	CreateTestCase(string, string, *testing.T) (*resource.TestCase, error)
	Anonymize(m map[string]interface{})
	URL(id string) string
}

type ManagementZoneTest struct {
	resourceKey string
}

func NewManagementZoneTest() ResourceTest {
	return &ManagementZoneTest{resourceKey: "dynatrace_management_zone"}
}

func (test *ManagementZoneTest) ResourceKey() string {
	return test.resourceKey
}

func (test *ManagementZoneTest) Anonymize(m map[string]interface{}) {
	delete(m, "id")
	delete(m, "name")
	delete(m, "metadata")
}

func (test *ManagementZoneTest) CreateTestCase(file string, localJSONFile string, t *testing.T) (*resource.TestCase, error) {
	var content []byte
	var err error
	if content, err = ioutil.ReadFile(file); err != nil {
		return nil, err
	}
	config := string(content)
	name := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	resourceName := test.ResourceKey() + "." + name
	config = strings.ReplaceAll(config, "#name#", name)
	return &resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		IDRefreshName:     resourceName,
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      test.CheckDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeAggregateTestCheckFunc(
					test.CheckExists(resourceName),
					compareLocalRemote(test, resourceName, localJSONFile, t),
				),
			},
		},
	}, nil
}

func TestAccManagementZones(t *testing.T) {
	if disabled, ok := testbase.DisabledTests["management_zones"]; ok && disabled {
		t.Skip()
	}

	test := NewManagementZoneTest()
	var err error
	var testCase *resource.TestCase
	if testCase, err = test.CreateTestCase(
		"test_data/management_zones/example_a.tf",
		"test_data/management_zones/example_a.json",
		t,
	); err != nil {
		t.Fatal(err)
		return
	}
	resource.Test(t, *testCase)
}

func (test *ManagementZoneTest) CheckDestroy(s *terraform.State) error {
	providerConf := testAccProvider.Meta().(*config.ProviderConfiguration)
	restClient := managementzones.NewService(providerConf.DTenvURL, providerConf.APIToken)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "dynatrace_management_zone" {
			continue
		}

		id := rs.Primary.ID

		if _, err := restClient.Get(id, false); err != nil {
			// HTTP Response "404 Not Found" signals a success
			if strings.Contains(err.Error(), `"code": 404`) {
				return nil
			}
			// any other error should fail the test
			return err
		}
		return fmt.Errorf("Management Zone still exists: %s", rs.Primary.ID)
	}

	return nil
}

func (test *ManagementZoneTest) URL(id string) string {
	envURL := testAccProvider.Meta().(*config.ProviderConfiguration).DTenvURL
	reqPath := "%s/managementZones/%s?includeProcessGroupReferences=false"
	return fmt.Sprintf(reqPath, envURL, id)
}

func (test *ManagementZoneTest) CheckExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		providerConf := testAccProvider.Meta().(*config.ProviderConfiguration)
		restClient := managementzones.NewService(providerConf.DTenvURL, providerConf.APIToken)

		if rs, ok := s.RootModule().Resources[n]; ok {
			if _, err := restClient.Get(rs.Primary.ID, false); err != nil {
				return err
			}
			return nil
		}

		return fmt.Errorf("not found: %s", n)
	}
}

func compareLocalRemote(test ResourceTest, n string, localJSONFile string, t *testing.T) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		var err error
		var localMap map[string]interface{}
		var remoteMap map[string]interface{}

		if rs, ok := s.RootModule().Resources[n]; ok {
			token := testAccProvider.Meta().(*config.ProviderConfiguration).APIToken
			url := test.URL(rs.Primary.ID)
			if remoteMap, err = loadHTTP(url, token); err != nil {
				return err
			}
			if localMap, err = loadLocal(localJSONFile); err != nil {
				return err
			}
			test.Anonymize(localMap)
			test.Anonymize(remoteMap)
			if !deepEqual(localMap, remoteMap) {
				sLocalMap, _ := json.Marshal(localMap)
				sRemoteMap, _ := json.Marshal(remoteMap)
				return fmt.Errorf("--LOCAL--\n%v\n\n\n--REMOTE--\n%v", string(sLocalMap), string(sRemoteMap))
			}
			return nil
		}

		return fmt.Errorf("not found: %s", n)
	}
}

func deepEqual(a interface{}, b interface{}) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil && b != nil {
		return false
	}
	if a != nil && b == nil {
		return false
	}
	if reflect.TypeOf(a) != reflect.TypeOf(b) {
		return false
	}
	switch ta := a.(type) {
	case map[string]interface{}:
		return deepEqualMap(ta, b.(map[string]interface{}))
	case bool:
		return ta == b.(bool)
	case string:
		return ta == b.(string)
	case float64:
		return ta == b.(float64)
	case []interface{}:
		return deepEqualSlice(ta, b.([]interface{}))
	default:
		panic(fmt.Errorf("unsupported type %T", ta))
	}
}

func deepEqualSlice(a []interface{}, b []interface{}) bool {
	if len(a) != len(b) {
		return false
	}
	for _, va := range a {
		found := false
		for _, vb := range b {
			if deepEqual(va, vb) {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}

func deepEqualMap(a map[string]interface{}, b map[string]interface{}) bool {
	for k, va := range a {
		vb, found := b[k]
		if !found {
			return false
		}
		if !deepEqual(va, vb) {
			return false
		}
	}
	return true
}
