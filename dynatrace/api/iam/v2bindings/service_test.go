//go:build integration

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

package v2bindings_test

import (
	"net/url"
	"os"
	"strings"
	"testing"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/testing/api"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/provider/envutils"
	"github.com/stretchr/testify/require"
)

func TestAccTestCasesV2Bindings(t *testing.T) {
	setAccountEnv(t)
	api.TestAccTestCases(t)
}

func TestAccV2Bindings(t *testing.T) {
	setAccountEnv(t)
	api.TestAcc(t)
}

func setAccountEnv(t *testing.T) {
	accountID := os.Getenv("DT_ACCOUNT_ID")
	//fallback to DYNATRACE_ACCOUNT_ID
	if accountID == "" {
		accountID = os.Getenv("DYNATRACE_ACCOUNT_ID")
	}
	t.Setenv("TF_VAR_ENVIRONMENT_ID", getEnvIdFromURL(t, envutils.DynatraceEnvURL.Get()))
	t.Setenv("TF_VAR_ACCOUNT_ID", accountID)
}

func getEnvIdFromURL(t *testing.T, envURL string) string {
	t.Helper()

	u, err := url.Parse(envURL)
	require.NoError(t, err)

	host := u.Hostname() // strips scheme, port, trailing slash
	envId, _, found := strings.Cut(host, ".")
	require.True(t, found)
	require.NotEmpty(t, envId)

	return envId
}
