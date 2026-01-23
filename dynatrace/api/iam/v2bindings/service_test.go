//go:build integration

/*
 * @license
 * Copyright 2026 Dynatrace LLC
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package v2bindings_test

import (
	"os"
	"testing"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/testing/api"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/provider"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestAccV2Bindings(t *testing.T) {
	if !api.AccEnvsGiven(t) {
		return
	}

	accountID := os.Getenv("DT_ACCOUNT_ID")
	//fallback to DYNATRACE_ACCOUNT_ID
	if accountID == "" {
		accountID = os.Getenv("DYNATRACE_ACCOUNT_ID")
	}
	t.Setenv("TF_VAR_ACCOUNT_ID", accountID)
	configCreate, _ := api.ReadTfConfig(t, "testdata/create.tf")
	configUpdate, _ := api.ReadTfConfig(t, "testdata/update.tf")

	providerFactories := map[string]func() (*schema.Provider, error){
		"dynatrace": func() (*schema.Provider, error) {
			return provider.Provider(), nil
		},
	}

	t.Run("Removal of policies works", func(t *testing.T) {
		testCase := resource.TestCase{
			ProviderFactories: providerFactories,
			Steps: []resource.TestStep{
				{Config: configCreate},
				{Config: configUpdate}, // removes one policy, adds a new policy and updates an existing policy
			},
		}
		resource.Test(t, testCase)
	})
}
