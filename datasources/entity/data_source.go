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

package entity

import (
	srv "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v2/entities"
	entities "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v2/entities/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/provider/config"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/provider/logging"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSource() *schema.Resource {
	return &schema.Resource{
		Read: logging.EnableDS(DataSourceRead),
		Schema: map[string]*schema.Schema{
			"type": {
				Type:         schema.TypeString,
				Optional:     true,
				RequiredWith: []string{"type", "name"},
			},
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				RequiredWith: []string{"type", "name"},
			},
			"entity_selector": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"type", "name"},
			},
			"properties": {
				Type:        schema.TypeMap,
				Description: "Properties defining the entity.",
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
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
	var name string
	if v, ok := d.GetOk("name"); ok {
		name = v.(string)
	}
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
	// service := cache.Read(srv.Service(entityType, entitySelector, creds), true)
	service := srv.Service(entityType, entitySelector, from, to, creds)
	if err := service.Get(service.SchemaID(), &settings); err != nil {
		return err
	}
	if len(settings.Entities) != 0 {
		if len(name) > 0 {
			// When looking for a specific name
			// The first entity that matches that name will be returned
			for _, entity := range settings.Entities {
				if name == *entity.DisplayName {
					d.SetId(*entity.EntityId)
					if len(entity.Properties) > 0 {
						d.Set("properties", entity.Properties)
					}
					return nil
				}
			}
			// No entity with that name -> ID empty string signals not found
			d.SetId("")
			return nil
		}
		// When looking via entity_selector the first result
		// will be returned
		d.SetId(*settings.Entities[0].EntityId)
		if len(settings.Entities[0].Properties) > 0 {
			d.Set("properties", settings.Entities[0].Properties)
		}
		return nil
	}
	// Without any matches we're setting the ID to an empty string
	d.SetId("")
	return nil
}
