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

package service

import (
	dscommon "github.com/dynatrace-oss/terraform-provider-dynatrace/datasources"
	services "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/topology/services"
	servsettings "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/topology/services/settings"
	tagapi "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/topology/tag"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/provider/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSource() *schema.Resource {
	return &schema.Resource{
		Read: DataSourceRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"tags": {
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Description: "Required tags of the service to find",
				MinItems:    1,
			},
		},
	}
}

func DataSourceRead(d *schema.ResourceData, m any) (err error) {
	var name string
	if v, ok := d.GetOk("name"); ok {
		name = v.(string)
	}

	var tagList []any
	var tags []tagapi.Tag
	if v, ok := d.GetOk("tags"); ok {
		sTags := v.(*schema.Set)
		tagList = sTags.List()
		dscommon.StringsToTags(tagList, &tags)
	}

	service := services.Service(config.Credentials(m))
	var stubs settings.Stubs
	if stubs, err = service.List(); err != nil {
		return err
	}
	if len(stubs) > 0 {
		for _, stub := range stubs {
			if name == stub.Name {
				var record servsettings.Settings
				if err = service.Get(stub.ID, &record); err != nil {
					return err
				}
				if dscommon.TagSubsetCheck(record.Tags, tags) {
					d.SetId(record.EntityId)
					return nil
				}
			}
		}
	}

	d.SetId("")
	return nil
}
