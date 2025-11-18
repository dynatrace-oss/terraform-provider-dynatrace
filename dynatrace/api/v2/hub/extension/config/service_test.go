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

package extension_config_test

import (
	"testing"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/testing/api"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/provider"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestAccExtensionConfiguration(t *testing.T) {
	// this also tests "no scope" => "no change after apply & plan"
	api.TestAcc(t)
}

func TestAccExtensionsCustomScope(t *testing.T) {
	if !api.AccEnvsGiven(t) {
		return
	}

	const resourceName = "dynatrace_hub_extension_configuration.this"
	configNoScope, _ := api.ReadTfConfig(t, "testdata/terraform/noscope.tf")
	configWithScope, _ := api.ReadTfConfig(t, "testdata/terraform/withscope.tf")
	scope := "environment"

	providerFactories := map[string]func() (*schema.Provider, error){
		"dynatrace": func() (*schema.Provider, error) {
			return provider.Provider(), nil
		},
	}

	t.Run("No scope to default custom scope yields empty plan", func(t *testing.T) {
		testCase := resource.TestCase{
			ProviderFactories: providerFactories,
			Steps: []resource.TestStep{
				{Config: configNoScope},
				{
					Config: configWithScope,
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr(resourceName, "scope", scope),
					),
					PlanOnly: true,
				},
			},
		}
		resource.Test(t, testCase)
	})

	t.Run("Removal of default custom scope yields empty plan", func(t *testing.T) {
		testCase := resource.TestCase{
			ProviderFactories: providerFactories,
			Steps: []resource.TestStep{
				{Config: configWithScope},
				{
					Config: configNoScope,
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr(resourceName, "scope", scope),
					),
					PlanOnly: true,
				},
			},
		}
		resource.Test(t, testCase)
	})
}
