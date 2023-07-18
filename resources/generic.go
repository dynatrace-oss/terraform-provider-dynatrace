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
	"encoding/json"
	"errors"
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

func NewGeneric(resourceType export.ResourceType, credVal ...int) *Generic {
	descriptor := export.AllResources[resourceType]
	cv := CredValDefault
	if len(credVal) > 0 {
		cv = credVal[0]
	}
	return &Generic{Type: resourceType, Descriptor: descriptor, CredentialValidation: cv}
}

type Computer interface {
	IsComputer() bool
}

const (
	CredValDefault = iota
	CredValIAM
	CredValNone
)

type Generic struct {
	Type                 export.ResourceType
	Descriptor           export.ResourceDescriptor
	CredentialValidation int
}

type Deprecated interface {
	Deprecated() string
}

func VisitResource(res *schema.Resource) {
	if res == nil {
		return
	}
	VisitSchemaMap(res.Schema)
}

func VisitSchema(sch *schema.Schema) {
	if sch == nil {
		return
	}
	if !sch.Computed && sch.Type == schema.TypeString {
		if sch.DiffSuppressFunc != nil {
			storedDiffSuppressFunc := sch.DiffSuppressFunc
			sch.DiffSuppressFunc = func(k, oldValue, newValue string, d *schema.ResourceData) bool {
				if hcl.SuppressJSONorEOT(k, oldValue, newValue, d) {
					return true
				}
				return storedDiffSuppressFunc(k, oldValue, newValue, d)
			}
		} else {
			sch.DiffSuppressFunc = hcl.SuppressJSONorEOT
		}
	}
	if res, ok := sch.Elem.(*schema.Resource); ok {
		VisitResource(res)
	}
}

func VisitSchemaMap(schemata map[string]*schema.Schema) map[string]*schema.Schema {
	if len(schemata) == 0 {
		return schemata
	}
	for _, sch := range schemata {
		VisitSchema(sch)
	}
	return schemata
}

func (me *Generic) Resource() *schema.Resource {
	stngs := me.Descriptor.NewSettings()
	sch := VisitSchemaMap(stngs.Schema())

	if dep, ok := stngs.(Deprecated); ok {
		return &schema.Resource{
			Schema:             sch,
			CreateContext:      logging.Enable(me.Create),
			UpdateContext:      logging.Enable(me.Update),
			ReadContext:        logging.Enable(me.Read),
			DeleteContext:      logging.Enable(me.Delete),
			Importer:           &schema.ResourceImporter{StateContext: schema.ImportStatePassthroughContext},
			DeprecationMessage: dep.Deprecated(),
		}
	}

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
		Token:      conf.APIToken,
		URL:        conf.EnvironmentURL,
		IAM:        conf.IAM,
		Automation: conf.Automation,
	}
}

func (me *Generic) validateCredentials(m any) diag.Diagnostics {
	if me.CredentialValidation != CredValDefault {
		return diag.Diagnostics{}
	}
	conf := m.(*config.ProviderConfiguration)
	if len(conf.EnvironmentURL) == 0 {
		return diag.Errorf("No Environment URL has been specified. Use either the environment variable `DYNATRACE_ENV_URL` or the configuration attribute `dt_env_url` of the provider for that.")
	}
	if !strings.HasPrefix(conf.EnvironmentURL, "https://") && !strings.HasPrefix(conf.EnvironmentURL, "http://") {
		return diag.Errorf("The Environment URL `%s` neither starts with `https://` nor with `http://`. Please check your configuration.\nFor SaaS environments: `https://######.live.dynatrace.com`.\nFor Managed environments: `https://############/e/########-####-####-####-############`", conf.EnvironmentURL)
	}
	if len(conf.APIToken) == 0 {
		return diag.Errorf("No API Token has been specified. Use either the environment variable `DYNATRACE_API_TOKEN` or the configuration attribute `dt_api_token` of the provider for that.")
	}
	return diag.Diagnostics{}
}

func (me *Generic) Settings() settings.Settings {
	return me.Descriptor.NewSettings()
}

func (me *Generic) Service(m any) settings.CRUDService[settings.Settings] {
	return me.Descriptor.Service(me.createCredentials(m))
}

func (me *Generic) Create(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	if diags := me.validateCredentials(m); len(diags) > 0 {
		return diags
	}
	sttngs := me.Settings()
	if err := hcl.UnmarshalHCL(sttngs, hcl.DecoderFrom(d)); err != nil {
		return diag.FromErr(err)
	}
	stub, err := me.Service(m).Create(sttngs)
	if err != nil {
		if restError, ok := err.(rest.Error); ok {
			vm := restError.ViolationMessage()
			if len(vm) > 0 {
				return diag.FromErr(errors.New(vm))
			}
			return diag.FromErr(errors.New(restError.Message))
		}
		return diag.FromErr(err)
	}
	if stub == nil {
		return diag.FromErr(errors.New("stub was nil"))
	}
	// If the stub returned by the Service contains a value
	// and that value also contains information about how to restore
	// that setting to the original state we persist it right away
	restore := settings.GetRestoreOnDelete(stub.Value)
	if restore != nil {
		d.Set("_restore_", *restore)
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
	if _, ok := sttngs.(Computer); ok {
		return me.ReadForSettings(ctx, d, m, sttngs)
	}
	return me.Read(ctx, d, m)
}

func (me *Generic) Update(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	if diags := me.validateCredentials(m); len(diags) > 0 {
		return diags
	}
	sttngs := me.Settings()
	if strings.HasSuffix(d.Id(), "---flawed----") {
		return me.Create(ctx, d, m)
	}
	if err := hcl.UnmarshalHCL(sttngs, hcl.DecoderFrom(d)); err != nil {
		return diag.FromErr(err)
	}
	if err := me.Service(m).Update(d.Id(), sttngs); err != nil {
		if restError, ok := err.(rest.Error); ok {
			vm := restError.ViolationMessage()
			if len(vm) > 0 {
				return diag.FromErr(errors.New(vm))
			}
			return diag.FromErr(errors.New(restError.Message))
		}
		return diag.FromErr(err)
	}
	if _, ok := sttngs.(Computer); ok {
		return me.ReadForSettings(ctx, d, m, sttngs)
	}
	return me.Read(ctx, d, m)
}

func (me *Generic) ReadForSettings(ctx context.Context, d *schema.ResourceData, m any, sttngs settings.Settings) diag.Diagnostics {
	if diags := me.validateCredentials(m); len(diags) > 0 {
		return diags
	}
	var err error
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

func (me *Generic) Read(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	if diags := me.validateCredentials(m); len(diags) > 0 {
		return diags
	}
	if strings.HasSuffix(d.Id(), "---flawed----") {
		return diag.Diagnostics{}
	}
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
	return me.ReadForSettings(ctx, d, m, sttngs)
}

func (me *Generic) Delete(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	if diags := me.validateCredentials(m); len(diags) > 0 {
		return diags
	}
	if strings.HasSuffix(d.Id(), "---flawed----") {
		d.SetId("")
		return diag.Diagnostics{}
	}
	// if the state offers an attribute _restore_ deletion essentially means
	// to restore the settings persisted within that attribute.
	// If the attribute doesn't contain usable data we're deleting as usual
	restorev, ok := d.GetOk("_restore_")
	if ok {
		restore := restorev.(string)
		if len(restore) > 0 {
			sttngs := me.Settings()
			if err := json.Unmarshal([]byte(restore), sttngs); err != nil {
				return diag.FromErr(err)
			}
			if err := me.Service(m).Update(d.Id(), sttngs); err != nil {
				return diag.FromErr(err)
			}
			d.SetId("")
			return diag.Diagnostics{}
		}
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
