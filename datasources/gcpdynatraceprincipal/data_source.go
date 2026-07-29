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

// Package dynatraceprincipal exposes the data source
// `dynatrace_gcp_dynatrace_principal`, which surfaces the read-only Settings
// 2.0 object exposing Dynatrace's own GCP service account
// (`dynatrace-<tenant-id>@dtp-prod-gcp-auth.iam.gserviceaccount.com`).
//
// If the principal object does not yet exist on the tenant, the data source
// creates it (mirroring dtctl's `EnsureDynatracePrincipalWithResult`) and
// re-reads to capture the materialised principal email.
package dynatraceprincipal

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/provider/config"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/provider/logging"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const PrincipalSchemaID = "builtin:hyperscaler-authentication.connections.gcp-dynatrace-principal"

func DataSource() *schema.Resource {
	return &schema.Resource{
		ReadContext: logging.EnableDSCtx(DataSourceRead),
		Schema: map[string]*schema.Schema{
			"principal": {
				Type:        schema.TypeString,
				Description: "Dynatrace-side GCP service-account email (e.g. `dynatrace-<tenant-id>@dtp-prod-gcp-auth.iam.gserviceaccount.com`). Use this to wire `roles/iam.serviceAccountTokenCreator` on the customer SA.",
				Computed:    true,
			},
			"object_id": {
				Type:        schema.TypeString,
				Description: "Settings 2.0 objectId of the principal entry.",
				Computed:    true,
			},
		},
	}
}

type settingsListResponse struct {
	Items []struct {
		ObjectID string          `json:"objectId"`
		Value    json.RawMessage `json:"value"`
	} `json:"items"`
}

type principalValue struct {
	Principal string `json:"principal"`
}

func DataSourceRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	creds, err := config.Credentials(m, config.CredValDefault)
	if err != nil {
		return diag.FromErr(err)
	}
	client := rest.HybridClient(creds)

	principal, objectID, err := findPrincipal(ctx, client)
	if err != nil {
		return diag.FromErr(err)
	}
	if principal == "" {
		if err := createPrincipal(ctx, client); err != nil {
			return diag.FromErr(err)
		}
		principal, objectID, err = findPrincipal(ctx, client)
		if err != nil {
			return diag.FromErr(err)
		}
	}
	if principal == "" {
		return diag.Errorf("dynatrace_gcp_dynatrace_principal: tenant did not return a principal email after ensuring schema %s", PrincipalSchemaID)
	}

	d.SetId(objectID)
	_ = d.Set("principal", principal)
	_ = d.Set("object_id", objectID)
	return diag.Diagnostics{}
}

func findPrincipal(ctx context.Context, client rest.Client) (string, string, error) {
	url := fmt.Sprintf("/api/v2/settings/objects?schemaIds=%s&fields=objectId,value&pageSize=100", PrincipalSchemaID)
	var out settingsListResponse
	if err := client.Get(ctx, url).Expect(200).Finish(&out); err != nil {
		return "", "", fmt.Errorf("failed to list %s: %w", PrincipalSchemaID, err)
	}
	if len(out.Items) == 0 {
		return "", "", nil
	}
	var v principalValue
	if err := json.Unmarshal(out.Items[0].Value, &v); err != nil {
		return "", "", fmt.Errorf("failed to parse principal value: %w", err)
	}
	return v.Principal, out.Items[0].ObjectID, nil
}

func createPrincipal(ctx context.Context, client rest.Client) error {
	body := []map[string]any{
		{
			"schemaId": PrincipalSchemaID,
			"value":    map[string]any{},
		},
	}
	if err := client.Post(ctx, "/api/v2/settings/objects", body).Expect(200, 201).Finish(); err != nil {
		return fmt.Errorf("failed to create %s: %w", PrincipalSchemaID, err)
	}
	return nil
}
