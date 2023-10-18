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

package entities

import (
	srv "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v2/entities"
	entities "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v2/entities/settings"
	entity "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v2/entity/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/provider/config"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/provider/logging"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSource() *schema.Resource {
	return &schema.Resource{
		Read: logging.EnableDS(DataSourceRead),
		Schema: map[string]*schema.Schema{
			"type": {
				Type:          schema.TypeString,
				Optional:      true,
				Description:   "The type of the entities to find, e.g. `HOST`. You cannot use `type` and `entity_selector` at the same time",
				ConflictsWith: []string{"entity_selector"},
			},
			"entity_selector": {
				Type:          schema.TypeString,
				Description:   "An entity selector that filters the entities of interest. You cannot use `type` and `entity_selector` at the same time",
				Optional:      true,
				ConflictsWith: []string{"type"},
			},
			"entities": {
				Type:     schema.TypeList,
				Elem:     &schema.Resource{Schema: new(entity.Entity).Schema()},
				Optional: true,
				Computed: true,
			},
			"from": {
				Type:        schema.TypeString,
				Description: "Limits the time frame entities are queried for - specifically the start of the requested timeframe. Defaults to `now-3y`. You can use one of the following formats:\n  * Timestamp in UTC milliseconds\n  * Human-readable format of `2021-01-25T05:57:01.123+01:00`. If no time zone is specified, `UTC` is used. You can use a space character instead of the `T`. Seconds and fractions of a second are optional\n  * Relative timeframe, back from now. The format is `now-NU/A`, where `N` is the amount of time, `U` is the unit of time, and `A` is an alignment. The alignment rounds all the smaller values to the nearest zero in the past. For example, `now-1y/w` is one year back, aligned by a week. You can also specify relative timeframe without an alignment: `now-NU`. Supported time units for the relative timeframe are:\n    - `m` for minutes\n    - `h` for hours\n    - `d` for days\n    - `w` for weeks\n    - `M` for months\n    - `y` for years",
				Optional:    true,
				Default:     "now-3y",
			},
			"to": {
				Type:        schema.TypeString,
				Description: "Limits the time frame entities are queried for - specifically the end of the requested timeframe. Defaults to `now`. You can use one of the following formats:\n  * Timestamp in UTC milliseconds\n  * Human-readable format of `2021-01-25T05:57:01.123+01:00`. If no time zone is specified, `UTC` is used. You can use a space character instead of the `T`. Seconds and fractions of a second are optional\n  * Relative timeframe, back from now. The format is `now-NU/A`, where `N` is the amount of time, `U` is the unit of time, and `A` is an alignment. The alignment rounds all the smaller values to the nearest zero in the past. For example, `now-1y/w` is one year back, aligned by a week. You can also specify relative timeframe without an alignment: `now-NU`. Supported time units for the relative timeframe are:\n    - `m` for minutes\n    - `h` for hours\n    - `d` for days\n    - `w` for weeks\n    - `M` for months\n    - `y` for years",
				Optional:    true,
				Default:     "now",
			},
		},
	}
}

func DataSourceRead(d *schema.ResourceData, m any) error {
	var entityType string
	if v, ok := d.GetOk("type"); ok {
		entityType = v.(string)
	}

	var entitySelector string
	if v, ok := d.GetOk("entity_selector"); ok {
		entitySelector = v.(string)
	}
	var from string
	if v, ok := d.GetOk("from"); ok {
		from = v.(string)
	}
	var to string
	if v, ok := d.GetOk("to"); ok {
		to = v.(string)
	}
	if to == "now" {
		to = ""
	}

	creds, err := config.Credentials(m, config.CredValDefault)
	if err != nil {
		return err
	}

	var settings entities.Settings
	service := srv.Service(entityType, entitySelector, from, to, creds)
	if err := service.Get(service.SchemaID(), &settings); err != nil {
		return err
	}
	d.SetId(service.SchemaID())
	if len(settings.Entities) != 0 {
		marshalled := hcl.Properties{}
		err := marshalled.Encode("settings", &settings)
		if err != nil {
			return err
		}
		d.Set("entities", marshalled["settings"].([]any)[0].(hcl.Properties)["entities"])
	}
	return nil
}
