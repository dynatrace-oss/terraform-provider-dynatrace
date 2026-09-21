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

package admin_access

import (
	"testing"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/testing/api"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/provider/envutils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAdminAccessExport(t *testing.T) {
	api.AccEnvsGiven(t)
	exportEnvsGiven(t)

	resourceType := "dynatrace_openpipeline_v2_logs_pipelines"
	resourceName := resourceType + ".pipeline"

	setup, identifier := api.ReadTfConfig(t, "./testdata/example.tf")

	resource.Test(t, resource.TestCase{
		ProviderFactories: api.GetProviderFactories(),
		Steps: []resource.TestStep{
			{
				// Create the config with user1
				Config: setup,
				// Export the config with user2
				Check: func(state *terraform.State) error {
					objectId := state.RootModule().Resources[resourceName].Primary.ID

					// adminAccess: false => not found
					targetDir, err := api.ExportResource(t, resourceType, objectId, false)
					require.NoError(t, err)
					found, err := api.FindExportedResource(targetDir, resourceType, identifier)
					require.NoError(t, err)
					require.Falsef(t, found, "Configuration for resource %s and identifier %s found, even though it should be hidden", objectId, identifier)

					// adminAccess: true => found
					targetDir, err = api.ExportResource(t, resourceType, objectId, true)
					require.NoError(t, err)
					found, err = api.FindExportedResource(targetDir, resourceType, identifier)
					require.NoError(t, err)
					assert.Truef(t, found, "No configuration found for resource %s and identifier %s", objectId, identifier)

					return nil
				},
			},
		},
	})
}

func exportEnvsGiven(t *testing.T) {
	t.Helper()

	if cs := envutils.DynatraceSourceClientSecret.Get(); cs == "" {
		t.Skip("Source client secret has not been set for acceptance tests")
	}
	if cId := envutils.DynatraceSourceClientID.Get(); cId == "" {
		t.Skip("Source client ID has not been set for acceptance tests")
	}
}
