/**
* @license
* Copyright 2025 Dynatrace LLC
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

package appdetection

import (
	"context"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api"
	appdetectionsrv "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/rum/web/appdetection"
	appdetection "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/rum/web/appdetection/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings/services/cache"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/provider/config"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/provider/logging"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSource() *schema.Resource {
	return &schema.Resource{
		ReadContext: logging.EnableDSCtx(DataSourceRead),
		Schema: map[string]*schema.Schema{
			"values": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Application Detection Rule ID",
						},
						"application_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Application ID",
						},
						"matcher": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Matcher",
						},
						"pattern": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Pattern",
						},
					},
				},
			},
		},
	}
}

func DataSourceRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	d.SetId("dynatrace_application_detection_rules")
	creds, err := config.Credentials(m, config.CredValDefault)
	if err != nil {
		return diag.FromErr(err)
	}

	service := cache.Read[*appdetection.Settings](appdetectionsrv.Service(creds), true)
	var stubs api.Stubs
	if stubs, err = service.List(ctx); err != nil {
		return diag.FromErr(err)
	}
	values := []map[string]any{}
	for _, stub := range stubs {
		stubValue := stub.Value.(*appdetection.Settings)
		values = append(values, map[string]any{
			"id":             stub.ID,
			"application_id": stubValue.ApplicationID,
			"matcher":        stubValue.Matcher,
			"pattern":        stubValue.Pattern,
		})
	}
	d.Set("values", values)
	return diag.Diagnostics{}
}
