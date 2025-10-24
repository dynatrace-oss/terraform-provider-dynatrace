/*
 * @license
 * Copyright 2025 Dynatrace LLC
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package featureflags_test

import (
	"fmt"
	"testing"

	featureflags "github.com/dynatrace-oss/terraform-provider-dynatrace/provider/featureflag"
	"github.com/stretchr/testify/assert"
)

func TestFeatureFlagEnabledOnTrueValues(t *testing.T) {
	ff := featureflags.FeatureFlag{Name: "env", DefaultValue: false}

	for _, v := range []string{"1", "t", "T", "TRUE", "true", "tRuE", "True"} {
		t.Run(fmt.Sprintf("'%s' should be true", v), func(t *testing.T) {
			t.Setenv(ff.Name, v)
			assert.True(t, ff.Enabled())
		})
	}
}

func TestFeatureFlagDisabledOnFalseValues(t *testing.T) {
	ff := featureflags.FeatureFlag{Name: "env", DefaultValue: true}

	for _, v := range []string{"0", "f", "F", "FALSE", "false", "False", "fAlSe"} {
		t.Run(fmt.Sprintf("'%s' should be false", v), func(t *testing.T) {
			t.Setenv(ff.Name, v)
			assert.False(t, ff.Enabled())
		})
	}
}

func TestFeatureFlagDefaultValues(t *testing.T) {
	ff := featureflags.FeatureFlag{Name: "env", DefaultValue: true}

	t.Run("Uses the default value if the env is not set", func(t *testing.T) {
		assert.True(t, ff.Enabled())
	})

	t.Run("Uses the default value if the env is not boolean", func(t *testing.T) {
		t.Setenv(ff.Name, "invalidBool")
		assert.True(t, ff.Enabled())
	})
}

func TestFeatureFlagString(t *testing.T) {
	ff := featureflags.FeatureFlag{Name: "env", DefaultValue: true}
	assert.Equal(t, "env", ff.String())
}
