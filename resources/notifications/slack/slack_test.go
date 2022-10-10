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

package slack_test

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/dtcookie/dynatrace/api/config/v2/notifications"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/config"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/testbase"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

const ResourceName = "dynatrace_slack_notification"
const TestDataFolder = "../../../test_data/notifications/slack"
const RequestPath = "%s/settings/objects/%s"

type TestStruct struct {
	resourceKey string
}

func (test *TestStruct) Anonymize(m map[string]interface{}) {
	delete(m, "author")
	delete(m, "created")
	delete(m, "modified")
	delete(m, "objectId")
	delete(m, "schemaId")
	delete(m, "schemaVersion")
	delete(m, "scope")
	delete(m, "summary")
	delete(m, "updateToken")
	value := m["value"].(map[string]interface{})
	slackNotification := value["slackNotification"].(map[string]interface{})
	delete(slackNotification, "url")
	delete(value, "displayName")
}

func (test *TestStruct) ResourceKey() string {
	return test.resourceKey
}

func (test *TestStruct) CreateTestCase(file string, localJSONFile string, t *testing.T) (*resource.TestCase, error) {
	var content []byte
	var err error
	if content, err = os.ReadFile(file); err != nil {
		return nil, err
	}
	config := string(content)
	name := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	resourceName := test.ResourceKey() + "." + name
	config = strings.ReplaceAll(config, "#name#", name)
	return &resource.TestCase{
		PreCheck:          func() { testbase.TestAccPreCheck(t) },
		IDRefreshName:     resourceName,
		ProviderFactories: testbase.TestAccProviderFactories,
		CheckDestroy:      test.CheckDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeAggregateTestCheckFunc(
					CheckExists(resourceName, t),
					testbase.CompareLocalRemoteExt(test, resourceName, localJSONFile, t, false),
				),
			},
		},
	}, nil
}

func TestSlackNotifications(t *testing.T) {
	test := &TestStruct{resourceKey: ResourceName}
	var err error
	var testCase *resource.TestCase
	if testCase, err = test.CreateTestCase(
		TestDataFolder+"/example_a.tf",
		TestDataFolder+"/example_a.json",
		t,
	); err != nil {
		t.Fatal(err)
		return
	}
	resource.Test(t, *testCase)
}

func (test *TestStruct) URL(id string) string {
	envURL := testbase.TestAccProvider.Meta().(*config.ProviderConfiguration).DTApiV2URL
	reqPath := RequestPath
	return fmt.Sprintf(reqPath, envURL, id)
}

func (test *TestStruct) CheckDestroy(s *terraform.State) error {
	providerConf := testbase.TestAccProvider.Meta().(*config.ProviderConfiguration)
	restClient := notifications.NewService(providerConf.DTApiV2URL, providerConf.APIToken)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != ResourceName {
			continue
		}

		id := rs.Primary.ID
		var cfg notifications.Notification
		if err := restClient.Get(id, &cfg); err != nil {
			// HTTP Response "404 Not Found" signals a success
			if strings.Contains(err.Error(), `"code": 404`) {
				return nil
			}
			// any other error should fail the test
			return err
		}
		return fmt.Errorf("Configuration still exists: %s", rs.Primary.ID)
	}

	return nil
}

func CheckExists(n string, t *testing.T) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		providerConf := testbase.TestAccProvider.Meta().(*config.ProviderConfiguration)
		restClient := notifications.NewService(providerConf.DTApiV2URL, providerConf.APIToken)

		var err error
		var cnt = 0

		for cnt < 10 {
			if rs, ok := s.RootModule().Resources[n]; ok {
				var cfg notifications.Notification
				if err = restClient.Get(rs.Primary.ID, &cfg); err != nil {
					cnt++
				}
				return nil
			}
		}
		return err

		// return fmt.Errorf("Not found: %s", n)
	}
}
