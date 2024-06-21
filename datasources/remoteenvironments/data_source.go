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

package remoteenvironments

import (
	"context"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api"
	settings "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/remote/environment/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/export"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/provider/config"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/provider/logging"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSource() *schema.Resource {
	return &schema.Resource{
		ReadContext: logging.EnableDSCtx(DataSourceRead),
		Schema: map[string]*schema.Schema{
			"remote_environments": {
				Type:     schema.TypeList,
				Elem:     &schema.Resource{Schema: new(settings.Settings).Schema()},
				Computed: true,
			},
		},
	}
}

func DataSourceRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	creds, err := config.Credentials(m, config.CredValDefault)
	if err != nil {
		return diag.FromErr(err)
	}

	service := export.Service(creds, export.ResourceTypes.RemoteEnvironments)
	var stubs api.Stubs
	if stubs, err = service.List(ctx); err != nil {
		return diag.FromErr(err)
	}

	remoteEnvironments := []*settings.Settings{}
	for _, stub := range stubs {
		value := stub.Value.(*settings.Settings)
		remoteEnvironments = append(remoteEnvironments, value)
	}
	marshalled := hcl.Properties{}
	if marshalled.EncodeSlice("remote_environments", remoteEnvironments); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(service.SchemaID())
	d.Set("remote_environments", marshalled["remote_environments"])

	return diag.Diagnostics{}
}
