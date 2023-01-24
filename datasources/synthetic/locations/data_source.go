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
	locations "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/synthetic/locations"
	locsettings "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/synthetic/locations/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/opt"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/provider/config"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSource() *schema.Resource {
	return &schema.Resource{
		Read: DataSourceRead,
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

func DataSourceRead(d *schema.ResourceData, m any) (err error) {
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

	var stubs settings.Stubs
	if stubs, err = locations.Service(config.Credentials(m)).List(); err != nil {
		return err
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
		return err
	}
	d.Set("locations", marshalled["locations"])

	return nil
}
