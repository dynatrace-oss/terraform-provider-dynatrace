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

package users

import (
	"sort"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/iam/users"
	usr "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/iam/users/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/provider/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSource() *schema.Resource {
	return &schema.Resource{
		Read:        DataSourceRead,
		Description: "Fetches the groups the user identified with the given email is assigned to",
		Schema: map[string]*schema.Schema{
			"email": {
				Type:     schema.TypeString,
				Required: true,
			},
			"groups": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
		},
	}
}

func DataSourceRead(d *schema.ResourceData, m any) error {
	var email string
	if v, ok := d.GetOk("email"); ok {
		email = v.(string)
	}
	if len(email) == 0 {
		d.SetId("")
		return nil
	}
	d.SetId(email)

	creds, err := config.Credentials(m, config.CredValIAM)
	if err != nil {
		return err
	}

	var user usr.User
	service := users.Service(creds)
	if err := service.Get(email, &user); err != nil {
		return err
	}
	if len(user.Groups) > 0 {
		sort.Strings(user.Groups)
		d.Set("groups", user.Groups)
	} else {
		d.Set("groups", []string{})
	}
	return nil
}
