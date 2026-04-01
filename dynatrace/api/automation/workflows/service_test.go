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

package workflows_test

import (
	"testing"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/testing/api"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/provider"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestAccWorkflows(t *testing.T) {
	api.TestAcc(t)
}

// TestAccWorkflows_Update checks if set items can be updated to empty, like guide or description
func TestAccWorkflows_Update(t *testing.T) {
	if !api.AccEnvsGiven(t) {
		return
	}

	create, _ := api.ReadTfConfig(t, "testdata/create.tf")
	update, _ := api.ReadTfConfig(t, "testdata/update.tf") // name change because of #name# doesn't matter here. Important is the non-empty plan for guide and description

	providerFactories := map[string]func() (*schema.Provider, error){
		"dynatrace": func() (*schema.Provider, error) {
			return provider.Provider(), nil
		},
	}

	resource.Test(t, resource.TestCase{
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{Config: create},
			{Config: update},
		},
	})
}
