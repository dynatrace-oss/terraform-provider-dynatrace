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
	"context"
	"fmt"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api"
	documents "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/documents/document/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/export"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/provider/config"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/provider/logging"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSource() *schema.Resource {
	return &schema.Resource{
		Description: "Retrieve a list of all documents.",
		ReadContext: logging.EnableDSCtx(DataSourceRead),
		Schema: map[string]*schema.Schema{
			"type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The type of documents to query for. Leave empty if you want to query for all kinds of documents. Possible values are `dashboard` or `notebook`",
			},
			"values": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unique identifier of the document.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the document.",
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The type of the document. This could be a specific format or category the document belongs to.",
						},
						"owner": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The owner of the document. This could be a user or a group that has ownership rights over the document.",
						},
					},
				},
			},
		},
	}
}

func DataSourceRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	values := []map[string]any{}
	var docType string
	if v, ok := d.GetOk("type"); ok {
		docType = v.(string)
	}

	creds, err := config.Credentials(m, config.CredValAutomation)
	if err != nil {
		return diag.FromErr(err)
	}

	service := export.Service(creds, export.ResourceTypes.Documents)
	var stubs api.Stubs
	if stubs, err = service.List(ctx); err != nil {
		return diag.FromErr(err)
	}

	if len(stubs) > 0 {
		for _, stub := range stubs {
			if len(docType) > 0 {
				stubDocType := stub.Value.(*documents.Document).Type
				if stubDocType != docType {
					continue
				}
			}
			m := map[string]any{
				"id":    stub.ID,
				"name":  stub.Name,
				"type":  stub.Value.(*documents.Document).Type,
				"owner": stub.Value.(*documents.Document).Owner,
			}

			values = append(values, m)
		}

		if len(docType) > 0 {
			d.SetId(fmt.Sprintf("documents[%s]", docType))
		} else {
			d.SetId("documents")
		}
		d.Set("values", values)
		return diag.Diagnostics{}
	}

	d.SetId("")
	return diag.Diagnostics{}
}
