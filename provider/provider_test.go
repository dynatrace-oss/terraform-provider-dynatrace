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

package provider_test

import (
	"context"
	"os"
	"testing"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/provider"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/provider/config"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIAMClientID(t *testing.T) {
	if len(os.Getenv("TF_ACC")) > 0 {
		return
	}
	provider := provider.Provider()

	for _, envVarName := range []string{
		"IAM_CLIENT_ID", "DYNATRACE_IAM_CLIENT_ID", "DT_IAM_CLIENT_ID", "DT_CLIENT_ID", "DYNATRACE_CLIENT_ID",
		"AUTOMATION_CLIENT_ID", "DYNATRACE_AUTOMATION_CLIENT_ID", "DT_AUTOMATION_CLIENT_ID",
	} {
		t.Run(envVarName, func(t *testing.T) {
			iam_client_id := uuid.NewString()
			t.Setenv(envVarName, iam_client_id)

			credentials := createCredentials(&config.ConfigGetter{Provider: provider})
			if credentials == nil {
				return
			}
			require.Equal(t, iam_client_id, credentials.IAM.ClientID)
			assert.Equal(t, iam_client_id, credentials.OAuth.ClientID)
		})
	}
}

func TestIAMClientSecret(t *testing.T) {
	if len(os.Getenv("TF_ACC")) > 0 {
		return
	}
	provider := provider.Provider()

	for _, envVarName := range []string{
		"IAM_CLIENT_SECRET", "DYNATRACE_IAM_CLIENT_SECRET", "DT_IAM_CLIENT_SECRET", "DYNATRACE_CLIENT_SECRET", "DT_CLIENT_SECRET",
		"AUTOMATION_CLIENT_SECRET", "DYNATRACE_AUTOMATION_CLIENT_SECRET", "DT_AUTOMATION_CLIENT_SECRET",
	} {
		t.Run(envVarName, func(t *testing.T) {
			iam_client_secret := uuid.NewString()
			t.Setenv(envVarName, iam_client_secret)

			credentials := createCredentials(&config.ConfigGetter{Provider: provider})
			if credentials == nil {
				return
			}
			require.Equal(t, iam_client_secret, credentials.IAM.ClientSecret)
			assert.Equal(t, iam_client_secret, credentials.OAuth.ClientSecret)
		})
	}
}
func TestSSOTokenURL(t *testing.T) {
	if len(os.Getenv("TF_ACC")) > 0 {
		return
	}
	provider := provider.Provider()

	for _, envURL := range []string{
		"https://foo.live.dynatrace.com",
		"https://foo.apps.dynatrace.com",
		"https://foo.live.dynatrace.com ",
		"https://foo.apps.dynatrace.com ",
		"https://foo.live.dynatrace.com/",
		"https://foo.apps.dynatrace.com/",
		"https://foo.live.dynatrace.com/ ",
		"https://foo.apps.dynatrace.com/ ",
	} {
		t.Run(envURL, func(t *testing.T) {
			t.Setenv("DYNATRACE_ENV_URL", envURL)

			credentials := createCredentials(&config.ConfigGetter{Provider: provider})
			if credentials == nil {
				return
			}
			require.Equal(t, rest.ProdTokenURL, credentials.IAM.TokenURL)
			require.Equal(t, rest.ProdIAMEndpointURL, credentials.IAM.EndpointURL)
			assert.Equal(t, rest.ProdTokenURL, credentials.OAuth.TokenURL)
		})
	}
	for _, envURL := range []string{
		"https://foo.sprint.dynatracelabs.com",
		"https://foo.sprint.apps.dynatracelabs.com",
		"https://foo.sprint.dynatracelabs.com ",
		"https://foo.sprint.apps.dynatracelabs.com ",
		"https://foo.sprint.dynatracelabs.com/",
		"https://foo.sprint.apps.dynatracelabs.com/",
		"https://foo.sprint.dynatracelabs.com/ ",
		"https://foo.sprint.apps.dynatracelabs.com/ ",
	} {
		t.Run(envURL, func(t *testing.T) {
			t.Setenv("DYNATRACE_ENV_URL", envURL)

			credentials := createCredentials(&config.ConfigGetter{Provider: provider})
			if credentials == nil {
				return
			}
			require.Equal(t, rest.SprintTokenURL, credentials.IAM.TokenURL)
			require.Equal(t, rest.SprintIAMEndpointURL, credentials.IAM.EndpointURL)
			assert.Equal(t, rest.SprintTokenURL, credentials.OAuth.TokenURL)
		})
	}
	for _, envURL := range []string{
		"https://foo.dev.dynatracelabs.com",
		"https://foo.dev.apps.dynatracelabs.com",
		"https://foo.dev.dynatracelabs.com/",
		"https://foo.dev.apps.dynatracelabs.com/",
		"https://foo.dev.dynatracelabs.com ",
		"https://foo.dev.apps.dynatracelabs.com ",
		"https://foo.dev.dynatracelabs.com/ ",
		"https://foo.dev.apps.dynatracelabs.com/ ",
	} {
		t.Run(envURL, func(t *testing.T) {
			t.Setenv("DYNATRACE_ENV_URL", envURL)

			credentials := createCredentials(&config.ConfigGetter{Provider: provider})
			if credentials == nil {
				return
			}
			require.Equal(t, rest.DevTokenURL, credentials.IAM.TokenURL)
			require.Equal(t, rest.DevIAMEndpointURL, credentials.IAM.EndpointURL)
			assert.Equal(t, rest.DevTokenURL, credentials.OAuth.TokenURL)
		})
	}
}

func createCredentials(getter config.Getter) *rest.Credentials {
	configResult, _ := config.ProviderConfigureGeneric(context.Background(), getter)
	creds, _ := config.Credentials(configResult, config.CredValNone)
	return creds
}
