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

package detection

import (
	"context"
	"time"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api"
	detection "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/applications/web/detection/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings"
)

const SchemaID = "v1:config:applications:detection"
const BasePath = "/api/config/v1/applicationDetectionRules"

func Service(credentials *rest.Credentials) settings.CRUDService[*detection.Rule] {
	return &service{settings.NewAPITokenService(
		credentials,
		SchemaID,
		settings.DefaultServiceOptions[*detection.Rule](BasePath),
	)}
}

type service struct {
	service settings.CRUDService[*detection.Rule]
}

func (s *service) List(ctx context.Context) (api.Stubs, error) {
	return s.service.List(ctx)
}

func (s *service) Get(ctx context.Context, id string, v *detection.Rule) error {
	return s.service.Get(ctx, id, v)
}

func (s *service) SchemaID() string {
	return s.service.SchemaID()
}

func (s *service) Create(ctx context.Context, v *detection.Rule) (*api.Stub, error) {
	stub, err := s.service.Create(ctx, v)
	if err == nil {
		return stub, err
	}
	numRetriesLeft := 6
	for err != nil && err.Error() == "Unable to persist application detection rule" {
		numRetriesLeft--
		if numRetriesLeft < 0 {
			break
		}
		time.Sleep(10 * time.Second)
		stub, err = s.service.Create(ctx, v)
	}
	return stub, err
}

func (s *service) Update(ctx context.Context, id string, v *detection.Rule) error {
	return s.service.Update(ctx, id, v)
}

func (s *service) Delete(ctx context.Context, id string) error {
	return s.service.Delete(ctx, id)
}
