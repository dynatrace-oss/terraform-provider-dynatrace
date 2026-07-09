/**
* @license
* Copyright 2026 Dynatrace LLC
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

package platformslo

import (
	"context"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api"
	sloservice "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/slo"
	slosettings "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/slo/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/provider/config"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/provider/logging"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSource() *schema.Resource {
	scm := hcl.SetComputedSchema(new(slosettings.SLO).Schema())
	scm["name"] = &schema.Schema{
		Type:        schema.TypeString,
		Description: "Name of the SLO",
		Required:    true,
	}

	return &schema.Resource{
		ReadContext: logging.EnableDSCtx(DataSourceRead),
		Schema:      scm,
	}
}

func DataSourceRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	name := d.Get("name").(string)

	creds, err := config.Validate(m, config.CredValDefault)
	if err != nil {
		return diag.FromErr(err)
	}

	service := sloservice.Service(creds)
	var stubs api.Stubs
	if stubs, err = service.List(ctx); err != nil {
		return diag.FromErr(err)
	}

	stub := findStubByName(stubs, name)
	if stub == nil {
		d.SetId("")
		return diag.Diagnostics{}
	}

	var v slosettings.SLO
	if err = service.Get(ctx, stub.ID, &v); err != nil {
		return diag.FromErr(err)
	}

	marshalled := hcl.Properties{}
	if err = v.MarshalHCL(marshalled); err != nil {
		return diag.FromErr(err)
	}

	for k, val := range marshalled {
		err = d.Set(k, val)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	d.SetId(stub.ID)
	return diag.Diagnostics{}
}

func findStubByName(stubs api.Stubs, name string) *api.Stub {
	for _, stub := range stubs {
		if name == stub.Name {
			return stub
		}
	}
	return nil
}
