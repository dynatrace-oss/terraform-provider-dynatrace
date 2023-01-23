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

package vault

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings"
)

// CredentialsList A list of credentials sets for Synthetic monitors.
type CredentialsList struct {
	Credentials []CredentialsResponseElement `json:"credentials"` // A list of credentials sets for Synthetic monitors.
}

func (me CredentialsList) ToStubs() settings.Stubs {
	stubs := settings.Stubs{}
	for _, elem := range me.Credentials {
		stubs = append(stubs, &settings.Stub{ID: *elem.ID, Name: elem.Name})
	}
	return stubs
}

// CredentialsResponseElement Metadata of the credentials set.
type CredentialsResponseElement struct {
	Name            string                         `json:"name"`            // The name of the credentials set.
	ID              *string                        `json:"id,omitempty"`    // The ID of the credentials set.
	Description     string                         `json:"description"`     // A short description of the credentials set.
	Owner           string                         `json:"owner"`           // The owner of the credential.
	OwnerAccessOnly bool                           `json:"ownerAccessOnly"` // Flag indicating that this credential is visible only to the owner.
	Type            CredentialsResponseElementType `json:"type"`            // The type of the credentials set.
	Scope           Scope                          `json:"scope,omitempty"` // The scope of the credentials set
}
