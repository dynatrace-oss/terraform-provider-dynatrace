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

package settings

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Credential links the GCP monitoring configuration to a HAS connection and
// pins the customer service account that Dynatrace will impersonate.
type Credential struct {
	ConnectionID   string
	ServiceAccount string
	Description    string
	Enabled        bool
}

type Credentials []*Credential

func (me *Credential) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"connection_id": {
			Type:        schema.TypeString,
			Description: "ObjectId of the `dynatrace_gcp_connection` resource (or the manually-created connection).",
			Required:    true,
		},
		"service_account": {
			Type:        schema.TypeString,
			Description: "The customer GCP service account email Dynatrace impersonates (e.g. `dynatrace-integration@<project>.iam.gserviceaccount.com`).",
			Required:    true,
		},
		"description": {
			Type:        schema.TypeString,
			Description: "Free-form description for this credential. Defaults to the top-level `name`.",
			Optional:    true,
			Computed:    true,
		},
		"enabled": {
			Type:        schema.TypeBool,
			Description: "Per-credential enable flag. Defaults to true. Distinct from the top-level `enabled`.",
			Optional:    true,
			Default:     true,
		},
	}
}

func (me *Credential) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"connection_id":   me.ConnectionID,
		"service_account": me.ServiceAccount,
		"description":     me.Description,
		"enabled":         me.Enabled,
	})
}

func (me *Credential) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"connection_id":   &me.ConnectionID,
		"service_account": &me.ServiceAccount,
		"description":     &me.Description,
		"enabled":         &me.Enabled,
	})
}
