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

package policies

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/iam/policies"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/provider/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceSingle() *schema.Resource {
	return &schema.Resource{
		Read: DataSourceSingleRead,
		Schema: map[string]*schema.Schema{
			"account": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The account that policiy is defined for. Omit if the policy is not defined for an account but for an environment or is global",
			},
			"environment": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The environment that policiy is defined for. Omit if the policy is not defined for an environment but for an account or is global",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the policy",
			},
			"uuid": {
				Type:        schema.TypeString,
				Description: "The UUID of the policy",
				Optional:    true,
			},
		},
	}
}

func DataSourceSingleRead(d *schema.ResourceData, m any) error {
	var name string
	if v, ok := d.GetOk("name"); ok {
		name = v.(string)
	}
	var environment LevelID
	if v, ok := d.GetOk("environment"); ok {
		environment = LevelID(v.(string))
	}
	var account LevelID
	if v, ok := d.GetOk("account"); ok {
		account = LevelID(v.(string))
	}
	var global LevelID
	if len(environment) == 0 && len(account) == 0 {
		global = LevelID("*")
	}

	creds, err := config.Credentials(m, config.CredValIAM)
	if err != nil {
		return err
	}

	service := policies.Service(creds)
	stubs, err := service.List()
	if err != nil {
		return err
	}
	if len(stubs) > 0 {
		for _, stub := range stubs {
			if stub.Name != name {
				continue
			}
			uuid, levelType, levelID, _ := policies.SplitIDNoDefaults(stub.ID)
			switch levelType {
			case "global":
				if global.Matches(levelID) {
					d.SetId(stub.ID)
					d.Set("uuid", uuid)
					return nil
				}
			case "environment":
				if environment.Matches(levelID) {
					d.SetId(stub.ID)
					d.Set("uuid", uuid)
					return nil
				}
			case "account":
				if account.Matches(levelID) {
					d.SetId(stub.ID)
					d.Set("uuid", uuid)
					return nil
				}
			}
		}
	}
	d.SetId("")
	return nil
}
