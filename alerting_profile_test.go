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

	"github.com/dtcookie/dynatrace/api/config/alertingprofiles"
	"github.com/dtcookie/terraform-provider-dynatrace/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

type AlertingProfileTest struct {
	resourceKey string
}

func NewAlertingProfileTest() ResourceTest {
	return &AlertingProfileTest{resourceKey: "dynatrace_alerting_profile"}
}

func (test *AlertingProfileTest) Anonymize(m map[string]interface{}) {
	delete(m, "id")
	delete(m, "displayName")
	delete(m, "metadata")
}

func (test *AlertingProfileTest) ResourceKey() string {
	return test.resourceKey
}

func (test *AlertingProfileTest) CreateTestCase(file string, localJSONFile string, t *testing.T) (*resource.TestCase, error) {
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

func TestAccAlertingProfileExampleA(t *testing.T) {
	test := NewAlertingProfileTest()
	var err error
	var testCase *resource.TestCase
	if testCase, err = test.CreateTestCase(
		"test_data/alerting_profiles/example_a.tf",
		"test_data/alerting_profiles/example_a.json",
		t,
	); err != nil {
		t.Fatal(err)
		return
	}
	resource.Test(t, *testCase)
}

func (test *AlertingProfileTest) URL(id string) string {
	envURL := testAccProvider.Meta().(*config.ProviderConfiguration).DTenvURL
	reqPath := "%s/alertingProfiles/%s"
	return fmt.Sprintf(reqPath, envURL, id)
}

func (test *AlertingProfileTest) CheckDestroy(s *terraform.State) error {
	providerConf := testAccProvider.Meta().(*config.ProviderConfiguration)
	restClient := alertingprofiles.NewService(providerConf.DTenvURL, providerConf.APIToken)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "dynatrace_custom_service" {
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
		return fmt.Errorf("Custom Service still exists: %s", rs.Primary.ID)
	}

	return nil
}

func (test *AlertingProfileTest) CheckExists(n string, t *testing.T) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		providerConf := testAccProvider.Meta().(*config.ProviderConfiguration)
		restClient := alertingprofiles.NewService(providerConf.DTenvURL, providerConf.APIToken)

		if rs, ok := s.RootModule().Resources[n]; ok {
			if _, err := restClient.Get(rs.Primary.ID); err != nil {
				return err
			}
			return nil
		}

		return fmt.Errorf("Not found: %s", n)
	}
}
