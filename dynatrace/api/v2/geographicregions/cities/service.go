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

package cities

import (
	"context"
	"fmt"
	"strings"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings"

	cities "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v2/geographicregions/cities/settings"
)

const SchemaID = "v2:geographicregions:cities"

func Service(credentials *rest.Credentials) settings.RService[*cities.Settings] {
	return &service{client: rest.APITokenClient(credentials)}
}

type service struct {
	client rest.Client
}

func (me *service) Get(ctx context.Context, id string, v *cities.Settings) (err error) {
	var countryCode, regionCode string
	parts := strings.Split(id, "-")
	if len(parts) == 2 {
		countryCode = parts[0]
		regionCode = parts[1]
	} else {
		return fmt.Errorf("invalid ID format: {countrycode}-{regioncode}")
	}
	return me.client.Get(ctx, fmt.Sprintf("/api/v2/rum/cities/%s/%s", countryCode, regionCode), 200).Finish(v)
}

func (me *service) List(ctx context.Context) (api.Stubs, error) {
	return api.Stubs{&api.Stub{ID: me.SchemaID(), Name: me.SchemaID()}}, nil
}

func (me *service) SchemaID() string {
	return SchemaID
}
