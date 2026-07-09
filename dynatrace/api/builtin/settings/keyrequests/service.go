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

package keyrequests

import (
	"context"
	"fmt"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api"
	keyrequests "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/settings/keyrequests/settings"
	srv "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v2/entities"
	entities "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v2/entities/settings"
	toposervices "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v2/entity"
	entity "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v2/entity/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings/services/settings20"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/provider/config"
)

const SchemaVersion = "0.1.9"
const SchemaID = "builtin:settings.subscriptions.service"

func Service(credentials *config.ProviderConfiguration) settings.CRUDService[*keyrequests.Settings] {
	var topologyService settings.RService[*entity.Entity]
	if settings.ExportRunning {
		topologyService = toposervices.DataSourceService(credentials)
	} else {
		topologyService = toposervices.Service(credentials)
	}
	return &service{
		service: settings20.Service[*keyrequests.Settings](credentials, SchemaID, SchemaVersion, &settings20.ServiceOptions[*keyrequests.Settings]{
			Name: func(ctx context.Context, id string, v *keyrequests.Settings) (string, error) {
				service := settings.NewSettings(topologyService)
				if err := topologyService.Get(ctx, v.ServiceID, service); err != nil {
					return "", err
				}
				return "Key Requests for " + *service.DisplayName, nil
			},
		}),
		credentials: credentials,
		client:      rest.HybridClient(credentials.EnvironmentURL, credentials.APIToken, credentials.Platform),
	}
}

type service struct {
	service     settings.CRUDService[*keyrequests.Settings]
	credentials *config.ProviderConfiguration
	client      rest.Client
}

func (me *service) Create(ctx context.Context, v *keyrequests.Settings) (*api.Stub, error) {
	return me.service.Create(ctx, v)
}

func (me *service) Update(ctx context.Context, id string, v *keyrequests.Settings) error {
	return me.service.Update(ctx, id, v)
}

func (me *service) Get(ctx context.Context, id string, v *keyrequests.Settings) error {
	if err := me.service.Get(ctx, id, v); err != nil {
		return err
	}

	v.KeyRequestIDs = me.getKeyRequestIDs(ctx, v.ServiceID, v.Names)

	return nil
}

// getKeyRequestIDs fetches the service_method of a given service and maps via the given names the IDs to a map {name: ID}
// if the request errors the state value is returned
func (me *service) getKeyRequestIDs(ctx context.Context, serviceID string, names []string) map[string]string {
	var entitySettings entities.Settings
	service := srv.Service("", "", fmt.Sprintf("type(\"SERVICE_METHOD\"),fromRelationships.isServiceMethodOf(type(\"SERVICE_METHOD_GROUP\"),fromRelationships.isGroupOf(type(\"SERVICE\"),entityId(\"%s\")))", serviceID), "", "", me.credentials)
	if err := service.Get(ctx, service.SchemaID(), &entitySettings); err == nil {
		keyRequestIDs := map[string]string{}
		for _, name := range names {
			for _, entity := range entitySettings.Entities {
				if entity.DisplayName != nil && *entity.DisplayName == name {
					keyRequestIDs[*entity.DisplayName] = *entity.EntityId
				}
			}
		}
		return keyRequestIDs
	}
	cfg := ctx.Value(settings.ContextKeyStateConfig)
	if keyRequestConfig, ok := cfg.(*keyrequests.Settings); ok && keyRequestConfig != nil {
		return keyRequestConfig.KeyRequestIDs
	}
	return nil
}

func (me *service) Delete(ctx context.Context, id string) error {
	return me.service.Delete(ctx, id)
}

func (me *service) List(ctx context.Context) (api.Stubs, error) {
	return me.service.List(ctx)
}

func (me *service) SchemaID() string {
	return me.service.SchemaID()
}
