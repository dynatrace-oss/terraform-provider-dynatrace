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

package gcpdynatraceprincipal

import (
	"context"
	"fmt"

	gcpdynatraceprincipal "github.com/dynatrace-oss/terraform-provider-dynatrace/datasources/hyperscalerauthentication/connections/gcpdynatraceprincipal/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings/services/settings20"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/provider/config"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/provider/logging"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const SchemaVersion = "0.0.3"
const SchemaID = "builtin:hyperscaler-authentication.connections.gcp-dynatrace-principal"

func ensurePrincipalExists(ctx context.Context, credentials *rest.Credentials) (string, error) {
	service := settings20.Service[*gcpdynatraceprincipal.Settings](credentials, SchemaID, SchemaVersion)

	stubs, err := service.List(ctx)
	if err != nil {
		return "", err
	}

	if len(stubs) == 0 {
		return "", fmt.Errorf("no Dynatrace GCP Principal found")
		// TODO: ensure principal exists
	}

	// There can only be one principal, so we take the first (and only one) in the list.
	v := stubs[0].Value.(*gcpdynatraceprincipal.Settings)
	return v.Principal, nil
}

// Attribute keys
const (
	attrPrincipal = "principal"
)

func DataSource() *schema.Resource {
	return &schema.Resource{
		ReadContext: logging.EnableDSCtx(DataSourceRead),
		Description: "Returns the Dynatrace GCP Principal",
		Schema: map[string]*schema.Schema{
			attrPrincipal: {
				Type:        schema.TypeString,
				Description: "Dynatrace GCP Principal",
				Computed:    true,
			},
		},
	}
}

func DataSourceRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	creds, err := config.Credentials(m, config.CredValDefault)
	if err != nil {
		return diag.FromErr(err)
	}

	principal, err := ensurePrincipalExists(ctx, creds)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(principal)
	if err := d.Set(attrPrincipal, principal); err != nil {
		return diag.FromErr(err)
	}

	return diag.Diagnostics{}
}
