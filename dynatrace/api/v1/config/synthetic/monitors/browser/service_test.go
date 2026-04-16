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

package browser_test

import (
	"testing"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/testing/api"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccBrowserMonitors(t *testing.T) {
	api.TestAcc(t, api.TestAccOptions{ExternalProviders: map[string]resource.ExternalProvider{"time": {Source: "hashicorp/time"}}})
}

func TestAccTestCasesBrowserMonitors(t *testing.T) {
	t.Skip("Skipping because the eventual consistency time is way too high")

	api.TestAccTestCases(t, api.TestCaseAccOptions{ExternalProviders: map[string]resource.ExternalProvider{"time": {Source: "hashicorp/time"}}})
}
