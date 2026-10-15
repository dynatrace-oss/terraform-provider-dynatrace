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

package pipelines_test

import (
	"strings"
	"testing"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/testing/api"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/stretchr/testify/assert"
)

func TestAccOpenPipelineDataSources(t *testing.T) {
	// Note: events.security: feature not enabled

	api.TestAcc(t, api.TestAccOptions{
		Check: func(s *terraform.State) error {
			// every example has one resource and one data source. We get the data via using the first match.
			dataSourceName := getDataSourceName(s)
			resourceId := getResourceID(s)

			assert.NotEmpty(t, dataSourceName)
			assert.NotEmpty(t, resourceId)

			return resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttr(dataSourceName, "items.0.id", resourceId),
				resource.TestCheckResourceAttr(dataSourceName, "items.#", "1"),
			)(s)
		},
	})
}

func TestAccOpenPipelineDataSourcesTwoConfigs(t *testing.T) {
	api.TestAccSingle(t, "testdata/example-two-configs.tf", api.TestAccOptions{
		Check: func(s *terraform.State) error {
			dataSourceName := getDataSourceName(s)
			resourceIds := getResourceIDs(s)

			assert.NotEmpty(t, dataSourceName)
			assert.NotEmpty(t, resourceIds)

			return resource.TestCheckResourceAttr(dataSourceName, "items.#", "2")(s)
		},
	})
}

func TestAccOpenPipelineDataSourcesOneConfig(t *testing.T) {
	api.TestAccSingle(t, "testdata/example-one-config.tf", api.TestAccOptions{
		Check: func(s *terraform.State) error {
			// every example has one resource and one data source. We get the data via using the first match.
			dataSourceName := getDataSourceName(s)
			resourceId := getResourceID(s)

			assert.NotEmpty(t, dataSourceName)
			assert.NotEmpty(t, resourceId)

			return resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttr(dataSourceName, "items.0.id", resourceId),
				resource.TestCheckResourceAttr(dataSourceName, "items.#", "1"),
			)(s)
		},
	})
}

// Returns the ID of the first resource
func getResourceID(s *terraform.State) string {
	for name, rs := range s.RootModule().Resources {
		if !strings.HasPrefix(name, "data.") {
			return rs.Primary.ID
		}
	}
	return ""
}

// Returns the IDs of the all resources
func getResourceIDs(s *terraform.State) []string {
	ids := make([]string, 0, len(s.RootModule().Resources))
	for name, rs := range s.RootModule().Resources {
		if !strings.HasPrefix(name, "data.") {
			ids = append(ids, rs.Primary.ID)
		}
	}
	return ids
}

// Returns the resource name of the first data source
func getDataSourceName(s *terraform.State) string {
	for name := range s.RootModule().Resources {
		if strings.HasPrefix(name, "data.") {
			return name
		}
	}
	return ""
}
