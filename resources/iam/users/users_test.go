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

package users_test

import (
	"os"
	"strings"
	"testing"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/testbase"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const ResourceName = "dynatrace_iam_user"
const TestDataFolder = "../../../test_data/iam/users"
const RequestPath = "%s/settings/objects/%s"

type TestStruct struct {
	resourceKey string
}

func (test *TestStruct) ResourceKey() string {
	return test.resourceKey
}

func (test *TestStruct) CreateTestCase(file string, t *testing.T) (*resource.TestCase, error) {
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
		Steps:             []resource.TestStep{{Config: config}},
	}, nil
}

func TestIAMUsers(t *testing.T) {
	test := &TestStruct{resourceKey: ResourceName}
	var err error
	var testCase *resource.TestCase
	if testCase, err = test.CreateTestCase(
		TestDataFolder+"/example_a.tf",
		t,
	); err != nil {
		t.Fatal(err)
		return
	}
	resource.Test(t, *testCase)
}
