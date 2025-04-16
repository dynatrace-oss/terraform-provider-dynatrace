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

package objectivetemplates

import (
	"context"
	"encoding/json"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/shutdown"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/provider/config"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/provider/logging"
	"github.com/google/uuid"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceMulti() *schema.Resource {
	return &schema.Resource{
		ReadContext: logging.EnableDSCtx(DataSourceMultiRead),
		Schema: map[string]*schema.Schema{
			"templates": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of the SLO objective template",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of the SLO objective template",
						},
					},
				},
			},
		},
	}
}

var staticID = uuid.NewString()

func DataSourceMultiRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	creds, err := config.Credentials(m, config.CredValDefault)
	if err != nil {
		return diag.FromErr(err)
	}

	restClient, _ := rest.CreatePlatformClient(ctx, creds.OAuth.EnvironmentURL, creds)

	client := NewClient(restClient)

	templates := []map[string]any{}
	nextPage := true
	var nextPageKey *string
	for nextPage {
		var sol TemplatesResponse
		queryParams := map[string][]string{}
		if nextPageKey != nil {
			queryParams["page-key"] = []string{*nextPageKey}
		} else {
			queryParams["page-size"] = []string{"500"}
		}
		response, err := client.List(ctx, queryParams)
		if err != nil {
			return diag.FromErr(err)
		}
		if shutdown.System.Stopped() {
			diag.Errorf("execution interrupted")
		}
		if err := json.Unmarshal(response.Data, &sol); err != nil {
			return diag.FromErr(err)
		}

		if len(sol.Items) > 0 {
			for _, item := range sol.Items {
				templates = append(templates, map[string]any{
					"id":   item.ID,
					"name": item.Name,
				})
			}
		}
		nextPageKey = sol.NextPageKey
		nextPage = (nextPageKey != nil)
	}

	d.Set("templates", templates)
	d.SetId(staticID)
	return diag.Diagnostics{}
}
