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

package apitoken

import (
	"context"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/provider/config"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/provider/logging"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	srv "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v2/apitokens"
	apitokens "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v2/apitokens/settings"
)

func DataSource() *schema.Resource {
	return &schema.Resource{
		ReadContext: logging.EnableDSCtx(DataSourceRead),
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"enabled": {
				Type:        schema.TypeBool,
				Description: "The token is enabled (true) or disabled (false), default disabled (false).",
				Computed:    true,
			},
			"personal_access_token": {
				Type:        schema.TypeBool,
				Description: "The token is a personal access token (true) or an API token (false).",
				Computed:    true,
			},
			"expiration_date": {
				Type:        schema.TypeString,
				Description: "The expiration date of the token.",
				Computed:    true,
			},
			"owner": {
				Type:        schema.TypeString,
				Description: "The owner of the token",
				Computed:    true,
			},
			"creation_date": {
				Type:        schema.TypeString,
				Description: "Token creation date in ISO 8601 format (yyyy-MM-dd'T'HH:mm:ss.SSS'Z')",
				Computed:    true,
			},
			"scopes": {
				Type:        schema.TypeSet,
				Description: "A list of the scopes to be assigned to the token.",
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func DataSourceRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	var name string
	if v, ok := d.GetOk("name"); ok {
		name = v.(string)
	}
	creds, err := config.Credentials(m, config.CredValDefault)
	if err != nil {
		return diag.FromErr(err)
	}

	service := srv.Service(creds)
	var stubs api.Stubs
	if stubs, err = service.List(ctx); err != nil {
		return diag.FromErr(err)
	}

	if len(stubs) > 0 {
		for _, stub := range stubs {
			if name == stub.Name {
				d.SetId(stub.ID)

				value := stub.Value.(*apitokens.APIToken)
				d.Set("enabled", value.Enabled)
				d.Set("personal_access_token", value.PersonalAccessToken)
				d.Set("expiration_date", value.ExpirationDate)
				d.Set("owner", value.Owner)
				d.Set("creation_date", value.CreationDate)
				d.Set("scopes", value.Scopes)

				return diag.Diagnostics{}
			}
		}
	}
	d.SetId("")
	return diag.Diagnostics{}
}
