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

package credentialsv2

import (
	"context"
	"encoding/json"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/provider/config"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/provider/logging"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/dynatrace/dynatrace-configuration-as-code-core/clients/credentialvault"
)

func DataSource() *schema.Resource {
	return &schema.Resource{
		Description: "Retrieve credentials stored within the Credential Vault.",
		ReadContext: logging.EnableDSCtx(DataSourceRead),
		Schema: map[string]*schema.Schema{
			"type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The type of the credential. Possible values are `CERTIFICATE`, `PUBLIC_CERTIFICATE`, `SNMPV3`, `TOKEN`, `USERNAME_PASSWORD`, `AWS_MONITORING_KEY_BASED`, `AWS_MONITORING_ROLE_BASED`. If not specified all credential types will match",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The name of the credential as shown within the Dynatrace WebUI, matched case-insensitively. If not specified all names will match",
			},
			"scope": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The scope of the credential. Possible values are `APP_ENGINE`, `EXTENSION`, `EXTENSION_AUTHENTICATION` and `SYNTHETIC`. If not specified all scopes will match.",
			},
			"user": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Filters the result to credentials accessible to the specified user (those they own or those shared with everyone). Entries owned by other users are listed only with `credential-vault:entries:admin`",
			},
			"items": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The credentials matching the filters. As credentials aren't unique by name, scope, user or type, multiple credentials can be returned",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of the credential",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the credential",
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The type of the credential",
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The description of the credential",
						},
						"scopes": {
							Type:        schema.TypeList,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "The scopes of the credential",
						},
						"owner_access_only": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "The credential is available to the owner only (`true`) or to every user (`false`)",
						},
					},
				},
			},
		},
	}
}

type entry struct {
	ID              string   `json:"id"`
	Name            string   `json:"name"`
	Type            string   `json:"type"`
	Description     string   `json:"description"`
	Scopes          []string `json:"scopes"`
	OwnerAccessOnly bool     `json:"ownerAccessOnly"`
}

func DataSourceRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	clientSet, err := config.ClientSet(m, config.CredValPlatform)
	if err != nil {
		return diag.FromErr(err)
	}
	return dataSourceRead(ctx, d, clientSet)
}

func dataSourceRead(ctx context.Context, d *schema.ResourceData, clientSet rest.ClientSet) diag.Diagnostics {
	options := credentialvault.ListOptions{
		Type:  d.Get("type").(string),
		Scope: d.Get("scope").(string),
		User:  d.Get("user").(string),
		Name:  d.Get("name").(string),
	}

	platformClient, err := clientSet.PlatformClient()
	if err != nil {
		return diag.FromErr(err)
	}

	result, err := credentialvault.NewClient(platformClient).List(ctx, options)
	if err != nil {
		return diag.FromErr(err)
	}

	items := make([]any, 0, len(result.All()))
	for _, raw := range result.All() {
		var credential entry
		if err = json.Unmarshal(raw, &credential); err != nil {
			return diag.FromErr(err)
		}
		items = append(items, map[string]any{
			"id":                credential.ID,
			"name":              credential.Name,
			"type":              credential.Type,
			"description":       credential.Description,
			"scopes":            credential.Scopes,
			"owner_access_only": credential.OwnerAccessOnly,
		})
	}

	if err = d.Set("items", items); err != nil {
		return diag.FromErr(err)
	}
	d.SetId("credentials-v2")
	return diag.Diagnostics{}
}
