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

package main

import (
	"context"
	"reflect"

	api "github.com/dtcookie/dynatrace/api/config"
	"github.com/dtcookie/dynatrace/api/config/maintenancewindows"
	"github.com/dtcookie/dynatrace/rest"
	"github.com/dtcookie/dynatrace/terraform"
	"github.com/dtcookie/opt"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/config"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/logging"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type maintWins struct {
}

// Resource produces the schema definition for Maintenance Windows
func (mzs *maintWins) Resource() *schema.Resource {
	resource := terraform.ResourceFor(new(maintenancewindows.MaintenanceWindow))
	resource.CreateContext = logging.Enable(mzs.Create)
	resource.UpdateContext = logging.Enable(mzs.Update)
	resource.ReadContext = logging.Enable(mzs.Read)
	resource.DeleteContext = logging.Enable(mzs.Delete)

	return resource
}

// Create expects the configuration of a Maintenance Window within the given ResourceData
// and sends them to the Dynatrace Server in order to create that resource
func (mzs *maintWins) Create(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	rest.Verbose = config.HTTPVerbose

	var err error
	var resolver terraform.Resolver
	if resolver, err = terraform.NewResolver(d); err != nil {
		return diag.FromErr(err)
	}
	var untypedInstance interface{}
	if untypedInstance, err = resolver.Resolve(reflect.TypeOf(maintenancewindows.MaintenanceWindow{})); err != nil {
		return diag.FromErr(err)
	}

	instance := untypedInstance.(maintenancewindows.MaintenanceWindow)
	instance.ID = nil
	conf := m.(*config.ProviderConfiguration)
	apiService := maintenancewindows.NewService(conf.DTenvURL, conf.APIToken)
	var stub *api.EntityShortRepresentation
	if stub, err = apiService.Create(&instance); err != nil {
		return diag.FromErr(err)
	}
	d.SetId(stub.ID)

	return mzs.Read(ctx, d, m)
}

// Update expects the configuration of a Maintenance Window within the given ResourceData
// and sends them to the Dynatrace Server in order to modify that resource
func (mzs *maintWins) Update(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	rest.Verbose = config.HTTPVerbose

	var err error
	var resolver terraform.Resolver
	if resolver, err = terraform.NewResolver(d); err != nil {
		return diag.FromErr(err)
	}
	var untypedInstance interface{}
	if untypedInstance, err = resolver.Resolve(reflect.TypeOf(maintenancewindows.MaintenanceWindow{})); err != nil {
		return diag.FromErr(err)
	}

	instance := untypedInstance.(maintenancewindows.MaintenanceWindow)
	instance.ID = opt.NewString(d.Id())
	conf := m.(*config.ProviderConfiguration)
	apiService := maintenancewindows.NewService(conf.DTenvURL, conf.APIToken)
	if err = apiService.Update(&instance); err != nil {
		return diag.FromErr(err)
	}
	return mzs.Read(ctx, d, m)
}

// Read queries the Dynatrace REST API for the configuration of a Maintenance Window
// defined by the ID contained within the given ResourceData
func (mzs *maintWins) Read(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	rest.Verbose = config.HTTPVerbose

	var err error
	conf := m.(*config.ProviderConfiguration)
	apiService := maintenancewindows.NewService(conf.DTenvURL, conf.APIToken)
	instance := new(maintenancewindows.MaintenanceWindow)
	if instance, err = apiService.Get(d.Id()); err != nil {
		return diag.FromErr(err)
	}
	if err = terraform.ToTerraform(instance, d); err != nil {
		return diag.FromErr(err)
	}

	return diag.Diagnostics{}
}

// Delete deletes the configuration of a Maintenance Window from the Dynatrace Server
func (mzs *maintWins) Delete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	rest.Verbose = config.HTTPVerbose

	var err error
	conf := m.(*config.ProviderConfiguration)
	apiService := maintenancewindows.NewService(conf.DTenvURL, conf.APIToken)
	if err = apiService.Delete(d.Id()); err != nil {
		return diag.FromErr(err)
	}
	return diag.Diagnostics{}
}
