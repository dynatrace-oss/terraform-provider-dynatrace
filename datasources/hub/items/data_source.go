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

package items

import (
	srv "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v2/hub/items"
	items "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v2/hub/items/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/provider/config"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/provider/logging"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func DataSource() *schema.Resource {
	return &schema.Resource{
		Read: logging.EnableDS(DataSourceRead),
		Schema: map[string]*schema.Schema{
			"items": {
				Type:        schema.TypeList,
				Description: "The items within this list",
				Computed:    true,
				Elem:        &schema.Resource{Schema: new(items.HubItem).Schema()},
			},
			"type": {
				Type:         schema.TypeString,
				Description:  "Represents the type of item. It can be `TECHNOLOGY`, `EXTENSION1` or `EXTENSION2`. If not specified, no restriction regarding type happens",
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{`TECHNOLOGY`, `EXTENSION1`, `EXTENSION2`}, false),
			},
		},
	}
}

func DataSourceRead(d *schema.ResourceData, m any) error {
	creds, err := config.Credentials(m, config.CredValDefault)
	if err != nil {
		return err
	}
	var settings items.HubItemList
	var opts srv.Options
	if v, ok := d.GetOk("type"); ok {
		opts.Type = v.(string)
	}

	service := srv.Service(creds, opts)
	if err := service.Get(service.SchemaID(), &settings); err != nil {
		return err
	}
	d.SetId(service.SchemaID())
	if len(settings.Items) != 0 {
		marshalled := hcl.Properties{}
		if err := settings.MarshalHCL(marshalled); err != nil {
			return err
		}
		d.Set("items", marshalled["items"])
	}
	if len(opts.Type) > 0 {
		d.Set("type", opts.Type)
	}
	return nil
}
