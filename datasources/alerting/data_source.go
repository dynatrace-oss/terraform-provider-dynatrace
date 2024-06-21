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

package alerting

import (
	"context"
	"sort"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api"
	alertingsrv "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/alerting/profile"
	alerting "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/alerting/profile/settings"
	managementzonessrv "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/managementzones"
	managementzones "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/managementzones/settings"
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
			"profiles": {
				Type:     schema.TypeMap,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
			},
			"values": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of the Alerting Profile when referred to as a Settings 2.0 resource (e.g. from within `dynatrace_slack_notification`)",
						},
						"legacy_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of the Alerting Profile when referred to as a Configuration API resource (e.g. from within `dynatrace_notification`)",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the Alerting Profile",
						},
						"management_zone_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "The ID of the management zone to which the alerting profile applies (Settings 2.0)",
						},
						"management_zone_legacy_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "The ID of the management zone to which the alerting profile applies (Configuration API)",
						},
					},
				},
			},
		},
	}
}

func DataSourceRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	d.SetId("dynatrace_alerting_profiles")
	creds, err := config.Credentials(m, config.CredValDefault)
	if err != nil {
		return diag.FromErr(err)
	}
	service := cache.Read[*alerting.Profile](alertingsrv.Service(creds), true)
	var stubs api.Stubs
	if stubs, err = service.List(ctx); err != nil {
		return diag.FromErr(err)
	}
	mgmzService := cache.Read[*managementzones.Settings](managementzonessrv.Service(creds), true)
	var mgmzStubs api.Stubs
	if mgmzStubs, err = mgmzService.List(ctx); err != nil {
		return diag.FromErr(err)
	}
	mgms := map[string]*api.Stub{}
	for _, mgmzStub := range mgmzStubs {
		mgms[*mgmzStub.LegacyID] = mgmzStub
	}

	profiles := map[string]any{}
	for _, stub := range stubs {
		profiles[stub.Name] = stub.LegacyID
	}
	d.Set("profiles", profiles)
	values := []map[string]any{}
	sort.SliceStable(stubs, func(i, j int) bool {
		return stubs[i].Name < stubs[j].Name
	})
	for _, stub := range stubs {
		stubValue := stub.Value.(*alerting.Profile)
		var mgmzLegacyID *string
		if stubValue != nil && stubValue.ManagementZone != nil {
			mgmzLegacyID = stubValue.ManagementZone
		}
		var mgmzID string
		if mgmzLegacyID != nil {
			if mgmzStub, ok := mgms[*mgmzLegacyID]; ok && mgmzStub != nil {
				mgmzID = mgmzStub.ID
			}

		}
		m := map[string]any{
			"id":                        stub.ID,
			"legacy_id":                 stubValue.LegacyID,
			"name":                      stub.Name,
			"management_zone_id":        "",
			"management_zone_legacy_id": "",
		}
		if mgmzLegacyID != nil {
			m["management_zone_id"] = mgmzID
			m["management_zone_legacy_id"] = *mgmzLegacyID
		}
		values = append(values, m)
	}
	d.Set("values", values)
	return diag.Diagnostics{}
}
