/**
* @license
* Copyright 2025 Dynatrace LLC
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

package serviceusers

import (
	"context"
	"sort"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/iam/serviceusers"
	svu "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/iam/serviceusers/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/provider/config"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/provider/logging"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSource() *schema.Resource {
	return &schema.Resource{
		ReadContext: logging.EnableDSCtx(DataSourceRead),
		Description: "Fetches the groups the service user identified with the given UUID is assigned to",
		Schema: map[string]*schema.Schema{
			"uuid": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The UUID of the service user",
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of the service user",
			},
			"groups": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Computed:    true,
				Description: "Groups that the service user belongs to",
			},
		},
	}
}

func DataSourceRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	var uuid string
	if v, ok := d.GetOk("uuid"); ok {
		uuid = v.(string)
	}
	if len(uuid) == 0 {
		d.SetId("")
		return diag.Diagnostics{}
	}
	d.SetId(uuid)

	creds, err := config.Credentials(m, config.CredValIAM)
	if err != nil {
		return diag.FromErr(err)
	}

	var serviceUser svu.ServiceUser
	service := serviceusers.Service(creds)
	if err := service.Get(ctx, uuid, &serviceUser); err != nil {
		return diag.FromErr(err)
	}
	if len(serviceUser.ServiceUserName) > 0 {
		d.Set("name", serviceUser.ServiceUserName)
	}
	if len(serviceUser.Groups) > 0 {
		sort.Strings(serviceUser.Groups)
		d.Set("groups", serviceUser.Groups)
	} else {
		d.Set("groups", []string{})
	}
	return diag.Diagnostics{}
}
