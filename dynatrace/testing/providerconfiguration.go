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

package testing

import (
	"os"
	"testing"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/provider/config"
	"github.com/stretchr/testify/require"
)

var _ config.Getter = &environmentGetter{}

type environmentGetter struct {
	envURL   string
	apiToken string
}

// EnvironmentGetter returns a config.Getter that retrieves the Dynatrace environment URL and API token from environment variables.
func EnvironmentGetter(t *testing.T) *environmentGetter {
	envURL := os.Getenv("DYNATRACE_ENV_URL")
	require.NotEmpty(t, envURL, "Environment variable DYNATRACE_ENV_URL must be specified")

	apiToken := os.Getenv("DYNATRACE_API_TOKEN")
	require.NotEmpty(t, apiToken, "Environment variable DYNATRACE_API_TOKEN must be specified")

	return &environmentGetter{
		envURL:   envURL,
		apiToken: apiToken,
	}
}

func (t *environmentGetter) Get(key string) any {
	switch key {
	case "dt_env_url":
		return t.envURL
	case "dt_api_token":
		return t.apiToken
	}
	return ""
}
