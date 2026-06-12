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
	"fmt"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Credential links the Azure monitoring configuration to a HAS connection
// and pins the Azure Service Principal that Dynatrace will impersonate.
type Credential struct {
	ConnectionID       string
	ServicePrincipalID string
	Type               string // FEDERATED | SECRET
	Description        string
	Enabled            bool
}

type Credentials []*Credential

func (me *Credential) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"connection_id": {
			Type:        schema.TypeString,
			Description: "ObjectId of the `dynatrace_azure_connection` resource.",
			Required:    true,
		},
		"service_principal_id": {
			Type:        schema.TypeString,
			Description: "Azure application (client) id of the Service Principal that Dynatrace impersonates. Typically `azuread_application.<name>.client_id`.",
			Required:    true,
		},
		"type": {
			Type:        schema.TypeString,
			Description: "`FEDERATED` (workload-identity federation) or `SECRET` (client secret). Defaults to `FEDERATED` when omitted.",
			Optional:    true,
			Computed:    true,
			ValidateFunc: func(i any, k string) (warnings []string, errs []error) {
				v, _ := i.(string)
				if v != "" && v != "FEDERATED" && v != "SECRET" {
					errs = append(errs, fmt.Errorf("%s must be FEDERATED or SECRET, got %q", k, v))
				}
				return
			},
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
		"connection_id":        me.ConnectionID,
		"service_principal_id": me.ServicePrincipalID,
		"type":                 me.Type,
		"description":          me.Description,
		"enabled":              me.Enabled,
	})
}

func (me *Credential) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"connection_id":        &me.ConnectionID,
		"service_principal_id": &me.ServicePrincipalID,
		"type":                 &me.Type,
		"description":          &me.Description,
		"enabled":              &me.Enabled,
	})
}
