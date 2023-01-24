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

package resources

import (
	"context"
	"os"
	"strings"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/provider/config"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/provider/logging"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/export"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func NewGeneric(resourceType export.ResourceType) *Generic {
	descriptor := export.AllResources[resourceType]
	return &Generic{Type: resourceType, Descriptor: descriptor}
}

type Generic struct {
	Type       export.ResourceType
	Descriptor export.ResourceDescriptor
}

func (me *Generic) Resource() *schema.Resource {
	stngs := me.Descriptor.NewSettings()
	sch := stngs.Schema()
	// implicitUpate := false
	// stnt := reflect.ValueOf(stngs).Elem().Type()
	// for idx := 0; idx < stnt.NumField(); idx++ {
	// 	field := stnt.Field(idx)
	// 	if field.Type == implicitUpdateType {
	// 		implicitUpate = true
	// 		break
	// 	}
	// }
	// if implicitUpate {
	// 	sch["replaced_value"] = &schema.Schema{
	// 		Type:        schema.TypeString,
	// 		Description: "for internal use only",
	// 		Optional:    true,
	// 		Computed:    true,
	// 	}
	// }

	return &schema.Resource{
		Schema:        sch,
		CreateContext: logging.Enable(me.Create),
		UpdateContext: logging.Enable(me.Update),
		ReadContext:   logging.Enable(me.Read),
		DeleteContext: logging.Enable(me.Delete),
		Importer:      &schema.ResourceImporter{StateContext: schema.ImportStatePassthroughContext},
	}
}

func (me *Generic) createCredentials(m any) *settings.Credentials {
	conf := m.(*config.ProviderConfiguration)
	return &settings.Credentials{
		Token: conf.APIToken,
		URL:   conf.EnvironmentURL,
	}
}

func (me *Generic) Settings() settings.Settings {
	return me.Descriptor.NewSettings()
}

func (me *Generic) Service(m any) settings.CRUDService[settings.Settings] {
	return me.Descriptor.Service(me.createCredentials(m))
}

func (me *Generic) Create(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	sttngs := me.Settings()
	if err := sttngs.UnmarshalHCL(hcl.DecoderFrom(d)); err != nil {
		return diag.FromErr(err)
	}
	stub, err := me.Service(m).Create(sttngs)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(stub.ID)
	// if settings.SupportsFlawedReasons(sttngs) {
	// 	flawedReasons := settings.GetFlawedReasons(sttngs)
	// 	if len(flawedReasons) > 0 {
	// 		d.Set("flawed_reasons", flawedReasons)
	// 	} else {
	// 		d.Set("flawed_reasons", []string{})
	// 	}
	// }
	return me.Read(ctx, d, m)
}

func (me *Generic) Update(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	sttngs := me.Settings()
	if strings.HasSuffix(d.Id(), "---flawed----") {
		return me.Create(ctx, d, m)
	}
	if err := sttngs.UnmarshalHCL(hcl.DecoderFrom(d)); err != nil {
		return diag.FromErr(err)
	}
	if err := me.Service(m).Update(d.Id(), sttngs); err != nil {
		return diag.FromErr(err)
	}
	return me.Read(ctx, d, m)
}

func (me *Generic) Read(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	if strings.HasSuffix(d.Id(), "---flawed----") {
		return diag.Diagnostics{}
	}
	var err error
	sttngs := me.Settings()
	// if os.Getenv("CACHE_OFFLINE_MODE") != "true" {
	// 	if _, ok := settings.(*vault.Credentials); ok {
	// 		return diag.Diagnostics{}
	// 	}
	// 	if _, ok := settings.(*notifications.Notification); ok {
	// 		return diag.Diagnostics{}
	// 	}
	// }
	service := me.Service(m)
	if err := service.Get(d.Id(), sttngs); err != nil {
		if restError, ok := err.(rest.Error); ok {
			if restError.Code == 404 {
				d.SetId("")
				return diag.Diagnostics{}
			}
		}
		return diag.FromErr(err)
	}
	if preparer, ok := sttngs.(MarshalPreparer); ok {
		preparer.PrepareMarshalHCL(hcl.DecoderFrom(d))
	}
	if os.Getenv("DT_TERRAFORM_IMPORT") == "true" {
		if demoSettings, ok := sttngs.(settings.DemoSettings); ok {
			demoSettings.FillDemoValues()
		}
	}
	marshalled := hcl.Properties{}
	err = sttngs.MarshalHCL(marshalled)
	if os.Getenv("DT_TERRAFORM_IMPORT") != "true" {
		attributes := Attributes{}
		attributes.collect("", map[string]any(marshalled))
		stateAttributes := NewAttributes(d.State().Attributes)
		for key, value := range attributes {
			if value == "${state.secret_value}" {
				stored := false
				matches := stateAttributes.MatchingKeys(key)
				siblings := attributes.Siblings(key)
				for _, m := range matches {
					sibs := stateAttributes.Siblings(m)
					if sibs.Contains(siblings...) {
						stored = true
						store(marshalled, key, stateAttributes[m])
						break
					}
				}
				if !stored {
					remove(marshalled, key)
				}
			}
		}
	}
	// if settings.SupportsFlawedReasons(sttngs) {
	// 	flawedReasons := settings.GetFlawedReasons(sttngs)
	// 	if len(flawedReasons) > 0 {
	// 		d.Set("flawed_reasons", flawedReasons)
	// 	} else {
	// 		d.Set("flawed_reasons", []string{})
	// 		delete(marshalled, "flawed_reasons")
	// 	}
	// }
	if err != nil {
		return diag.FromErr(err)
	}
	for k, v := range marshalled {
		d.Set(k, v)
	}

	return diag.Diagnostics{}
}

func (me *Generic) Delete(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	if strings.HasSuffix(d.Id(), "---flawed----") {
		d.SetId("")
		return diag.Diagnostics{}
	}
	if err := me.Service(m).Delete(d.Id()); err != nil {
		if restError, ok := err.(rest.Error); ok {
			if restError.Code == 404 {
				d.SetId("")
				return diag.Diagnostics{}
			}
		}
		return diag.FromErr(err)
	}
	return diag.Diagnostics{}
}

type MarshalPreparer interface {
	PrepareMarshalHCL(hcl.Decoder) error
}
