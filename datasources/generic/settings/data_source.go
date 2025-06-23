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

package settings

import (
	"context"
	"fmt"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api"
	srv "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/generic"
	generic "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/generic/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/provider/config"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/provider/logging"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSource() *schema.Resource {
	return &schema.Resource{
		ReadContext: logging.EnableDSCtx(DataSourceRead),
		Schema: map[string]*schema.Schema{
			"schema": {
				Type:         schema.TypeString,
				Description:  "Schema IDs to which the requested objects belong",
				Optional:     true,
				Computed:     true,
				AtLeastOneOf: []string{"schema", "scope"},
			},
			"scope": {
				Type:         schema.TypeString,
				Description:  "Scope that the requested objects target",
				Optional:     true,
				Computed:     true,
				AtLeastOneOf: []string{"schema", "scope"},
			},
			"filter": {
				Type:        schema.TypeString,
				Description: "Filter of the requested objects",
				Optional:    true,
			},
			"value": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func DataSourceRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	var query srv.QueryParams
	if v, ok := d.GetOk("schema"); ok {
		query.Schema = v.(string)
	}
	if v, ok := d.GetOk("scope"); ok {
		query.Scope = v.(string)
	}
	if v, ok := d.GetOk("filter"); ok {
		query.Filter = v.(string)
	}

	d.SetId(fmt.Sprintf("generic_settings[%s][%s][%s]", query.Schema, query.Scope, query.Filter))

	creds, err := config.Credentials(m, config.CredValDefault)
	if err != nil {
		return diag.FromErr(err)
	}

	var stubs api.Stubs
	myService := srv.Service(creds)
	if spService, ok := myService.(SpecificLister); ok {
		stubs, err = spService.ListSpecific(ctx, query)
	}
	if err != nil {
		return diag.FromErr(err)
	}

	if len(stubs) > 0 {
		value := stubs[0].Value.(*generic.Settings)

		marshalled := hcl.Properties{}
		if err := value.MarshalHCL(marshalled); err != nil {
			return diag.FromErr(err)
		}

		for k, v := range marshalled {
			d.Set(k, v)
		}
	}

	return diag.Diagnostics{}
}
