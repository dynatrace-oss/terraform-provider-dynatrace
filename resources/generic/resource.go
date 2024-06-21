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

package generic

import (
	"context"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/generic"
	settings "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/generic/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/provider/config"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/provider/logging"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/confighcl"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Resource() *schema.Resource {
	return &schema.Resource{
		Schema:        new(settings.Settings).Schema(),
		CreateContext: logging.Enable(Create),
		UpdateContext: logging.Enable(Update),
		ReadContext:   logging.Enable(Read),
		DeleteContext: logging.Enable(Delete),
		Importer:      &schema.ResourceImporter{StateContext: schema.ImportStatePassthroughContext},
	}
}

func Create(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	var err error
	creds, err := config.Credentials(m, config.CredValDefault)
	if err != nil {
		return diag.FromErr(err)
	}

	config := new(settings.Settings)
	if err := config.UnmarshalHCL(confighcl.DecoderFrom(d, Resource())); err != nil {
		return diag.FromErr(err)
	}

	var stub *api.Stub
	if stub, err = generic.Service(creds).Create(ctx, config); err != nil {
		return diag.FromErr(err)
	}
	d.SetId(stub.ID)

	config.LocalStorage = config.Value

	marshalled := hcl.Properties{}
	if err := config.MarshalHCL(marshalled); err != nil {
		return diag.FromErr(err)
	}

	for k, v := range marshalled {
		d.Set(k, v)
	}
	return diag.Diagnostics{}
}

func Update(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	creds, err := config.Credentials(m, config.CredValDefault)
	if err != nil {
		return diag.FromErr(err)
	}
	config := new(settings.Settings)
	if err := config.UnmarshalHCL(confighcl.DecoderFrom(d, Resource())); err != nil {
		return diag.FromErr(err)
	}

	if err := generic.Service(creds).Update(ctx, d.Id(), config); err != nil {
		return diag.FromErr(err)
	}

	config.LocalStorage = config.Value
	marshalled := hcl.Properties{}
	if err := config.MarshalHCL(marshalled); err != nil {
		return diag.FromErr(err)
	}

	for k, v := range marshalled {
		d.Set(k, v)
	}
	return diag.Diagnostics{}
}

func Read(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	creds, err := config.Credentials(m, config.CredValDefault)
	if err != nil {
		return diag.FromErr(err)
	}

	var configToMarshal *settings.Settings

	apiConfig := new(settings.Settings)
	if err := generic.Service(creds).Get(ctx, d.Id(), apiConfig); err != nil {
		return diag.FromErr(err)
	}

	// in case the configuration contains secrets
	stateConfig := new(settings.Settings)
	if err := stateConfig.UnmarshalHCL(confighcl.StateDecoderFrom(d, Resource())); err != nil {
		return diag.FromErr(err)
	}
	if len(stateConfig.Value) > 0 {
		stateConfig.Merge(apiConfig)
		stateConfig.LocalStorage = stateConfig.Value
		configToMarshal = stateConfig
	} else {
		apiConfig.LocalStorage = apiConfig.Value
		configToMarshal = apiConfig
	}

	marshalled := hcl.Properties{}
	if err := configToMarshal.MarshalHCL(marshalled); err != nil {
		return diag.FromErr(err)
	}
	for k, v := range marshalled {
		d.Set(k, v)
	}

	return diag.Diagnostics{}
}

func Delete(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	creds, err := config.Credentials(m, config.CredValDefault)
	if err != nil {
		return diag.FromErr(err)
	}
	if err := generic.Service(creds).Delete(ctx, d.Id()); err != nil {
		return diag.FromErr(err)
	}

	return diag.Diagnostics{}
}
