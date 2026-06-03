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

package gcp_test

import (
	"regexp"
	"testing"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/hyperscalerauthentication/connections/gcp"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/testing/api"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/provider"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/provider/envutils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestGcpConnection(t *testing.T) {
	t.Skip("Enable this test as soon as GCP test environment setup is done.")

	api.TestAcc(t, api.TestAccOptions{
		ExternalProviders: map[string]resource.ExternalProvider{
			"google": {VersionConstraint: "7.33.0", Source: "hashicorp/google"},
			"time":   {VersionConstraint: "0.14.0", Source: "hashicorp/time"},
		},
	})
}

// TestAccGcpConnectionAuthenticationFailure verifies that the Dynatrace API still rejects an
// unauthorized service account impersonation with a constraint violation whose message contains
// gcp.RetryableAuthenticationErrorMessage. The Create retry only kicks in on that exact message, so
// this test binds the asserted message and the message used by the retry together: if the API ever
// changes the wording, this test fails and forces the constant to be updated in lockstep.
//
// KNOWN FLAKE (accepted for now): because the message is classified as retryable, the
// permanently-failing service account is retried until the create timeout elapses, and the retry
// budget currently coincides with the resource context deadline. If a request is in-flight when the
// deadline fires, the final error can be a context-cancellation error instead of the API message,
// causing ExpectError to miss (~6% of runs). A follow-up PR decouples the retry budget from the
// context deadline (reserving a buffer so the final attempt completes against a live context),
// which makes this deterministic. Until then this opt-in test (it requires acceptance env vars plus
// DT_GCP_TEST_UNIMPERSONABLE_SERVICE_ACCOUNT) is skipped, pending the GCP connection resource being
// fully wired up in a follow-up PR.
func TestAccGcpConnectionAuthenticationFailure(t *testing.T) {
	t.Skip("Enable this test as soon as the GCP connection resource is wired.")

	if !api.AccEnvsGiven(t) {
		return
	}

	if envutils.DTGCPTestUnimpersonableServiceAccount.Get() == "" {
		t.Skipf("%s is not set; skipping GCP authentication failure test", envutils.DTGCPTestUnimpersonableServiceAccount.Key)
	}

	config, _ := api.ReadTfConfig(t, "testdata-auth-failure/gcp_connection.tf")
	providerFactories := map[string]func() (*schema.Provider, error){
		"dynatrace": func() (*schema.Provider, error) {
			return provider.Provider(), nil
		},
	}

	resource.Test(t, resource.TestCase{
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config:      config,
				ExpectError: regexp.MustCompile(regexp.QuoteMeta(gcp.RetryableAuthenticationErrorMessage)),
			},
		},
	})
}
