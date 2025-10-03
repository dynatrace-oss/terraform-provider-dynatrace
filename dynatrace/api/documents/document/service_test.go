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

package documents_test

import (
	"regexp"
	"testing"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/testing/api"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/provider"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestAccDocuments(t *testing.T) {
	// this also tests "no custom ID" => "no change after apply & plan"
	api.TestAcc(t)
}

func TestAccDocumentsCustomID(t *testing.T) {
	if !api.AccEnvsGiven(t) {
		return
	}

	const resourceName = "dynatrace_document.this"
	configNoID, _ := api.ReadTfConfig(t, "testdata-2/without-custom-id.tf")
	configCustomID, customID := api.ReadTfConfig(t, "testdata-2/with-custom-id.tf")
	configCustomID2, customID2 := api.ReadTfConfig(t, "testdata-2/with-custom-id.tf")
	configUUID, _ := api.ReadTfConfig(t, "testdata-2/with-custom-uuid.tf")

	providerFactories := map[string]func() (*schema.Provider, error){
		"dynatrace": func() (*schema.Provider, error) {
			return provider.Provider(), nil
		},
	}

	t.Run("No ID to custom ID works", func(t *testing.T) {
		testCase := resource.TestCase{
			ProviderFactories: providerFactories,
			Steps: []resource.TestStep{
				{Config: configNoID},
				{Config: configCustomID, ExpectNonEmptyPlan: true, PlanOnly: true},
				{
					Config: configCustomID,
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr(resourceName, "custom_id", customID),
						resource.TestCheckResourceAttr(resourceName, "id", customID),
					),
				},
			},
		}
		resource.Test(t, testCase)
	})

	t.Run("Custom ID to a different custom ID works", func(t *testing.T) {
		testCase := resource.TestCase{
			ProviderFactories: providerFactories,
			Steps: []resource.TestStep{
				{Config: configCustomID},
				{Config: configCustomID2, ExpectNonEmptyPlan: true, PlanOnly: true},
				{
					Config: configCustomID2,
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr(resourceName, "custom_id", customID2),
					),
				},
			},
		}
		resource.Test(t, testCase)
	})

	t.Run("Custom ID to nothing results in no change", func(t *testing.T) {
		testCase := resource.TestCase{
			ProviderFactories: providerFactories,
			Steps: []resource.TestStep{
				{Config: configCustomID},
				{
					Config: configNoID,
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr(resourceName, "custom_id", customID),
					),
					PlanOnly: true,
				},
			},
		}
		resource.Test(t, testCase)
	})

	t.Run("A UUID as a custom ID results in an error", func(t *testing.T) {
		testCase := resource.TestCase{
			ProviderFactories: providerFactories,
			Steps: []resource.TestStep{
				{Config: configUUID, ExpectError: regexp.MustCompile("not be a UUID"), PlanOnly: true},
			},
		}
		resource.Test(t, testCase)
	})
}
