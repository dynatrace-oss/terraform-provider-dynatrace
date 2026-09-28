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

package locations_test

import (
	"path"
	"testing"
	"time"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/testing/api"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccSyntheticLocations(t *testing.T) {
	api.TestAcc(t)
}

func TestAccSyntheticLocation_UseNewKubernetesVersion(t *testing.T) {
	if !api.AccEnvsGiven(t) {
		return
	}

	const testcaseFolder = "testcases/use-new-kubernetes-version"
	const resourceNameIdentifier = "dynatrace_synthetic_location.location"

	// A freshly created location rejects updates with a misleading `404 Location <id> not found`
	// for a short while, so every step but the first waits for the location to settle.
	const settleDelay = 5 * time.Second
	settle := func() { time.Sleep(settleDelay) }

	configOmitted, identifier := api.ReadTfConfig(t, path.Join(testcaseFolder, "omitted.tf"))
	configExplicitTrue := api.ReadTfConfigWithIdentifier(t, path.Join(testcaseFolder, "explicit_true.tf"), identifier)
	configExplicitFalse := api.ReadTfConfigWithIdentifier(t, path.Join(testcaseFolder, "explicit_false.tf"), identifier)

	resource.Test(t, resource.TestCase{
		ProviderFactories: api.GetProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: configOmitted,
				Check:  resource.TestCheckResourceAttr(resourceNameIdentifier, "use_new_kubernetes_version", "true"),
			},
			{
				PreConfig: settle,
				Config:    configExplicitTrue,
				Check:     resource.TestCheckResourceAttr(resourceNameIdentifier, "use_new_kubernetes_version", "true"),
			},
			{
				PreConfig: settle,
				Config:    configExplicitFalse,
				Check:     resource.TestCheckResourceAttr(resourceNameIdentifier, "use_new_kubernetes_version", "false"),
			},
			{
				PreConfig: settle,
				Config:    configOmitted,
				Check:     resource.TestCheckResourceAttr(resourceNameIdentifier, "use_new_kubernetes_version", "true"),
			},
		},
	})
}
