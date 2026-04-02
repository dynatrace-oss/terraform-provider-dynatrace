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

package monitoringconfigurations_test

import (
	"testing"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/testing/api"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/provider"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/stretchr/testify/assert"
)

func TestMonitoringConfiguration(t *testing.T) {
	api.TestAcc(t)
}

func TestMonitoringConfiguration_Update_NoRecreate(t *testing.T) {
	if !api.AccEnvsGiven(t) {
		return
	}

	configCreate, _ := api.ReadTfConfig(t, "testcases/update-no-recreate/create.tf")
	configUpdate, _ := api.ReadTfConfig(t, "testcases/update-no-recreate/update.tf")

	providerFactories := map[string]func() (*schema.Provider, error){
		"dynatrace": func() (*schema.Provider, error) {
			return provider.Provider(), nil
		},
	}
	const resourceNameIdentifier = "dynatrace_hub_extension_v2_config.config"
	var resourceIDCreate, resourceIDUpdate string
	var resourceCount int

	resource.Test(t, resource.TestCase{
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: configCreate,
				Check: func(s *terraform.State) (err error) {
					resourceIDCreate = s.Modules[0].Resources[resourceNameIdentifier].Primary.ID
					return err
				},
			},
			{
				Config: configUpdate,
				Check: func(s *terraform.State) (err error) {
					resourceIDUpdate = s.Modules[0].Resources[resourceNameIdentifier].Primary.ID
					resourceCount = len(s.Modules[0].Resources)
					return err
				},
			},
		},
	})
	assert.Equal(t, resourceIDCreate, resourceIDUpdate, "Resource was recreated instead of updated")
	assert.Equal(t, 1, resourceCount, "Expected exactly one resource in the state after update")
}

func TestMonitoringConfiguration_Update_Recreate(t *testing.T) {
	if !api.AccEnvsGiven(t) {
		return
	}

	configCreate, _ := api.ReadTfConfig(t, "testcases/recreate-scope-change/create.tf")
	configUpdate, _ := api.ReadTfConfig(t, "testcases/recreate-scope-change/update.tf")

	providerFactories := map[string]func() (*schema.Provider, error){
		"dynatrace": func() (*schema.Provider, error) {
			return provider.Provider(), nil
		},
	}
	const resourceNameIdentifier = "dynatrace_hub_extension_v2_config.config"
	var resourceIDCreate, resourceIDUpdate string
	var resourceCount int

	resource.Test(t, resource.TestCase{
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: configCreate,
				Check: func(s *terraform.State) (err error) {
					resourceIDCreate = s.Modules[0].Resources[resourceNameIdentifier].Primary.ID
					return err
				},
			},
			{
				Config: configUpdate,
				Check: func(s *terraform.State) (err error) {
					resourceIDUpdate = s.Modules[0].Resources[resourceNameIdentifier].Primary.ID
					resourceCount = len(s.Modules[0].Resources)
					return err
				},
			},
		},
	})
	assert.NotEqual(t, resourceIDCreate, resourceIDUpdate, "Resource was updated instead of recreated")
	assert.Equal(t, 2 /*data source + resource*/, resourceCount, "Expected exactly two resources in the state after update")
}
