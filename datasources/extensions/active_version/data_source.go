/*
 * @license
 * Copyright 2026 Dynatrace LLC
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package active_version

import (
	"context"
	"encoding/json"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/provider/config"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/provider/logging"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	coreapi "github.com/dynatrace/dynatrace-configuration-as-code-core/api"
	coreextensions "github.com/dynatrace/dynatrace-configuration-as-code-core/clients/extensions"
)

func DataSource() *schema.Resource {
	return &schema.Resource{
		ReadContext: logging.EnableDSCtx(dataSourceRead),
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Description: "The name of the extension",
				Required:    true,
			},
			"active_version": {
				Type:        schema.TypeString,
				Description: "The active version of the extension",
				Computed:    true,
			},
		},
	}
}

// ExtensionClient defines the interface for interacting with the Extensions API.
type ExtensionClient interface {
	// GetEnvironmentConfiguration returns the active version of an extension in an environment
	GetEnvironmentConfiguration(ctx context.Context, extensionName string) (coreapi.Response, error)
}

func createCoreClient(ctx context.Context, credentials *rest.Credentials) (ExtensionClient, error) {
	platformClient, err := rest.CreatePlatformClient(ctx, credentials.OAuth.EnvironmentURL, credentials)
	if err != nil {
		return nil, err
	}
	return coreextensions.NewClient(platformClient), nil
}

func dataSourceRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	creds, err := config.Credentials(m, config.CredValDefault)
	if err != nil {
		return diag.FromErr(err)
	}

	client, err := createCoreClient(ctx, creds)
	if err != nil {
		return diag.FromErr(err)
	}

	return DataSourceReadWithClient(ctx, d, client)
}

type environmentConfiguration struct {
	Version string `json:"version"`
}

func DataSourceReadWithClient(ctx context.Context, d *schema.ResourceData, client ExtensionClient) diag.Diagnostics {
	extName := d.Get("name").(string)

	response, err := client.GetEnvironmentConfiguration(ctx, extName)
	if err != nil {
		return diag.FromErr(err)
	}

	var conf environmentConfiguration
	err = json.Unmarshal(response.Data, &conf)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(extName)
	err = d.Set("active_version", conf.Version)
	if err != nil {
		return diag.FromErr(err)
	}
	return diag.Diagnostics{}
}
