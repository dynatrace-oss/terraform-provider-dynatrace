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

package apitoken

import (
	"context"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/provider/config"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/provider/logging"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	srv "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v2/apitokens"
	apitokens "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v2/apitokens/settings"
)

func DataSourceMultiple() *schema.Resource {
	return &schema.Resource{
		ReadContext: logging.EnableDSCtx(DataSourceReadMultiple),
		Schema: map[string]*schema.Schema{
			"api_tokens": {
				Type:     schema.TypeList,
				Elem:     &schema.Resource{Schema: new(apitokens.APIToken).Schema()},
				Computed: true,
			},
		},
	}
}

func DataSourceReadMultiple(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	d.SetId("dynatrace_api_tokens")

	creds, err := config.Credentials(m, config.CredValDefault)
	if err != nil {
		return diag.FromErr(err)
	}

	service := srv.Service(creds)
	var stubs api.Stubs
	if stubs, err = service.List(ctx); err != nil {
		return diag.FromErr(err)
	}

	if len(stubs) > 0 {
		tokens := []*apitokens.APIToken{}
		for _, stub := range stubs {
			value := stub.Value.(*apitokens.APIToken)
			tokens = append(tokens, value)
		}

		marshalled := hcl.Properties{}
		if marshalled.EncodeSlice("api_tokens", tokens); err != nil {
			return diag.FromErr(err)
		}

		d.Set("api_tokens", marshalled["api_tokens"])
	}

	return diag.Diagnostics{}
}
