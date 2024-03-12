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

package document

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api"
	documents "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/documents/document/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/export"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/provider/config"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/provider/logging"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSource() *schema.Resource {
	return &schema.Resource{
		Read: logging.EnableDS(DataSourceRead),
		Schema: map[string]*schema.Schema{
			"values": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "",
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "",
						},
						"owner": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "",
						},
					},
				},
			},
		},
	}
}

func DataSourceRead(d *schema.ResourceData, m any) (err error) {
	values := []map[string]any{}

	creds, err := config.Credentials(m, config.CredValAutomation)
	if err != nil {
		return err
	}

	service := export.Service(creds, export.ResourceTypes.Documents)
	var stubs api.Stubs
	if stubs, err = service.List(); err != nil {
		return err
	}

	if len(stubs) > 0 {
		for _, stub := range stubs {
			m := map[string]any{
				"id":    stub.ID,
				"name":  stub.Name,
				"type":  stub.Value.(*documents.Document).Type,
				"owner": stub.Value.(*documents.Document).Owner,
			}

			values = append(values, m)
		}

		d.SetId("documents")
		d.Set("values", values)
		return nil
	}

	d.SetId("")
	return nil
}
