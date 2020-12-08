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

	"github.com/dtcookie/terraform-provider-dynatrace/config"
	customservice "github.com/dtcookie/terraform-provider-dynatrace/custom_service"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ResourceSpecification has no documentation
type ResourceSpecification interface {
	Resource() *schema.Resource
	Create(context.Context, *schema.ResourceData, interface{}) diag.Diagnostics
	Update(context.Context, *schema.ResourceData, interface{}) diag.Diagnostics
	Read(context.Context, *schema.ResourceData, interface{}) diag.Diagnostics
	Delete(context.Context, *schema.ResourceData, interface{}) diag.Diagnostics
}

// Provider function for Dynatrace API
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"dt_env_url": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{"DYNATRACE_ENV_URL", "DT_ENV_URL"}, nil),
			},
			"dt_api_token": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{"DYNATRACE_API_TOKEN", "DT_API_TOKEN"}, nil),
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"dynatrace_custom_service":     customservice.Resource(),
			"dynatrace_dashboard":          new(resDashboards).Resource(),
			"dynatrace_management_zone":    new(managementZones).Resource(),
			"dynatrace_maintenance_window": new(maintWins).Resource(),
			"dynatrace_request_attribute":  new(reqAtts).Resource(),
			"dynatrace_alerting_profile":   new(alertingProfiles).Resource(),
		},
		ConfigureContextFunc: config.ProviderConfigure,
	}
}
