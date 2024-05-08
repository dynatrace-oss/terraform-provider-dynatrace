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
	"fmt"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/iam/policies"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/provider/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSource() *schema.Resource {
	return &schema.Resource{
		Read: DataSourceRead,
		Schema: map[string]*schema.Schema{
			"accounts": {
				Type:        schema.TypeList,
				MinItems:    1,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "The results will contain policies defined for the given accountID. If one of the entries contains `*` the results will contain policies for all accounts",
			},
			"environments": {
				Type:        schema.TypeList,
				MinItems:    1,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "The results will contain policies defined for the given environments. If one of the entries contains `*` the results will contain policies for all environments",
			},
			"global": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "If `true` the results will contain global policies",
			},
			"policies": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Description: "The ID of the policy as it is known by Terraform. It contains the UUID, LevelType and LevelID of the policy in concatenated form",
							Required:    true,
						},
						"uuid": {
							Type:        schema.TypeString,
							Description: "The UUID of the policy",
							Required:    true,
						},
						"name": {
							Type:        schema.TypeString,
							Description: "The name of the policy",
							Required:    true,
						},
						"account": {
							Type:        schema.TypeString,
							Description: "The account UUID the policy is defined for",
							Optional:    true,
						},
						"environment": {
							Type:        schema.TypeString,
							Description: "The environment ID the policy is defined for",
							Optional:    true,
						},
						"global": {
							Type:        schema.TypeBool,
							Description: "`true` if this is a global policy`",
							Optional:    true,
						},
					},
				},
			},
		},
	}
}

type LevelID string

func (l LevelID) Matches(value string) bool {
	return string(l) == "*" || string(l) == value
}

func DataSourceRead(d *schema.ResourceData, m any) error {
	var global bool
	if v, ok := d.GetOk("global"); ok {
		global = v.(bool)
	}
	var environments []LevelID
	if v, ok := d.GetOk("environments"); ok {
		for _, elem := range v.([]any) {
			environments = append(environments, LevelID(elem.(string)))
		}
	}
	var accounts []LevelID
	if v, ok := d.GetOk("accounts"); ok {
		for _, elem := range v.([]any) {
			accounts = append(accounts, LevelID(elem.(string)))
		}
	}
	dataSourceID := fmt.Sprintf("%#v.%#v.%#v", global, environments, accounts)
	creds, err := config.Credentials(m, config.CredValIAM)
	if err != nil {
		return err
	}

	service := policies.ServiceWithGloabals(creds)
	stubs, err := service.ListWithGlobals()
	if err != nil {
		return err
	}
	results := []map[string]any{}
	if len(stubs) > 0 {
		for _, stub := range stubs {
			uuid, levelType, levelID, _ := policies.SplitIDNoDefaults(stub.ID)
			switch levelType {
			case "global":
				if global {
					results = append(results, map[string]any{
						"id":     stub.ID,
						"name":   stub.Name,
						"uuid":   uuid,
						"global": true,
					})
				}
			case "environment":
				for _, environment := range environments {
					if environment.Matches(levelID) {
						results = append(results, map[string]any{
							"id":          stub.ID,
							"name":        stub.Name,
							"uuid":        uuid,
							"global":      false,
							"environment": levelID,
						})
						break
					}
				}
			case "account":
				for _, account := range accounts {
					if account.Matches(levelID) {
						results = append(results, map[string]any{
							"id":      stub.ID,
							"name":    stub.Name,
							"uuid":    uuid,
							"global":  false,
							"account": levelID,
						})
						break
					}
				}
			}
		}
	}
	d.Set("policies", results)
	d.SetId(dataSourceID)
	return nil
}
