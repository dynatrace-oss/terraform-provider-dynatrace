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

package alerting

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/export"
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
		},
	}
}

func DataSourceRead(d *schema.ResourceData, m any) (err error) {
	name := d.Get("name").(string)

	d.SetId("dynatrace_v2_alerting_profiles")
	service := export.Service(config.Credentials(m), export.ResourceTypes.Alerting)
	var stubs settings.Stubs
	if stubs, err = service.List(); err != nil {
		return err
	}

	if len(stubs) == 0 {
		d.SetId("")
		return nil
	}
	for _, stub := range stubs {
		if stub.Name == name {
			d.SetId(stub.ID)
			return nil
		}
	}
	d.SetId("")
	return nil
}
