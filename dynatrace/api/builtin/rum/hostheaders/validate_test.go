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

package hostheaders_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	hostheadersvc "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/rum/hostheaders"
	hostheaders "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/rum/hostheaders/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings/services/settings20"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/provider/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestValidateDoesNotCreateObjects guards against settings20.Validate accidentally
// omitting the validateOnly=true query parameter and creating a real settings object.
// builtin:rum.host-headers is used as the test vehicle because it has no custom
// validators, no side-effects, and a single plain-text field.
func TestValidateDoesNotCreateObjects(t *testing.T) {
	envURL := os.Getenv("DYNATRACE_ENV_URL")
	apiToken := os.Getenv("DYNATRACE_API_TOKEN")
	if envURL == "" || apiToken == "" {
		t.Skip("Environment Variables DYNATRACE_ENV_URL and DYNATRACE_API_TOKEN must be specified")
	}

	headerName := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)

	credentials := &config.ProviderConfiguration{EnvironmentURL: envURL, APIToken: apiToken}
	service := settings20.Service[*hostheaders.Settings](credentials, hostheadersvc.SchemaID, hostheadersvc.SchemaVersion)

	// If Validate accidentally creates an object, clean it up so subsequent
	// test runs are not polluted.
	t.Cleanup(func() {
		cleanupContext := context.Background()
		stubsAfterTest, err := service.List(cleanupContext)
		if err != nil {
			return
		}
		for _, stub := range stubsAfterTest {
			if stub.Name == headerName {
				err = service.Delete(cleanupContext, stub.ID)
				require.NoError(t, err, fmt.Sprintf("Test cleanup failed for settings object %s. Manual cleanup necessary", stub.ID))
			}
		}
	})

	validator, ok := service.(settings.Validator[*hostheaders.Settings])
	require.True(t, ok, "settings20 service must implement settings.Validator")

	err := validator.Validate(t.Context(), &hostheaders.Settings{HeaderName: headerName})
	assert.NoError(t, err)

	stubsAfter, err := service.List(t.Context())
	require.NoError(t, err)

	for _, stub := range stubsAfter {
		assert.NotEqual(t, headerName, stub.Name, "Validate must not create any settings objects")
	}
}
