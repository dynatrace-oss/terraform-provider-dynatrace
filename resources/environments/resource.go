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

package environments

import (
	"context"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/cluster/v2/envs"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/opt"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/provider/config"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/provider/logging"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Resource produces terraform resource definition for Management Zones
func Resource() *schema.Resource {
	return &schema.Resource{
		Schema:        new(envs.Environment).Schema(),
		CreateContext: logging.Enable(Create),
		UpdateContext: logging.Enable(Update),
		ReadContext:   logging.Enable(Read),
		DeleteContext: logging.Enable(Delete),
		Importer:      &schema.ResourceImporter{StateContext: schema.ImportStatePassthroughContext},
	}
}

func NewService(m any) (*envs.ServiceClient, error) {
	creds, err := config.Credentials(m, config.CredValCluster)
	if err != nil {
		return nil, err
	}

	apiService := envs.NewService(creds)
	return apiService, nil
}

// Create expects the configuration within the given ResourceData and sends it to the Dynatrace Server in order to create that resource
func Create(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	_, err := config.Credentials(m, config.CredValCluster)
	if err != nil {
		return diag.FromErr(err)
	}
	config := new(envs.Environment)
	if err := config.UnmarshalHCL(hcl.DecoderFrom(d)); err != nil {
		return diag.FromErr(err)
	}

	service, err := NewService(m)
	if err != nil {
		return diag.FromErr(err)
	}

	config.ID = nil
	objStub, err := service.Create(ctx, config)
	if err != nil {
		if is404(err) {
			d.SetId("")
			return diag.Diagnostics{}
		}
		return diag.FromErr(err)
	}
	d.SetId(objStub.ID)
	return Read(ctx, d, m)
}

// Update expects the configuration within the given ResourceData and send them to the Dynatrace Server in order to update that resource
func Update(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	_, err := config.Credentials(m, config.CredValCluster)
	if err != nil {
		return diag.FromErr(err)
	}
	config := new(envs.Environment)
	if err := config.UnmarshalHCL(hcl.DecoderFrom(d)); err != nil {
		return diag.FromErr(err)
	}

	service, err := NewService(m)
	if err != nil {
		return diag.FromErr(err)
	}

	config.ID = opt.NewString(d.Id())
	if err := service.Update(ctx, config); err != nil {
		if is404(err) {
			d.SetId("")
			return diag.Diagnostics{}
		}
		return diag.FromErr(err)
	}
	return Read(ctx, d, m)
}

// Read queries the Dynatrace Server for the configuration
func Read(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	_, err := config.Credentials(m, config.CredValCluster)
	if err != nil {
		return diag.FromErr(err)
	}

	service, err := NewService(m)
	if err != nil {
		return diag.FromErr(err)
	}

	config, err := service.Get(ctx, d.Id())
	if err != nil {
		if is404(err) {
			d.SetId("")
			return diag.Diagnostics{}
		}
		return diag.FromErr(err)
	}
	marshalled := hcl.Properties{}
	err = config.MarshalHCL(marshalled)
	if err != nil {
		return diag.FromErr(err)
	}
	for k, v := range marshalled {
		d.Set(k, v)
	}
	return diag.Diagnostics{}
}

// Delete the configuration
func Delete(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	_, err := config.Credentials(m, config.CredValCluster)
	if err != nil {
		return diag.FromErr(err)
	}

	service, err := NewService(m)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := service.Delete(ctx, d.Id()); err != nil {
		if is404(err) {
			return diag.Diagnostics{}
		}
		return diag.FromErr(err)
	}
	return diag.Diagnostics{}
}

func is404(err error) bool {
	switch rerr := err.(type) {
	case *rest.Error:
		return rerr.Code == 404
	case rest.Error:
		return rerr.Code == 404
	}
	return false
}
