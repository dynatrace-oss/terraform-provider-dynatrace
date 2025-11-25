/**
* @license
* Copyright 2025 Dynatrace LLC
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

package customservices_order

import (
	"context"
	"fmt"
	"slices"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api"
	order "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/customservices/order/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings"
)

const SchemaID = "v1:config:custom-services:order"
const StaticID = "a67b38c9-1e0b-4c5f-8320-f1713df613d3"
const StaticName = "custom_service_order"
const BasePath = "/api/config/v1/service/customServices/%s"

func Service(credentials *rest.Credentials) settings.CRUDService[*order.Settings] {
	return &service{client: rest.APITokenClient(credentials)}
}

type service struct {
	client rest.Client
}

func (s *service) Create(ctx context.Context, v *order.Settings) (*api.Stub, error) {
	// Creating this resource doesn't really add a setting on the remote side
	// It just defines a specific order
	return &api.Stub{ID: StaticID, Name: StaticName}, s.Update(ctx, "", v)
}

func (s *service) Get(ctx context.Context, id string, v *order.Settings) (err error) {
	cfg := ctx.Value(settings.ContextKeyStateConfig)
	var currentConfig *order.Settings
	if orderConfig, ok := cfg.(*order.Settings); ok {
		currentConfig = orderConfig
	}
	v.IDs = map[string][]string{}
	for techID, technology := range order.Technologies {
		var ids []string
		if ids, err = s.get(ctx, technology.RESTName); err != nil {
			return err
		}
		currentIDs := []string{}
		if currentConfig != nil {
			if techIDs, found := currentConfig.IDs[techID]; found {
				currentIDs = techIDs
			}
		}
		techIDs := []string{}
		for _, id := range ids {
			if len(currentIDs) == 0 || slices.Contains(currentIDs, id) {
				techIDs = append(techIDs, id)
			}
		}
		v.IDs[techID] = techIDs
	}
	return nil
}

func (s *service) get(ctx context.Context, technology string) (ids []string, err error) {
	var listResponse struct {
		Values []struct {
			ID string `json:"id"`
		} `json:"values"`
	}
	req := s.client.Get(ctx, fmt.Sprintf(BasePath, technology), 200)
	if err = req.Finish(&listResponse); err != nil {
		return nil, err
	}
	for _, entry := range listResponse.Values {
		ids = append(ids, settings.JoinID(entry.ID, technology))
	}
	return ids, nil
}

type ListResponse struct {
	Values []struct {
		ID string `json:"id"`
	} `json:"values"`
}

func (s *service) List(ctx context.Context) (stubs api.Stubs, err error) {
	return append(stubs, &api.Stub{ID: StaticID, Name: StaticName}), nil
}

func (s *service) SchemaID() string {
	return SchemaID
}

func (s *service) Update(ctx context.Context, id string, v *order.Settings) (err error) {
	if len(v.IDs) == 0 {
		return nil
	}
	for techKey, technology := range order.Technologies {
		if ids, found := v.IDs[techKey]; found {
			if err := s.update(ctx, technology.RESTName, ids); err != nil {
				return err
			}
		}
	}
	return nil
}

func (s *service) update(ctx context.Context, technology string, ids []string) (err error) {
	if len(ids) == 0 {
		return nil
	}
	var payload = struct {
		Values []struct {
			ID string `json:"id"`
		} `json:"values"`
	}{}
	for _, id := range ids {
		uuid, _, _ := settings.SplitID(id)
		payload.Values = append(payload.Values, struct {
			ID string `json:"id"`
		}{ID: uuid})
	}
	return s.client.Put(ctx, fmt.Sprintf(BasePath, technology)+"/order", &payload, 204).Finish()
}

func (s *service) Delete(ctx context.Context, id string) error {
	// Deleting this resource does not actually destroy anything on the remote side
	// We simply forget about the order of things
	return nil
}
