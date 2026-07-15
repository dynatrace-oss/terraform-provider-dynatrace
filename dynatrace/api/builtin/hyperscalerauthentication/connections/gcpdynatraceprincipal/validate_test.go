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

package gcpdynatraceprincipal_test

import (
	"os"
	"testing"

	gcpservice "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/hyperscalerauthentication/connections/gcp"
	gcpsettings "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/hyperscalerauthentication/connections/gcp/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/hyperscalerauthentication/connections/gcpdynatraceprincipal"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings/services/settings20"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/provider/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Guards against schema drift in ValidConnection. The dynatrace_gcp_principal resource triggers
// Dynatrace GCP Principal creation as a side effect of submitting ValidConnection to the validate
// endpoint; if the payload is rejected, no principal will ever be created.
func TestValidateIsSuccessfulForGcpConnection(t *testing.T) {
	envURL := os.Getenv("DYNATRACE_ENV_URL")
	apiToken := os.Getenv("DYNATRACE_API_TOKEN")
	if envURL == "" || apiToken == "" {
		t.Skip("Environment Variables DYNATRACE_ENV_URL and DYNATRACE_API_TOKEN must be specified")
	}

	clientSet := &config.ProviderConfiguration{EnvironmentURL: envURL, APIToken: apiToken}
	service := settings20.Service[*gcpsettings.Settings](clientSet, gcpservice.SchemaID, gcpservice.SchemaVersion)

	validator, ok := service.(settings.Validator[*gcpsettings.Settings])
	require.True(t, ok, "settings20 service must implement settings.Validator")

	err := validator.Validate(t.Context(), &gcpdynatraceprincipal.ValidConnection)
	assert.NoError(t, err, "ValidConnection in service.go was rejected by the API — update the fixture or principal creation might never be triggered.")
}
