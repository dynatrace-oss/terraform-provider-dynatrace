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

package locations

import (
	"context"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api"
	locations "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/synthetic/locations"
	locsettings "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/synthetic/locations/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/opt"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/provider/config"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/provider/logging"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSource() *schema.Resource {
	return &schema.Resource{
		ReadContext: logging.EnableDSCtx(DataSourceRead),
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"locations": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Elem:     &schema.Resource{Schema: new(locsettings.SyntheticLocation).Schema()},
				Optional: true,
				Computed: true,
			},
		},
	}
}

func DataSourceRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	var id *string
	var name *string

	if v, ok := d.GetOk("id"); ok {
		d.SetId(v.(string))
		id = opt.NewString(v.(string))
	} else {
		d.SetId(uuid.New().String())
	}

	if v, ok := d.GetOk("name"); ok {
		name = opt.NewString(v.(string))
	}

	creds, err := config.Credentials(m, config.CredValDefault)
	if err != nil {
		return diag.FromErr(err)
	}
	var stubs api.Stubs
	if stubs, err = locations.Service(creds).List(ctx); err != nil {
		return diag.FromErr(err)
	}
	locs := locsettings.SyntheticLocations{}
	for _, stub := range stubs {
		if id != nil {
			if *id != stub.ID {
				continue
			}
		}
		if name != nil {
			if *name != stub.Name {
				continue
			}
		}
		value := stub.Value.(*locsettings.SyntheticLocation)
		locs.Locations = append(locs.Locations, value)
	}
	marshalled := hcl.Properties{}
	if err = locs.MarshalHCL(marshalled); err != nil {
		return diag.FromErr(err)
	}
	d.Set("locations", marshalled["location"])

	return diag.Diagnostics{}
}
