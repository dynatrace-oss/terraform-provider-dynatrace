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
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const platformIDPrefix = "CREDENTIAL_VAULT_ENTRY-"
const classicIDPrefix = "CREDENTIALS_VAULT-"

// ToClassicID converts a platform credential ID (`CREDENTIAL_VAULT_ENTRY-…`) into its classic form (`CREDENTIALS_VAULT-…`).
func ToClassicID(id string) string {
	return replaceIDPrefix(id, platformIDPrefix, classicIDPrefix)
}

// ToPlatformID converts a classic credential ID (`CREDENTIALS_VAULT-…`) into its platform form (`CREDENTIAL_VAULT_ENTRY-…`).
func ToPlatformID(id string) string {
	return replaceIDPrefix(id, classicIDPrefix, platformIDPrefix)
}

func replaceIDPrefix(id string, from string, to string) string {
	if suffix, found := strings.CutPrefix(id, from); found {
		return to + suffix
	}
	return id
}

// suppressEquivalentCredentialID treats the classic and the platform form of the same credential ID as equal,
// because the configuration may reference either a `dynatrace_credentials` or a `dynatrace_credentials_v2` resource.
func suppressEquivalentCredentialID(_ string, oldValue string, newValue string, _ *schema.ResourceData) bool {
	return ToClassicID(oldValue) == ToClassicID(newValue)
}
