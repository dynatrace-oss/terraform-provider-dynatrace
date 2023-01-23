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

package alerting

import (
	"context"
	"os"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/provider/config"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/provider/logging"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/export"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var ResourceType = export.ResourceTypes.Alerting

// Resource produces terraform resource definition for Management Zones
func Resource() *schema.Resource {
	return &schema.Resource{
		Schema:        export.AllResources[ResourceType].NewSettings().Schema(),
		CreateContext: logging.Enable(Create),
		UpdateContext: logging.Enable(Update),
		ReadContext:   logging.Enable(Read),
		DeleteContext: logging.Enable(Delete),
		Importer:      &schema.ResourceImporter{StateContext: schema.ImportStatePassthroughContext},
	}
}

func createCredentials(m any) *settings.Credentials {
	conf := m.(*config.ProviderConfiguration)
	return &settings.Credentials{
		Token: conf.APIToken,
		URL:   conf.EnvironmentURL,
	}
}

func Settings() settings.Settings {
	return export.AllResources[ResourceType].NewSettings()
}

func Service(m any) settings.CRUDService[settings.Settings] {
	return export.AllResources[ResourceType].Service(createCredentials(m))
}

// Create expects the configuration within the given ResourceData and sends it to the Dynatrace Server in order to create that resource
func Create(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	settings := Settings()
	if err := settings.UnmarshalHCL(hcl.DecoderFrom(d)); err != nil {
		return diag.FromErr(err)
	}
	stub, err := Service(m).Create(settings)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(stub.ID)
	return Read(ctx, d, m)
}

// Update expects the configuration within the given ResourceData and send them to the Dynatrace Server in order to update that resource
func Update(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	settings := Settings()
	if err := settings.UnmarshalHCL(hcl.DecoderFrom(d)); err != nil {
		return diag.FromErr(err)
	}
	if err := Service(m).Update(d.Id(), settings); err != nil {
		return diag.FromErr(err)
	}
	return Read(ctx, d, m)
}

// Read queries the Dynatrace Server for the configuration
func Read(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	var err error
	var restLogFile *os.File
	restLogFileName := os.Getenv("DT_REST_DEBUG_LOG")
	if len(restLogFileName) > 0 {
		if restLogFile, err = os.Create(restLogFileName); err != nil {
			return diag.FromErr(err)
		}
		rest.SetLogWriter(restLogFile)
	}
	settings := Settings()
	if err := Service(m).Get(d.Id(), settings); err != nil {
		return diag.FromErr(err)
	}
	marshalled := hcl.Properties{}
	if err := settings.MarshalHCL(marshalled); err != nil {
		return diag.FromErr(err)
	}
	for k, v := range marshalled {
		d.Set(k, v)
	}
	return diag.Diagnostics{}
}

// Delete the configuration
func Delete(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	if err := Service(m).Delete(d.Id()); err != nil {
		return diag.FromErr(err)
	}
	return diag.Diagnostics{}
}
