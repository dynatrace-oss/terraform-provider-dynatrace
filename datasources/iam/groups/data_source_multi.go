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

package groups

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/iam/groups"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/provider/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceMulti() *schema.Resource {
	return &schema.Resource{
		Read: DataSourceMultiRead,
		Schema: map[string]*schema.Schema{
			"groups": {
				Type:     schema.TypeMap,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
		},
	}
}

func DataSourceMultiRead(d *schema.ResourceData, m any) error {
	creds, err := config.Credentials(m, config.CredValIAM)
	if err != nil {
		return err
	}

	service := groups.Service(creds)
	stubs, err := service.List()
	if err != nil {
		return err
	}
	d.SetId("iam-groups")
	groups := map[string]string{}
	if len(stubs) > 0 {
		for _, stub := range stubs {
			groups[stub.ID] = stub.Name
		}
	}
	d.Set("groups", groups)
	return nil
}
