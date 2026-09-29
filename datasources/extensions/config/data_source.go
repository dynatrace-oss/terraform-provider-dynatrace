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

package extension_config

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/provider/config"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/provider/logging"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	coreapi "github.com/dynatrace/dynatrace-configuration-as-code-core/api"
	coreextensions "github.com/dynatrace/dynatrace-configuration-as-code-core/clients/extensions"

	extensionconfig "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/extensions/monitoringconfigurations/settings"
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
			"description": {
				Type:        schema.TypeString,
				Description: "The description of the configuration",
				Optional:    true,
			},
			"version": {
				Type:        schema.TypeString,
				Description: "The extension version of the configuration",
				Optional:    true,
			},
			"activation_context": {
				Type:        schema.TypeString,
				Description: "The activation context of the configuration",
				Optional:    true,
			},
			"items": {
				Type:        schema.TypeList,
				Description: "An array of the found configurations",
				Computed:    true,
				Elem:        &schema.Resource{Schema: itemSchema()},
			},
		},
	}
}

func itemSchema() map[string]*schema.Schema {
	scm := hcl.SetComputedSchema(new(extensionconfig.Settings).Schema())
	delete(scm, "name") // remove the foreign-key (extension name)
	scm["id"] = &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "The ID of the configuration",
	}
	return scm
}

// ExtensionClient defines the interface for interacting with the Extensions API.
type ExtensionClient interface {
	// ListMonitoringConfigurations returns all configurations of an extension
	ListMonitoringConfigurations(ctx context.Context, extensionName string, filter string) (coreapi.PagedListResponse, error)
}

func dataSourceRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	clientSet, err := config.ClientSet(m, config.CredValDefault)
	if err != nil {
		return diag.FromErr(err)
	}

	platformClient, err := clientSet.PlatformClient()
	if err != nil {
		return diag.FromErr(err)
	}

	return DataSourceReadWithClient(ctx, d, coreextensions.NewClient(platformClient))
}

func DataSourceReadWithClient(ctx context.Context, d *schema.ResourceData, client ExtensionClient) diag.Diagnostics {
	extName := d.Get("name").(string)
	filter := getFilter(d)

	response, err := client.ListMonitoringConfigurations(ctx, extName, filter)
	if err != nil {
		return diag.FromErr(err)
	}

	responses := response.All()
	items := make([]any, 0, len(responses))

	type idResponse struct {
		ID string `json:"objectId"`
	}

	for _, cfg := range responses {
		extConf := extensionconfig.Settings{}
		if err = json.Unmarshal(cfg, &extConf); err != nil {
			return diag.FromErr(err)
		}
		var idResp idResponse
		if err = json.Unmarshal(cfg, &idResp); err != nil {
			return diag.FromErr(err)
		}

		item := hcl.Properties{}
		if err = extConf.MarshalHCL(item); err != nil {
			return diag.FromErr(err)
		}
		// remove the foreign-key (extension name)
		delete(item, "name")

		item["id"] = idResp.ID
		items = append(items, item)
	}

	d.SetId(extName)
	if err = d.Set("items", items); err != nil {
		return diag.FromErr(err)
	}
	return diag.Diagnostics{}
}

func getFilter(d *schema.ResourceData) string {
	filters := make([]string, 0, 3)

	if val, ok := d.GetOk("version"); ok {
		version := val.(string)
		filters = append(filters, fmt.Sprintf("version='%s'", version))
	}
	if val, ok := d.GetOk("description"); ok {
		description := val.(string)
		filters = append(filters, fmt.Sprintf("description='%s'", description))
	}
	if val, ok := d.GetOk("activation_context"); ok {
		activationContext := val.(string)
		filters = append(filters, fmt.Sprintf("activationContext='%s'", activationContext))
	}

	return strings.Join(filters, " and ")
}
