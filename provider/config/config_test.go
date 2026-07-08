//go:build unit

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

package config_test

import (
	"testing"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/provider/config"
	"github.com/stretchr/testify/assert"
)

type mockResourceData map[string]any

func (mrd mockResourceData) Get(k string) any {
	return mrd[k]
}

func TestProviderConfigure(t *testing.T) {
	d := mockResourceData{
		"dt_env_url":   "https://something.live.dynatrace.com",
		"dt_api_token": "faketoken",
	}

	result, _ := config.ProviderConfigureGeneric(t.Context(), d)
	configuration := result.(*config.ProviderConfiguration)
	assert.Equal(t, "https://something.live.dynatrace.com", configuration.EnvironmentURL)

}
