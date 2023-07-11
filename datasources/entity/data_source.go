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

	var settings entities.Settings
	// service := cache.Read(srv.Service(entityType, entitySelector, config.Credentials(m)), true)
	service := srv.Service(entityType, entitySelector, config.Credentials(m))
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
