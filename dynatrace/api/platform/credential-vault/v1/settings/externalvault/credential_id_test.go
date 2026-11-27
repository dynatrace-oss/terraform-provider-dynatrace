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

package externalvault

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToClassicID_ConvertsPlatformID(t *testing.T) {
	assert.Equal(t, "CREDENTIALS_VAULT-56405E761B90CAEA", ToClassicID("CREDENTIAL_VAULT_ENTRY-56405E761B90CAEA"))
}

func TestToClassicID_KeepsClassicID(t *testing.T) {
	assert.Equal(t, "CREDENTIALS_VAULT-56405E761B90CAEA", ToClassicID("CREDENTIALS_VAULT-56405E761B90CAEA"))
}

func TestToPlatformID_ConvertsClassicID(t *testing.T) {
	assert.Equal(t, "CREDENTIAL_VAULT_ENTRY-56405E761B90CAEA", ToPlatformID("CREDENTIALS_VAULT-56405E761B90CAEA"))
}

func TestToPlatformID_KeepsPlatformID(t *testing.T) {
	assert.Equal(t, "CREDENTIAL_VAULT_ENTRY-56405E761B90CAEA", ToPlatformID("CREDENTIAL_VAULT_ENTRY-56405E761B90CAEA"))
}

func TestToClassicID_OnlyReplacesPrefix(t *testing.T) {
	assert.Equal(t, "my-CREDENTIAL_VAULT_ENTRY-id", ToClassicID("my-CREDENTIAL_VAULT_ENTRY-id"))
}
