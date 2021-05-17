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

package customservices

import (
	"context"

	"github.com/dtcookie/dynatrace/rest"
	"github.com/dtcookie/hcl"
	"github.com/dtcookie/opt"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/config"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/hcl2sdk"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/logging"

	"github.com/dtcookie/dynatrace/api/config/customservices"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Resource produces terraform resource definition for Management Zones
func Resource() *schema.Resource {
	return &schema.Resource{
		Schema:        hcl2sdk.Convert(new(customservices.CustomService).Schema()),
		CreateContext: logging.Enable(Create),
		UpdateContext: logging.Enable(Update),
		ReadContext:   logging.Enable(Read),
		DeleteContext: logging.Enable(Delete),
		Importer:      &schema.ResourceImporter{StateContext: schema.ImportStatePassthroughContext},
	}
}

func NewService(m interface{}) *customservices.ServiceClient {
	conf := m.(*config.ProviderConfiguration)
	apiService := customservices.NewService(conf.DTenvURL, conf.APIToken)
	rest.Verbose = config.HTTPVerbose
	return apiService
}

// Create expects the configuration within the given ResourceData and sends it to the Dynatrace Server in order to create that resource
func Create(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	config := new(customservices.CustomService)
	if err := config.UnmarshalHCL(hcl.DecoderFrom(d)); err != nil {
		return diag.FromErr(err)
	}
	config.ID = nil
	config.Metadata = nil
	technology := customservices.Technology(d.Get("technology").(string))
	objStub, err := NewService(m).Create(config, technology)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(objStub.ID)
	return Read(ctx, d, m)
}

// Update expects the configuration within the given ResourceData and send them to the Dynatrace Server in order to update that resource
func Update(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	config := new(customservices.CustomService)
	if err := config.UnmarshalHCL(hcl.DecoderFrom(d)); err != nil {
		return diag.FromErr(err)
	}
	config.ID = opt.NewString(d.Id())
	config.Metadata = nil
	technology := customservices.Technology(d.Get("technology").(string))
	if err := NewService(m).Update(config, technology); err != nil {
		return diag.FromErr(err)
	}
	return Read(ctx, d, m)
}

// Read queries the Dynatrace Server for the configuration
func Read(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var err error
	srvc := NewService(m)

	var customService *customservices.CustomService

	technology := customservices.Technology(d.Get("technology").(string))
	if technology == "" {
		err = nil
		for _, technology = range []customservices.Technology{customservices.Technologies.DotNet, customservices.Technologies.Java, customservices.Technologies.NodeJS, customservices.Technologies.PHP, customservices.Technologies.Go} {
			if customService, err = srvc.Get(d.Id(), technology, true); err != nil {
				if restErr, ok := err.(*rest.Error); ok {
					if restErr.Code != 404 {
						return diag.FromErr(err)
					}
				} else {
					return diag.FromErr(err)
				}
			}
			if customService != nil {
				break
			}
		}
		if customService == nil {
			return diag.FromErr(err)
		}
	} else {
		if customService, err = srvc.Get(d.Id(), technology, true); err != nil {
			return diag.FromErr(err)
		}
	}

	marshalled, err := customService.MarshalHCL()
	if err != nil {
		return diag.FromErr(err)
	}
	for k, v := range marshalled {
		d.Set(k, v)
	}
	return diag.Diagnostics{}
}

// Delete the configuration
func Delete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	technology := customservices.Technology(d.Get("technology").(string))
	if err := NewService(m).Delete(d.Id(), technology); err != nil {
		return diag.FromErr(err)
	}
	return diag.Diagnostics{}
}
