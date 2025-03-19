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

package defaultversion

import (
	"context"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api"
	defaultversion "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/deployment/oneagent/defaultversion/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings"
	"github.com/google/uuid"
)

const SchemaVersion = "1.0.2"
const SchemaID = "builtin:deployment.oneagent.default-version"

func Service(credentials *rest.Credentials) settings.CRUDService[*defaultversion.Settings] {
	return &service{}
}

type service struct{}

func (me *service) List(ctx context.Context) (api.Stubs, error) {
	return api.Stubs{}, nil
}

func (me *service) Get(ctx context.Context, id string, v *defaultversion.Settings) error {
	return nil
}

func (me *service) SchemaID() string {
	return SchemaID
}

func (me *service) Create(ctx context.Context, v *defaultversion.Settings) (*api.Stub, error) {
	return &api.Stub{ID: uuid.NewString()}, nil
}

func (me *service) Update(ctx context.Context, id string, v *defaultversion.Settings) error {
	return nil
}

func (me *service) Delete(ctx context.Context, id string) error {
	return nil
}
