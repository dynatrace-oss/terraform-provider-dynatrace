//go:build integration

/**
* @license
* Copyright 2026 Dynatrace LLC
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

package configuration_test

import (
	"os"
	"testing"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/testing/api"
)

// TestAccSyntheticConfigurationV2 runs the acceptance tests without an API token, as the resource
// only needs platform credentials. The OAuth leg disables platform tokens and WIF, the platform token
// leg disables the OAuth client secret and WIF.
func TestAccSyntheticConfigurationV2(t *testing.T) {
	t.Run("OAuth", func(t *testing.T) {
		t.Setenv("DYNATRACE_API_TOKEN", "")
		t.Setenv("DT_API_TOKEN", "")
		t.Setenv("DYNATRACE_PLATFORM_TOKEN", "")
		t.Setenv("DT_PLATFORM_TOKEN", "")
		t.Setenv("DYNATRACE_WIF_AUDIENCE", "")
		api.TestAcc(t)
		api.TestAccTestCases(t)
	})

	t.Run("PlatformToken", func(t *testing.T) {
		if os.Getenv("DYNATRACE_PLATFORM_TOKEN") == "" && os.Getenv("DT_PLATFORM_TOKEN") == "" {
			t.Skip("no platform token configured")
		}
		t.Setenv("DYNATRACE_API_TOKEN", "")
		t.Setenv("DT_API_TOKEN", "")
		t.Setenv("DT_CLIENT_SECRET", "")
		t.Setenv("DYNATRACE_CLIENT_SECRET", "")
		t.Setenv("DYNATRACE_WIF_AUDIENCE", "")
		api.TestAcc(t)
		api.TestAccTestCases(t)
	})
}
