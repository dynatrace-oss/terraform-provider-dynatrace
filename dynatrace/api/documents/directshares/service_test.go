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

package directshares_test

import (
	"testing"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/testing/api"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/provider"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestAccDirectShares(t *testing.T) {
	api.TestAcc(t)
}

func TestUpdateRecipients(t *testing.T) {
	if !api.AccEnvsGiven(t) {
		return
	}

	directShareWithServiceUserAndGroup, _ := api.ReadTfConfig(t, "testdata/create.tf")
	directShareWithTwoGroups, _ := api.ReadTfConfig(t, "testdata/update.tf")

	providerFactories := map[string]func() (*schema.Provider, error){
		"dynatrace": func() (*schema.Provider, error) {
			return provider.Provider(), nil
		},
	}

	t.Run("Removing a user and adding a group as recipients", func(t *testing.T) {
		t.Setenv("TF_VAR_random_name", acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
		testCase := resource.TestCase{
			ProviderFactories: providerFactories,
			Steps: []resource.TestStep{
				{Config: directShareWithServiceUserAndGroup},
				{Config: directShareWithTwoGroups}, // with an empty one, the API would reject it
			},
		}
		resource.Test(t, testCase)
	})
}
