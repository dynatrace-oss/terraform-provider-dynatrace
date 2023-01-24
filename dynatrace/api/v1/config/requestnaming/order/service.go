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

package order

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings"

	order "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/requestnaming/order/settings"
)

const SchemaID = "v1:config:service:request-naming:order"

func Service(credentials *settings.Credentials) settings.CRUDService[*order.Order] {
	return &service{client: rest.DefaultClient(credentials.URL, credentials.Token)}
}

type service struct {
	client rest.Client
}

func (me *service) List() (settings.Stubs, error) {
	return settings.Stubs{&settings.Stub{ID: "dynatrace_request_namings", Name: "dynatrace_request_namings"}}, nil
}

func (me *service) Get(id string, v *order.Order) error {
	return me.client.Get("/api/config/v1/service/requestNaming", 200).Finish(v)
}

func (me *service) SchemaID() string {
	return SchemaID
}

func (me *service) Create(v *order.Order) (*settings.Stub, error) {
	if err := me.client.Put("/api/config/v1/service/requestNaming/order", v, 204).Finish(); err != nil {
		return nil, err
	}
	return &settings.Stub{ID: "dynatrace_request_namings", Name: "dynatrace_request_namings"}, nil
}

func (me *service) Update(id string, v *order.Order) error {
	return me.client.Put("/api/config/v1/service/requestNaming/order", v, 204).Finish()
}

func (me *service) Delete(id string) error {
	return nil
}
