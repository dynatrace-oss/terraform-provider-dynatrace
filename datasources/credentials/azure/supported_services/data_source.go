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

package supported_services

import (
	"context"
	"strings"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/credentials/azure/services"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/provider/config"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/provider/logging"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSource() *schema.Resource {
	return &schema.Resource{
		ReadContext: logging.EnableDSCtx(DataSourceRead),
		Schema: map[string]*schema.Schema{
			"except": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Services with the given names won't be included in the results",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"services": {
				Type:        schema.TypeMap,
				Computed:    true,
				Description: "The keys are the names of the supported services. The values provide information whether that service is built in or not.",
				Elem:        &schema.Schema{Type: schema.TypeBool},
			},
		},
	}
}

func DataSourceRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	theID := "b06436b3-f9e1-46c4-b0eb-ed6c9438673c"
	except := map[string]string{}
	if iExcept, ok := d.GetOk("except"); ok && iExcept != nil {
		for _, elem := range iExcept.(*schema.Set).List() {
			el := strings.TrimSpace(strings.ToLower(elem.(string)))
			except[el] = el
			theID = theID + ":" + el
		}
	}
	creds, err := config.Credentials(m, config.CredValDefault)
	if err != nil {
		return diag.FromErr(err)
	}
	srvc := services.NewSupportedServicesService(creds)
	all, err := srvc.List()
	if err != nil {
		return diag.FromErr(err)
	}
	srvmap := map[string]bool{}
	for k, v := range all {
		lk := strings.ToLower(k)
		if _, found := except[lk]; !found {
			srvmap[lk] = v.BuiltIn
		}
	}
	d.Set("services", srvmap)
	d.SetId(theID)
	return diag.Diagnostics{}
}
