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

package mgmz

import (
	"context"
	"sort"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api"
	managementzonessrv "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/managementzones"
	managementzones "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/managementzones/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings/services/cache"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/provider/config"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/provider/logging"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const staticID = "46465fe6-70cb-4564-864f-c3344caae5c0"

func DataSourceMultiple() *schema.Resource {
	return &schema.Resource{
		ReadContext: logging.EnableDSCtx(DataSourceReadMultiple),
		Schema: map[string]*schema.Schema{
			"values": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of the Management Zone when referred to as a Settings 2.0 resource (e.g. from within `dynatrace_slack_notification`)",
						},
						"legacy_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of the Management Zone when referred to as a Configuration API resource (e.g. from within `dynatrace_notification`)",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the Management Zone",
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The description of the Management Zone",
						},
					},
				},
			},
		},
	}
}

func DataSourceReadMultiple(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	creds, err := config.Credentials(m, config.CredValDefault)
	if err != nil {
		return diag.FromErr(err)
	}

	service := cache.Read[*managementzones.Settings](managementzonessrv.Service(creds), true)
	var stubs api.Stubs
	if stubs, err = service.List(ctx); err != nil {
		return diag.FromErr(err)
	}
	d.SetId(staticID)
	values := []map[string]any{}
	sort.SliceStable(stubs, func(i, j int) bool {
		return stubs[i].Name < stubs[j].Name
	})
	for _, stub := range stubs {
		stubValue := stub.Value.(*managementzones.Settings)
		description := ""
		if stubValue != nil && stubValue.Description != nil {
			description = *stubValue.Description
		}
		values = append(values, map[string]any{
			"id":          stub.ID,
			"legacy_id":   stubValue.LegacyID,
			"name":        stub.Name,
			"description": description,
		})
	}
	d.Set("values", values)
	return diag.Diagnostics{}
}
