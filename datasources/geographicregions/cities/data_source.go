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

package cities

import (
	"context"
	"fmt"

	srv "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v2/geographicregions/cities"
	cities "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v2/geographicregions/cities/settings"
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
			"region_code": {
				Type:        schema.TypeString,
				Description: "The code of the required region",
				Required:    true,
			},
			"cities": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Resource{Schema: new(cities.City).Schema()},
			},
		},
	}
}

func DataSourceRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	d.SetId("dynatrace_geo_cities")

	var countryCode string
	if v, ok := d.GetOk("country_code"); ok {
		countryCode = v.(string)
	}
	var regionCode string
	if v, ok := d.GetOk("region_code"); ok {
		regionCode = v.(string)
	}
	countryRegionCode := fmt.Sprintf("%s-%s", countryCode, regionCode)

	creds, err := config.Credentials(m, config.CredValDefault)
	if err != nil {
		return diag.FromErr(err)
	}
	service := srv.Service(creds)

	var settings cities.Settings
	if err := service.Get(ctx, countryRegionCode, &settings); err != nil {
		return diag.FromErr(err)
	}

	marshalled := hcl.Properties{}
	if err := marshalled.EncodeSlice("cities", settings.Regions[0].Cities); err != nil {
		return diag.FromErr(err)
	}
	d.Set("cities", marshalled["cities"])

	return diag.Diagnostics{}
}
