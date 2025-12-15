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

package anomalydetectors

import (
	"context"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api"
	anomalydetectors "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/davis/anomalydetectors/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings/services/settings20"
)

const SchemaVersion = "1.0.10"
const SchemaID = "builtin:davis.anomaly-detectors"

func Service(credentials *rest.Credentials) settings.CRUDService[*anomalydetectors.Settings] {
	return &service{
		credentials: credentials,
		service:     settings20.Service[*anomalydetectors.Settings](credentials, SchemaID, SchemaVersion),
	}
}

type service struct {
	credentials *rest.Credentials
	service     settings.CRUDService[*anomalydetectors.Settings]
}

func (me *service) Create(ctx context.Context, v *anomalydetectors.Settings) (*api.Stub, error) {
	stub, err := me.service.Create(ctx, v)
	if err == nil {
		return stub, nil
	}

	// if the anomaly detector requires OAuth, try again with OAuth (without an API token)
	if rest.IsRequiresOAuthError(err) && me.credentials.ContainsOAuthOrPlatformToken() {
		ctx := rest.NewPreferOAuthContext(ctx)
		return me.service.Create(ctx, v)
	}

	return nil, err
}

func (me *service) Update(ctx context.Context, id string, v *anomalydetectors.Settings) error {
	err := me.service.Update(ctx, id, v)
	if err == nil {
		return nil
	}

	// if the anomaly detector requires OAuth, try again with OAuth (without an API token)
	if rest.IsRequiresOAuthError(err) && me.credentials.ContainsOAuthOrPlatformToken() {
		ctx := rest.NewPreferOAuthContext(ctx)
		return me.service.Update(ctx, id, v)
	}
	return err
}

func (me *service) Validate(v *anomalydetectors.Settings) error {
	return nil // Settings 2.0 doesn't offer validation
}

func (me *service) Delete(ctx context.Context, id string) error {
	return me.service.Delete(ctx, id)
}

func (me *service) Get(ctx context.Context, id string, v *anomalydetectors.Settings) error {
	return me.service.Get(ctx, id, v)
}

func (me *service) List(ctx context.Context) (api.Stubs, error) {
	return me.service.List(ctx)
}

func (me *service) SchemaID() string {
	return SchemaID
}
