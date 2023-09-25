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

package slo

import (
	"strings"

	common "github.com/dynatrace-oss/terraform-provider-dynatrace/datasources"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/monitoring/slo"
	slosetting "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/monitoring/slo/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/provider/config"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/provider/logging"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSource() *schema.Resource {
	return &schema.Resource{
		Read: logging.EnableDS(DataSourceRead),
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The custom description of the SLO",
			},
			"metric_expression": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The percentage-based metric expression for the calculation of the SLO",
			},
			"enabled": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "The SLO is enabled (`true`) or disabled (`false`)",
			},
			"evaluation_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The evaluation type of the SLO. Currently only `AGGREGATE` is supported",
			},
			"filter": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The entity filter for the SLO evaluation. See [syntax of entity selector](https://dt-url.net/entityselector) for details",
			},
			"target_success": {
				Type:        schema.TypeFloat,
				Computed:    true,
				Description: "The target value of the SLO",
			},
			"target_warning": {
				Type:        schema.TypeFloat,
				Computed:    true,
				Description: "The warning value of the SLO. At warning state the SLO is still fulfilled but is getting close to failure",
			},
			"evaluation_window": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The timeframe during which the SLO is to be evaluated. For the timeframe you can enter expressions like -1h (last hour), -1w (last week) or complex expressions like -2d to now (last two days), -1d/d to now/d (beginning of yesterday to beginning of today).",
			},
			"metric_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "No documentation available",
			},
			"fast_burn_threshold": {
				Type:        schema.TypeFloat,
				Computed:    true,
				Description: "The threshold defines when a burn rate is marked as fast-burning (high-emergency). Burn rates lower than this threshold (and greater than 1) are highlighted as slow-burn (low-emergency)",
			},
			"burn_rate_visualization_enabled": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Burn rate visualization is enabled (`true`) or disabled (`false`)",
			},
			"legacy_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of this setting when referred to by the Config REST API V1",
			},
		},
	}
}

func DataSourceRead(d *schema.ResourceData, m any) (err error) {
	var name string
	if v, ok := d.GetOk("name"); ok {
		name = v.(string)
	}

	creds, err := config.Credentials(m, config.CredValDefault)
	if err != nil {
		return err
	}
	service := slo.Service(creds)
	var stubs api.Stubs
	if stubs, err = service.List(); err != nil {
		return err
	}
	if len(stubs) > 0 {
		for _, stub := range stubs {
			if name == stub.Name {
				var v slosetting.Settings
				if err := service.Get(stub.ID, &v); err != nil {
					return err
				}
				if v.CustomDescription != nil {
					d.Set("description", *v.CustomDescription)
				}
				d.Set("metric_expression", v.MetricExpression)
				d.Set("enabled", v.Enabled)
				d.Set("evaluation_type", v.EvaluationType)
				d.Set("filter", v.Filter)
				d.Set("target_success", v.TargetSuccess)
				d.Set("target_warning", v.TargetWarning)
				d.Set("evaluation_window", v.EvaluationWindow)
				d.Set("metric_name", v.MetricName)
				if v.LegacyID != nil {
					d.Set("legacy_id", *v.LegacyID)
				}
				if v.ErrorBudgetBurnRate != nil {
					d.Set("burn_rate_visualization_enabled", v.ErrorBudgetBurnRate.BurnRateVisualizationEnabled)
					if v.ErrorBudgetBurnRate.FastBurnThreshold != nil {
						d.Set("fast_burn_threshold", v.ErrorBudgetBurnRate.FastBurnThreshold)
					}
				}
				d.SetId(stub.ID)
				return nil
			}
		}
	}

	d.SetId(common.NotFoundID(strings.ToLower(strings.ReplaceAll(name, " ", ""))))
	return nil
}
