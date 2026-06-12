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

// Package settings models the typed Terraform resource
// `dynatrace_gcp_connection`, mapping HCL attributes onto the Settings 2.0
// schema `builtin:hyperscaler-authentication.connections.gcp`.
package settings

import (
	"time"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ConnectionType is always `serviceAccountImpersonation` — the wire schema
// rejects anything else for this provider.
const ConnectionType = "serviceAccountImpersonation"

// DefaultConsumer matches what dtctl writes when the user does not pass an
// explicit consumer list.
const DefaultConsumer = "SVC:com.dynatrace.da"

const DefaultTimeout = 2 * time.Minute

// ServiceAccountImpersonation is the nested wire object that carries the
// customer service account email plus the consumers allowed to mint tokens
// for it.
type ServiceAccountImpersonation struct {
	ServiceAccountID string   `json:"serviceAccountId,omitempty"`
	Consumers        []string `json:"consumers"`
}

// Settings is the typed model exposed to Terraform. JSON tags mirror the
// `value.*` shape sent to /api/v2/settings/objects.
type Settings struct {
	Name                        string                       `json:"name"`
	Type                        string                       `json:"type"`
	ServiceAccountImpersonation *ServiceAccountImpersonation `json:"serviceAccountImpersonation,omitempty"`
}

func (me *Settings) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:        schema.TypeString,
			Description: "The name of the connection.",
			Required:    true,
		},
		"service_account_id": {
			Type:        schema.TypeString,
			Description: "Customer GCP service account email (e.g. `dynatrace-integration@<project>.iam.gserviceaccount.com`). Optional at create time — dtctl allows patching it later via `dtctl update gcp connection --serviceAccountId …`; Terraform users normally set it directly.",
			Optional:    true,
			Computed:    true,
		},
		"oauth_consumers": {
			Type:        schema.TypeSet,
			Description: "Consumers that can use the connection. Defaults to `[\"SVC:com.dynatrace.da\"]`.",
			Optional:    true,
			Computed:    true,
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
	}
}

func (me *Settings) Timeouts() *schema.ResourceTimeout {
	return &schema.ResourceTimeout{
		Create: schema.DefaultTimeout(DefaultTimeout),
		Update: schema.DefaultTimeout(DefaultTimeout),
	}
}

func (me *Settings) MarshalHCL(properties hcl.Properties) error {
	out := map[string]any{
		"name": me.Name,
	}
	if me.ServiceAccountImpersonation != nil {
		out["service_account_id"] = me.ServiceAccountImpersonation.ServiceAccountID
		out["oauth_consumers"] = me.ServiceAccountImpersonation.Consumers
	}
	return properties.EncodeAll(out)
}

func (me *Settings) UnmarshalHCL(decoder hcl.Decoder) error {
	var serviceAccountID string
	var consumers []string
	if err := decoder.DecodeAll(map[string]any{
		"name":               &me.Name,
		"service_account_id": &serviceAccountID,
		"oauth_consumers":    &consumers,
	}); err != nil {
		return err
	}
	me.Type = ConnectionType
	if len(consumers) == 0 {
		consumers = []string{DefaultConsumer}
	}
	me.ServiceAccountImpersonation = &ServiceAccountImpersonation{
		ServiceAccountID: serviceAccountID,
		Consumers:        consumers,
	}
	return nil
}
