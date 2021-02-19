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
	"fmt"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/config"

	"github.com/dtcookie/dynatrace/api/config/autotags"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

type AutoTagTest struct {
	resourceKey string
	testDataFld string
	reqPath     string
	display     string
}

func NewAutoTagTest() *AutoTagTest {
	return &AutoTagTest{
		resourceKey: "dynatrace_autotag",
		testDataFld: "auto_tags",
		reqPath:     "autoTags",
		display:     "Auto Tag",
	}
}

func (test *AutoTagTest) Anonymize(m map[string]interface{}) {
	delete(m, "id")
	delete(m, "name")
	delete(m, "metadata")
}

func (test *AutoTagTest) ResourceKey() string {
	return test.resourceKey
}

func (test *AutoTagTest) CreateTestCase(file string, localJSONFile string, t *testing.T) (*resource.TestCase, error) {
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
					test.CheckExists(resourceName, t),
					compareLocalRemote(test, resourceName, localJSONFile, t),
				),
			},
		},
	}, nil
}

func TestAccAutoTagExampleA(t *testing.T) {
	test := NewAutoTagTest()
	var err error
	var testCase *resource.TestCase
	if testCase, err = test.CreateTestCase(
		fmt.Sprintf("test_data/%s/example_a.tf", test.testDataFld),
		fmt.Sprintf("test_data/%s/example_a.json", test.testDataFld),
		t,
	); err != nil {
		t.Fatal(err)
		return
	}
	resource.Test(t, *testCase)
}

func (test *AutoTagTest) URL(id string) string {
	envURL := testAccProvider.Meta().(*config.ProviderConfiguration).DTenvURL
	reqPath := "%s/" + test.reqPath + "/%s"
	return fmt.Sprintf(reqPath, envURL, id)
}

func (test *AutoTagTest) CheckDestroy(s *terraform.State) error {
	providerConf := testAccProvider.Meta().(*config.ProviderConfiguration)
	restClient := autotags.NewService(providerConf.DTenvURL, providerConf.APIToken)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != test.resourceKey {
			continue
		}

		id := rs.Primary.ID

		if _, err := restClient.Get(id); err != nil {
			// HTTP Response "404 Not Found" signals a success
			if strings.Contains(err.Error(), `"code": 404`) {
				return nil
			}
			// any other error should fail the test
			return err
		}
		return fmt.Errorf("%s still exists: %s", test.display, rs.Primary.ID)
	}

	return nil
}

func (test *AutoTagTest) CheckExists(n string, t *testing.T) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		providerConf := testAccProvider.Meta().(*config.ProviderConfiguration)
		restClient := autotags.NewService(providerConf.DTenvURL, providerConf.APIToken)

		if rs, ok := s.RootModule().Resources[n]; ok {
			if _, err := restClient.Get(rs.Primary.ID); err != nil {
				return err
			}
			return nil
		}

		return fmt.Errorf("Not found: %s", n)
	}
}
