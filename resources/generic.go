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

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/export/sensitive"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/provider/config"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/provider/logging"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/confighcl"
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

func NewGenericWithAlwaysPrintingViolationPath(resourceType export.ResourceType, credVal ...int) *Generic {
	descriptor := export.AllResources[resourceType]
	cv := CredValDefault
	if len(credVal) > 0 {
		cv = credVal[0]
	}
	return &Generic{Type: resourceType, Descriptor: descriptor, CredentialValidation: cv, AlwaysPrintViolationPath: true}
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
	Type                     export.ResourceType
	Descriptor               export.ResourceDescriptor
	CredentialValidation     int
	AlwaysPrintViolationPath bool
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

type DiffCustomizer interface {
	CustomizeDiff(ctx context.Context, rd *schema.ResourceDiff, i any) error
}

func (me *Generic) Resource() *schema.Resource {
	stngs := me.Descriptor.NewSettings()
	sch := VisitSchemaMap(stngs.Schema())
	nonUpdateableAttrs := make([]string, 0)
	for k, v := range sch {
		if v.ForceNew || v.Computed && !v.Optional {
			nonUpdateableAttrs = append(nonUpdateableAttrs, k)
		}
	}
	updateableAttrs := len(sch) - len(nonUpdateableAttrs)

	if dep, ok := stngs.(Deprecated); ok {
		resRes := &schema.Resource{
			Schema:             sch,
			CreateContext:      logging.Enable(me.Create),
			ReadContext:        logging.Enable(me.Read),
			DeleteContext:      logging.Enable(me.Delete),
			Importer:           &schema.ResourceImporter{StateContext: schema.ImportStatePassthroughContext},
			DeprecationMessage: dep.Deprecated(),
		}
		if updateableAttrs > 0 {
			resRes.UpdateContext = logging.Enable(me.Update)
		}
		if dc, ok := stngs.(DiffCustomizer); ok {
			resRes.CustomizeDiff = dc.CustomizeDiff
		}
		return resRes
	}

	resRes := &schema.Resource{
		Schema:        sch,
		CreateContext: logging.Enable(me.Create),
		ReadContext:   logging.Enable(me.Read),
		DeleteContext: logging.Enable(me.Delete),
		Importer:      &schema.ResourceImporter{StateContext: schema.ImportStatePassthroughContext},
	}
	if updateableAttrs > 0 {
		resRes.UpdateContext = logging.Enable(me.Update)
	}
	if dc, ok := stngs.(DiffCustomizer); ok {
		resRes.CustomizeDiff = dc.CustomizeDiff
	}
	return resRes
}

func (me *Generic) createCredentials(m any) (*rest.Credentials, error) {
	// By default credential validation follows the default route
	// (EnvURL, APIToken)
	cv := config.CredValDefault
	// Unless `me.CredentialValidation` vetoes it
	// Example `dynatrace_iam_*` resources don't require environment URL
	// But instead need OAuth Credentials
	if me.CredentialValidation != cv {
		cv = me.CredentialValidation
	}
	conf := m.(*config.ProviderConfiguration)
	if _, err := config.Credentials(m, cv); err != nil {
		return nil, err
	}
	return &rest.Credentials{
		Token: conf.APIToken,
		URL:   conf.EnvironmentURL,
		IAM:   conf.IAM,
		OAuth: conf.Automation,
	}, nil
}

func (me *Generic) Settings() settings.Settings {
	return me.Descriptor.NewSettings()
}

func (me *Generic) Service(m any) (settings.CRUDService[settings.Settings], error) {
	creds, err := me.createCredentials(m)
	if err != nil {
		return nil, err
	}
	return me.Descriptor.Service(creds), nil
}

func (me *Generic) Create(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	sttngs := me.Settings()
	if err := hcl.UnmarshalHCL(sttngs, hcl.DecoderFrom(d)); err != nil {
		return diag.FromErr(err)
	}
	service, err := me.Service(m)
	if err != nil {
		return diag.FromErr(err)
	}
	var stub *api.Stub

	stub, err = service.Create(ctx, sttngs)
	if err != nil {
		if restWarning, ok := err.(rest.Warning); ok {
			if stub != nil && len(stub.ID) > 0 {
				d.SetId(stub.ID)
			}
			return diag.Diagnostics{diag.Diagnostic{Severity: diag.Warning, Summary: restWarning.Message}}
		}
		if restError, ok := err.(rest.Error); ok {
			vm := restError.ViolationMessage(me.AlwaysPrintViolationPath)
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

	if settings.RefersToMissingID(sttngs) {
		settingName := settings.Name(sttngs, "")
		logging.File.Printf("The resource `%s` with name `%s` and ID `%s` is not in its final state. It refers to resources that don't exist yet.", me.Type, settingName, stub.ID)
	}

	if _, ok := sttngs.(Computer); ok {
		return me.ReadForSettings(context.WithValue(ctx, settings.ContextKeyStateConfig, sttngs), d, m, sttngs)
	}
	// because of the current rate limitations of api.dynatrace.com we simply trust
	// that the results on the remote side are correct
	// and therefore avoid unnecessary GET calls
	if me.Type == export.ResourceTypes.IAMGroup || me.Type == export.ResourceTypes.IAMPermission || me.Type == export.ResourceTypes.IAMPolicy || me.Type == export.ResourceTypes.IAMPolicyBindings || me.Type == export.ResourceTypes.IAMPolicyBindingsV2 || me.Type == export.ResourceTypes.IAMUser {
		return diag.Diagnostics{}
	}
	return me.Read(context.WithValue(ctx, settings.ContextKeyStateConfig, sttngs), d, m)
}

func (me *Generic) Update(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	sttngs := me.Settings()
	if strings.HasSuffix(d.Id(), "---flawed----") {
		return me.Create(ctx, d, m)
	}
	if err := hcl.UnmarshalHCL(sttngs, hcl.DecoderFrom(d)); err != nil {
		return diag.FromErr(err)
	}
	service, err := me.Service(m)
	if err != nil {
		return diag.FromErr(err)
	}

	if ctx.Value(settings.ContextKeyStateConfig) == nil {
		stateConfig := me.Settings()
		if err := stateConfig.UnmarshalHCL(confighcl.StateDecoderFrom(d, me.Resource())); err == nil {
			ctx = context.WithValue(ctx, settings.ContextKeyStateConfig, stateConfig)
		}
	}
	err = service.Update(ctx, d.Id(), sttngs)
	if err != nil {
		if restWarning, ok := err.(rest.Warning); ok {
			return diag.Diagnostics{diag.Diagnostic{Severity: diag.Warning, Summary: restWarning.Message}}
		}
		if restError, ok := err.(rest.Error); ok {
			vm := restError.ViolationMessage(me.AlwaysPrintViolationPath)
			if len(vm) > 0 {
				return diag.FromErr(errors.New(vm))
			}
			return diag.FromErr(errors.New(restError.Message))
		}
		return diag.FromErr(err)
	}
	if settings.RefersToMissingID(sttngs) {
		settingName := settings.Name(sttngs, "")
		logging.File.Printf("The resource `%s` with name `%s` and ID `%s` is not in its final state. It refers to resources that don't exist yet.", me.Type, settingName, d.Id())
	}
	if _, ok := sttngs.(Computer); ok {
		return me.ReadForSettings(ctx, d, m, sttngs)
	}
	// because of the current rate limitations of api.dynatrace.com we simply trust
	// that the results on the remote side are correct
	// and therefore avoid unnecessary GET calls
	if me.Type == export.ResourceTypes.IAMGroup || me.Type == export.ResourceTypes.IAMPermission || me.Type == export.ResourceTypes.IAMPolicy || me.Type == export.ResourceTypes.IAMPolicyBindings || me.Type == export.ResourceTypes.IAMPolicyBindingsV2 || me.Type == export.ResourceTypes.IAMUser || me.Type == export.ResourceTypes.WebApplication {
		return diag.Diagnostics{}
	}
	// dynatrace_hub_extension_config may contain credentials
	// That service matches against a regex (***457***) for them
	// and replaces the values with what's currently in the state
	// ==>
	// if we execute a Read right now, it uses the "old" state
	// Like with other resources we need to trust that the contents
	// we just applied are correct.
	//
	// to be discussed: should we go that route for all resources
	//                  executing a Read right after an update
	//                  is not common practice anyways
	if me.Type == export.ResourceTypes.HubExtensionConfig {
		return diag.Diagnostics{}
	}
	return me.Read(ctx, d, m)
}

type IDGenerator interface {
	GenID() string
}

func (me *Generic) ReadForSettings(ctx context.Context, d *schema.ResourceData, m any, sttngs settings.Settings) diag.Diagnostics {
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
		stateAttributes := NewAttributes(d.State().Attributes)

		// Replacing Algorithm A
		//   Looks for attributes that contain the value `${state.secret_value.exact}`
		//   As opposed to Algorithm B it doesn't even TRY to take any sets within the configuration into consideration.
		//
		//   As a consequence `${state.secret_value.exact}` should only be used
		//   for attributes of schema.TypeString or schema.TypeList (with elems of schema.TypeString).
		//
		//   If that algorithm finds an attribute with the EXACT SAME address within the state
		//   it replaces the value with the value found within the state.
		//
		//   In case the state does NOT contain an attribute with the exact same address
		//   it assumes that one of these situations
		//
		//   * The value at a specific index within the list of string has indeed changed on the remote side
		//   * The number of entries within the list of strings has changed (i.e. has grown)
		//   * The attribute is part of a SET of blocks within the resource block - where there's exists no
		//     predictable (generic) way to determine the matching address. Replacing with a value from the
		//     state is therefore not possible.
		secretExactAttributes := Attributes{}
		secretExactAttributes.collectForSecretExact("", map[string]any(marshalled))

		for key, value := range secretExactAttributes {
			if value == sensitive.SecretValueExact {
				if stateAttributeValue, found := stateAttributes[key]; found {
					store(marshalled, key, stateAttributeValue)
				}
			}
		}

		// -- Replacing Algorithm B --
		//   Looks for attributes that contain the value `${state.secret_value}`
		//   This algorithm unfortunately behaves unpredictably in situations where
		//   multiple "sensitive" attributes are contained within the same block
		//   The flaw emerges because it's trying to deal with sets of blocks
		//
		//   TODO: Identify all resources that are currently using `${state.secret_value}`
		//         for their sensitive attributes.
		//         For resources wheree the address(es) of these attributes don't contain any
		//         hash codes (i.e. there are no schema.TypeSet involved) change that value
		//         to `${state.secret_value.exact}`. That allows Algorithm A to kick in.
		//
		//         For resources where a switch to Algorithm A isn't possible there is currently
		//         no predictable and bullet proof way to identify the exact match for the secret
		//         value from the state.
		//         Here it's likely better to generate `lifecycle` blocks containing `ignore_changes`
		//         for the relevant addresses.
		//         In other words, the author of the HCL code simply needs to be aware of the fact
		//         that by default the plan for these attributes cannot be empty.
		attributes := Attributes{}
		attributes.collect("", map[string]any(marshalled))

		for key, value := range attributes {
			if value == sensitive.SecretValue {
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
	if err != nil {
		return diag.FromErr(err)
	}
	for k, v := range marshalled {
		err := d.Set(k, v)
		// currently behind an env variable, because it could potentially break existing resources
		// TODO: Remove this once we have proper integration tests
		if err != nil && os.Getenv("DYNATRACE_DEBUG") == "true" {
			return diag.FromErr(err)
		}
	}
	return diag.Diagnostics{}
}

func (me *Generic) Read(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	if strings.HasSuffix(d.Id(), "---flawed----") {
		return diag.Diagnostics{}
	}
	sttngs := me.Settings()
	service, err := me.Service(m)
	if err != nil {
		return diag.FromErr(err)
	}

	if ctx.Value(settings.ContextKeyStateConfig) == nil {
		stateConfig := me.Settings()
		if err := stateConfig.UnmarshalHCL(confighcl.StateDecoderFrom(d, me.Resource())); err == nil {
			ctx = context.WithValue(ctx, settings.ContextKeyStateConfig, stateConfig)
		}
	}
	if err = service.Get(ctx, d.Id(), sttngs); err != nil {
		if err.Error() == "inaccessible" {
			return diag.Diagnostics{}
		}
		if strings.Contains(err.Error(), "re-run with confighcl") {
			tfConfig := me.Settings()
			if err = tfConfig.UnmarshalHCL(confighcl.DecoderFrom(d, me.Resource())); err == nil {
				ctx = context.WithValue(ctx, settings.ContextKeyStateConfig, tfConfig)
			}
			err = service.Get(ctx, d.Id(), sttngs)
		} else {
			if restError, ok := err.(rest.Error); ok {
				if restError.Code == 404 {
					d.SetId("")
					return diag.Diagnostics{}
				}
			}
			return diag.FromErr(err)
		}
	}
	if err != nil {
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
			srv, err := me.Service(m)
			if err != nil {
				return diag.FromErr(err)
			}
			if err = srv.Update(ctx, d.Id(), sttngs); err != nil {
				return diag.FromErr(err)
			}
			d.SetId("")
			return diag.Diagnostics{}
		}
	}
	service, err := me.Service(m)
	if err != nil {
		return diag.FromErr(err)
	}

	if ctx.Value(settings.ContextKeyStateConfig) == nil {
		stateConfig := me.Settings()
		if err := stateConfig.UnmarshalHCL(confighcl.StateDecoderFrom(d, me.Resource())); err == nil {
			ctx = context.WithValue(ctx, settings.ContextKeyStateConfig, stateConfig)
		}
	}
	err = service.Delete(ctx, d.Id())

	if err != nil {
		if restWarning, ok := err.(rest.Warning); ok {
			return diag.Diagnostics{diag.Diagnostic{Severity: diag.Warning, Summary: restWarning.Message}}
		}
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
