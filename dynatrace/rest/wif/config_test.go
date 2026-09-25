//go:build unit

/*
 * @license
 * Copyright 2026 Dynatrace LLC
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

package wif

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// The audience is the only value the configuration contributes, so asserting on the whole config is
// what shows it reaches the vendor it was inferred for.
func TestInferVendorConfigDetectsGitHub(t *testing.T) {
	injectGitHubCredentials(t, "https://token.service.invalid/", "request-token")

	config, err := inferVendorConfig("dynatrace")

	require.NoError(t, err)
	assert.Equal(t, gitHubConfig{
		audience:          "dynatrace",
		tokenRequestURL:   "https://token.service.invalid/",
		tokenRequestToken: "request-token",
	}, config)
}

func TestInferVendorConfigReportsThatNoVendorWasDetected(t *testing.T) {
	injectGitHubCredentials(t, "", "")

	_, err := inferVendorConfig("dynatrace")

	assert.ErrorIs(t, err, errNoVendorDetected)
}
