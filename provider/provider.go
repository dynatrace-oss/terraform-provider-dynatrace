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

package provider

import (
	"context"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/config"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/logging"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/resources/alerting"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/resources/anomalies/applications"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/resources/anomalies/databases"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/resources/anomalies/disks"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/resources/anomalies/hosts"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/resources/anomalies/metrics"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/resources/anomalies/services"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/resources/autotags"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/resources/credentials/aws"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/resources/credentials/azure"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/resources/credentials/k8s"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/resources/customservices"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/resources/dashboards"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/resources/maintenance"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/resources/metrics/calculated/service"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/resources/mgmz"
	hostnaming "github.com/dynatrace-oss/terraform-provider-dynatrace/resources/naming/hosts"
	processgroupnaming "github.com/dynatrace-oss/terraform-provider-dynatrace/resources/naming/processgroups"
	servicenaming "github.com/dynatrace-oss/terraform-provider-dynatrace/resources/naming/services"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/resources/notifications"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/resources/requestattributes"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/resources/slo"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/resources/spans/capture"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/resources/spans/ctxprop"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/resources/spans/entrypoints"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/resources/spans/resattr"

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
	logging.SetOutput()
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
		DataSourcesMap: map[string]*schema.Resource{
			"dynatrace_alerting_profiles": alerting.DataSource(),
		},
		ResourcesMap: map[string]*schema.Resource{
			"dynatrace_custom_service":            customservices.Resource(),
			"dynatrace_dashboard":                 dashboards.Resource(),
			"dynatrace_management_zone":           mgmz.Resource(),
			"dynatrace_maintenance_window":        maintenance.Resource(),
			"dynatrace_request_attribute":         requestattributes.Resource(),
			"dynatrace_alerting_profile":          alerting.Resource(),
			"dynatrace_notification":              notifications.Resource(),
			"dynatrace_autotag":                   autotags.Resource(),
			"dynatrace_aws_credentials":           aws.Resource(),
			"dynatrace_azure_credentials":         azure.Resource(),
			"dynatrace_k8s_credentials":           k8s.Resource(),
			"dynatrace_service_anomalies":         services.Resource(),
			"dynatrace_application_anomalies":     applications.Resource(),
			"dynatrace_host_anomalies":            hosts.Resource(),
			"dynatrace_database_anomalies":        databases.Resource(),
			"dynatrace_custom_anomalies":          metrics.Resource(),
			"dynatrace_disk_anomalies":            disks.Resource(),
			"dynatrace_calculated_service_metric": service.Resource(),
			"dynatrace_service_naming":            servicenaming.Resource(),
			"dynatrace_host_naming":               hostnaming.Resource(),
			"dynatrace_processgroup_naming":       processgroupnaming.Resource(),
			"dynatrace_slo":                       slo.Resource(),
			"dynatrace_span_entry_point":          entrypoints.Resource(),
			"dynatrace_span_capture_rule":         capture.Resource(),
			"dynatrace_span_context_propagation":  ctxprop.Resource(),
			"dynatrace_resource_attributes":       resattr.Resource(),
		},
		ConfigureContextFunc: config.ProviderConfigure,
	}
}
