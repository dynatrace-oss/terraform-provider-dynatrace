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
	"fmt"

	alertingapi "github.com/dtcookie/dynatrace/api/config/v2/alerting"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/config"
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

func DataSourceRead(d *schema.ResourceData, m interface{}) error {
	name := d.Get("name").(string)

	conf := m.(*config.ProviderConfiguration)
	apiService := alertingapi.NewService(conf.DTApiV2URL, conf.APIToken)
	ids, err := apiService.List()
	if err != nil {
		return err
	}
	if ids == nil {
		return fmt.Errorf("No alerting profile with name `%s` found", name)
	}
	for _, id := range ids {
		profile, err := apiService.Get(id)
		if err != nil {
			return err
		}
		if profile.Name == name {
			d.SetId(id)
			return nil
		}
	}
	return fmt.Errorf("No alerting profile with name `%s` found", name)
}
