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
	"encoding/json"
	"log"
	"reflect"

	api "github.com/dtcookie/dynatrace/api/config"
	"github.com/dtcookie/dynatrace/api/config/dashboards"
	"github.com/dtcookie/dynatrace/rest"
	"github.com/dtcookie/dynatrace/terraform"
	"github.com/dtcookie/opt"
	"github.com/dtcookie/terraform-provider-dynatrace/config"
	"github.com/dtcookie/terraform-provider-dynatrace/logging"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type resDashboards struct{}

// ResourceDashboards produces resource definitions for Dashboards
func (db *resDashboards) Resource() *schema.Resource {
	resource := terraform.ResourceFor(new(dashboards.Dashboard))
	resource.CreateContext = logging.Enable(db.Create)
	resource.UpdateContext = logging.Enable(db.Update)
	resource.ReadContext = logging.Enable(db.Read)
	resource.DeleteContext = logging.Enable(db.Delete)

	return resource
}

// Create expects configuration for a Dashboard within the given ResourceData
// and sends it via POST request to the Dynatrace Server in order to create that resource.
func (db *resDashboards) Create(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	if config.HTTPVerbose {
		log.Println("Dashboards.Create(...)")
	}

	var err error
	var resolver terraform.Resolver
	if resolver, err = terraform.NewResolver(d); err != nil {
		return diag.FromErr(err)
	}
	var untypedDashboard interface{}
	if untypedDashboard, err = resolver.Resolve(reflect.TypeOf(dashboards.Dashboard{})); err != nil {
		return diag.FromErr(err)
	}

	dashboard := untypedDashboard.(dashboards.Dashboard)
	var data []byte
	if data, err = json.MarshalIndent(dashboard, "", "  "); err != nil {
		return diag.FromErr(err)
	}
	log.Println(string(data))
	dashboard.ID = nil
	dashboard.Metadata.Owner = nil
	conf := m.(*config.ProviderConfiguration)
	dashboardService := dashboards.NewService(conf.DTenvURL, conf.APIToken)
	rest.Verbose = config.HTTPVerbose
	var dashboardStub *api.EntityShortRepresentation
	if dashboardStub, err = dashboardService.Create(&dashboard); err != nil {
		return diag.FromErr(err)
	}
	d.SetId(dashboardStub.ID)

	return db.Read(ctx, d, m)
}

// Update expects configuration for a Dashboard within the given ResourceData
// and sends it via POST request to the Dynatrace Server in order to update that resource
func (db *resDashboards) Update(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	if config.HTTPVerbose {
		log.Println("Dashboards.Update(...)")
	}
	var err error
	var resolver terraform.Resolver
	if resolver, err = terraform.NewResolver(d); err != nil {
		return diag.FromErr(err)
	}
	var untypedDashboard interface{}
	if untypedDashboard, err = resolver.Resolve(reflect.TypeOf(dashboards.Dashboard{})); err != nil {
		return diag.FromErr(err)
	}

	dashboard := untypedDashboard.(dashboards.Dashboard)
	dashboard.ID = opt.NewString(d.Id())
	conf := m.(*config.ProviderConfiguration)
	dashboardService := dashboards.NewService(conf.DTenvURL, conf.APIToken)
	rest.Verbose = config.HTTPVerbose
	if err = dashboardService.Update(&dashboard); err != nil {
		return diag.FromErr(err)
	}
	return db.Read(ctx, d, m)
}

// Read queries the Dynatrace REST API for the configuration of a Dashboard
func (db *resDashboards) Read(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	if config.HTTPVerbose {
		log.Println("Dashboards.Read(...)")
	}
	var err error
	conf := m.(*config.ProviderConfiguration)
	dashboardService := dashboards.NewService(conf.DTenvURL, conf.APIToken)
	rest.Verbose = config.HTTPVerbose
	dashboard := new(dashboards.Dashboard)
	if dashboard, err = dashboardService.Get(d.Id()); err != nil {
		return diag.FromErr(err)
	}
	if err = terraform.ToTerraform(dashboard, d); err != nil {
		return diag.FromErr(err)
	}

	return diag.Diagnostics{}
}

// Delete deletes the dashboard with the given ID
func (db *resDashboards) Delete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	if config.HTTPVerbose {
		log.Println("Dashboards.Delete(...)")
	}

	var err error
	conf := m.(*config.ProviderConfiguration)
	dashboardService := dashboards.NewService(conf.DTenvURL, conf.APIToken)
	rest.Verbose = config.HTTPVerbose
	if err = dashboardService.Delete(d.Id()); err != nil {
		return diag.FromErr(err)
	}
	return diag.Diagnostics{}
}
