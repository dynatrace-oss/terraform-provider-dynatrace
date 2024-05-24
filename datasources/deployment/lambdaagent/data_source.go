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

package lambdaagent

import (
	srv "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/deployment/lambdaagent"
	latest "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/deployment/lambdaagent/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/provider/config"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/provider/logging"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSource() *schema.Resource {
	return &schema.Resource{
		Read: logging.EnableDS(DataSourceRead),
		Schema: map[string]*schema.Schema{
			"java": {
				Type:        schema.TypeString,
				Description: "Latest version name of Java code module",
				Computed:    true,
				Optional:    true,
			},
			"java_with_collector": {
				Type:        schema.TypeString,
				Description: "Latest version name of Java code module with log collector",
				Computed:    true,
				Optional:    true,
			},
			"python": {
				Type:        schema.TypeString,
				Description: "Latest version name of Python code module",
				Computed:    true,
				Optional:    true,
			},
			"python_with_collector": {
				Type:        schema.TypeString,
				Description: "Latest version name of Python code module with log collector",
				Computed:    true,
				Optional:    true,
			},
			"nodejs": {
				Type:        schema.TypeString,
				Description: "Latest version name of NodeJS code module",
				Computed:    true,
				Optional:    true,
			},
			"nodejs_with_collector": {
				Type:        schema.TypeString,
				Description: "Latest version name of NodeJS code module with log collector",
				Computed:    true,
				Optional:    true,
			},
			"collector": {
				Type:        schema.TypeString,
				Description: "Latest version name of standalone log collector",
				Computed:    true,
				Optional:    true,
			},
		},
	}
}

func DataSourceRead(d *schema.ResourceData, m any) error {
	creds, err := config.Credentials(m, config.CredValDefault)
	if err != nil {
		return err
	}

	var latest latest.Latest
	service := srv.Service(creds)
	if err := service.Get("", &latest); err != nil {
		return err
	}

	d.SetId(service.SchemaID())
	d.Set("java", latest.Java)
	d.Set("java_with_collector", latest.JavaWithCollector)
	d.Set("python", latest.Python)
	d.Set("python_with_collector", latest.PythonWithCollector)
	d.Set("nodejs", latest.NodeJS)
	d.Set("nodejs_with_collector", latest.NodeJSWithCollector)
	d.Set("collector", latest.Collector)

	return nil
}
