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

package monitors_test

import (
	"testing"
	"time"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/testing/api"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/stretchr/testify/require"
)

// Note: Eventual consistency was still seen with 50s
var eventualConsistencyTime = 70 * time.Second

func TestAccNetworkAvailabilityMonitors(t *testing.T) {
	api.TestAcc(t, api.TestAccOptions{Check: waitForTagsToBeInSync(t)})
}

func TestAccTestCasesNetworkMonitors(t *testing.T) {
	api.TestAccTestCases(t, api.TestCaseAccOptions{Check: waitForTagsToBeInSync(t)})
}

// waitForTagsToBeInSync waits a fixed duration for the API to be consistent (eventual consistency)
func waitForTagsToBeInSync(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		res := api.GetFirstResource(s)
		require.NotNil(t, res)
		if val, ok := res.Primary.Attributes["tags.#"]; ok && val != "" && val != "0" {
			time.Sleep(eventualConsistencyTime)
		}
		return nil
	}
}
