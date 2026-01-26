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

package permissions_test

import (
	"os"
	"testing"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/testing/api"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/provider"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestAccPermissions(t *testing.T) {
	if !api.AccEnvsGiven(t) {
		return
	}

	accountID := os.Getenv("DT_ACCOUNT_ID")
	//fallback to DYNATRACE_ACCOUNT_ID
	if accountID == "" {
		accountID = os.Getenv("DYNATRACE_ACCOUNT_ID")
	}
	t.Setenv("TF_VAR_ACCOUNT_ID", accountID)

	groupName := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	t.Setenv("TF_VAR_GROUP_NAME", groupName)

	configCreate, _ := api.ReadTfConfig(t, "testdata/create.tf")

	providerFactories := map[string]func() (*schema.Provider, error){
		"dynatrace": func() (*schema.Provider, error) {
			return provider.Provider(), nil
		},
	}

	t.Run("Create group with permission", func(t *testing.T) {
		testCase := resource.TestCase{
			ProviderFactories: providerFactories,
			Steps: []resource.TestStep{
				{Config: configCreate}, // creates group and permission
			},
		}
		resource.Test(t, testCase)
	})
}
