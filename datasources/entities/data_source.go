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
	creds, err := config.Credentials(m, config.CredValDefault)
	if err != nil {
		return err
	}

	var settings entities.Settings
	service := srv.Service(entityType, entitySelector, creds)
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
