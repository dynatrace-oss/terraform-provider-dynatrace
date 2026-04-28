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

package latest_version

import (
	"context"
	"encoding/json"

	"github.com/Masterminds/semver/v3"
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
			"latest_version": {
				Type:        schema.TypeString,
				Description: "The latest installed version of the extension",
				Computed:    true,
			},
		},
	}
}

// ExtensionClient defines the interface for interacting with the Extensions API.
type ExtensionClient interface {
	// ListExtensionVersions returns all installed versions of an extension
	ListExtensionVersions(ctx context.Context, extensionName string) (coreapi.PagedListResponse, error)
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

type listExtensionVersions struct {
	Version string `json:"version"`
}

func DataSourceReadWithClient(ctx context.Context, d *schema.ResourceData, client ExtensionClient) diag.Diagnostics {
	extName := d.Get("name").(string)

	response, err := client.ListExtensionVersions(ctx, extName)
	if err != nil {
		return diag.FromErr(err)
	}

	var latestVersion *semver.Version
	for _, body := range response.All() {
		var extensionInfo listExtensionVersions
		if err := json.Unmarshal(body, &extensionInfo); err != nil {
			return diag.FromErr(err)
		}
		v, err := semver.NewVersion(extensionInfo.Version)
		if err != nil {
			// skip entries that are not valid semver
			continue
		}
		if latestVersion == nil || v.GreaterThan(latestVersion) {
			latestVersion = v
		}
	}
	if latestVersion == nil {
		d.SetId("")
		return nil
	}

	d.SetId(extName)
	err = d.Set("latest_version", latestVersion.String())
	if err != nil {
		return diag.FromErr(err)
	}
	return diag.Diagnostics{}
}
