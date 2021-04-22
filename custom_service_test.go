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
	"net/http"
	"strings"
	"testing"

	"github.com/dtcookie/dynatrace/api/config/customservices"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/config"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/testbase"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func loadHTTP(url string, token string) (map[string]interface{}, error) {
	var err error
	var request *http.Request
	var response *http.Response
	var data []byte

	if request, err = http.NewRequest("GET", url, nil); err != nil {
		return nil, err
	}
	request.Header.Set("Authorization", "Api-Token "+token)

	client := &http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(nil)}}
	if response, err = client.Do(request); err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if data, err = ioutil.ReadAll(response.Body); err != nil {
		return nil, err
	}

	m := map[string]interface{}{}
	if err = json.Unmarshal(data, &m); err != nil {
		return nil, err
	}
	return m, nil
}

func loadLocal(file string) (map[string]interface{}, error) {
	var err error
	var data []byte
	if data, err = ioutil.ReadFile(file); err != nil {
		return nil, err
	}
	m := map[string]interface{}{}
	if err = json.Unmarshal(data, &m); err != nil {
		return nil, err
	}
	return m, nil
}

type CustomServiceTest struct {
	resourceKey string
}

func NewCustomServiceTest() ResourceTest {
	return &CustomServiceTest{resourceKey: "dynatrace_custom_service"}
}

func (test *CustomServiceTest) Anonymize(m map[string]interface{}) {
	delete(m, "id")
	delete(m, "name")
	delete(m, "metadata")
	if rules, found := m["rules"]; found {
		for _, rule := range rules.([]interface{}) {
			typedRule := rule.(map[string]interface{})
			delete(typedRule, "id")
			if methodRules, found := typedRule["methodRules"]; found {
				for _, methodRule := range methodRules.([]interface{}) {
					typedMethodRule := methodRule.(map[string]interface{})
					delete(typedMethodRule, "id")
				}
			}
		}
	}
}

func (test *CustomServiceTest) ResourceKey() string {
	return test.resourceKey
}

func (test *CustomServiceTest) CreateTestCase(file string, localJSONFile string, t *testing.T) (*resource.TestCase, error) {
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

func TestAccCustomServices(t *testing.T) {
	if disabled, ok := testbase.DisabledTests["custom_services"]; ok && disabled {
		t.Skip()
	}
	test := NewCustomServiceTest()
	var err error
	var testCase *resource.TestCase
	if testCase, err = test.CreateTestCase(
		"test_data/custom_services/example_a.tf",
		"test_data/custom_services/example_a.json",
		t,
	); err != nil {
		t.Fatal(err)
		return
	}
	resource.Test(t, *testCase)
}

func (test *CustomServiceTest) URL(id string) string {
	envURL := testAccProvider.Meta().(*config.ProviderConfiguration).DTenvURL
	reqPath := "%s/service/customServices/%v/%s?includeProcessGroupReferences=false"
	return fmt.Sprintf(reqPath, envURL, customservices.Technologies.Java, id)
}

func (test *CustomServiceTest) CheckDestroy(s *terraform.State) error {
	providerConf := testAccProvider.Meta().(*config.ProviderConfiguration)
	restClient := customservices.NewService(providerConf.DTenvURL, providerConf.APIToken)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "dynatrace_custom_service" {
			continue
		}

		id := rs.Primary.ID

		if _, err := restClient.Get(id, customservices.Technologies.Java, false); err != nil {
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

func (test *CustomServiceTest) CheckExists(n string, t *testing.T) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		providerConf := testAccProvider.Meta().(*config.ProviderConfiguration)
		restClient := customservices.NewService(providerConf.DTenvURL, providerConf.APIToken)

		if rs, ok := s.RootModule().Resources[n]; ok {
			if _, err := restClient.Get(rs.Primary.ID, customservices.Technologies.Java, false); err != nil {
				return err
			}
			return nil
		}

		return fmt.Errorf("Not found: %s", n)
	}
}
