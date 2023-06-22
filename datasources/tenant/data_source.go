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

package tenant

import (
	"strings"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/provider/config"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/provider/logging"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSource() *schema.Resource {
	return &schema.Resource{
		Read: logging.EnableDS(DataSourceRead),
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func DataSourceRead(d *schema.ResourceData, m any) error {
	creds := config.Credentials(m)
	if len(creds.URL) == 0 {
		d.SetId("")
		return nil
	}
	var tenant string
	if strings.Contains(creds.URL, "/e/") {
		idx := strings.Index(creds.URL, "/e/")
		tenant = strings.TrimSuffix(strings.TrimPrefix(creds.URL[idx:], "/e/"), "/")
	} else if strings.HasPrefix(creds.URL, "http://") {
		tenant = strings.TrimPrefix(creds.URL, "http://")
		if idx := strings.Index(tenant, "."); idx != -1 {
			tenant = tenant[:idx]
		}
	} else if strings.HasPrefix(creds.URL, "https://") {
		tenant = strings.TrimPrefix(creds.URL, "https://")
		if idx := strings.Index(tenant, "."); idx != -1 {
			tenant = tenant[:idx]
		}
	}
	d.SetId(tenant)
	d.Set("name", tenant)
	return nil
}
