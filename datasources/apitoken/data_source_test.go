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

package apitoken_test

import (
	"testing"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/testing/api"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccApiTokenDataSource(t *testing.T) {
	t.Skip("API doesn't return tokens when queried via OAuth")
	api.AccEnvsGiven(t)

	setup, identifier := api.ReadTfConfig(t, "./testdata/api_token_setup.tf")

	resource.Test(t, resource.TestCase{
		ProviderFactories: api.GetProviderFactories(),
		Steps: []resource.TestStep{
			{
				// First we have to create the API token with a classic token (create not possible with OAuth)
				Config: setup,
				Check: func(state *terraform.State) error {
					// then we can check the data sources with platform and classic credentials
					api.TestAccClassicHybrid(t, api.TestAccOptions{Identifier: identifier})
					return nil
				},
			},
		},
	})
}
