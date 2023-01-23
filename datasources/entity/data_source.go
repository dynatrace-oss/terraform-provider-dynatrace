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

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSource() *schema.Resource {
	return &schema.Resource{
		Read: DataSourceRead,
		Schema: map[string]*schema.Schema{
			"type": {Type: schema.TypeString, Required: true},
			"name": {Type: schema.TypeString, Required: true},
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

	var settings entities.Settings
	service := srv.Service(entityType, config.Credentials(m))
	if err := service.Get(service.SchemaID(), &settings); err != nil {
		return err
	}
	d.SetId(service.SchemaID())
	if len(settings.Entities) != 0 {
		for _, entity := range settings.Entities {
			if name == *entity.DisplayName {
				d.SetId(*entity.EntityId)
				return nil
			}
		}
	}
	d.SetId("")
	return nil
}
