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

package requestnaming

import (
	"strings"
	"time"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings/services/httpcache"

	requestnaming "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/requestnaming/settings"
)

const SchemaID = "v1:config:service:request-naming"
const BasePath = "/api/config/v1/service/requestNaming"

func Service(credentials *settings.Credentials) settings.CRUDService[*requestnaming.RequestNaming] {
	return &service{service: settings.NewCRUDService(
		credentials,
		SchemaID,
		settings.DefaultServiceOptions[*requestnaming.RequestNaming](BasePath),
	), client: httpcache.DefaultClient(credentials.URL, credentials.Token, SchemaID)}
}

type service struct {
	service settings.CRUDService[*requestnaming.RequestNaming]
	client  rest.Client
}

func (me *service) List() (api.Stubs, error) {
	return me.service.List()
}

func (me *service) Get(id string, v *requestnaming.RequestNaming) error {
	return me.service.Get(id, v)
}

func (me *service) SchemaID() string {
	return me.service.SchemaID()
}

func (me *service) Create(v *requestnaming.RequestNaming) (*api.Stub, error) {
	var err error
	var req rest.Request
	var stub api.Stub
	retries := 30

	for i := 0; i < retries; i++ {
		req = me.client.Post(BasePath, v).Expect(201)
		if err = req.Finish(&stub); err != nil {
			if !strings.Contains(err.Error(), "Unknown management zone.") {
				return nil, err
			}
		} else {
			break
		}
		time.Sleep(2 * time.Second)
	}
	if err != nil {
		return nil, err
	}

	return &stub, nil
}

func (me *service) Update(id string, v *requestnaming.RequestNaming) error {
	return me.service.Update(id, v)
}

func (me *service) Delete(id string) error {
	return me.service.Delete(id)
}

func (me *service) Name() string {
	return me.service.Name()
}
