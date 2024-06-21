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

package processgroup

import (
	"context"

	dscommon "github.com/dynatrace-oss/terraform-provider-dynatrace/datasources"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api"
	processgroups "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/topology/processgroups"
	pgsettings "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/topology/processgroups/settings"
	tagapi "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/topology/tag"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/provider/config"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/provider/logging"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSource() *schema.Resource {
	return &schema.Resource{
		ReadContext: logging.EnableDSCtx(DataSourceRead),
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"tags": {
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Description: "Required tags of the process group to find",
				MinItems:    1,
			},
		},
	}
}

func DataSourceRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	var name string
	if v, ok := d.GetOk("name"); ok {
		name = v.(string)
	}

	var tagList []any
	var tags []tagapi.Tag
	if v, ok := d.GetOk("tags"); ok {
		sTags := v.(*schema.Set)
		tagList = sTags.List()
		dscommon.StringsToTags(tagList, &tags)
	}
	creds, err := config.Credentials(m, config.CredValDefault)
	if err != nil {
		return diag.FromErr(err)
	}

	service := processgroups.Service(creds)
	var stubs api.Stubs
	if stubs, err = service.List(ctx); err != nil {
		return diag.FromErr(err)
	}
	if len(stubs) > 0 {
		for _, stub := range stubs {
			if name == stub.Name {
				var processGroup pgsettings.ProcessGroup
				if err = service.Get(ctx, stub.ID, &processGroup); err != nil {
					return diag.FromErr(err)
				}
				if dscommon.TagSubsetCheck(processGroup.Tags, tags) {
					d.SetId(processGroup.EntityId)
					return diag.Diagnostics{}
				}
			}
		}
	}

	d.SetId("")
	return diag.Diagnostics{}
}
