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

package entity

import (
	"fmt"
	"net/url"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings"

	entity "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v2/entity/settings"
)

const SchemaID = "v2:environment:entity"

func Service(credentials *settings.Credentials) settings.RService[*entity.Entity] {
	return &service{client: rest.DefaultClient(credentials.URL, credentials.Token)}
}

type service struct {
	client rest.Client
}

func (me *service) Get(id string, v *entity.Entity) error {
	return me.client.Get(fmt.Sprintf(`/api/v2/entities/%s`, url.QueryEscape(id)), 200).Finish(v)
}

func (me *service) SchemaID() string {
	return SchemaID
}

func (me *service) List() (settings.Stubs, error) {
	return settings.Stubs{&settings.Stub{ID: me.SchemaID(), Name: me.SchemaID()}}, nil
}
