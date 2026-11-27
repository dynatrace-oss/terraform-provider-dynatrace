//go:build unit

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
	"testing"

	settings "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/platform/synthetic/configuration/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func decode(t *testing.T, raw map[string]any) *settings.Settings {
	t.Helper()
	s := new(settings.Settings)
	data := schema.TestResourceDataRaw(t, s.Schema(), raw)
	require.NoError(t, s.UnmarshalHCL(hcl.DecoderFrom(data)))
	return s
}

func TestUnmarshalHCL(t *testing.T) {
	t.Run("Decodes all fields", func(t *testing.T) {
		result := decode(t, map[string]any{
			"bm_monitor_timeout": 120000,
			"bm_step_timeout":    45000,
			"maintenance_windows_excluded_from_availability": true,
		})

		require.NotNil(t, result.BmMonitorTimeout)
		require.NotNil(t, result.BmStepTimeout)
		require.NotNil(t, result.MaintenanceWindowsExcludedFromAvailability)
		assert.Equal(t, int64(120000), *result.BmMonitorTimeout)
		assert.Equal(t, int64(45000), *result.BmStepTimeout)
		assert.True(t, *result.MaintenanceWindowsExcludedFromAvailability)
	})

	t.Run("Decodes an explicit false", func(t *testing.T) {
		result := decode(t, map[string]any{"maintenance_windows_excluded_from_availability": false})

		require.NotNil(t, result.MaintenanceWindowsExcludedFromAvailability)
		assert.False(t, *result.MaintenanceWindowsExcludedFromAvailability)
	})

	t.Run("Leaves unset fields nil", func(t *testing.T) {
		result := decode(t, map[string]any{})

		assert.Nil(t, result.BmMonitorTimeout)
		assert.Nil(t, result.BmStepTimeout)
		assert.Nil(t, result.MaintenanceWindowsExcludedFromAvailability)
	})
}
