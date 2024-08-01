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

package autotagging

import (
	"context"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api"
	autotagging "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/tags/autotagging/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings/services/settings20"
)

const SchemaVersion = "1.0.16"
const SchemaID = "builtin:tags.auto-tagging"

func Service(credentials *settings.Credentials) settings.CRUDService[*autotagging.Settings] {
	return &service{settings20.Service[*autotagging.Settings](credentials, SchemaID, SchemaVersion)}
}

/*
Reason for that wrapper:
	When the attribute `rules_maintained_externally` is set
	then the rules specified with in that resource instance
	are irrelevant. They are getting maintained externally
	or via `dynatrace_autotag_rules`.
*/

type service struct {
	service settings.ListIDCRUDService[*autotagging.Settings]
}

func (me *service) ListIDs(ctx context.Context) (api.Stubs, error) {
	return me.service.ListIDs(ctx)
}

func (me *service) List(ctx context.Context) (api.Stubs, error) {
	return me.service.List(ctx)
}

func (me *service) Get(ctx context.Context, id string, v *autotagging.Settings) error {
	cfg := ctx.Value(settings.ContextKeyStateConfig)

	err := me.service.Get(ctx, id, v)

	if stateConfig, stateConfigFound := cfg.(*autotagging.Settings); stateConfig != nil && stateConfigFound {
		if stateConfig.RulesMaintainedExternally {
			v.Rules = stateConfig.Rules
			v.RulesMaintainedExternally = true
		}
	}

	return err
}

func (me *service) SchemaID() string {
	return me.service.SchemaID()
}

func (me *service) Create(ctx context.Context, v *autotagging.Settings) (*api.Stub, error) {
	if v.RulesMaintainedExternally {
		v.Rules = nil
	}
	return me.service.Create(ctx, v)
}

func (me *service) Update(ctx context.Context, id string, v *autotagging.Settings) error {
	if v.RulesMaintainedExternally {
		var remoteSettings autotagging.Settings
		if err := me.service.Get(ctx, id, &remoteSettings); err != nil {
			return err
		}
		v.Rules = remoteSettings.Rules
	}
	return me.service.Update(ctx, id, v)
}

func (me *service) Delete(ctx context.Context, id string) error {
	return me.service.Delete(ctx, id)
}
