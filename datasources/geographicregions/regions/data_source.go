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

package regions

import (
	"context"

	srv "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v2/geographicregions/regions"
	regions "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v2/geographicregions/regions/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/provider/config"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/provider/logging"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSource() *schema.Resource {
	return &schema.Resource{
		ReadContext: logging.EnableDSCtx(DataSourceRead),
		Schema: map[string]*schema.Schema{
			"country_code": {
				Type:        schema.TypeString,
				Description: "The ISO code of the required country",
				Required:    true,
			},
			"regions": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Resource{Schema: new(regions.Region).Schema()},
			},
		},
	}
}

func DataSourceRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	d.SetId("dynatrace_geo_regions")

	var countryCode string
	if v, ok := d.GetOk("country_code"); ok {
		countryCode = v.(string)
	}

	creds, err := config.Credentials(m, config.CredValDefault)
	if err != nil {
		return diag.FromErr(err)
	}
	service := srv.Service(creds)

	var settings regions.Settings
	if err := service.Get(ctx, countryCode, &settings); err != nil {
		return diag.FromErr(err)
	}

	marshalled := hcl.Properties{}
	if err := marshalled.EncodeSlice("regions", settings.Regions); err != nil {
		return diag.FromErr(err)
	}
	d.Set("regions", marshalled["regions"])

	return diag.Diagnostics{}
}
